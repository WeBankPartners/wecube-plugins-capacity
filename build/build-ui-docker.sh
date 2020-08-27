#!/bin/bash
set -e -x
cd $(dirname $0)/../ui
npm install
npm run build
npm run plugin