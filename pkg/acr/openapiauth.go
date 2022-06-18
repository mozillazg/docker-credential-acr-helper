package acr

import (
	"os"
	"path/filepath"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/mozillazg/docker-credential-acr-helper/pkg/version"
)

const (
	envRoleArn       = "ALIBABA_CLOUD_ROLE_ARN"
	envOidcArn       = "ALIBABA_CLOUD_OIDC_PROVIDER_ARN"
	envOidcTokenFile = "ALIBABA_CLOUD_OIDC_TOKEN_FILE"
)

var defaultProfilePath = filepath.Join("~", ".alibabacloud", "credentials")

func getOpenapiAuth() (credentials.Credential, error) {
	path, err := expandPath(defaultProfilePath)
	if err == nil {
		_ = os.Setenv(credentials.ENVCredentialFile, path)
	}
	var conf *credentials.Config

	roleArn := os.Getenv(envRoleArn)
	oidcArn := os.Getenv(envOidcArn)
	tokenFile := os.Getenv(envOidcTokenFile)
	if roleArn != "" && oidcArn != "" && tokenFile != "" {
		conf = new(credentials.Config).
			SetType("oidc_role_arn").
			SetOIDCProviderArn(oidcArn).
			SetOIDCTokenFilePath(tokenFile).
			SetRoleArn(roleArn).
			SetRoleSessionName(version.ProjectName)
	}

	cred, err := credentials.NewCredential(conf)
	return cred, err
}

func expandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}
	return path, nil
}
