#!/bin/bash

docker build -t drmaa/drmaajobtrackertest:latest -f ./Dockerfiles/libdrmaa/Dockerfile .
docker run --rm -it drmaa/drmaajobtrackertest:latest /bin/bash
#docker run --rm -it drmaa/drmaajobtrackertest:latest
