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

### Updating jest snapshots

After stopping all previous frontend-test containers:

```
docker rm $(docker-compose ps -aq frontend-test )
```

We can recreate jest snapshots in the frontend-test container:

```
docker-compose run frontend-test yarn test -u
```

Now that we only have one frontend-test container, stopped, with the updated snapshot, we can copy the snapshots to our local:

```
docker cp "$(docker-compose ps -aq frontend-test)":/app/src/components/ frontend/src/
```

Or just do it all in one line:

```
docker rm $(docker-compose ps -aq frontend-test) && docker-compose run frontend-test yarn test -u && docker cp "$(docker-compose ps -aq frontend-test)":/app/src/components/ frontend/src/
```
