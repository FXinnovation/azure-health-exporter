package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
)

// AvailabilityStatusIDSuffix is the common suffix of all the AvailabilityStatus IDs
const AvailabilityStatusIDSuffix = "/providers/Microsoft.ResourceHealth/availabilityStatuses/current"

// ResourceHealthClient is the client implementation to ResourceHealth API
type ResourceHealthClient struct {
	Session                *AzureSession
	Client                 *resourcehealth.AvailabilityStatusesClient
	LastRatelimitRemaining string
}

// ResourceHealth client interface
type ResourceHealth interface {
	GetAvailabilityStatus(resourceURI string) (*resourcehealth.AvailabilityStatus, error)
	GetAllAvailabilityStatuses() (*[]resourcehealth.AvailabilityStatus, error)
	GetSubscriptionID() string
	GetLastRatelimitRemaining() string
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

// GetLastRatelimitRemaining return last ratelimit remaining value
func (rc *ResourceHealthClient) GetLastRatelimitRemaining() string {
	return rc.LastRatelimitRemaining
}

// GetAllAvailabilityStatuses fetch all Resources Health availability statuses of the subscription
func (rc *ResourceHealthClient) GetAllAvailabilityStatuses() (*[]resourcehealth.AvailabilityStatus, error) {
	var asList []resourcehealth.AvailabilityStatus

	for it, err := rc.Client.ListBySubscriptionIDComplete(context.Background(), "", ""); it.NotDone(); err = it.Next() {
		if err != nil {
			return nil, err
		}
		asList = append(asList, it.Value())
		rc.LastRatelimitRemaining = it.Response().Header.Get("X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests")
	}
	return &asList, nil
}

// GetAvailabilityStatus fetch all Resources Health availability statuses of the subscription
func (rc *ResourceHealthClient) GetAvailabilityStatus(resourceURI string) (*resourcehealth.AvailabilityStatus, error) {
	as, err := rc.Client.GetByResource(context.Background(), resourceURI, "", "")
	if err != nil {
		return nil, err
	}
	rc.LastRatelimitRemaining = as.Response.Header.Get("X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests")

	return &as, nil
}
