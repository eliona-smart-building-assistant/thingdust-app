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

package conf

import (
	"context"
	"thingdust/apiserver"
	dbthingdust "thingdust/db/thingdust"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/boil"
	//"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/null/v8"
)

//
// Todo: Define anything for configuration like structures and methods to read and process configuration
//
func GetSpaces(ctx context.Context, configId int64) ([]apiserver.AssetMapping, error) {
	
	var mods []qm.QueryMod // Declare array of datatype Querymod, which queries can be applied to.

	// configId >0 if not null
	if configId > 0 {
		//append this array which can be queried with the configId
		mods = append(mods, dbthingdust.SpaceWhere.ConfigID.EQ(configId))
		//log.Info("main", "Mods: %v", mods)
	}
	dbAssetMappings, err := dbthingdust.Spaces(mods...).All(ctx, db.Database("thingdust")) //returns space slice and error
	if err != nil {
		return nil, err
	}
	var apiAssetMappings []apiserver.AssetMapping

	// Convert asset mappings into apiAssetMapping type.....why?
	for _, dbAssetMapping := range dbAssetMappings {
		apiAssetMappings = append(apiAssetMappings, *apiAssetMappingFromDbAssetMapping(dbAssetMapping))
	}
	return apiAssetMappings, nil
}



// DeleteConfig reads configured endpoints to Thingdust space
func DeleteConfig(ctx context.Context, configId int64) (int64, error) {
	return dbthingdust.Configs(dbthingdust.ConfigWhere.ConfigID.EQ(configId)).DeleteAll(ctx, db.Database("thingdust"))
}

// GetConfig reads configured endpoints to a Thingdust space
func GetConfig(ctx context.Context, configId int64) (*apiserver.Configuration, error) {
	dbConfigs, err := dbthingdust.Configs(dbthingdust.ConfigWhere.ConfigID.EQ(configId)).All(ctx, db.Database("thingdust"))
	if err != nil {
		return nil, err
	}
	if len(dbConfigs) == 0 {
		return nil, err
	}
	return apiConfigFromDbConfig(dbConfigs[0]), nil
}

// GetConfigs reads all configured endpoints for a Hailo Digital Hub
func GetConfigs(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := dbthingdust.Configs().All(ctx, db.Database("thingdust"))
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		apiConfigs = append(apiConfigs, *apiConfigFromDbConfig(dbConfig))
	}
	return apiConfigs, nil
}



// InsertConfig inserts or updates
func InsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig := dbConfigFromApiConfig(&config)
	err := dbConfig.Insert(ctx, db.Database("thingdust"), boil.Blacklist(dbthingdust.ConfigColumns.ConfigID))
	if err != nil {
		return apiserver.Configuration{}, err
	}
	config.ConfigId = dbConfig.ConfigID
	return config, err
}



// UpsertConfigById inserts or updates
func UpsertConfigById(ctx context.Context, configId int64, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig := dbConfigFromApiConfig(&config)
	dbConfig.ConfigID = configId
	err := dbConfig.Upsert(ctx, db.Database("thingdust"), true,
		[]string{dbthingdust.ConfigColumns.ConfigID},
		boil.Blacklist(dbthingdust.ConfigColumns.ConfigID),
		boil.Infer(),
	)
	config.ConfigId = dbConfig.ConfigID
	return config, err
}



 ///// API to DB Mappings //////


func apiAssetMappingFromDbAssetMapping(dbAssetMapping *dbthingdust.Space) *apiserver.AssetMapping {
	var apiAssetMapping apiserver.AssetMapping
	apiAssetMapping.AssetId = dbAssetMapping.AssetID
	apiAssetMapping.SpaceName = dbAssetMapping.SpaceName
	apiAssetMapping.ConfigId = int32(dbAssetMapping.ConfigID)
	apiAssetMapping.ProjId = dbAssetMapping.ProjectID
	return &apiAssetMapping
}

func apiConfigFromDbConfig(dbConfig *dbthingdust.Config) *apiserver.Configuration {
	var apiConfig apiserver.Configuration
	apiConfig.ConfigId = dbConfig.ConfigID
	apiConfig.ApiEndpoint = dbConfig.APIEndpoint
	apiConfig.ApiKey = dbConfig.APIKey
	apiConfig.Enable = &dbConfig.Enable.Bool
	apiConfig.RefreshInterval = dbConfig.RefreshInterval.Int32
	apiConfig.RequestTimeout = dbConfig.RequestTimeout.Int32
	apiConfig.Active = dbConfig.Active.Bool
	apiConfig.ProjIds = common.Ptr[[]string](dbConfig.ProjectIds)
	return &apiConfig
}

func dbConfigFromApiConfig(apiConfig *apiserver.Configuration) *dbthingdust.Config {
	var dbConfig dbthingdust.Config
	dbConfig.ConfigID = null.Int64FromPtr(&apiConfig.ConfigId).Int64
	dbConfig.APIEndpoint = apiConfig.ApiEndpoint
	dbConfig.APIKey = apiConfig.ApiKey
	dbConfig.Enable = null.BoolFromPtr(apiConfig.Enable)
	dbConfig.RefreshInterval = null.Int32FromPtr(&apiConfig.RefreshInterval)
	dbConfig.RequestTimeout = null.Int32FromPtr(&apiConfig.RequestTimeout)
	dbConfig.Enable = null.BoolFromPtr(&apiConfig.Active)
	if apiConfig.ProjIds != nil {
		dbConfig.ProjectIds = *apiConfig.ProjIds
	}
	return &dbConfig
}