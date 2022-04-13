#!/bin/bash

function gen() {
  oapi-codegen -config oapi-codegen.yml -package $1 -o $2 $3
}

# Components
gen auth ../pkg/apis/auth/auth.gen.go ../spec/components/auth.yml
gen process ../pkg/apis/process/process.gen.go ../spec/components/process.yml
gen common ../pkg/apis/common/common.gen.go ../spec/components/common.yml

# Endpoints
gen process ../services/process-api/controllers.gen.go ../spec/process-api.openapi.yml
