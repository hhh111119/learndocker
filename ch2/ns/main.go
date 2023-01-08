package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // hostname 隔离
			syscall.CLONE_NEWIPC | // ipc 隔离
			syscall.CLONE_NEWPID | // 进程 id, 让当前进程 pid 1
			syscall.CLONE_NEWNS | // mount 隔离
			syscall.CLONE_NEWUSER | // 用户组 id 和 group id 重新映射
			syscall.CLONE_NEWNET, // 网络 ns
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 1234,
				HostID:      0,
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 1234,
				HostID:      0,
				Size:        1,
			},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
