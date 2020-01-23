package main

import (
	"testing"
)

func TestNewResources_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	_ = NewResources(session)
}
