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

cp ./daemon/daemon "$DESTINATION/daemon"
cp ./daemon/org.erneker.personalmanager.plist "/Users/joe/Library/LaunchAgents/org.erneker.personalmanager.plist"
cp ./personalmanager "$DESTINATION/personalmanager"
cp -r ./icons "$DESTINATION"

launchctl unload /Users/joe/Library/LaunchAgents/org.erneker.personalmanager.plist
launchctl load /Users/joe/Library/LaunchAgents/org.erneker.personalmanager.plist