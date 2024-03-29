/*
 * Thingdust app API
 *
 * API to access and configure the Thingdust app
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiservices

import (
	"context"
	//"errors"
	"net/http"
	"thingdust/apiserver"
	"thingdust/conf"
)

// SpacesApiService is a service that implements the logic for the SpacesApiServicer
// This service should implement the business logic for every endpoint for the SpacesApi API.
// Include any external packages or services that will be required by this service.
type SpacesApiService struct {
}

// NewSpacesApiService creates a default api service
func NewSpacesApiService() apiserver.SpacesApiServicer {
	return &SpacesApiService{}
}

// GetSpaces - List all spaces mapped to eliona assets
func (s *SpacesApiService) GetSpaces(ctx context.Context, configId int64) (apiserver.ImplResponse, error) {

	spaces, err := conf.GetSpaces(ctx, configId)
	// If error then return error response
	if err != nil {
		// Code: http.StatusInternalServerError = 500
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	// Return successful ImplResponse with mapped spaces to given configId
	return apiserver.Response(http.StatusOK, spaces), nil
}
