package device

import (
	"fmt"
	"pismo-dev/internal/types"
	"regexp"
	"strconv"
	"strings"

	"github.com/ua-parser/uap-go/uaparser"
)

func DetectOldDevice(ua string, source string) types.UserAgentDetails {
	parser := uaparser.NewFromSaved()
	client := parser.Parse(ua)
	uaDetails := types.UserAgentDetails{}
	uaDetails.Os = ""
	uaDetails.OsVersion = ""
	uaDetails.DeviceType = ""
	// Skip user agent version for non-web source
	if source != "web" {
		uaDetails.Version = ""
	}
	if client.UserAgent.Family != "" && client.Device.Family != "" && client.Os.Family != "" {
		if source == "web" && client.UserAgent.Family == "Android Browser" {
			uaDetails.Browser = ""
		} else {
			uaDetails.Browser = strings.TrimSpace(strings.Replace(client.UserAgent.Family, "Webview", "", 1))
		}
		uaDetails.Os = strings.TrimSpace(client.Os.Family)
		uaDetails.OsVersion = strings.TrimSpace(client.Os.ToVersionString())
		if client.Device.Brand != "" && client.Device.Model != "" {
			uaDetails.DeviceType = fmt.Sprintf("%s %s", client.Device.Brand, client.Device.Model)
			if client.Device.Model == "iPhone" {
				re := regexp.MustCompile(`iPhone/(\d+)`)
				match := re.FindStringSubmatch(ua)
				if len(match) > 1 {
					uaDetails.DeviceType = fmt.Sprintf("%s %s", uaDetails.DeviceType, match[1])
				}
			}
			uaDetails.Browser = strings.TrimSpace(strings.Replace(uaDetails.Browser, "Android Browser", "", 1))
		}
	}
	if uaDetails.Browser == "" || uaDetails.Browser == "unknown" {
		uaDetails.Browser = ""
	}
	if uaDetails.Version == "" || uaDetails.Version == "unknown" {
		uaDetails.Version = strings.TrimSpace(client.UserAgent.ToVersionString())
		if uaDetails.Version == "" || uaDetails.Version == "unknown" {
			uaDetails.Version = ""
		}
	}
	if uaDetails.Os == "" || uaDetails.Os == "unknown" {
		uaDetails.Os = ""
	}
	if uaDetails.OsVersion == "" || uaDetails.OsVersion == "unknown" {
		uaDetails.OsVersion = ""
	}
	if uaDetails.DeviceType == "" || uaDetails.DeviceType == "unknown" || uaDetails.DeviceType == "unknown unknown" {
		uaDetails.DeviceType = ""
	}
	return uaDetails
}
func SafeMinMax(array []string) (int, int, error) {
	if array == nil || len(array) < 1 {
		return 0, 0, fmt.Errorf("Invalid array")
	}
	min, err := strconv.Atoi(array[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid array")
	}
	max := min
	for i := 1; i < len(array); i++ {
		value, err := strconv.Atoi(array[i])
		if err != nil {
			return 0, 0, fmt.Errorf("Invalid array")
		}
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max, nil
}

func DetectNewDevice(ua string, source string) types.UserAgentDetails {
	parser := uaparser.NewFromSaved()
	client := parser.Parse(ua)
	uaDetails := types.UserAgentDetails{}
	uaDetails.Browser = ""
	uaDetails.Version = ""
	uaDetails.Os = ""
	uaDetails.OsVersion = ""
	uaDetails.DeviceType = ""
	// Skip user agent version for non-web source
	if source != "web" {
		uaDetails.Version = ""
	}
	if client != nil {
		uaDetails.Browser = strings.TrimSpace(strings.Replace(client.UserAgent.Family, "Webview", "", 1))
		uaDetails.Version = strings.TrimSpace(client.UserAgent.ToVersionString())
		uaDetails.Os = strings.TrimSpace(client.Os.Family)
		uaDetails.OsVersion = strings.TrimSpace(client.Os.ToVersionString())
		uaDetails.DeviceType = fmt.Sprintf("%s %s", strings.TrimSpace(client.Device.Brand), strings.TrimSpace(client.Device.Model))
		if client.Device.Model == "iPhone" {
			re := regexp.MustCompile(`iPhone/(\d+)`)
			match := re.FindStringSubmatch(ua)
			if len(match) > 1 {
				uaDetails.DeviceType = fmt.Sprintf("%s %s", uaDetails.DeviceType, match[1])
			}
		}
		uaDetails.Browser = strings.TrimSpace(strings.Replace(uaDetails.Browser, "Android Browser", "", 1))
	}
	if uaDetails.Browser == "" || uaDetails.Browser == "unknown" {
		uaDetails.Browser = ""
	}
	if uaDetails.Version == "" || uaDetails.Version == "unknown" {
		uaDetails.Version = ""
	}
	if uaDetails.Os == "" || uaDetails.Os == "unknown" {
		uaDetails.Os = ""
	}
	if uaDetails.OsVersion == "" || uaDetails.OsVersion == "unknown" {
		uaDetails.OsVersion = ""
	}
	if uaDetails.DeviceType == "" || uaDetails.DeviceType == "unknown" || uaDetails.DeviceType == "unknown unknown" {
		uaDetails.DeviceType = ""
	}
	return uaDetails
}

func buildDeviceFromDetails(ua types.UserAgentDetails) string {
	_array := []string{}

	parts := strings.Split(ua.Version, ".")
	if len(parts) > 2 {
		ua.Version = strings.Join(parts[:2], ".")
	}

	browser := strings.TrimSpace(ua.Browser + " " + ua.Version)
	deviceType := strings.TrimSpace(ua.DeviceType)
	platform := strings.TrimSpace(ua.Os + " " + ua.OsVersion)

	if len(browser) > 0 {
		_array = append(_array, browser)
	}
	if len(deviceType) > 0 {
		_array = append(_array, deviceType)
	}
	if len(platform) > 0 {
		_array = append(_array, platform)
	}

	return strings.TrimSpace(strings.Join(_array, " - "))
}
