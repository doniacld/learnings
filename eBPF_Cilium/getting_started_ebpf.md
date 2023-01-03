# Getting started with eBPF

This page is a recap of my notes while doing the [Getting started with eBPF](https://isovalent.com/labs/getting-started-with-ebpf/) lab from Isovalent.

## 🐝 eBPF

> BPF is a revolutionary technology with origins in the Linux kernel that can run sandboxed programs in an operating system kernel. It is used to safely and efficiently extend the capabilities of the kernel 
> without requiring to change kernel source code or load kernel modules.

What you will do in this course
To get first hand experiences with eBPF, in this lab we will:

- Build and use `opensnoop`, an eBPF based tool that reports whenever a file is opened
- Use `readelf` to compare a BPF object file with its source code
- Use `bpftool` to see how your tool is loaded into the kernel
- Add additional “hello world”-style tracing to the source code
- Build and run it again to see your own customized eBPF tracing 

The example tool we're using, `opensnoop`, is one of a collection of eBPF-based tools from the BCC project

## 🕵️ opensnoop

In `~/bcc/libbpf-tools` there are all the available bcc source code

```shell
root@ubuntu-2110:~# cd ~/bcc/libbpf-tools
root@ubuntu-2110:~/bcc/libbpf-tools# ls
Makefile            biostacks.h          execsnoop.bpf.c    gethostlatency.bpf.c  mountsnoop.h      runqlen.h               tcpconnect.c
README.md           bitesize.bpf.c       execsnoop.c        gethostlatency.c      numamove.bpf.c    runqslower.bpf.c        tcpconnect.h
arm64               bitesize.c           execsnoop.h        gethostlatency.h      numamove.c        runqslower.c            tcpconnlat.bpf.c
bashreadline.bpf.c  bitesize.h           exitsnoop.bpf.c    hardirqs.bpf.c        offcputime.bpf.c  runqslower.h            tcpconnlat.c
bashreadline.c      bits.bpf.h           exitsnoop.c        hardirqs.c            offcputime.c      runqslower_example.txt  tcpconnlat.h
bashreadline.h      blk_types.h          exitsnoop.h        hardirqs.h            offcputime.h      softirqs.bpf.c          tcprtt.bpf.c
bin                 cachestat.bpf.c      filelife.bpf.c     kernel.config         oomkill.bpf.c     softirqs.c              tcprtt.c
bindsnoop.bpf.c     cachestat.c          filelife.c         klockstat.bpf.c       oomkill.c         softirqs.h              tcprtt.h
bindsnoop.c         core_fixes.bpf.h     filelife.h         klockstat.c           oomkill.h         solisten.bpf.c          tcpsynbl.bpf.c
bindsnoop.h         cpudist.bpf.c        filetop.bpf.c      klockstat.h           opensnoop.bpf.c   solisten.c              tcpsynbl.c
biolatency.bpf.c    cpudist.c            filetop.c          ksnoop.bpf.c          opensnoop.c       solisten.h              tcpsynbl.h
biolatency.c        cpudist.h            filetop.h          ksnoop.c              opensnoop.h       stat.h                  trace_helpers.c
biolatency.h        cpufreq.bpf.c        fsdist.bpf.c       ksnoop.h              powerpc           statsnoop.bpf.c         trace_helpers.h
biopattern.bpf.c    cpufreq.c            fsdist.c           llcstat.bpf.c         readahead.bpf.c   statsnoop.c             uprobe_helpers.c
biopattern.c        cpufreq.h            fsdist.h           llcstat.c             readahead.c       statsnoop.h             uprobe_helpers.h
biopattern.h        drsnoop.bpf.c        fsslower.bpf.c     llcstat.h             readahead.h       syscall_helpers.c       vfsstat.bpf.c
biosnoop.bpf.c      drsnoop.c            fsslower.c         map_helpers.c         runqlat.bpf.c     syscall_helpers.h       vfsstat.c
biosnoop.c          drsnoop.h            fsslower.h         map_helpers.h         runqlat.c         syscount.bpf.c          vfsstat.h
biosnoop.h          drsnoop_example.txt  funclatency.bpf.c  maps.bpf.h            runqlat.h         syscount.c              x86
biostacks.bpf.c     errno_helpers.c      funclatency.c      mountsnoop.bpf.c      runqlen.bpf.c     syscount.h
biostacks.c         errno_helpers.h      funclatency.h      mountsnoop.c          runqlen.c         tcpconnect.bpf.c
root@ubuntu-2110:~/bcc/libbpf-tools# 
```

Build one of the program with make:
```shell
make opensnoop
```

The output of the build for `opensnoop`:
```shell
root@ubuntu-2110:~/bcc/libbpf-tools# make opensnoop
  MKDIR    libbpf
  LIB      libbpf.a
make[1]: pkg-config: No such file or directory
  MKDIR    /root/bcc/libbpf-tools/.output/libbpf/staticobjs
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/bpf.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/btf.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/libbpf.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/libbpf_errno.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/netlink.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/nlattr.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/str_error.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/libbpf_probes.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/bpf_prog_linfo.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/xsk.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/btf_dump.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/hashmap.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/ringbuf.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/strset.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/linker.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/gen_loader.o
  CC       /root/bcc/libbpf-tools/.output/libbpf/staticobjs/relo_core.o
  AR       /root/bcc/libbpf-tools/.output/libbpf/libbpf.a
  INSTALL  bpf.h libbpf.h btf.h libbpf_common.h libbpf_legacy.h xsk.h bpf_helpers.h bpf_helper_defs.h bpf_tracing.h bpf_endian.h bpf_core_read.h skel_internal.h libbpf_version.h
  INSTALL  /root/bcc/libbpf-tools/.output/libbpf/libbpf.pc
  INSTALL  /root/bcc/libbpf-tools/.output/libbpf/libbpf.a 
  MKDIR    .output
  BPF      opensnoop.bpf.o
  GEN-SKEL opensnoop.skel.h
  CC       opensnoop.o
  CC       trace_helpers.o
  CC       syscall_helpers.o
  CC       errno_helpers.o
  CC       map_helpers.o
  CC       uprobe_helpers.o
  BINARY   opensnoop
```

We can see at the end of the output that the binary `opensnoop` is available.

Now we can run opensnoop. This needs CAP_BPF privileges — but since in this demo environment we are running as root anyway, this is not a problem:

```shell
./opensnoop
```

In terminal 2, we do some operations:
```shell
root@ubuntu-2110:~# cat /etc/os-release
PRETTY_NAME="Ubuntu 21.10"
NAME="Ubuntu"
VERSION_ID="21.10"
VERSION="21.10 (Impish Indri)"
VERSION_CODENAME=impish
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=impish
root@ubuntu-2110:~# ls -lart /etc/os-release
lrwxrwxrwx 1 root root 21 Oct 11  2021 /etc/os-release -> ../usr/lib/os-release
root@ubuntu-2110:~# chmod 777 /etc/os-release
root@ubuntu-2110:~# ls -lart /etc/os-release
lrwxrwxrwx 1 root root 21 Oct 11  2021 /etc/os-release -> ../usr/lib/os-release
root@ubuntu-2110:~# 
```

While in terminal 1, we run `opensnoop`:
```shell
root@ubuntu-2110:~/bcc/libbpf-tools# ./opensnoop 
PID    COMM              FD ERR PATH
1      systemd           22   0 /proc/1435/cgroup
1      systemd           22   0 /proc/197/cgroup
1      systemd           22   0 /proc/442/cgroup
2882   opensnoop         17   0 /etc/localtime
519    google_osconfig    8   0 /etc/hosts
1435   systemd-resolve   17   0 /run/systemd/netif/links/2
1435   systemd-resolve   18   0 /run/systemd/netif/links/2
626    google_guest_ag    8   0 /etc/hosts
1435   systemd-resolve   17   0 /run/systemd/netif/links/2
1435   systemd-resolve   18   0 /run/systemd/netif/links/2
2883   cat                3   0 /etc/ld.so.cache
2883   cat                3   0 /lib/x86_64-linux-gnu/libc.so.6
2883   cat                3   0 /usr/lib/locale/locale-archive
2883   cat                3   0 /usr/share/locale/locale.alias
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_IDENTIFICATION
2883   cat                3   0 /usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_MEASUREMENT
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_TELEPHONE
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_ADDRESS
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_NAME
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_PAPER
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES/SYS_LC_MESSAGES
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_MONETARY
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_COLLATE
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_TIME
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_NUMERIC
2883   cat                3   0 /usr/lib/locale/C.UTF-8/LC_CTYPE
2883   cat                3   0 /etc/os-release
2732   bash               3   0 /root/.bash_history
2884   ls                 3   0 /etc/ld.so.cache
2884   ls                 3   0 /lib/x86_64-linux-gnu/libselinux.so.1
2884   ls                 3   0 /lib/x86_64-linux-gnu/libc.so.6
2884   ls                 3   0 /lib/x86_64-linux-gnu/libpcre2-8.so.0
2884   ls                 3   0 /proc/filesystems
2884   ls                 3   0 /usr/lib/locale/locale-archive
2884   ls                 3   0 /usr/share/locale/locale.alias
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_IDENTIFICATION
2884   ls                 3   0 /usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MEASUREMENT
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_TELEPHONE
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_ADDRESS
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_NAME
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_PAPER
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES/SYS_LC_MESSAGES
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MONETARY
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_COLLATE
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_TIME
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_NUMERIC
2884   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_CTYPE
2884   ls                -1   2 /usr/share/locale/C.UTF-8/LC_TIME/coreutils.mo
2884   ls                -1   2 /usr/share/locale/C.utf8/LC_TIME/coreutils.mo
2884   ls                -1   2 /usr/share/locale/C/LC_TIME/coreutils.mo
2884   ls                 3   0 /etc/nsswitch.conf
2884   ls                 3   0 /etc/passwd
2884   ls                 3   0 /etc/group
2884   ls                 3   0 /etc/localtime
2732   bash               3   0 /root/.bash_history
1      systemd           22   0 /proc/849/cgroup
2885   chmod              3   0 /etc/ld.so.cache
2885   chmod              3   0 /lib/x86_64-linux-gnu/libc.so.6
2885   chmod              3   0 /usr/lib/locale/locale-archive
2885   chmod              3   0 /usr/share/locale/locale.alias
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_IDENTIFICATION
2885   chmod              3   0 /usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_MEASUREMENT
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_TELEPHONE
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_ADDRESS
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_NAME
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_PAPER
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES/SYS_LC_MESSAGES
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_MONETARY
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_COLLATE
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_TIME
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_NUMERIC
2885   chmod              3   0 /usr/lib/locale/C.UTF-8/LC_CTYPE
2732   bash               3   0 /root/.bash_history
2886   ls                 3   0 /etc/ld.so.cache
2886   ls                 3   0 /lib/x86_64-linux-gnu/libselinux.so.1
2886   ls                 3   0 /lib/x86_64-linux-gnu/libc.so.6
2886   ls                 3   0 /lib/x86_64-linux-gnu/libpcre2-8.so.0
2886   ls                 3   0 /proc/filesystems
2886   ls                 3   0 /usr/lib/locale/locale-archive
2886   ls                 3   0 /usr/share/locale/locale.alias
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_IDENTIFICATION
2886   ls                 3   0 /usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MEASUREMENT
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_TELEPHONE
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_ADDRESS
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_NAME
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_PAPER
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MESSAGES/SYS_LC_MESSAGES
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_MONETARY
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_COLLATE
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_TIME
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_NUMERIC
2886   ls                 3   0 /usr/lib/locale/C.UTF-8/LC_CTYPE
2886   ls                -1   2 /usr/share/locale/C.UTF-8/LC_TIME/coreutils.mo
2886   ls                -1   2 /usr/share/locale/C.utf8/LC_TIME/coreutils.mo
2886   ls                -1   2 /usr/share/locale/C/LC_TIME/coreutils.mo
2886   ls                 3   0 /etc/nsswitch.conf
2886   ls                 3   0 /etc/passwd
2886   ls                 3   0 /etc/group
2886   ls                 3   0 /etc/localtime
2732   bash               3   0 /root/.bash_history
1824   snapd              4   0 /var/lib/snapd/assertions/asserts-v0/model/16/canonical/gcp-classic
1824   snapd             12   0 /var/lib/snapd/assertions/asserts-v0/model/16/canonical/gcp-classic/active
1824   snapd              4   0 /var/lib/snapd/assertions/asserts-v0/serial/canonical/gcp-classic/bd2e3bb7-99c4-46b8-8d63-e6743f19ee16
1824   snapd             12   0 /var/lib/snapd/assertions/asserts-v0/serial/canonical/gcp-classic/bd2e3bb7-99c4-46b8-8d63-e6743f19ee16/active
1824   snapd              4   0 /var/lib/snapd/assertions/asserts-v0/model/16/canonical/gcp-classic
1824   snapd             12   0 /var/lib/snapd/assertions/asserts-v0/model/16/canonical/gcp-classic/active
1824   snapd              4   0 /var/lib/snapd/assertions/asserts-v0/serial/canonical/gcp-classic/bd2e3bb7-99c4-46b8-8d63-e6743f19ee16
1824   snapd             12   0 /var/lib/snapd/assertions/asserts-v0/serial/canonical/gcp-classic/bd2e3bb7-99c4-46b8-8d63-e6743f19ee16/active
1824   snapd              4   0 /var/lib/snapd/assertions/asserts-v0/model/16/canonical/gcp-classic
1824   snapd             12   0 /var/lib/snapd/assertions/asserts-v0/model/16/canonical/gcp-classic/active
1      systemd           22   0 /proc/1824/cgroup
```

_NB: If you leave opensnoop running, from time to time you might see output generated by other processes running on the VM, such as systemd._

## 🔍 Deeper look into an eBPF program

Use the command `readelf` to examine the BPF program file:

```shell
root@ubuntu-2110:~/bcc/libbpf-tools# readelf --section-details --headers .output/opensnoop.bpf.o
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00 
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              REL (Relocatable file)
  Machine:                           Linux BPF
  Version:                           0x1
  Entry point address:               0x0
  Start of program headers:          0 (bytes into file)
  Start of section headers:          11304 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           0 (bytes)
  Number of program headers:         0
  Size of section headers:           64 (bytes)
  Number of section headers:         20
  Section header string table index: 19

Section Headers:
  [Nr] Name
       Type              Address          Offset            Link
       Size              EntSize          Info              Align
       Flags
  [ 0] 
       NULL             0000000000000000  0000000000000000  0
       0000000000000000 0000000000000000  0                 0
       [0000000000000000]: 
  [ 1] .text
       PROGBITS         0000000000000000  0000000000000040  0
       0000000000000000 0000000000000000  0                 4
       [0000000000000006]: ALLOC, EXEC
  [ 2] tracepoint/syscalls/sys_enter_open
       PROGBITS         0000000000000000  0000000000000040  0
       0000000000000178 0000000000000000  0                 8
       [0000000000000006]: ALLOC, EXEC
  [ 3] tracepoint/syscalls/sys_enter_openat
       PROGBITS         0000000000000000  00000000000001b8  0
       0000000000000178 0000000000000000  0                 8
       [0000000000000006]: ALLOC, EXEC
  [ 4] tracepoint/syscalls/sys_exit_open
       PROGBITS         0000000000000000  0000000000000330  0
       00000000000002d0 0000000000000000  0                 8
       [0000000000000006]: ALLOC, EXEC
  [ 5] tracepoint/syscalls/sys_exit_openat
       PROGBITS         0000000000000000  0000000000000600  0
       00000000000002d0 0000000000000000  0                 8
       [0000000000000006]: ALLOC, EXEC
  [ 6] .rodata
       PROGBITS         0000000000000000  00000000000008d0  0
       000000000000000d 0000000000000000  0                 4
       [0000000000000002]: ALLOC
  [ 7] .maps
       PROGBITS         0000000000000000  00000000000008e0  0
       0000000000000038 0000000000000000  0                 8
       [0000000000000003]: WRITE, ALLOC
  [ 8] license
       PROGBITS         0000000000000000  0000000000000918  0
       0000000000000004 0000000000000000  0                 1
       [0000000000000003]: WRITE, ALLOC
  [ 9] .BTF
       PROGBITS         0000000000000000  000000000000091c  0
       0000000000000d2b 0000000000000000  0                 1
       [0000000000000000]: 
  [10] .BTF.ext
       PROGBITS         0000000000000000  0000000000001647  0
       00000000000007ec 0000000000000000  0                 1
       [0000000000000000]: 
  [11] .symtab
       SYMTAB           0000000000000000  0000000000001e38  19
       00000000000002d0 0000000000000018  19                8
       [0000000000000000]: 
  [12] .reltracepoint/syscalls/sys_enter_open
       REL              0000000000000000  0000000000002108  11
       0000000000000040 0000000000000010  2                 8
       [0000000000000000]: 
  [13] .reltracepoint/syscalls/sys_enter_openat
       REL              0000000000000000  0000000000002148  11
       0000000000000040 0000000000000010  3                 8
       [0000000000000000]: 
  [14] .reltracepoint/syscalls/sys_exit_open
       REL              0000000000000000  0000000000002188  11
       0000000000000040 0000000000000010  4                 8
       [0000000000000000]: 
  [15] .reltracepoint/syscalls/sys_exit_openat
       REL              0000000000000000  00000000000021c8  11
       0000000000000040 0000000000000010  5                 8
       [0000000000000000]: 
  [16] .rel.BTF
       REL              0000000000000000  0000000000002208  11
       0000000000000070 0000000000000010  9                 8
       [0000000000000000]: 
  [17] .rel.BTF.ext
       REL              0000000000000000  0000000000002278  11
       0000000000000780 0000000000000010  10                8
       [0000000000000000]: 
  [18] .llvm_addrsig
       LOOS+0xfff4c03   0000000000000000  00000000000029f8  0
       000000000000000b 0000000000000000  0                 1
       [0000000080000000]: EXCLUDE
  [19] .strtab
       STRTAB           0000000000000000  0000000000002a03  0
       0000000000000224 0000000000000000  0                 1
       [0000000000000000]: 

There are no program headers in this file.
```

There are 4 BPF programs we can see on the line `.text` above.
We find them in the source code as `tracepoint__syscalls__`.

Source code from file `opensnoop.bpf.c`:
```cgo
// SPDX-License-Identifier: GPL-2.0
// Copyright (c) 2019 Facebook
// Copyright (c) 2020 Netflix
#include <vmlinux.h>
#include <bpf/bpf_helpers.h>
#include "opensnoop.h"

const volatile pid_t targ_pid = 0;
const volatile pid_t targ_tgid = 0;
const volatile uid_t targ_uid = 0;
const volatile bool targ_failed = false;

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 10240);
	__type(key, u32);
	__type(value, struct args_t);
} start SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
	__uint(key_size, sizeof(u32));
	__uint(value_size, sizeof(u32));
} events SEC(".maps");

static __always_inline bool valid_uid(uid_t uid) {
	return uid != INVALID_UID;
}

static __always_inline
bool trace_allowed(u32 tgid, u32 pid)
{
	u32 uid;

	/* filters */
	if (targ_tgid && targ_tgid != tgid)
		return false;
	if (targ_pid && targ_pid != pid)
		return false;
	if (valid_uid(targ_uid)) {
		uid = (u32)bpf_get_current_uid_gid();
		if (targ_uid != uid) {
			return false;
		}
	}
	return true;
}

SEC("tracepoint/syscalls/sys_enter_open")
int tracepoint__syscalls__sys_enter_open(struct trace_event_raw_sys_enter* ctx)
{
	u64 id = bpf_get_current_pid_tgid();
	/* use kernel terminology here for tgid/pid: */
	u32 tgid = id >> 32;
	u32 pid = id;

	/* store arg info for later lookup */
	if (trace_allowed(tgid, pid)) {
		struct args_t args = {};
		args.fname = (const char *)ctx->args[0];
		args.flags = (int)ctx->args[1];
		bpf_map_update_elem(&start, &pid, &args, 0);
	}
	return 0;
}

SEC("tracepoint/syscalls/sys_enter_openat")
int tracepoint__syscalls__sys_enter_openat(struct trace_event_raw_sys_enter* ctx)
{
	u64 id = bpf_get_current_pid_tgid();
	/* use kernel terminology here for tgid/pid: */
	u32 tgid = id >> 32;
	u32 pid = id;

	/* store arg info for later lookup */
	if (trace_allowed(tgid, pid)) {
		struct args_t args = {};
		args.fname = (const char *)ctx->args[1];
		args.flags = (int)ctx->args[2];
		bpf_map_update_elem(&start, &pid, &args, 0);
	}
	return 0;
}

static __always_inline
int trace_exit(struct trace_event_raw_sys_exit* ctx)
{
	struct event event = {};
	struct args_t *ap;
	int ret;
	u32 pid = bpf_get_current_pid_tgid();

	ap = bpf_map_lookup_elem(&start, &pid);
	if (!ap)
		return 0;	/* missed entry */
	ret = ctx->ret;
	if (targ_failed && ret >= 0)
		goto cleanup;	/* want failed only */

	/* event data */
	event.pid = bpf_get_current_pid_tgid() >> 32;
	event.uid = bpf_get_current_uid_gid();
	bpf_get_current_comm(&event.comm, sizeof(event.comm));
	bpf_probe_read_user_str(&event.fname, sizeof(event.fname), ap->fname);
	event.flags = ap->flags;
	event.ret = ret;

	/* emit event */
	bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU,
			      &event, sizeof(event));

cleanup:
	bpf_map_delete_elem(&start, &pid);
	return 0;
}

SEC("tracepoint/syscalls/sys_exit_open")
int tracepoint__syscalls__sys_exit_open(struct trace_event_raw_sys_exit* ctx)
{
	return trace_exit(ctx);
}

SEC("tracepoint/syscalls/sys_exit_openat")
int tracepoint__syscalls__sys_exit_openat(struct trace_event_raw_sys_exit* ctx)
{
	return trace_exit(ctx);
}

char LICENSE[] SEC("license") = "GPL";
```

## 🌰 BPF programs in the Kernel

We are going to use bpftool to see what we have loaded into the Kernel.

```shell
root@ubuntu-2110:~/bcc/libbpf-tools# bpftool prog list
[...]
219: tracepoint  name tracepoint__sys  tag 9f196d70d0c1964b  gpl
        loaded_at 2023-01-03T09:37:12+0000  uid 0
        xlated 248B  jited 140B  memlock 4096B  map_ids 5,2
        btf_id 34
221: tracepoint  name tracepoint__sys  tag 47b06acd3f9a5527  gpl
        loaded_at 2023-01-03T09:37:12+0000  uid 0
        xlated 248B  jited 140B  memlock 4096B  map_ids 5,2
        btf_id 34
222: tracepoint  name tracepoint__sys  tag 387291c2fb839ac6  gpl
        loaded_at 2023-01-03T09:37:12+0000  uid 0
        xlated 696B  jited 475B  memlock 4096B  map_ids 2,5,3
        btf_id 34
223: tracepoint  name tracepoint__sys  tag 387291c2fb839ac6  gpl
        loaded_at 2023-01-03T09:37:12+0000  uid 0
        xlated 696B  jited 475B  memlock 4096B  map_ids 2,5,3
        btf_id 34
```

```cgo
root@ubuntu-2110:~#  bpftool map list
2: hash  name start  flags 0x0
        key 4B  value 16B  max_entries 10240  memlock 245760B
        btf_id 34
3: perf_event_array  name events  flags 0x0
        key 4B  value 4B  max_entries 1  memlock 4096B
5: array  name opensnoo.rodata  flags 0x480
        key 4B  value 13B  max_entries 1  memlock 4096B
        btf_id 34  frozen
```

map_ids corresponding 2,3,5

At the start of each line, you see the ID of the corresponding BPF program. Take an ID of a tracepoint program, and dump the bytecode (in our run the number was 48, your number might be different):
```shell
bpftool prog dump xlated id 223 linum
```

Output of the dump command:
```cgo
root@ubuntu-2110:~# bpftool prog dump xlated id 223 linum
int tracepoint__syscalls__sys_exit_openat(struct trace_event_raw_sys_exit * ctx):
; int tracepoint__syscalls__sys_exit_openat(struct trace_event_raw_sys_exit* ctx) [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:124 line_col:0]
   0: (bf) r6 = r1
   1: (b7) r1 = 0
; struct event event = {}; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:88 line_col:15]
   2: (7b) *(u64 *)(r10 -8) = r1
   3: (7b) *(u64 *)(r10 -16) = r1
   4: (7b) *(u64 *)(r10 -24) = r1
   5: (7b) *(u64 *)(r10 -32) = r1
   6: (7b) *(u64 *)(r10 -40) = r1
   7: (7b) *(u64 *)(r10 -48) = r1
   8: (7b) *(u64 *)(r10 -56) = r1
   9: (7b) *(u64 *)(r10 -64) = r1
  10: (7b) *(u64 *)(r10 -72) = r1
  11: (7b) *(u64 *)(r10 -80) = r1
  12: (7b) *(u64 *)(r10 -88) = r1
  13: (7b) *(u64 *)(r10 -96) = r1
  14: (7b) *(u64 *)(r10 -104) = r1
  15: (7b) *(u64 *)(r10 -112) = r1
  16: (7b) *(u64 *)(r10 -120) = r1
  17: (7b) *(u64 *)(r10 -128) = r1
  18: (7b) *(u64 *)(r10 -136) = r1
  19: (7b) *(u64 *)(r10 -144) = r1
  20: (7b) *(u64 *)(r10 -152) = r1
  21: (7b) *(u64 *)(r10 -160) = r1
  22: (7b) *(u64 *)(r10 -168) = r1
  23: (7b) *(u64 *)(r10 -176) = r1
  24: (7b) *(u64 *)(r10 -184) = r1
  25: (7b) *(u64 *)(r10 -192) = r1
  26: (7b) *(u64 *)(r10 -200) = r1
  27: (7b) *(u64 *)(r10 -208) = r1
  28: (7b) *(u64 *)(r10 -216) = r1
  29: (7b) *(u64 *)(r10 -224) = r1
  30: (7b) *(u64 *)(r10 -232) = r1
  31: (7b) *(u64 *)(r10 -240) = r1
  32: (7b) *(u64 *)(r10 -248) = r1
  33: (7b) *(u64 *)(r10 -256) = r1
  34: (7b) *(u64 *)(r10 -264) = r1
  35: (7b) *(u64 *)(r10 -272) = r1
  36: (7b) *(u64 *)(r10 -280) = r1
  37: (7b) *(u64 *)(r10 -288) = r1
  38: (7b) *(u64 *)(r10 -296) = r1
; u32 pid = bpf_get_current_pid_tgid(); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:91 line_col:12]
  39: (85) call bpf_get_current_pid_tgid#139360
; u32 pid = bpf_get_current_pid_tgid(); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:91 line_col:6]
  40: (63) *(u32 *)(r10 -300) = r0
  41: (bf) r2 = r10
;  [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:0 line_col:0]
  42: (07) r2 += -300
; ap = bpf_map_lookup_elem(&start, &pid); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:93 line_col:7]
  43: (18) r1 = map[id:2]
  45: (85) call __htab_map_lookup_elem#154288
  46: (15) if r0 == 0x0 goto pc+1
  47: (07) r0 += 56
  48: (bf) r7 = r0
; if (!ap) [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:94 line_col:6]
  49: (15) if r7 == 0x0 goto pc+35
; ret = ctx->ret; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:96 line_col:13]
  50: (79) r8 = *(u64 *)(r6 +16)
; if (targ_failed && ret >= 0) [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:97 line_col:6]
  51: (18) r1 = map[id:5][0]+12
  53: (71) r1 = *(u8 *)(r1 +0)
; event.pid = bpf_get_current_pid_tgid() >> 32; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:101 line_col:14]
  54: (85) call bpf_get_current_pid_tgid#139360
; event.pid = bpf_get_current_pid_tgid() >> 32; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:101 line_col:41]
  55: (77) r0 >>= 32
; event.pid = bpf_get_current_pid_tgid() >> 32; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:101 line_col:12]
  56: (63) *(u32 *)(r10 -288) = r0
; event.uid = bpf_get_current_uid_gid(); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:102 line_col:14]
  57: (85) call bpf_get_current_uid_gid#139680
; event.uid = bpf_get_current_uid_gid(); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:102 line_col:12]
  58: (63) *(u32 *)(r10 -284) = r0
; bpf_get_current_comm(&event.comm, sizeof(event.comm)); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:103 line_col:23]
  59: (bf) r1 = r10
  60: (07) r1 += -272
; bpf_get_current_comm(&event.comm, sizeof(event.comm)); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:103 line_col:2]
  61: (b7) r2 = 16
  62: (85) call bpf_get_current_comm#139776
; bpf_probe_read_user_str(&event.fname, sizeof(event.fname), ap->fname); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:104 line_col:65]
  63: (79) r3 = *(u64 *)(r7 +0)
; bpf_probe_read_user_str(&event.fname, sizeof(event.fname), ap->fname); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:104 line_col:26]
  64: (bf) r1 = r10
  65: (07) r1 += -256
; bpf_probe_read_user_str(&event.fname, sizeof(event.fname), ap->fname); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:104 line_col:2]
  66: (b7) r2 = 255
  67: (85) call bpf_probe_read_user_str#-64480
; event.flags = ap->flags; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:105 line_col:20]
  68: (61) r1 = *(u32 *)(r7 +8)
; event.flags = ap->flags; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:105 line_col:14]
  69: (63) *(u32 *)(r10 -276) = r1
; event.ret = ret; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:106 line_col:12]
  70: (63) *(u32 *)(r10 -280) = r8
  71: (bf) r4 = r10
; event.pid = bpf_get_current_pid_tgid() >> 32; [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:101 line_col:14]
  72: (07) r4 += -296
; bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:109 line_col:2]
  73: (bf) r1 = r6
  74: (18) r2 = map[id:3]
  76: (18) r3 = 0xffffffff
  78: (b7) r5 = 296
  79: (85) call bpf_perf_event_output_tp#-63120
  80: (bf) r2 = r10
;  [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:0 line_col:0]
  81: (07) r2 += -300
; bpf_map_delete_elem(&start, &pid); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:113 line_col:2]
  82: (18) r1 = map[id:2]
  84: (85) call htab_map_delete_elem#159440
; return trace_exit(ctx); [file:/root/bcc/libbpf-tools/opensnoop.bpf.c line_num:126 line_col:2]
  85: (b7) r0 = 0
  86: (95) exit
```

## ✍🏼 Write our own code

Now we know how to run eBPF tools, how to observe their behaviour, check what is loaded in the Kernel, even get information about what is actually happening compared to the actual source code.
Next, we want to write actual code! To do so, we will add our own tracing message right into the code!

Add in the source code from `opensnoop.bpf.c`:
```cgo
	/* emit event */
	bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU,
			      &event, sizeof(event));

	/* add our own trace */
	bpf_printk("Hello world");

cleanup:
	bpf_map_delete_elem(&start, &pid);
	return 0;
}
```

After rebuilding `opensnoop` program.

Let's follow the traces on one side while opening some files:
```shell
root@ubuntu-2110:~# cat /sys/kernel/debug/tracing/trace_pipe
       opensnoop-3178    [000] d...  1340.238800: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.761487: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.761987: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762007: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762068: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762118: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762132: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762146: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762159: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762173: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762187: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762199: bpf_trace_printk: Hello world
   sandbox-agent-1494    [000] d...  1355.762212: bpf_trace_printk: Hello world
```

On the otherside running opensnoop
```shell
root@ubuntu-2110:~/bcc/libbpf-tools# ./opensnoop 
PID    COMM              FD ERR PATH
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.gitignore
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf/bpf.h
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf/bpf_core_read.h
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf/bpf_endian.h
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf/bpf_helper_defs.h
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf/bpf_helpers.h
1492   sandbox-agent     10   0 /root/bcc/libbpf-tools/.output/bpf/bpf_tracing.h
[...]
3178   opensnoop         17   0 /etc/localtime
```

🙌🏽 I ran my first eBPF program! 

Note that as well as showing the string you defined, the trace line includes other useful contextual information - for example, the name of the executable and the process ID that triggered the event that ran the program - in this example, sandbox-agent running as PID 9061.

## Summary

An eBPF program can gather useful information about the event that triggered it - which could be reported to user space for observability purposes, for example.

## To go further

- 📙 [What is eBPF](https://isovalent.com/data/liz-rice-what-is-ebpf.pdf) by Liz Rice
- 🔗 https://ebpf.io/
