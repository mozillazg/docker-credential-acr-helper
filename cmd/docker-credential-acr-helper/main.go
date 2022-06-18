package main

import (
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/mozillazg/docker-credential-acr-helper/pkg/credhelper"
)

func main() {
	credentials.Serve(credhelper.NewACRHelper())
}
