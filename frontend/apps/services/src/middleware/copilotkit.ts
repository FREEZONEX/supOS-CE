import {
  CopilotRuntime,
  OpenAIAdapter,
  LangChainAdapter,
  copilotRuntimeNodeHttpEndpoint,
  // langGraphPlatformEndpoint,
} from '@copilotkit/runtime';
import { ChatOpenAI } from '@langchain/openai';
import { ChatOllama } from '@langchain/ollama';
import { ChatAnthropic } from '@langchain/anthropic';
import type { Context, Next } from 'hono';
import OpenAI from 'openai';
import { config } from '@/config';

const ollamaModel = new ChatOllama({
  model: config.ollamaModal,
  baseUrl: config.ollamaBaseurl,
});

const openaiModel = new ChatOpenAI(
  {
    model: config.openAiModel,
    apiKey: config.openAiKey,
    // baseUrl: OLLAMA_BASEURL,
  }
  // {
  //   httpAgent: agentProxy,
  // }
);

const anthropicModel = new ChatAnthropic({
  model: config.anthropicAiModel,
  apiKey: config.anthropicAiKey,
});

const alibabaTongyiModel: any = new OpenAI({
  apiKey: config.tongyiAiKey,
  baseURL: config.tongyiBaseurl,
});

const serviceAdapterByllama = new LangChainAdapter({
  chainFn: async ({ messages, tools }) => {
    return ollamaModel.bindTools(tools).stream(messages);
  },
});

const serviceAdapterByOpenai = new LangChainAdapter({
  chainFn: async ({ messages, tools }) => {
    return openaiModel.bindTools(tools).stream(messages);
  },
});

const serviceAdapterByAnthropic = new LangChainAdapter({
  chainFn: async ({ messages, tools }) => {
    return anthropicModel.bindTools(tools).stream(messages);
  },
});

const serviceAdapterByTongyi = new OpenAIAdapter({ openai: alibabaTongyiModel, model: config.tongyiModal });

const llmType: any = {
  ollama: serviceAdapterByllama,
  openai: serviceAdapterByOpenai,
  anthropic: serviceAdapterByAnthropic,
  tongyi: serviceAdapterByTongyi,
};

export const copilotkitHandler = async (c: Context, next: Next): Promise<void | Response> => {
  // 直连mcpclient的agent
  // const runtime = new CopilotRuntime({
  //   remoteEndpoints: [
  //     langGraphPlatformEndpoint({
  //       deploymentUrl: config.agentDeploymentUrl,
  //       langsmithApiKey: config.langsmithAiKey,
  //       agents: [
  //         {
  //           name: 'sample_agent',
  //           description: 'A helpful LLM agent.',
  //         },
  //       ],
  //     }),
  //   ],
  // });

  const runtime = new CopilotRuntime();

  const handler = copilotRuntimeNodeHttpEndpoint({
    endpoint: '/copilotkit',
    runtime,
    serviceAdapter: llmType?.[config.llmType],
  });
  // 将 Hono 请求转换为标准的 Node.js 请求格式
  const request = new Request(c.req.url, {
    method: c.req.method,
    headers: c.req.raw.headers,
    body: c.req.raw.body,
  });

  // handler 期望接收一个包含 request 属性的对象
  const result = await handler({ request });

  // 如果 handler 返回了响应，则返回它，否则继续执行下一个中间件
  if (result) {
    return result;
  }

  // 如果没有返回响应，则继续执行下一个中间件
  return next();
};
