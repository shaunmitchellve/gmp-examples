# GO GMP Places (New) Text Search Example

## This example utilizes the new golang grpc libraries to perform a text search
https://pkg.go.dev/cloud.google.com/go/maps

### Authentication
- This exmaple use the Google Clouds Application Default Credentials (https://pkg.go.dev/cloud.google.com/go#hdr-Authentication_and_Authorization)

The example code is pretty simple in that it reads in a JSON file and performs a text search against the POI name and Location. It applies a field mask to we only
return a sub-set for fields vs all places fields.