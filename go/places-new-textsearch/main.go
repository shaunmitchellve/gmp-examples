package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	places "cloud.google.com/go/maps/places/apiv1"
	placespb "cloud.google.com/go/maps/places/apiv1/placespb"
	"google.golang.org/grpc/metadata"
)

type LLMResponse []struct {
	POI			string	`json:"POI"`
	Location	string	`json:"Location"`
}

func main() {
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx,
	metadata.Pairs("X-Goog-FieldMask", "places.displayName,places.id,places.formattedAddress,places.rating,places.types,places.businessStatus"))
	placesClient, err := places.NewClient(ctx)

	if err != nil {
		fmt.Printf("Places Client Error: %v", err)
		os.Exit(1)
	}

	defer placesClient.Close()

	jsonInput, err := os.Open("POI_Search.json")

	if err != nil {
		fmt.Printf("Unable to open Input JSON file: %v", err)
		os.Exit(1)
	}

	jsonOutput, err := os.Create("Places_response.json")

	if err != nil {
		fmt.Printf("Unable to create output json file for writing: %v", err)
		os.Exit(1)
	}

	jsonOutput.WriteString("[\n")

	defer jsonInput.Close()
	defer jsonOutput.Close()

	jsonBytes, _ := io.ReadAll(jsonInput)

	var llmResponse LLMResponse

	json.Unmarshal(jsonBytes, &llmResponse)

	for i := 0; i < len(llmResponse); i++ {
		outPutLine := ""

		if i < len(llmResponse) && i != 0 {
			outPutLine += ",\n"
		}

		req := &placespb.SearchTextRequest{
			TextQuery: llmResponse[i].POI + ", " +llmResponse[i].Location,
		}

		resp, err := placesClient.SearchText(ctx, req)

		if err != nil {
			fmt.Printf("Places Text Search Error %v", err)
			os.Exit(1)
		}

		places := resp.GetPlaces()

		for a := 0; a < len(places); a++ {
			placesJson, err := json.MarshalIndent(places[a],"","	")

			if err != nil {
				fmt.Printf("Unable to marshal Places return: %v", err)
			} else {
				outPutLine += (string(placesJson))
				if (a+1) < len(places) {
					outPutLine += ",\n"
				}
			}
		}

		jsonOutput.WriteString(outPutLine)
	}

	jsonOutput.WriteString("\n]")
}