#!/bin/bash

oapi-codegen -generate types,chi-server,spec -package genserver ./jobtracker_1_0_0_openapi_v3.yaml > server/generated/jobtracker_generated_1_0_0.go
oapi-codegen -generate types,client,spec -package genclient ./jobtracker_1_0_0_openapi_v3.yaml > client/generated/jobtracker_generated_1_0_0.go
