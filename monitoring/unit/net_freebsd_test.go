//go:build freebsd

package monitoring

import (
	"net"
	"testing"
)

func TestCountSockstat(t *testing.T) {
	const header = "USER COMMAND PID FD PROTO LOCAL ADDRESS FOREIGN ADDRESS\n"
	for _, tc := range []struct {
		name, output string
		tcp, udp     int
		wantErr      bool
	}{
		{"mixed", header + "root sshd 123 3 tcp4 *:22 *:*\nroot sshd 123 4 tcp6 *:22 *:*\nroot ntpd 234 5 udp46 *:123 *:*\n", 2, 1, false},
		{"no sockets", header, 0, 0, false},
		{"missing header", "unexpected output", 0, 0, true},
		{"truncated", header + "root sshd 123\n", 0, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tcp, udp, err := countSockstat(tc.output)
			if (err != nil) != tc.wantErr || tcp != tc.tcp || udp != tc.udp {
				t.Fatalf("countSockstat = (%d, %d, %v), want (%d, %d, error=%v)", tcp, udp, err, tc.tcp, tc.udp, tc.wantErr)
			}
		})
	}
}

func TestFreeBSDMonitoring(t *testing.T) {
	cpu := CpuStaticInfo()
	if cpu.CPUCores < 1 || cpu.CPUName == "Unknown" {
		t.Fatalf("invalid CPU info: %+v", cpu)
	}
	if ram := Ram(); ram.Total == 0 || ram.Used > ram.Total {
		t.Fatalf("invalid memory info: %+v", ram)
	}
	previousRaw := flags.MemoryReportRawUsed
	flags.MemoryReportRawUsed = true
	raw := Ram()
	flags.MemoryReportRawUsed = previousRaw
	if raw.Total == 0 || raw.Used > raw.Total {
		t.Fatalf("invalid fallback for Linux memory mode on FreeBSD: %+v", raw)
	}
	if disk := Disk(); disk.Total == 0 || disk.Used > disk.Total {
		t.Fatalf("invalid disk info: %+v", disk)
	}
	if up, err := Uptime(); err != nil || up == 0 {
		t.Fatalf("invalid uptime: %d, %v", up, err)
	}
	if nics, err := InterfaceList(); err != nil || len(nics) == 0 {
		t.Fatalf("no network interfaces: %v, %v", nics, err)
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	udp, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer udp.Close()
	if tcpCount, udpCount, err := ConnectionsCount(); err != nil || tcpCount < 1 || udpCount < 1 {
		t.Fatalf("socket counters: tcp=%d udp=%d err=%v", tcpCount, udpCount, err)
	}
	t.Logf("CPU=%+v RAM=%+v disk=%+v OS=%s", cpu, Ram(), Disk(), OSName())
}
