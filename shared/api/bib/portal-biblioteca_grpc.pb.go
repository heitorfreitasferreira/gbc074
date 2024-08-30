// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: api/portal-biblioteca.proto

package api_bib

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PortalBibliotecaClient is the client API for PortalBiblioteca service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PortalBibliotecaClient interface {
	RealizaEmprestimo(ctx context.Context, opts ...grpc.CallOption) (PortalBiblioteca_RealizaEmprestimoClient, error)
	RealizaDevolucao(ctx context.Context, opts ...grpc.CallOption) (PortalBiblioteca_RealizaDevolucaoClient, error)
	BloqueiaUsuarios(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (*Status, error)
	LiberaUsuarios(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (*Status, error)
	ListaUsuariosBloqueados(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (PortalBiblioteca_ListaUsuariosBloqueadosClient, error)
	ListaLivrosEmprestados(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (PortalBiblioteca_ListaLivrosEmprestadosClient, error)
	ListaLivrosEmFalta(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (PortalBiblioteca_ListaLivrosEmFaltaClient, error)
	PesquisaLivro(ctx context.Context, in *Criterio, opts ...grpc.CallOption) (PortalBiblioteca_PesquisaLivroClient, error)
}

type portalBibliotecaClient struct {
	cc grpc.ClientConnInterface
}

func NewPortalBibliotecaClient(cc grpc.ClientConnInterface) PortalBibliotecaClient {
	return &portalBibliotecaClient{cc}
}

func (c *portalBibliotecaClient) RealizaEmprestimo(ctx context.Context, opts ...grpc.CallOption) (PortalBiblioteca_RealizaEmprestimoClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortalBiblioteca_ServiceDesc.Streams[0], "/biblioteca.PortalBiblioteca/RealizaEmprestimo", opts...)
	if err != nil {
		return nil, err
	}
	x := &portalBibliotecaRealizaEmprestimoClient{stream}
	return x, nil
}

type PortalBiblioteca_RealizaEmprestimoClient interface {
	Send(*UsuarioLivro) error
	CloseAndRecv() (*Status, error)
	grpc.ClientStream
}

type portalBibliotecaRealizaEmprestimoClient struct {
	grpc.ClientStream
}

func (x *portalBibliotecaRealizaEmprestimoClient) Send(m *UsuarioLivro) error {
	return x.ClientStream.SendMsg(m)
}

func (x *portalBibliotecaRealizaEmprestimoClient) CloseAndRecv() (*Status, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *portalBibliotecaClient) RealizaDevolucao(ctx context.Context, opts ...grpc.CallOption) (PortalBiblioteca_RealizaDevolucaoClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortalBiblioteca_ServiceDesc.Streams[1], "/biblioteca.PortalBiblioteca/RealizaDevolucao", opts...)
	if err != nil {
		return nil, err
	}
	x := &portalBibliotecaRealizaDevolucaoClient{stream}
	return x, nil
}

type PortalBiblioteca_RealizaDevolucaoClient interface {
	Send(*UsuarioLivro) error
	CloseAndRecv() (*Status, error)
	grpc.ClientStream
}

type portalBibliotecaRealizaDevolucaoClient struct {
	grpc.ClientStream
}

func (x *portalBibliotecaRealizaDevolucaoClient) Send(m *UsuarioLivro) error {
	return x.ClientStream.SendMsg(m)
}

func (x *portalBibliotecaRealizaDevolucaoClient) CloseAndRecv() (*Status, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *portalBibliotecaClient) BloqueiaUsuarios(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/biblioteca.PortalBiblioteca/BloqueiaUsuarios", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portalBibliotecaClient) LiberaUsuarios(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/biblioteca.PortalBiblioteca/LiberaUsuarios", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portalBibliotecaClient) ListaUsuariosBloqueados(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (PortalBiblioteca_ListaUsuariosBloqueadosClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortalBiblioteca_ServiceDesc.Streams[2], "/biblioteca.PortalBiblioteca/ListaUsuariosBloqueados", opts...)
	if err != nil {
		return nil, err
	}
	x := &portalBibliotecaListaUsuariosBloqueadosClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PortalBiblioteca_ListaUsuariosBloqueadosClient interface {
	Recv() (*UsuarioBloqueado, error)
	grpc.ClientStream
}

type portalBibliotecaListaUsuariosBloqueadosClient struct {
	grpc.ClientStream
}

func (x *portalBibliotecaListaUsuariosBloqueadosClient) Recv() (*UsuarioBloqueado, error) {
	m := new(UsuarioBloqueado)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *portalBibliotecaClient) ListaLivrosEmprestados(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (PortalBiblioteca_ListaLivrosEmprestadosClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortalBiblioteca_ServiceDesc.Streams[3], "/biblioteca.PortalBiblioteca/ListaLivrosEmprestados", opts...)
	if err != nil {
		return nil, err
	}
	x := &portalBibliotecaListaLivrosEmprestadosClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PortalBiblioteca_ListaLivrosEmprestadosClient interface {
	Recv() (*Livro, error)
	grpc.ClientStream
}

type portalBibliotecaListaLivrosEmprestadosClient struct {
	grpc.ClientStream
}

func (x *portalBibliotecaListaLivrosEmprestadosClient) Recv() (*Livro, error) {
	m := new(Livro)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *portalBibliotecaClient) ListaLivrosEmFalta(ctx context.Context, in *Vazia, opts ...grpc.CallOption) (PortalBiblioteca_ListaLivrosEmFaltaClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortalBiblioteca_ServiceDesc.Streams[4], "/biblioteca.PortalBiblioteca/ListaLivrosEmFalta", opts...)
	if err != nil {
		return nil, err
	}
	x := &portalBibliotecaListaLivrosEmFaltaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PortalBiblioteca_ListaLivrosEmFaltaClient interface {
	Recv() (*Livro, error)
	grpc.ClientStream
}

type portalBibliotecaListaLivrosEmFaltaClient struct {
	grpc.ClientStream
}

func (x *portalBibliotecaListaLivrosEmFaltaClient) Recv() (*Livro, error) {
	m := new(Livro)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *portalBibliotecaClient) PesquisaLivro(ctx context.Context, in *Criterio, opts ...grpc.CallOption) (PortalBiblioteca_PesquisaLivroClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortalBiblioteca_ServiceDesc.Streams[5], "/biblioteca.PortalBiblioteca/PesquisaLivro", opts...)
	if err != nil {
		return nil, err
	}
	x := &portalBibliotecaPesquisaLivroClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PortalBiblioteca_PesquisaLivroClient interface {
	Recv() (*Livro, error)
	grpc.ClientStream
}

type portalBibliotecaPesquisaLivroClient struct {
	grpc.ClientStream
}

func (x *portalBibliotecaPesquisaLivroClient) Recv() (*Livro, error) {
	m := new(Livro)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PortalBibliotecaServer is the server API for PortalBiblioteca service.
// All implementations must embed UnimplementedPortalBibliotecaServer
// for forward compatibility
type PortalBibliotecaServer interface {
	RealizaEmprestimo(PortalBiblioteca_RealizaEmprestimoServer) error
	RealizaDevolucao(PortalBiblioteca_RealizaDevolucaoServer) error
	BloqueiaUsuarios(context.Context, *Vazia) (*Status, error)
	LiberaUsuarios(context.Context, *Vazia) (*Status, error)
	ListaUsuariosBloqueados(*Vazia, PortalBiblioteca_ListaUsuariosBloqueadosServer) error
	ListaLivrosEmprestados(*Vazia, PortalBiblioteca_ListaLivrosEmprestadosServer) error
	ListaLivrosEmFalta(*Vazia, PortalBiblioteca_ListaLivrosEmFaltaServer) error
	PesquisaLivro(*Criterio, PortalBiblioteca_PesquisaLivroServer) error
	mustEmbedUnimplementedPortalBibliotecaServer()
}

// UnimplementedPortalBibliotecaServer must be embedded to have forward compatible implementations.
type UnimplementedPortalBibliotecaServer struct {
}

func (UnimplementedPortalBibliotecaServer) RealizaEmprestimo(PortalBiblioteca_RealizaEmprestimoServer) error {
	return status.Errorf(codes.Unimplemented, "method RealizaEmprestimo not implemented")
}
func (UnimplementedPortalBibliotecaServer) RealizaDevolucao(PortalBiblioteca_RealizaDevolucaoServer) error {
	return status.Errorf(codes.Unimplemented, "method RealizaDevolucao not implemented")
}
func (UnimplementedPortalBibliotecaServer) BloqueiaUsuarios(context.Context, *Vazia) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BloqueiaUsuarios not implemented")
}
func (UnimplementedPortalBibliotecaServer) LiberaUsuarios(context.Context, *Vazia) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiberaUsuarios not implemented")
}
func (UnimplementedPortalBibliotecaServer) ListaUsuariosBloqueados(*Vazia, PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	return status.Errorf(codes.Unimplemented, "method ListaUsuariosBloqueados not implemented")
}
func (UnimplementedPortalBibliotecaServer) ListaLivrosEmprestados(*Vazia, PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	return status.Errorf(codes.Unimplemented, "method ListaLivrosEmprestados not implemented")
}
func (UnimplementedPortalBibliotecaServer) ListaLivrosEmFalta(*Vazia, PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	return status.Errorf(codes.Unimplemented, "method ListaLivrosEmFalta not implemented")
}
func (UnimplementedPortalBibliotecaServer) PesquisaLivro(*Criterio, PortalBiblioteca_PesquisaLivroServer) error {
	return status.Errorf(codes.Unimplemented, "method PesquisaLivro not implemented")
}
func (UnimplementedPortalBibliotecaServer) mustEmbedUnimplementedPortalBibliotecaServer() {}

// UnsafePortalBibliotecaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PortalBibliotecaServer will
// result in compilation errors.
type UnsafePortalBibliotecaServer interface {
	mustEmbedUnimplementedPortalBibliotecaServer()
}

func RegisterPortalBibliotecaServer(s grpc.ServiceRegistrar, srv PortalBibliotecaServer) {
	s.RegisterService(&PortalBiblioteca_ServiceDesc, srv)
}

func _PortalBiblioteca_RealizaEmprestimo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PortalBibliotecaServer).RealizaEmprestimo(&portalBibliotecaRealizaEmprestimoServer{stream})
}

type PortalBiblioteca_RealizaEmprestimoServer interface {
	SendAndClose(*Status) error
	Recv() (*UsuarioLivro, error)
	grpc.ServerStream
}

type portalBibliotecaRealizaEmprestimoServer struct {
	grpc.ServerStream
}

func (x *portalBibliotecaRealizaEmprestimoServer) SendAndClose(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *portalBibliotecaRealizaEmprestimoServer) Recv() (*UsuarioLivro, error) {
	m := new(UsuarioLivro)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PortalBiblioteca_RealizaDevolucao_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PortalBibliotecaServer).RealizaDevolucao(&portalBibliotecaRealizaDevolucaoServer{stream})
}

type PortalBiblioteca_RealizaDevolucaoServer interface {
	SendAndClose(*Status) error
	Recv() (*UsuarioLivro, error)
	grpc.ServerStream
}

type portalBibliotecaRealizaDevolucaoServer struct {
	grpc.ServerStream
}

func (x *portalBibliotecaRealizaDevolucaoServer) SendAndClose(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *portalBibliotecaRealizaDevolucaoServer) Recv() (*UsuarioLivro, error) {
	m := new(UsuarioLivro)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PortalBiblioteca_BloqueiaUsuarios_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Vazia)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortalBibliotecaServer).BloqueiaUsuarios(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/biblioteca.PortalBiblioteca/BloqueiaUsuarios",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortalBibliotecaServer).BloqueiaUsuarios(ctx, req.(*Vazia))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortalBiblioteca_LiberaUsuarios_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Vazia)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortalBibliotecaServer).LiberaUsuarios(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/biblioteca.PortalBiblioteca/LiberaUsuarios",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortalBibliotecaServer).LiberaUsuarios(ctx, req.(*Vazia))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortalBiblioteca_ListaUsuariosBloqueados_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Vazia)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PortalBibliotecaServer).ListaUsuariosBloqueados(m, &portalBibliotecaListaUsuariosBloqueadosServer{stream})
}

