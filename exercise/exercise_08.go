package exercise

// Ex08 /* เขียน function หาราคารวมของ Product ทั้งหมด
func Ex08(productList []Product) float64 {
	total := 0.0
	for _, product := range productList {
		total += *product.ProductPrice
	}
	return total
}
