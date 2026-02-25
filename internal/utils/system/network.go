// =============================================================================
// FILE: internal/utils/system/network.go
// PURPOSE: Network utility functions. Provides helpers for proxy detection,
//          local IP resolution, and network interface enumeration. Ports
//          Python utils/system/network.py.
// =============================================================================

package system

import (
	"fmt"
	"net"
	"os"
	"time"
)

// ---------------------------------------------------------------------------
// Network utilities
// ---------------------------------------------------------------------------

// LocalIP returns the local machine's preferred outbound IP address.
//
// Returns:
//   - The local IP as a string, or an error.
func LocalIP() (string, error) {
	conn, err := net.DialTimeout("udp", "8.8.8.8:80", 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to determine local IP: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GetProxyFromEnv reads proxy configuration from environment variables.
// Checks HTTP_PROXY, HTTPS_PROXY, and ALL_PROXY in order.
//
// Returns:
//   - The proxy URL string, or empty if no proxy is configured.
func GetProxyFromEnv() string {
	for _, key := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy", "ALL_PROXY", "all_proxy"} {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return ""
}

// IsPortOpen checks if a TCP port is open on the given host.
//
// Parameters:
//   - host: The host to check (e.g. "localhost").
//   - port: The port number.
//   - timeout: Maximum wait time.
//
// Returns:
//   - true if the port is open and accepting connections.
func IsPortOpen(host string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ListInterfaces returns the names of all active (up) network interfaces.
//
// Returns:
//   - Slice of interface names, and any error.
func ListInterfaces() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %w", err)
	}

	var names []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 {
			names = append(names, iface.Name)
		}
	}
	return names, nil
}
