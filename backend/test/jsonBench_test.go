package test

import (
	"backend/internal/common/utils/datetimeutils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func TestJsonRaw(t *testing.T) {
	jsonStr := `[1221, {"id":1},{"double1":-20.791168813592314,"quality":0,"timeStamp":1768285507000 }]`
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jsonStr)))
	token, err := decoder.Token()
	if err != nil {
		t.Errorf("json parse error: %s", err)
		return
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		t.Errorf("json parse error: token=%s", token)
		return
	}
	for i := 0; decoder.More(); i++ {
		var raw json.RawMessage
		err := decoder.Decode(&raw)
		if err != nil {
			t.Errorf("Error decoding: %v", err)
		} else {
			t.Logf("[%d]: %s\n", i, string(raw))
		}
	}
	out := bytes.NewBuffer(make([]byte, 0, 128))
	err = stripFirstAndLastCharStream(bytes.NewBuffer([]byte(jsonStr)), out)
	if err != nil {
		t.Errorf("Error copying: %v", err)
	} else {
		t.Log("复制结果：", out.String())
	}
}

// 流式版本，适合大文件
func stripFirstAndLastCharStream(src io.Reader, dst io.Writer) error {
	// 跳过第一个字符
	if _, err := io.CopyN(io.Discard, src, 1); err != nil {
		return err
	}

	// 使用缓冲区读取
	buf := make([]byte, 32*1024)
	var prevChunk []byte

	for {
		n, err := src.Read(buf)
		if n > 0 {
			// 如果有上一次的chunk，先写入
			if prevChunk != nil {
				if _, werr := dst.Write(prevChunk); werr != nil {
					return werr
				}
			}
			// 保存当前chunk，等待下一次读取
			prevChunk = make([]byte, n)
			copy(prevChunk, buf[:n])
		}

		if err == io.EOF {
			// 写入除了最后一个字节的所有数据
			if prevChunk != nil && len(prevChunk) > 0 {
				dst.Write(prevChunk[:len(prevChunk)-1])
			}
			return nil
		}

		if err != nil {
			return err
		}
	}
}
func copySkipFirstAndLast(src io.Reader, dst io.Writer) error {
	// 跳过第一个字符
	if _, err := io.CopyN(ioutil.Discard, src, 1); err != nil {
		return err
	}

	// 创建一个管道，在管道中处理最后一个字符
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		buf := make([]byte, 32*1024)
		var lastChunk []byte
		var totalBytes int64

		for {
			n, err := src.Read(buf)
			if n > 0 {
				if lastChunk != nil {
					// 写入上一次的chunk
					if _, werr := pw.Write(lastChunk); werr != nil {
						return
					}
				}
				lastChunk = make([]byte, n)
				copy(lastChunk, buf[:n])
				totalBytes += int64(n)
			}

			if err == io.EOF {
				// 写入除了最后一个字节的所有数据
				if lastChunk != nil && len(lastChunk) > 0 {
					pw.Write(lastChunk[:len(lastChunk)-1])
				}
				return
			}

			if err != nil {
				return
			}
		}
	}()

	_, err := io.Copy(dst, pr)
	return err
}

const SYS_FIELD_CREATE_TIME = "timeStamp"

func TestJsonGetTs(t *testing.T) {
	jsonStr := `{"double1":-20.791168813592314,"quality":0,"timeStamp":1768285507000 }`
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		t.Fatal(err)
	}
	ts := time.Now()
	m[SYS_FIELD_CREATE_TIME] = &ts
	t.Log("ts:", parseTimestamp(m[SYS_FIELD_CREATE_TIME]))
}
func newErr() (er error) {
	return er
}
func parseTimestamp(curT any) (ct int64) {
	if curT == nil {
		return -1
	} else if Float, isFloat := curT.(float64); isFloat { // json unmarshal 来的都是 float64 类型
		ct = int64(Float)
	} else if Long, isLong := curT.(int64); isLong {
		ct = Long
	} else {
		str := fmt.Sprint(curT)
		Double, err := strconv.ParseFloat(str, 64)
		if err != nil {
			ct = -1
			if dt, dtEr := datetimeutils.ParseDate(str); dtEr == nil && dt.Year() > 1970 {
				ct = dt.UnixMilli()
			}
		} else if Int := int64(Double); Int > 1100000000000 {
			ct = Int
		}
	}
	if ct < 1100000000000 || ct > 11000000000001 {
		return -1
	}
	return ct
}
func BenchmarkJson(b *testing.B) {
	jsonStr := `{"timeStamp":1768288504446, "wet":168.931,"qos":0}`
	b.StartTimer()
	//BenchmarkJson-12    	  683923	      1636 ns/op
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &m)
		if err != nil {
			b.Fatal(err)
		}
		ts, hasTs := m[SYS_FIELD_CREATE_TIME]
		value, hasValue := m["wet"].(float64)
		if !hasTs || !hasValue {
			b.Fatal("wet", ts, value)
		}
	}
}
func Benchmark_jsoniter(b *testing.B) {
	jsonStr := `{"timeStamp":1768288504446, "wet":168.931,"qos":0}`
	b.StartTimer()
	//Benchmark_jsoniter-12    	 1504185	       783.6 ns/op
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &m)
		if err != nil {
			b.Fatal(err)
		}
		ts, hasTs := m[SYS_FIELD_CREATE_TIME]
		value, hasValue := m["wet"].(float64)
		if !hasTs || !hasValue {
			b.Fatal("wet", ts, value)
		}
	}
}
