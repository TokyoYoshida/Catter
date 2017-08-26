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

	cgroup.Init()

	cg := cgroup.NewCgroup("bar")
	cg.AddController("cpu")
	cg.AddController("memory")
	cg.Create()
	cg.AttachTask()
	// cg.AddController("memory")

	ctls, err := cgroup.GetAllControllers()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range ctls {
		fmt.Println(ctls[i])
	}

	err = syscall.Chroot("./rootfs")
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("/bin/bash")
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	printBanner()
}
