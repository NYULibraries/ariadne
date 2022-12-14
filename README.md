# Ariadne

## Overview

The NYU Libraries [OpenURL](https://en.wikipedia.org/wiki/OpenURL) link resolver

## Application architecture

* [backend/](backend/README.md): API server written in Go
* [frontend/](frontend/README.md): client SPA built with React

## Deployment

TODO: architecture and workflows, deployment examples, etc.

## E2E tests

Run tests:

```
cd e2e/
yarn test:e2e
```

Run E2E tests in a container:

```
docker-compose run --rm e2e
```
