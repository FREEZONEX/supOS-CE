package types

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestFt(t *testing.T) {
	fs := []string{"integer", "boolean", "double", "int", "FloatX"}
	for _, f := range fs {
		rs, Ok := GetFieldTypeByNameIgnoreCase(f)
		t.Log(f, rs, Ok)
	}
}
func TestJsonEncodePoint(t *testing.T) {
	type UserResp struct {
		ID       *int64 `json:"id,string"` // 注意：这个标签在某些情况下有效
		Username string `json:"name"`
	}
	id := int64(1234567890123456789)
	usr := &UserResp{
		ID:       &id,
		Username: "Lucy",
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(usr, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("序列化结果:")
	fmt.Println(string(jsonData))

	{
		var decodedProduct UserResp
		jsonStr := `{
        "id": "1876543210987654321",
        "name": "Tablet"
    }`

		if err := json.Unmarshal([]byte(jsonStr), &decodedProduct); err != nil {
			panic(err)
		}

		fmt.Printf("\n反序列化结果1:\n")
		fmt.Printf("Id: %d (类型: %T)\n", *decodedProduct.ID, decodedProduct.ID)
		fmt.Printf("Price: %v (类型: %T)\n", decodedProduct.Username, decodedProduct.Username)
	}

	{
		var decodedProduct UserResp
		jsonStr := `{
        "id": "0",
        "name": "Tablet"
    }`

		if err := json.Unmarshal([]byte(jsonStr), &decodedProduct); err != nil {
			log.Println(err)
		}

		fmt.Printf("\n反序列化结果2:\n")
		if decodedProduct.ID != nil {
			fmt.Printf("Id: %d (类型: %T)\n", *decodedProduct.ID, decodedProduct.ID)
		}
		fmt.Printf("Price: %v (类型: %T)\n", decodedProduct.Username, decodedProduct.Username)
	}
}
func TestJsonEncodePointErr(t *testing.T) {
	type UnsTreeCondition struct {
		ParentId *int64 `json:"parentId,optional"`
		Keyword  string `json:"keyword,optional"`
	}
	id := int64(1234567890123456789)
	usr := &UnsTreeCondition{
		ParentId: &id,
		Keyword:  "Lucy",
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(usr, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("序列化结果:")
	fmt.Println(string(jsonData))

	{
		var decodedProduct UnsTreeCondition
		jsonStr := `{
        "parentId": "1876543210987654321",
        "keyword": "Tablet"
    }`

		if err := json.Unmarshal([]byte(jsonStr), &decodedProduct); err != nil {
			panic(err)
		}

		fmt.Printf("\n反序列化结果1:\n")
		fmt.Printf("Id: %d (类型: %T)\n", *decodedProduct.ParentId, decodedProduct.ParentId)
		fmt.Printf("Price: %v (类型: %T)\n", decodedProduct.Keyword, decodedProduct.Keyword)
	}

	{
		var decodedProduct UnsTreeCondition
		jsonStr := `{
        "parentId": "0",
        "keyword": "Tablet"
    }`

		if err := json.Unmarshal([]byte(jsonStr), &decodedProduct); err != nil {
			log.Println(err)
		}

		fmt.Printf("\n反序列化结果2:\n")
		if decodedProduct.ParentId != nil {
			fmt.Printf("Id: %d (类型: %T)\n", *decodedProduct.ParentId, decodedProduct.ParentId)
		}
		fmt.Printf("Price: %v (类型: %T)\n", decodedProduct.Keyword, decodedProduct.Keyword)
	}
}
func TestJsonEncodeStr(t *testing.T) {
	type UserResp struct {
		ID       int64  `json:"id,string"` // 注意：这个标签在某些情况下有效
		Username string `json:"name"`
	}
	usr := &UserResp{
		ID:       1234567890123456789,
		Username: "Lucy",
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(usr, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("序列化结果:")
	fmt.Println(string(jsonData))

	var decodedProduct UserResp
	jsonStr := `{
        "id": "1876543210987654321",
        "name": "Tablet"
    }`

	if err := json.Unmarshal([]byte(jsonStr), &decodedProduct); err != nil {
		panic(err)
	}

	fmt.Printf("\n反序列化结果:\n")
	fmt.Printf("Id: %d (类型: %T)\n", decodedProduct.ID, decodedProduct.ID)
	fmt.Printf("Price: %v (类型: %T)\n", decodedProduct.Username, decodedProduct.Username)
}
func TestJsonCodec(t *testing.T) {

	// 使用自定义类型的结构体
	type Product struct {
		ID          Int64   `json:"id"`
		Price       Float64 `json:"price"`
		Weight      Float32 `json:"weight"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
	}
	// 创建示例数据
	product := Product{
		ID:          Int64(1234567890123456789),
		Price:       Float64(299.99),
		Weight:      Float32(15.5),
		Name:        "Laptop",
		Description: "High-performance laptop",
	}

	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("序列化结果:")
	fmt.Println(string(jsonData))

	// 输出:
	// {
	//   "id": "1234567890123456789",
	//   "price": "299.99",
	//   "weight": "15.5",
	//   "name": "Laptop",
	//   "description": "High-performance laptop"
	// }

	// 反序列化测试
	var decodedProduct Product
	jsonStr := `{
        "id": "1876543210987654321",
        "price": "199.50",
        "weight": "12.25",
        "name": "Tablet",
        "description": "Portable tablet"
    }`

	if err := json.Unmarshal([]byte(jsonStr), &decodedProduct); err != nil {
		panic(err)
	}

	fmt.Printf("\n反序列化结果:\n")
	fmt.Printf("Id: %d (类型: %T)\n", decodedProduct.ID, decodedProduct.ID)
	fmt.Printf("Price: %f (类型: %T)\n", decodedProduct.Price, decodedProduct.Price)
	fmt.Printf("Weight: %f (类型: %T)\n", decodedProduct.Weight, decodedProduct.Weight)
}
