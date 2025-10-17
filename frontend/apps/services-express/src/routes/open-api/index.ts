import express from 'express';
import { healthRouter } from './health';

const openApiRouter = express.Router();

openApiRouter.use('/open-api', healthRouter);

export { openApiRouter };
