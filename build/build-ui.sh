#!/bin/bash
set -e -x
npm -v
if [ $? -eq 0 ]
then
    cd $1/ui
    npm install
    npm run build
    npm run plugin
else
    docker run --rm -v $1:/app/wecube-plugins-capacity --name capacity-node-build node:12.13.1 /bin/bash /app/wecube-plugins-capacity/build/build-ui-docker.sh
fi