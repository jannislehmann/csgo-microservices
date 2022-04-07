# Shared

This module contains shared utility code, which is used
by all the other microservices.

## API client

The api client packages contains a small http client wrapper to ease the testing of the API consumers.

## Config

The config package uses [viper](https://github.com/spf13/viper) to load an application config using either a json file or env variables.
It also contains a global configuration struct, which is most-likely to be used by every service. However, the config will be converted to a struct
that has to be passed to the service. Thus, each service can define its own struct and therefore require different configuration variables.

### Global configuration

The global configuration of each project has the following properties:

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `logLevel` | `info` | The log level can be either set to `info`, `error`, `debug` or `trace`. |

### Broker configuration

Each application uses the AMQP broker for communication. The following configuration parameters are available:

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `uri` | `amqp://guest:guest@localhost:5672/` |  The connection uri for the broker. |

## Entity

The entity package contains an ID type, which is basically a UUID.

## Metrics

The metrics package offers helper methods for the creation of metrics.
Currently, it has a method to create a Prometheus http server.

## Queue

The queue package is able to connect to a RabbitMQ AMQP broker and publish messages or consume these.

## ShareCode

The share code package contains the logic to decode CSGO share codes to their match id, outcome id and token.

## Util

This package contains smaller util functions.
