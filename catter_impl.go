package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	err := syscall.Chroot("./rootfs")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = os.MkdirAll("/proc", 0770); err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll("/dev/pts", 0770); err != nil {
		fmt.Println(err)
		return
	}

	syscall.Unmount("/dev/pts", 0)

	if err := syscall.Mount("devpts", "/dev/pts", "devpts", 0, ""); err != nil {
		fmt.Println(err)
		return
	}

	syscall.Unmount("/proc", 0)

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("/bin/bash")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Wait()
}
