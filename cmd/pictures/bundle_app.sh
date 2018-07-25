#!/bin/sh

wd=$PWD
cd ../../app
polymer build --add-service-worker --bundle
mv build $wd
cd $wd
go-bindata -prefix "build/default" "build/default/..."
rm -rf build
