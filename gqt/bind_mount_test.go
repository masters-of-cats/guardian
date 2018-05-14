package gqt_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gqt/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Bind Mount", func() {
	var (
		client *runner.RunningGarden
	)

	BeforeEach(func() {
		client = runner.Start(config)
	})

	AfterEach(func() {
		Expect(client.DestroyAndStop()).To(Succeed())
	})

	DescribeTable("privileged container",
		func(descrFactory func() (bindMountDescriptor, error), verification ...verificationFunc) {
			for _, v := range verification {
				descr, err := descrFactory()
				Expect(err).NotTo(HaveOccurred())

				container := createContainer(client, descr, true)
				Expect(descr.postCreateFunc()).To(Succeed())
				Expect(err).NotTo(HaveOccurred())

				v(container, descr)
				Expect(client.Destroy(container.Handle())).To(Succeed())
				destroy(descr)
			}
		},
		Entry("source is a file and mode is RO", buildDescriptor(fileSource, readOnly), canRead("root"), canRead("alice"), cannotWrite("root"), cannotWrite("alice")),
		Entry("source is a file and mode is RW", buildDescriptor(fileSource, readWrite), canRead("root"), canRead("alice"), canWrite("root"), canWrite("alice")),
		Entry("source is a file that is a bind mount and mode is RO", buildDescriptor(fileSource, bindMountToSelf, readOnly), canRead("root"), canRead("alice"), cannotWrite("root"), cannotWrite("alice")),
		Entry("source is a file that is a bind mount and mode is RW", buildDescriptor(fileSource, bindMountToSelf, readWrite), canRead("root"), canRead("alice"), canWrite("root"), canWrite("alice")),

		Entry("source is a directory and mode is RO", buildDescriptor(directorySource, readOnly), canRead("root"), canRead("alice"), cannotWrite("root"), cannotWrite("alice")),
		// Alice should not be able to write since host dir is owned by root
		Entry("source is a directory and mode is RW", buildDescriptor(directorySource, readWrite), canRead("root"), canRead("alice"), canWrite("root"), cannotWrite("alice")),

		// We test that everyone can read/write into the nested bind mount
		Entry("source is a directory with a nested mount and mode is RO", buildDescriptor(directorySource, bindMountToSelf, readOnly, createNestedMount), canRead("root"), canRead("alice"), canWrite("root"), canWrite("alice")),
		Entry("source is a directory with a nested mount and mode is RW", buildDescriptor(directorySource, bindMountToSelf, readWrite, createNestedMount), canRead("root"), canRead("alice"), canWrite("root"), canWrite("alice")),
	)
})

type bindMountDescriptor struct {
	garden.BindMount
	fileToCheck    string
	mountpoints    []string
	postCreateFunc func() error
}

func destroy(descriptor bindMountDescriptor) {
	for i := len(descriptor.mountpoints) - 1; i >= 0; i-- {
		unmount(descriptor.mountpoints[i])
	}

	Expect(os.RemoveAll(descriptor.BindMount.SrcPath)).To(Succeed())
}

func buildDescriptor(enhancers ...bindMountDescriptorEnhancer) func() (bindMountDescriptor, error) {
	return func() (bindMountDescriptor, error) {
		descriptor := bindMountDescriptor{postCreateFunc: func() error { return nil }}
		for _, e := range enhancers {
			var err error
			descriptor, err = e(descriptor)
			if err != nil {
				return bindMountDescriptor{}, err
			}
		}

		return descriptor, nil
	}
}

// Enhancers
type bindMountDescriptorEnhancer func(descr bindMountDescriptor) (bindMountDescriptor, error)

func fileSource(descr bindMountDescriptor) (bindMountDescriptor, error) {
	file, err := ioutil.TempFile("", "file-")
	if err != nil {
		return bindMountDescriptor{}, err
	}
	defer file.Close()

	err = os.Chmod(file.Name(), 0777)
	if err != nil {
		return bindMountDescriptor{}, err
	}

	descr.BindMount.SrcPath = file.Name()

	descr.BindMount.DstPath = "/home/alice/destination"
	descr.fileToCheck = "/home/alice/destination"

	return descr, nil
}

