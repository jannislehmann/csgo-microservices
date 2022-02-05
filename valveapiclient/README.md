# Valve API client microservice

This api consuming microservice has a set of users with match history authentication code and queries Valve's API to get the latest match share codes.
These will then be published on the broker.

## Config

The service has the following extensions for the config.

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
