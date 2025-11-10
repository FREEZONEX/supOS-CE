import { TransportConfig } from '@/types';

const getParams = (urlString: string) => {
  try {
    const url = new URL(urlString);
    const params = new URLSearchParams(url.search);
    return Object.fromEntries(params);
  } catch (e) {
    console.log(e);
    return {};
  }
};
/**
 * 解析自定义传输协议URL
 * @param urlString 要解析的自定义协议URL字符串
 * @returns 包含传输类型和配置的对象
 */
export function parseTransportUrl(urlString: string): TransportConfig {
  // 使用URL类解析整个字符串[1,7](@ref)
  try {
    const url = new URL(urlString);
    const protocol = url.protocol.replace(':', ''); // 移除冒号，得到协议部分

    switch (protocol) {
      case 'sse': {
        const serverUrl = urlString.replace('sse://', '');
        const params = getParams(serverUrl);
        return {
          transportType: 'sse',
          clientName: params.clientName || 'see',
          serverUrl, // 提取see://后面的完整URL[7](@ref)
        };
      }

      case 'streamable-http': {
        const serverUrl = urlString.replace('streamable-http://', '');
        const params = getParams(serverUrl);
        return {
          transportType: 'streamable-http',
          clientName: params.clientName || 'streamable-http',
          serverUrl, // 提取streamable-http://后面的完整URL[7](@ref)
        };
      }

      case 'stdio': {
        // 提取命令（主机名部分）
        const command = url.hostname;

        // 提取参数（路径名部分，需要特殊处理包含'/'的参数）
        const pathParts = url.pathname.split('/').filter((part) => part !== '');
        const args: string[] = [];

        // 智能合并参数，处理包含'/'的参数（如包名）
        for (let i = 0; i < pathParts.length; i++) {
          const part = pathParts[i];
          // 如果参数以'@'开头，说明是包名，需要合并后续的路径部分
          if (part && part.startsWith('@') && i + 1 < pathParts.length) {
            // 合并包名和版本/路径部分
            const packageName = `${part}/${pathParts[i + 1]}`;
            args.push(packageName);
            i++; // 跳过下一个部分，因为已经合并了
          } else {
            if (part) args.push(part);
          }
        }

        // 提取环境变量（从查询参数中获取'env'的值）
        const envParam = url.searchParams.get('env');
        const env: Record<string, string> = {};

        if (envParam) {
          // 将env参数按逗号分割成多个键值对，然后按第一个冒号分割[6](@ref)
          const envPairs = envParam.split(',');
          for (const pair of envPairs) {
            const colonIndex = pair.indexOf(':');
            if (colonIndex !== -1) {
              const key = pair.substring(0, colonIndex);
              const value = pair.substring(colonIndex + 1);
              env[key] = value;
            }
          }
        }

        return {
          transportType: 'stdio',
          clientName: args[1] || args[0] || 'stdio',
          serverUrl: 'http://localhost:3000',
          stdioConfig: {
            command,
            args,
            env,
          },
        };
      }

      default:
        throw new Error(`不支持的协议类型: ${protocol}`);
    }
  } catch (error) {
    throw new Error(`URL解析错误: ${error instanceof Error ? error.message : '未知错误'}`);
  }
}
