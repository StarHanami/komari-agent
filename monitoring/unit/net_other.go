//go:build !freebsd

package monitoring

func platformConnectionsCount() (int, int, error) {
	return gopsutilConnectionsCount()
}
