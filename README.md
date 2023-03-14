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

Run E2E tests in a container:

```
docker-compose run --rm e2e
```

Update golden files and *-linux.png screenshots in a container (won't update *-darwin.png screenshots):

```
docker-compose run e2e-update-screenshots
```

Run tests:

```
cd e2e/
yarn install
yarn test:e2e
```

Update golden files for E2E tests (run in _e2e/_):

```
UPDATE_GOLDEN_FILES=true yarn test:e2e
```
(Developer note: it's current not possible to use a custom flag like `--update-golden-files`
with `playwright`: [\[Feature\] Add support for test\.each / describe\.each \#7036](https://github.com/microsoft/playwright/issues/7036))


