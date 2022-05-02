package main

import (
	"fmt"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

func commitContainer(imageName string) {
	mntURL := "/home/howu/go/src/my_docker/rootfs/mnt/"
	imageTar := "/home/howu/study/docker_stu/" + imageName + ".tar"
	fmt.Printf("%s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil { // 最后-C是切换到指定目录，.代表打包当前目录
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
