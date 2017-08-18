#!/usr/bin/env bash
rm githorse
#make bindata
go build -gcflags='-N -l' github.com/gityflow/githorse &&  dlv --listen=:8000 --headless=true --api-version=2 exec ./githorse web

