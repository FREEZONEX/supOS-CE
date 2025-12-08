
const q = require('./lib/queue');


class ModbusBridge {

    queue;

    timer;

    constructor(node, mappings, interval) {
        this.queue = q.newQueue();
        // 往平台推送数据的topic
        let plantTopic = `4174348a-9222-4e81-b33e-5d72d2fd7f1e`

        let _t = node.envs.field_t_var;
        let _q = node.envs.field_q_var;
        this.timer = setInterval(() => {
            let payloads = aggregate(this.queue, mappings, _t, _q);
            if (payloads.length > 0) {
                let newMsg = {
                    "topic": plantTopic,
                    "payload": payloads
                }
                node.send([newMsg])
            }
        }, interval); 

        node.log(`MODBUS ${interval}ms定时推送任务开启, 节点ID=${node.id}`)
    }

    refreshMappings(deviceName, newMapping) {
        // this.mappings = newMapping;
        console.log(`===> MODBUS device=${deviceName}: 本地缓存刷新成功`)
    }


    receive(inputMsg) {
        if (!Array.isArray(inputMsg.payload))  {
            return "supmodel.error_modbus_data_not_array";
        }
        var msg = {
            "timestamp": Date.now(),
            "address": inputMsg.input?.payload?.address || 0,
            "quantity": inputMsg.input?.payload?.quantity || 0,
            "data": inputMsg.payload
        }
        if (this.queue) {
            this.queue.offer(msg);
        }
    }

    destroy(nodeId) {
        clearInterval(this.timer);
        this.queue = null;
        console.log(new Date(), `: MODBUS节点被销毁, 节点ID=${nodeId}`);
    }
}

function aggregate(queue, mapping, _t, _q) {
    let aggregations = [];
    // 每次聚合最多1000个元素, 每组数据最大255b
    for (let i = 0; i < 1000; i++) {
        let msg = queue.poll();
        if (msg == null) {
            break;
        }
        let data = transfer(msg, mapping, _t, _q);
        if (data.length > 0) {
            aggregations.push(...data);
        }
    }
    return aggregations;
}


/**
 * modbus数据转换
 */
function transfer(msg, mappings, _t, _q) {

    let values = [];
    let timestamp = msg.timestamp;

    let address = msg.address || 0; // 数组下标起始位置

    for (let index in msg.data) {
        let realIndex = plus(index, address);
        var tags = mappings[realIndex];
        if (tags) {
            for (let i in tags) {
                let propName = tags[i].propName;
                let value = {
                    "alias": tags[i].alias,
                    "data": {
                    }
                }
                value.data[_t] = timestamp;
                value.data[_q] = 0;
                value.data[propName] = msg.data[index];
                values.push(value)
            }
        }
        
    }
    return values;
}

function plus(a, b) {
    return Number(a) + Number(b);
}

//1: int; 2: double; 3: bool; 4: string
function unwrapModbusVariant(val) {
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




function newModbusBridge(node, mappings, interval) {
    return new ModbusBridge(node, mappings, interval);
}

module.exports = { newModbusBridge, transfer }; // 导出函数