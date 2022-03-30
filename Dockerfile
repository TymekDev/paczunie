FROM golang:1.17 as build-env
WORKDIR /src
COPY . .
# Reference: https://7thzero.com/blog/golang-w-sqlite3-docker-scratch-image
RUN CGO_ENABLED=1 go build -ldflags '-linkmode external -extldflags "-static"' -o /app

FROM scratch
COPY --from=build-env /app /app
COPY packages.db /
CMD [ "/app" ]
