import { config as envConfig } from 'dotenv';

// 加载环境
envConfig();

// 环境配置
export const config = {
  // 服务器配置
  port: parseInt(process.env.PORT || '4000', 10),
  nodeEnv: process.env.NODE_ENV || 'development',
  // Docker配置
  dockerHost: process.env.DOCKER_HOST || 'host.docker.internal',
  dockerPort: parseInt(process.env.DOCKER_PORT || '2376', 10),
  llmType: process.env.LLM_TYPE || 'openai',
  ollamaModal: process.env.LLM_MODEL || 'llama3.2',
  ollamaBaseurl: process.env.LLM_BASEURL || 'http://office.unibutton.com:11435',
  openAiKey: process.env.LLM_API_KEY || 'openApiKey',
  openAiModel: process.env.LLM_MODEL || 'gpt-4o',
  anthropicAiKey: process.env.LLM_API_KEY || 'anthropicAiKey',
  anthropicAiModel: process.env.LLM_MODEL || 'claude-3-7-sonnet-thinking',
  tongyiAiKey: process.env.LLM_API_KEY || 'tongyiAiKey',
  tongyiModal: process.env.LLM_MODEL || 'gpt-4o',
  tongyiBaseurl: process.env.LLM_BASEURL || 'https://dashscope.aliyuncs.com/compatible-mode/v1',
  agentDeploymentUrl: process.env.AGENT_DEPLOYMENT_URL || 'http://localhost:8123',
  langsmithAiKey: process.env.LANGSMITH_API_KEY || 'langsmithAiKey',
} as const;

export type Config = typeof config;
