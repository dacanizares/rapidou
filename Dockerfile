FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY src ./src
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /rapidou ./src

FROM golang:1.26-bookworm AS functional-test
RUN apt-get update \
    && apt-get install -y --no-install-recommends chromium ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY src ./src
COPY tst ./tst
CMD ["go", "test", "./...", "-v", "-count=1"]

FROM ubuntu:24.04
RUN useradd --system --uid 10001 --create-home rapidou && mkdir /data && chown rapidou:rapidou /data
USER rapidou
COPY --from=build /rapidou /usr/local/bin/rapidou
ENV APP_ADDRESS=:8080 APP_DATABASE=/data/app.db
EXPOSE 8080
ENTRYPOINT ["rapidou"]
