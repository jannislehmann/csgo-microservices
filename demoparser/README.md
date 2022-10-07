# Demo parser microservice

This service parses the demos, which are stored in the configured `demosDir`.

TODO: Add database

## Config

The service has the following extensions for the config.

### Parser

This service requires a directory to read demos from.

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `demosDir` | `/home/demo` | The directoy where the demos are stored. |
| `workerCount` | `5` | The amount of parallel workers to parse demos. |

### Database

This services uses a MongoDB, which requires permissions to create and write to collections.

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `host` | `localhost` | The database host. |
| `port` | `27017` | The database port. |
| `username` | `csgo` | Username of the database user. |
| `password` | `secret` | Secret password of the database user. |
| `database` | `csgo` | The database name to store the data in. |
