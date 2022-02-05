# CSGO Microservice Suite

The CSGO tool suit can automatically detect and download new CSGO official matchmaking demos using the GameCoordinator.
In order to do this, a few API credentials and a **separate** Steam account is needed. The application creates a CSGO game sessions using the separate account
and uses the Steam Web API to check whether a new demo can be fetched. If that's the case, the application sends a full match info request to the game's GameCoordinator.
The GC then returns information about the match which also contain a download link.

## Services

Currently, the microservice suite consists of the following services.

### Valve API client

The API client consumes Valve's game history API and saves the game share codes in the database.
In order to add a new steam / csgo user, whose demos should be monitored, a user must be manually created in the database.

## Infrastructure

All services require a common RabbitMQ broker. Each service may define its own dependencies as well, which will be described in the project's README file.

Each deployable service will have its own `docker-compose` file and maybe a `.env` file, which needs to be updated with your own configuration.

## Building

Each project has a `Makefile` configuration. Commonly shared options are:

* `init` to initialize the local environment
* `proto` to generate protos
* `update` to update dependencies
* `tidy` to tidy the `go.mod` files
* `build` to build the binary
* `test` to test the project and generate coverage
* `docker` to build the Docker image
* `mock` to generate service mocks

## Testing

This project tries to have a high test coverage. Most services have an interface definition and also an implementing interface.
With the use of [mockery](https://github.com/vektra/mockery) each service can easily generate mock definitions for each service class.
Thus, unit tests are easy to write and should be implemented accordingly.

## Configuration

Copy the `config.json.example` of the service in the `config` dir and rename it to `config.json` in the same dir.

You can also use ENV vars to override single or set all configuration variables. The formatting for the configuration is as with the JSON configuration. The ENV base is `CSGO`. The Steam two factor secret turns into `STEAM_TWOFACTORSECRET`.

Details about the global configuration can be found [here](https://github.com/Cludch/csgo-microservices/blob/main/shared/README.md#config).
## CSGO Demo Tools

This project used to be a monolithic microservice. The old source code can be seen [here](https://github.com/Cludch/csgo-tools/). This repository also has a fully functional suite with all the services,
which are planned to be migrated into a stateless and indepentend microservice.

## Disclaimer

If you have suggestions please feel free to create an issue. This would help me a lot!

This tool is not affiliated with Valve Software or Steam.

## Other projects that helped me a lot

* [go-steam](https://github.com/Philipp15b/go-steam)
* [cs-go](https://github.com/Gacnt/cs-go)
* [go-dota2](https://github.com/paralin/go-dota2)
* [csgo](https://github.com/ValvePython/csgo)
* [Uber style guide](https://github.com/uber-go/guide/blob/master/style.md)
