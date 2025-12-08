
const q = require('./lib/queue');


class ModbusBridge {

    queue;

    timer;

    mappings;

    constructor(node, mappings, interval) {
        this.mappings = mappings;
        this.queue = q.newQueue();

        this.timer = setInterval(() => {
            let newMsg = this.queue.poll();
            if (newMsg != null) {
                node.send([newMsg])
            }
        }, interval); // 5毫秒轮询队列里的数据，每次只取一个topic发送
        node.log(`MODBUS ${interval}ms定时推送任务开启, 节点ID=${node.id}`)
    }

    refreshMappings(deviceName, newMapping) {
        // this.mappings = newMapping;
        console.log(`===> MODBUS device=${deviceName}: 本地缓存刷新成功`)
    }

    receive(inputMsg, envs) {
        if (!this.queue) {
            return;
        }
        if (!Array.isArray(inputMsg.payload))  {
            return "supmodel.error_modbus_data_not_array";
        }
        const address = inputMsg.input?.payload?.address || 0;
        let topicMap = transfer(inputMsg.payload, address, this.mappings, envs);
        for (let key in topicMap) {
            this.queue.offer({
                topic: key,
                payload: topicMap[key]
            });
        }
    }

    destroy(nodeId) {
        clearInterval(this.timer);
        this.queue = null;
        this.mappings = null;
        console.log(new Date(), `: MODBUS节点被销毁, 节点ID=${nodeId}`);
    }
}

/**
 * modbus数据转换
 */
function transfer(dataArray, address, mappings, envs) {

    let timestamp = Date.now();

    let topicResult = {};

    if (!mappings == null) {
        return topicResult;
    }

    let _t = envs.field_t_var;
    let _q = envs.field_q_var;
    let useAlias = envs.use_alias;

    for (let index in dataArray) {
        let realIndex = plus(index, address);
        var fields = mappings[realIndex];
        if (fields) {
            for (let i in fields) {
                let propName = fields[i].propName;
                let vqt = {};
                vqt[_t] = timestamp;
                vqt[_q] = 0;
                vqt[propName] = dataArray[index];
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