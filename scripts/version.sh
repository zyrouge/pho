#!/bin/bash

set -eu

sed -nr 's/const AppVersion = "([[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+)"/\1/p' ./core/meta.go
