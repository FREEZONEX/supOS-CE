import express, { Request, Response } from 'express';
import { mcpManager } from '@/utils';

const mcpRouter = express.Router();

// 详细健康检查
mcpRouter.get('/list', (_: Request, res: Response) => {
  try {
    res.status(200).json({
      status: 'ok',
      list: mcpManager.getMCPClientCache(),
      size: mcpManager.getMCPClientCount(),
    });
  } catch (e) {
    res.status(500).json({
      error: e,
    });
  }
});

// 刷新某个mcp
mcpRouter.post('/refresh', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        success: false,
        message: '缺少必需的参数: endpoint',
      });
    }

    const result = await mcpManager.refreshMCPClient(endpoint);

    if (result.success) {
      return res.status(200).json({
        success: true,
        message: result.message,
      });
    } else {
      return res.status(400).json({
        success: false,
        message: result.message,
      });
    }
  } catch (e) {
    console.error('刷新MCP客户端失败:', e);
    return res.status(500).json({
      success: false,
      message: `刷新MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

// 重启某个mcp
mcpRouter.post('/restart', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        success: false,
        message: '缺少必需的参数: endpoint',
      });
    }

    const result = await mcpManager.restartMCPClient(endpoint);

    if (result.success) {
      return res.status(200).json({
        success: true,
        message: result.message,
      });
    } else {
      return res.status(400).json({
        success: false,
        message: result.message,
      });
    }
  } catch (e) {
    console.error('重启MCP客户端失败:', e);
    return res.status(500).json({
      success: false,
      message: `重启MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

// 停止某个mcp
mcpRouter.post('/stop', async (req: Request, res: Response) => {
  try {
    const { endpoint } = req.body;

    if (!endpoint) {
      return res.status(400).json({
        success: false,
        message: '缺少必需的参数: endpoint',
      });
    }

    const result = await mcpManager.stopMCPClient(endpoint);

    if (result.success) {
      return res.status(200).json({
        success: true,
        message: result.message,
      });
    } else {
      return res.status(400).json({
        success: false,
        message: result.message,
      });
    }
  } catch (e) {
    console.error('停止MCP客户端失败:', e);
    return res.status(500).json({
      success: false,
      message: `停止MCP客户端失败: ${e instanceof Error ? e.message : String(e)}`,
    });
  }
});

export { mcpRouter };
