#!/usr/bin/env bash

coverdir=$(mktemp -d /tmp/coverage.XXXXXXXXXX)
profile="${coverdir}/cover.out"

go test -coverprofile=${profile} ./...
go tool cover -func ${profile}
