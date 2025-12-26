FROM golang:1.23.12-alpine3.22 as build
WORKDIR /src
RUN apk add --no-cache make
COPY cmd /src/cmd
COPY internal /src/internal
COPY go.* /src/
COPY Makefile /src/
RUN make vendor test build

FROM alpine:3.23.2
RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY --from=build /src/build/dynd /bin/dynd
ENTRYPOINT ["/bin/dynd"]
