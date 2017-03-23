package kawasaki

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter . HostFileCompiler
type HostFileCompiler interface {
	Compile(log lager.Logger, containerIp net.IP, handle string) ([]byte, error)
}

//go:generate counterfeiter . NameserversDeterminer
type NameserversDeterminer interface {
	Determine(resolvContents string, hostIP net.IP, pluginNameservers, operatorNameservers, additionalNameservers []net.IP) []net.IP
}

//go:generate counterfeiter . NameserversSerializer
type NameserversSerializer interface {
	Serialize([]net.IP) []byte
}

type ResolvConfigurer struct {
	HostsFileCompiler     HostFileCompiler
	NameserversDeterminer NameserversDeterminer
	NameserversSerializer NameserversSerializer
	ResolvFilePath        string
	DepotDir              string
}

func (d *ResolvConfigurer) Configure(log lager.Logger, cfg NetworkConfig, pid int) error {
	log = log.Session("dns-resolve-configure")

	containerHostsContents, err := d.HostsFileCompiler.Compile(log, cfg.ContainerIP, cfg.ContainerHandle)
	if err != nil {
		log.Error("compiling-hosts-file", err)
		return err
	}

	if err := writeExistingFile(filepath.Join(d.DepotDir, cfg.ContainerHandle, "hosts"), containerHostsContents); err != nil {
		log.Error("writing-hosts-file", err)
		return err
	}

	hostResolvContents, err := ioutil.ReadFile(d.ResolvFilePath)
	if err != nil {
		log.Error("reading-host-resolv-file", err)
		return err
	}
	nameservers := d.NameserversDeterminer.Determine(string(hostResolvContents), cfg.BridgeIP, cfg.PluginNameservers, cfg.OperatorNameservers, cfg.AdditionalNameservers)
	containerResolvContents := d.NameserversSerializer.Serialize(nameservers)

	if err := writeExistingFile(filepath.Join(d.DepotDir, cfg.ContainerHandle, "resolv.conf"), containerResolvContents); err != nil {
		log.Error("writing-resolv-file", err)
		return err
	}

	return nil
}

func writeExistingFile(path string, contents []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(contents); err != nil {
		return err
	}
	return nil
}
