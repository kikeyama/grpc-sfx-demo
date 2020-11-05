package main

import (
	"context"
	"log"
	"os"
	"net"
	"fmt"
	"strconv"
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

	grpctrace "github.com/signalfx/signalfx-go-tracing/contrib/google.golang.org/grpc"
	mongotrace "github.com/signalfx/signalfx-go-tracing/contrib/mongodb/mongo-go-driver/mongo"
	"github.com/signalfx/signalfx-go-tracing/tracing"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

var collection *mongo.Collection
var client *mongo.Client

const (
	serviceName = "kikeyama_grpc_server"
)

type AnimalInfo struct {
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
	mongoPort := getEnv("MONGO_PORT", "27017")
	mongoDatabase := getEnv("MONGO_DATABASE", "test")
	mongoUser := getEnv("MONGO_USER", "appuser")
	mongoPassword := getEnv("MONGO_PASSWORD", "password")
	mongoAuthMechanism := getEnv("MONGO_AUTH_MECHANISM", "SCRAM-SHA-256")

//	opts := options.Client()
	// SignalFx Instrumentation
//	opts.SetMonitor(mongotrace.NewMonitor(mongotrace.WithServiceName("kikeyama_mongo")))
	opts := options.Client().SetMonitor(mongotrace.NewMonitor(mongotrace.WithServiceName("kikeyama_mongo")))

	credential := options.Credential{
		AuthMechanism: mongoAuthMechanism,
		AuthSource:    mongoDatabase,
		Username:      mongoUser,
		Password:      mongoPassword,
	}

	ctx := context.TODO()
	
	client, err := mongo.NewClient(opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)).SetAuth(credential))
	if err != nil {
		logger.Printf("level=error message=\"failed to open client %v\"", err)
		return err
	}
	err = client.Connect(ctx)
	if err != nil {
		logger.Printf("level=error message=\"failed to connect mongodb %v\"", err)
		return err
	}
	collection = client.Database(mongoDatabase).Collection("animals")

	return nil
}

func listAnimals(ctx context.Context, in *pb.Empty) (*pb.Animals, error) {
	logger.Printf("level=info message=\"List Animals\"")

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
		id := result.Id
		animalUuid, err := uuid.FromBytes(id.Data)
		if err != nil {
			logger.Printf("level=error message=failed to parse UUID from bytes[]: %v", err)
		}
		animal = &pb.AnimalInfo{
			Id:       animalUuid.String(),
			Type:     result.Type,
			Name:     result.Name,
			Height:   result.Height,
			Weight:   result.Weight,
			Region:   result.Region,
			IsCattle: result.IsCattle,
		}

		animals = append(animals, animal)
	}

	animalsJson, err := json.Marshal(animals)
	if err != nil {
		logger.Printf("level=error message=\"unable to marshall animals to json\"")
	}
	logger.Printf("level=info message=\"list animals from mongodb\" data=%s", string(animalsJson))

	// gRPC response
	return &pb.Animals{Animals: animals}, nil
}

func getAnimal(ctx context.Context, in *pb.AnimalId) (*pb.AnimalInfo, error) {
	logger.Printf(fmt.Sprintf("level=info message=\"Get Animal for id: %s\"", in.GetId()))

	animalUuid, err := uuid.Parse(in.GetId())
	id, err := animalUuid.MarshalBinary()
	if err != nil {
		logger.Printf("level=error message=\"unable to parse uuid: %v\"", err)
	}

	var animal AnimalInfo
	var pbAnimal pb.AnimalInfo

	// Mongo Query
	err = collection.FindOne(ctx, bson.M{"_id": primitive.Binary{
		Subtype: 0x04,
		Data:    id,
	}}).Decode(&animal)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Printf("level=error message=\"Document not found: %v\"", err)
			return &pbAnimal, status.Error(codes.NotFound, "document not found")
//			return &pbAnimal, nil
		}
		logger.Printf("level=error message\"failed to decode reuslt: %v\"", err)
		return &pbAnimal, err
	}

	resultUuid, err := uuid.FromBytes(animal.Id.Data)
	if err != nil {
		logger.Printf("level=error message=failed to parse UUID from bytes[]: %v", err)
	}

	pbAnimal = pb.AnimalInfo{
		Id:       resultUuid.String(),
		Type:     animal.Type,
		Name:     animal.Name,
		Height:   animal.Height,
		Weight:   animal.Weight,
		Region:   animal.Region,
		IsCattle: animal.IsCattle,
	}

	animalJson, err := json.Marshal(pbAnimal)
	if err != nil {
		logger.Printf("level=error message=\"unable to marshall animal to json\"")
	}
	logger.Printf("level=info message=\"get animal from mongodb\" data=%s", string(animalJson))

	return &pbAnimal, nil
}

