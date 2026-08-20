export const namespaceExample = `[
  {
    "name": "v1",
    "children": [
      {
        "name": "Plant_Name",
        "children": [
          {
            "name": "SMT-Area-1",
            "children": [
              {
                "name": "State",
                "children": [
                  {
                    "name": "current_job",
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
                "children": [
                  {
                    "name": "start_job",
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
                "children": [
                  {
                    "name": "board_cycle_time",
                    "enableHistory": "TRUE",
                    "mockData": "FALSE",
                    "extendProperties": {
                      "ext1": "value1"
                    },
                    "fields": [
                      {
                        "name": "cycle_time",
                        "type": "DOUBLE"
                      }
                    ]
                  },
                  {
                    "name": "boards_metrics",
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
]`;

export const placeholder = `{
  "namespace": [
${namespaceExample
  .split('\n')
  .slice(1, -1)
  .map((line) => `    ${line}`)
  .join('\n')}
  ]
}
`;

export const template = `{
  "notes": {
    "path": "Path nodes organize the namespace tree and can contain other path nodes or system topic folders.",
    "topic": "Topic nodes are created under State, Action, or Metric folders and store data definitions.",
    "topicFields": {
      "fields": "Optional payload schema for State and Action topics. Metric topics must include at least one field. Each field supports name and type. type supports INTEGER, STRING, FLOAT, DOUBLE, BOOLEAN, LONG, and DATETIME.",
      "enableHistory": "TRUE enables historical data persistence for the topic; FALSE or omitted keeps the default behavior.",
      "mockData": "TRUE creates and binds a mock source flow that publishes sample data; FALSE or omitted does not create mock data.",
      "extendProperties": "Optional custom key-value metadata for integrations."
    }
  },
  "namespace": [
${namespaceExample
  .split('\n')
  .slice(1, -1)
  .map((line) => `    ${line}`)
  .join('\n')}
  ]
}
`;
