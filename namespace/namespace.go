/*
UTS Namespace主要用来隔离nodename和domainname两个系统标识。在UTS namespace里，每个namespace允许有自己的hostname。
系统API 中的clone()创建新的进程。根据填入的参数来判断哪些namesapce会被创建，而且它们的子进程也会被包含到这些namespace中。
*/
package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh") // 指定被fork出来的新进程内的初始命令
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, //使用CLONE_NEWUTS标识来创建一个UTC namesapce
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil { //go封装了对于系统clone()函数的调用，这段代码执行后会进入一个sh运行环境中
		log.Fatal(err)
	}
}
