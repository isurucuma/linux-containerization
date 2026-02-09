package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: run <cmd> [args...] | child <cmd> [args...]")
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
	if len(os.Args) < 3 {
		panic("usage: run <cmd> [args...]")
	}

	fmt.Printf("[parent] cmd=%v pid=%d\n", os.Args[2:], os.Getpid())

	args := append([]string{"child"}, os.Args[2:]...)
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
	if len(os.Args) < 3 {
		panic("usage: child <cmd> [args...]")
	}

	command := os.Args[2]
	commandArgs := os.Args[3:]

	fmt.Printf("[child] cmd=%v pid=%d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(command, commandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
