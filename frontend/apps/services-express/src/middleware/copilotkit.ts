import { Request, Response, NextFunction } from 'express';
import {
  CopilotRuntime,
  LangChainAdapter,
  copilotRuntimeNodeHttpEndpoint,
  // langGraphPlatformEndpoint,
} from '@copilotkit/runtime';
import { ChatOpenAI } from '@langchain/openai';
import { ChatOllama } from '@langchain/ollama';
import { ChatAnthropic } from '@langchain/anthropic';
import { config } from '@/config';
import { mcpManager } from '@/utils';

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

const alibabaTongyiModel: any = new ChatOpenAI({
  apiKey: config.tongyiAiKey,
  model: config.tongyiModal,
  configuration: {
    baseURL: config.tongyiBaseurl,
  },
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

const serviceAdapterByTongyi = new LangChainAdapter({
  chainFn: async ({ messages, tools }) => {
    return alibabaTongyiModel.bindTools(tools).stream(messages);
  },
});

const llmType: any = {
  ollama: serviceAdapterByllama,
  openai: serviceAdapterByOpenai,
  anthropic: serviceAdapterByAnthropic,
  tongyi: serviceAdapterByTongyi,
};

let globalCopilotRuntime: CopilotRuntime | null = null;
// 存储老的mcpservers记录
let oldMcpServers: any = [];
/**
 * 创建或更新CopilotRuntime实例
 */
function createOrUpdateCopilotRuntime(): CopilotRuntime {
  const currentMcpServers = mcpManager?.getMCPClientCache()?.map((m) => ({ endpoint: m.endpoint })) || [];
  if (!globalCopilotRuntime || JSON.stringify(currentMcpServers) !== JSON.stringify(oldMcpServers)) {
    oldMcpServers = currentMcpServers;
    globalCopilotRuntime = new CopilotRuntime(
      currentMcpServers?.length > 0
        ? {
            mcpServers: currentMcpServers,
            createMCPClient: async (config) => {
              return await mcpManager.getOrCreateMCPClient(config);
            },
          }
        : {}
    );
  }
  return globalCopilotRuntime;
}

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
    const handler = copilotRuntimeNodeHttpEndpoint({
      endpoint: '/copilotkit',
      runtime: createOrUpdateCopilotRuntime(),
      serviceAdapter: llmType?.[config.llmType],
    });
    return handler(req, res);
  })().catch((e) => {
    console.log(e);
    next(e);
  });
};
