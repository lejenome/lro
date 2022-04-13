#!/bin/bash

function genTypes() {
  echo "Generate for $3"
  oapi-codegen -config oapi-codegen.types.yml -package $1 -o $2 $3
}
function genAPI() {
  echo "Generate for $3"
  oapi-codegen -config oapi-codegen.api.yml -package $1 -o $2 $3
}

# Components
genTypes auth ../pkg/apis/auth/auth.gen.go ../spec/components/auth.yml
genTypes process ../pkg/apis/process/process.gen.go ../spec/components/process.yml
genTypes common ../pkg/apis/common/common.gen.go ../spec/components/common.yml

# Endpoints
genAPI process ../services/process-api/controllers.gen.go ../spec/process-api.openapi.yml
