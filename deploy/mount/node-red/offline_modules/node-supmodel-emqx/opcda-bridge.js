
const q = require('./lib/queue');

class OpcdaBridge {

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
        
        node.log(`OPCDA ${interval}ms定时推送任务开启, 节点ID=${node.id}`)
    }

    destroy(nodeId) {
        clearInterval(this.timer);
        this.queue = null;
        this.mappings = null;
        console.log(new Date(), `: OPCDA节点被销毁, 节点ID=${nodeId}`)
    }

    refreshMappings(deviceName, newMapping) {
        // this.mappings = newMapping;
        console.log(`===> OPCDA device=${deviceName}: 本地缓存刷新成功`)
    }

    receive(inputMsg, envs) {
        if (!this.queue) {
            return;
        }
        let topicMap = transfer(inputMsg.payload, this.mappings, envs);
        for (let key in topicMap) {
            this.queue.offer({
                topic: key,
                payload: topicMap[key]
            });
        }
    }
}

/**
 * opcda数据转换
 * 
 * msg数据格式为数组
 */
function transfer(dataArray, mappings, envs) {
    let topicResult = {};

    if (!mappings == null) {
        return topicResult;
    }

    let _t = envs.field_t_var;
    let _q = envs.field_q_var;
    let useAlias = envs.use_alias;
    
    for (let i in dataArray) {
        let tag = dataArray[i].itemID; // opcda的位号
        let fields = mappings[tag]; // uns的位号， 可以是多个
        if (!fields) {
            console.log(`opcda ${tag} 映射uns字段不存在`);
            continue;
        }
        for (let ii in fields) {
            let vqt = {};
            vqt[_t] = new Date(dataArray[i].timestamp).getTime();
            vqt[_q] = dataArray[i].errorCode;
            let propName = fields[ii].propName;
            vqt[propName] = dataArray[i].value;

            if (useAlias === true) {
                let values = topicResult[fields[ii].alias] || [];
                values.push(vqt);
                topicResult[fields[ii].alias] = values;
            } else {
                let values = topicResult[fields[ii].path] || [];
                values.push(vqt);
                topicResult[fields[ii].path] = values;
            }
        }
    }
    return topicResult;
}

function newOpcdaBridge(node, mappings, interval) {
    return new OpcdaBridge(node, mappings, interval);
}

module.exports = { newOpcdaBridge }; // 导出函数