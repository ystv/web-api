# YSTV web-api

A Go based backend that should be able to handle website queries? Hopefully having a supportive subroutine to keep everything in order. Designed kind of like a monolith, but we'll see where we get to with it.

Generation of usable JWT's is available in [web-auth](https://github.com/ystv/web-auth) currently.

## Functionality

### REST API

- [ ] API for public website
  - [ ] VOD
    - [x] Videos
    - [x] Series
    - [ ] Playlists
    - [ ] Search
    - [x] Breadcrumbs
    - [x] Path to video
    - [x] Path to series
    - [ ] Thumbnail inheritance
  - [ ] Live
    - [x] Streams (currently static)
  - [x] Teams
    - [x] Get current team
    - [x] List teams
    - [x] Get a team by year
    - [x] List all officers
- [ ] Internal (Secured)
  - [ ] Creator Studio
    - [ ] Videos
      - [x] Listing
      - [x] Getting with files
      - [x] Uploading
      - [ ] Updating
      - [ ] Deleting
    - [x] Series (very experimental)
    - [ ] Playlists
      - [x] List meta
      - [x] Create meta
      - [x] Update meta
      - [ ] Create video list
      - [ ] Update video list
    - [ ] Encoding
      - [ ] Presets
        - [x] Get
        - [x] Create (experimental)
        - [x] Update (experimental)
        - [ ] Delete
      - [ ] Formats
        - [x] Get
        - [ ] Create
        - [ ] Update
        - [ ] Delete
    - [x] Calendar
    - [x] Stats
  - [ ] People
    - [x] User by ID
    - [x] User by JWT
    - [x] Permissions
    - [ ] Create user
    - [ ] Create roles
    - [x] List all users
    - [x] List all users by role
  - [ ] Clapper
    - [ ] Events
      - [x] List by month
      - [ ] List by term
      - [ ] Signups
        - [x] Listing including positions
        - [x] Creating
        - [x] Updating
    - [ ] Positions
      - [x] List
      - [x] Create
      - [x] Update
      - [ ] Delete
  - [ ] Encoder
    - [x] Video upload auth hook
  - [x] Stream auth (experimental)
  - [ ] Misc internal services

### Services

- [ ] Encode management
- [x] Mailer

### Connections

- [x] Postgres
- [x] RabbitMQ
- [x] web-auth integration
- [ ] Authstack

## Dependencies

- Go (For developing only)
- A database (Although, it might be hardcoded for postgres?)
- A "CDN", I just mean a connection to an S3 like interface.
- AMQP compatible broker i.e. RabbitMQ.

## Installations

### Static binary

#### Building

`go get github.com/ystv/web-api`

if that doesn't work try `git clone https://github.com/ystv/web-api`

#### Installing

Run `go build -o web-api`, hopefully you've got the binary.

Copy the `.env` file from `configs` and place it in the root directory  
`cp configs/.env .`

Run `./web-api` and hopefully it should be running.

### Docker

Unlike the static binary this uses docker's environment variables instead of a `.env` file. The API's port will be exposed on port `8080`.

To start either method you will need to clone the repo.

`git clone https://github.com/ystv/web-api`

#### Dockerfile

Build the image.

`docker image build -t web-api .`

You will then need to set the environment variables like how it is set up in `docker-compose.yml` or you might be able to create a container and then use `docker export` so you can run a static binary built in docker.

#### Docker-compose

Copy the `docker-compose.yml` from `configs` and place it in the root directory.  
`mv configs/docker-compose.yml .`

Then fill in the docker-compose file with your credentials.

`docker-compose up -d --build`

### Jenkins

Checkout this [document](ci.md).

## Developing

Developed on Go version 1.14+

Clone the repo and create a `.env.local` for an easy config setup. Use the `debug` flag to disable auth.

I recommend not using the production environment for testing and recommend running postgres and rabbitmq in Docker.

Updating the DB schema use `goose` to migrate safely.

When ran with the `debug` flag set to true. 500 server errors will be returned to the browser, not just logged. Otherwise, it will only return the 500 code and not the actual error.

Try to keep all "business" logic with in the `/services` and try to keep the imports local to that package, but you'll probably need the `utils` package, but we're trying to keep it modular so web-api can theoretically be split up and keeping it separate would likely make it a lot easier.

Generate docs by running `go generate` in the root of the repo.

### Layout

> For further reading checkout the [architecture](architecture.md) document.

- `/configs`  
  Example configurations.
- `/controllers`  
  Internal API to REST API.
- `/docs`  
  echo-swagger autogenerated docs.
- `/middleware`  
  REST middleware, implements some URL sanitization and logging.
- `/routes`  
  Handles how the REST api is arranged and applies extra middleware if necessary.
- `/services`  
  Internal API that is the main "business logic" for each service web-api provides.
- `/utils`  
  Provides access to commonly used functions, i.e. cdn, mq, sql

### Database info

This is currently build using the new schema available in the [planning repo](https://github.com/ystv/Website2ElectricBoogaloo), it is still lacking functionality compared to the current implementation.
