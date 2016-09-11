#!/usr/bin/env bash

set -e
set -x

go test github.com/7joe7/personalmanager/operations
go test github.com/7joe7/personalmanager/db
go test github.com/7joe7/personalmanager/resources
go test github.com/7joe7/personalmanager/alfred
go test github.com/7joe7/personalmanager/checks