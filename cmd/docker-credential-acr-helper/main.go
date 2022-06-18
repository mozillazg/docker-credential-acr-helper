package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/mozillazg/docker-credential-acr-helper/pkg/credhelper"
	"github.com/mozillazg/docker-credential-acr-helper/pkg/version"
)

const versionInfo = `docker-credential-acr-helper
Version:    %s
Git commit: %s
`

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, "print version and exit")
	flag.Parse()
	if versionFlag {
		fmt.Printf(versionInfo, version.Version, version.GitCommit)
		os.Exit(0)
	}

	credentials.Serve(credhelper.NewACRHelper())
}
