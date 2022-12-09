# datasink

A simple Go package that acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database (ex: InfluxDB).

This project uses [mystique](https://github.com/TheThingsIndustries/mystique) for the MQTT server implementation.

## Usage

```
datasink is tool that acts as acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database. More documentation at https://krishnaiyer.dev/golang/datasink

Usage:
  datasink [flags]
  datasink [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Display config information
  help        Help about any command
  init-db     Initialize the database
  version     Display version information

Flags:
  -c, --config string                                              config file (Default; config.yml in the current directory) (default "./config.yml")
      --database.influxdb.address string                           server address
      --database.influxdb.bucket string                            data bucket
      --database.influxdb.non_blocking_writes.batch_size int       batch size
      --database.influxdb.non_blocking_writes.enabled              enable non-blocking writes
      --database.influxdb.non_blocking_writes.flush_interval int   flush interval
      --database.influxdb.organization string                      organization
      --database.influxdb.setup.password string                    password
      --database.influxdb.setup.retention_period_hrs int           retention period in hours
      --database.influxdb.setup.username string                    username
      --database.influxdb.token string                             auth token. Generate a random one using 'openssl rand -hex 32'
      --database.influxdb.write_timeout int                        write timeout in seconds (for blocking writes)
      --database.type string                                       The type of database to use. Supported values are 'influxdb'
      --devices.smart-meter.values strings                         Values to record and the corresponding data type
  -h, --help                                                       help for datasink
      --http.address string                                        server address
      --mqtt.address string                                        server address
      --mqtt.allowed-topic-prefix strings                          allowed topic prefix per username
      --mqtt.auth.htpasswd-file string                             location of the htpasswd file
      --mqtt.auth.type string                                      authentication file type. Supported values are 'htpasswd'
      --mqtt.debug                                                 enable debug mode

Use "datasink [command] --help" for more information about a command.
```

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
