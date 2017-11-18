#!/bin/sh
fly -t ci set-pipeline --config pipeline.yml --pipeline drmaa2os --load-vars-from params.yml 
