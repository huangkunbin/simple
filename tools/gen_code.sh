#!/usr/bin/env bash -v

cd "../../$(dirname "$0")"

workdir=$(pwd)

cd $workdir/simple

go run ./cmd/game_server/ -gencode=true -genpath=$workdir