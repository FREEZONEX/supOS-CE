const crypto = require('crypto');

function generateGetSignature(secretKey, url) {
    const queryPart = url.split('?')[1] || '';
    let baseUrl = url.split('?')[0];

    let canonicalRequestString = "GET\n" + baseUrl + "\n" + "application/json\n";
    if (queryPart) {
        var sortedQueryString = buildSortedQueryString(queryPart);
        canonicalRequestString += sortedQueryString + "\n";
    }
    return hexEncodeHMACSHA256(secretKey, canonicalRequestString);
}

function hexEncodeHMACSHA256(secretKey, canonicalRequestString) {
    // 创建 HMAC-SHA256 哈希对象
    const hmac = crypto.createHmac('sha256', secretKey);
    
    // 更新哈希对象的内容
    hmac.update(canonicalRequestString);
    
    // 计算签名并转换为十六进制字符串
    const signature = hmac.digest('hex');
    
    return signature;
}

function buildSortedQueryString(queryPart) {
    
    // 使用 URLSearchParams 解析参数
    const params = new URLSearchParams(queryPart);
    const paramMap = new Map();
  
    // 转换键为小写并存储值
    for (const [key, value] of params) {
      paramMap.set(key.toLowerCase(), value);
    }
  
    // 按小写键的字母序排序
    const sortedKeys = [...paramMap.keys()].sort();
  
    // 拼接排序后的参数
    return sortedKeys
      .map(k => `${k}=${paramMap.get(k)}`)
      .join('&');
  }

module.exports = { generateGetSignature }; 