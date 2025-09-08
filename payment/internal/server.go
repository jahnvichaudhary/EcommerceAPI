package internal

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	order "github.com/rasadov/EcommerceAPI/order/client"
	"github.com/rasadov/EcommerceAPI/payment/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartServers runs both gRPC and HTTP webhook servers concurrently
func StartServers(service Service, orderURL string, grpcPort, webhookPort int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ListenGRPC(service, orderURL, grpcPort); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Start webhook HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := listenWebhook(service, orderURL, webhookPort); err != nil {
			errCh <- fmt.Errorf("webhook server error: %w", err)
		}
	}()

	// Wait for first error or all servers to complete
	go func() {
		wg.Wait()
		close(errCh)
	}()

	return <-errCh
}

func ListenGRPC(service Service, orderURL string, port int) error {
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		orderClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterPaymentServiceServer(serv, &grpcServer{
		pb.UnimplementedPaymentServiceServer{},
		service,
		orderClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func listenWebhook(service Service, orderURL string, port int) error {
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		return err
	}
	defer orderClient.Close()

	webhookServer := &WebhookServer{
		service:     service,
		orderClient: orderClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/payment", webhookServer.HandlePaymentWebhook)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Webhook server listening on port %d", port)
	return server.ListenAndServe()
}
