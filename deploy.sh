#!/bin/sh

docker-compose -f docker-compose-build.yaml   up --build -d

#docker cp buildtmp:/go/bin/master.exe ./ && docker cp buildtmp:/go/bin/worker.exe ./
#nohup ./master.exe &
#nohup ./worker.exe &

docker cp buildtmp:/go/bin/master ./ && docker cp buildtmp:/go/bin/worker ./
nohup ./master &
nohup ./worker &