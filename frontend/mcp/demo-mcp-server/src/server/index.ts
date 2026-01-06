import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { formatErrorMessage, formatWeatherResult, WeatherTool } from '../tools/weather.js';
import { z } from 'zod';
import { logger } from '../utils';

export const createServer = () => {
  const server = new McpServer(
    {
      name: 'demo-mcp-server',
      title: 'Demo MCP Server',
      version: '1.0.0',
    },
    {
      capabilities: {
        // prompts: {},
        // resources: { subscribe: true },
        tools: {},
        // logging: {},
        // completions: {},
      },
      // instructions,
    }
  );
  // 创建天气工具实例
  // 注册天气工具
  server.registerTool(
    'get_weather',
    {
      description: '获取指定城市的天气信息',
      inputSchema: {
        city: z.string().describe('城市名称（例如：北京、上海、杭州）'),
      },
    },
    async ({ city }) => {
      const weatherTool = new WeatherTool();
      try {
        // 执行天气查询
        const result = await weatherTool.getWeatherByCity(city);
        const structuredContent = formatWeatherResult(result);
        logger.log(`[WeatherMCPTool] 天气查询成功完成`);
        return {
          content: [
            {
              type: 'text',
              text: structuredContent,
            },
          ],
        };
      } catch (error) {
        logger.error(`[WeatherMCPTool] 执行天气查询失败:`, error);
        const errorMessage = formatErrorMessage(error);

        return {
          content: [
            {
              type: 'text',
              text: errorMessage,
            },
          ],
        };
      }
    }
  );

  return { server };
};
