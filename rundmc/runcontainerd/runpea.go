package runcontainerd

import (
	"io"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
)

type Creator interface {
	Create(log lager.Logger, bundlePath, id string, io garden.ProcessIO) error
}

type RunContainerPea struct {
	Creator Creator
}

func (r *RunContainerPea) Run(
	log lager.Logger, processID, processPath, sandboxHandle, sandboxBundlePath string,
	pio garden.ProcessIO, tty bool, procJSON io.Reader, extraCleanup func() error,
) (garden.Process, error) {

	if err := r.Creator.Create(log, processPath, processID, pio); err != nil {
		return &process{}, err
	}

	return &process{}, nil
}

func (r *RunContainerPea) Attach(log lager.Logger, processID string, io garden.ProcessIO, processesPath string) (garden.Process, error) {
	return &process{}, nil
}