func directorySource(descr bindMountDescriptor) (bindMountDescriptor, error) {
	dir, err := ioutil.TempDir("", "dir-")
	if err != nil {
		return bindMountDescriptor{}, err
	}

	if err := os.Chown(dir, 0, 0); err != nil {
		return bindMountDescriptor{}, err
	}

	if err := os.Chmod(dir, 0755); err != nil {
		return bindMountDescriptor{}, err
	}

	filePath := filepath.Join(dir, "testfile")
	if err := ioutil.WriteFile(filePath, []byte{}, os.ModePerm); err != nil {
		return bindMountDescriptor{}, err
	}

	descr.BindMount.SrcPath = dir
	descr.BindMount.DstPath = "/home/alice/destination"
	descr.fileToCheck = "/home/alice/destination/testfile"

	return descr, nil
}

func readOnly(descr bindMountDescriptor) (bindMountDescriptor, error) {
	descr.BindMount.Mode = garden.BindMountModeRO
	return descr, nil
}

func readWrite(descr bindMountDescriptor) (bindMountDescriptor, error) {
	descr.BindMount.Mode = garden.BindMountModeRW
	return descr, nil
}

func bindMountToSelf(descr bindMountDescriptor) (bindMountDescriptor, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("mount", "--bind", descr.BindMount.SrcPath, descr.BindMount.SrcPath)
	if err := cmd.Run(); err != nil {
		return bindMountDescriptor{}, err
	}

	cmd = exec.Command("mount", "--make-shared", descr.BindMount.SrcPath)
	if err := cmd.Run(); err != nil {
		return bindMountDescriptor{}, err
	}

	descr.mountpoints = []string{descr.BindMount.SrcPath}
	return descr, nil
}

func createNestedMount(desc bindMountDescriptor) (bindMountDescriptor, error) {
	nestedBindPath := filepath.Join(desc.SrcPath, "nested-bind")
	desc.postCreateFunc = func() error {
		if err := os.MkdirAll(nestedBindPath, os.FileMode(0755)); err != nil {
			return err
		}

		cmd := exec.Command("mount", "-t", "tmpfs", "tmpfs", nestedBindPath)
		if err := cmd.Run(); err != nil {
			return err
		}

		filePath := filepath.Join(nestedBindPath, "nested-file")
		if err := ioutil.WriteFile(filePath, []byte{}, os.ModePerm); err != nil {
			return err
		}
		if err := os.Chmod(filePath, 0777); err != nil {
			return err
		}
		return nil
	}

	desc.mountpoints = append(desc.mountpoints, nestedBindPath)
	desc.fileToCheck = filepath.Join(desc.BindMount.DstPath, "nested-bind", "nested-file")

	return desc, nil
}

// Verification functions
type verificationFunc func(container garden.Container, descr bindMountDescriptor)

func canRead(user string) func(container garden.Container, descr bindMountDescriptor) {
	return func(container garden.Container, descr bindMountDescriptor) {
		Expect(containerReadFile(container, descr.fileToCheck, user)).To(Succeed(), "User %s cannot read %s in container %s", user, descr.fileToCheck, container.Handle())
	}
}

func canWrite(user string) func(container garden.Container, descr bindMountDescriptor) {
	return func(container garden.Container, descr bindMountDescriptor) {
		if isDir(descr.BindMount.SrcPath) {
			Expect(touchFile(container, filepath.Dir(descr.fileToCheck), user)).To(Succeed(), "User %s cannot touch files in directory %s in container %s", user, filepath.Dir(descr.fileToCheck), container.Handle())
		} else {
			Expect(writeFile(container, descr.fileToCheck, user)).To(Succeed(), "User %s cannot write file %s in container %s", user, descr.fileToCheck, container.Handle())
		}
	}
}

func cannotWrite(users ...string) func(container garden.Container, descr bindMountDescriptor) {
	return func(container garden.Container, descr bindMountDescriptor) {
		for _, u := range users {
			if isDir(descr.BindMount.SrcPath) {
				Expect(touchFile(container, filepath.Dir(descr.fileToCheck), u)).NotTo(Succeed(), "User %s should not be able to touch files in directory %s in container %s", u, filepath.Dir(descr.fileToCheck), container.Handle())
			} else {
				Expect(writeFile(container, descr.fileToCheck, u)).NotTo(Succeed(), "User %s should not be able to write file %s in container %s", u, descr.fileToCheck, container.Handle())
			}
		}
	}
}

// Utilities
func createContainer(client *runner.RunningGarden, descr bindMountDescriptor, privileged bool) garden.Container {
	container, err := client.Create(garden.ContainerSpec{
		Privileged: privileged,
		BindMounts: []garden.BindMount{descr.BindMount},
		Network:    fmt.Sprintf("10.0.%d.0/24", GinkgoParallelNode()),
	})
	Expect(err).NotTo(HaveOccurred())
	return container
}

