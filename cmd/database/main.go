package main

import (
	"context"
	"fmt"
	"time"
	"training/database"
	"training/persistence"
	"training/pointer_func"

	"github.com/google/uuid"
)

func main() {

	ctx := context.Background()

	postgresUrl := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	db := database.NewPostgresDB(postgresUrl)
	defer db.Close()

	productId := uuid.Must(uuid.NewV7())

	// **************** Insert ******************
	product := persistence.Product{
		ProductId:   productId,
		ProductName: "Coca Cola",
		Price:       15.50,
		CreatedBy:   "ADMIN",
		CreatedAt:   time.Now(),
	}
	cmd, err := db.Exec(ctx, "insert into products(product_id, product_name, price, created_by, created_at) values($1, $2, $3, $4, $5)",
		product.ProductId,
		product.ProductName,
		product.Price,
		product.CreatedBy,
		product.CreatedAt,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted affected: ", cmd.RowsAffected())

	// **************** Update ******************
	product.ProductName = "Pepsi"
	product.Price = 16.20
	product.UpdatedBy = pointer_func.ToPointer("ADMIN")
	product.UpdatedAt = pointer_func.ToPointer(time.Now())

	cmd, err = db.Exec(ctx, "update products set product_name = $1, price = $2, updated_by = $3, updated_at = $4 where product_id = $5",
		product.ProductName,
		product.Price,
		product.UpdatedBy,
		product.UpdatedAt,
		product.ProductId,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Updated affected: ", cmd.RowsAffected())

	// **************** Select ******************

	var getProduct persistence.Product
	err = db.QueryRow(ctx, "SELECT product_id, product_name, price, created_by, created_at, updated_by, updated_at FROM products WHERE product_id = $1",
		productId,
	).Scan(&getProduct.ProductId,
		&getProduct.ProductName,
		&getProduct.Price,
		&getProduct.CreatedBy,
		&getProduct.CreatedAt,
		&getProduct.UpdatedBy,
		&getProduct.UpdatedAt,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Select %v \n", getProduct)

	// **************** Delete ******************
	cmd, err = db.Exec(ctx,
		"delete from products where product_id = $1",
		productId,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted affected: ", cmd.RowsAffected())
}