func createAnimal(ctx context.Context, in *pb.Animal) (*pb.AnimalInfo, error) {
	logger.Printf("level=info message=\"Create Animal\"")

	var pbAnimalInfo pb.AnimalInfo
	var animalInfo AnimalInfo

	animalUuid, err := uuid.NewRandom()
	id, err := animalUuid.MarshalBinary()
	if err != nil {
		logger.Printf("level=error message=failed to create a new UUID: %v", err)
		return &pbAnimalInfo, err
	}

	logger.Printf(fmt.Sprintf("level=info message=\"insert animal data with uuid: %s\"", animalUuid.String()))

	animalInfo = AnimalInfo{
		Id: primitive.Binary{
			Subtype: 0x04,
			Data:    id,
		},
		Type:     in.Type,
		Name:     in.Name,
		Height:   in.Height,
		Weight:   in.Weight,
		Region:   in.Region,
		IsCattle: in.IsCattle,
	}

	_, err = collection.InsertOne(ctx, animalInfo)
	if err != nil {
		logger.Printf("level=error message=failed to insert new animal data: %v", err)
		return &pbAnimalInfo, err
	}

	pbAnimalInfo = pb.AnimalInfo{
		Id:       animalUuid.String(),
		Type:     in.Type,
		Name:     in.Name,
		Height:   in.Height,
		Weight:   in.Weight,
		Region:   in.Region,
		IsCattle: in.IsCattle,
	}
	animalJson, err := json.Marshal(pbAnimalInfo)
	if err != nil {
		logger.Printf("level=error message=\"unable to marshall animal to json\"")
	}
	logger.Printf("level=info message=\"create animal into mongodb\" data=%s", string(animalJson))

	return &pbAnimalInfo, nil
}

func deleteAnimal(ctx context.Context, in *pb.AnimalId) (*pb.Empty, error) {
	logger.Printf(fmt.Sprintf("level=info message=\"Delete Animal for id: %s\"", in.GetId()))

	animalUuid, err := uuid.Parse(in.GetId())
	id, err := animalUuid.MarshalBinary()
	if err != nil {
		logger.Printf("level=error message=\"unable to parse uuid: %v\"", err)
	}

	var pbEmpty pb.Empty

	// Mongo Query
	result, err := collection.DeleteOne(ctx, bson.M{"_id": primitive.Binary{
		Subtype: 0x04,
		Data:    id,
	}})

	if err != nil {
		logger.Printf("level=error message\"failed to delete record: %v\"", err)
		return &pbEmpty, err
	}

	deletedCount := result.DeletedCount
	if deletedCount < 1 {
		logger.Printf("level=info message=\"Document not found\"")
		return &pbEmpty, status.Error(codes.NotFound, "document not found")
//		return &pbEmpty, nil
	}

	logger.Printf(fmt.Sprintf("level=info message=\"deleted %s record\"", strconv.FormatInt(deletedCount, 10)))

	return &pbEmpty, nil
}

func main() {
	// Create a listener for the server.
	grpcPort := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":" + grpcPort)
	if err != nil {
		logger.Fatalf("level=fatal message=\"failed to listen: %v\"", err)
	}

	// Use signalfx tracing
	tracing.Start(tracing.WithGlobalTag("stage", "demo"), tracing.WithServiceName(serviceName))
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
	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))

	pb.RegisterAnimalServiceService(s, &pb.AnimalServiceService{
		ListAnimals: listAnimals,
		GetAnimal: getAnimal,
		CreateAnimal: createAnimal,
		DeleteAnimal: deleteAnimal,
	})
	if err = s.Serve(lis); err != nil {
		logger.Fatalf("level=fatal message=\"failed to serve at AnimalService.ListAnimals: %v\"", err)
	}
}
