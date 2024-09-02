#!/bin/bash

EOE_TESTER_WORKING_DIR=$(pwd) go test ./cmd/internal/tester -v
