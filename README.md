# docker-credential-acr-helper

A [credential helper](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers) for the Docker daemon
that makes it easier to use [Alibaba Cloud Container Registry](https://www.alibabacloud.com/product/container-registry).

## Installation

Download the latest release from [Releases](https://github.com/mozillazg/docker-credential-acr-helper/releases) page.

## Configuration

### Credentials

If the `ALIBABA_CLOUD_ACCESS_KEY_ID` and `ALIBABA_CLOUD_ACCESS_KEY_SECRET` environment variables
are defined and are not empty, the program will use them to create the default credential. 
If not, the program loads and looks for the client in the configuration file(`~/.alibabacloud/credentials`).

For more information about configuring credentials, see [Provider](https://github.com/aliyun/credentials-go#provider).

### Docker

Place the `docker-credential-acr-helper` binary on your `PATH` and set the contents of your
`~/.docker/config.json` file to be:

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
see [Credential helpers](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers) in Docker Documentation.
