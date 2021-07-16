# Black Desert social REST API

An unofficial JSON API server for Black Desert Online that gets guild and player data via scraping the official website.

## API
See [OpenAPI documentation](https://gitlab.com/man90/black-desert-social-rest-api/-/tree/master/doc/api/openapi.json).

## Projects using this API
- [BDO Leaderboards](https://man90.gitlab.io/bdo-leader-boards): a web-based leaderboard application for Black Desert Online guilds.



## Getting the binary
### Prebuilt
You can download prebuilt binaries [here](https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines).

### Building from source
Prerequisites: <abbr title="Not tested on other platforms.">GNU/Linux</abbr>, Go >=1.15

Command:
```bash
go build
```

By default, scraped results are cached in memory and stored for up to 2 hours, it helps to ease the pressure on BDO servers and speeds up the response time in some situations. If you don't want to cache scraped results (for example if you just want to create a dump or need extra fresh data), use this command instead:
```bash
go build -tags "cacheless"
```

## Environment variables
Catch requests on a specific port (8001 by default):
```bash
export PORT=3000
```

Use a proxy to make requests to BDO servers (none by default):
```bash
export PROXY=http://123.123.123.123:8080
# or
export PROXY="http://123.123.123.123:8080 http://124.124.124.124:8081"
```

## Flags
Flags override environment variables
```
-cachettl int
	Cache TTL in minutes (default 180)
-port int
	Port to catch requests on (default 8001)
-proxy string
	Open proxy address to make requests to BDO servers
```
Use them like this:
```bash
./black-desert-social-rest-api -proxy="http://192.168.0.0.1:8080" -cachettl=30
```

## Known bugs
A vast majority of bugs comes from the original BDO website, where data is taken from. You can find a list of known bugs and workarounds [here](doc/brokenStuff.md).
