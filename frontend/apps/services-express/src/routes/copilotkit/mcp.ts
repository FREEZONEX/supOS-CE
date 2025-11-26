import express, { Request, Response } from 'express';
import { convertConfigToUrl, mcpManager } from '@/utils';

const mcpRouter = express.Router();

// mcp列表
mcpRouter.get('/list', (_: Request, res: Response) => {
  try {
    res.status(200).json({
      code: 200,
      data: mcpManager.getMCPClientCache(),
      msg: 'success',
    });
  } catch (e) {
    res.status(500).json({
      error: e,
    });
  }
});

// 新增mcp客户端（支持单个或批量添加）
mcpRouter.post('/add', async (req: Request, res: Response) => {
  try {
    const configs = Array.isArray(req.body) ? req.body : [req.body];

    if (configs.length === 0) {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: '缺少MCP配置参数',
      });
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
      return res.status(400).json({
        code: 400,
        data: null,
        msg: '所有MCP配置添加失败',
        errors: errors.map((e) => e.error),
      });
    }

    const successCount = results.length;
    const errorCount = errors.length;

    return res.status(200).json({
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
    });
  } catch (e) {
    return res.status(500).json({
      code: 500,
      data: null,
      msg: `创建MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

// 删除某个mcp
mcpRouter.post('/delete', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: '缺少必需的参数: endpoint',
      });
    }

    await mcpManager.removeMCPClient(endpoint);

    return res.status(200).json({
      code: 200,
      data: null,
      msg: `MCP客户端 ${endpoint} 删除成功`,
    });
  } catch (e) {
    return res.status(500).json({
      code: 500,
      data: null,
      msg: `删除MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

// 刷新某个mcp
mcpRouter.post('/refresh', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: '缺少必需的参数: endpoint',
      });
    }

    const result = await mcpManager.refreshMCPClient(endpoint);

    if (result.success) {
      return res.status(200).json({
        code: 200,
        data: null,
        msg: result.message,
      });
    } else {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: result.message,
      });
    }
  } catch (e) {
    return res.status(500).json({
      code: 500,
      data: null,
      msg: `刷新MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

// 重启某个mcp
mcpRouter.post('/restart', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: '缺少必需的参数: endpoint',
      });
    }

    const result = await mcpManager.restartMCPClient(endpoint);

    if (result.success) {
      return res.status(200).json({
        code: 200,
        data: null,
        msg: result.message,
      });
    } else {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: result.message,
      });
    }
  } catch (e) {
    return res.status(500).json({
      code: 500,
      data: null,
      msg: `重启MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

// 停止某个mcp
mcpRouter.post('/stop', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: '缺少必需的参数: endpoint',
      });
    }

    const result = await mcpManager.stopMCPClient(endpoint);

    if (result.success) {
      return res.status(200).json({
        code: 200,
        data: null,
        msg: result.message,
      });
    } else {
      return res.status(400).json({
        code: 400,
        data: null,
        msg: result.message,
      });
    }
  } catch (e) {
    return res.status(500).json({
      code: 500,
      data: null,
      msg: `停止MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

export { mcpRouter };
