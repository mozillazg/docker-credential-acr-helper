# docker-credential-acr-helper

A [credential helper](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers) for the Docker daemon
that makes it easier to use [Alibaba Cloud Container Registry (ACR)](https://www.alibabacloud.com/product/container-registry).

## Installation

Download the latest release from the [Releases](https://github.com/mozillazg/docker-credential-acr-helper/releases) page.

## Configuration

### ACR Credentials

By default, the helper searches for ACR credentials in the following order:

1. It fetches the credentials via [RAM Roles for Service Accounts (RRSA) OIDC Token](https://www.alibabacloud.com/help/en/container-service-for-kubernetes/latest/use-rrsa-to-enforce-access-control)
   when the `ALIBABA_CLOUD_ROLE_ARN`, `ALIBABA_CLOUD_OIDC_PROVIDER_ARN`, and
   `ALIBABA_CLOUD_OIDC_TOKEN_FILE` environment variables are defined and are not empty.
2. Use access key id and access key secret that are specified by the `ALIBABA_CLOUD_ACCESS_KEY_ID` and
   `ALIBABA_CLOUD_ACCESS_KEY_SECRET` environment variables.
3. A profile file whose path is specified by the `ALIBABA_CLOUD_CREDENTIALS_FILE` environment variable.
4. A profile file in a default location:
   * On Windows, this is `C:\Users\USER_NAME\.alibabacloud\credentials`.
   * On other systems, it is `~/.alibabacloud/credentials`.
5. It fetches the credentials of the RAM Role associated with the VM from the metadata server when
   the environment variable `ALIBABA_CLOUD_ECS_METADATA` is defined and not empty.

For more information about configuring credentials, see [Provider](https://github.com/aliyun/credentials-go#provider)
in the @aliyun/credentials-go.

### Docker

Place the `docker-credential-acr-helper` binary on your `PATH` and
add a `credHelpers` entry to the Docker config file (`~/.docker/config.json`)
for each ACR registry that you care about.
Keys specify the registry domain (**without** the `https://`), and values specify the suffix of the credential helper binary (everything after `docker-credential-`). 
For example:

```
{
	"credHelpers": {
		"registry.cn-beijing.aliyuncs.com": "acr-helper",
		"registry-intl.ap-southeast-1.aliyuncs.com": "acr-helper",
		"registry.<region>.aliyuncs.com": "acr-helper",
		"<acr-ee-instance-name>-registry.<region>.cr.aliyuncs.com": "acr-helper"
	}
}
```

For more information about configuring Docker,
see [Credential helpers](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers) in the Docker Documentation.
