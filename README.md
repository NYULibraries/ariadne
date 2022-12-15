# Ariadne

## Overview

The NYU Libraries [OpenURL](https://en.wikipedia.org/wiki/OpenURL) link resolver

## Application architecture

* [backend/](backend/README.md): API server written in Go
* [frontend/](frontend/README.md): client SPA built with React

## Quickstart

To run a local instance of the Ariadne system, start the API server by following
the **Start the API server** instructions in [backend/README.md](backend/README.md),
and then build and start the frontend client by following the **Usage** instructions
in [frontend/README.md](frontend/README.md).

## Deployment

TODO: architecture and workflows, deployment examples, etc.

## E2E tests

Run tests:

```
cd e2e/
yarn install
yarn test:e2e
```

Run E2E tests in a container:

```
docker-compose run --rm e2e
```
