export const placeholder = `{
  "notes": "type:PATH|TOPIC,topicType:STATE|ACTION|METRIC,fields.type:INTEGER|STRING|FLOAT|DOUBLE|BOOLEAN|LONG|DATETIME",
  "templates": [
    {
      "alias": "basic_order",
      "name": "BasicTemplate",
      "fields": [
        {
          "name": "order_id",
          "type": "STRING"
        }
      ]
    }
  ],
  "labels": [
    {
      "name": "alpha"
    },
    {
      "name": "beta"
    }
  ],
  "namespace": [
    {
      "name": "v1",
      "type": "PATH",
      "children": [
        {
          "name": "Plant_Name",
          "type": "PATH",
          "children": [
            {
              "name": "SMT-Area-1",
              "type": "PATH",
              "children": [
                {
                  "name": "State",
                  "type": "PATH",
                  "topicType": "STATE",
                  "labels": "alpha,beta",
                  "children": [
                    {
                      "name": "current_job",
                      "type": "TOPIC",
                      "extendProperties": {
                        "ext1": "value1"
                      },
                      "fields": [
                        {
                          "name": "job_id",
                          "type": "INTEGER"
                        },
                        {
                          "name": "status",
                          "type": "FLOAT"
                        },
                        {
                          "name": "created_at",
                          "type": "DATETIME"
                        }
                      ]
                    },
                    {
                      "name": "productA",
                      "type": "TOPIC",
                      "description": "desc",
                      "enableHistory": "TRUE",
                      "fields": [
                        {
                          "name": "product_id",
                          "type": "INTEGER"
                        },
                        {
                          "name": "status",
                          "type": "STRING"
                        }
                      ]
                    }
                  ]
                },
                {
                  "name": "Action",
                  "type": "PATH",
                  "topicType": "ACTION",
                  "children": [
                    {
                      "name": "start_job",
                      "type": "TOPIC",
                      "fields": [
                        {
                          "name": "light_turn_on",
                          "type": "BOOLEAN"
                        },
                        {
                          "name": "height_adjust",
                          "type": "FLOAT"
                        }
                      ]
                    },
                    {
                      "name": "stop_job",
                      "type": "TOPIC",
                      "fields": [
                        {
                          "name": "water_pump_off",
                          "type": "BOOLEAN"
                        },
                        {
                          "name": "window_close",
                          "type": "BOOLEAN"
                        }
                      ]
                    }
                  ]
                },
                {
                  "name": "Metric",
                  "type": "PATH",
                  "topicType": "METRIC",
                  "children": [
                    {
                      "name": "board_cycle_time",
                      "type": "TOPIC",
                      "enableHistory": "TRUE",
                      "mockData": "FALSE",
                      "templateAlias": "basic_order",
                      "labels": "alpha",
                      "extendProperties": {
                        "ext1": "value1"
                      }
                    },
                    {
                      "name": "boards_metrics",
                      "type": "TOPIC",
                      "enableHistory": "TRUE",
                      "mockData": "FALSE",
                      "fields": [
                        {
                          "name": "water",
                          "type": "LONG"
                        },
                        {
                          "name": "temp",
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
`;
