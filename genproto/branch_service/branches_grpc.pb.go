// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: branches.proto

package branch_service

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BranchServiceClient is the client API for BranchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BranchServiceClient interface {
	Create(ctx context.Context, in *CreateBranch, opts ...grpc.CallOption) (*GetBranch, error)
	GetByID(ctx context.Context, in *BranchPrimaryKey, opts ...grpc.CallOption) (*GetBranch, error)
	GetList(ctx context.Context, in *GetListBranchRequest, opts ...grpc.CallOption) (*GetListBranchResponse, error)
	Update(ctx context.Context, in *UpdateBranch, opts ...grpc.CallOption) (*GetBranch, error)
	Delete(ctx context.Context, in *BranchPrimaryKey, opts ...grpc.CallOption) (*empty.Empty, error)
}

type branchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBranchServiceClient(cc grpc.ClientConnInterface) BranchServiceClient {
	return &branchServiceClient{cc}
}

func (c *branchServiceClient) Create(ctx context.Context, in *CreateBranch, opts ...grpc.CallOption) (*GetBranch, error) {
	out := new(GetBranch)
	err := c.cc.Invoke(ctx, "/branch_service_go.BranchService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *branchServiceClient) GetByID(ctx context.Context, in *BranchPrimaryKey, opts ...grpc.CallOption) (*GetBranch, error) {
	out := new(GetBranch)
	err := c.cc.Invoke(ctx, "/branch_service_go.BranchService/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *branchServiceClient) GetList(ctx context.Context, in *GetListBranchRequest, opts ...grpc.CallOption) (*GetListBranchResponse, error) {
	out := new(GetListBranchResponse)
	err := c.cc.Invoke(ctx, "/branch_service_go.BranchService/GetList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *branchServiceClient) Update(ctx context.Context, in *UpdateBranch, opts ...grpc.CallOption) (*GetBranch, error) {
	out := new(GetBranch)
	err := c.cc.Invoke(ctx, "/branch_service_go.BranchService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *branchServiceClient) Delete(ctx context.Context, in *BranchPrimaryKey, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/branch_service_go.BranchService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BranchServiceServer is the server API for BranchService service.
// All implementations should embed UnimplementedBranchServiceServer
// for forward compatibility
type BranchServiceServer interface {
	Create(context.Context, *CreateBranch) (*GetBranch, error)
	GetByID(context.Context, *BranchPrimaryKey) (*GetBranch, error)
	GetList(context.Context, *GetListBranchRequest) (*GetListBranchResponse, error)
	Update(context.Context, *UpdateBranch) (*GetBranch, error)
	Delete(context.Context, *BranchPrimaryKey) (*empty.Empty, error)
}

// UnimplementedBranchServiceServer should be embedded to have forward compatible implementations.
type UnimplementedBranchServiceServer struct {
}

func (UnimplementedBranchServiceServer) Create(context.Context, *CreateBranch) (*GetBranch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedBranchServiceServer) GetByID(context.Context, *BranchPrimaryKey) (*GetBranch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedBranchServiceServer) GetList(context.Context, *GetListBranchRequest) (*GetListBranchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedBranchServiceServer) Update(context.Context, *UpdateBranch) (*GetBranch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedBranchServiceServer) Delete(context.Context, *BranchPrimaryKey) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// UnsafeBranchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BranchServiceServer will
// result in compilation errors.
type UnsafeBranchServiceServer interface {
	mustEmbedUnimplementedBranchServiceServer()
}

func RegisterBranchServiceServer(s grpc.ServiceRegistrar, srv BranchServiceServer) {
	s.RegisterService(&BranchService_ServiceDesc, srv)
}

func _BranchService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBranch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BranchServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/branch_service_go.BranchService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BranchServiceServer).Create(ctx, req.(*CreateBranch))
	}
	return interceptor(ctx, in, info, handler)
}

func _BranchService_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BranchPrimaryKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BranchServiceServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/branch_service_go.BranchService/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BranchServiceServer).GetByID(ctx, req.(*BranchPrimaryKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _BranchService_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BranchServiceServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/branch_service_go.BranchService/GetList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BranchServiceServer).GetList(ctx, req.(*GetListBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BranchService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBranch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BranchServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/branch_service_go.BranchService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BranchServiceServer).Update(ctx, req.(*UpdateBranch))
	}
	return interceptor(ctx, in, info, handler)
}

func _BranchService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BranchPrimaryKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BranchServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/branch_service_go.BranchService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BranchServiceServer).Delete(ctx, req.(*BranchPrimaryKey))
	}
	return interceptor(ctx, in, info, handler)
}

// BranchService_ServiceDesc is the grpc.ServiceDesc for BranchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BranchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "branch_service_go.BranchService",
	HandlerType: (*BranchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _BranchService_Create_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _BranchService_GetByID_Handler,
		},
		{
			MethodName: "GetList",
			Handler:    _BranchService_GetList_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _BranchService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _BranchService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "branches.proto",
}
