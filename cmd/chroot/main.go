package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: run <rootfs> <cmd> [args...] | child <rootfs> <cmd> [args...]")
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("unknown command")
	}
}

func run() {
	if len(os.Args) < 4 {
		panic("usage: run <rootfs> <cmd> [args...]")
	}

	rootfs := os.Args[2]
	fmt.Printf("[parent] rootfs=%s cmd=%v pid=%d\n", rootfs, os.Args[3:], os.Getpid())

	args := append([]string{"child", rootfs}, os.Args[3:]...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func child() {
	if len(os.Args) < 4 {
		panic("usage: child <rootfs> <cmd> [args...]")
	}

	rootfs := os.Args[2]
	command := os.Args[3]
	commandArgs := os.Args[4:]

	fmt.Printf("[child] rootfs=%s cmd=%v pid=%d\n", rootfs, os.Args[3:], os.Getpid())

	if err := syscall.Sethostname([]byte("alpine-container")); err != nil {
		panic(err)
	}

	if err := syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC), ""); err != nil {
		panic(err)
	}

	if err := syscall.Chroot(rootfs); err != nil {
		panic(err)
	}

	if err := os.Chdir("/"); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("/proc", 0o755); err != nil {
		panic(err)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		panic(err)
	}
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command(command, commandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
