# Migration guide
Version v0 of this API is no longer supported and is not fully functional due to changes in the official BDO website. This guide explains how to migrate an application designed to work with v0 of this API to using v1 instead.

## API paths
All available API paths got renamed to follow the REST standards more closely:

| Previously (v0)        | Now (v1)              |
| ---------------------- | --------------------- |
| /v0/profile            | /v1/adventurer        |
| /v0/profileSearch      | /v1/adventurer/search |
| /v0/guildProfile       | /v1/guild             |
| /v0/guildProfileSearch | /v1/guild/search      |

## Redundant data
There was some tautology in the v0 responses. Hence, some attributes got removed/renamed in v1:

### GET /v1/adventurer
- `response.guild.region` got removed, use `response.region` instead.

### GET /v1/adventurer/search
- `response[i].guild.region` got removed, use `response[i].region` instead.

### GET /v1/guild
- `response.guildMaster.region` got removed, use `response.region` instead.
- `response.guildMaster` got renamed into `response.master`.

### GET /v1/guild/search
- `response[i].guildMaster.region` got removed, use `response[i].region` instead.
- `response[i].guildMaster` got renamed into `response[i].master`.

## Error handling
v0 didn't have a single error communication standard, so applications using it didn't have a reliable way to handle them. With v1, you can always rely on the response [HTTP status code](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status). Codes that are currently in use:
| Status                  | Meaning                                                    |
| ----------------------- | ---------------------------------------------------------- |
| 200 OK                  | Enjoy your data                                            |
| 400 Bad Request         | Some required request parameters are missing               |
| 404 Not Found           | The data you requested is not available on the BDO website |
| 503 Service Unavailable | BDO servers are currently under maintenance                |
