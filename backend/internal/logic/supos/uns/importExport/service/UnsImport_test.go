package service

import (
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	"backend/internal/types"
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestImportFile(t *testing.T) {
	file, er := os.Open("F:\\data\\uns-uat.json")
	if er != nil {
		t.Fatal(er)
	}
	defer file.Close()
	err := jsonstream.DecodeJsonTreeToFlat(file, 100, node2vo, func(rd int64, propName string, ns []*types.CreateTopicDto) {
		jsonBytes, _ := json.Marshal(ns)
		t.Logf("readSize=%d, prop=%s , Nodes[%d]: %v", rd, propName, len(ns), string(jsonBytes))
	}, func(node *FileData) {
		err := node.Error
		node.Error = ""
		jsonBytes, _ := json.Marshal(node)
		t.Log("ErrorNode: ", err, string(jsonBytes))
	})
	if err != nil {
		t.Fatalf("Error parsing JSON: %v\n", err)
	}
}

func TestDecodeStreamedJson(t *testing.T) {
	bigJson := bytes.NewBuffer([]byte(__realJson))
	err := jsonstream.DecodeJsonTreeToFlat(bigJson, 10, node2vo, func(rd int64, propName string, ns []*types.CreateTopicDto) {
		jsonBytes, _ := json.Marshal(ns)
		t.Logf("readSize=%d, prop=%s , Nodes[%d]: %v", rd, propName, len(ns), string(jsonBytes))
	}, func(node *FileData) {
		err := node.Error
		node.Error = ""
		jsonBytes, _ := json.Marshal(node)
		t.Log("ErrorNode: ", err, string(jsonBytes))
	})
	if err != nil {
		t.Fatalf("Error parsing JSON: %v\n", err)
	}
}
func TestLabelImport(t *testing.T) {
	__json := `{
"Label":[
{"name":"L1"},
{"name":"seqx"},
{"name":"fas"}
],"Template":[
{"name":"wet_it","alias":"__1c200aefe8b248e6b640","fields":[{"name":"wet","type":"FLOAT"},{"name":"d111","type":"INTEGER"}],"description":"湿度模版9","frequency":"1m"},
{"name":"wet_it","alias":"__76b5c7a847454fc4b77a","fields":[{"name":"wet","type":"FLOAT"},{"name":"csdw","type":"LONG"}],"description":"湿度2"},
{"name":"312312","alias":"__69bd7f8b4b594a708f62","fields":[{"name":"wet","type":"FLOAT"},{"name":"val","type":"LONG"}],"description":"湿度2fads"},
{"name":"aaa","alias":"T__9ede38958a724b1993a0","fields":[{"name":"g","type":"FLOAT","displayName":"a","remark":"b"},{"name":"val11","type":"INTEGER"}]},
{"name":"文件夹","alias":"Twenjianjia_293658187a4d43cea27f","fields":[{"name":"a","type":"DATETIME"}]},
{"name":"文件2","alias":"Twenjian_8ff63f10b0d044e792ba","fields":[{"name":"a","type":"DATETIME"}]},
{"name":"AlarmTemplate","alias":"_alarm_model","fields":[{"name":"current_value","type":"FLOAT"},{"name":"is_alarm","type":"BOOLEAN"},{"name":"limit_value","type":"FLOAT"},{"name":"read_status","type":"BOOLEAN"},{"name":"uns","type":"LONG"},{"name":"uns_path","type":"STRING","maxLen":512}],"description":"Alarm Model"}
],"UNS":[
 {"type":"folder","name":"状态","alias":"_state_","displayName":"状态","children":[{"type":"file","name":"121","alias":"_121_77ffd11e5c784f3b8396","fields":[{"name":"json","type":"STRING"}],"dataType":"JSONB_TYPE","generateDashboard":"TRUE","enableHistory":"FALSE","mockData":"FALSE","topicType":"STATE"}]}
 ]
}
`
	bigJson := bytes.NewBuffer([]byte(__json))
	err := jsonstream.DecodeJsonTreeToFlat(bigJson, 100, node2vo, func(rd int64, propName string, ns []*types.CreateTopicDto) {
		jsonBytes, _ := json.Marshal(ns)
		t.Logf("readSize=%d, prop=%s , Nodes[%d]: %v", rd, propName, len(ns), string(jsonBytes))
	}, func(node *FileData) {
		err := node.Error
		node.Error = ""
		jsonBytes, _ := json.Marshal(node)
		t.Log("ErrorNode: ", err, string(jsonBytes))
	})
	if err != nil {
		t.Fatalf("Error parsing JSON: %v\n", err)
	}
}

var __realJson = `
{
  "notes": "type:folder|file,topicType:STATE|ACTION|METRIC,dataType:TEMPLATE_TYPE|TIME_SEQUENCE_TYPE|RELATION_TYPE|CALCULATION_REAL_TYPE|CALCULATION_HIST_TYPE|MERGE_TYPE|CITING_TYPE|JSONB_TYPE|,fields.type:INTEGER|LONG|FLOAT|DOUBLE|BOOLEAN|DATETIME|STRING",
  "Template": [],
  "Label": [],
  "UNS": [
    {
      "name": "v1",
      "type": "folder",
      "children": [
        {
          "name": "Plant_Name",
          "type": "folder",
          "children": [
            {
              "name": "SMT-Area-1",
              "type": "folder",
              "children": [
                {
                  "name": "SMT-Line-1",
                  "type": "folder",
                  "children": [
                    {
                      "name": "Printer-Cell",
                      "type": "folder",
                      "children": [
                        {
                          "name": "Printer01",
                          "type": "folder",
                          "children": [
                            {
                              "name": "State",
                              "type": "folder",
                              "topicType": "STATE",
                              "children": [
                                {
                                  "name": "current_job",
                                  "type": "file",
                                  "topicType": "STATE",
                                  "dataType": "RELATION_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "job_id",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "product_id",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "planned_quantity",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "completed_quantity",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "status",
                                      "type": "LONG"
                                    }
                                  ]
                                },
                                {
                                  "name": "alarm_status",
                                  "type": "file",
                                  "topicType": "STATE",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE"
                                }
                              ]
                            },
                            {
                              "name": "Action",
                              "type": "folder",
                              "topicType": "ACTION",
                              "children": [
                                {
                                  "name": "start_job",
                                  "type": "file",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                },
                                {
                                  "name": "stop_job",
                                  "type": "file",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                }
                              ]
                            },
                            {
                              "name": "Metric",
                              "type": "folder",
                              "topicType": "METRIC",
                              "children": [
                                {
                                  "name": "board_cycle_time",
                                  "type": "file",
                                  "topicType": "METRIC",
                                  "dataType": "TIME_SEQUENCE_TYPE",
                                  "generateDashboard": "TRUE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "cycle_time_ms",
                                      "type": "LONG"
                                    }
                                  ]
                                },
                                {
                                  "name": "boards_count",
                                  "type": "file",
                                  "topicType": "METRIC",
                                  "dataType": "TIME_SEQUENCE_TYPE",
                                  "generateDashboard": "TRUE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "good_count",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "ng_count",
                                      "type": "LONG"
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
`
var __unsJson = `[
                    {
                      "alias": "wenjianjiaA_686e7fad643b4db0b11e",
                      "displayName": "指标",
                      "type": "folder",
                      "name": "指标",
                      "children": [
                        {
                          "alias": "wenjianjiaA_96a9ab047699462faa27",
                          "fields": [
                            {
                              "name": "aaa",
                              "type": "INTEGER"
                            },
                            {
                              "name": "status",
                              "type": "LONG"
                            }
                          ],
                          "dataType": "TIME_SEQUENCE_TYPE",
                          "generateDashboard": "FALSE",
                          "enableHistory": "FALSE",
                          "mockData": "FALSE",
                          "type": "file",
                          "name": "时序1",
                          "topicType": "METRIC"
                        },
                        {
                          "alias": "wenjianjiaA_c515c3d4e61a4a8f9024",
                          "fields": [
                            {
                              "name": "a",
                              "type": "LONG"
                            }
                          ],
                          "dataType": "TIME_SEQUENCE_TYPE",
                          "generateDashboard": "TRUE",
                          "enableHistory": "FALSE",
                          "mockData": "FALSE",
                          "type": "file",
                          "name": "asdfa",
                          "topicType": "METRIC"
                        },
                        {
                          "alias": "wenjianjiaA_ce18b74fd6c74ef2b5a4",
                          "fields": [
                            {
                              "name": "a",
                              "type": "LONG"
                            }
                          ],
                          "dataType": "TIME_SEQUENCE_TYPE",
                          "generateDashboard": "TRUE",
                          "enableHistory": "FALSE",
                          "mockData": "FALSE",
                          "type": "file",
                          "name": "asdfa",
                          "topicType": "METRIC"
                        }
                      ],
                      "topicType": "METRIC"
                    }
]
`
