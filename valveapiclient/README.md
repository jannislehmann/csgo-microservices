# Valve API client microservice

This api consuming microservice has a set of users with match history authentication code and queries Valve's API to get the latest match share codes.
These will then be published on the broker.

## Configuration

Copy the `config.json.example` in the `configs` dir and rename it to `config.json` in the same dir.

You can also use ENV vars to override single or set all configuration variables. The formatting for the configuration is as with the JSON configuration. The ENV base is `CSGO`. The Steam two factor secret turns into `STEAM_TWOFACTORSECRET`.

Details about the global configuration can be found [here](https://github.com/Cludch/csgo-microservices/blob/main/shared/README.md#config).

### Database

This services uses a Postgres database, which requires permissions to create tables and read/write data.

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `host` | `localhost` | The database host. |
| `port` | `5432` | The database port. |
| `username` | `csgo` | Username of the database user. |
| `password` | `secret` | Secret password of the database user. |
| `database` | `csgo` | The database name to store the data in. |

## Dependencies

This service requires a Postgres database.
