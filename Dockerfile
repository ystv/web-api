FROM golang:1.23.4-alpine3.21 AS build

LABEL site="api"
LABEL stage="builder"

# Create webapiuser.
ENV USER=webapiuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src/

# Stores our dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Generate documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate

# Download git
RUN apk update && apk upgrade && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

ARG WAPI_VERSION_ARG

# Set build variables
RUN echo -n "-X 'main.Version=$WAPI_VERSION_ARG" > ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "' -X 'main.Commit=$(git log --format="%H" -n 1)" >> ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "'" >> ./ldflags

RUN GOOS=linux GOARCH=amd64 go build -ldflags="$(cat ./ldflags)" -o /bin/api

FROM scratch
LABEL site="api"

# Import the user and group files from the builder.
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /bin/api /bin/api
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Use an unprivileged user.
USER webapiuser:webapiuser

ENTRYPOINT ["/bin/api"]
