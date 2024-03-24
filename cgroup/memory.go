package cgroup

type Memory struct {
	Swap *int64
	Min  *int64
	Max  *int64
	Low  *int64
	High *int64
}

func (m *Memory) Values() (v []Value) {
	if m.Swap != nil {
		v = append(v, Value{
			filename: "memory.swap.max",
			value:    *m.Swap,
		})
	}
	if m.Min != nil {
		v = append(v, Value{
			filename: "memory.min",
			value:    *m.Min,
		})
	}
	if m.Max != nil {
		v = append(v, Value{
			filename: "memory.max",
			value:    *m.Max,
		})
	}
	if m.Low != nil {
		v = append(v, Value{
			filename: "memory.low",
			value:    *m.Low,
		})
	}
	if m.High != nil {
		v = append(v, Value{
			filename: "memory.high",
			value:    *m.High,
		})
	}
	return
}
