package gqt_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gqt/runner"

	uuid "github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cgroup bug testing", func() {
	var (
		cgroupPathTemplate string
		client             *runner.RunningGarden
	)

	BeforeEach(func() {
		id, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())

		cgroupPathTemplate = fmt.Sprintf("/tmp/cgroups-%d/%%s/test_%s", GinkgoParallelNode(), id.String())

		client = runner.Start(config)
		_, err = client.Create(garden.ContainerSpec{})
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		Expect(client.DestroyAndStop()).To(Succeed())
	})

	FIt("keeps adding processes to cgroup.procs", func() {
		// subsystems := []string{"cpuset", "cpu", "cpuacct", "blkio", "memory", "devices", "freezer", "net_cls", "perf_event", "net_prio", "hugetlb", "pids"}
		subsystems := []string{"cpu", "cpuacct", "memory"}
		allLogFiles := []*os.File{}

		errorChan := make(chan error)

		parallelRoutinesCount := 10
		wg := &sync.WaitGroup{}
		wg.Add(parallelRoutinesCount)
		done := doneChannel(wg)

		processCount := 10

		for i := 0; i < parallelRoutinesCount; i++ {
			logFile, err := os.OpenFile(fmt.Sprintf("/tmp/cgroup-test-%d", i), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			Expect(err).NotTo(HaveOccurred())
			defer logFile.Close()

			allLogFiles = append(allLogFiles, logFile)

			go testRoutine(wg, errorChan, processCount, logFile, cgroupPathTemplate+fmt.Sprintf("/routine-%d", i), subsystems)()
		}

		select {
		case err := <-errorChan:
			Expect(err).NotTo(HaveOccurred())
		case <-done:
			// allContents := ""
			// for _, lf := range allLogFiles {
			// 	contents, err := ioutil.ReadFile(lf.Name())
			// 	Expect(err).NotTo(HaveOccurred())
			// 	allContents += "\n\n\n -------------------- \n\n\n"
			// 	allContents += string(contents)
			// }
			// Expect(errors.New(allContents)).NotTo(HaveOccurred())
			Expect(errors.New("DONE")).NotTo(HaveOccurred())
		}
	})
})

func doneChannel(wg *sync.WaitGroup) chan struct{} {
	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}

func testRoutine(wg *sync.WaitGroup, errChan chan error, processCount int, logFile *os.File, cgroupPathTemplate string, subsystems []string) func() {
	return func() {
		defer wg.Done()

		for _, subs := range subsystems {
			path := fmt.Sprintf(cgroupPathTemplate, subs)
			log(logFile, fmt.Sprintf("Creating cgroup %s\n", path))
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				errChan <- fmt.Errorf("Failed to mkdriall: %v", err)
				return
			}
		}

		for i := 0; i < processCount; i++ {
			log(logFile, fmt.Sprintf("<<<<< Starting interation %d <<<<<<\n", i))

			cmd := exec.Command("sleep", "9999")
			log(logFile, fmt.Sprintf("Starting command %#v\n", cmd))
			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("failed to execute sleep:  %v", err)
				return
			}

			for _, subs := range subsystems {
				path := filepath.Join(fmt.Sprintf(cgroupPathTemplate, subs), "cgroup.procs")
				log(logFile, fmt.Sprintf("Adding pid %d to %s\n", cmd.Process.Pid, path))
				if err := ioutil.WriteFile(path, []byte(strconv.Itoa(cmd.Process.Pid)), 0700); err != nil {
					errChan <- fmt.Errorf("Failed to write pid %d to %s:  %v", cmd.Process.Pid, path, err)
					return
				}
			}

			log(logFile, fmt.Sprintf("Killing process with pid %d\n", cmd.Process.Pid))
			if err := cmd.Process.Kill(); err != nil {
				errChan <- fmt.Errorf("failed to kill process:  %v", err)
				return
			}
			log(logFile, fmt.Sprintf("Waiting for process with pid %d\n", cmd.Process.Pid))
			_, err := cmd.Process.Wait()
			if err != nil {
				errChan <- fmt.Errorf("failed to wait for process:  %v", err)
				return
			}

			log(logFile, fmt.Sprintf("<<<<< End of iteration %d <<<<<<\n", i))
		}
	}
}

func log(logFile *os.File, message string) {
	_, err := logFile.WriteString(message)
	if err != nil {
		panic(fmt.Errorf("Failed to log into %s: %v", logFile.Name(), err))
	}

}

// func addToCgroup(cgroupPath string, pid int) error {
// 	Expect(ioutil.WriteFile(filepath.Join(cgroupPath, "cpu,cpuacct", "cgroup.procs"), []byte(strconv.Itoa(pid)), 0700))
// 	return nil
// }
