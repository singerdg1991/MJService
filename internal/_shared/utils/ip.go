package utils

import (
	"net"
	"net/http"
)

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip, nil
	}
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip, nil
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}
