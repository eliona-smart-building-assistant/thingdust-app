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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/null/v8"
)

func GetSpaces(ctx context.Context, configId int64) ([]apiserver.Space, error) {
	var mods []qm.QueryMod // Declare array of datatype Querymod, which queries can be applied to.
	// configId >0 if not null
	if configId > 0 {
		mods = append(mods, dbthingdust.SpaceWhere.ConfigID.EQ(configId))
	}
	dbSpaces, err := dbthingdust.Spaces(mods...).All(ctx, db.Database("thingdust")) //returns space slice and error
	if err != nil {
		return nil, err
	}
	var apiSpaces []apiserver.Space
	for _, dbSpaces := range dbSpaces {
		apiSpaces = append(apiSpaces, *apiSpacesFromDbSpaces(dbSpaces))
	}
	return apiSpaces, nil
}

// Return space if space already exists in app database (with assetid etc.) otherwise return nil
func GetSpace(ctx context.Context, configId int64, projectId string, spaceName string) (*apiserver.Space, error) {
	var mods []qm.QueryMod // Declare array of datatype Querymod, which queries can be applied to.
	if configId > 0 {
		mods = append(mods, dbthingdust.SpaceWhere.ConfigID.EQ(configId))
		mods = append(mods, dbthingdust.SpaceWhere.ProjectID.EQ(projectId))
		mods = append(mods, dbthingdust.SpaceWhere.SpaceName.EQ(spaceName))
	}
	//dbSpaces is an array of length 1, as configId, projectId and space_name form a primary key
	dbSpaces, err := dbthingdust.Spaces(mods...).All(ctx, db.Database("thingdust")) //returns space slice and error
	if err != nil {
		return nil, err
	}
	if len(dbSpaces)!= 1 {
		return nil, nil
	}
	return apiSpacesFromDbSpaces(dbSpaces[0]), nil
}

func InsertSpace(ctx context.Context, configId int64, projectId string, spaceName string, assetId int32) error {
	var dbSpace dbthingdust.Space
	dbSpace.ConfigID = configId
	dbSpace.ProjectID = projectId
	dbSpace.AssetID = assetId
	dbSpace.SpaceName = spaceName
	return dbSpace.Insert(ctx, db.Database("thingdust"), boil.Infer())
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

func apiSpacesFromDbSpaces(dbSpaces *dbthingdust.Space) *apiserver.Space {
	var apiSpaces apiserver.Space
	apiSpaces.AssetId = dbSpaces.AssetID
	apiSpaces.SpaceName = dbSpaces.SpaceName
	apiSpaces.ConfigId = int32(dbSpaces.ConfigID)
	apiSpaces.ProjId = dbSpaces.ProjectID
	return &apiSpaces
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

func SetConfigActiveState(ctx context.Context, config apiserver.Configuration, state bool) (int64, error) {
	return dbthingdust.Configs(
		dbthingdust.ConfigWhere.ConfigID.EQ(null.Int64FromPtr(&config.ConfigId).Int64),
	).UpdateAll(ctx, db.Database("thingdust"), dbthingdust.M{
		dbthingdust.ConfigColumns.Active: state,
	})
}

