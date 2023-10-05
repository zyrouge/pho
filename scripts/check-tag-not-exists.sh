#!/bin/bash

git ls-remote --exit-code --tags origin $1
status=$?
if [[ status -eq 0 ]] then
    exit 1
fi

exit 0
