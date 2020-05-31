# YSTV web-api

A Go based backend that should be able to handle website queries?

This is currently built to handle a few of the tables in `ystv`. If there is a future table needed, you'll probably need to re-run [sqlboiler](https://github.com/volatiletech/sqlboiler) so it can generate the new models. (You'll probably want to update `sqlboiler.toml` so your table isn't blacklisted)

## Dependencies

- Golang (Although if it's compiled maybe not?)
- A database (Although, it might be hardcoded for postgres?)
- A "CDN", I just mean a connection to an S3 like interface.

# Static binary

## Building

`go get github.com/ystv/web-api`

if that doesn't work try `git clone https://github.com/ystv/web-api`

## Installing

Run `go build -o web-api`, hopefully you've got the binary.

Copy the `.env` file from `configs` and place it in the root directory  
`mv configs/.env .`

Run `./web-api` and hopefully it should be running.

# Docker

Unlike the static binary this uses docker's environment variables instead of a `.env` file. The API's port will be exposed on port `8080`.

To start either method you will need to clone the repo.

`git clone https://github.com/ystv/web-api`

## Dockerfile

Build the image.

`docker image build -t web-api .`

You will then need to set the environment variables like how it is setup in `docker-compose.yml` or you might be able to create a container and then use `docker export` so you can run a static binary built in docker.

## Docker-compose

Copy the `docker-compose.yml` from `configs` and place it in the root directory.  
`mv configs/docker-compose.yml .`

Then fill in the docker-compose file with your credentials.

`docker-compose up -d --build`
