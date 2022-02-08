package registry

import (
	"github.com/containrrr/watchtower/pkg/registry/helpers"
	watchtowerTypes "github.com/containrrr/watchtower/pkg/types"
	ref "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// GetPullOptions creates a struct with all options needed for pulling images from a registry
func GetPullOptions(container watchtowerTypes.Container) (types.ImagePullOptions, error) {
	imageName := container.ImageName()
	registryUsername := container.RegistryUsername()
	registryPassword := container.RegistryPassword()

	auth, err := EncodedAuth(imageName, registryUsername, registryPassword)
	log.Debugf("Got image name: %s", imageName)
	if err != nil {
		return types.ImagePullOptions{}, err
	}

	if auth == "" {
		return types.ImagePullOptions{}, nil
	}
	log.Tracef("Got auth value: %s", auth)

	return types.ImagePullOptions{
		RegistryAuth:  auth,
		PrivilegeFunc: DefaultAuthHandler,
	}, nil
}

// DefaultAuthHandler will be invoked if an AuthConfig is rejected
// It could be used to return a new value for the "X-Registry-Auth" authentication header,
// but there's no point trying again with the same value as used in AuthConfig
func DefaultAuthHandler() (string, error) {
	log.Debug("Authentication request was rejected. Trying again without authentication")
	return "", nil
}

// WarnOnAPIConsumption will return true if the registry is known-expected
// to respond well to HTTP HEAD in checking the container digest -- or if there
// are problems parsing the container hostname.
// Will return false if behavior for container is unknown.
func WarnOnAPIConsumption(container watchtowerTypes.Container) bool {

	normalizedName, err := ref.ParseNormalizedNamed(container.ImageName())
	if err != nil {
		return true
	}

	containerHost, err := helpers.NormalizeRegistry(normalizedName.String())
	if err != nil {
		return true
	}

	if containerHost == "index.docker.io" || containerHost == "ghcr.io" {
		return true
	}

	return false
}
