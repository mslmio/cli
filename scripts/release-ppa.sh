#!/bin/bash

# Build and to upload mslm offical PPA.

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1
KEY=$2

if [ -z "$VSN" ]; then
    echo "require version as first parameter" 2>&1
    exit 1
fi

if [ -z "$KEY" ]; then
    echo "require gpg key as second parameter" 2>&1
    exit 1
fi

# building the package
cd $ROOT
debuild -us -uc -S -d

# signing the package
cd $ROOT/..
debsign -k $KEY mslm${VSN}.dsc mslm${VSN}_source.changes

# uploading the package to ppa
dput ppa:mslm/ppa mslm${VSN}_source.changes