func containerReadFile(container garden.Container, filePath, user string) error {
	process, err := container.Run(garden.ProcessSpec{
		Path: "cat",
		Args: []string{filePath},
		User: user,
	}, ginkgoIO)
	Expect(err).ToNot(HaveOccurred())

	exitCode, err := process.Wait()
	Expect(err).ToNot(HaveOccurred())
	if exitCode != 0 {
		return fmt.Errorf("Could not read file %s in container %s, exit code %d", filePath, container.Handle(), exitCode)
	}

	return nil
}

func unmount(mountpoint string) {
	cmd := exec.Command("umount", "-f", mountpoint)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Printf("Command: umount -f [%s]\n%v", mountpoint, string(output))
	}
	Expect(err).NotTo(HaveOccurred())
}

func isDir(path string) bool {
	stat, err := os.Stat(path)
	Expect(err).NotTo(HaveOccurred())
	return stat.IsDir()
}

func writeFile(container garden.Container, dstPath, user string) error {
	process, err := container.Run(garden.ProcessSpec{
		Path: "/bin/sh",
		Args: []string{"-c", fmt.Sprintf("echo i-can-write > %s", dstPath)},
		User: user,
	}, ginkgoIO)
	Expect(err).ToNot(HaveOccurred())

	exitCode, err := process.Wait()
	Expect(err).ToNot(HaveOccurred())
	if exitCode != 0 {
		return fmt.Errorf("Could not write to file %s in container %s, exit code %d", dstPath, container.Handle(), exitCode)
	}

	return nil
}

func touchFile(container garden.Container, dstDir, user string) error {
	fileToTouch := filepath.Join(dstDir, "can-touch-this")
	process, err := container.Run(garden.ProcessSpec{
		Path: "touch",
		Args: []string{fileToTouch},
		User: user,
	}, ginkgoIO)
	Expect(err).ToNot(HaveOccurred())

	exitCode, err := process.Wait()
	Expect(err).ToNot(HaveOccurred())
	if exitCode != 0 {
		return fmt.Errorf("Could not touch file %s in container %s, exit code %d", fileToTouch, container.Handle(), exitCode)
	}

	return nil
}

