package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/TokyoYoshida/go-cgroup"
)

func printBanner() {
	banner := `
   _|_|_|    _|_|    _|_|_|_|_|  _|_|_|_|_|  _|_|_|_|  _|_|_|
 _|        _|    _|      _|          _|      _|        _|    _|
 _|        _|_|_|_|      _|          _|      _|_|_|    _|_|_|
 _|        _|    _|      _|          _|      _|        _|    _|
   _|_|_|  _|    _|      _|          _|      _|_|_|_|  _|    _|
	 `
	fmt.Println(banner)
}

func main() {

	printBanner()

	cgroup.Init()

	cg := cgroup.NewCgroup("bar")
	ccpu := cg.AddController("cpu")
	if ccpu == nil {
		fmt.Println("error")
		return
	}
	cmem := cg.AddController("memory")
	if ccpu == nil {
		fmt.Println("error")
		return
	}
	err := ccpu.AddValueInt64("cpu.shares", 500)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cmem.AddValueInt64("memory.limit_in_bytes", 1048576)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cg.Create()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cg.AttachTask()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctls, err := cgroup.GetAllControllers()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range ctls {
		fmt.Println(ctls[i])
	}

	if err = syscall.Unshare(syscall.CLONE_NEWPID); err != nil {
		fmt.Println(err)
		return
	}

	err = syscall.Chroot("./rootfs")
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("/bin/bash", "-c", "/catter_impl")
	if err != nil {
		fmt.Println(err)
		return
	}
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
