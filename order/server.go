package order

import (
	"context"
	"fmt"
	"github.com/deckarep/golang-set/v2"
	"github.com/rasadov/EcommerceMicroservices/account"
	"github.com/rasadov/EcommerceMicroservices/order/pb"
	"github.com/rasadov/EcommerceMicroservices/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	productClient *product.Client
}

func ListenGRPC(s Service, accountURL string, productURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	productClient, err := product.NewClient(productURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		productClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		pb.UnimplementedOrderServiceServer{},
		s,
		accountClient,
		productClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Println("Error getting account", err)
		return nil, err
	}
	var productIDs []string
	for _, p := range request.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting ordered products", err)
		return nil, err
	}

	var products []OrderedProduct

	for _, p := range orderedProducts {
		productId, _ := strconv.ParseInt(p.ID, 10, 64)
		productObj := OrderedProduct{
			ID:          uint(productId),
			Name:        p.Name,
			Description: p.Description,
			Price:       product.StringToFloat(p.Price),
			Quantity:    0,
		}
		for _, requestProduct := range request.Products {
			if requestProduct.ProductId == p.ID {
				productObj.Quantity = requestProduct.Quantity
				break
			}
		}

		if productObj.Quantity != 0 {
			products = append(products, productObj)
		}
	}

	order, err := s.service.PostOrder(ctx, request.AccountId, products)
	if err != nil {
		log.Println("Error posting order", err)
		return nil, err
	}

	orderProto := &pb.Order{
		Id:         strconv.Itoa(int(order.ID)),
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          strconv.Itoa(int(p.ID)),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, request *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, request.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Taking unique products. We use set to avoid repeating
	productIDsSet := mapset.NewSet[string]()
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDsSet.Add(strconv.Itoa(int(p.ID)))
		}
	}

	productIDs := productIDsSet.ToSlice()

	products, err := s.productClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}

	// Collecting orders

	var orders []*pb.Order
	for _, o := range accountOrders {
		// Encode order
		op := &pb.Order{
			AccountId:  o.AccountID,
			Id:         strconv.Itoa(int(o.ID)),
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Decorate orders with products
		for _, orderedProduct := range o.Products {
			// Populate product fields
			for _, p := range products {
				if p.ID == strconv.Itoa(int(orderedProduct.ID)) {
					orderedProduct.Name = p.Name
					orderedProduct.Description = p.Description
					orderedProduct.Price = product.StringToFloat(p.Price)
					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          strconv.Itoa(int(orderedProduct.ID)),
				Name:        orderedProduct.Name,
				Description: orderedProduct.Description,
				Price:       orderedProduct.Price,
				Quantity:    orderedProduct.Quantity,
			})
		}

		orders = append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
