package main

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
)

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

// NewAzureSession create a new Azure session
func NewAzureSession(subscriptionID string) (*AzureSession, error) {

	if subscriptionID == "" {
		return nil, errors.New("Invalid subscription ID")
	}

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize authorizer")
	}

	session := AzureSession{
		SubscriptionID: subscriptionID,
		Authorizer:     authorizer,
	}

	return &session, nil
}
