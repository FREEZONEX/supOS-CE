import express from 'express';
import { mcpRouter } from './mcp';

const mcpApiRouter = express.Router();

mcpApiRouter.use('/', mcpRouter);

export { mcpApiRouter };
