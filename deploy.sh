#!/usr/bin/env bash

set -x
set -e

go build

DESTINATION="/Users/Joe/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/user.workflow.7925D680-5674-4BC2-9CA8-B7019A650147"

cp ./personalmanager "$DESTINATION/personalmanager"
cp -r ./icons "$DESTINATION/icons"