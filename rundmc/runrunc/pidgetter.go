package runrunc

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type FilePidGetter func(string, string) (uint32, error)

func (p FilePidGetter) GetPid(bundlePath, containerID string) (uint32, error) {
	return p(bundlePath, containerID)
}

func GetFilePid(bundlePath, containerID string) (uint32, error) {
	ctrInitPid, err := ioutil.ReadFile(filepath.Join(bundlePath, "pidfile"))
	if err != nil {
		return 0, err
	}
	pid, err := strconv.ParseUint(string(ctrInitPid), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(pid), nil
}
