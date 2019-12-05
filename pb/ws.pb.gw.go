// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: ws.proto

/*
Package pb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package pb

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = descriptor.ForMessage

func request_WebsocketService_Echo_0(ctx context.Context, marshaler runtime.Marshaler, client WebsocketServiceClient, req *http.Request, pathParams map[string]string) (WebsocketService_EchoClient, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	stream, err := client.Echo(ctx, &protoReq)
	if err != nil {
		return nil, metadata, err
	}
	header, err := stream.Header()
	if err != nil {
		return nil, metadata, err
	}
	metadata.HeaderMD = header
	return stream, metadata, nil

}

// RegisterWebsocketServiceHandlerServer registers the http handlers for service WebsocketService to "mux".
// UnaryRPC     :call WebsocketServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
func RegisterWebsocketServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server WebsocketServiceServer) error {

	mux.Handle("GET", pattern_WebsocketService_Echo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		err := status.Error(codes.Unimplemented, "streaming calls are not yet supported in the in-process transport")
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
		return
	})

	return nil
}

// RegisterWebsocketServiceHandlerFromEndpoint is same as RegisterWebsocketServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterWebsocketServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterWebsocketServiceHandler(ctx, mux, conn)
}

// RegisterWebsocketServiceHandler registers the http handlers for service WebsocketService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterWebsocketServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterWebsocketServiceHandlerClient(ctx, mux, NewWebsocketServiceClient(conn))
}

// RegisterWebsocketServiceHandlerClient registers the http handlers for service WebsocketService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "WebsocketServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "WebsocketServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "WebsocketServiceClient" to call the correct interceptors.
func RegisterWebsocketServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client WebsocketServiceClient) error {

	mux.Handle("GET", pattern_WebsocketService_Echo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_WebsocketService_Echo_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_WebsocketService_Echo_0(ctx, mux, outboundMarshaler, w, req, func() (proto.Message, error) { return resp.Recv() }, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_WebsocketService_Echo_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"ws", "echo"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_WebsocketService_Echo_0 = runtime.ForwardResponseStream
)
