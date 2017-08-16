#!/usr/bin/env bash
docker rm -f postgres-db
docker-compose -f compose-postgres.yml up -d
