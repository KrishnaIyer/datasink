# datasink

A simple Go package that acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database (ex: InfluxDB).

This project uses [mystique](https://github.com/TheThingsIndustries/mystique) for the MQTT server implementation.

## Usage

1. Clone this repository.

2. Initialize it locally.

```bash
$ make init
```
3. The example [docker-compose.yml](./docker-compose.yml) file starts Influx with some defaults. Change these values if necessary and run the container.

```bash
$ docker-compose up -d influxdb
```

4. Extract the Auth token required for the Go InfluxDB client to connect to the database.

```bash
$ docker exec datasink-influxdb-1 influx auth list | awk '/testuser/ {print $4 " "}'
```
> Note: Change the username (`testuser`) to the required value.

This outputs a token similar to the following
```
-siG5kFijoSri_J9h8NundRx1rXWbwHTAjNv2Mc07-PPZGo1rZw6a9pzY0yEPDvfncbkqwxS_X4I0xJOBejs9Q==
```

### Auth

Create an htpasswd file.
```
$ htpasswd -c test.htpasswd test
```

## References

1. Setting up InfluxDB via Docker: https://medium.com/geekculture/deploying-influxdb-2-0-using-docker-6334ced65b6c
2. Getting started with InfluxDB and Go: https://www.influxdata.com/blog/getting-started-with-the-influxdb-go-client/

## License

The contents of this repository are packaged under the terms of the [Apache 2.0 License](./LICENSE).
