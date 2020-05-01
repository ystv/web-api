FROM golang:1.14-alpine AS build
WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/api

FROM scratch
COPY --from=build /bin/api /bin/api
ENTRYPOINT ["/bin/api"]
