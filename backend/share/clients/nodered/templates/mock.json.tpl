[
    {
        "id": "$id_inject",
        "type": "inject",
        "z": "",
        "name": "",
        "props": [
            {
                "p": "payload"
            }
        ],
        "repeat": "",
        "crontab": "",
        "once": false,
        "onceDelay": 0.1,
        "topic": "",
        "payload": "",
        "payloadType": "date",
        "x": 220,
        "y": 420,
        "wires": [
            [
                "$id_model_selector"
            ]
        ]
    },
    {
        "id": "$id_model_selector",
        "type": "supmodel",
        "z": "",
        "name": "",
        "protocol": "mock",
        "selectedModel": "$uns_path",
        "selectedModelAlias": "$model_alias",
        "modelShowName": "",
        "tableValid": true,
        "x": 430,
        "y": 420,
        "wires": [
            [
                "$id_mqtt"
            ],
            []
        ]
    },
    {
        "id": "$id_mqtt",
        "type": "mqtt out",
        "z": "",
        "name": "",
        "topic": "",
        "qos": "",
        "retain": "",
        "respTopic": "",
        "contentType": "",
        "userProps": "",
        "correl": "",
        "expiry": "",
        "broker": "85bb67b2dbefe3ba",
        "x": 670,
        "y": 420,
        "wires": []
    },
    {
        "id": "85bb67b2dbefe3ba",
        "type": "mqtt-broker",
        "name": "emqx:1883",
        "broker": "emqx",
        "port": 1883,
        "clientid": "",
        "autoConnect": true,
        "usetls": false,
        "protocolVersion": "4",
        "keepalive": "60",
        "cleansession": true,
        "autoUnsubscribe": true,
        "birthTopic": "",
        "birthQos": "0",
        "birthRetain": "false",
        "birthPayload": "",
        "birthMsg": {},
        "closeTopic": "",
        "closeQos": "0",
        "closeRetain": "false",
        "closePayload": "",
        "closeMsg": {},
        "willTopic": "",
        "willQos": "0",
        "willRetain": "false",
        "willPayload": "",
        "willMsg": {},
        "userProps": "",
        "sessionExpiry": ""
    }
]