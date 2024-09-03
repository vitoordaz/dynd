FROM golang:1.23.0-alpine3.20 as build
WORKDIR /src
RUN apk add --no-cache make
COPY cmd /src/cmd
COPY internal /src/internal
COPY go.* /src/
COPY Makefile /src/
RUN make vendor test build

FROM alpine
COPY --from=build /src/build/dynd /bin/dynd
ENTRYPOINT ["/bin/dynd"]
