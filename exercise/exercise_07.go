package exercise

// Ex07 /* เขียนฟังก์ชัน ลบ product จาก slice ของ productCodeToDeleted ออกจาก productList ที่มีอยู่ แล้ว return slice of Product ที่่ไมถูกลบออก
func Ex07(productList []Product, productCodeToDeleted []string) []Product {
	toDelete := make(map[string]bool)
	for _, code := range productCodeToDeleted {
		toDelete[code] = true
	}

	result := make([]Product, 0, len(productList))

	for _, product := range productList {
		if !toDelete[product.ProductCode] {
			result = append(result, product)
		}
	}

	return result
}
