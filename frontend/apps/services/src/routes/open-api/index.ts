import { Hono } from 'hono';
import { healthRoutes } from './health';

const openApi = new Hono();

openApi.route('/open-api', healthRoutes);

export { openApi };
