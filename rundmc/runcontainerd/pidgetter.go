package runcontainerd

import (
	"context"

	"github.com/containerd/containerd"
)

type PidGetter struct {
	Client  *containerd.Client
	Context context.Context
}

func (p *PidGetter) GetPid(bundlePath, containerID string) (uint32, error) {
	container, err := p.Client.LoadContainer(p.Context, containerID)
	if err != nil {
		return 0, err
	}

	task, err := container.Task(p.Context, nil)
	if err != nil {
		return 0, err
	}
	return task.Pid(), err
}
