package sourceflow

import (
	dao "backend/internal/repo/relationDB"
	"bytes"
	"io"
	"testing"

	"gitee.com/unitedrhino/share/conf"
)

func TestNodeRedFlowExport(t *testing.T) {
	dao.InitDbConfig(conf.Database{DSN: "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos"})
	defer removeMock()
	initMock(t)

	groupIds := []int64{3, 4}
	ids := []int64{2004190137830871040, 2004190506053013504, 1963148065022947328, 1910955436975697920}
	{
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		NodeRedFlowExport(t.Context(), groupIds, ids, true, "http://nodered:1880")(buf)
		t.Log(buf.String())
	}
}
func Test_exportNodesByFlows(t *testing.T) {
	jsonWriter := bytes.NewBuffer(nil)
	allFlows := `[
{"type":"tab","label":"Metric/jkjkjkjk","id":"7ed2a8bc2cb737ca","info":"","disabled":false},{"crontab":"","id":"1f54f0b41a854102","name":"","once":false,"onceDelay":1,"payload":"","payloadType":"date","props":[{"p":"payload"}],"repeat":"10","topic":"","type":"inject","wires":[["ad39c9112de34497"]],"x":320,"y":160,"z":"7ed2a8bc2cb737ca"},{"finalize":"","func":"// 随机字符串\nfunction randomString() {\n    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';\n    let result = '';\n    for (let i = 0; i < 20; i++) {\n        result += chars.charAt(Math.floor(Math.random() * chars.length));\n    }\n    return result;\n}\n// 100以内随机数字\nfunction generateRandomNumber() {\n    return Math.floor(Math.random() * 100);\n}\n// 随机生成100以内浮点数，保留2位小数\nfunction generateRandomFloatWithTwoDecimals() {\n    const randomFloat = Math.random() * 100;\n    return randomFloat.toFixed(2);\n}\n// 对当前时间格式化\nfunction formatCurDate() {return Date.now();\n}\n\nfunction getBool() {\n    var randomInt = generateRandomNumber();\n    return randomInt > 50;\n}\nmsg.topic='Metric/jkjkjkjk';\nmsg.payload = {\n'iii': generateRandomNumber() \n};\n\nreturn msg;","id":"ad39c9112de34497","initialize":"","libs":[],"name":"mock data","noerr":0,"outputs":1,"timeout":0,"type":"function","wires":[["5d2841b045ae4720"]],"x":520,"y":160,"z":"7ed2a8bc2cb737ca"},{"broker":"85bb67b2dbefe3ba","contentType":"","correl":"","expiry":"","id":"d56d84e17dc0471b","name":"","qos":"","respTopic":"","retain":"","topic":"","type":"mqtt out","userProps":"","wires":[],"x":990,"y":160,"z":"7ed2a8bc2cb737ca"},{"id":"5d2841b045ae4720","modelShowName":"","name":"","protocol":"mock","selectedModel":"Metric/jkjkjkjk","selectedModelAlias":"_jkjkjkjk_1df6e992e3ce4bba8bb0","tableValid":true,"type":"supmodel","wires":[["d56d84e17dc0471b"],[]],"x":750,"y":160,"z":"7ed2a8bc2cb737ca"},{"autoConnect":true,"autoUnsubscribe":true,"birthMsg":{},"birthPayload":"","birthQos":"0","birthRetain":"false","birthTopic":"","broker":"emqx","cleansession":true,"clientid":"","closeMsg":{},"closePayload":"","closeQos":"0","closeRetain":"false","closeTopic":"","id":"85bb67b2dbefe3ba","keepalive":"60","name":"","port":"1883","protocolVersion":"4","sessionExpiry":"","type":"mqtt-broker","userProps":"","usetls":false,"willMsg":{},"willPayload":"","willQos":"0","willRetain":"false","willTopic":""},{"type":"tab","label":"AIgc2","id":"d7ac0a1e7833fd91","info":"auto mock for AIgc2","disabled":false},{"crontab":"","id":"84f14bd8c68143cf","name":"","once":false,"onceDelay":1,"payload":"","payloadType":"date","props":[{"p":"payload"}],"repeat":"10","topic":"","type":"inject","wires":[["70f8e9f9a57a45f1"]],"x":320,"y":160,"z":"d7ac0a1e7833fd91"},{"finalize":"","func":"// 随机字符串\nfunction randomString() {\n    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';\n    let result = '';\n    for (let i = 0; i < 20; i++) {\n        result += chars.charAt(Math.floor(Math.random() * chars.length));\n    }\n    return result;\n}\n// 100以内随机数字\nfunction generateRandomNumber() {\n    return Math.floor(Math.random() * 100);\n}\n// 随机生成100以内浮点数，保留2位小数\nfunction generateRandomFloatWithTwoDecimals() {\n    const randomFloat = Math.random() * 100;\n    return randomFloat.toFixed(2);\n}\n// 对当前时间格式化\nfunction formatCurDate() {return Date.now();\n}\n\nfunction getBool() {\n    var randomInt = generateRandomNumber();\n    return randomInt > 50;\n}\nmsg.topic='hktest/State/2AIgc';\nmsg.payload = {\n'json': randomString() \n};\n\nreturn msg;","id":"70f8e9f9a57a45f1","initialize":"","libs":[],"name":"mock data","noerr":0,"outputs":1,"timeout":0,"type":"function","wires":[["40610463973448b8"]],"x":520,"y":160,"z":"d7ac0a1e7833fd91"},{"broker":"85bb67b2dbefe3ba","contentType":"","correl":"","expiry":"","id":"0c84eb6631b64829","name":"","qos":"","respTopic":"","retain":"","topic":"","type":"mqtt out","userProps":"","wires":[],"x":990,"y":160,"z":"d7ac0a1e7833fd91"},{"id":"40610463973448b8","modelShowName":"","name":"","protocol":"mock","selectedModel":"hktest/State/2AIgc","selectedModelAlias":"AIgc2","tableValid":true,"type":"supmodel","wires":[["0c84eb6631b64829"],[]],"x":750,"y":160,"z":"d7ac0a1e7833fd91"}
]`
	exportNodesByFlows(t.Context(), jsonWriter, io.NopCloser(bytes.NewBuffer([]byte(allFlows))), map[string]bool{"d7ac0a1e7833fd91": true}, nil)
	t.Log(jsonWriter.String())
}
