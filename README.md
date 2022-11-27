# datasink

A simple Go package that acts as a server with multiple protocols (ex: mqtt, websocket) for incoming traffic and writes to a time series database (ex: InfluxDB).

This project uses [mystique](https://github.com/TheThingsIndustries/mystique) for the MQTT server implementation.

## Usage


### Auth

Create an htpasswd file.
```
$ htpasswd -c test.htpasswd test
```

## License

The contents of this repository are packaged under the terms of the [Apache 2.0 License](./LICENSE).
