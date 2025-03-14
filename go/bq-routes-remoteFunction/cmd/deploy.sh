#!/bin/bash

gcloud functions deploy RF_ROUTEMATRIXJSON \
--gen2 \
--region=us-central1 \
--runtime=go123 \
--source=/home/user/gmp-rf \
--entry-point=computeRouteMatrixJSON \
--trigger-http