package device

import (
	"fmt"
	"net"
	"pismo-dev/internal/types"
	"pismo-dev/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

// Parse the IP address
func parseIPAddress(c *gin.Context) string {
	var ip string

	if xForwardedFor := c.Request.Header.Get("x-forwarded-for"); xForwardedFor != "" {
		ip = xForwardedFor
	} else if xRealIP := c.Request.Header.Get("x-real-ip"); xRealIP != "" {
		ip = xRealIP
	} else {
		ip, _, _ = net.SplitHostPort(c.Request.RemoteAddr)
	}

	parts := strings.Split(ip, ",")
	ipAddress := strings.TrimSpace(parts[len(parts)-1])

	return ipAddress
}

// Detect the device details
func detectDevice(ua *user_agent.UserAgent, xSource string) (string, string, string, string) {
	browser, version := ua.Browser()
	os := ua.OS()
	platform := ua.Platform()
	return browser, version, os, platform
}

// Build the device
func buildDevice(browser string, version string, os string, platform string) string {
	return browser + " " + version + " " + os + " " + platform
}

func SetDevice(c *gin.Context) (types.DeviceDetailInterface, error) {
	headers := c.Request.Header
	if headers == nil {
		errMsg := "Invalid request headers"
		logger.Error(errMsg)
		return types.DeviceDetailInterface{}, fmt.Errorf(errMsg)
	}
	skipFurther := false
	agentDetailString := c.Request.UserAgent()
	userDeviceDetails := DetectOldDevice(agentDetailString, headers.Get("x-source"))
	builtDeviceDetails := buildDeviceFromDetails(userDeviceDetails)
	if builtDeviceDetails != "" {
		skipFurther = true
	}
	if !skipFurther {
		userDeviceDetails = DetectNewDevice(agentDetailString, headers.Get("x-source"))
		builtDeviceDetails = buildDeviceFromDetails(userDeviceDetails)
		if builtDeviceDetails != "" {
			skipFurther = true
		}
	}
	if !skipFurther {
		userAgent := user_agent.New(agentDetailString)
		if userAgent == nil {
			errMsg := "User Agent not detected"
			logger.Error(errMsg)
			return types.DeviceDetailInterface{}, fmt.Errorf(errMsg)
		}

		browser, version, os, platform := detectDevice(userAgent, headers.Get("x-source"))
		if browser == "" || version == "" || os == "" || platform == "" {
			errMsg := "No details of device found"
			logger.Error(errMsg, map[string]interface{}{"context": c, "userAgent": agentDetailString})
			return types.DeviceDetailInterface{}, fmt.Errorf(errMsg)
		}
		userDeviceDetails = types.UserAgentDetails{
			Browser:  browser,
			Version:  version,
			Os:       os,
			Platform: platform,
		}
		builtDeviceDetails = buildDeviceFromDetails(userDeviceDetails)
	}

	ipAddress := parseIPAddress(c)
	if ipAddress == "" {
		errMsg := "Ip address not found"
		logger.Error(errMsg, map[string]interface{}{"context": c, "forward": c.Request.Header.Get("x-forwarded-for"), "realIp": c.Request.Header.Get("x-real-ip")})
		return types.DeviceDetailInterface{}, fmt.Errorf(errMsg)
	}

	deviceDetails := types.DeviceDetailInterface{
		IPAddress: ipAddress,
		Device:    builtDeviceDetails,
		Source:    buildSource(headers.Get("x-source")),
	}
	return deviceDetails, nil

}

func buildSource(source string) string {
	if source == "app" {
		source = "android"
	}
	return source
}