type PortalBiblioteca_ListaUsuariosBloqueadosServer interface {
	Send(*UsuarioBloqueado) error
	grpc.ServerStream
}

type portalBibliotecaListaUsuariosBloqueadosServer struct {
	grpc.ServerStream
}

func (x *portalBibliotecaListaUsuariosBloqueadosServer) Send(m *UsuarioBloqueado) error {
	return x.ServerStream.SendMsg(m)
}

func _PortalBiblioteca_ListaLivrosEmprestados_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Vazia)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PortalBibliotecaServer).ListaLivrosEmprestados(m, &portalBibliotecaListaLivrosEmprestadosServer{stream})
}

type PortalBiblioteca_ListaLivrosEmprestadosServer interface {
	Send(*Livro) error
	grpc.ServerStream
}

type portalBibliotecaListaLivrosEmprestadosServer struct {
	grpc.ServerStream
}

func (x *portalBibliotecaListaLivrosEmprestadosServer) Send(m *Livro) error {
	return x.ServerStream.SendMsg(m)
}

func _PortalBiblioteca_ListaLivrosEmFalta_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Vazia)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PortalBibliotecaServer).ListaLivrosEmFalta(m, &portalBibliotecaListaLivrosEmFaltaServer{stream})
}

type PortalBiblioteca_ListaLivrosEmFaltaServer interface {
	Send(*Livro) error
	grpc.ServerStream
}

