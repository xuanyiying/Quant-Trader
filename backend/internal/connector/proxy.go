package connector

import (
	"net/url"
	"os"
	"strings"
)

func getProxyURL() *url.URL {
	proxyAddrs := []string{
		os.Getenv("HTTPS_PROXY"),
		os.Getenv("HTTP_PROXY"),
		os.Getenv("SOCKS5_PROXY"),
	}
	for _, addr := range proxyAddrs {
		if addr == "" {
			continue
		}
		if !strings.Contains(addr, "://") {
			addr = "socks5://" + addr
		}
		proxyURL, err := url.Parse(addr)
		if err == nil {
			return proxyURL
		}
	}
	return nil
}
