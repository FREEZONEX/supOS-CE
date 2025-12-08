

module.exports = function (RED) {

    function FieldMapping(config) {

        RED.nodes.createNode(this, config);
        const node = this;

        node.fsMapping = buildFieldMapping(config.fieldMapping);

        node.selectModelTopic = config.selectedModelTopic;
        node.dataRoot = config.dataRoot;

        node.on('input', function (msg) {
            var inputData = msg.payload;
            let dataArray = getInnerData(inputData, node.dataRoot);
            let unsDataArray = [];
            if (dataArray && dataArray.length > 0) {
                // 先将数据源扁平化
                let flatArray = flatDataArray(node, dataArray);
                // 根据节点的字段映射关系，将数据源转换为uns字段对应的数据
                for (let index in flatArray) {
                    var unsObj = {};
                    var dataObj = flatArray[index];
                    for (let key in dataObj) {
                        var mapValue = node.fsMapping[key]; // 映射关系uns字段名
                        if (mapValue) {
                            unsObj[mapValue] = dataObj[key];
                        }
                    }
                    unsDataArray.push(unsObj);
                }
            }
            msg.topic = node.selectModelTopic;
            msg.payload = unsDataArray;
            node.send(msg);
        });

    }

    function buildFieldMapping(fieldMapping) {
        if (!fieldMapping) {
            return {};
        }
        let fieldMappingCache = {};
        for (let i in fieldMapping) {
            var o = fieldMapping[i];
            let key  = o['keyValue'];
            let value = o['valueValue'];
            fieldMappingCache[key] = value;
        }
        return fieldMappingCache;
    }

    // 将数据扁平化
    function flatDataArray(node, dataArray) {
        let flatArray = [];
        for (let i in dataArray) {
            var sd = dataArray[i];
            let single = {};
            if (isJSON(sd)) {
                // 遍历单个对象的所有属性，并将其扁平化成一个map对象
                for (let k in sd) {
                    var flatArr = recurve(k, sd[k], {});
                    Object.assign(single, flatArr);
                }
                flatArray.push(single);
            } else {
                node.error(`${sd} not a json object, cannot be parsed`);
            }
        }
        return flatArray;
    }

    function recurve(k, vObj, container) {
        if (!container) {
            container = {};
        }
        if (isJSON(vObj)) {
            for (let t in vObj) {
                recurve(k + "." + t, vObj[t], container);
            }
        } else {
            container[k] = vObj;
        }
        return container;
    }

    function isJSON(obj) {
        try {
            var objStr = JSON.stringify(obj); // 尝试解析字符串
            if (/^\s*(\{.*\}|\[.*\])\s*$/.test(objStr)) {
                return true;
            }
        } catch (e) {
        }
        return false; 
    }

    // 获取里层的数据，例如{page: 1, list: {data: []}}, dataRoot为list.data
    function getInnerData(inputData, dataRoot) {
        if (Array.isArray(inputData)) {
            return inputData;
        }
        if (dataRoot) {
            var sp = dataRoot.split(".");
            var innerData = inputData;
            for (let i in sp) {
                innerData = innerData[sp[i]];
            }
            return getInnerData(innerData, "");
        } else {
            return [inputData];
        }

    }

    RED.nodes.registerType("modelConverter", FieldMapping);
    
    
}