// // TODO - stuff below is for reference, delete it!
// // Option a)
// // * have two tables (one for privileged and one for unprivileged
// // * the table function creates and destroys a container based on the bind mount descriptor and performs the verification function supplied
// // * each entry represents all of the combinations of source type, RO vs RW bind mount mode, whether the source is a nested mount, whether the source is a symbolic link, etc.
// // * there are unlimited verification functions which can be applied on each bind mount descriptor, `canRead`, `canWrite` are just examples
// //
// // * cons: almost identical tables that differ on whether the container is privileged and verification functions, both tables would have the same bind mount descriptors (a bit of duplication
// // * pros:
// //		- readable tests
// //		- can focus a single entry
// //		- easy to add more test cases if needed
// //
// var _ = Describe("Math", func() {
// 	DescribeTable("privileged",
// 		func(descr bindMountDescriptor, verification ...verificationFunc) {
// 			for _, v := range verification {
// 				container := createPrivilegedContainer(descr)
// 				v(container, descr)
// 				destroyContainer(container)
// 			}
// 			destroy(descr)
// 		},
// 		Entry("description", createMountpointDirSrcBindMount("/home/alice/ro-mountpoint-dir", garden.BindMountModeRO), canRead("root"), canWrite("root"), canRead("alice"), cannotWrite("alice")),
// 		Entry("description", createMountpointDirSrcBindMount("/home/alice/ro-mountpoint-dir", garden.BindMountModeRO), rootCanRead, rootCanWrite, aliceCanRead, aliceCannotWrite),
// 	)
//
// 	DescribeTable("unprivileged",
// 		func(descr bindMountDescriptor, verification ...verificationFunc) {
// 			for _, v := range verification {
// 				container := createUnprivilegedContainer(descr)
// 				v(container, descr)
// 				destroyContainer(container)
// 			}
// 			destroy(descr)
// 		},
// 		Entry("description", createMountpointDirSrcBindMount("/home/alice/ro-mountpoint-dir", garden.BindMountModeRO), canRead, rootCanWrite, aliceCanRead, aliceCannotWrite),
// 		Entry("description", createMountpointDirSrcBindMount("/home/alice/ro-mountpoint-dir", garden.BindMountModeRO), rootCanRead, rootCanWrite, aliceCanRead, aliceCannotWrite),
// 	)
// })
//
// // Option b
// // * have a single container created with all possible bind mounts
// // * cons: multiple tests in a single container -> test polution
// // * pros: concise and understandable tests
// var _ = FDescribe("Bind mount", func() {
// 	var (
// 		client        *runner.RunningGarden
// 		container     garden.Container
// 		containerSpec garden.ContainerSpec
//
// 		bindMounts map[string]bindMountDescriptor
// 	)
//
// 	BeforeEach(func() {
// 		bindMounts = make(map[string]bindMountDescriptor)
//
// 		// bindMounts["ro-file"] = createFileSrcBindMount("/home/alice/file/ro", garden.BindMountModeRO)
// 		// bindMounts["rw-file"] = createFileSrcBindMount("/home/alice/file/rw", garden.BindMountModeRW)
// 		//
// 		// bindMounts["ro-dir"] = createDirSrcBindMount("/home/alice/dir/ro", garden.BindMountModeRO)
// 		// bindMounts["rw-dir"] = createDirSrcBindMount("/home/alice/dir/rw", garden.BindMountModeRW)
//
// 		bindMounts["ro-mountpoint-dir"] = createMountpointDirSrcBindMount("/home/alice/ro-mountpoint-dir", garden.BindMountModeRO)
// 		bindMounts["rw-mountpoint-dir"] = createMountpointDirSrcBindMount("/home/alice/rw-mountpoint-dir", garden.BindMountModeRW)
//
// 		bindMounts["ro-nested-mountpoint-dir"] = createNestedMountpointDirSrcBindMount("/home/alice/ro-nested-mountpoint-dir", garden.BindMountModeRO)
// 		bindMounts["rw-nested-mountpoint-dir"] = createNestedMountpointDirSrcBindMount("/home/alice/rw-nested-mountpoint-dir", garden.BindMountModeRW)
//
// 		containerSpec = garden.ContainerSpec{
// 			BindMounts: gardenBindMounts(bindMounts),
// 			Network:    fmt.Sprintf("10.0.%d.0/24", GinkgoParallelNode()),
// 		}
// 	})
//
// 	JustBeforeEach(func() {
// 		client = runner.Start(config)
//
// 		var err error
// 		container, err = client.Create(containerSpec)
// 		Expect(err).NotTo(HaveOccurred())
// 	})
//
// 	AfterEach(func() {
// 		for _, desc := range bindMounts {
// 			destroy(desc)
// 		}
//
// 		Expect(client.DestroyAndStop()).To(Succeed())
// 	})
//
// 	Context("when the continer is privileged", func() {
// 		BeforeEach(func() {
// 			containerSpec.Privileged = true
// 		})
//
// 		Context("User root", func() {
// 			It("can read", func() {
// 				canRead("root", container, all(bindMounts)...)
// 			})
//
// 			It("can write", func() {
// 				canWrite("root", container, readWrite(bindMounts)...)
// 			})
//
// 			It("cannot write", func() {
// 				canNotWrite("root", container, readOnly(bindMounts)...)
// 			})
// 		})
//
// 		Context("User alice", func() {
// 			It("can read", func() {
// 				canRead("alice", container, all(bindMounts)...)
// 			})
//
// 			It("can write", func() {
// 				canWrite("alice", container, readWrite(bindMounts)...)
// 			})
//
// 			It("cannot write", func() {
// 				canNotWrite("alice", container, readOnly(bindMounts)...)
// 			})
// 		})
// 	})
//
// 	Context("when the container is not privileged", func() {
// 		BeforeEach(func() {
// 			containerSpec.Privileged = false
// 		})
//
// 		Context("User root", func() {
// 			It("can read", func() {
// 				canRead("root", container, all(bindMounts)...)
// 			})
//
// 			It("can write", func() {
// 				canWrite("root", container, readWrite(bindMounts)...)
// 			})
//
// 			It("cannot write", func() {
// 				canNotWrite("root", container, readOnly(bindMounts)...)
// 			})
// 		})
//
// 		Context("User alice", func() {
// 			It("can read", func() {
// 				canRead("alice", container, all(bindMounts)...)
// 			})
//
// 			It("can write", func() {
// 				canWrite("alice", container, readWrite(bindMounts)...)
// 			})
//
// 			It("cannot write", func() {
// 				canNotWrite("alice", container, readOnly(bindMounts)...)
// 			})
// 		})
//
// 	})
//
// })
//
// func unmount(mountpoint string) {
// 	cmd := exec.Command("umount", "-f", mountpoint)
// 	output, err := cmd.CombinedOutput()
// 	if len(output) > 0 {
// 		fmt.Printf("Command: umount -f [%s]\n%v", mountpoint, string(output))
// 	}
// 	Expect(err).NotTo(HaveOccurred())
// }
//
// func createFileSrcBindMount(dstPath string, mode garden.BindMountMode) bindMountDescriptor {
// 	file, err := ioutil.TempFile("", fmt.Sprintf("file-%s-", humanise(mode)))
// 	Expect(err).NotTo(HaveOccurred())
// 	defer file.Close()
//
// 	return bindMountDescriptor{
// 		BindMount: garden.BindMount{
// 			SrcPath: file.Name(),
// 			DstPath: dstPath,
// 			Mode:    mode,
// 		},
// 		fileToCheck: file.Name(),
// 	}
// }
//
// func createDirSrcBindMount(dstPath string, mode garden.BindMountMode) bindMountDescriptor {
// 	dir, err := ioutil.TempDir("", fmt.Sprintf("dir-%s-", humanise(mode)))
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(os.Chown(dir, 0, 0)).To(Succeed())
// 	Expect(os.Chmod(dir, 0755)).To(Succeed())
// 	filePath := filepath.Join(dir, "testfile")
// 	Expect(ioutil.WriteFile(filePath, []byte{}, os.ModePerm)).To(Succeed())
// 	// Expect(os.Chmod(filePath, 0777)).To(Succeed())
//
// 	return bindMountDescriptor{
// 		BindMount: garden.BindMount{
// 			SrcPath: dir,
// 			DstPath: dstPath,
// 			Mode:    mode,
// 		},
// 		fileToCheck: filepath.Join(dstPath, "testfile"),
// 	}
// }
//
// func createMountpointDirSrcBindMount(dstPath string, mode garden.BindMountMode) bindMountDescriptor {
// 	desc := createDirSrcBindMount(dstPath, mode)
//
// 	var cmd *exec.Cmd
// 	cmd = exec.Command("mount", "--bind", desc.BindMount.SrcPath, desc.BindMount.SrcPath)
// 	Expect(cmd.Run()).To(Succeed())
//
// 	cmd = exec.Command("mount", "--make-shared", desc.BindMount.SrcPath)
// 	Expect(cmd.Run()).To(Succeed())
//
// 	desc.mountpoints = []string{desc.BindMount.SrcPath}
//
// 	return desc
// }
//
// func createNestedMountpointDirSrcBindMount(dstPath string, mode garden.BindMountMode) bindMountDescriptor {
// 	desc := createMountpointDirSrcBindMount(dstPath, mode)
//
// 	nestedBindPath := filepath.Join(desc.SrcPath, "nested-bind")
// 	Expect(os.MkdirAll(nestedBindPath, os.FileMode(0755))).To(Succeed())
//
// 	cmd := exec.Command("mount", "-t", "tmpfs", "tmpfs", nestedBindPath)
// 	Expect(cmd.Run()).To(Succeed())
//
// 	filePath := filepath.Join(nestedBindPath, "nested-file")
// 	Expect(ioutil.WriteFile(filePath, []byte{}, os.ModePerm)).To(Succeed())
// 	Expect(os.Chmod(filePath, 0777)).To(Succeed())
//
// 	desc.mountpoints = append(desc.mountpoints, nestedBindPath)
// 	desc.BindMount.SrcPath = nestedBindPath
// 	desc.fileToCheck = filepath.Join(desc.BindMount.DstPath, "nested-file")
//
// 	return desc
// }
//
// func gardenBindMounts(bindMounts map[string]bindMountDescriptor) []garden.BindMount {
// 	mounts := []garden.BindMount{}
// 	for _, v := range bindMounts {
// 		mounts = append(mounts, v.BindMount)
// 	}
//
// 	return mounts
// }
//
// // TODO: maybe we need canTouch/canNotTouch as well to make sure we can create new files in the bind mount
//
// func canRead(user string, container garden.Container, descriptors ...bindMountDescriptor) {
// 	for _, d := range descriptors {
// 		Expect(containerReadFile(container, d.fileToCheck, user)).To(Succeed())
// 	}
// }
//
// func canNotRead(user string, container garden.Container, descriptors ...bindMountDescriptor) {
// 	for _, d := range descriptors {
// 		Expect(containerReadFile(container, d.fileToCheck, user)).NotTo(Succeed())
// 	}
// }
//
// func canWrite(user string, container garden.Container, descriptors ...bindMountDescriptor) {
// 	for _, d := range descriptors {
// 		if isDir(d.BindMount.SrcPath) {
// 			Expect(touchFile(container, filepath.Dir(d.fileToCheck), user)).To(Succeed())
// 		} else {
// 			Expect(writeFile(container, d.fileToCheck, user)).To(Succeed())
// 		}
// 	}
// }
//
// func canNotWrite(user string, container garden.Container, descriptors ...bindMountDescriptor) {
// 	for _, d := range descriptors {
// 		if isDir(d.BindMount.SrcPath) {
// 			Expect(touchFile(container, filepath.Dir(d.fileToCheck), user)).NotTo(Succeed())
// 		} else {
// 			Expect(writeFile(container, d.fileToCheck, user)).NotTo(Succeed())
// 		}
// 	}
// }
//
// func all(bindMounts map[string]bindMountDescriptor) []bindMountDescriptor {
// 	descs := []bindMountDescriptor{}
// 	for _, v := range bindMounts {
// 		descs = append(descs, v)
// 	}
//
// 	return descs
// }
//
// func readWrite(bindMounts map[string]bindMountDescriptor) []bindMountDescriptor {
// 	rw := []bindMountDescriptor{}
// 	for _, v := range bindMounts {
// 		if v.BindMount.Mode == garden.BindMountModeRW {
// 			rw = append(rw, v)
// 		}
// 	}
//
// 	return rw
// }
//
// func readOnly(bindMounts map[string]bindMountDescriptor) []bindMountDescriptor {
// 	ro := []bindMountDescriptor{}
// 	for _, v := range bindMounts {
// 		if v.BindMount.Mode == garden.BindMountModeRO {
// 			ro = append(ro, v)
// 		}
// 	}
//
// 	return ro
// }
//
// func containerReadFile(container garden.Container, filePath, user string) error {
// 	process, err := container.Run(garden.ProcessSpec{
// 		Path: "cat",
// 		Args: []string{filePath},
// 		User: user,
// 	}, ginkgoIO)
// 	Expect(err).ToNot(HaveOccurred())
//
// 	exitCode, err := process.Wait()
// 	Expect(err).ToNot(HaveOccurred())
// 	if exitCode != 0 {
// 		return fmt.Errorf("Could not read file %s in container %s, exit code %d", filePath, container.Handle(), exitCode)
// 	}
//
// 	return nil
// }
//
// func writeFile(container garden.Container, dstPath, user string) error {
// 	process, err := container.Run(garden.ProcessSpec{
// 		Path: "/bin/sh",
// 		Args: []string{"-c", fmt.Sprintf("echo i-can-write > %s", dstPath)},
// 		User: user,
// 	}, ginkgoIO)
// 	Expect(err).ToNot(HaveOccurred())
//
// 	exitCode, err := process.Wait()
// 	Expect(err).ToNot(HaveOccurred())
// 	if exitCode != 0 {
// 		return fmt.Errorf("Could not write to file %s in container %s, exit code %d", dstPath, container.Handle(), exitCode)
// 	}
//
// 	return nil
// }
//
// func touchFile(container garden.Container, dstDir, user string) error {
// 	fileToTouch := filepath.Join(dstDir, "can-touch-this")
// 	process, err := container.Run(garden.ProcessSpec{
// 		Path: "touch",
// 		Args: []string{fileToTouch},
// 		User: user,
// 	}, ginkgoIO)
// 	Expect(err).ToNot(HaveOccurred())
//
// 	exitCode, err := process.Wait()
// 	Expect(err).ToNot(HaveOccurred())
// 	if exitCode != 0 {
// 		return fmt.Errorf("Could not touch file %s in container %s, exit code %d", fileToTouch, container.Handle(), exitCode)
// 	}
//
// 	return nil
// }
//
// func humanise(mode garden.BindMountMode) string {
// 	if mode == garden.BindMountModeRW {
// 		return "rw"
// 	}
// 	return "ro"
// }
//
// func isDir(path string) bool {
// 	stat, err := os.Stat(path)
// 	Expect(err).NotTo(HaveOccurred())
// 	return stat.IsDir()
// }
