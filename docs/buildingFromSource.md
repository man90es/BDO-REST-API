# Building from source
Building the scraper from source may be preferable in some cases. This way, the app has a smaller footprint than a [Docker image](https://hub.docker.com/r/man90/bdo-rest-api).

## Prerequisites:
- GNU/Linux (other platforms should work as well but I haven't tested them)
- Go compiler >=v1.15

## Compilation
By default, scraped results are cached in memory and stored for 3 hours. This helps to ease the pressure on BDO servers and decreases the response time tremendously (for cached responses). Use this command to compile the app:
```bash
go build
```

If you don't want to cache scraped results (e.g., if you are 100% sure that there will be no similar requests sent to the API), you can also use this command instead:
```bash
go build -tags "cacheless"
```
