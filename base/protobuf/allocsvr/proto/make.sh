#!/usr/bin/env bash

echo `pwd`
protoc --go_out=plugins=grpc:../ *.proto