package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/rasadov/EcommerceAPI/product/client"

	account "github.com/rasadov/EcommerceAPI/account/client"
	order "github.com/rasadov/EcommerceAPI/order/client"
	"github.com/rasadov/EcommerceAPI/recommender"
)

type Server struct {
	accountClient     *account.Client
	productClient     *client.Client
	orderClient       *order.Client
	recommenderClient *recommender.Client
}

func NewGraphQLServer(accountUrl, productUrl, orderUrl string) (*Server, error) {
	accClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	prodClient, err := client.NewClient(productUrl)
	if err != nil {
		accClient.Close()
		return nil, err
	}

	ordClient, err := order.NewClient(orderUrl)
	if err != nil {
		accClient.Close()
		prodClient.Close()
		return nil, err
	}

	recClient, err := recommender.NewClient("http://localhost:8080")
	if err != nil {
		accClient.Close()
		prodClient.Close()
		ordClient.Close()
		return nil, err
	}

	return &Server{
		accountClient:     accClient,
		productClient:     prodClient,
		orderClient:       ordClient,
		recommenderClient: recClient,
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

func (server *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: server,
	})
}
