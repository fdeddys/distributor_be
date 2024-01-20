#!/usr/bin/env bash

env GOOS=linux GOARCH=amd64 go build -o apotik_be
# env GOOS=windows GOARCH=amd64 go build -o apotik_be.exe
# pg_restore -U username -d dbname -1 filename.dump
