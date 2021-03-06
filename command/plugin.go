package command

import (
	"github.com/proffer/components"
	awscopyamiresource "github.com/proffer/resources/aws/copyami"
	awsshareamiresource "github.com/proffer/resources/aws/shareami"
)

// Resources represents the map of available resources.
var Resources = map[string]components.Resourcer{
	"aws-copyami":  new(awscopyamiresource.Resource),
	"aws-shareami": new(awsshareamiresource.Resource),
}
