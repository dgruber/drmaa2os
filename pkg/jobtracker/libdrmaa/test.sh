#!/bin/bash

docker build -t drmaa/drmaajobtrackertest:latest .
docker run --rm -it drmaa/drmaajobtrackertest:latest

