import { Hono } from 'hono';
import { convertConfigToUrl, mcpManager } from '@/utils';

const mcpRouter = new Hono();

// mcp列表（包含工具信息）
mcpRouter.get('/list', async (c) => {
  try {
    const mcpCache = mcpManager.getMCPClientCache();
    // 合并MCP客户端信息和工具信息
    const enhancedData = await Promise.all(
      mcpCache.map(async (item) => {
        const toolsResult = await mcpManager.getToolsListByEndpoint(item.endpoint);
        const tools = toolsResult.success ? toolsResult.data.flatMap((item) => item.tools) : [];
        return {
          ...item,
          tools: tools,
        };
      })
    );
    return c.json(
      {
        code: 200,
        data: enhancedData,
        msg: 'success',
      },
      200
    );
  } catch (e) {
    console.error('获取MCP列表时发生错误:', e);
    return c.json(
      {
        code: 500,
        error: e,
        msg: `获取MCP列表失败: ${e instanceof Error ? e.message : String(e)}`,
      },
      500
    );
  }
});

// 新增mcp客户端（支持单个或批量添加）
mcpRouter.post('/add', async (c) => {
  try {
    const body = await c.req.json();
    const configs = Array.isArray(body) ? body : [body];

    if (configs.length === 0) {
      return c.json(
        {
          code: 400,
          data: null,
          msg: '缺少MCP配置参数',
        },
        400
      );
    }

    const results = [];
    const errors = [];

    for (const config of configs) {
      try {
        const endpoint = convertConfigToUrl(config);
        if (!endpoint) {
          errors.push({
            config,
            error: '缺少必需的参数: name',
          });
          continue;
        }

        const clientConfig = {
          endpoint,
        };

        await mcpManager.getOrCreateMCPClient(clientConfig);

        results.push({
          endpoint: endpoint,
          isConnected: true,
          lastUsed: Date.now(),
        });
      } catch (e) {
        errors.push({
          config,
          error: e instanceof Error ? e.message : String(e),
        });
      }
    }

    if (results.length === 0 && errors.length > 0) {
      return c.json(
        {
          code: 400,
          data: null,
          msg: '所有MCP配置添加失败',
          errors: errors.map((e) => e.error),
        },
        400
      );
    }

    const successCount = results.length;
    const errorCount = errors.length;

    return c.json(
      {
        code: 200,
        data: {
          results,
          errors: errorCount > 0 ? errors : undefined,
          summary: {
            total: configs.length,
            success: successCount,
            failed: errorCount,
          },
        },
        msg: `成功添加 ${successCount} 个MCP客户端${errorCount > 0 ? `，失败 ${errorCount} 个` : ''}`,
      },
      200
    );
  } catch (e) {
    return c.json(
      {
        code: 500,
        data: null,
        msg: `创建MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
      },
      500
    );
  }
});

// 删除某个mcp
mcpRouter.post('/delete', async (c) => {
  try {
    const body = await c.req.json();
    const { endpoint } = body;

    if (!endpoint) {
      return c.json(
        {
          code: 400,
          data: null,
          msg: '缺少必需的参数: endpoint',
        },
        400
      );
    }

    await mcpManager.removeMCPClient(endpoint);

    return c.json(
      {
        code: 200,
        data: null,
        msg: `MCP客户端 ${endpoint} 删除成功`,
      },
      200
    );
  } catch (e) {
    return c.json(
      {
        code: 500,
        data: null,
        msg: `删除MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
      },
      500
    );
  }
});

// 刷新某个mcp
mcpRouter.post('/refresh', async (c) => {
  try {
    const body = await c.req.json();
    const { endpoint } = body;

    if (!endpoint) {
      return c.json(
        {
          code: 400,
          data: null,
          msg: '缺少必需的参数: endpoint',
        },
        400
      );
    }

    const result = await mcpManager.refreshMCPClient(endpoint);

    if (result.success) {
      return c.json(
        {
          code: 200,
          data: null,
          msg: result.message,
        },
        200
      );
    } else {
      return c.json(
        {
          code: 400,
          data: null,
          msg: result.message,
        },
        400
      );
    }
  } catch (e) {
    return c.json(
      {
        code: 500,
        data: null,
        msg: `刷新MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
      },
      500
    );
  }
});

// 重启某个mcp
mcpRouter.post('/restart', async (c) => {
  try {
    const body = await c.req.json();
    const { endpoint } = body;

    if (!endpoint) {
      return c.json(
        {
          code: 400,
          data: null,
          msg: '缺少必需的参数: endpoint',
        },
        400
      );
    }

    const result = await mcpManager.restartMCPClient(endpoint);

    if (result.success) {
      return c.json(
        {
          code: 200,
          data: null,
          msg: result.message,
        },
        200
      );
    } else {
      return c.json(
        {
          code: 400,
          data: null,
          msg: result.message,
        },
        400
      );
    }
  } catch (e) {
    return c.json(
      {
        code: 500,
        data: null,
        msg: `重启MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
      },
      500
    );
  }
});

// 停止某个mcp
mcpRouter.post('/stop', async (c) => {
  try {
    const body = await c.req.json();
    const { endpoint } = body;

    if (!endpoint) {
      return c.json(
        {
          code: 400,
          data: null,
          msg: '缺少必需的参数: endpoint',
        },
        400
      );
    }

    const result = await mcpManager.stopMCPClient(endpoint);

    if (result.success) {
      return c.json(
        {
          code: 200,
          data: null,
          msg: result.message,
        },
        200
      );
    } else {
      return c.json(
        {
          code: 400,
          data: null,
          msg: result.message,
        },
        400
      );
    }
  } catch (e) {
    return c.json(
      {
        code: 500,
        data: null,
        msg: `停止MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
      },
      500
    );
  }
});

export { mcpRouter };
