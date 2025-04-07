package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

// GraphQLRequest is a helper struct for the request body
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse is a helper struct for the response body
type GraphQLResponse struct {
	Data   interface{}   `json:"data,omitempty"`
	Errors []interface{} `json:"errors,omitempty"`
}

var (
	serverURL = "http://localhost:8080/graphql"
	AuthToken string
)

// doRequest is a helper that executes a GraphQL query against our test server
func doRequest(t *testing.T, serverURL, query string, variables map[string]interface{}) GraphQLResponse {
	body := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Encode request to JSON
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	var req *http.Request
	var resp *http.Response

	// Execute request
	if AuthToken != "" {
		req, err = http.NewRequest("POST", serverURL, bytes.NewBuffer(b))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+AuthToken)
		resp, err = http.DefaultClient.Do(req)
	} else {
		resp, err = http.Post(serverURL, "application/json", bytes.NewBuffer(b))
	}
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Decode response
	var gqlResp GraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&gqlResp)
	assert.NoError(t, err)

	return gqlResp
}

func TestRegister(t *testing.T) {
	query := `
        mutation Register($account: RegisterInput!) {
          register(account: $account) {
            token
          }
        }
    `
	// Variables
	variables := map[string]interface{}{
		"account": map[string]interface{}{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "mypass123",
		},
	}

	// Make the request
	resp := doRequest(t, serverURL, query, variables)

	// We expect "token" in response data
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Nil(t, resp.Errors)

	reg, ok := data["register"].(map[string]interface{})
	AuthToken = reg["token"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, reg["token"], "expected a token in register response")
}

func TestLogin(t *testing.T) {
	query := `
        mutation Login($account: LoginInput!) {
          login(account: $account) {
            token
          }
        }
    `
	variables := map[string]interface{}{
		"account": map[string]interface{}{
			"email":    "john@example.com",
			"password": "mypass123",
		},
	}

	resp := doRequest(t, serverURL, query, variables)
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Nil(t, resp.Errors)

	login, ok := data["login"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, login["token"], "expected a token in login response")
}

func TestCreateProduct(t *testing.T) {
	query := `
        mutation CreateProduct($product: CreateProductInput!) {
          createProduct(product: $product) {
            id
            name
            description
            price
            accountId
          }
        }
    `
	variables := map[string]interface{}{
		"product": map[string]interface{}{
			"name":        "Test Product",
			"description": "A test description",
			"price":       12.99,
		},
	}

	resp := doRequest(t, serverURL, query, variables)
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Nil(t, resp.Errors)

	p, ok := data["createProduct"].(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, p["id"])
	assert.Equal(t, "Test Product", p["name"])
	assert.Equal(t, "A test description", p["description"])
	assert.EqualValues(t, 12.99, p["price"])
}

func TestCreateOrder(t *testing.T) {
	query := `
        mutation CreateOrder($order: OrderInput!) {
          createOrder(order: $order) {
            id
            createdAt
            totalPrice
            products {
                id
                name
                price
                quantity
            }
          }
        }
    `
	variables := map[string]interface{}{
		"order": map[string]interface{}{
			"products": []interface{}{
				map[string]interface{}{
					"id":       "product1",
					"quantity": 2,
				},
				map[string]interface{}{
					"id":       "product2",
					"quantity": 1,
				},
			},
		},
	}

	resp := doRequest(t, serverURL, query, variables)
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Nil(t, resp.Errors)

	createdOrder, ok := data["createOrder"].(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, createdOrder["id"])
	assert.NotEmpty(t, createdOrder["createdAt"])
	assert.NotEmpty(t, createdOrder["totalPrice"])

	products, ok := createdOrder["products"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, products, 2, "Expected 2 products in the order")
}

func TestQueryAccounts(t *testing.T) {
	query := `
        query GetAccounts($pagination: PaginationInput) {
          accounts(pagination: $pagination) {
            id
            name
            email
            orders {
              id
              totalPrice
            }
          }
        }
    `
	variables := map[string]interface{}{
		"pagination": map[string]interface{}{
			"skip": 0,
			"take": 10,
		},
	}

	resp := doRequest(t, serverURL, query, variables)
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Nil(t, resp.Errors)

	accounts, ok := data["accounts"].([]interface{})
	assert.True(t, ok)
	log.Println("Accounts:", accounts)
}

func TestQueryProducts(t *testing.T) {
	query := `
        query GetProducts($pagination: PaginationInput, $query: String, $id: String, $recommended: Boolean) {
          product(pagination: $pagination, query: $query, id: $id, recommended: $recommended) {
            id
            name
            description
            price
            accountId
          }
        }
    `
	variables := map[string]interface{}{
		"pagination": map[string]interface{}{
			"skip": 0,
			"take": 5,
		},
		// "query":       "",
		// "id":         "",
		"recommended": false,
	}

	resp := doRequest(t, serverURL, query, variables)
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Nil(t, resp.Errors)

	products, ok := data["product"].([]interface{})
	assert.True(t, ok)
	log.Println("Products:", products)
}
