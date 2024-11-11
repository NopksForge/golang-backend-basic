package exercise

import "encoding/json"

// Ex04 /* แปลง json string ทีรับเข้ามา เป็น slice of Product แล้ว return ออกไป
func Ex04(jsonString string) []Product {
	var products []Product
	json.Unmarshal([]byte(jsonString), &products)
	return products
}
