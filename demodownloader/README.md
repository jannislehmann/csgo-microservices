# Demo downloader microservice

This service receives demo download details through the queue and downloads these.

## Config

The service has the following extensions for the config.

### Downloader

This service requires a directory to store the downloaded demos

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `demosDir` | `/home/demo` | The directoy where the downloaded demos should be stored. |
