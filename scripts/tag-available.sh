#!/bin/bash

tag=$1

if [ "$tag" == "" ]; then
    echo "No tag name specified"
    exit 1
fi

output=$(git ls-remote --exit-code --tags origin "$tag")

if [ "$output" == "" ]; then
    echo "Tag $tag is available"
    exit 0
fi

echo "Tag $tag is not available"
exit 1
