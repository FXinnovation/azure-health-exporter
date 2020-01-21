package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
)

// ResourceHealthClient is the client implementation to ResourceHealth API
type ResourceHealthClient struct {
	Session *AzureSession
	Client  *resourcehealth.AvailabilityStatusesClient
}

// ResourceHealth client interface
type ResourceHealth interface {
	GetAllAvailabilityStatuses() (*[]resourcehealth.AvailabilityStatus, error)
	GetSubscriptionID() string
}

// NewResourceHealth returns a new ResourceHealth client
func NewResourceHealth(session *AzureSession) ResourceHealth {

	client := resourcehealth.NewAvailabilityStatusesClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	return &ResourceHealthClient{
		Session: session,
		Client:  &client,
	}
}

// GetSubscriptionID return the client's Subscription ID
func (rc *ResourceHealthClient) GetSubscriptionID() string {
	return rc.Session.SubscriptionID
}

// GetAllAvailabilityStatuses fetch all Resources Health availability statuses of the subscription
func (rc *ResourceHealthClient) GetAllAvailabilityStatuses() (*[]resourcehealth.AvailabilityStatus, error) {
	var asList []resourcehealth.AvailabilityStatus
	for it, err := rc.Client.ListBySubscriptionIDComplete(context.Background(), "", ""); it.NotDone(); err = it.Next() {
		if err != nil {
			return nil, err
		}
		asList = append(asList, it.Value())
	}
	return &asList, nil
}
