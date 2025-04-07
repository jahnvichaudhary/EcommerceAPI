package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
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

// Change this if needed:
var (
	serverURL = "http://localhost:8080/graphql"
	Email     string
	Password  string
	AuthToken string
)

// doRequest is a helper that executes a GraphQL mutation/query
// against our server, attaching the JWT token as a *cookie*
// if AuthToken is set.
func doRequest(t *testing.T, serverURL, query string, variables map[string]interface{}) GraphQLResponse {
	body := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Encode request to JSON
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	// Build the request
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(b))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// If we have a token, set it as a cookie named "token"
	if AuthToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: AuthToken,
			Path:  "/",
		})
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Decode response
	var gqlResp GraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&gqlResp)
	assert.NoError(t, err)

	return gqlResp
}

// 1) REGISTER
func Test01Register(t *testing.T) {
	query := `
        mutation Register($account: RegisterInput!) {
          register(account: $account) {
            token
          }
        }
    `
	Email = fmt.Sprintf("random%d@example.com", rand.Intn(100000))
	Password = fmt.Sprintf("password%d", rand.Intn(100000))
	variables := map[string]interface{}{
		"account": map[string]interface{}{
			"name":     "John Doe",
			"email":    Email,
			"password": Password,
		},
	}

	resp := doRequest(t, serverURL, query, variables)
	assert.Nil(t, resp.Errors, "unexpected GraphQL errors during Register")

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok, "response Data should be a map")

	reg, ok := data["register"].(map[string]interface{})
	assert.True(t, ok, "register field should be a map")

	token, ok := reg["token"].(string)
	assert.True(t, ok, "token should be a string")
	assert.NotEmpty(t, token, "expected a token in register response")

	AuthToken = token // store the token globally for subsequent tests
	log.Println("Got token from Register:", AuthToken)
}

// 2) LOGIN
func Test02Login(t *testing.T) {
	query := `
        mutation Login($account: LoginInput!) {
          login(account: $account) {
            token
          }
        }
    `
	variables := map[string]interface{}{
		"account": map[string]interface{}{
			"email":    Email,
			"password": Password,
		},
	}

	resp := doRequest(t, serverURL, query, variables)
	assert.Nil(t, resp.Errors, "unexpected GraphQL errors during Login")

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok, "response Data should be a map")

	login, ok := data["login"].(map[string]interface{})
	assert.True(t, ok, "login field should be a map")

	token, ok := login["token"].(string)
	assert.True(t, ok, "token should be a string")
	assert.NotEmpty(t, token, "expected a token in login response")

	AuthToken = token // refresh the token from login (optional)
	log.Println("Got token from Login:", AuthToken)
}

// 3) CREATE PRODUCT
func Test03CreateProduct(t *testing.T) {
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
	assert.Nil(t, resp.Errors, "unexpected GraphQL errors during CreateProduct")

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok, "response Data should be a map")

	p, ok := data["createProduct"].(map[string]interface{})
	assert.True(t, ok, "createProduct field should be a map")

	assert.NotEmpty(t, p["id"], "expected product ID to be returned")
	assert.Equal(t, "Test Product", p["name"])
	assert.Equal(t, "A test description", p["description"])
	assert.EqualValues(t, 12.99, p["price"])
	log.Println("Created product:", p)
}

func Test04CreateOrder(t *testing.T) {
	// 1) Query products to get a list of available product IDs
	productQuery := `
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
		// Use "query" if you want to filter products by name or something, or just leave it blank
		"query":       nil,
		"id":          nil,
		"recommended": false,
	}

	productsResp := doRequest(t, serverURL, productQuery, variables)

	// 2) Debug logs
	log.Printf("ProductsResp errors: %#v", productsResp.Errors)
	log.Printf("ProductsResp data: %#v", productsResp.Data)

	// 3) If there are GraphQL errors, fail immediately so we see the cause
	if len(productsResp.Errors) > 0 {
		t.Fatalf("unexpected GraphQL errors during product query: %v", productsResp.Errors)
	}

	// 4) Parse the data
	productsData, ok := productsResp.Data.(map[string]interface{})
	assert.True(t, ok, "expected product query data to be a map")

	productList, ok := productsData["product"].([]interface{})
	assert.True(t, ok, "expected 'product' field to be a slice in the response")
	assert.True(t, len(productList) >= 2, "need at least 2 products to create an order")

	// 5) Pick 2 random products
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(productList), func(i, j int) {
		productList[i], productList[j] = productList[j], productList[i]
	})
	product1 := productList[0].(map[string]interface{})
	product2 := productList[1].(map[string]interface{})

	id1, _ := product1["id"].(string)
	id2, _ := product2["id"].(string)
	assert.NotEmpty(t, id1, "product 1 id is empty")
	assert.NotEmpty(t, id2, "product 2 id is empty")
	log.Printf("Randomly selected product IDs: %q and %q", id1, id2)

	// 6) Now, call CreateOrder using the 2 random IDs
	createOrderQuery := `
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
	orderVariables := map[string]interface{}{
		"order": map[string]interface{}{
			"products": []interface{}{
				map[string]interface{}{
					"id":       id1,
					"quantity": 2,
				},
				map[string]interface{}{
					"id":       id2,
					"quantity": 1,
				},
			},
		},
	}
	resp := doRequest(t, serverURL, createOrderQuery, orderVariables)

	// 7) Debug logs for order creation
	log.Printf("CreateOrderResp errors: %#v", resp.Errors)
	log.Printf("CreateOrderResp data: %#v", resp.Data)

	// 8) Check for GraphQL errors before parsing
	assert.Nil(t, resp.Errors, "unexpected GraphQL errors during CreateOrder")

	// 9) Assert the response is valid
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok, "createOrder response data should be a map")

	createdOrder, ok := data["createOrder"].(map[string]interface{})
	assert.True(t, ok, "createOrder field should be a map")

	assert.NotEmpty(t, createdOrder["id"], "expected an order ID")
	assert.NotEmpty(t, createdOrder["createdAt"], "expected a createdAt timestamp")
	assert.NotEmpty(t, createdOrder["totalPrice"], "expected a totalPrice")

	products, ok := createdOrder["products"].([]interface{})
	assert.True(t, ok, "expected products to be a list")
	assert.Len(t, products, 2, "Expected 2 products in the order")

	log.Println("Created order:", createdOrder)
}

// 5) QUERY ACCOUNTS
func Test05QueryAccounts(t *testing.T) {
	query := `
        query GetAccounts($pagination: PaginationInput) {
          accounts(pagination: $pagination) {
            id
            name
            email
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
	assert.Nil(t, resp.Errors)

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)

	accounts, ok := data["accounts"].([]interface{})
	assert.True(t, ok)
	log.Println("Accounts:", accounts)
	// Add additional assertions as needed
}

// 6) QUERY PRODUCTS
func Test06QueryProducts(t *testing.T) {
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
	assert.Nil(t, resp.Errors)

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)

	products, ok := data["product"].([]interface{})
	assert.True(t, ok)

	log.Println("Products:", products)
	// Add additional assertions as needed
}
