
const fs = require('fs');
const q = require('./lib/queue');

class OpcuaBridge {

    queue;

    timer;

    mappings;

    topologyTimer;

    latestStatus; // 保留最近一次opcua-server的状态值
    topologyFile = "/data/log/topology.log";

    constructor(node, mappings, interval) {
        this.queue = q.newQueue();
        this.mappings = mappings;

        this.timer = setInterval(() => {
            let newMsg = null;
            while ((newMsg = this.queue.poll()) != null) {
                node.send([newMsg])
            }
        }, interval); // 轮询队列里的数据

        // this.topologyTimer = setInterval(() => {
        //     fs.unlink(this.topologyFile, (err) => {
        //         if (err) {
        //             console.error('删除文件失败:', this.topologyFile, err);
        //         } else {
        //             console.log('文件已删除:', this.topologyFile);
        //         }
        //     });
        // }, 8 * 60 * 60 * 1000); // 8小时删除一次文件
        node.log(`OPCUA ${interval}ms定时推送任务开启, 节点ID=${node.id}`)
    }

    destroy(nodeId) {
        this.queue = null;
        this.mappings = null;
        clearInterval(this.timer);
        clearInterval(this.topologyTimer);
        console.log(new Date(), `: OPCUA节点被销毁, 节点ID=${nodeId}`)
    }

    refreshMappings(deviceName, newMapping) {
        // this.mappings = newMapping;
        console.log(`===> OPCUA device=${deviceName}: 本地缓存刷新成功`)
    }

    receive(inputMsg, envs) {
        if (!this.queue || inputMsg.topic == 'subscriptionId') {
            return;
        }
        let topicMap = transfer(inputMsg, this.mappings, envs);
        for (let key in topicMap) {
            this.queue.offer({
                topic: key,
                payload: topicMap[key]
            });
        }
    }

}

function reportStatusToFile(mappings, status, filePath) {
    let current = Date.now();
    let errorJsons = "";

    if (!mappings) {
        return;
    }

    let statusObj = getErrorByStatus(status);
    for (let key in mappings) {
        for (let i in mappings[key]) {
            if (statusObj.code !== 200) {
                let errJson = `{"instanceTopic":"${mappings[key][i].alias}","topologyNode":"pushOriginalData","eventCode":"1","eventMessage":"${statusObj.message}","eventTime":${current}},\n`
                errorJsons += errJson;
            } else {
                let errJson = `{"instanceTopic":"${mappings[key][i].alias}","topologyNode":"pushOriginalData","eventCode":"0","eventMessage":"","eventTime":${current}},\n`
                errorJsons += errJson;
            }
        }
    }

    fs.appendFile(filePath, errorJsons, (error) => {
        if (error) {
            console.log(error);
        }
    });
}

function getErrorByStatus(status) {
    let aliveStatus = ['keepalive', 'connected', 'subscribed'];
    if (status.indexOf('active') !== -1 || status.indexOf('connected') !== -1 || aliveStatus.includes(status)) {
        return {
            "code": 200,
            "message": 'opcua server is connected.'
        }
    } else if (status == 'create client' || status == 'connecting' ) {
        return {
            "code": 400,
            "message": 'opcua client is connecting...'
        }
    } else {
        return {
            "code": 400,
            "message": 'opcua server connect failed, status: ' + status
        }
    }
}


/**
 * opcua数据转换
 */
function transfer(msg, mappings, envs) {
    // opcua input分为批量和单个，数据格式有所不同
    let isBatch = (msg.payload.value?.dataType != null && msg.payload.value?.dataType != undefined) ? true : false;
    let topicResult = {};

    if (mappings == null || mappings == undefined) {
        return topicResult;
    }
    let fields = mappings[msg.topic];

    let _t = envs.field_t_var;
    let _q = envs.field_q_var;
    let useAlias = envs.use_alias;

    if (fields) {
        if (isBatch) {
            //批量订阅的数据结构
            /**
             * {
                topic: 'ns=6;s=DataItem_0962',
                payload: DataValue {
                    statusCode: ConstantStatusCode {
                        _value: 0,
                        _description: 'The operation succeeded.',
                        _name: 'Good'
                    },
                    sourceTimestamp: 2025-03-27T01:13:32.000Z,
                    sourcePicoseconds: 0,
                    serverTimestamp: 2025-03-27T01:13:32.000Z,
                    serverPicoseconds: 0,
                    value: Variant {
                        dataType: 11,
                        arrayType: 0,
                        value: 54.46349034227433,
                        dimensions: null
                    }
                },
                _msgid: 'a2ced847173eeec8'
                }
            */
            // 组装平台数据格式
            for (let i in fields) {
                let propName = fields[i].propName;
                let vqt = {};
                vqt[_t] = new Date(msg.payload.sourceTimestamp).getTime();
                vqt[_q] = msg.payload.statusCode._value;
                vqt[propName] = msg.payload.value.value;

                if (useAlias === true) {
                    let values = topicResult[fields[i].alias] || [];
                    values.push(vqt);
                    topicResult[fields[i].alias] = values;
                } else {
                    let values = topicResult[fields[i].path] || [];
                    values.push(vqt);
                    topicResult[fields[i].path] = values;
                }
            }
        } else {
            // 单个数据结构
            /**
             * {
                    "_msgid": "5a137f3ea9ec78f9",
                    "payload": 71.93397996788191,
                    "topic": "ns=6;s=DataItem_0393",
                    "datatype": "Double",
                    "browseName": "",
                    "statusCode": {
                        "value": 0
                    },
                    "serverTimestamp": "2025-03-08T02:14:56.561Z",
                    "sourceTimestamp": "2025-03-08T02:14:56.000Z"
                }
            */
            // 组装平台数据格式
            for (let i in fields) {
                let propName = fields[i].propName;
                let vqt = {};
                vqt[_t] = new Date(msg.sourceTimestamp).getTime();
                vqt[_q] = msg.statusCode.value;
                vqt[propName] = msg.payload;
                
                if (useAlias === true) {
                    let values = topicResult[fields[i].alias] || [];
                    values.push(vqt);
                    topicResult[fields[i].alias] = values;
                } else {
                    let values = topicResult[fields[i].path] || [];
                    values.push(vqt);
                    topicResult[fields[i].path] = values;
                }
            }
        }
    }
    return topicResult;
}

//1: int; 2: double; 3: bool; 4: string
function unwrapOpcUaVariant(opcuaDataType, value) {
    switch (opcuaDataType) {
        case "4": 
        case "6": 
        case "8": 
        case "27": 
        case "5": 
        case "7": 
        case "9": 
        case "28": 
            return 1;
        case "10": 
        case "11": 
        case "26":
            return 2;
        case "1": 
            return 3;
        case "3": 
        case "15": 
        case "2": 
        case "12": 
        case "21":
            return 4;
        default: return unwrapUnknownVariant(value);
    }
}

function unwrapUnknownVariant(val) {
    var type = typeof(val);
    switch (type) {
        case "number": 
            if (val.toString().indexOf(".") !== -1) { // 判断是否为浮点数
                return 2; // 浮点
            } else {
                return 1; // 整数
            }
        case "boolean": 
            return 3;
        case "string": 
            return 4;
        
        default: return 0;
    }
}

function newOpcuaBridge(node, mappings, interval) {
    return new OpcuaBridge(node, mappings, interval);
}

module.exports = { newOpcuaBridge }; // 导出函数