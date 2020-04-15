# YSTV web-api

A Go based backend that should be able to handle website queries?

This is currently built to handle the majority of the tables in `ystv`. If there is a future table needed, you'll probably need to re-run [sqlboiler](https://github.com/volatiletech/sqlboiler) so it can generate the new models.

## Dependencies

- Golang (Although if it's compiled maybe not?)
- A database (Although, it might be hardcoded for postgres?)

# Installing

`go get github.com/ystv/web-api`

Fill in the `.env` with the database details.

Run `go build`, hopefully you've got the binary.
