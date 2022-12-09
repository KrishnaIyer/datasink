# datasink

A simple Go package that acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database (ex: InfluxDB).

This project uses [mystique](https://github.com/TheThingsIndustries/mystique) for the MQTT server implementation.

## Usage

1. Create a configuration file based on the [provided default](./config.yml).

2. Pull Docker images

```bash
$ docker-compose pull
```

3. Start the InfluxDB container.

The example [docker-compose.yml](./docker-compose.yml) file starts InfluxDB with some defaults. Change these values if necessary and run the container.

```bash
$ docker-compose up -d influxdb
```

4. Initialize the database.

```bash
$ docker-compose run datasink datasink -c /etc/config.yml init-db
```

5. Create an `htpasswd` file with the MQTT login credentials.

This process is different for each OS so look for the solution online. For unix-based OS, `htpasswd` is usually already installed.

```bash
$ htpasswd -c test.htpasswd <username>
```

6. Start the containers

```bash
$ docker-compose up -d
```

7. Connect a device via MQTT and send some data. The access credentials are the same as what goes into the `htpasswd` file.

8. Login to Grafana at http://localhost:3000. This assumes the default configuration. If using a different port, that should reflect here.

9. Add InfluxDB as a data source and use the `flux` option. For more details, check the [grafana docs](https://grafana.com/docs/grafana/latest/datasources/influxdb/).

10. At this point, you should be able to create InfluxDB (flux) queries on your measurements.

## Development

1. Clone this repository.

2. Initialize it locally.

```bash
$ make init
```

## License

The contents of this repository are packaged under the terms of the [Apache 2.0 License](./LICENSE).
