#!/bin/bash

echo ""
echo "For slurm drmaa2os -> drmaa:"
echo "----------------------------"
echo "cd /root/go/src/github.com/dgruber/drmaa2os/examples/libdrmaa"
echo "export CGO_LDFLAGS="-L/usr/local/lib"
echo "export CGO_CFLAGS="-DSLURM -I/usr/local/include""
echo "export LD_LIBRARY_PATH=/usr/local/lib"
echo "./libdrmaa"
echo ""
echo "Enable accounting for slurm (sacct):"
echo "------------------------------------"
echo "sacctmgr --immediate add cluster name=linux"
echo "supervisorctl restart slurmdbd"
echo "supervisorctl restart slurmctld"
echo "yes | sacctmgr add account none,test Cluster=linux Description=\"none\" Organization=\"none\""

echo ""
echo "Starting container"

docker run --rm -it -h ernie slurm-drmaa2-dev
