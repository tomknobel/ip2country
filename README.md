# ip2country

## Configuration

all configuration can can be found [.env](./.env).
for simplicity, I only configure csv db wit connection string 
but the services can be extract to more db ( need to add username password etc to config )
the connection string in the [.env](./.env) is set for running in docker if one want to run the services without
docker container need to change the connection string tio full path of the cvs db now locate
[here](./dataBase/ipToCountryDB.csv)

## Running the service 
can go build and go run main just need to change the db connection string to full path

or

can run 
```bash
# Build the service

docker-compose build
# Run the service
docker-compose up
```
the [Dockerfile](./Dockerfile) is multi-stage Dockerfile and can be run directly by docker commend

### Unit test

unit test can be run 
```bash
# By docker / docker-compose
docker-compose run --rm test

# locally 
go test ./... -v

```

Because of time constraints I only did unit test for one package in this service :( 
 ### example 
working query  with ip that persist on this small db 
v1/find-country?ip=1.0.0.255