
const q = require('./lib/queue');
const modbus = require('./modbus-bridge');

class CustomProtocolBridge {

    queue;

    timer;

    constructor(node, mappings, interval) {
        this.queue = q.newQueue();
        // 往平台推送数据的topic
        let plantTopic = `4174348a-9222-4e81-b33e-5d72d2fd7f1e`

        this.timer = setInterval(() => {
            let payloads = aggregate(this.queue, mappings);
            if (payloads.length > 0) {
                let newMsg = {
                    "topic": plantTopic,
                    "payload": payloads
                }
                node.send([newMsg])
            }
        }, interval); // 100毫秒推送一次

        node.log(`自定义协议(${node.protocol}) ${interval}ms定时推送任务开启, 节点ID=${node.id}`)
    }

    destroy(nodeId) {
        clearInterval(this.timer);
        this.queue = null;
        console.log(new Date(), `: 自定义协议节点被销毁, 节点ID=${nodeId}`)
    }

    refreshMappings(deviceName, newMapping) {
        // this.mappings = newMapping;
        console.log(`===> 自定义协议 device=${deviceName}: 本地缓存刷新成功`)
    }

    receive(inputMsg) {
        if (this.queue) {
            var msg = {
                "timestamp": Date.now(),
                "data": inputMsg.payload
            }
            this.queue.offer(msg);
        }
    }
}

function aggregate(queue, mappings) {
    let aggregations = [];
    // 每次聚合最多2000个元素
    for (let i = 0; i < 2000; i++) {
        let msg = queue.poll();
        
        if (msg == null) {
            break;
        }
        let data = transfer(msg, mappings);
        if (data.length > 0) {
            aggregations.push(...data);
        }
    }
    return aggregations;
}


/**
 * opcda数据转换
 * 
 * msg数据格式为数组
 */
function transfer(msg, mappings) {

    if (Array.isArray(msg.data)) {
        return modbus.transfer(msg, mappings);
    } 

    if (isPrimitive(msg.data)) {
        msg.data = [msg.data];
        return modbus.transfer(msg, mappings);
    }

    let timestamp = msg.timestamp;
    let values = [];
    let payloadJsonObj = null;

    try {
        payloadJsonObj = JSON.parse(msg.data);
    } catch (e) {
        console.log("msg.data is not json string, which is ", msg.data);
        // ignore
        return [];
    }

    for (let key in payloadJsonObj) {
        let props = mappings[key]; // key 为位号名  props为数组结构 [{name:"", alias:""}]
        if (!props) {
            console.log(`自定义协议位号(${key})映射uns字段不存在`);
            continue;
        }
        
        for (let i in props) {
            let value = {
                "alias": props[i].alias,
                "data": {
                    "quality": 0,
                    "timeStamp": timestamp
                }
            }
            let propName = props[i].propName;
            value.data[propName] = payloadJsonObj[key];

            values.push(value);
        }
    }
 
    return values;
}

// 判断是否为基础数据类型
function isPrimitive(value) {
    return (
      typeof value === 'string' ||
      typeof value === 'number' ||
      typeof value === 'boolean' ||
      typeof value === 'symbol' ||
      value === null
    );
}


function newCustomProtocolBridge(node, mappings, interval) {
    return new CustomProtocolBridge(node, mappings, interval);
}

module.exports = { newCustomProtocolBridge }; // 导出函数