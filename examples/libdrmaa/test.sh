#!/bin/bash

docker build -t drmaa/drmaa2oslibdrmaaexample:latest .
docker run --rm -it drmaa/drmaa2oslibdrmaaexample:latest /bin/bash

