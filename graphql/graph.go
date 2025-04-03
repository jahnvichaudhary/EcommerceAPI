package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/rasadov/EcommerceMicroservices/account"
	"github.com/rasadov/EcommerceMicroservices/order"
	"github.com/rasadov/EcommerceMicroservices/product"
	"github.com/rasadov/EcommerceMicroservices/recommender"
)

type Server struct {
	accountClient     *account.Client
	productClient     *product.Client
	orderClient       *order.Client
	recommenderClient *recommender.Client
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
	recommenderClient, err := recommender.NewClient("http://localhost:8080")
	if err != nil {
		accountClient.Close()
		productClient.Close()
		orderClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		productClient,
		orderClient,
		recommenderClient,
	}, nil
}

func (server *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: server,
	}
}

func (server *Server) Query() QueryResolver {
	return &queryResolver{
		server: server,
	}
}

func (server *Server) Account() AccountResolver {
	return &accountResolver{
		server: server,
	}
}

func (server *Server) toExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: server,
	})
}
