package cni

import (
	"os"

	"github.com/pkg/errors"
	"github.com/projecteru2/docker-cni/utils"
)

func (p netnsManager) GetNetns(ID string) (netnsPath string, err error) {
	_, err = os.Stat(p.getNetnsPath(ID))
	return p.getNetnsPath(ID), errors.WithStack(err)
}

func (p netnsManager) CreateNetns(ID string) (netnsPath string, _ *utils.Process, _ error) {
	proc, err := utils.NewProcess("ip", []string{"net", "a", p.getID(ID)}, nil, nil)
	return p.getNetnsPath(ID), proc, err
}

func (p netnsManager) DeleteNetns(ID string) (*utils.Process, error) {
	return utils.NewProcess("ip", []string{"net", "d", p.getID(ID)}, nil, nil)
}

func (p netnsManager) getNetnsPath(ID string) string {
	return "/var/run/netns/" + p.getID(ID)
}

func (p netnsManager) getID(ID string) string {
	return ID[:12]
}
