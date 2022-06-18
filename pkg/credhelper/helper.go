package credhelper

import (
	"errors"
	"fmt"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/mozillazg/docker-credential-acr-helper/pkg/acr"
	"github.com/sirupsen/logrus"
)

var errNotImplemented = errors.New("not implemented")

type ACRHelper struct {
	client *acr.Client
}

func NewACRHelper() *ACRHelper {
	return &ACRHelper{client: &acr.Client{}}
}

func (a *ACRHelper) Get(serverURL string) (string, string, error) {
	// TODO: add cache
	cred, err := a.client.GetCredentials(serverURL)
	if err != nil {
		logrus.WithField("serverURL", serverURL).WithError(err).Error("get credentials failed")
		return "", "", fmt.Errorf("get credentials for %q failed: %+v", serverURL, err)
	}
	return cred.UserName, cred.Password, nil
}

func (a *ACRHelper) Add(creds *credentials.Credentials) error {
	return errNotImplemented
}

func (a *ACRHelper) Delete(serverURL string) error {
	return errNotImplemented
}

func (a *ACRHelper) List() (map[string]string, error) {
	return nil, errNotImplemented
}
