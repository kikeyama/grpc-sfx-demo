package main

import (
	"context"
	"log"
	"os"
	"net"
	"fmt"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/kikeyama/grpc-sfx-demo/pb"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.mongodb.org/mongo-driver/mongo/readpref"
//	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

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

type AnimalInfo struct {
//	ID       primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Id       primitive.Binary  `bson:"_id,omitempty" json:"id"`
	Type     string    `bson:"type" json:"type"`
	Name     string    `bson:"name" json:"name"`
	Height   int32     `bson:"height" json:"height"`
	Weight   int32     `bson:"weight" json:"weight"`
	Region   []string  `bson:"region" json:"region"`
	IsCattle bool      `bson:"isCattle" json:"isCattle"`
}

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
		logger.Printf("level=error message=\"failed to open client %v\"", err)
		return err
	}
	err = client.Connect(ctx)
	if err != nil {
		logger.Printf("level=error message=\"failed to connect mongodb %v\"", err)
		return err
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
		logger.Printf("level=error message=\"collection.Find failed: %v\"", err)
	}
	defer cur.Close(ctx)

	var animals []*pb.AnimalInfo

	for cur.Next(ctx) {
		var result AnimalInfo
		var animal *pb.AnimalInfo

		if err = cur.Decode(&result); err != nil {
			logger.Printf("level=error message=\"failed to decode cursor: %v\"", err)
		}
//		resultJson, err := json.Marshal(result)
//		if err != nil {
//			logger.Printf("level=error message=\"unable to marshal animalinfo to json: %v\"", err)
//		}
//		logger.Printf(fmt.Sprintf("AnimalInfo: %s", string(resultJson)))
		id := result.Id
		animalUuid, err := uuid.FromBytes(id.Data)
		if err != nil {
			logger.Printf("level=error message=failed to parse UUID from bytes[]: %v", err)
		}
//		logger.Printf(fmt.Sprintf("id=%s", id.Hex()))
//		logger.Printf(fmt.Sprintf("id=%s", animalUuid.String()))
		animal = &pb.AnimalInfo{
//			Id:       id.Hex(),
			Id:       animalUuid.String(),
//			Type:     result.Type,
//			Name:     result.Name,
//			Height:   result.Height,
//			Weight:   result.Weight,
//			Region:   result.Region,
//			IsCattle: result.IsCattle,
		}
//		animalJson, err := json.Marshal(animal)
//		if err != nil {
//			logger.Printf("level=error message=unable to marshal animal to json: %v", err)
//		}
//		logger.Printf(fmt.Sprintf("pb.AnimalInfo: %s", string(animalJson)))
//		id, ok := result["_id"]
//		if ok {
//			animal["id"] = id.String()
//		} else {
//			animal["id"] = ""
//		}

		animals = append(animals, animal)
	}
//
//	animalsJson, err := json.Marshal(d)

	err = cur.All(ctx, &animals)
	if err != nil {
		logger.Printf("level=error message=\"unable to put data into animals: %v\"", err)
	}

	animalsJson, err := json.Marshal(animals)
	if err != nil {
		logger.Printf("level=error message=\"unable to marshall animals to json\"")
	}
	logger.Printf("level=info message=\"list animals from mongodb\" data=%s", string(animalsJson))

	// gRPC response
	return &pb.Animals{Animals: animals}, nil
}

//func (s *server) GetAnimal(ctx context.Context, in *pb.AnimalId) (*pb.AnimalInfo, error) {
func getAnimal(ctx context.Context, in *pb.AnimalId) (*pb.AnimalInfo, error) {
	logger.Printf(fmt.Sprintf("level=info message=\"Get Animal for id: %s\"", in.GetId()))

	animalUuid, err := uuid.Parse(in.GetId())
	id, err := animalUuid.MarshalBinary()
//	id, err := animalUuid.MarshalText()
	if err != nil {
		logger.Printf("level=error message=\"unable to parse uuid: %v\"", err)
	}

	var animal pb.AnimalInfo
	err = collection.FindOne(ctx, bson.M{"_id": primitive.Binary{
		Subtype: 0x04,
		Data:    id,
	}}).Decode(&animal)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Printf("level=info message=\"Document not found: %v\"", err)
			return &animal, status.Error(codes.NotFound, "document not found")
		}
		logger.Printf("level=error message\"failed to decode reuslt: %v\"", err)
	}

//	if err = res.Decode(&animal); err != nil {
//		logger.Printf("level=error message\"failed to decode reuslt:%v\"", err)
//	}
	animal.Id = in.GetId()

	animalJson, err := json.Marshal(animal)
	if err != nil {
		logger.Printf("level=error message=\"unable to marshall animal to json\"")
	}
	logger.Printf("level=info message=\"get animal from mongodb\" data=%s", string(animalJson))

	return &animal, nil
}

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

	// Connect MongoDB
	err = connectMongo()
	if err != nil {
		logger.Fatalf("level=fatal message=\"cannot connect to MongoDB: %v\"", err)
	}

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

//	pb.RegisterAnimalServiceService(s, &pb.AnimalServiceService{ListAnimals: listAnimals})
	pb.RegisterAnimalServiceService(s, &pb.AnimalServiceService{
		ListAnimals: listAnimals,
		GetAnimal: getAnimal,
	})
	if err = s.Serve(lis); err != nil {
		logger.Fatalf("level=fatal message=\"failed to serve at AnimalService.ListAnimals: %v\"", err)
	}
}
