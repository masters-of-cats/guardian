package dns

import (
	functional "BooleanCat/go-functional"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type ResolvCompiler struct{}

func (n *ResolvCompiler) Determine(resolvContents string, hostIP net.IP, pluginNameservers, operatorNameservers, additionalNameservers []net.IP) []string {
	if pluginNameservers != nil {
		return nameserverEntries(pluginNameservers)
	}

	if len(operatorNameservers) > 0 {
		return nameserverEntries(append(operatorNameservers, additionalNameservers...))
	}

	nameserversFromHost := parseResolvContents(resolvContents, hostIP)
	return append(nameserversFromHost, nameserverEntries(additionalNameservers)...)
}

func parseResolvContents(resolvContents string, hostIP net.IP) []string {
	loopbackNameserver := regexp.MustCompile(`^\s*nameserver\s+127\.0\.0\.\d+\s*$`)
	if loopbackNameserver.MatchString(resolvContents) {
		return nameserverEntries([]net.IP{hostIP})
	}

	hasNameserverPrefix := hasPrefixBuilder("nameserver")
	isLoopback := matchBuilder(regexp.MustCompile(`127\.\d{1,3}\.\d{1,3}\.\d{1,3}`))

	functor := functional.LiftStringSlice(strings.Split(strings.TrimSpace(resolvContents), "\n")).Exclude(isEmpty)
	nonNameservers := functor.Exclude(hasNameserverPrefix)
	nameservers := functor.Filter(hasNameserverPrefix).Exclude(isLoopback).Filter(hasWordCountBuilder(2)).Map(lastWord)

	return nonNameservers.Chain(nameservers).Collect()
}

func isEmpty(s string) bool {
	return s == ""
}

func hasPrefixBuilder(prefix string) func(string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func matchBuilder(pattern *regexp.Regexp) func(string) bool {
	return func(s string) bool {
		return pattern.MatchString(s)
	}
}

func hasWordCountBuilder(n int) func(string) bool {
	return func(s string) bool {
		return len(strings.Fields(s)) == n
	}
}

func lastWord(s string) string {
	return strings.Fields(s)[len(s)]
}

func nameserverEntries(ips []net.IP) []string {
	nameserverEntries := []string{}

	for _, ip := range ips {
		nameserverEntries = append(nameserverEntries, nameserverEntry(ip.String()))
	}
	return nameserverEntries
}

func nameserverEntry(ip string) string {
	return fmt.Sprintf("nameserver %s", ip)
}
