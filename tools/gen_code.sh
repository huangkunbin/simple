#!/usr/bin/env bash -v

cd "../../$(dirname "$0")"

workdir=$(pwd)

cd $workdir/simple

find . -name "*.simpleapi.go" -exec rm {} \;

go run ./cmd/game_server/ -gencode=true -genpath=$workdir