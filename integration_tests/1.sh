#!/usr/bin/env bash

set -x
set -e

source ./build.sh

./personalmanager -action "create-habit" -name "First habit"
./personalmanager -action "print-habits"
./personalmanager -action "delete-habit" -id "0"
./personalmanager -action "create-habit" -name "Second habit"
./personalmanager -action "print-habits"
./personalmanager -action "modify-habit" -id "1" -active -repetition "Daily"
./personalmanager -action "print-habits"
./personalmanager -action "modify-habit" -id "1" -name "Modified second habit"
./personalmanager -action "print-habits"
./personalmanager -action "modify-habit" -id "1" -deadline "2.1.2016"
./personalmanager -action "print-habits"
./personalmanager -action "modify-habit" -id "1" -basePoints 14
./personalmanager -action "print-habits"
./personalmanager -action "modify-habit" -id "1" -done
./personalmanager -action "print-habits"
./personalmanager -action "delete-habit" -id "1"
./personalmanager -action "print-habits"

source ./teardown.sh