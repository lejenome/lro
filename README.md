Long-Running Operations (LRO) Service
=====================================

Implementation of a Long-Running Operations (LRO) Service on Golang.

## Directories:
- `pkg`: contains the core code used by different microservices
- `services`: contains the internal implementation of diffrent microservices.
- `spec`: contains the specification files of the REST API (OpenAPI), the
		Pub/Sub topics (AsyncAPI) and the data structures (JsonSchema).
- `dockerfiles`: contains the different dockerfiles and other config files
		needed to build and run the containers defined on `docker-compose.yaml`
- `builds`: various scripts and tools used on the dev, build and test process
