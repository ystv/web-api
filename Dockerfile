FROM golang:1.15-alpine AS build
LABEL site="api"
LABEL stage="builder"

WORKDIR /src/

# Stores our dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Download git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Set build variables
RUN echo -n "-X 'main.Version=$(git describe --abbrev=0)" > ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "' -X 'main.Commit=$(git log --format="%H" -n 1)" >> ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "'" >> ./ldflags

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="$(cat ./ldflags)" -o /bin/api

FROM scratch
LABEL site="api"
COPY --from=build /bin/api /bin/api
ENTRYPOINT ["/bin/api"]
