# BDO-REST-API
A collector for Black Desert Online player in-game data that provides an unofficial REST API.

## Projects using this API
- [BDO Leaderboards](https://bdo.hemlo.cc/leaderboards/) ([GitHub](https://github.com/octoman90/BDO-Leaderboards)): web-based leaderboards for Black Desert Online.

## How to start using it
There are two ways to use this API in your apps:
1. https://bdo-rest-api.herokuapp.com/v1 is an instance that I host. Keep in mind that it may be slow to respond, and might just stop working one day. It's also rolling-release: if the API changes in the master branch on GitHub, this instance will reflect it immediately and your app may break. The API documentation can be viewed [here](https://gitlab.com/man90/black-desert-social-rest-api/-/tree/master/doc/api/openapi.json).
2. Host it yourself. The server may cost some money, but the process is trivial. This approach will give you more stability and freedom. There are four easy steps to it:
	1. Build the server from the source code following [this guide](doc/buildingFromSource.md) or download a prebuilt Linux binary from [here](https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines).
	2. Set the environment variables if you want. The list is in a section below.
	3. Run the binary. Possible flags are described in a section below.
	4. Use the API as described in the [documentation](https://gitlab.com/man90/black-desert-social-rest-api/-/tree/master/doc/api/openapi.json).

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
These flags override environment variables.
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
There is a number of bugs that the official BDO website has. This collector does not do anything about them for the sake of simplicity, so your apps may need to use the [workarounds](doc/brokenStuff.md).

## API migration
There used to be a v0 version that is no longer supported, so if you use it, use the  [migration guide](doc/migrationGuilde.md).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](LICENSE)
