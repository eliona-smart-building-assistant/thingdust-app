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
	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	utilshttp "github.com/eliona-smart-building-assistant/go-utils/http"
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

func initialization() {
	ctx := context.Background()

	// Necessary to close used init resources
	conn := db.NewInitConnectionWithContextAndApplicationName(ctx, app.AppName())
	defer conn.Close(ctx)

	// Init the app before the first run.
	app.Init(db.Pool(), app.AppName(),
		asset.InitAssetTypeFile("eliona/asset-type-thingdust_space.json"),
		app.ExecSqlFile("conf/init.sql"),
		conf.InitConfiguration,
		eliona.InitEliona,
	)
}

func CheckConfigsandSetActiveState() {
	configs, err := conf.GetConfigs(context.Background())
	if err != nil {
		log.Fatal("conf", "Couldn't read configs from DB: %v", err)
		return
	}

	for _, config := range configs {
		// Skip config if disabled and set inactive
		if !conf.IsConfigEnabled(config) {
			if conf.IsConfigActive(config) {
				conf.SetConfigActiveState(config.ConfigId, false)
			}
			continue
		}

		// Signals that this config is active
		if !conf.IsConfigActive(config) {
			conf.SetConfigActiveState(config.ConfigId, true)
			log.Info("conf", "Collecting initialized with Configuration %d:\n"+
				"API Endpoint: %s\n"+
				"API Key: %s\n"+
				"Enable: %t\n"+
				"Refresh Interval: %d\n"+
				"Request Timeout: %d\n"+
				"Active: %t\n"+
				"Project IDs: %v\n",
				config.ConfigId,
				config.ApiEndpoint,
				config.ApiKey,
				*config.Enable,
				config.RefreshInterval,
				config.RequestTimeout,
				*config.Active,
				config.ProjIds)
		}

		// Runs the ReadNode. If the current node is currently running, skip the execution
		// After the execution sleeps the configured timeout. During this timeout no further
		// process for this config is started to read the data.
		common.RunOnceWithParam(func(config apiserver.Configuration) {
			log.Info("main", "Processing Spaces for Configuration with configId %d started", config.ConfigId)

			processSpaces(config)

			log.Info("main", "Processing Spaces for Configuration with configId %d finished", config.ConfigId)

			time.Sleep(time.Second * time.Duration(config.RefreshInterval))
		}, config, config.ConfigId)
	}
}

// For each enabled configuration, processSpaces() performs Continuous Asset Creation
// for each project_id and space pair, and writes corresponding data to each asset.
func processSpaces(config apiserver.Configuration) {
	spaces, err := fetchSpaces(config)
	if err != nil {
		return
	}
	if config.ProjIds != nil {
		for _, projId := range *config.ProjIds {
			log.Debug("projectid", "ProjId: %v", projId)
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

func fetchSpaces(config apiserver.Configuration) (thingdust.Spaces, error) {
	log.Debug("spaces", "Processing space with configID: %v", config.ConfigId)
	request, err := http.NewRequestWithApiKey(config.ApiEndpoint+"/get_space_states", "X-API-KEY", config.ApiKey)
	if err != nil {
		log.Error("spaces", "Error with request: %v", err)
		return nil, err
	}
	spaces, err := http.Read[thingdust.Spaces](request, time.Duration(time.Duration.Seconds(1)), true)
	if err != nil {
		log.Error("spaces", "Error reading spaces: %v", err)
		return nil, err
	}
	return spaces, nil
}

func sendData(confSpace *apiserver.Space, spaces thingdust.Spaces, spaceName string) {
	err := asset.UpsertData(api.Data{
		AssetId: confSpace.AssetId,
		Subtype: api.SUBTYPE_INPUT,
		Data: common.StructToMap(eliona.Data{
			Temperature: spaces[spaceName].Temperature,
			Occupancy:   occupancyToInt(spaces[spaceName].Occupancy),
			Humidity:    spaces[spaceName].Humidity,
		}),
		AssetTypeName: *api.NewNullableString(common.Ptr("thingdust_space")),
	})
	log.Debug("spaces", "Sending data for space %v", spaceName)
	if err != nil {
		log.Error("spaces", "Error sending data %v", err)
	}
}

func getOrCreateMappingIfNecessary(config apiserver.Configuration, projId string, spaceName string) (*apiserver.Space, error) {
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

func createAssetAndMapping(projId string, spaceName string, config apiserver.Configuration) (*apiserver.Space, error) {
	assetId, err := eliona.CreateNewAsset(projId, spaceName)
	if err != nil {
		log.Error("spaces", "Error when creating new asset: %v", err)
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

func occupancyToInt(occupancy string) int64 {
	if occupancy == "occupied" {
		return 1
	} else {
		return 0
	}
}

func listenApiRequests() {
	err := nethttp.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), utilshttp.NewCORSEnabledHandler(
		apiserver.NewRouter(
			apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
			apiserver.NewVersionApiController(apiservices.NewVersionApiService()),
			apiserver.NewCustomizationApiController(apiservices.NewCustomizationApiService()),
		)))
	log.Fatal("main", "Error in API Server: %v", err)
}
