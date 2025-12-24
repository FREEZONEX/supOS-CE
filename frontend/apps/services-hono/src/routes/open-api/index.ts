import { Hono } from 'hono';
import { healthRouter } from './health';

const openApiRouter = new Hono();

openApiRouter.route('/open-api', healthRouter);

export { openApiRouter };
