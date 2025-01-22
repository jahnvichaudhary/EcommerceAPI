package main

import "context"

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {

}

func (r *queryResolver) Product(ctx context.Context, pagination *PaginationInput, query, id *string) ([]*Product, error) {

}
