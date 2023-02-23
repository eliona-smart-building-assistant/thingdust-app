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
	//"fmt"
	"context"
	nethttp "net/http"

	"thingdust/apiserver"
	"thingdust/apiservices"
	"thingdust/conf"

	//"os"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/http"

	//"encoding/json"

	"github.com/eliona-smart-building-assistant/go-utils/log"

	//"io/ioutil"
	"time"
)

type Spaces map[string]Space


type Space struct {
	Humidity float64 `json:"humidity"`
	Occupancy string `json:"occupancy"`
	Temperature float64 `json:"temperature"`
}

// doAnything is the main app function which is called periodically
func doAnything() {


	request, err := http.NewRequestWithApiKey("https://demo.cust.prod.thingdust.io/api/v2/get_space_states", "X-API-KEY", "UEKNEYKACORWF9JMYBGLPOCPIBHNJUHYIAADBRQCEHQM2V7YJUSCVBFUNOWW" )
	
	if err != nil {
		log.Error("spaces", "Error with request: %v", err)
		return
	}

	spaces, err1 := http.Read[Spaces](request, time.Duration(time.Duration.Seconds(1)), true)

	if err1 != nil {
		log.Error("spaces", "Error reading spaces: %v",err1)
		return
	}
	log.Debug("spaces", "Read %d spaces", len(spaces))

	for spaceName:= range spaces {
		// Mapping exists?
		confSpace, err2 :=conf.GetSpace(context.Background(), 1, "", spaceName)
		if err2 != nil {
			log.Error("spaces", "Error when reading spaces from configurations")
			return
		}

		if confSpace == nil {
			// Create Asset Mapping
		}
		else{
			//Asset exists
		}
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
