package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// 在挂载点信息中查找某个subsystem的hierarchy的根节点
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreat bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreat && os.IsNotExist(err)) {
		if os.IsNotExist(err) { // 文件不存在
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
			} else {
				return "", fmt.Errorf("err creat cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error: %v", err)
	}
}
