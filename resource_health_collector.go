package main

import (
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// ResourceHealthCollector collect ResourceHealth metrics
type ResourceHealthCollector struct {
	resourceHealth ResourceHealth
	resources      Resources
}

// NewResourceHealthCollector returns the collector
func NewResourceHealthCollector(session *AzureSession) *ResourceHealthCollector {
	resourceHealth := NewResourceHealth(session)
	resources := NewResources(session)

	return &ResourceHealthCollector{
		resourceHealth: resourceHealth,
		resources:      resources,
	}
}

// Describe to satisfy the collector interface.
func (c *ResourceHealthCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("ResourceHealthCollector", "dummy", nil, nil)
}

// Collect metrics from Resource Health API
func (c *ResourceHealthCollector) Collect(ch chan<- prometheus.Metric) {

	// In order to avoid the very low resource health API rate limit,
	// all availability statuses are fetched in 1 query and then parsed to lookup configured resources
	asList, err := c.resourceHealth.GetAllAvailabilityStatuses()
	if err != nil {
		log.Errorf("Failed to get all availability status: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	for _, resourceConfiguration := range config.ResourceConfigurations {
		for _, resourceType := range resourceConfiguration.ResourceTypes {
			resourceList, err := c.resources.GetResources(resourceType, resourceConfiguration.ResourceTags)
			if err != nil {
				log.Errorf("Failed to get resource list: %v", err)
				ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
				return
			}

			for _, resource := range *resourceList {
				for _, as := range *asList {
					if strings.ToLower(*as.ID) == strings.ToLower(*resource.ID+AvailabilityStatusIDSuffix) {
						c.CollectAvailabilityUp(ch, &as, &resource)
					}
				}
			}
		}
	}

	c.CollectRateLimitRemaining(ch)
}

// CollectAvailabilityUp converts Resource Health Availability status as an UP metric
func (c *ResourceHealthCollector) CollectAvailabilityUp(ch chan<- prometheus.Metric, as *resourcehealth.AvailabilityStatus,
	resource *resources.GenericResource) {

	// Only the `Unavailable` status can be used with confidence to consider availability "down"
	up := 1.0
	if as.Properties.AvailabilityState == resourcehealth.Unavailable {
		up = 0
	}

	labels, err := ParseResourceID(*resource.ID)
	if err != nil {
		log.Errorf("Failed to parse resource ID: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	labels["subscription_id"] = c.resourceHealth.GetSubscriptionID()
	labels["resource_type"] = *resource.Type

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("azure_resource_health_availability_up", "Resource health availability that relies on signals from different Azure services to assess whether a resource is healthy", nil, labels),
		prometheus.GaugeValue,
		up,
	)

	if config.ExposeAzureTagInfo {
		ExportAzureTagInfo(ch, resource.Tags, resource.Type, labels)
	}
}

// CollectRateLimitRemaining converts X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests header as metric
func (c *ResourceHealthCollector) CollectRateLimitRemaining(ch chan<- prometheus.Metric) {

	labels := make(map[string]string)
	labels["subscription_id"] = c.resourceHealth.GetSubscriptionID()

	ratelimitRemaining, err := strconv.ParseFloat(c.resourceHealth.GetLastRatelimitRemaining(), 64)
	if err != nil {
		log.Errorf("Failed to parse ratelimit remaining: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("azure_resource_health_ratelimit_remaining_count", "Azure subscription scoped Resource Health requests remaining (based on X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests header)", nil, labels),
		prometheus.GaugeValue,
		ratelimitRemaining,
	)

}
