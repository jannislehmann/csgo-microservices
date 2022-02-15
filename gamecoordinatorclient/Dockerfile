FROM golang:alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/gamecoordinatorclient
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc
COPY Makefile go.mod go.sum ./
RUN make init && go mod download
COPY . .
RUN make proto tidy build

FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/gamecoordinatorclient/gamecoordinatorclient /gamecoordinatorclient
# See https://github.com/Philipp15b/go-steam/issues/124 for further information
ENV GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn
ENTRYPOINT ["/gamecoordinatorclient"]
CMD []
