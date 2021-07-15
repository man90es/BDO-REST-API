# Migration guide
Version v0 of this API is no longer supported and is not fully funcional due to changes in the official BDO website. This guide explains how to migrate an application that was designed to work with v0 of this API to using v1 instead.

## API paths
All available API paths got renamed to follow the REST standards more closely:

| Previously (v0)        | New (v1)              |
| ---------------------- | --------------------- |
| /v0/profile            | /v1/adventurer        |
| /v0/profileSearch      | /v1/adventurer/search |
| /v0/guildProfile       | /v1/guild             |
| /v0/guildProfileSearch | /v1/guild/search      |

## Redundant data
Some data was included in responses more than once in v0. Hence, some attributes got removed in v1:

### GET /v1/adventurer
- If you were accessing `response.guild.region`, replace it with `response.region`.

### GET /v1/adventurer/search
- If you were accessing `response[i].guild.region`, replace it with `response[i].region`.

### GET /v1/guild
- If you were accessing `response.guildMaster.region`, replace it with `response.region`.
- If you were accessing `response.guildMaster`, replace it with `response.master`.

### GET /v1/guild/search
- If you were accessing `response[i].guildMaster.region`, replace it with `response.region`.
- If you were accessing `response[i].guildMaster`, replace it with `response.master`.

## Error handling
v0 didn't have a single error standard, so applications using it didn't have a reliable way to handle them. With v1, you can always rely on the response [HTTP status code](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status). Codes that are currently in use:
| Status                  | Meaning                                                    |
| ----------------------- | ---------------------------------------------------------- |
| 200 OK                  | Enjoy your data                                            |
| 400 Bad Request         | Some required request parameters are missing               |
| 404 Not Found           | The data you requested is not available on the BDO website |
| 503 Service Unavailable | BDO servers are currently under maintenance                |
