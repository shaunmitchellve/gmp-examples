#!/bin/bash

# origin: -33.47683, -70.61738 destination: -33.527268, -70.613372
curl -d '{"requestId": "124ab1c","caller": "//bigquery.googleapis.com/projects/myproject/jobs/myproject:US.bquxjob_5b4c112c_17961fafeaf","sessionUser": "test-user@test-company.com","userDefinedContext": {"mode": "distance"},"calls": [[-33.47683, -70.61738,-33.527268,-70.613372],[-33.47683, -70.61738,-33.527268,-70.613372]]}' \
-H "Content-Type: application/json" \
-X POST \
localhost:8080