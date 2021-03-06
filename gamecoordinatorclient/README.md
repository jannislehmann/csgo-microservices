# Gamecoordinator Client

This service will consume messages from the broker, which indicate that a new share code / match was found.
With that, it will then query the CSGO / Valve Gamecoordinator to receive the download link of the match and publish this information
to the broker.

In order to scale this service, each service *must* have a different steam account to use.
Steam only allows one session per account to the csgo gamecoordinator / being logged in.

## Protobuf issue

To run the application and allow protobufs to be definied multiple times, use the following command:
`GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go run main.go`.

This should be fixed eventually to prevent unwanted side-effects. See the issue [here](https://github.com/Philipp15b/go-steam/issues/124).

## Config

The service has the following extensions for the config.

### Steam

This services requires a Steam account with the game to be able to talk to the Gamecoordinator.

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `username` | `csgo` | Username of the steam account. |
| `password` | `secret` | Secret password of the steam account. |
| `twoFactorSecret` | `xyz` | The 2FA secret of the steam account.  |

## Dependencies

This service works without external dependencies.
