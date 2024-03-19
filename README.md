# cgroup

Go package for creating, managing, inspecting, and destorying cgrouops.My project documentation and code refer to this [project](https://github.com/containerd/cgroups).

## Introduction

This project uses cgroup v2 version, why don’t I use v1? Because I think v2 is more common and convenient.

| Distro                           | Version          |
| -------------------------------- | ---------------- |
| Fedora                           | since 31         |
| Arch Linux                       | since April 2021 |
| openSUSE Tumbleweed              | since c. 2021    |
| Debian GNU/Linux                 | since 11         |
| Ubuntu                           | since 21.10      |
| RHEL and RHEL-like distributions | since 9          |

You can enter the following command `stat -fc %T /sys/fs/cgroup` to check your version number. If the output is `cgroup2fs`, it is v2. Otherwise it is v1.
