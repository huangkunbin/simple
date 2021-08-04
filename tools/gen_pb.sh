#!/usr/bin/env bash -v

protoc -I=../proto/ --go_out=../api ../proto/*