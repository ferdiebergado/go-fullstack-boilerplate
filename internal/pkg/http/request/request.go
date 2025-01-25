package request

import (
	"net"
	"net/http"
	"strings"
)

// getIPAddress extracts the client's IP address from the request.
func GetIPAddress(r *http.Request) string {
	// Get the IP from the X-REAL-IP header
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Get the IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")
	if len(splitIps) > 0 {
		// Remove whitespace and return the first IP
		ip = strings.TrimSpace(splitIps[0])
		return ip
	}

	// Fallback to the remote address. This will include the port number
	// that you will need to remove.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // If there is an error simply return the RemoteAddr
	}

	return ip
}
