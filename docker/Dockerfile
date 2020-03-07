FROM golang:1.13 AS build-environment

ARG GOARCH=amd64
ARG SEMVER=1.0.0

ADD . /files
WORKDIR /files

ENV GOOS=linux GOARCH=$GOARCH CGO_ENABLED=0

RUN make build SEMVER=$SEMVER BUILD_PATH=/build/goto build

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=build-environment /build/goto /app/goto
ENTRYPOINT [ "./goto" ]
