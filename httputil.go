package misc

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func ReadBytes(req *http.Request) ([]byte, error) {
	return ioutil.ReadAll(req.Body)
}
func ReadContent(req *http.Request) (content string, err error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return content, err
	}
	return string(b), nil
}

func URLEncode(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}
	return u.Query().Encode()
}

func GetUrlPath(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}
	return u.Path
}

func GetRealIp(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if "" == ip || "unknown" == ip {
		ip = req.Header.Get("Proxy-Client-IP")
	}
	if "" == ip || "unknown" == ip {
		ip = req.Header.Get("WL-Proxy-Client-IP")
	}
	if "" == ip || "unknown" == ip {
		ip = req.Header.Get("HTTP_CLIENT_IP")
	}
	if "" == ip || "unknown" == ip {
		ip = req.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if "" == ip || "unknown" == ip {
		ip = req.RemoteAddr
	}
	return ip
}

func GetRefererDomain(req *http.Request) string {
	referer := req.Header.Get("Referer")
	u, _ := url.Parse(referer)
	if u != nil {
		return u.Host
	}
	return ""
}
