
const opcua = require('./opcua-bridge');
const opcda = require('./opcda-bridge');
const modbus = require('./modbus-bridge');
const custom = require('./custom-bridge');
const mqtt = require('./mqtt-bridge');
const mock = require('./mock-bridge');
const fs = require("fs");
const formidable = require("formidable");

module.exports = function (RED) {

    const httpNode = RED.httpNode;

    const storage = '/data/cache/context/global/';
    const storagePath = storage + 'global.json';

    const templateJSON = '/data/template/template.json';

    // 下载模板文件
    httpNode.get('/nodered-api/download/template', (req, res) => {
        res.download(templateJSON);
    });

    httpNode.get('/nodered-api/load/tags', (req, res) => {
        const nodeId = req.query.nodeId; // 查询条件
        // tag array
        let tags = loadStorage(nodeId);
        res.status(200).json({data: tags});

    });

    httpNode.get('/nodered-api/export/tags', (req, res) => {
        const nodeId = req.query.nodeId;
        const tags = loadStorage(nodeId) || [];
        const mapping = toMapping(tags);
        const exportJsonName = `UNS-Mapper-${nodeId || 'unknown'}.json`;
        res.setHeader('Content-Disposition', `attachment; filename="${exportJsonName}"`);
        res.setHeader('Content-Type', 'application/json');
        res.send(JSON.stringify(mapping, null, 2));
    });

    // 导入/上传 JSON
    httpNode.post('/nodered-api/upload/tags', (req, res) => {
        const contentType = (req.headers['content-type'] || '').toLowerCase();
        const nodeId = (req.query && req.query.nodeId) || (req.body && req.body.nodeId);

        const respond = (rows) => {
            if (!rows || rows.length === 0) {
                res.status(400).end('mapping is empty');
                return;
            }
            if (nodeId) {
                saveGlobalStorage(nodeId, rows);
            }
            res.status(200).json({ data: rows });
        };

        if (contentType.indexOf('application/json') !== -1) {
            try {
                const rows = normalizeMappingPayload(req.body);
                respond(rows);
            } catch (err) {
                console.log(err);
                res.status(400).end("Invalid JSON payload");
            }
            return;
        }

        const form = new formidable.IncomingForm();
        form.uploadDir = "/data/uploads"; // 设置上传目录
        form.keepExtensions = true; // 保留文件扩展名

        form.parse(req, (err, fields, files) => {
            if (err) {
                console.log(err);
                res.status(500).end("上传失败");
                return;
            }
            res.setHeader('Content-Type', 'application/json');
            const uploadFile = files.file && files.file[0];
            const jsonFilePath = uploadFile && uploadFile.filepath;
            
            if (jsonFilePath && fs.existsSync(jsonFilePath)) {
                try {
                    const buffer = fs.readFileSync(jsonFilePath);
                    const jsonData = JSON.parse(buffer.toString('utf-8'));
                    const rows = normalizeMappingPayload(jsonData);
                    fs.unlinkSync(jsonFilePath);
                    respond(rows);
                } catch (error) {
                    console.log("解析json异常", error);
                    res.status(500).end("JSON parse failed, please check format!");
                }
            } else {
                RED.log.error('Please confirm whether the "/data/uploads" directory exists in container.');
                res.status(500).end("JSON upload failed!");
            }
        });
        
    });

    httpNode.get('/nodered-api/tags', (req, res) => {
        const nodeIds = req.query.nodeId; // 查询条件
        let idsArray = [];
        let data = [];
        if (Array.isArray(nodeIds)) {
            idsArray = nodeIds;
        } else if (nodeIds) {
            idsArray = [nodeIds];
        } else {
            res.status(200).json(data);
            return;
        }
        for (let key in idsArray) {
            let nodeId = idsArray[key];
            let tags = loadStorage(nodeId);
            if (tags && tags.length > 0) {
                let ss = {};
                ss[nodeId] = tags;
                data.push(ss);
            }
        }
        res.status(200).json(data);
    });

    httpNode.post('/nodered-api/batchSave/tags', (req, res) => {
        const itemsArray = req.body;
        itemsArray.forEach((item, index) => {
            Object.entries(item).forEach(([nodeId, tags]) => {
                if (tags && tags.length > 0) {
                    saveGlobalStorage(nodeId, tags);
                }
            });
        });
        res.status(200).end("success");
    });

    httpNode.post('/nodered-api/save/tags', (req, res) => {
        
        let success = saveGlobalStorage(req.body.nodeId, req.body.tags);
        if (success === true) {
            res.status(200).end("success");
        } else {
            res.status(500).json({msg: "Tag save failed"});
        }
    });

    httpNode.get('/nodered-api/query/tags', (req, res) => {

        let pageNo = req.query.pageNo || 1; // 起始页
        let nodeId = req.query.nodeId; // 查询条件

        const mappingData = loadStorage(nodeId);

        res.status(200).json({data: mappingData});
    });

    function saveGlobalStorage(key, data) {
        try {
            if (!fs.existsSync(storagePath)) {
                fs.mkdirSync(storage, { recursive: true });
                fs.writeFileSync(storagePath, JSON.stringify({key: []}));
            }
            let buffer = fs.readFileSync(storagePath);
            if (buffer[0] === 0xEF && buffer[1] === 0xBB && buffer[2] === 0xBF) {
                buffer = buffer.subarray(3);
            }
            let content = buffer.toString('utf-8');
            if (!content) {
                content = "{}";
            }
            let o = JSON.parse(content);
            o[key] = data || [];
            fs.writeFileSync(storagePath, JSON.stringify(o));
            return true;
        } catch (err) {
            console.log(err)
            RED.log.error(`Tag save failed`);
            return false;
        }
        
    }

    function loadStorage(key) {
        try {
            if (!fs.existsSync(storagePath)) {
                fs.mkdirSync(storage, { recursive: true });
                fs.writeFileSync(storagePath, JSON.stringify({key: []}));
                return [];
            }

            let buffer = fs.readFileSync(storagePath);
            if (buffer[0] === 0xEF && buffer[1] === 0xBB && buffer[2] === 0xBF) {
                buffer = buffer.subarray(3);
            }
            let content = buffer.toString('utf-8');
            if (!content) {
                content = "{}";
            }
            let o = JSON.parse(content);
            return o[key] || [];
        } catch (err) {
            console.log(err);
            return [];
        }
    }

    function buildMappings(mappingData) {
        let mappings = {};
        if (!mappingData || mappingData.length == 0) {
            return mappings;
        }
        
        mappingData.map(row => {
            // path, alias, propName, propType, tag
            let path, alias, propName, propType, tag;
            if (Array.isArray(row)) {
                [path, alias, propName, propType, tag] = row;
            } else {
                path = row?.topic || row?.Topic || row?.targetTopic || row?.path;
                alias = row?.alias || row?.targetAlias;
                propName = row?.attributeName || row?.AttributeName || row?.targetField || row?.propName;
                propType = row?.attributeType || row?.targetType || row?.propType || row?.type;
                tag = row?.tagConfiguration || row?.TagConfiguration || row?.selector || row?.tag;
            }
            if (tag === undefined || tag === null) {
                return;
            }
            let props = mappings[tag] || [];

            props.push({
                path: path,
                alias: alias,
                propName: propName,
                propType: propType
            });
            mappings[tag] = props;
        });
        return mappings;
    }
    
    function normalizeMappingPayload(payload) {
        if (payload == null) {
            return [];
        }
        if (Array.isArray(payload)) {
            return payload.map(normalizeSingleMapping).filter(Boolean);
        }
        if (Array.isArray(payload.mapping)) {
            return payload.mapping.map(normalizeSingleMapping).filter(Boolean);
        }
        if (Array.isArray(payload.tags)) {
            return payload.tags.map(normalizeSingleMapping).filter(Boolean);
        }
        if (Array.isArray(payload.data)) {
            return payload.data.map(normalizeSingleMapping).filter(Boolean);
        }
        return [];
    }

    function normalizeSingleMapping(item) {
        if (!item) {
            return null;
        }
        if (Array.isArray(item)) {
            return item;
        }
        const selector = toSafeString(item.tagConfiguration || item.TagConfiguration || item.selector || item.tag || item.target || item.topicIndex);
        const targetTopic = toSafeString(item.topic || item.Topic || item.targetTopic || item.path || item.topic);
        const targetField = toSafeString(item.attributeName || item.AttributeName || item.targetField || item.propName || item.field || item.property);
        const alias = toSafeString(item.alias || item.targetAlias);
        const propType = toSafeString(item.attributeType || item.targetType || item.propType || item.type);
        if (!selector && !targetTopic && !targetField) {
            return null;
        }
        return [targetTopic, alias, targetField, propType, selector];
    }

    function toMapping(rows) {
        if (!Array.isArray(rows)) {
            return [];
        }
        return rows.map(row => {
            if (Array.isArray(row)) {
                return {
                    tagConfiguration: row[4] !== undefined ? String(row[4]) : "",
                    topic: row[0] || "",
                    attributeName: row[2] || "",
                    attributeType: row[3] || ""
                };
            }
            if (row && typeof row === 'object') {
                return {
                    tagConfiguration: toSafeString(row.tagConfiguration || row.TagConfiguration || row.selector || row.tag),
                    topic: toSafeString(row.topic || row.Topic || row.targetTopic || row.path),
                    attributeName: toSafeString(row.attributeName || row.AttributeName || row.targetField || row.propName),
                    attributeType: toSafeString(row.attributeType || row.targetType || row.propType || row.type)
                };
            }
        }).filter(item => item && (item.tagConfiguration || item.topic || item.attributeName));
    }

    function toSafeString(value) {
        if (value === undefined || value === null) {
            return "";
        }
        return String(value).trim();
    }

    function SelectModel(config) {

        RED.nodes.createNode(this, config);
        const node = this;

        if (!fs.existsSync(storagePath)) {
            node.context().global.set(node.id, []);
        }

        node.protocol = config.protocol;
        node.selectedModel = config.selectedModel;
        node.selectedModelAlias = config.selectedModelAlias;
        let models = loadStorage(node.id);
        if (!models || models.length === 0) {
            models = normalizeMappingPayload(config.mapping);
            if (models && models.length > 0) {
                saveGlobalStorage(node.id, models);
            }
        }
        node.models = models;
        node.envs = {
            "field_t_var": process.env.TIMESTAMP_NAME,
            "field_q_var": process.env.QUALITY_NAME,
            "use_alias": process.env.USE_ALIAS_AS_TOPIC
        }
        node.mappings = buildMappings(node.models);

        node.bridge = null;
        node.interval = config.interval || 100; // 推送频率 默认100ms

        switch (node.protocol) {
            case "opcua": node.bridge = opcua.newOpcuaBridge(node, node.mappings, 5); break;
            case "opcda": node.bridge = opcda.newOpcdaBridge(node, node.mappings, 5); break;
            case "modbus": node.bridge = modbus.newModbusBridge(node, node.mappings, 5); break;
            case "mqtt": node.bridge = mqtt.newMqttBridge(node, node.mappings, 5); break;
            case "mock": node.bridge = mock.newMockDataBridge(node, 1000); break;
            case "custom": node.bridge = custom.newCustomProtocolBridge(node, node.mappings, 5); break;
            default: {
                node.error('节点未生效：请选择协议');
                return;
            }
        }

        this.on("close", function (done) {
            if (node.bridge) {
                node.bridge.destroy(node.id);
                node.bridge = null;
            }
            done(); // 必须调用以完成清理
        });

        node.on('input', function (msg) {
            var errorMsg = node.bridge.receive(msg, node.envs);
            if (errorMsg) {
                var errMsg = RED._(errorMsg);
                node.error(errMsg, msg);
            }
        });


    }
    RED.nodes.registerType("supmodel", SelectModel);
    
    
}
