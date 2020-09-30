// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// DemoService is the service API for Demo service.
// Fields should be assigned to their respective handler implementations only before
// RegisterDemoService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type DemoService struct {
	// request demo message
	GetMessageService func(context.Context, *DemoRequest) (*DemoReply, error)
}

func (s *DemoService) getMessageService(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DemoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetMessageService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/pb.Demo/GetMessageService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetMessageService(ctx, req.(*DemoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterDemoService registers a service implementation with a gRPC server.
func RegisterDemoService(s grpc.ServiceRegistrar, srv *DemoService) {
	srvCopy := *srv
	if srvCopy.GetMessageService == nil {
		srvCopy.GetMessageService = func(context.Context, *DemoRequest) (*DemoReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetMessageService not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "pb.Demo",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetMessageService",
				Handler:    srvCopy.getMessageService,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "pb/demo.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewDemoService creates a new DemoService containing the
// implemented methods of the Demo service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewDemoService(s interface{}) *DemoService {
	ns := &DemoService{}
	if h, ok := s.(interface {
		GetMessageService(context.Context, *DemoRequest) (*DemoReply, error)
	}); ok {
		ns.GetMessageService = h.GetMessageService
	}
	return ns
}

// UnstableDemoService is the service API for Demo service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableDemoService interface {
	// request demo message
	GetMessageService(context.Context, *DemoRequest) (*DemoReply, error)
}

// AnimalServiceService is the service API for AnimalService service.
// Fields should be assigned to their respective handler implementations only before
// RegisterAnimalServiceService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type AnimalServiceService struct {
	GetAnimal    func(context.Context, *AnimalId) (*AnimalInfo, error)
	ListAnimals  func(context.Context, *Empty) (*Animals, error)
	CreateAnimal func(context.Context, *Animal) (*AnimalInfo, error)
}

func (s *AnimalServiceService) getAnimal(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnimalId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetAnimal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/pb.AnimalService/GetAnimal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetAnimal(ctx, req.(*AnimalId))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *AnimalServiceService) listAnimals(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.ListAnimals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/pb.AnimalService/ListAnimals",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.ListAnimals(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *AnimalServiceService) createAnimal(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Animal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.CreateAnimal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/pb.AnimalService/CreateAnimal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.CreateAnimal(ctx, req.(*Animal))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterAnimalServiceService registers a service implementation with a gRPC server.
func RegisterAnimalServiceService(s grpc.ServiceRegistrar, srv *AnimalServiceService) {
	srvCopy := *srv
	if srvCopy.GetAnimal == nil {
		srvCopy.GetAnimal = func(context.Context, *AnimalId) (*AnimalInfo, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetAnimal not implemented")
		}
	}
	if srvCopy.ListAnimals == nil {
		srvCopy.ListAnimals = func(context.Context, *Empty) (*Animals, error) {
			return nil, status.Errorf(codes.Unimplemented, "method ListAnimals not implemented")
		}
	}
	if srvCopy.CreateAnimal == nil {
		srvCopy.CreateAnimal = func(context.Context, *Animal) (*AnimalInfo, error) {
			return nil, status.Errorf(codes.Unimplemented, "method CreateAnimal not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "pb.AnimalService",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetAnimal",
				Handler:    srvCopy.getAnimal,
			},
			{
				MethodName: "ListAnimals",
				Handler:    srvCopy.listAnimals,
			},
			{
				MethodName: "CreateAnimal",
				Handler:    srvCopy.createAnimal,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "pb/demo.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewAnimalServiceService creates a new AnimalServiceService containing the
// implemented methods of the AnimalService service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewAnimalServiceService(s interface{}) *AnimalServiceService {
	ns := &AnimalServiceService{}
	if h, ok := s.(interface {
		GetAnimal(context.Context, *AnimalId) (*AnimalInfo, error)
	}); ok {
		ns.GetAnimal = h.GetAnimal
	}
	if h, ok := s.(interface {
		ListAnimals(context.Context, *Empty) (*Animals, error)
	}); ok {
		ns.ListAnimals = h.ListAnimals
	}
	if h, ok := s.(interface {
		CreateAnimal(context.Context, *Animal) (*AnimalInfo, error)
	}); ok {
		ns.CreateAnimal = h.CreateAnimal
	}
	return ns
}

// UnstableAnimalServiceService is the service API for AnimalService service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableAnimalServiceService interface {
	GetAnimal(context.Context, *AnimalId) (*AnimalInfo, error)
	ListAnimals(context.Context, *Empty) (*Animals, error)
	CreateAnimal(context.Context, *Animal) (*AnimalInfo, error)
}
