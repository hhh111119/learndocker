package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgMemoryHierachyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程
		fmt.Printf("子进程current pid %v", os.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = new(syscall.SysProcAttr)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // hostname 隔离
			syscall.CLONE_NEWIPC | // ipc 隔离
			syscall.CLONE_NEWPID | // 进程 id, 让当前进程 pid 1
			syscall.CLONE_NEWNS | // mount 隔离
			syscall.CLONE_NEWUSER | // 用户组 id 和 group id 重新映射
			syscall.CLONE_NEWNET, // 网络 ns
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("父进程获取子进程id %v", cmd.Process.Pid)
		// os.Mkdir("tttt", 0755)

		// 创建 子 cg
		os.Mkdir(path.Join(cgMemoryHierachyMount, "testmemorylimit"), 0755)
		// 将 子进程的 pid 加入 cg
		ioutil.WriteFile(path.Join(cgMemoryHierachyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		// 限制内存为 100 m
		ioutil.WriteFile(path.Join(cgMemoryHierachyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
		cmd.Process.Wait()
	}
}
