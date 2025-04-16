package main

import (
	"context"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	pb "github.com/DinukaKaveen/Golang-gRPC-Microservices/proto/order/generated"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type orderServer struct {
	pb.UnimplementedOrderServiceServer
}

func (s *orderServer) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	orderId := uuid.New().String()
	return &pb.OrderResponse{
		OrderId: orderId,
		Status: "CREATED",
	}, nil
}

func main() {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &orderServer{})

	// Start Fiber app
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Order service is healthy")
	})

	// Run servers concurrently
	go func() {
		log.Fatal(app.Listen(":3001"))
	}()

	log.Printf("Order service running on :50051 (gRPC) and :3001 (HTTP)")
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}