package main

import (
	"my_docker/run_container/cgroups"
	"my_docker/run_container/cgroups/subsystems"
	"my_docker/run_container/container"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker_cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	if tty {
		parent.Wait()
		mntURL := "/home/howu/go/src/my_docker/rootfs/mnt/"
		rootURL := "/home/howu/go/src/my_docker/rootfs/"
		container.DeleteWorkSpace(rootURL, mntURL, volume)
	}
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
