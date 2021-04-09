package oci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/projecteru2/docker-cni/utils"
	log "github.com/sirupsen/logrus"
)

func (c *ContainerMeta) ID() string {
	// the best guess I can get, otherwise we have to count on extra env passed from docker
	parts := strings.Split(c.Linux.CgroupsPath, "/")
	return parts[len(parts)-1]
}

func (c *ContainerMeta) Labels() map[string]string {
	return nil
}

func (c *ContainerMeta) UpdateNetns(netnsPath string) {
	for idx, ns := range c.Linux.Namespaces {
		if ns.Type == specs.NetworkNamespace {
			if ns.Path != "" {
				log.Warnf("netns path existed and have been replaced: %s", ns.Path)
			}
			c.Linux.Namespaces[idx] = specs.LinuxNamespace{
				Type: specs.NetworkNamespace,
				Path: netnsPath,
			}
		}
	}
}

func (c *ContainerMeta) AppendPoststopHook(process *utils.Process) {
	cmd := fmt.Sprintf("%s %s", process.Path, strings.Join(process.Args, " "))
	if process.Stdio != nil && process.Stdio.StdinBytes != nil {
		cmd += " <<<'" + strings.ReplaceAll(strings.ReplaceAll(string(process.StdinBytes), "\n", ""), " ", "") + "'"
	}
	c.Hooks.Poststop = append(c.Hooks.Poststop, specs.Hook{
		Path: "/bin/bash",
		Args: []string{"bash", "-c", cmd},
		Env:  process.Env,
	})
}

func (c *ContainerMeta) Save() (err error) {
	data, err := json.Marshal(c.Spec)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Debugf("save config")
	return errors.WithStack(ioutil.WriteFile(c.bundlePath, data, 0644))
}
