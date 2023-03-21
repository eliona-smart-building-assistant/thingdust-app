//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	nethttp "net/http"
	"thingdust/apiserver"
	"thingdust/apiservices"
	"thingdust/conf"
	"thingdust/eliona"
	"thingdust/thingdust"
	"time"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// For each enabled configuration, processSpaces() performs Continuous Asset Creation
// for each project_id and space pair, and writes corresponding data to each asset.
func processSpaces(configId int64) {
	config, spaces, err := fetchSpacesAndSetActiveState(configId)
	if err != nil {
		return
	}
	if config.ProjIds != nil {
		for _, projId := range *config.ProjIds {
			for spaceName := range spaces {
				confSpace, err := getOrCreateMappingIfNecessary(config, projId, spaceName)
				if err != nil {
					return
				}
				if confSpace != nil {
					sendData(confSpace, spaces, spaceName)
				}
			}
		}
	}
}

func fetchSpacesAndSetActiveState(configId int64) (*apiserver.Configuration, thingdust.Spaces, error) {
	config, err := conf.GetConfig(context.Background(), configId)
	if err != nil {
		log.Error("spaces", "Error reading configuration")
		return nil, nil, err
	}
	if config.Enable == nil || !*config.Enable {
		conf.SetConfigActiveState(config.ConfigId, false)
		return nil, nil, err
	}
	conf.SetConfigActiveState(config.ConfigId, true)
	log.Debug("spaces", "Processing space with configID: %v", config.ConfigId)
	request, err := http.NewRequestWithApiKey(config.ApiEndpoint + "/get_space_states", "X-API-KEY", config.ApiKey)
	if err != nil {
		log.Error("spaces", "Error with request: %v", err)
		return nil, nil, err
	}
	spaces, err := http.Read[thingdust.Spaces](request, time.Duration(time.Duration.Seconds(1)), true)
	if err != nil {
		log.Error("spaces", "Error reading spaces: %v", err)
		return nil, nil, err
	}
	return config, spaces, nil
}

func sendData(confSpace *apiserver.Space, spaces thingdust.Spaces, spaceName string) {
	err := asset.UpsertData(api.Data{
		AssetId: confSpace.AssetId,
		Subtype: "input",
		Data: common.StructToMap(eliona.Data{
			Temperature: spaces[spaceName].Temperature,
			Occupancy:   occupancyToBool(spaces[spaceName].Occupancy),
			Humidity:    spaces[spaceName].Humidity,
		}),
		AssetTypeName: *api.NewNullableString(common.Ptr("thingdust_space")),
	})
	if err != nil {
		log.Error("spaces", "Error sending data %v", err)
	}
}

func getOrCreateMappingIfNecessary(config *apiserver.Configuration, projId string, spaceName string) (*apiserver.Space, error) {
	var confSpace *apiserver.Space
	confSpace, err := conf.GetSpace(context.Background(), config.ConfigId, projId, spaceName)
	if err != nil {
		log.Error("spaces", "Error when reading spaces from configurations")
		return nil, err
	}
	if confSpace == nil {
		confSpace, err = createAssetAndMapping(projId, spaceName, config)
		if err != nil {
			return nil, err
		}
	} else {
		exists, err := asset.ExistAsset(confSpace.AssetId)
		if err != nil {
			log.Error("spaces", "Error when checking if asset already exists")
			return nil, err
		}
		if exists {
			log.Debug("spaces", "Asset already exists for space %v with AssetId %v", spaceName, confSpace.AssetId)
		} else {
			log.Debug("spaces", "Asset with AssetId %v does no longer exist in eliona", confSpace.AssetId)
			return nil, nil
		}
	}
	return confSpace, nil
}

func createAssetAndMapping(projId string, spaceName string, config *apiserver.Configuration) (*apiserver.Space, error) {
	assetId, err := eliona.CreateNewAsset(projId, spaceName)
	if err != nil {
		log.Error("spaces", "Error when creating new asset")
		return nil, err
	}
	log.Debug("spaces", "AssetId %v assigned to space %v", assetId, spaceName)
	err = conf.InsertSpace(context.Background(), config.ConfigId, projId, spaceName, assetId)
	if err != nil {
		log.Error("spaces", "Error when inserting space into database")
		return nil, err
	}
	log.Debug("spaces", "Asset with AssetId %v corresponding to space %v inserted into eliona database", assetId, spaceName)
	confSpace, err := conf.GetSpace(context.Background(), config.ConfigId, projId, spaceName)
	if err != nil {
		log.Error("spaces", "Error when reading spaces from configurations")
		return nil, err
	}
	return confSpace, nil
}

func occupancyToBool(occupancy string) bool {
	if occupancy == "occupied" {
		return true
	} else {
		return false
	}
}

// listenApi starts the API server and listen for requests
func listenApi() {
	http.ListenApiWithOs(&nethttp.Server{Addr: ":" + common.Getenv("API_SERVER_PORT", "3000"), Handler: apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
		apiserver.NewVersionApiController(apiservices.NewVersionApiService()),
		apiserver.NewCustomizationApiController(apiservices.NewCustomizationApiService()),
		apiserver.NewSpacesApiController(apiservices.NewSpacesApiService()),
	)})
}
