package main

import (
	"log"

	pb "github.com/your-repo-name/your-project-name/proto/order/generated"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Create new gRPC client to connect to Order Service gRPC
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to order service: %v", err)
	}

	// Dont forget to close it
	defer conn.Close()

	// Create a new order service client from generated code and pass in the connection created above
	orderClient := pb.NewOrderServiceClient(conn)

	// Initialize Fiber app
	app := fiber.New()

	// Sample REST endpoint
	app.Post("/users/:id/order", func(c *fiber.Ctx) error {
		userId := c.Params("id")

		// Call Order Service via gRPC
		resp, err := orderClient.CreateOrder(c.Context(), &pb.OrderRequest{
			UserId: userId,
			Amount: 100.00,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"order_id": resp.OrderId,
			"status":   resp.Status,
		})
	})

	log.Fatal(app.Listen(":3000"))
}