type portalBibliotecaListaLivrosEmFaltaServer struct {
	grpc.ServerStream
}

func (x *portalBibliotecaListaLivrosEmFaltaServer) Send(m *Livro) error {
	return x.ServerStream.SendMsg(m)
}

func _PortalBiblioteca_PesquisaLivro_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Criterio)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PortalBibliotecaServer).PesquisaLivro(m, &portalBibliotecaPesquisaLivroServer{stream})
}

type PortalBiblioteca_PesquisaLivroServer interface {
	Send(*Livro) error
	grpc.ServerStream
}

type portalBibliotecaPesquisaLivroServer struct {
	grpc.ServerStream
}

func (x *portalBibliotecaPesquisaLivroServer) Send(m *Livro) error {
	return x.ServerStream.SendMsg(m)
}

// PortalBiblioteca_ServiceDesc is the grpc.ServiceDesc for PortalBiblioteca service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PortalBiblioteca_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "biblioteca.PortalBiblioteca",
	HandlerType: (*PortalBibliotecaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BloqueiaUsuarios",
			Handler:    _PortalBiblioteca_BloqueiaUsuarios_Handler,
		},
		{
			MethodName: "LiberaUsuarios",
			Handler:    _PortalBiblioteca_LiberaUsuarios_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RealizaEmprestimo",
			Handler:       _PortalBiblioteca_RealizaEmprestimo_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RealizaDevolucao",
			Handler:       _PortalBiblioteca_RealizaDevolucao_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ListaUsuariosBloqueados",
			Handler:       _PortalBiblioteca_ListaUsuariosBloqueados_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListaLivrosEmprestados",
			Handler:       _PortalBiblioteca_ListaLivrosEmprestados_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListaLivrosEmFalta",
			Handler:       _PortalBiblioteca_ListaLivrosEmFalta_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PesquisaLivro",
			Handler:       _PortalBiblioteca_PesquisaLivro_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/portal-biblioteca.proto",
}