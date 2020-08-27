#!/bin/bash

mkdir -p /app/capacity/public/r_images

sed -i "s~{{CAPACITY_MYSQL_HOST}}~$CAPACITY_MYSQL_HOST~g" /app/capacity/conf/default.json
sed -i "s~{{CAPACITY_MYSQL_PORT}}~$CAPACITY_MYSQL_PORT~g" /app/capacity/conf/default.json
sed -i "s~{{CAPACITY_MYSQL_USER}}~$CAPACITY_MYSQL_USER~g" /app/capacity/conf/default.json
sed -i "s~{{CAPACITY_MYSQL_PWD}}~$CAPACITY_MYSQL_PWD~g" /app/capacity/conf/default.json
sed -i "s~{{CAPACITY_LOG_LEVEL}}~$CAPACITY_LOG_LEVEL~g" /app/capacity/conf/default.json
sed -i "s~{{GATEWAY_URL}}~$GATEWAY_URL~g" /app/capacity/conf/default.json

./server