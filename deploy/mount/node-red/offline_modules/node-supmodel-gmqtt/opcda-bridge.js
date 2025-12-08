
const q = require('./lib/queue');

class OpcdaBridge {

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
        }, interval); // 100毫秒推送一次

        node.log(`OPCDA ${interval}ms定时推送任务开启, 节点ID=${node.id}`)
    }

    destroy(nodeId) {
        clearInterval(this.timer);
        this.queue = null;
        console.log(new Date(), `: OPCDA节点被销毁, 节点ID=${nodeId}`)
    }

    refreshMappings(deviceName, newMapping) {
        // this.mappings = newMapping;
        console.log(`===> OPCDA device=${deviceName}: 本地缓存刷新成功`)
    }

    receive(inputMsg) {
        if (this.queue) {
            this.queue.offer(inputMsg.payload);
        }
    }
}

function aggregate(queue, mapping, _t, _q) {
    let aggregations = [];
    // 每次聚合最多2000个元素
    for (let i = 0; i < 2000; i++) {
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
 * opcda数据转换
 * 
 * msg数据格式为数组
 */
function transfer(msg, mappings, _t, _q) {
    let values = [];
    
    for (let i in msg) {
        let tag = msg[i].itemID; // opcda的位号
        let fields = mappings[tag]; // uns的位号， 可以是多个
        if (!fields) {
            console.log(`opcda ${tag} 映射uns字段不存在`);
            continue;
        }
        for (let j in fields) {
            let value = {
                "alias": fields[j].alias,
                "data": {
                }
            }
            value.data[_t] = new Date(msg[i].timestamp).getTime();
            value.data[_q] = msg[i].errorCode;
            let propName = fields[j].propName;
            value.data[propName] = msg[i].value;

            values.push(value);
        }
    }
    return values;
}

function newOpcdaBridge(node, mappings, interval) {
    return new OpcdaBridge(node, mappings, interval);
}

module.exports = { newOpcdaBridge }; // 导出函数