#!/bin/bash

CGO_ENABLED=0 go build . && \
    sudo docker build -t b.gcr.io/txregistry/random . && \
    sudo gcloud docker push b.gcr.io/txregistry/random

