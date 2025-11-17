package middleware

import (
	"fmt"
	"net"
	"net/url"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthParams struct {
	IP      string
	Referer string
	Domain  string
}

const ContextAuthParamsKey = "authParams"

func isIPInList(clientIP string, ipList string) bool {
	if ipList == "" {
		return false
	}

	ipAddr := net.ParseIP(clientIP)
	if ipAddr == nil {
		return false
	}

	ips := strings.Split(ipList, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}

		if strings.Contains(ip, "/") {
			_, ipNet, err := net.ParseCIDR(ip)
			if err == nil && ipNet.Contains(ipAddr) {
				return true
			}
		} else {
			if ip == clientIP {
				return true
			}
		}
	}

	return false
}

func isDomainInList(domain string, domainList string) bool {
	if domainList == "" || domain == "" {
		return false
	}

	domains := strings.Split(domainList, ",")
	for _, d := range domains {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}

		if strings.HasPrefix(d, "*.") {
			suffix := d[1:]
			if strings.HasSuffix(domain, suffix) {
				return true
			}
		} else {
			if domain == d {
				return true
			}
		}
	}

	return false
}

func extractDomainFromReferer(referer string) string {
	if referer == "" {
		return ""
	}

	u, err := url.Parse(referer)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

func isFromBaseUrl(clientIP string, domain string) bool {
	baseUrl := utils.GetBaseUrl()
	if baseUrl == "" {
		return true
	}

	baseUrlObj, err := url.Parse(baseUrl)
	if err != nil {
		return true
	}

	baseUrlDomain := baseUrlObj.Hostname()

	if domain != "" && domain == baseUrlDomain {
		return true
	}

	if net.ParseIP(baseUrlDomain) != nil && baseUrlDomain == clientIP {
		return true
	}

	ips, err := net.LookupIP(baseUrlDomain)
	if err == nil {
		for _, ip := range ips {
			if ip.To4() != nil && ip.String() == clientIP {
				return true
			}
		}
	}

	return false
}

/* IpRefererMiddleware 获取用户IP和Referer信息并存储到上下文中，并进行IP白黑名单检测 */
func IpRefererMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if clientIP == "" || clientIP == "::1" || clientIP == "127.0.0.1" {
			forwardedIPs := c.GetHeader("X-Forwarded-For")
			if forwardedIPs != "" {
				ips := strings.Split(forwardedIPs, ",")
				if len(ips) > 0 {
					clientIP = strings.TrimSpace(ips[0])
				}
			}

			if clientIP == "" || clientIP == "::1" || clientIP == "127.0.0.1" {
				realIP := c.GetHeader("X-Real-IP")
				if realIP != "" {
					clientIP = realIP
				}
			}
		}

		referer := c.GetHeader("Referer")

		domain := extractDomainFromReferer(referer)

		authParams := &AuthParams{
			IP:      clientIP,
			Referer: referer,
			Domain:  domain,
		}

		c.Set(ContextAuthParamsKey, authParams)
		if isFromBaseUrl(clientIP, domain) {
			c.Next()
			return
		}

		securitySettings, err := setting.GetSettingsByGroupAsMap("security")
		if err == nil {
			var ipWhitelist string
			if val, ok := securitySettings.Settings["ip_whitelist"]; ok {
				if whitelistStr, ok := val.(string); ok {
					ipWhitelist = whitelistStr
				}
			}

			var ipBlacklist string
			if val, ok := securitySettings.Settings["ip_blacklist"]; ok {
				if blacklistStr, ok := val.(string); ok {
					ipBlacklist = blacklistStr
				}
			}

			var domainWhitelist string
			if val, ok := securitySettings.Settings["domain_whitelist"]; ok {
				if whitelistStr, ok := val.(string); ok {
					domainWhitelist = whitelistStr
				}
			}

			var domainBlacklist string
			if val, ok := securitySettings.Settings["domain_blacklist"]; ok {
				if blacklistStr, ok := val.(string); ok {
					domainBlacklist = blacklistStr
				}
			}

			if ipWhitelist != "" && !isIPInList(clientIP, ipWhitelist) {
				errorMessage := fmt.Sprintf("您的IP(%s)不在访问白名单中", clientIP)
				err := errors.New(errors.CodeIPNotInWhitelist, errorMessage)
				c.JSON(errors.HTTPStatus(err), gin.H{
					"code":    int(errors.CodeIPNotInWhitelist),
					"message": err.Message,
					"ip":      clientIP,
				})
				c.Abort()
				return
			}

			if ipBlacklist != "" && isIPInList(clientIP, ipBlacklist) {
				errorMessage := fmt.Sprintf("您的IP(%s)已被列入黑名单", clientIP)
				err := errors.New(errors.CodeIPInBlacklist, errorMessage)
				c.JSON(errors.HTTPStatus(err), gin.H{
					"code":    int(errors.CodeIPInBlacklist),
					"message": err.Message,
					"ip":      clientIP,
				})
				c.Abort()
				return
			}

			if domain != "" && domainWhitelist != "" && !isDomainInList(domain, domainWhitelist) {
				errorMessage := fmt.Sprintf("您的域名(%s)不在访问白名单中", domain)
				err := errors.New(errors.CodeDomainNotInWhitelist, errorMessage)
				c.JSON(errors.HTTPStatus(err), gin.H{
					"code":    int(errors.CodeDomainNotInWhitelist),
					"message": err.Message,
					"domain":  domain,
				})
				c.Abort()
				return
			}

			if domain != "" && domainBlacklist != "" && isDomainInList(domain, domainBlacklist) {
				errorMessage := fmt.Sprintf("您的域名(%s)已被列入黑名单", domain)
				err := errors.New(errors.CodeDomainInBlacklist, errorMessage)
				c.JSON(errors.HTTPStatus(err), gin.H{
					"code":    int(errors.CodeDomainInBlacklist),
					"message": err.Message,
					"domain":  domain,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func GetAuthParams(c *gin.Context) *AuthParams {
	value, exists := c.Get(ContextAuthParamsKey)
	if !exists {
		return &AuthParams{}
	}

	if params, ok := value.(*AuthParams); ok {
		return params
	}

	return &AuthParams{}
}
