package exercise

import (
	"encoding/json"
)

// Ex05 /* แปลง map[string]any ทีรับเข้ามา เป็น Product แล้ว return ออกไป
func Ex05(productMap map[string]any) Product {
	var product Product
	jsonData, _ := json.Marshal(productMap)
	json.Unmarshal(jsonData, &product)
	return product
}
