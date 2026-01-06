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
        "repeat": "10",
        "crontab": "",
        "once": false,
        "onceDelay": 1,
        "topic": "",
        "payload": "",
        "payloadType": "date",
        "x": 320,
        "y": 160,
        "wires": [
            [
                "$id_func"
            ]
        ]
    },
    {
        "id": "$id_func",
        "type": "function",
        "z": "",
        "name": "mock data",
        "func": "// 随机字符串\nfunction randomString() {\n    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';\n    let result = '';\n    for (let i = 0; i < 20; i++) {\n        result += chars.charAt(Math.floor(Math.random() * chars.length));\n    }\n    return result;\n}\n// 100以内随机数字\nfunction generateRandomNumber() {\n    return Math.floor(Math.random() * 100);\n}\n// 随机生成100以内浮点数，保留2位小数\nfunction generateRandomFloatWithTwoDecimals() {\n    const randomFloat = Math.random() * 100;\n    return randomFloat.toFixed(2);\n}\n// 对当前时间格式化\nfunction formatCurDate() {return Date.now();\n}\n\nfunction getBool() {\n    var randomInt = generateRandomNumber();\n    return randomInt > 50;\n}\nmsg.topic='$alias_path_topic';\nmsg.payload = {$payload \n};\n\nreturn msg;",
        "outputs": 1,
        "timeout": 0,
        "noerr": 0,
        "initialize": "",
        "finalize": "",
        "libs": [],
        "x": 520,
        "y": 160,
        "wires": [
            [
                "$id_model_selector"
            ]
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
        "x": 990,
        "y": 160,
        "wires": []
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
        "x": 750,
        "y": 160,
        "wires": [
            [
                "$id_mqtt"
            ],
            []
        ]
    },
    {
        "id": "85bb67b2dbefe3ba",
        "type": "mqtt-broker",
        "name": "",
        "broker": "emqx",
        "port": "1883",
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