#!/bin/bash

go build -o opencall cmd/web/*.go
./opencall -dbname=opencall -dbuser=postgres -dbpass=root -cache=false -production=false