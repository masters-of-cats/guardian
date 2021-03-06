package goci

import "os/exec"

// RuncBinary is the path to a runc binary.
type RuncBinary struct {
	Path string
	Root string
}

// StartCommand returns an *exec.Cmd that, when run, will execute a given bundle.
func (runc RuncBinary) StartCommand(path, id string, detach bool, log string) *exec.Cmd {
	args := runc.addGlobalFlags([]string{"start"}, log)
	if detach {
		args = append(args, "-d")
	}

	args = append(args, id)

	cmd := exec.Command(runc.Path, args...)
	cmd.Dir = path
	return cmd
}

// ExecCommand returns an *exec.Cmd that, when run, will execute a process spec
// in a running container.
func (runc RuncBinary) ExecCommand(id, processJSONPath, pidFilePath string) *exec.Cmd {
	return exec.Command(
		runc.Path, runc.addRootFlagIfNeeded([]string{"exec", id, "--pid-file", pidFilePath, "-p", processJSONPath})...,
	)
}

// EventsCommand returns an *exec.Cmd that, when run, will retrieve events for the container
func (runc RuncBinary) EventsCommand(id string) *exec.Cmd {
	return exec.Command(runc.Path, runc.addRootFlagIfNeeded([]string{"events", id})...)
}

// KillCommand returns an *exec.Cmd that, when run, will signal the running
// container.
func (runc RuncBinary) KillCommand(id, signal, logFile string) *exec.Cmd {
	return exec.Command(
		runc.Path, runc.addGlobalFlags([]string{"kill", id, signal}, logFile)...,
	)
}

// StateCommand returns an *exec.Cmd that, when run, will get the state of the
// container.
func (runc RuncBinary) StateCommand(id, logFile string) *exec.Cmd {
	return exec.Command(runc.Path, runc.addGlobalFlags([]string{"state", id}, logFile)...)
}

// StatsCommand returns an *exec.Cmd that, when run, will get the stats of the
// container.
func (runc RuncBinary) StatsCommand(id, logFile string) *exec.Cmd {
	return exec.Command(runc.Path, runc.addGlobalFlags([]string{"events", "--stats", id}, logFile)...)
}

// DeleteCommand returns an *exec.Cmd that, when run, will signal the running
// container.
func (runc RuncBinary) DeleteCommand(id string, force bool, logFile string) *exec.Cmd {
	deleteArgs := runc.addGlobalFlags([]string{"delete"}, logFile)
	if force {
		deleteArgs = append(deleteArgs, "--force")
	}
	return exec.Command(runc.Path, append(deleteArgs, id)...)
}

func (runc RuncBinary) addGlobalFlags(args []string, logFile string) []string {
	baseFlags := append([]string{"--debug", "--log", logFile, "--log-format", "json"}, args...)
	return runc.addRootFlagIfNeeded(baseFlags)
}

func (runc RuncBinary) addRootFlagIfNeeded(args []string) []string {
	if runc.Root == "" {
		return args
	}
	return append([]string{"--root", runc.Root}, args...)
}
