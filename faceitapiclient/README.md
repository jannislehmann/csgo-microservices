# Faceit API client microservice

This api consuming microservice has a set of users with match history authentication code and queries Faceit's API to get the latest downloadable matches.
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

### Faceit

This service requires a Faceit API key for requesting the match list and details.

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `apiKey` | `secret` | The Faceit API key. Can be generate [here](https://developers.faceit.com) |

## Dependencies

This service requires a Postgres database.
