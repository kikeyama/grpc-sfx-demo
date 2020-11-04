module github.com/kikeyama/grpc-sfx-demo

go 1.14

require (
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/signalfx/golib v2.5.1+incompatible // indirect
	github.com/signalfx/signalfx-go-tracing v1.4.2
	github.com/tinylib/msgp v1.1.2 // indirect
	go.mongodb.org/mongo-driver v1.4.1
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/kikeyama/grpc-sfx-demo/pb => ./pb
