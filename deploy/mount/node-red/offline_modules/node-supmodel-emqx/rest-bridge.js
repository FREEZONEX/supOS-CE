
const q = require('./lib/queue');

class RestBridge {

    queue;

    timer;

    constructor(node, model, interval) {
        this.queue = q.newQueue();
        // 往平台推送数据的topic
        let plantTopic = `4174348a-9222-4e81-b33e-5d72d2fd7f1e`

        this.timer = setInterval(() => {
            let payloads = aggregate(node.queue, model);
            if (payloads.length > 0) {
                let newMsg = {
                    "topic": plantTopic,
                    "payload": payloads
                }
                node.send([newMsg])
            }
        }, interval); 

        console.log(new Date(), `: REST定时推送任务开启, 节点ID=${node.id},  设备名称=${model.alias}`)
    }

    refreshMappings(deviceName, newMapping) {
        
    }

    receive(inputMsg) {
        if (this.queue) {
            this.queue.offer(inputMsg.payload);
        }
    }

    destroy(nodeId, deviceName) {
        clearInterval(this.timer);
        this.queue = null;
        console.log(new Date(), `: REST节点被销毁, 节点ID=${nodeId}, 设备名称=${deviceName}`)
    }
}

function aggregate(queue, model) {
    let aggregations = [];
    // 每次聚合最多1000个元素, 每组数据最大255b
    for (let i = 0; i < 100; i++) {
        let payload = queue.poll();
        if (payload == null) {
            break;
        }
        let data = transfer(payload, model);

        aggregations.push(...data);
    }
    return aggregations;
}


function transfer(payload, model) {

    let values = [];
    let value = {
        "alias": model.alias,
        "data":  payload
    }
    values.push(value);
        
    return values;
}

function newRestBridge(node, model, interval) {
    return new RestBridge(node, model, interval);
}

module.exports = { newRestBridge }; // 导出函数