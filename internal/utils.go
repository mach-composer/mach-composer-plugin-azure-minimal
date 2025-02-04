package internal

import (
	"fmt"
)

// azureServicePlanResourceName Retrieve the resource name for an Azure app
// service plan.  The reason to make this conditional is because of backwards
// compatability; existing environments already have a `functionapp` resource.
// We want to keep that intact.
func azureServicePlanResourceName(value string) string {
	name := "functionapps"
	if value != "default" {
		name = fmt.Sprintf("functionapps_%s", value)
	}
	return name
}
