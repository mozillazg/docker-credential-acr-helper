package main

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/mozillazg/docker-credential-acr-helper/pkg/credhelper"
)

func main() {
	kc := authn.NewMultiKeychain(
		authn.DefaultKeychain,
		authn.NewKeychainFromHelper(credhelper.NewACRHelper()),
	)
	ref := os.Getenv("REPO_URL")
	digest, err := crane.Digest(ref, crane.WithAuthFromKeychain(kc))
	if err != nil {
		panic(err)
	}
	fmt.Printf("got digest for %q:\n%s\n", ref, digest)
}
