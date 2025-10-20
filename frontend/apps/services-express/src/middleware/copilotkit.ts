import { Request, Response, NextFunction } from 'express';
import {
  CopilotRuntime,
  OpenAIAdapter,
  LangChainAdapter,
  copilotRuntimeNodeHttpEndpoint,
  langGraphPlatformEndpoint,
} from '@copilotkit/runtime';
import { ChatOpenAI } from '@langchain/openai';
import { ChatOllama } from '@langchain/ollama';
import { ChatAnthropic } from '@langchain/anthropic';
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

export const copilotkitHandler = (req: Request, res: Response, next: NextFunction) => {
  (async () => {
    // 直连mcpclient的agent
    const runtime = new CopilotRuntime({
      remoteEndpoints: [
        langGraphPlatformEndpoint({
          deploymentUrl: config.agentDeploymentUrl,
          langsmithApiKey: config.langsmithAiKey,
          agents: [
            {
              name: 'sample_agent',
              description: 'A helpful LLM agent.',
            },
          ],
        }),
      ],
    });

    // const runtime = new CopilotRuntime();

    const handler = copilotRuntimeNodeHttpEndpoint({
      endpoint: '/copilotkit',
      runtime,
      serviceAdapter: llmType?.[config.llmType],
    });
    return handler(req, res);
  })().catch(next);
};
