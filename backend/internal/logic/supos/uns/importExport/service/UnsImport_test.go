package service

import (
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	"backend/internal/types"
	"bytes"
	"encoding/json"
	"testing"
)

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
