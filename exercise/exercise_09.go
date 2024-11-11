package exercise

import "sort"

// Ex09 /* เขียน function เพื่อเรียงลำดับ slice ของ Product ใหม่ โดยเรียงจาก ProductName, ProductCode ตามลำดับ
func Ex09(productList []Product) []Product {
	sort.Slice(productList, func(i, j int) bool {
		isNameLess := productList[i].ProductName < productList[j].ProductName
		isNameEqual := productList[i].ProductName == productList[j].ProductName
		isCodeLess := productList[i].ProductCode < productList[j].ProductCode

		return isNameLess || (isNameEqual && isCodeLess)
	})
	return productList
}
