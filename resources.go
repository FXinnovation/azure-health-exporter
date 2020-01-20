package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
)

// ResourcesClient is the client implementation to VirtualMachines API
type ResourcesClient struct {
	Session *AzureSession
	Client  *resources.Client
}

// Resources client interface
type Resources interface {
	GetResources(resourceTags map[string]string) (*[]resources.GenericResource, error)
}

// NewResources returns a new Resources client
func NewResources(session *AzureSession) Resources {
	client := resources.NewClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	return &ResourcesClient{
		Session: session,
		Client:  &client,
	}
}

// GetResources return resources by tags (if tags map is not empty).
// A resource must match all tag parameters in order to be fetched
func (rc *ResourcesClient) GetResources(resourceTags map[string]string) (*[]resources.GenericResource, error) {
	// TODO : filter by tag here
	filter := fmt.Sprintf("")
	resList, err := rc.list(filter)
	if err != nil {
		return nil, err
	}

	// Filtering by tag is done manually, as Azure does not support
	// to filter both by resource type and by tag name/value
	var filteredList []resources.GenericResource
	for _, resource := range *resList {
		include := true
		if resourceTags != nil {
			for name, value := range resourceTags {
				if *resource.Tags[name] != value {
					include = false
					break
				}
			}
		}
		if include {
			filteredList = append(filteredList, resource)
		}
	}

	return &filteredList, nil
}

func (rc *ResourcesClient) list(filter string) (*[]resources.GenericResource, error) {
	var resourceList []resources.GenericResource

	for it, err := rc.Client.ListComplete(context.Background(), filter, "", nil); it.NotDone(); err = it.Next() {
		if err != nil {
			return nil, err
		}
		resourceList = append(resourceList, it.Value())
	}

	return &resourceList, nil
}
