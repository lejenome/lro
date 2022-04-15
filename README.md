Long-Running Operations (LRO) Service
=====================================

Implementation of a Long-Running Operations (LRO) Service on Golang.

## Directories:
- `pkg`: contains the core code used by different microservices
- `services`: contains the internal implementation of diffrent microservices.

    - `process-api`: the implementation of Process API microservice, it will
				handles user facing REST API endpoints.

    - `process-executor`: the implementation of Process Executor microservice,
				it runs the scheduled jobs. Communication with the Process API
				microservice is atchived using Event-based architectures with NATS as a
				message broker, Redis as a shared in-memory key/value store and
				PostgreSQL as a persistent data store.
- `spec`: contains the specification files of the REST API (OpenAPI), the
		Pub/Sub topics (AsyncAPI) and the data structures (JsonSchema).
- `dockerfiles`: contains the different dockerfiles and other shared config files
		needed to build and run the containers defined on `docker-compose.yaml`.
- `configs`: contains the secret config files and envirement variables. To
		define and run your docker-compose envirement, copy the content of the
		folder to a subfolder `configs/dev` or `config/prod` and make the desired
		changes to the files.
- `builds`: various scripts and tools used on the dev, build and test process
