import { SSEServerTransport } from '@modelcontextprotocol/sdk/server/sse.js';
import express from 'express';
import { createServer } from '../server';
import cors from 'cors';
import { logger } from '../utils';

logger.error('Starting SSE server...');

const app = express();
app.use(
  cors({
    origin: '*',
    methods: 'GET,POST',
    preflightContinue: false,
    optionsSuccessStatus: 204,
  })
);

const transports: Map<string, SSEServerTransport> = new Map<string, SSEServerTransport>();

app.get('/sse', async (req, res) => {
  let transport: SSEServerTransport;
  const { server } = createServer();

  if (req?.query?.sessionId) {
    const sessionId = req?.query?.sessionId as string;
    transport = transports.get(sessionId) as SSEServerTransport;
    logger.error(
      "Client Reconnecting? This shouldn't happen; when client has a sessionId, GET /sse should not be called again.",
      transport.sessionId
    );
  } else {
    transport = new SSEServerTransport('/message', res);
    transports.set(transport.sessionId, transport);

    await server.connect(transport);
    logger.error('Client Connected: ', transport.sessionId);
  }
});

app.post('/message', async (req, res) => {
  const sessionId = req?.query?.sessionId as string;
  const transport = transports.get(sessionId);
  if (transport) {
    logger.error('Client Message from', sessionId);
    await transport.handlePostMessage(req, res);
  } else {
    logger.error(`No transport found for sessionId ${sessionId}`);
  }
});

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  logger.error(`MCP Sse Server listening on (http://localhost:${PORT}/sse)`);
});
