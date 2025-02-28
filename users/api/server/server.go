package server

import (
	"context"
	"fmt"
	"go-grpc-demo/users/proto"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	// gRPC server endpoint
	grpcAddr = fmt.Sprintf("%s:%d", "0.0.0.0", 50051)
	httpAddr = fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
)

// Server provides an implementation for the server
type Server struct {
	DB *gorm.DB
}

// NewServer returns a new server given the Options
func NewServer(DB *gorm.DB) (*Server, error) {

	// Return the server
	return &Server{
		DB: DB,
	}, nil
}

func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	//err = proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, *postEndpoint, dialOpts)
	//if err != nil {
	//	return nil, err
	//}

	return mux, nil
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter) {
	headers := []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

//allowCORS allows Cross Origin Resource Sharing from any origin.
//Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// Serve sets up the server and listens for requests
func (s *Server) Serve() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// var opts []grpc.ServerOption

	// Output some debugging
	log.Println("Starting server...")

	// Setup the gRPC server
	// unaryInterceptor checks AUTH
	g := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)
	proto.RegisterUserServiceServer(g, s)

	// start the grpc server on its own port
	go func() {
		conn, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			panic(err)
		}
		log.Println("grpc Server started on ", grpcAddr)
		g.Serve(conn)
	}()

	// Setup the mux
	mux := http.NewServeMux()

	// Setup the gateway mux
	gw, err := newGateway(ctx)
	if err != nil {
		return err
	}

	mux.Handle("/", gw)
	mux.HandleFunc("/auth", GetToken)

	// Start the http getaway mux server
	log.Println("Getaway Mux Server started on ", httpAddr)

	return http.ListenAndServe(httpAddr, allowCORS(mux))

}
