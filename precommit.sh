#!/usr/bin/env bash

set -x
set -e

echo $PATH
echo $GOPATH

export PROJECT_NAME="github.com/7joe7/personalmanager"

failed=false
PACKAGES=(db operations resources utils)

for i in ${PACKAGES[@]}; do
    i="$PROJECT_NAME/$i"
    go vet $i
    if [[ $(goimports -l $GOPATH/src/$i) ]]; then
        goimports -w $GOPATH/src/$i
        failed=true
    fi
done

if [ $failed == true ]; then echo "Found formatting issues"; exit 0; fi