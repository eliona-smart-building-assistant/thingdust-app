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
	"errors"
	"net/http"
	"thingdust/apiserver"
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
	// TODO - update GetSpaces with the required logic for this service method.
	// Add api_spaces_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []AssetMapping{}) or use other options such as http.Ok ...
	//return Response(200, []AssetMapping{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetSpaces method not implemented")
}