FROM golang:1.11-alpine as builder
RUN mkdir /src
ADD . /src
WORKDIR /src
ENV GOPATH=/:$GOPATH
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' ./...
FROM scratch
COPY --from=builder /src/server /app/
WORKDIR /app
CMD ["./server"]
