package main

import (
	"context"
	"log"
	"os"
	"net"
	"fmt"
	"encoding/json"

	"google.golang.org/grpc"
	pb "github.com/kikeyama/grpc-sfx-demo/pb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.mongodb.org/mongo-driver/mongo/readpref"

	grpctrace "github.com/signalfx/signalfx-go-tracing/contrib/google.golang.org/grpc"
	mongotrace "github.com/signalfx/signalfx-go-tracing/contrib/mongodb/mongo-go-driver/mongo"
	"github.com/signalfx/signalfx-go-tracing/tracing"
)

//var logger log.Logger
var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

var collection *mongo.Collection
var client *mongo.Client

const (
	port = ":50051"
	serviceName = "kikeyama_grpc_server"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); !exists {
		return defaultVal
	} else {
		return value
	}
}

func connectMongo() error {
	// MongoDB
	mongoHost := getEnv("MONGO_HOST", "localhost")
	opts := options.Client()

	ctx := context.TODO()
	
	// SignalFx Instrumentation
	opts.SetMonitor(mongotrace.NewMonitor(mongotrace.WithServiceName("kikeyama_mongo")))

	client, err := mongo.NewClient(opts.ApplyURI(fmt.Sprintf("mongodb://%s:27017", mongoHost)))
	if err != nil {
		logger.Fatalf("level=fatal message=\"failed to open client %v\"", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatalf("level=fatal message=\"failed to connect mongodb %v\"", err)
	}
	collection = client.Database("test").Collection("animals")

	return nil
}

func listAnimals(ctx context.Context, in *pb.EmptyRequest) (*pb.Animals, error) {
	logger.Printf("level=info message=\"List Animals\"")

//	// MongoDB
//	mongoHost := getEnv("MONGO_HOST", "localhost")
//	opts := options.Client()
//
//	// SignalFx Instrumentation
//	opts.SetMonitor(mongotrace.NewMonitor(mongotrace.WithServiceName("kikeyama_mongo")))
//
//	client, err := mongo.NewClient(opts.ApplyURI(fmt.Sprintf("mongodb://%s:27017", mongoHost)))
//	if err != nil {
//		logger.Fatalf("level=fatal message=\"failed to open client %v\"", err)
//	}
////	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
////	defer cancel()
//	err = client.Connect(ctx)
//	if err != nil {
//		logger.Fatalf("level=fatal message=\"failed to connect mongodb %v\"", err)
//	}
//	defer func() {
//		if err = client.Disconnect(ctx); err != nil {
//			logger.Fatalf("level=fatal message=\"failed to disconnect mongodb: %v\"", err)
//		}
//	}()
//
//	collection := client.Database("test").Collection("animals")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		logger.Printf("level=error message=\"failed to find collection: %v\"", err)
	}
	defer cur.Close(ctx)

	var animals []*pb.AnimalInfo
//	var animal pb.AnimalInfo
//
//	var d bson.D
//
//	for cur.Next(ctx) {
//		var result bson.M
//		if err = cursor.Decode(&result); err != nil {
//			logger.Printf("level=error message=\"failed to decode cursor %v\"", err)
//		}
//		id, ok := result["_id"]
//		if ok {
//			result["id"] = id.String()
//		} else {
//			result["id"] = ""
//		}
//		d = append(d, result)
//	}

	err = cur.All(ctx, &animals)
	animalsJson, err := json.Marshal(animals)
//	animalsJson, err := json.Marshal(d)
	if err != nil {
		logger.Printf("level=error message=\"unable to marshall animals to json\"")
	}
	logger.Printf("level=info message=\"retrieve list from mongodb\" data=%s", string(animalsJson))

	// gRPC response
	return &pb.Animals{Animals: animals}, nil
}

//func getAnimal(ctx context.Context, in *pb.AnimalId) (*pb.AnimalInfo, error) {
//	logger.Printf("level=info message=\"Get Animal\"")
//	return &pb.AnimalInfo{Id: in.GetId(), Animal: in.getAnimal()}, nil
//}

func main() {
	// Create a listener for the server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("level=fatal message=\"failed to listen: %v\"", err)
	}

	// Use signalfx tracing
	tracing.Start(tracing.WithGlobalTag("stage", "demo"), tracing.WithServiceName(serviceName))
//	tracing.Start()
	defer tracing.Stop()

	err = connectMongo()

	defer func() {
		ctx := context.Background()
		if err = client.Disconnect(ctx); err != nil {
			logger.Fatalf("level=fatal message=\"failed to disconnect mongodb: %v\"", err)
		}
	}()

	// Create the server interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(serviceName))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(serviceName))

	// Initialize the grpc server as normal, using the tracing interceptor.
	//s := grpc.NewServer()
	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))

	pb.RegisterAnimalServiceService(s, &pb.AnimalServiceService{ListAnimals: listAnimals})
	if err = s.Serve(lis); err != nil {
		logger.Fatalf("level=fatal message=\"failed to serve at AnimalService.ListAnimals: %v\"", err)
	}
}
