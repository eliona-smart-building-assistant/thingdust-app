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
	"thingdust/thingdust"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"thingdust/eliona"
	"time"
	//api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
)


// doAnything is the main app function which is called periodically
func doAnything() {
	// Fetch Spacs using thingdust api
	request, err := http.NewRequestWithApiKey("https://demo.cust.prod.thingdust.io/api/v2/get_space_states", "X-API-KEY", "UEKNEYKACORWF9JMYBGLPOCPIBHNJUHYIAADBRQCEHQM2V7YJUSCVBFUNOWW" )
	if err != nil {
		log.Error("spaces", "Error with request: %v", err)
		return
	}
	spaces, err1 := http.Read[thingdust.Spaces](request, time.Duration(time.Duration.Seconds(1)), true)
	if err1 != nil {
		log.Error("spaces", "Error reading spaces: %v",err1)
		return
	}
	for spaceName:= range spaces {
		// Mapping exists?
		confSpace, err2 :=conf.GetSpace(context.Background(), 1, "empty", spaceName)
		if err2 != nil {
			log.Error("spaces", "Error when reading spaces from configurations")
			return
		}
		if confSpace == nil {
			// Create Asset Mapping between newly generated asset and space with key spaceName
			assetId, err := eliona.CreateNewAsset("empty", spaceName)
			if err != nil {
				log.Error("spaces", "Error when creating new asset")
				return
			}
			if err == nil {
				log.Debug("spaces", "AssetId %v assigned to space %v", assetId, spaceName)
			}
			err = conf.InsertSpace(context.Background(), 1, "empty", spaceName, assetId)
			if err != nil {
				log.Error("spaces","Error when inserting space into database")
				return
			}
		} else {
			//Asset exists
			exists, err := asset.ExistAsset(confSpace.AssetId)
			if err != nil {
				log.Error("spaces","Error when checking if asset already exists")
				return
			}
			if exists {
				log.Debug("spaces", "Asset already exists for space %v with AssetId %v", spaceName, confSpace.AssetId)
			}else {
				continue
			}
			

		}
		//  asset.UpsertData(api.Data{
		// 	AssetId: confSpace.AssetId,
		// 	Subtype: "input",
		// 	Data: spaces[spaceName].(thingdust.Space),
		// 	AssetTypeName: *api.NewNullableString(common.Ptr("thingdust_space")),
		//  })

		// // })
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
