package main

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	helloworld "grpc_restful/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}

type middlewareFunc func(http.Handler) http.Handler

func chainMiddleware(handler http.Handler, middlewares ...middlewareFunc) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

type server struct{}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{Message: "Hello, " + in.Name}, nil
}

func main() {
	grpcAddress := ":50051"
	httpAddress := ":8080"

	// 启动gRPC服务
	go func() {
		lis, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		helloworld.RegisterHelloWorldServer(grpcServer, &server{})
		log.Printf("gRPC server listening on %s", grpcAddress)
		grpcServer.Serve(lis)
	}()

	// 启动RESTful服务
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := helloworld.RegisterHelloWorldHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
		if err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}

		log.Printf("RESTful server listening on %s", httpAddress)
		handlerWithMiddlewares := chainMiddleware(mux, GzipMiddleware)
		http.ListenAndServe(httpAddress, handlerWithMiddlewares)
	}()

	// 等待服务器结束
	select {}
}
