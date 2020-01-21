package main

import (
	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
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

	// TODO Get Resources (by tag eventually), and compare lists by lowercased ResourceGroup.ResourceName.SubResourceName
	// pass the 2 list to CollectAvailabilityUp: creates logic from asList and labels from Resources

	asList, err := c.resourceHealth.GetAllAvailabilityStatuses()
	if err != nil {
		log.Errorf("Failed to get availability statuses list: %v", err)
		ch <- prometheus.NewInvalidMetric(azureErrorDesc, err)
		return
	}

	c.CollectAvailabilityUp(ch, asList)
}

// CollectAvailabilityUp converts Resource Health Availability as an UP metric
func (c *ResourceHealthCollector) CollectAvailabilityUp(ch chan<- prometheus.Metric, asList *[]resourcehealth.AvailabilityStatus) {

	for _, as := range *asList {
		up := 1.0
		if as.Properties.AvailabilityState == resourcehealth.Unavailable {
			up = 0
		}

		labels, err := ParseResourceID(*as.ID)
		if err != nil {
			log.Errorf("Skipping availability statuses: %s", err)
			continue
		}

		labels["subscription_id"] = c.resourceHealth.GetSubscriptionID()

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("resource_health_availability_up", "Resource health availability that relies on signals from different Azure services to assess whether a resource is healthy", nil, labels),
			prometheus.GaugeValue,
			up,
		)

		// TODO add optional azure_tag_info
	}
}
