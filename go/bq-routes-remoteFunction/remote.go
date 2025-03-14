// Remote function logic copied from: https://github.com/salrashid123/bq_cloud_function_golang/tree/main

package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	routes "cloud.google.com/go/maps/routing/apiv2"
	routespb "cloud.google.com/go/maps/routing/apiv2/routingpb"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	latlng "google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc/metadata"
)

type bqRequest struct {
	RequestId          string            `json:"requestId"`
	Caller             string            `json:"caller"`
	SessionUser        string            `json:"sessionUser"`
	UserDefinedContext map[string]string `json:"userDefinedContext"`
	Calls              [][]interface{}   `json:"calls"`
}

type bqResponse struct {
	Replies      []int64`json:"replies,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

func init() {
	functions.HTTP("computeRouteMatrix", routeMatrix)
}

func routeMatrix(w http.ResponseWriter, r *http.Request) {
	wait := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bqReq := &bqRequest{}
	bqResp := &bqResponse{}

	if err := json.NewDecoder(r.Body).Decode(&bqReq); err != nil {
		bqResp.ErrorMessage = fmt.Sprintf("External Function error: can't read POST body %v", err)
	} else {

		fmt.Printf("caller %s\n", bqReq.Caller)
		fmt.Printf("sessionUser %s\n", bqReq.SessionUser)
		fmt.Printf("userDefinedContext %v\n", bqReq.UserDefinedContext)
		fmt.Printf("calls %s\n", bqReq.Calls)

		objs := make([]int64, len(bqReq.Calls))
		// Calls should be 1, not being batched
		// Each call should have 4 numbers: [origin_LAT, origin_LNG, DEST_LAT, DEST_LNG]
		for i, r := range bqReq.Calls {
			if len(r) != 4 {
				bqResp.ErrorMessage = fmt.Sprintf("Invalid number of input fields provided.  expected 4, got  %d", len(r))
			}

			orgLat, ok := r[0].(float64)
			if !ok {
				bqResp.ErrorMessage = fmt.Sprintf("Origin Lat is not a number: %d", r[0])
			}
			orgLng, ok := r[1].(float64)
			if !ok {
				bqResp.ErrorMessage = fmt.Sprintf("Origin Lng is not a number: %d", r[0])
			}
			destLat, ok := r[2].(float64)
			if !ok {
				bqResp.ErrorMessage = fmt.Sprintf("Destination Lat is not a number: %d", r[0])
			}
			destLng, ok := r[3].(float64)
			if !ok {
				bqResp.ErrorMessage = fmt.Sprintf("Destination Lng is not a number: %d", r[0])
			}

			if bqResp.ErrorMessage != "" {
				bqResp.Replies = []int64{0}
				break
			}

			wait.Add(1)
			go func(j int) {
				defer wait.Done()

				for {
					select {
					case <-ctx.Done():
						return
					default:
						ctx = metadata.NewOutgoingContext(ctx,
							metadata.Pairs("X-Goog-FieldMask", "*"))

						if routesClient, err := routes.NewRoutesClient(ctx); err != nil {
							bqResp.ErrorMessage = fmt.Sprintf("Routes Client Error: %v", err)
						} else {
							defer routesClient.Close()

							var origins []*routespb.RouteMatrixOrigin
							var destinations []*routespb.RouteMatrixDestination

							origins = append(origins, &routespb.RouteMatrixOrigin{
								Waypoint: &routespb.Waypoint{
									LocationType: &routespb.Waypoint_Location{
										Location: &routespb.Location{
											LatLng: &latlng.LatLng{
												Latitude:  orgLat,
												Longitude: orgLng,
											},
										},
									},
								},
							})

							destinations = append(destinations, &routespb.RouteMatrixDestination{
								Waypoint: &routespb.Waypoint{
									LocationType: &routespb.Waypoint_Location{
										Location: &routespb.Location{
											LatLng: &latlng.LatLng{
												Latitude:  destLat,
												Longitude: destLng,
											},
										},
									},
								},
							})

							req := routespb.ComputeRouteMatrixRequest{
								Origins:      origins,
								Destinations: destinations,
							}

							if res, err := routesClient.ComputeRouteMatrix(ctx, &req); err != nil {
								bqResp.ErrorMessage = fmt.Sprintf("Compute Route Matrix Error: %v", err)
							} else {
								if rtMatrix, err := res.Recv(); err != nil {
									bqResp.ErrorMessage = fmt.Sprintf("Compute Route Matrix Read Error: %v", err)
								} else if rtMatrix.Condition != *routespb.RouteMatrixElementCondition_ROUTE_EXISTS.Enum() {
									bqResp.ErrorMessage = "Route does not exist"
								} else {
									if bqReq.UserDefinedContext["mode"] == "duration"{
										objs[j] = rtMatrix.Duration.Seconds
									} else if bqReq.UserDefinedContext["mode"] == "distance" {
										objs[j] = int64(rtMatrix.DistanceMeters)
									}

									return
								}
							}
						}
					}
				}
			}(i)
		}

		wait.Wait()
		if bqResp.ErrorMessage != "" {
			bqResp.Replies = []int64{0}
		} else {
			bqResp.Replies = objs
		}

		b, err := json.Marshal(bqResp)
		if err != nil {
			http.Error(w, fmt.Sprintf("can't convert response to JSON %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
