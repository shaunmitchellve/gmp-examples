#!/bin/bash

while getopts e:r: flag
do
    case "${flag}" in
        e) entry=${OPTARG};;
        r) region=${OPTARG};;
    esac
done

if [ ${#entry} -eq 0 ] || [ ${#region} -eq 0 ]; then
    echo "missing required field(s)\n"
    echo "Usage: ./deploy.sh -e json | metadata -r <region>"
    exit 1
fi

if [ "$entry" == "json" ]; then
    gcloud functions deploy RF_ROUTEMATRIX \
    --gen2 \
    --region=$region \
    --runtime=go123 \
    --source=/home/user/gmp-examples/go/bq-routes-remoteFunction \
    --entry-point=computeRouteMatrixJSON \
    --trigger-http
elif [ "$entry" == "metadata" ]; then
    gcloud functions deploy RF_ROUTEMATRIX \
    --gen2 \
    --region=$region \
    --runtime=go123 \
    --source=/home/user/gmp-examples/go/bq-routes-remoteFunction \
    --entry-point=computeRouteMatrix \
    --trigger-http
fi