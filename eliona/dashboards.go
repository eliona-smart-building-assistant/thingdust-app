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

package eliona

import (
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

func ThingdustDashboard(projectId string) (api.Dashboard, error) {
	dashboard := api.Dashboard{}
	dashboard.Name = "Thingdust Spaces"
	dashboard.ProjectId = projectId
	

	dashboard.Widgets = []api.Widget{}

	// Process spaces
	assets, _, err := client.NewClient().AssetsApi.
		GetAssets(client.AuthenticationContext()).
		AssetTypeName("thingdust_space").
		ProjectId(projectId).
		Execute()
	if err != nil {
		return api.Dashboard{}, err
	}
	for _, asset := range assets {
		widget := api.Widget{
			WidgetTypeName: "GeneralDisplay",
			AssetId:        asset.Id,
			Details: map[string]interface{}{
				"size":     1,
				"timespan": 7,
			},
			Data: []api.WidgetData{
				{
					ElementSequence: nullableInt32(1),
					AssetId:         asset.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "temperature",
						"description":         "Current Temperature",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(1),
					AssetId:         asset.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "humidity",
						"description":         "Current Humidity",
						"key":                 "",
						"seq":                 1,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(1),
					AssetId:         asset.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "occupancy",
						"description":         "Current Occupancy Status",
						"key":                 "",
						"seq":                 2,
						"subtype":             "input",
					},
				},
			},
		}
		dashboard.Widgets = append(dashboard.Widgets, widget)
	}
	return dashboard, nil
}

func nullableInt32(val int32) api.NullableInt32 {
	return *api.NewNullableInt32(common.Ptr[int32](val))
}
