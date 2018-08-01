#!/bin/bash

echo "Building operator code"
go build -o sample-controller .
echo "Copying operator to openshift master"
docker cp sample-controller openshift-master:.


echo "Building operator container"
docker build -t sample-controller .
docker tag sample-controller 192.168.1.73:5000/sample-controller
docker push 192.168.1.73:5000/sample-controller

docker exec openshift-node-1 docker pull 192.168.1.73:5000/sample-controller
docker exec openshift-node-2 docker pull 192.168.1.73:5000/sample-controller
