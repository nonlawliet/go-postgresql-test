package main

func getProducts() ([]Product, error) {
	// #1 - Query all data
	rows, err := db.Query("SELECT id, name, price from products")
	if err != nil {
		return nil, err
	}

	var respProducts []Product

	// #2 - Load data
	for rows.Next() {
		var respProduct Product

		// Scan each row into respProduct
		err := rows.Scan(&respProduct.ID, &respProduct.Name, &respProduct.Price)
		if err != nil {
			return nil, err
		}

		// Append respProduct into respProducts
		respProducts = append(respProducts, respProduct)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// #3 - Make response
	return respProducts, nil
}

func createProduct(product *Product) error {
	// #1 - Insert data
	_, err := db.Exec("INSERT INTO products(name, price) VALUES ($1, $2);", product.Name, product.Price)

	// #2 - Make response
	return err
}

func getProduct(id int) (Product, error) {
	var resp Product

	// #1 - Query data
	row := db.QueryRow("SELECT id, name, price FROM products WHERE id=$1;", id)

	// #2 - Scan result from query into resp
	err := row.Scan(&resp.ID, &resp.Name, &resp.Price)

	// #3 - Make response
	if err != nil {
		return Product{}, err
	}

	return resp, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	var resp Product

	// #1 - Update data
	_, err := db.Exec("UPDATE products SET name=$2, price=$3 WHERE id=$1;", id, product.Name, product.Price)

	if err != nil {
		return Product{}, err
	}

	// #2 - Make response
	resp, err = getProduct(id)
	if err != nil {
		return Product{}, err
	}

	return resp, nil
}

func deleteProduct(id int) error {
	// #1 - Delete data
	_, err := db.Exec("DELETE FROM products WHERE id=$1;", id)

	// #2 - Make response
	return err
}
