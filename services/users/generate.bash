#!/bin/sh

protoc --proto_path=proto user.proto --go_out=plugins=grpc:userspb