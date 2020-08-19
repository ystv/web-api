FROM golang:1.15-alpine AS build
LABEL site="api"
LABEL stage="builder"
WORKDIR /src/
COPY . /src/

RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN echo -n "-X 'main.Version=" > ./ldflags
RUN git describe --abbrev=0 >> ./ldflags
RUN tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags 
RUN echo -n "' -X 'main.Commit=" >> ./ldflags
RUN git log --format="%H" -n 1 >> ./ldflags
RUN tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags
RUN echo -n "'" >> ./ldflags

RUN CGO_ENABLED=0 go build -ldflags="$(cat ./ldflags)" -o /bin/api cmd/main.go
RUN ls /bin

FROM scratch
LABEL site="api"
COPY --from=build /bin/api /bin/api
ENTRYPOINT ["/bin/api"]
