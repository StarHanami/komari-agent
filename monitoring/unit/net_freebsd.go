//go:build freebsd

package monitoring

import (
	"fmt"
	"os/exec"
	"strings"
)

// sockstat is part of the FreeBSD base system. The gopsutil implementation
// shells out to lsof, which is an optional package on FreeBSD.
func platformConnectionsCount() (int, int, error) {
	out, err := exec.Command("/usr/bin/sockstat", "-46", "-P", "tcp,udp").Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get sockets with sockstat: %w", err)
	}
	return countSockstat(string(out))
}

func countSockstat(output string) (tcp, udp int, err error) {
	protoIndex := -1
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if protoIndex < 0 {
			for i, field := range fields {
				if field == "PROTO" {
					protoIndex = i
					break
				}
			}
			if protoIndex < 0 {
				return 0, 0, fmt.Errorf("sockstat output has no PROTO header")
			}
			continue
		}
		if len(fields) <= protoIndex {
			return 0, 0, fmt.Errorf("malformed sockstat row: %q", line)
		}
		switch fields[protoIndex] {
		case "tcp4", "tcp6", "tcp46":
			tcp++
		case "udp4", "udp6", "udp46":
			udp++
		}
	}
	if protoIndex < 0 {
		return 0, 0, fmt.Errorf("empty sockstat output")
	}
	return tcp, udp, nil
}
