#!/bin/sh

protoc --proto_path=proto note.proto --go_out=plugins=grpc:notespb