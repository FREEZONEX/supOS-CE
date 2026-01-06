package mock

import (
	"backend/internal/types"
	"math/rand"
	"os"
	"time"
)

type MockWeatherDTO struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type MockOpcuaDTO struct {
	VariableName  string `json:"variableName"`
	DataType      string `json:"dataType"`
	VariableRange string `json:"variableRange"`
}

type MockOrderDTO struct {
	ID                 int64   `json:"id"`
	Customer           string  `json:"customer"`
	SubmitTime         string  `json:"submitTime"`
	SkuCount           int     `json:"skuCount"`
	Count              int     `json:"count"`
	DeliveryTime       string  `json:"deliveryTime"`
	OrderPrice         float64 `json:"orderPrice"`
	UserID             int     `json:"userID"`
	UserName           string  `json:"userName"`
	Process            int     `json:"process"`
	ProducedCount      int     `json:"producedCount"`
	ShippedCount       int     `json:"shippedCount"`
	OrderType          int     `json:"orderType"`
	UserMaterialStatus int     `json:"userMaterialStatus"`
	Status             int     `json:"status"`
	ReceivableAmount   float64 `json:"receivableAmount"`
	CollectedAmount    float64 `json:"collectedAmount"`
	OverdueAmount      float64 `json:"overdueAmount"`
}

var (
	customers   = []string{"河北雄安素水互联网科技有限公司", "上海万汇云建筑科技服务有限公司", "兰州美尔康生物科技有限公司", "甘肃知行合创生物科技有限公司", "兰州宁远化工有限责任公司", "乌鲁木齐智星宏元科技有限公司", "甘肃润康药业有限公司"}
	customersEN = []string{"Google", "Microsoft", "Apple", "Amazon", "Tesla"}
	userNames   = []string{"张学文", "李秋菊"}
	userNamesEN = []string{"Elon Musk", "Bill Gates", "Mark Zuckerberg"}
)

// NewMockOrderDTO creates a new MockOrderDTO with randomized data.
func NewMockOrderDTO() *MockOrderDTO {
	price := float64(rand.Intn(990001) + 10000)
	return &MockOrderDTO{
		ID:                 rand.Int63n(1000000),
		Customer:           getRandomCustomer(),
		SubmitTime:         getRandomDateTime(),
		SkuCount:           rand.Intn(21) + 10,
		Count:              rand.Intn(491) + 10,
		DeliveryTime:       getRandomDateTime(),
		OrderPrice:         price,
		UserID:             rand.Intn(10) + 1,
		UserName:           getRandomUserName(),
		Process:            rand.Intn(101),
		ProducedCount:      rand.Intn(291) + 10,
		ShippedCount:       rand.Intn(191) + 10,
		OrderType:          1,
		UserMaterialStatus: rand.Intn(2),
		Status:             1,
		ReceivableAmount:   price,
		CollectedAmount:    float64(rand.Intn(10001)),
		OverdueAmount:      float64(rand.Intn(10001)),
	}
}

func getRandomCustomer() string {
	if os.Getenv("SYS_OS_LANG") == "en-US" {
		return customersEN[rand.Intn(len(customersEN))]
	}
	return customers[rand.Intn(len(customers))]
}

func getRandomUserName() string {
	if os.Getenv("SYS_OS_LANG") == "en-US" {
		return userNamesEN[rand.Intn(len(userNamesEN))]
	}
	return userNames[rand.Intn(len(userNames))]
}

func getRandomDateTime() string {
	min := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format(time.RFC3339)
}

// MockDemoDTO corresponds to the Java MockDemoDTO class.
type MockDemoDTO struct {
	TimeStamp            int64   `json:"timeStamp"`
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	InstalledCapacity    float64 `json:"installedCapacity"`
	DailyPowerGeneration string  `json:"dailyPowerGeneration"`
	Owner                string  `json:"owner"`
}

// NewMockDemoDTO creates a new MockDemoDTO with randomized data.
func NewMockDemoDTO() *MockDemoDTO {
	name := "光伏基站-"
	owner := "张灿"
	if os.Getenv("SYS_OS_LANG") == "en-US" {
		name = "Photovoltaic base station-"
		owner = "Jack"
	}

	return &MockDemoDTO{
		TimeStamp:            time.Now().UnixMilli(),
		ID:                   rand.Int(),
		Name:                 name + randomString(4),
		InstalledCapacity:    float64(rand.Intn(1000000)) / 100.0, // Approximates randomDouble
		DailyPowerGeneration: randomNumbers(5),
		Owner:                owner,
	}
}

// ConvertOrderToParams converts the DTO to a slice of any for database operations.
func (m *MockDemoDTO) ConvertOrderToParams(srcJdbcType *types.SrcJdbcType) []any {
	if srcJdbcType.TypeCode() == types.SrcJdbcTypeTdEngine.TypeCode() {
		return []any{m.TimeStamp, m.ID, m.Name, m.InstalledCapacity, m.DailyPowerGeneration, m.Owner}
	} else if srcJdbcType.TypeCode() == types.SrcJdbcTypeTimeScaleDB.TypeCode() {
		return []any{time.UnixMilli(m.TimeStamp), m.ID, m.Name, m.InstalledCapacity, m.DailyPowerGeneration, m.Owner}
	}
	return nil
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomNumbers(n int) string {
	var letters = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
