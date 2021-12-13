# YSTV web-api

A Go based backend that should be able to handle website queries.
Hopefully having a supportive subroutine to keep everything in order.
Designed kind of like a monolith but we'll see where we get to with it.

Generation of usable JWT's is available in [web-auth](github.com/ystv/web-auth) currently.

## Dependencies

- Go (For developing only)
- A database (Postgres may be the hardcoded requirement)
- A "CDN" (a connection to an S3 like interface)
- An AMQP compatible broker (RabbitMQ_).

## Installing
See the [install guide](./docs/INSTALL.md).

## Developing
See the [development guide](./docs/DEVELOP.md).
