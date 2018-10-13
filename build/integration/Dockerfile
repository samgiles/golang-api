FROM golang:1.11-alpine as builder
RUN mkdir /src
ADD . /src
WORKDIR /src
ENV GOPATH=/:$GOPATH
RUN CGO_ENABLED=0 GOOS=linux go test -c -a -installsuffix cgo -ldflags '-extldflags "-static"' ./cmd/...
FROM alpine as test
RUN mkdir -p /app/cmd/server
COPY --from=builder /src/test/ /app/test/
COPY --from=builder /src/server.test /app/cmd/server/
WORKDIR /app/cmd/server
CMD ["./server.test", "-test.run=Integration", "-test.v"]
