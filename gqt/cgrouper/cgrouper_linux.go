package cgrouper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func CleanGardenCgroups(cgroupsRootPath, tag string) error {
	subsystems, err := ioutil.ReadDir(cgroupsRootPath)
	if err != nil {
		return fmt.Errorf("read cgroup root dir %s: %v", cgroupsRootPath, err)
	}

	for _, subsystem := range subsystems {
		if err = unmountIfExists(filepath.Join(cgroupsRootPath, subsystem.Name())); err != nil {
			return fmt.Errorf("unmount subsystem %s: %v", subsystem.Name(), err)
		}
	}

	if err = unmountIfExists(cgroupsRootPath); err != nil {
		return fmt.Errorf("unmount cgroup root %s: %v", cgroupsRootPath, err)
	}

	return nil
}

func unmountIfExists(unmountPath string) error {
	err := unix.Unmount(unmountPath, 0)
	// If no flags in unmount, EINVAL can only mean target does not exist
	if os.IsNotExist(err) || err == unix.EINVAL {
		return nil
	}
	return err
}
