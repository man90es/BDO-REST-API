# Black Desert social REST API

An unofficial REST API for Black Desert Online that scrapes guild and player data into convenient JSONs.

## Getting the binary
### Building from source

```bash
go build
```

### Prebuilt
[https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines](https://gitlab.com/man90/black-desert-social-rest-api/-/pipelines)

## Routes
### Guild data
Assuming that your guild's name is «TumblrGirls» and it's on the EU server:
```python
http://localhost:8001/v0/guildProfile?guildName=TumblrGirls&region=EU
```

### Player data
Assuming that your profileTarget is «reeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee» (you can get the profileTarget string from guild members data or from player search):
```python
http://localhost:8001/v0/profile?profileTarget=reeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
```

## To-do
* Guild search route
* Player search route
* Cache
* Error handling
