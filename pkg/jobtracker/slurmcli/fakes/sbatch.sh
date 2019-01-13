#!/bin/bash
if [  "${!#}X" = "failX" ]; then
	exit 1
fi
echo "77;clustername"
echo ""