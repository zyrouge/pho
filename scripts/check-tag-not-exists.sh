#!/bin/bash

tag=$1
git ls-remote --exit-code --tags origin $tag
status=$?
if [ $status == 0 ]; then
    echo "Tag $tag is not available"
    exit 1
fi

echo "Tag $tag is available"
exit 0
