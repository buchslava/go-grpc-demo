# go-grpc-demo

gRPC server along with getaway Mux for http proxy in Go with GWT authorization

## Install

```sh
go mod download
```

## Run server

```sh
cd cmd/server
./run
```

## Run under Docker

Stop Postgres if it's already installed: `sudo systemctl stop postgresql`

### Environment

Put the following `.env` file into the root of the project:

```
POSTGRES_USER=...
POSTGRES_PASSWORD=...
POSTGRES_DB=...
```

### Start & Check

```sh
docker-compose up --build
docker image ls
docker ps

# enter to the API server
docker exec -it <container id> /bin/sh
# ...
exit
# enter to the DB server
docker exec -it <postgres container id> /bin/sh
# Run PG console
psql postgresql://user:pwd@database:5432/db
select * from users;
# Run Auth test in a separate console and:
select * from users;
# you should see test2@email.com profile in the list
exit
```

## Auth test

```bash
# get the true token
AUTH_DATA_1=$(curl http://localhost:8080/auth?role=admin -H "Accept: application/json")
TOKEN_1=$(echo "$AUTH_DATA_1" | gawk '{ match($0, /:"(.+)"/, arr); if(arr[1] != "") print arr[1] }')

# insert test2@email.com profile
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN_1" -d '{"user": {"email": "test2@email.com", "name": "Test 2", "password": "111"}}' http://localhost:8080/users
# get the list of profiles
curl http://localhost:8080/users -H "Accept: application/json" -H "Authorization: Bearer $TOKEN_1"

# get the false token and an attempt to read the records
AUTH_DATA_2=$(curl http://localhost:8080/auth?role=user -H "Accept: application/json")
TOKEN_2=$(echo "$AUTH_DATA_2" | gawk '{ match($0, /:"(.+)"/, arr); if(arr[1] != "") print arr[1] }')
curl http://localhost:8080/users -H "Accept: application/json" -H "Authorization: Bearer $TOKEN_2"
```
