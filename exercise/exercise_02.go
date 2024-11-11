package exercise

import "time"

// Ex02 /* ให้ return slice of product ตามโครงสร้าง json ดังนี้
//
//	[{
//		"productId": 1111,
//		"productCode": "COKE",
//		"productName": "Coca cola",
//		"productPrice": 15.00,
//		"productUnit": "CAN",
//		"createdBy": "Josh",
//		"createdAt": <today>,
//		"updatedBy": "Sarah",
//		"updatedAt": <today>
//	}, {
//		"productId": 2222,
//		"productCode": "PEPSI",
//		"productName": "Pepsi",
//		"productPrice": 15.50,
//		"productUnit": "BOTTOM",
//		"createdBy": "John",
//		"createdAt": <today>,
//		"updatedBy": null,
//		"updatedAt": null
//	}, {
//		"productId": 3333,
//		"productCode": "SPRITE",
//		"productName": "Sprite",
//		"productPrice": null,
//		"productUnit": "GLASS",
//		"createdBy": "Peter",
//		"createdAt": <today>,
//		"updatedBy": null,
//		"updatedAt": null
//	}]
//
// */
func Ex02() []Product {
	now := time.Now()
	price1 := 15.00
	price2 := 15.50
	sarah := "Sarah"

	return []Product{
		{
			ProductId:    1111,
			ProductCode:  "COKE",
			ProductName:  "Coca cola",
			ProductPrice: &price1,
			ProductUnit:  "CAN",
			CreatedBy:    "Josh",
			CreatedAt:    now,
			UpdatedBy:    &sarah,
			UpdatedAt:    &now,
		},
		{
			ProductId:    2222,
			ProductCode:  "PEPSI",
			ProductName:  "Pepsi",
			ProductPrice: &price2,
			ProductUnit:  "BOTTOM",
			CreatedBy:    "John",
			CreatedAt:    now,
			UpdatedBy:    nil,
			UpdatedAt:    nil,
		},
		{
			ProductId:    3333,
			ProductCode:  "SPRITE",
			ProductName:  "Sprite",
			ProductPrice: nil,
			ProductUnit:  "GLASS",
			CreatedBy:    "Peter",
			CreatedAt:    now,
			UpdatedBy:    nil,
			UpdatedAt:    nil,
		},
	}
}
