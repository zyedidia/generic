#!/bin/bash

SCRIPT_DIR="$(realpath $(dirname $0))"
# useful for early termination with Ctrl+C
trap 'trap - SIGINT; kill -SIGINT $$' SIGINT;

cd $SCRIPT_DIR
# pass environment variables to support benchmark configuration
RANGES="$RANGES" MAPS="$MAPS" go test -bench=.  -benchtime=3x
