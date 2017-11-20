package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/TokyoYoshida/go-cgroup"
	"github.com/docker/libcontainer/system"
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

func setNSByPid(pid int) error {
	nsPath := "/proc/" + strconv.Itoa(pid) + "/ns/"
	ns, err := os.Open(nsPath)
	if err != nil {
		return err
	}
	return system.Setns(ns.Fd(), syscall.CLONE_NEWNET)
}

func setNSByOwnPid() error {
	pid := os.Getpid()
	error := setNSByPid(pid)
	return error
}

func initNetWork() {
	_, err := exec.Command("/bin/bash", "-c", "ip link add name veth0-host type veth peer name veth0-ct").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}

	_, err = exec.Command("/bin/bash", "-c", "ip addr add 10.10.10.10/24 dev veth0-host").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}

	_, err = exec.Command("/bin/bash", "-c", "ip link set up veth0-host").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}

	_, err = exec.Command("/bin/bash", "-c", "sudo ip netns add netns01").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}

	_, err = exec.Command("/bin/bash", "-c", "ip link set veth0-ct netns netns01").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}

	_, err = exec.Command("/bin/bash", "-c", "ip netns exec netns01 ip addr add 10.10.10.11/24 dev veth0-ct").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}

	_, err = exec.Command("/bin/bash", "-c", "ip netns exec netns01 ip link set veth0-ct up").Output()
	if err != nil {
		fmt.Print("error:")
		fmt.Println(err)
	}
}

func main() {

	printBanner()

	initNetWork()

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
	err = cmem.AddValueInt64("memory.limit_in_bytes", 1048576*2)
	if err != nil {
		fmt.Println(err)
		return
	}
	cdev := cg.AddController("devices")
	if cdev == nil {
		fmt.Println("error")
		return
	}
	cpid := cg.AddController("pids")
	if cpid == nil {
		fmt.Println("error")
		return
	}
	cnetc := cg.AddController("net_cls")
	if cnetc == nil {
		fmt.Println("error")
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

	cmd := exec.Command("ip", "netns", "exec", "netns01", "bash", "-c", "./catter_impl")
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
