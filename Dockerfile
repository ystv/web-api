FROM golang:1.15-alpine AS build
LABEL site="api"
LABEL stage="builder"
WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/api cmd/main.go

FROM scratch
LABEL site="api"
COPY --from=build /bin/api /bin/api
ENTRYPOINT ["/bin/api"]
