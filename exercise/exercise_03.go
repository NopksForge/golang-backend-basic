package exercise

import "time"

// Ex03 /* ให้ return slice of order ตามโครงสร้าง json ดังนี้
//
//	[{
//		"orderId": 12345,
//		"orderNo": "ORD0001",
//		"productList": [
//			{
//				"productId": 1111,
//				"productCode": "COKE",
//				"productName": "Coca cola",
//				"productPrice": 15.00,
//				"productUnit": "CAN",
//				"createdBy": "Josh",
//				"createdAt": <today>,
//				"updatedBy": "Sarah",
//				"updatedAt": <today>
//			}, {
//				"productId": 2222,
//				"productCode": "PEPSI",
//				"productName": "Pepsi",
//				"productPrice": 15.50,
//				"productUnit": "BOTTOM",
//				"createdBy": "John",
//				"createdAt": <today>,
//				"updatedBy": null,
//				"updatedAt": null
//			}
//		],
//		"createdBy": "Tony",
//		"createdAt": <today>,
//		"updatedBy": null,
//		"updatedAt": null
//	}]
//
// */
func Ex03() []Order {
	now := time.Now()
	price1 := 15.00
	price2 := 15.50
	sarah := "Sarah"

	return []Order{
		{
			OrderId: 12345,
			OrderNo: "ORD0001",
			ProductList: []Product{
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
			},
			CreatedBy: "Tony",
			CreatedAt: now,
			UpdatedBy: nil,
			UpdatedAt: nil,
		},
	}
}
