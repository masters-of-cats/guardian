package iodaemon_test

import (
	"time"

	"io/ioutil"
	"os"
	"path/filepath"

	"bytes"
	"io"

	"github.com/cloudfoundry-incubator/guardian/rundmc/iodaemon"
	linkpkg "github.com/cloudfoundry-incubator/guardian/rundmc/iodaemon/link"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type wc struct {
	*bytes.Buffer
}

func (b wc) Close() error {
	return nil
}

var _ = Describe("Iodaemon", func() {
	var (
		socketPath       string
		tmpdir           string
		fakeOut          wc
		fakeErr          wc
		expectedExitCode int

		wirer  *iodaemon.Wirer
		daemon *iodaemon.Daemon

		exited chan struct{}
	)

	const DEFAULT_TIMEOUT = "3s"

	BeforeEach(func() {
		var err error
		expectedExitCode = 0
		tmpdir, err = ioutil.TempDir("", "socket-dir")
		Expect(err).ToNot(HaveOccurred())

		socketPath = filepath.Join(tmpdir, "iodaemon.sock")

		exited = make(chan struct{})

		fakeOut = wc{
			bytes.NewBuffer([]byte{}),
		}
		fakeErr = wc{
			bytes.NewBuffer([]byte{}),
		}

		wirer = &iodaemon.Wirer{}
		daemon = &iodaemon.Daemon{}
	})

	AfterEach(func() {
		defer os.RemoveAll(tmpdir)

		Eventually(exited, DEFAULT_TIMEOUT).Should(BeClosed())

		By("tidying up the socket file")
		if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
			Fail("socket file not cleaned up")
		}
	})

	Context("spawning a process: when no listeners connect", func() {
		spawnProcess := func(args ...string) {
			go func() {
				iodaemon.Spawn(socketPath, args, time.Second, fakeOut, wirer, daemon)
				close(exited)
			}()
		}

		It("times out when no listeners connect", func() {
			spawnProcess("echo", "hello")

			Eventually(exited, DEFAULT_TIMEOUT).Should(BeClosed())
		})
	})

	Context("spawning a process: when listeners connect", func() {
		spawnProcess := func(args ...string) {
			go func() {
				defer GinkgoRecover()
				Expect(iodaemon.Spawn(socketPath, args, time.Second, fakeOut, wirer, daemon)).To(Succeed())
				close(exited)
			}()
		}

		It("reports back stdout", func() {
			spawnProcess("echo", "hello")

			_, linkStdout, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say("hello\n"))
		})

		It("supports re-linking to an iodaemon instance", func() {
			spawnProcess("bash")

			l, _, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())
			err = l.Writer.TerminateConnection()
			Expect(err).ToNot(HaveOccurred())

			m, _, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			_, err = m.Write([]byte("exit\n"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("reports back stderr", func() {
			spawnProcess("bash", "-c", "echo error 1>&2")

			_, _, linkStderr, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStderr, DEFAULT_TIMEOUT).Should(gbytes.Say("error\n"))
		})

		It("sends stdin to child", func() {
			spawnProcess("env", "-i", "bash", "--noprofile", "--norc")

			l, linkStdout, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			_, err = l.Write([]byte("echo hello\n"))
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say(".*hello.*"))

			_, err = l.Write([]byte("exit\n"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("exits when the child exits", func() {
			spawnProcess("bash")

			l, _, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			_, err = l.Write([]byte("exit\n"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("closes stdin when the link is closed", func() {
			spawnProcess("bash")

			l, _, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			Expect(l.Close()).To(Succeed()) //bash will normally terminate when it receives EOF on stdin
		})

		Context("when there is an existing socket file", func() {
			BeforeEach(func() {
				file, err := os.Create(socketPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(file.Close()).To(Succeed())
			})

			It("still creates the process", func() {
				spawnProcess("echo", "hello")

				_, linkStdout, _, err := createLink(socketPath)
				Expect(err).ToNot(HaveOccurred())
				Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say("hello\n"))
			})
		})
	})

	Context("spawning a tty", func() {
		spawnTty := func(args ...string) {
			go func() {
				defer GinkgoRecover()
				Expect(iodaemon.Spawn(socketPath, args, time.Second, fakeOut, wirer, daemon)).To(Succeed())
				close(exited)
			}()
		}

		BeforeEach(func() {
			wirer.WithTty = true
			wirer.WindowColumns = 200
			wirer.WindowRows = 80
			daemon.WithTty = true
		})

		It("reports back stdout", func() {
			spawnTty("echo", "hello")

			_, linkStdout, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say("hello"))
		})

		It("reports back stderr to stdout", func() {
			spawnTty("bash", "-c", "echo error 1>&2")

			_, linkStdout, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say("error"))
		})

		It("exits when the child exits", func() {
			spawnTty("bash")

			l, _, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			_, err = l.Write([]byte("exit\n"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("closes stdin when the link is closed", func() {
			spawnTty("bash")

			l, _, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			Expect(l.Close()).To(Succeed()) //bash will normally terminate when it receives EOF on stdin
		})

		It("sends stdin to child", func() {
			spawnTty("env", "-i", "bash", "--noprofile", "--norc")

			l, linkStdout, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			_, err = l.Write([]byte("echo hello\n"))
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say(".*hello.*"))

			_, err = l.Write([]byte("exit\n"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly sets the window size", func() {
			spawnTty("env", "-i", "bash", "--noprofile", "--norc")

			l, linkStdout, _, err := createLink(socketPath)
			Expect(err).ToNot(HaveOccurred())

			_, err = l.Write([]byte("echo $COLUMNS $LINES\n"))
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say(".*\\s200 80\\s.*"))

			Expect(l.SetWindowSize(100, 40)).To(Succeed())

			_, err = l.Write([]byte("echo $COLUMNS $LINES\n"))
			Expect(err).ToNot(HaveOccurred())
			Eventually(linkStdout, DEFAULT_TIMEOUT).Should(gbytes.Say(".*\\s100 40\\s.*"))

			_, err = l.Write([]byte("exit\n"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

})

func createLink(socketPath string) (*linkpkg.Link, io.WriteCloser, io.WriteCloser, error) {
	linkStdout := gbytes.NewBuffer()
	linkStderr := gbytes.NewBuffer()
	var l *linkpkg.Link
	var err error
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		l, err = linkpkg.Create(socketPath, linkStdout, linkStderr)
		if err == nil {
			break
		}
	}
	return l, linkStdout, linkStderr, err
}
