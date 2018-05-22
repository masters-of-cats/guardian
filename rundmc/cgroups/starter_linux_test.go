package cgroups_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/guardian/rundmc/cgroups"
	"code.cloudfoundry.org/guardian/rundmc/cgroups/fs/fsfakes"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

var _ = Describe("CgroupStarter", func() {
	var (
		starter                 *cgroups.CgroupStarter
		logger                  lager.Logger
		fakeFS                  *fsfakes.FakeFS
		procCgroupsContents     string
		procSelfCgroupsContents string
		tmpDir                  string
	)

	BeforeEach(func() {
		tmpDir = tempDir("", "gdncgroup")

		procSelfCgroupsContents = ""
		procCgroupsContents = "#subsys_name\thierarchy\tnum_cgroups\tenabled\n" +
			"devices\t1\t1\t1\n"

		logger = lagertest.NewTestLogger("test")
		fakeFS = new(fsfakes.FakeFS)
	})

	JustBeforeEach(func() {

		starter = cgroups.NewStarter(
			logger,
			ioutil.NopCloser(strings.NewReader(procCgroupsContents)),
			ioutil.NopCloser(strings.NewReader(procSelfCgroupsContents)),
			path.Join(tmpDir, "cgroup"),
			"garden",
			[]specs.LinuxDeviceCgroup{{
				Type:   "c",
				Major:  int64ptr(10),
				Minor:  int64ptr(200),
				Access: "rwm",
			}},
		)
		starter.FS = fakeFS
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("mkdirs the cgroup path", func() {
		starter.Start()
		Expect(path.Join(tmpDir, "cgroup")).To(BeADirectory())
	})

	It("adds the right content into devices.allow", func() {
		Expect(starter.Start()).To(Succeed())

		Expect(path.Join(tmpDir, "cgroup", "devices", "garden")).To(BeADirectory())

		content := readFile(path.Join(tmpDir, "cgroup", "devices", "garden", "devices.allow"))
		Expect(string(content)).To(Equal("c 10:200 rwm"))
	})

	It("adds the right content into devices.deny", func() {
		Expect(starter.Start()).To(Succeed())

		Expect(path.Join(tmpDir, "cgroup", "devices", "garden")).To(BeADirectory())

		content := readFile(path.Join(tmpDir, "cgroup", "devices", "garden", "devices.deny"))
		Expect(string(content)).To(Equal("a"))
	})

	Context("when there is already a child device cgroup", func() {
		JustBeforeEach(func() {
			Expect(os.MkdirAll(path.Join(tmpDir, "cgroup", "devices", "garden", "child"), 0777)).To(Succeed())
		})

		It("does not write to devices.deny", func() {
			Expect(starter.Start()).To(Succeed())
			Expect(path.Join(tmpDir, "cgroup", "devices", "garden")).To(BeADirectory())
			Expect(path.Join(tmpDir, "cgroup", "devices", "garden", "devices.deny")).NotTo(BeAnExistingFile())
		})

	})

	It("mounts the cgroup tmpfs path", func() {
		Expect(starter.Start()).To(Succeed())
		expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup"), "tmpfs", 0, "uid=0,gid=0,mode=0755"))
	})

	Context("with a sane /proc/cgroups and /proc/self/cgroup", func() {
		BeforeEach(func() {
			procCgroupsContents = "#subsys_name\thierarchy\tnum_cgroups\tenabled\n" +
				"devices\t1\t1\t1\n" +
				"memory\t2\t1\t1\n" +
				"cpu\t3\t1\t1\n" +
				"cpuacct\t4\t1\t1\n"

			procSelfCgroupsContents = "5:devices:/\n" +
				"4:memory:/\n" +
				"3:cpu,cpuacct:/\n"
		})

		It("succeeds", func() {
			Expect(starter.Start()).To(Succeed())
		})

		It("mounts the hierarchies which are not already mounted", func() {
			Expect(starter.Start()).To(Succeed())

			Expect(fakeFS.MountCallCount()).To(Equal(5))

			expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup"), "tmpfs", 0, "uid=0,gid=0,mode=0755"))
			expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "devices"), "cgroup", 0, "devices"))
			expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "memory"), "cgroup", 0, "memory"))
			expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "cpu"), "cgroup", 0, "cpu,cpuacct"))
			expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "cpuacct"), "cgroup", 0, "cpu,cpuacct"))
		})

		It("creates needed directories", func() {
			starter.Start()
			Expect(path.Join(tmpDir, "cgroup", "devices")).To(BeADirectory())
		})

		It("creates subdirectories owned by the specified user and group", func() {
			Expect(starter.WithUID(123).WithGID(987).Start()).To(Succeed())
			allChowns := []string{}
			for i := 0; i < fakeFS.ChownCallCount(); i++ {
				path, uid, gid := fakeFS.ChownArgsForCall(i)
				allChowns = append(allChowns, path)
				Expect(uid).To(Equal(123))
				Expect(gid).To(Equal(987))
			}

			for _, subsystem := range []string{"devices", "cpu", "memory"} {
				fullPath := path.Join(tmpDir, "cgroup", subsystem, "garden")
				Expect(fullPath).To(BeADirectory())
				Expect(allChowns).To(ContainElement(fullPath))
				Expect(stat(fullPath).Mode() & os.ModePerm).To(Equal(os.FileMode(0755)))
			}
		})

		Context("when the garden folder already exists", func() {
			BeforeEach(func() {
				for _, subsystem := range []string{"devices", "cpu", "memory"} {
					fullPath := path.Join(tmpDir, "cgroup", subsystem, "garden")
					Expect(fullPath).ToNot(BeADirectory())
					Expect(os.MkdirAll(fullPath, 0700)).To(Succeed())
				}
			})

			It("changes the permissions of the subdirectories", func() {
				starter.Start()
				for _, subsystem := range []string{"devices", "cpu", "memory"} {
					fullPath := path.Join(tmpDir, "cgroup", subsystem, "garden")
					Expect(stat(fullPath).Mode() & os.ModePerm).To(Equal(os.FileMode(0755)))
				}
			})
		})

		Context("when we are in the nested case", func() {
			BeforeEach(func() {
				procCgroupsContents = "#subsys_name\thierarchy\tnum_cgroups\tenabled\n" +
					"memory\t2\t1\t1\n"

				procSelfCgroupsContents = "4:memory:/461299e6-b672-497c-64e5-793494b9bbdb\n"
			})

			It("creates subdirectories owned by the specified user and group", func() {
				Expect(starter.WithUID(123).WithGID(987).Start()).To(Succeed())
				allChowns := []string{}
				for i := 0; i < fakeFS.ChownCallCount(); i++ {
					path, uid, gid := fakeFS.ChownArgsForCall(i)
					Expect(uid).To(Equal(123))
					Expect(gid).To(Equal(987))
					allChowns = append(allChowns, path)
				}

				for _, subsystem := range []string{"memory"} {
					fullPath := path.Join(tmpDir, "cgroup", subsystem, "461299e6-b672-497c-64e5-793494b9bbdb", "garden")
					Expect(fullPath).To(BeADirectory())
					Expect(allChowns).To(ContainElement(fullPath))
					Expect(stat(fullPath).Mode() & os.ModePerm).To(Equal(os.FileMode(0755)))
				}
			})
		})

		Context("when a subsystem is not yet mounted anywhere", func() {
			BeforeEach(func() {
				procCgroupsContents = "#subsys_name\thierarchy\tnum_cgroups\tenabled\n" +
					"freezer\t7\t1\t1\n"
			})

			It("mounts it as its own subsystem", func() {
				Expect(starter.Start()).To(Succeed())
				expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "freezer"), "cgroup", 0, "freezer"))
			})
		})

		Context("when a subsystem is disabled", func() {
			BeforeEach(func() {
				procCgroupsContents = "#subsys_name\thierarchy\tnum_cgroups\tenabled\n" +
					"freezer\t7\t1\t0\n"
			})

			It("skips it", func() {
				Expect(starter.Start()).To(Succeed())
				expectNotMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "freezer"), "cgroup", 0, "freezer"))
			})
		})

		Context("when /proc/self/cgroup contains named cgroup hierarchies", func() {
			BeforeEach(func() {
				procSelfCgroupsContents = procSelfCgroupsContents + "1:name=systemd:/\n"
			})

			It("mounts it with name option as its own subsystem", func() {
				Expect(starter.Start()).To(Succeed())
				expectMounted(fakeFS, newMountArgs("cgroup", filepath.Join(tmpDir, "cgroup", "systemd"), "cgroup", 0, "name=systemd"))
			})
		})

		Context("when a cgroup is already mounted", func() {
			BeforeEach(func() {
				fakeFS.MountReturns(unix.EBUSY)
			})

			It("succeeds", func() {
				Expect(starter.Start()).To(Succeed())
			})
		})
	})

	Context("when /proc/cgroups contains malformed entries", func() {
		BeforeEach(func() {
			procCgroupsContents = "#subsys_name\thierarchy\tnum_cgroups\tenabled\n" +
				"devices\tA ONE AND A\t1\t1\n" +
				"memory\tTWO AND A\t1\t1\n" +
				"cpu\tTHREE AND A\t1\t1\n" +
				"cpuacct\tFOUR\t1\t1\n"

			procSelfCgroupsContents = "5:devices:/\n" +
				"4:memory:/\n" +
				"3:cpu,cpuacct:/\n"
		})

		It("returns CgroupsFormatError", func() {
			err := starter.Start()
			Expect(err).To(Equal(cgroups.CgroupsFormatError{Content: "devices\tA ONE AND A\t1\t1"}))
		})
	})

	Context("when /proc/cgroups is empty", func() {
		BeforeEach(func() {
			procCgroupsContents = ""

			procSelfCgroupsContents = "5:devices:/\n" +
				"4:memory:/\n" +
				"3:cpu,cpuacct:/\n"
		})

		It("returns CgroupsFormatError", func() {
			err := starter.Start()
			Expect(err).To(Equal(cgroups.CgroupsFormatError{Content: "(empty)"}))
		})
	})

	Context("when /proc/cgroups contains an unknown header scheme", func() {
		BeforeEach(func() {
			procCgroupsContents = "#subsys_name\tsome\tbogus\tcolumns\n" +
				"devices\t1\t1\t1" +
				"memory\t2\t1\t1" +
				"cpu\t3\t1\t1" +
				"cpuacct\t4\t1\t1"

			procSelfCgroupsContents = "5:devices:/\n" +
				"4:memory:/\n" +
				"3:cpu,cpuacct:/\n"
		})

		It("returns CgroupsFormatError", func() {
			err := starter.Start()
			Expect(err).To(Equal(cgroups.CgroupsFormatError{Content: "#subsys_name\tsome\tbogus\tcolumns"}))
		})
	})
})

func expectNotMounted(fakeFS *fsfakes.FakeFS, mntArgs mountArgs) {
	Expect(allMountArguments(fakeFS)).NotTo(ContainElement(mntArgs))
}

func expectMounted(fakeFS *fsfakes.FakeFS, mntArgs mountArgs) {
	Expect(allMountArguments(fakeFS)).To(ContainElement(mntArgs))
}

func allMountArguments(fakeFS *fsfakes.FakeFS) []mountArgs {
	Expect(fakeFS.MountCallCount()).To(BeNumerically(">", 0))
	var allMntArgs []mountArgs
	for i := 0; i < fakeFS.MountCallCount(); i++ {
		allMntArgs = append(allMntArgs, newMountArgs(fakeFS.MountArgsForCall(i)))
	}
	return allMntArgs
}
