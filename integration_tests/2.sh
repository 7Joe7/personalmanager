#!/usr/bin/env bash

set -x
set -e

source ./build.sh

./personalmanager -action "create-task" -name "First task"
./personalmanager -action "print-tasks"
./personalmanager -action "modify-task" -id "0" -name "Modified first task"
./personalmanager -action "print-tasks"
./personalmanager -action "modify-task" -id "0" -active
./personalmanager -action "print-tasks"
./personalmanager -action "delete-task" -id "0"
./personalmanager -action "print-tasks"

source ./teardown.sh