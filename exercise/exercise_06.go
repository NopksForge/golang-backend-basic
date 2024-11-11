package exercise

// Ex06 /* เขียนฟังก์ชัน ค้นหา product code จาก productList ที่มีอยู่ แล้ว return index ของ slice ที่เจอออกไป
func Ex06(productList []Product, productCode string) int {
	for i, product := range productList {
		if product.ProductCode == productCode {
			return i
		}
	}
	return -1
}
