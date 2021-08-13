#!/bin/bash
source ./config.sh

cmds=(
    game
)

cd $(dirname $0)/..

cmd=$1
while true; do
    case "$cmd" in
    game)
        go run ./cmd/game_server
        exit 0
        ;;
    *)
        read -p "输入要执行的程序:" cmd
        continue
        ;;
    esac
done
