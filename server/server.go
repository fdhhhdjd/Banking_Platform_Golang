package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	error_response "github.com/fdhhhdjd/Banking_Platform_Golang/api/handler/error"
	"github.com/fdhhhdjd/Banking_Platform_Golang/api/middlewares"
	config "github.com/fdhhhdjd/Banking_Platform_Golang/configs"
	"github.com/fdhhhdjd/Banking_Platform_Golang/gapi"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/db"
	routes "github.com/fdhhhdjd/Banking_Platform_Golang/internals/routers"
	util "github.com/fdhhhdjd/Banking_Platform_Golang/internals/utils"
	"github.com/fdhhhdjd/Banking_Platform_Golang/pb"
	"github.com/fdhhhdjd/Banking_Platform_Golang/pkg"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	pkg.InitLog()
	config.LoadConfig("configs")

	// DB
	err := db.InitStore()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

}

func Server() {
	nodeEnv := os.Getenv("ENV")

	var server *gin.Engine

	if nodeEnv == constants.DEV {
		server = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		server = gin.New()
		server.Use(gin.Logger())
		server.Use(gin.Recovery())
	}

	server.Use(middlewares.CORSMiddleware())
	server.Use(middlewares.SecurityHeadersMiddleware())
	server.Use(pkg.LogRequest)

	port := os.Getenv("PORT")

	if port == "" {
		port = strconv.Itoa(config.AppConfig.Server.Port)
	}

	server.GET("/ping", Pong)

	routes.SetupRoutes(server)

	// 404 - Not Found
	server.NoRoute(NotFound())

	// 500 - Internal Server Error
	server.Use(ServerError)

	server.Run(":" + port)
}

func StartGRPCServer() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "5005"
	}

	//* gRPC
	store := db.GetStore()
	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterSimpleBankServer(grpcServer, gapi.NewSimpleBankServer(*store))

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", port)
	// log.Printf("Registered gRPC services: %v", grpcServer.GetServiceInfo())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func StartGateWayGRPCServer() {
	port := os.Getenv("GATE_WAY_GRPC_PORT")
	if port == "" {
		port = "5006"
	}

	// * Option get key for database not key to define
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	store := db.GetStore()
	server := gapi.NewSimpleBankServer(*store)

	gwmux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	pb.RegisterSimpleBankHandlerServer(ctx, gwmux, server)

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	//* Swagger API
	fs := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen Gateway: %v", err)
	}

	log.Printf("Starting gRPC Gateway server on %s", listener.Addr().String())
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func Pong(c *gin.Context) {
	currentTime := time.Now().Unix()
	signal := "Nguyen Tien Tai"
	message := util.P("Pong %s", signal)

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": message,
		"time":    currentTime,
	})
}

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		errorResponse := error_response.NotFoundError("Not Found")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
	}
}

func ServerError(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		err := c.Errors[0].Err
		status := http.StatusInternalServerError
		requestId, _ := c.Get("requestId")

		logrus.WithFields(logrus.Fields{
			"status":    status,
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"time":      time.Now(),
			"request":   c.Request.RequestURI,
			"requestId": requestId,
		}).Error(err)

		errorResponse := error_response.ServiceUnavailable("Service Unavailable")
		c.AbortWithStatusJSON(errorResponse.Status, errorResponse)
	}
}
