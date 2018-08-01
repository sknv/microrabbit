#!/usr/bin/env bash

# start the math service development server
# usage example: bin/devserver-math.sh

# load environment variables from .env file
set -a; [ -f .env ] && . .env; set +a

# start the dev server
go run app/services/math/main.go
