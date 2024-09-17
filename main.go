package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5432
	databaseName = "mydatabase"
	username     = "myuser"
	password     = "mypassword"
)

// Create global variable
// - db
var db *sql.DB

// - product struct
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName)

	// Connect with Database
	sdb, err := sql.Open("postgres", psqlInfo)

	// defer - close database (before the program close, defer will execute)
	defer sdb.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Verify a connection to db
	err = sdb.Ping()

	// If err still nil, put into the global variable and use its instead
	db = sdb

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection Database Successful")

	// Create app fiber
	app := fiber.New()

	// Create handler function
	app.Get("/product/:id", getProductHandler)
	app.Post("/product", createProductHandler)
	app.Put("/product/:id", updateProductHandler)
	app.Delete("/product/:id", deleteProductHandler)
	app.Get("/products", getProductsHandler)

	app.Listen(":8080")
}

func getProductHandler(c *fiber.Ctx) error {
	// #1 - Get id from request
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #2 - Get product
	product, err := getProduct(productId)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #3 - Make response
	return c.JSON(product)
}

func createProductHandler(c *fiber.Ctx) error {
	// #1 - Create product instance
	product := new(Product)

	// #2 - Get info from context into product instance
	if err := c.BodyParser(product); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #3 - Create product
	err := createProduct(product)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(product)
}

func updateProductHandler(c *fiber.Ctx) error {
	var resp Product

	// #1 - Get id from request
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #2 - Create product instance
	product := new(Product)

	// #3 - Get info from context into product instance
	if err := c.BodyParser(product); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #4 - Update product info
	resp, err = updateProduct(productId, product)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #5 - Make response
	return c.JSON(resp)
}

func deleteProductHandler(c *fiber.Ctx) error {
	// #1 - Get id from request
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #2 - Delete product
	err = deleteProduct(productId)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #3 - Make response
	return c.SendStatus(fiber.StatusNoContent)
}

func getProductsHandler(c *fiber.Ctx) error {
	// #1 - Get product all
	products, err := getProducts()
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// #2 - Make response
	return c.JSON(products)
}
