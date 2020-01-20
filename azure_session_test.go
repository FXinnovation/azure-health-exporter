package main

import (
	"testing"
)

func TestNewAzureSession_OK(t *testing.T) {
	wantSubscriptionID := "subscriptionID"
	session, err := NewAzureSession(wantSubscriptionID)
	if err != nil {
		t.Errorf("Error occured %s", err)
	}

	if session.SubscriptionID != wantSubscriptionID {
		t.Errorf("Unexpected SubscriptionID; got: %v, want: %v", session.SubscriptionID, wantSubscriptionID)
	}
}

func TestNewAzureSession_MissingSubscriptionID(t *testing.T) {
	_, err := NewAzureSession("")

	if err == nil {
		t.Errorf("Want an error, got none")
	}
}
