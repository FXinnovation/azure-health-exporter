package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// resource component positions in an AvailabilityStatus ID
	resourceGroupPosition   = 4
	resourceNamePosition    = 8
	subResourceNamePosition = 10

	invalidLabelChars = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

// ParseResourceID - Returns resource info from a given resource ID.
func ParseResourceID(resourceID string) (map[string]string, error) {
	info := make(map[string]string)
	resource := strings.Split(resourceID, "/")

	if len(resource) < resourceNamePosition+1 {
		return nil, errors.New("Error parsing resource ID, expected pattern is not matched for " + resourceID)
	}

	info["resource_group"] = resource[resourceGroupPosition]
	info["resource_name"] = resource[resourceNamePosition]
	if len(resource) > subResourceNamePosition && resource[resourceNamePosition+1] != "providers" {
		info["sub_resource_name"] = resource[subResourceNamePosition]
	}
	return info, nil
}
