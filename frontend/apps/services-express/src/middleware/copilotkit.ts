import { Request, Response, NextFunction } from 'express';
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
import OpenAI from 'openai';
import { config } from '@/config';
import { MCPClient, parseTransportUrl } from '@/utils';

const ollamaModel = new ChatOllama({
  model: config.ollamaModal,
  baseUrl: config.ollamaBaseurl,
});

const openaiModel = new ChatOpenAI({
  model: config.openAiModel,
  apiKey: config.openAiKey,
});

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

const serviceAdapterByTongyi = new OpenAIAdapter({
  openai: alibabaTongyiModel,
  model: config.tongyiModal,
  keepSystemRole: true,
});

const llmType: any = {
  ollama: serviceAdapterByllama,
  openai: serviceAdapterByOpenai,
  anthropic: serviceAdapterByAnthropic,
  tongyi: serviceAdapterByTongyi,
};

// MCP客户端缓存，避免重复创建
// const mcpClientCache: MCPClient | null = null;
// const mcpEndpointCache: string | null = null;

export const copilotkitHandler = (req: Request, res: Response, next: NextFunction) => {
  (async () => {
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
    const runtime = new CopilotRuntime({
      createMCPClient: async (config) => {
        const props = parseTransportUrl(config.endpoint);
        // 检查是否可以使用缓存的客户端
        // if (mcpClientCache && mcpEndpointCache === config.endpoint) {
        //   console.log('使用缓存的MCP客户端');
        //   return mcpClientCache;
        // }
        const mcpClient = new MCPClient(props);
        await mcpClient.connect();
        // 缓存客户端和端点
        // mcpClientCache = mcpClient;
        // mcpEndpointCache = config.endpoint;
        return mcpClient;
      },
    });
    const handler = copilotRuntimeNodeHttpEndpoint({
      endpoint: '/copilotkit',
      runtime,
      serviceAdapter: llmType?.[config.llmType],
    });
    return handler(req, res);
  })().catch((e) => {
    console.log(e);
    next(e);
  });
};
