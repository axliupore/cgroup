package cgroup

import (
	"errors"
	"fmt"
	"github.com/cilium/ebpf/asm"
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
	"math"
)

type program struct {
	insts       asm.Instructions
	hasWildCard bool
	blockID     int
}

const (
	// license string format is same as kernel MODULE_LICENSE macro
	license = "Apache"
)

func setDevices(path string, devices []specs.LinuxDeviceCgroup) error {
	if len(devices) == 0 {
		return nil
	}
	insts, license, err := DeviceFilter(devices)
	if err != nil {
		return err
	}
	dirFD, err := unix.Open(path, unix.O_DIRECTORY|unix.O_RDONLY|unix.O_CLOEXEC, 0o600)
	if err != nil {
		return fmt.Errorf("cannot get dir FD for %s", path)
	}
	defer unix.Close(dirFD)
	if _, err := LoadAttachCgroupDeviceFilter(insts, license, dirFD); err != nil {
		if !canSkipEBPFError(devices) {
			return err
		}
	}
	return nil
}

func DeviceFilter(devices []specs.LinuxDeviceCgroup) (asm.Instructions, string, error) {
	p := &program{}
	p.init()
	for i := len(devices) - 1; i >= 0; i-- {
		if err := p.appendDevice(devices[i]); err != nil {
			return nil, "", err
		}
	}
	insts, err := p.finalize()
	return insts, license, err
}

func (p *program) init() {
	// struct bpf_cgroup_dev_ctx: https://elixir.bootlin.com/linux/v5.3.6/source/include/uapi/linux/bpf.h#L3423
	/*
		u32 access_type
		u32 major
		u32 minor
	*/
	// R2 <- type (lower 16 bit of u32 access_type at R1[0])
	p.insts = append(p.insts,
		asm.LoadMem(asm.R2, asm.R1, 0, asm.Half))

	// R3 <- access (upper 16 bit of u32 access_type at R1[0])
	p.insts = append(p.insts,
		asm.LoadMem(asm.R3, asm.R1, 0, asm.Word),
		// RSh: bitwise shift right
		asm.RSh.Imm32(asm.R3, 16))

	// R4 <- major (u32 major at R1[4])
	p.insts = append(p.insts,
		asm.LoadMem(asm.R4, asm.R1, 4, asm.Word))

	// R5 <- minor (u32 minor at R1[8])
	p.insts = append(p.insts,
		asm.LoadMem(asm.R5, asm.R1, 8, asm.Word))
}

// appendDevice needs to be called from the last element of OCI linux.resources.devices to the head element.
func (p *program) appendDevice(dev specs.LinuxDeviceCgroup) error {
	if p.blockID < 0 {
		return errors.New("the program is finalized")
	}
	if p.hasWildCard {
		// All entries after wildcard entry are ignored
		return nil
	}

	bpfType := int32(-1)
	hasType := true
	switch dev.Type {
	case string('c'):
		bpfType = int32(unix.BPF_DEVCG_DEV_CHAR)
	case string('b'):
		bpfType = int32(unix.BPF_DEVCG_DEV_BLOCK)
	case string('a'):
		hasType = false
	default:
		// if not specified in OCI json, typ is set to DeviceTypeAll
		return fmt.Errorf("invalid DeviceType %q", dev.Type)
	}
	if *dev.Major > math.MaxUint32 {
		return fmt.Errorf("invalid major %d", *dev.Major)
	}
	if *dev.Minor > math.MaxUint32 {
		return fmt.Errorf("invalid minor %d", *dev.Major)
	}
	hasMajor := *dev.Major >= 0 // if not specified in OCI json, major is set to -1
	hasMinor := *dev.Minor >= 0
	bpfAccess := int32(0)
	for _, r := range dev.Access {
		switch r {
		case 'r':
			bpfAccess |= unix.BPF_DEVCG_ACC_READ
		case 'w':
			bpfAccess |= unix.BPF_DEVCG_ACC_WRITE
		case 'm':
			bpfAccess |= unix.BPF_DEVCG_ACC_MKNOD
		default:
			return fmt.Errorf("unknown device access %v", r)
		}
	}
	// If the access is rwm, skip the check.
	hasAccess := bpfAccess != (unix.BPF_DEVCG_ACC_READ | unix.BPF_DEVCG_ACC_WRITE | unix.BPF_DEVCG_ACC_MKNOD)

	blockSym := fmt.Sprintf("block-%d", p.blockID)
	nextBlockSym := fmt.Sprintf("block-%d", p.blockID+1)
	prevBlockLastIdx := len(p.insts) - 1
	if hasType {
		p.insts = append(p.insts,
			// if (R2 != bpfType) goto next
			asm.JNE.Imm(asm.R2, bpfType, nextBlockSym),
		)
	}
	if hasAccess {
		p.insts = append(p.insts,
			// if (R3 & bpfAccess == 0 /* use R1 as a temp var */) goto next
			asm.Mov.Reg32(asm.R1, asm.R3),
			asm.And.Imm32(asm.R1, bpfAccess),
			asm.JEq.Imm(asm.R1, 0, nextBlockSym),
		)
	}
	if hasMajor {
		p.insts = append(p.insts,
			// if (R4 != major) goto next
			asm.JNE.Imm(asm.R4, int32(*dev.Major), nextBlockSym),
		)
	}
	if hasMinor {
		p.insts = append(p.insts,
			// if (R5 != minor) goto next
			asm.JNE.Imm(asm.R5, int32(*dev.Minor), nextBlockSym),
		)
	}
	if !hasType && !hasAccess && !hasMajor && !hasMinor {
		p.hasWildCard = true
	}
	p.insts = append(p.insts, acceptBlock(dev.Allow)...)
	// set blockSym to the first instruction we added in this iteration
	p.insts[prevBlockLastIdx+1] = p.insts[prevBlockLastIdx+1].WithSymbol(blockSym)
	p.blockID++
	return nil
}

func (p *program) finalize() (asm.Instructions, error) {
	if p.hasWildCard {
		// acceptBlock with asm.Return() is already inserted
		return p.insts, nil
	}
	blockSym := fmt.Sprintf("block-%d", p.blockID)
	p.insts = append(p.insts,
		// R0 <- 0
		asm.Mov.Imm32(asm.R0, 0).WithSymbol(blockSym),
		asm.Return(),
	)
	p.blockID = -1
	return p.insts, nil
}

func acceptBlock(accept bool) asm.Instructions {
	v := int32(0)
	if accept {
		v = 1
	}
	return []asm.Instruction{
		// R0 <- v
		asm.Mov.Imm32(asm.R0, v),
		asm.Return(),
	}
}
