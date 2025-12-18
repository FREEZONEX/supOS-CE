package postgresql

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestPaseDbUrl(t *testing.T) {
	// 测试示例
	testUrl := "postgres://root:Supos123@100.100.100.22:34099/postgres?search_path=supos"
	result := ParseDbUrlProperties(testUrl)

	fmt.Printf("解析结果:\n")
	fmt.Printf("Url: %s\n", result.Url)
	fmt.Printf("UserName: %s\n", result.UserName)
	fmt.Printf("Password: %s\n", result.Password)
	fmt.Printf("Schema: %s\n", result.Schema)
	fmt.Printf("Host: %s\n", result.HostPort)

	conf, er := pgx.ParseConfig(testUrl)
	if er != nil {
		t.Fatal(er)
	}
	t.Logf("conf: %+v\n", *conf)
}
