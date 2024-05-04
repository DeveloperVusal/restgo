package device

import (
	"fmt"
	"strings"
)

// This type for store information about browser
type browserInfo struct {
	name   string
	search string
}

// This type for store information about os
type osInfo struct {
	name   string
	search string
}

func DetectDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") {
		return "Mobile"
	} else if strings.Contains(ua, "tablet") {
		return "Tablet"
	} else if strings.Contains(ua, "bot") {
		return "Bot"
	} else {
		return "Desktop"
	}
}

func DetectOS(userAgent string) string {
	ua := strings.ToLower(userAgent)

	oses := []osInfo{
		{"Windows", "windows"},
		{"Linux", "linux"},
		{"macOS", "macintosh"},
		{"macOS", "mac os"},
		{"Android", "android"},
		{"MIUI (Xiaomi)", "xiaomi"},
		{"HyperOS", "hyperos"},
		{"iOS", "iphone"},
		{"iOS", "ipad"},
		{"iOS", "ipod"},
	}

	for _, o := range oses {
		if strings.Contains(ua, o.search) {
			versionStart := strings.Index(ua, o.search) + len(o.search)

			if versionStart >= len(ua) {
				return o.name
			}

			versionEnd := versionStart

			for versionEnd < len(ua) && (ua[versionEnd] >= '0' && ua[versionEnd] <= '9' || ua[versionEnd] == '.') {
				versionEnd++
			}

			if versionEnd > versionStart {
				version := ua[versionStart:versionEnd]
				return fmt.Sprintf("%s %s", o.name, version)
			} else {
				return o.name
			}
		}
	}

	return "Unknown OS"
}

func DetectBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	browsers := []browserInfo{
		{"Chrome", "chrome/"},
		{"Safari", "version/"},
		{"Firefox", "firefox/"},
		{"Opera", "opr/"},
		{"Edge", "edge/"},
		{"Brave", "brave/"},
		{"DuckDuckGo", "duckduckgo/"},
		{"Tor", "tor/"},
		{"Internet Explorer", "msie"},
		{"Internet Explorer", "trident"},
	}

	// Find browser in User-Agent
	for _, b := range browsers {
		if strings.Contains(ua, b.search) {
			versionStart := strings.Index(ua, b.search) + len(b.search)
			versionEnd := versionStart

			for versionEnd < len(ua) && (ua[versionEnd] >= '0' && ua[versionEnd] <= '9' || ua[versionEnd] == '.') {
				versionEnd++
			}

			version := ua[versionStart:versionEnd]

			return fmt.Sprintf("%s %s", b.name, version)
		}
	}

	return "Unknown Browser"
}
