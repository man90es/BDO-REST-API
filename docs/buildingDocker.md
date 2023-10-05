# Building Docker image
Building a Docker image yourself can be preferable over using a Docker Image from [Docker Hub](https://hub.docker.com/r/man90/bdo-rest-api) if you want to disable caching.

## Building an image
By default, scraped results are cached in memory and stored for up to 2 hours. This helps to ease the pressure on BDO servers and decreases the response time tremendously (for cached responses). Use this command to compile the app:
```bash
sudo docker build -t bdo-rest-api .
```

If you don't want to cache scraped results (e.g., if you are 100% sure that there will be no similar requests sent to the API), you can also use this command instead:
```bash
sudo docker build -t bdo-rest-api --build-arg tags=cacheless .
```

## Running a container
You can run a Docker container you just built by executing this command:
```bash
sudo docker container run -p 8001:8001 bdo-rest-api
```
