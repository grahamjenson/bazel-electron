#!/bin/bash
set -ex

export TMPDIR=`mktemp -d`
trap 'rm -rf "$TMPDIR"' EXIT

tar -xf {{app}} -C $TMPDIR
open -W $TMPDIR/{{name}}.app