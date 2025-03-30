package product

import (
	"context"
	"github.com/rasadov/EcommerceMicroservices/product/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewProductServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.ProductByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		r.Product.Id,
		r.Product.Name,
		r.Product.Description,
		FloatToString(r.Product.Price),
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, ids []string, query string) ([]Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, p := range r.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       FloatToString(p.Price),
		})
	}
	return products, nil
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.service.PostProduct(ctx, &pb.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       FloatToString(price),
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		r.Product.Id,
		r.Product.Name,
		r.Product.Description,
		FloatToString(r.Product.Price),
	}, nil
}

func (c *Client) UpdateProduct(ctx context.Context, id, name, description, price string) (*Product, error) {
	res, err := c.service.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		res.Product.Id,
		res.Product.Name,
		res.Product.Description,
		FloatToString(res.Product.Price),
	}, nil
}

func (c *Client) DeleteProduct(ctx context.Context, id string) error {
	_, err := c.service.DeleteProduct(ctx, &pb.ProductByIdRequest{Id: id})
	return err
}
