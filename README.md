# golang-api

Build pre-requisites: Make, sh, Docker + Docker Compose

- To build:
    - `make`
- To test:
    - `make test`
- To build, test, and integration test:
    - `make full`

## Testing

[![Build Status](https://travis-ci.com/samgiles/golang-api.svg?branch=master)](https://travis-ci.com/samgiles/golang-api)

On Travis we spin up an instance of the test process, and a postgres instance.
The full test suite, including high level integration/behaviour tests are then
run against the database instance.
See `build/integration/int-test.sh` for details.  

## Deployment notes

In `build/prod` there is a Dockerfile that can be used to build an image to
deploy.  The image is a [FROM scratch image](https://hub.docker.com/_/scratch/  ), and only contains the binary!
At the moment the prod image is only about 8.23MB. Woohoo!

The listener's port is defined by the `PORT` environment variable.  If the
`PORT` environment variable is empty, then `3000` is used as a default.

The server has two healthcheck routes, `/__readiness` and `__liveness`, these
are intended to match [Kubernetes' definitions and semantics of readiness and
liveness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/). Healthchecks are asynchronous, querying a healthcheck endpoint will
result only in a read of a cached healthcheck result. See the library I wrote to support async healthchecks with readiness and liveness semantics https://github.com/samgiles/health.

Postgres should be configured through environment variables: `DB_NAME`,
`DB_USERNAME`, `DB_PASSWORD`, and `DB_HOST`.  At application startup database
migrations are executed. We could use a readiness check to ensure that traffic
only reaches an instance of the server if the database we connect to has the
correct version for the binary. But details like this really depend on the
deployment pipeline.

If an instance is unable to check database migrations or connect it will not start.
