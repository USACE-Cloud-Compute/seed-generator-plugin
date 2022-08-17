#!/bin/bash

TAG=v0.0.5
REPO=williamlehman

docker build . --target prod  -t $REPO/eventgeneratorplugin:$TAG

docker push $REPO/eventgeneratorplugin:$TAG