package acr

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var errUnknownDomain = errors.New("unknown domain")
var domainPattern = regexp.MustCompile(
	`^(?:(?P<instanceName>[^.\s]+)-)?registry(?:-intl)?(?:-vpc)?(?:-internal)?(?:\.distributed)?\.(?P<region>[^.]+)\.(?:cr\.)?aliyuncs\.com`)

const (
	urlPrefix      = "https://"
	hostNameSuffix = ".aliyuncs.com"
)

type Registry struct {
	IsEE         bool
	InstanceId   string
	InstanceName string
	Region       string
	Domain       string
}

func parseServerURL(rawURL string) (*Registry, error) {
	if !strings.Contains(rawURL, hostNameSuffix) {
		return nil, errUnknownDomain
	}
	if !strings.HasPrefix(rawURL, urlPrefix) {
		rawURL = urlPrefix + rawURL
	}
	serverURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	domain := serverURL.Hostname()
	if !strings.HasSuffix(domain, hostNameSuffix) {
		return nil, errUnknownDomain
	}

	subItems := domainPattern.FindStringSubmatch(domain)
	if len(subItems) != 3 {
		return nil, errUnknownDomain
	}
	instanceName := subItems[1]
	region := subItems[2]
	isEE := instanceName != ""

	return &Registry{
		IsEE:         isEE,
		InstanceId:   "",
		InstanceName: instanceName,
		Region:       region,
		Domain:       domain,
	}, nil
}
