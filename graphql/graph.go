package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/rasadov/EcommerceMicroservices/account"
	"github.com/rasadov/EcommerceMicroservices/order"
	"github.com/rasadov/EcommerceMicroservices/product"
)

type Server struct {
	accountClient *account.Client
	productClient *product.Client
	orderClient   *order.Client
}

func NewGraphQLServer(
	accountUrl, productUrl, orderUrl string) (
	*Server, error) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}
	productClient, err := product.NewClient(productUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}
	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		productClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		productClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) toExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
