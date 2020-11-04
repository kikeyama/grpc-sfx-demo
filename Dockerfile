FROM golang:1.14 as build
ENV GO111MODULE on
WORKDIR /go/src/work
COPY main.go /go/src/work/
COPY go.mod /go/src/work/
COPY go.sum /go/src/work/
COPY pb/ /go/src/work/
#RUN go mod init
#RUN go mod edit -require github.com/opentracing/opentracing-go@v1.1.0
RUN CGO_ENABLED=0 go build -o /bin/grpc-sfx-demo

FROM scratch
COPY --from=build /bin/grpc-sfx-demo /bin/grpc-sfx-demo
COPY --from=busybox:1.32 /bin/busybox /bin/busybox
ENV GRPC_GO_RETRY on
CMD ["/bin/grpc-sfx-demo"]