#!/usr/bin/env bash

set -x
set -e

sh ./precommit.sh
sh ./test.sh

go build
cd ./daemon/
go build
cd ..

DESTINATION="/Users/joe/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/user.workflow.7925D680-5674-4BC2-9CA8-B7019A650147"

./personalmanager -action debug-database
