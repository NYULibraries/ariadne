# Ariadne frontend

## Usage

Clone the repo and then:

```
cd frontend && yarn install
```

Start the client locally:

```
cd frontend && yarn start
```

Run in a container:

```
docker-compose up frontend
```

## Testing

Run all tests:

```
cd frontend/
yarn test
```

Run tests in a container:

```
docker-compose run --rm frontend-test
```
