#!/bin/bash

mkdir -p "build"
go get .
go build -o build/app.exe
"%~dp0/build/gaga.exe"


set ACTION=%1