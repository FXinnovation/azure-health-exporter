package main

import (
	"errors"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// resource component positions in a resource ID
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

// CreateAllLabels creates label from Tags map and existing labels map
func CreateAllLabels(tags map[string]*string, resourceType *string, labels map[string]string) map[string]string {
	labels["resource_type"] = *resourceType
	for k, v := range tags {
		k = strings.ToLower(k)
		k = "tag_" + k
		k = invalidLabelChars.ReplaceAllString(k, "_")
		labels[k] = *v
	}

	return labels
}

// ExportAzureTagInfo exports azure_tag_info metric
func ExportAzureTagInfo(ch chan<- prometheus.Metric, tags map[string]*string, resourceType *string, labels map[string]string) {
	allLabels := CreateAllLabels(tags, resourceType, labels)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("azure_tag_info", "Tags of the Azure resource", nil, allLabels),
		prometheus.GaugeValue,
		1,
	)
}
