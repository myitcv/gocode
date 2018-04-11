#!/usr/bin/env bash

set -eu

if [ "${CI:-}" == "true" ]
then
	go get -u golang.org/x/vgo
	pushd $(go list -f "{{.Dir}}" golang.org/x/vgo) > /dev/null

	git fetch -q https://go.googlesource.com/vgo refs/changes/55/105855/3 && git checkout -qf FETCH_HEAD
	go install

	popd > /dev/null
fi

# for the tests in ./vgo
go get golang.org/x/net/html
