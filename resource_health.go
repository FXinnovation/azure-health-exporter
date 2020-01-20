package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
)

// ResourceHealthClient is the client implementation to ResourceHealth API
type ResourceHealthClient struct {
	Session   *AzureSession
	Client    *resourcehealth.AvailabilityStatusesClient
	Resources Resources
}

// ResourceHealth client interface
type ResourceHealth interface {
	GetResourcesHealth(resourceTags map[string]string) (*[]resourcehealth.AvailabilityStatus, error)
	GetSubscriptionID() string
}

// NewResourceHealth returns a new ResourceHealth client
func NewResourceHealth(session *AzureSession) ResourceHealth {

	client := resourcehealth.NewAvailabilityStatusesClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer
	resources := NewResources(session)

	return &ResourceHealthClient{
		Session:   session,
		Client:    &client,
		Resources: resources,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (rc *ResourceHealthClient) GetSubscriptionID() string {
	return rc.Session.SubscriptionID
}

// GetResourcesHealth fetch Resources Health
func (rc *ResourceHealthClient) GetResourcesHealth(resourceTags map[string]string) (*[]resourcehealth.AvailabilityStatus, error) {
	var rhList []resourcehealth.AvailabilityStatus

	resources, err := rc.Resources.GetResources(resourceTags)
	if err != nil {
		return nil, err
	}

	for _, resource := range *resources {
		rh, err := rc.Client.GetByResource(context.Background(), *resource.ID, "", "")
		if err != nil {
			return nil, err
		}

		rhList = append(rhList, rh)
	}

	return &rhList, nil
}
