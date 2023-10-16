#!/bin/bash

# Build and upload (to GitHub) for all platforms for cli $1 & version $2.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI="mslm"
VSN=$1

if [ -z "$VSN" ]; then
    echo "require version as first parameter" 2>&1
    exit 1
fi

# build
$ROOT/scripts/build-archive-all.sh "$CLI" "$VSN"

# release
gh release create ${CLI}-${VSN}                                               \
    -R ipinfo/mslm                                                             \
    -t "${CLI}-${VSN}"                                                        \
    $ROOT/build/${CLI}_${VSN}*.tar.gz                                         \
    $ROOT/build/${CLI}_${VSN}*.zip                                            \
    $ROOT/build/${CLI}_${VSN}*.deb                                            \
    $ROOT/${CLI}/macos.sh                                                     \
    $ROOT/${CLI}/windows.ps1                                                  \
    $ROOT/${CLI}/deb.sh
