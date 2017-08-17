#!/usr/bin/env bash
rm githorse
#make bindata
go build
./githorse web
