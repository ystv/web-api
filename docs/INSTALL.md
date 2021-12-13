# Installation

First, download the repository: `git clone https://github.com/ystv/web-api`

Alternatively, you can use go to download the code: `go get github.com/ystv/web-api`

Once you have the code, there are multiple options for running it locally.

## 1) Static binary

Run `go build -o web-api`, hopefully you've got the binary.

Copy the `.env` file from `configs` and place it in the root directory:
`cp configs/.env .env.local`

Fill in this file as appropriate (for YSTV use internal credentials)

Run `./web-api` and hopefully it should be running.

To fully build the API documentation, do `swag init --pd -o swagger/`.
You may have to use something like `~/go/bin/swag` instead, after installing it from [its github repo](github.com/swaggo/swag).

## 2) Docker

Unlike the static binary, this uses docker's environment variables instead of an `.env` file.
The API's port will be exposed on port `8080`.

### Dockerfile

Build the image: `docker image build -t web-api .`

You will then need to set the environment variables like how it is setup in `docker-compose.yml`.
Or you might be able to create a container and then use `docker export` so you can run a static binary built in docker.

### Docker-compose

Copy the `docker-compose.yml` from `configs` and place it in the root directory.  
`cp configs/docker-compose.yml .`

Fill in the docker-compose file with your credentials.

Then run docker compose: `docker-compose up -d --build`
