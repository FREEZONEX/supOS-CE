import { Hono } from 'hono';
import { mcpRouter } from './mcp';

const mcpApiRouter = new Hono();

mcpApiRouter.route('/copilotkit/mcp', mcpRouter);

export { mcpApiRouter };
