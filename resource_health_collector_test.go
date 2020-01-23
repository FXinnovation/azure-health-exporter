package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/mock"
)

type MockedResourceHealth struct {
	mock.Mock
}

type MockedResources struct {
	mock.Mock
}

func (mock *MockedResourceHealth) GetAvailabilityStatus(resourceURI string) (*resourcehealth.AvailabilityStatus, error) {
	args := mock.Called(resourceURI)
	return args.Get(0).(*resourcehealth.AvailabilityStatus), args.Error(1)
}

func (mock *MockedResourceHealth) GetAllAvailabilityStatuses() (*[]resourcehealth.AvailabilityStatus, error) {
	args := mock.Called()
	return args.Get(0).(*[]resourcehealth.AvailabilityStatus), args.Error(1)
}

func (mock *MockedResourceHealth) GetSubscriptionID() string {
	args := mock.Called()
	return args.Get(0).(string)
}

func (mock *MockedResources) GetResources(resourceType string, resourceTags map[string]string) (*[]resources.GenericResource, error) {
	args := mock.Called(resourceType, resourceTags)
	return args.Get(0).(*[]resources.GenericResource), args.Error(1)
}

func CallExporter(collector ResourceHealthCollector) *httptest.ResponseRecorder {
	loadConfig("config/config_example.yml")
	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()
	registry := prometheus.NewRegistry()
	registry.MustRegister(&collector)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(rr, req)
	return rr
}

func TestNewResourceHealthCollector_OK(t *testing.T) {
	session, err := NewAzureSession("subscriptionID")
	if err != nil {
		t.Errorf("Error occured %s", err)
	}
	_ = NewResourceHealthCollector(session)
}

func TestCollect_GetResources_Error(t *testing.T) {
	r := MockedResources{}
	rh := MockedResourceHealth{}
	collector := ResourceHealthCollector{
		resourceHealth: &rh,
		resources:      &r,
	}

	var resList []resources.GenericResource
	r.On("GetResources", mock.Anything, mock.Anything).Return(&resList, errors.New("Unit test Error"))

	rr := CallExporter(collector)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestCollect_GetAvailabilityStatus_Error(t *testing.T) {
	r := MockedResources{}
	rh := MockedResourceHealth{}
	collector := ResourceHealthCollector{
		resourceHealth: &rh,
		resources:      &r,
	}

	var resList []resources.GenericResource
	resourceID := "id"
	resourceType := "type"
	resList = append(resList, resources.GenericResource{
		ID:   &resourceID,
		Type: &resourceType,
	})
	r.On("GetResources", mock.Anything, mock.Anything).Return(&resList, nil)

	var as resourcehealth.AvailabilityStatus
	rh.On("GetAvailabilityStatus", mock.Anything).Return(&as, errors.New("Unit test Error"))
	rh.On("GetSubscriptionID").Return("my_subscription")

	rr := CallExporter(collector)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Wrong status code: got %v, want %v", status, http.StatusInternalServerError)
	}
}
