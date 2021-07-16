# Building from source
Prerequisites: <abbr title="Not tested on other platforms.">GNU/Linux</abbr>, Go >=1.15

Command:
```bash
go build
```

By default, scraped results are cached in memory and stored for up to 2 hours. It helps to ease the pressure on BDO servers and speeds up the response time in some situations. If you don't want to cache scraped results (e.g., if you only want to create a dump or need the most fresh data), use this command instead:
```bash
go build -tags "cacheless"
```
