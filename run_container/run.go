package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"my_docker/run_container/cgroups"
	"my_docker/run_container/cgroups/subsystems"
	"my_docker/run_container/container"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume, containerName, imageName string, envSlice []string) {
	//生成容器名，防止用户没有传容器名
	containerID := randStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envSlice)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 记录容器信息
	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName, containerID, volume)
	if err != nil {
		log.Errorf("record container info error %v", err)
		return
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker_cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	if tty {
		parent.Wait()
		container.DeleteWorkSpace(volume, containerName)
		deleteContainerInfo(containerName)
	}
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func recordContainerInfo(containerPID int, commandArray []string, containerName, containerId, volume string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, " ")
	containerInfo := &container.ContainerInfo{
		Id:          containerId,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
		Volume:      volume,
	}

	// go json转string方式
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}
	defer file.Close()
	return containerName, nil
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}
