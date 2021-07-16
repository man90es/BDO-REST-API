# Black Desert social REST API
An unofficial REST API server for Black Desert Online community data.

## Projects using this API
- [BDO Leaderboards](https://bdo.hemlo.cc/leaderboards/): a web-based leaderboard application for BDO guilds.

## API routes
The current version of API is v1. **v0 is no longer supported** and is not fully functional due to changes in the official BDO website.
- The migration guide can be found [here](doc/migrationGuilde.md).
- The v1 OpenAPI documentation can be found [here](https://gitlab.com/man90/black-desert-social-rest-api/-/tree/master/doc/api/openapi.json).

## Getting the binary
You can either download prebuilt binaries from [here](https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines) or build them yourself following [this guide](doc/buildingFromSource.md).

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
./bdo-rest-api -proxy="http://192.168.0.0.1:8080" -cachettl=30
```

## Known bugs
A vast majority of bugs comes from the original BDO website, where data is taken from. You can find a list of known bugs and workarounds [here](doc/brokenStuff.md).
