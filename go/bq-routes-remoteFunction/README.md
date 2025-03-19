<a href="https://idx.google.com/import?url=https://github.com/shaunmitchellve/gmp-examples">
  <img
    height="32"
    alt="Open in IDX"
    src="https://cdn.idx.dev/btn/open_dark_32.svg">
</a>

# BigQuery Remote Function - Routes: Compute Routes Matrix API Example

This example will show you how to set up a BigQuery remote function in order to call a Google Maps Platform API. In this case we are using the 
Routes API: Compute Route Matrix to calculate the distance / duration between an origin and destination. This repo has 2 different examples of how to return the data:

- JSON for multiple values
- metadata for single value

The JSON response is pretty self explanatory, the Cloud Run Function will return the following JSON response: {"Distance": 0, "Duration": 0}.

The metadata response utilizes the userDefinedContext parameter in the POST call to return a single value, either distance or duration. "userDefinedContext": {"mode": "distance"}

# Disclaimer

In accordance with the [Google Maps Platform Terms of Service](https://cloud.google.com/maps-platform/terms/maps-service-terms?hl=en) you can **NOT** save the returned data of this remote function into any
table. You can only view this data.

**Contributors**

- [jbranigan](https://github.com/jbranigan)

**Reference**

- [Big Query Remote Functions Docs](https://cloud.google.com/bigquery/docs/remote-functions)
- [Big Query Remote Function in GO inspiration](https://github.com/salrashid123/bq_cloud_function_golang)

## Setup

### 1. Enable the following APIs

- BigQuery API
- Cloud Run Admin API
- Cloud Build API
- Artifact Registry API
- Cloud Functions API
- Cloud Storage API
- Routes API

### 2. Setup the IAM permissions

Given a default setup, these permissions are applied to the default compute engine service.

- Artifact Registry Writer
- Cloud Run Developer
- Logs Writer
- Service Account User
- Storage Admin

### 3. Deploy the Cloud Run Function

To deploy the metadata (single value return) Cloud Run Function:

- cmd/deploy.sh -e metadata -r $REGION

To deploy the JSON (multiple values) Cloud Run Function:

- cmd/deploy.sh -e json -r $REGION

### 4. Create a Big Query connection to Cloud Run Functions and setup the IAM to execute the Cloud Run Function

[Google Cloud Docs](https://cloud.google.com/bigquery/docs/remote-functions#create_a_connection)

### 5. Create a Big Query User Defined Function

- Single value (metadata) function: (replace $PROJECT, $REGION, $BQ-DATASET and $CONNECTION-NAME)

```
CREATE FUNCTION `$PROJECT.$BQ-DATASET`.routeMatrixDuration(orgLat FLOAT64, orgLng FLOAT64, destLat FLOAT64 , destLng FLOAT64) RETURNS INT64
REMOTE WITH CONNECTION `$PROJECT.us.$CONNECTION-NAME`
OPTIONS (
  endpoint = 'https://$REGION-$PROJECT.cloudfunctions.net/RF_ROUTEMATRIX',
  user_defined_context = [("mode", "duration")],
  max_batching_rows = 3000
);

CREATE FUNCTION `$PROJECT.$BQ-DATASET`.routeMatrixDistance(orgLat FLOAT64, orgLng FLOAT64, destLat FLOAT64 , destLng FLOAT64) RETURNS INT64
REMOTE WITH CONNECTION `$PROJECT.us.$CONNECTION-NAME`
OPTIONS (
  endpoint = 'https://$REGION-$PROJECT.cloudfunctions.net/RF_ROUTEMATRIX',
  user_defined_context = [("mode", "distance")],
  max_batching_rows = 3000
```

- Multiple values (JSON) function: (replace $PROJECT, $REGION, BQ-DATASET and $CONNECTION-NAME)

```
CREATE FUNCTION `$PROJECT.$BQ-DATASET`.routeMatrix(orgLat FLOAT64, orgLng FLOAT64, destLat FLOAT64 , destLng FLOAT64) RETURNS JSON
REMOTE WITH CONNECTION `$PROJECT.us.$CONNECTION-NAME`
OPTIONS (
  endpoint = 'https://$REGION-$PROJECT.cloudfunctions.net/RF_ROUTEMATRIX',
  max_batching_rows = 3000
);
```
**NOTE:** We are trying to protect the Routes API: Compute Route Matrix endpoint quota by setting the limit to 3000 EPM (Entities Per Minute). 3000 EPM is the max quota per minute so if a call maxes out at 3000 rows for the batch then a second call immedatily after will trigger the quota limit (depending on the time it takes to process 3000 calls, this is not fully tested yet.)

### Example queries

- Single value (metadata)

```
 with calculated_distance as (select
 `$PROJECT.$BQ_DATASET.routeMatrixDuration`(ST_Y(origin),ST_X(origin),ST_Y(geom),ST_X(geom)) as duration,
 `$PROJECT.$BQ_DATASET.routeMatrixDistance`(ST_Y(origin),ST_X(origin),ST_Y(geom),ST_X(geom)) as distance,
 from `$BQ_DATASET.$TABLE`)

 select distance, duration
 from calculated_distance;
```

- Multiple value (json)

```
 with calculated_distance as (select
 `$PROJECT.$BQ_DATASET.routeMatrixJSON`(ST_Y(origin),ST_X(origin),ST_Y(geom),ST_X(geom)) as matrix_result
 from `$BQ_DATASET.$TABLE`)

 select
 CAST(JSON_VALUE(matrix_result, '$.distance') as numeric) as distance,
 CAST(JSON_VALUE(matrix_result, '$.duration') as numeric) as duration,
 from calculated_distance;
```