const _Env = Bun.env;
// 环境配置
export const config = {
  // 服务器配置
  port: parseInt(_Env.PORT || '4000', 10),
  nodeEnv: _Env.NODE_ENV || 'development',
  // Docker配置
  dockerHost: _Env.DOCKER_HOST || 'host.docker.internal',
  dockerPort: parseInt(_Env.DOCKER_PORT || '2376', 10),
  // ai配置
  llmType: _Env.LLM_TYPE || 'openai',
  ollamaModal: _Env.LLM_MODEL || 'llama3.2',
  ollamaBaseurl: _Env.LLM_BASEURL || 'http://office.unibutton.com:11435',
  openAiKey: _Env.LLM_API_KEY || 'openApiKey',
  openAiModel: _Env.LLM_MODEL || 'gpt-4o',
  anthropicAiKey: _Env.LLM_API_KEY || 'anthropicAiKey',
  anthropicAiModel: _Env.LLM_MODEL || 'claude-3-7-sonnet-thinking',
  tongyiAiKey: _Env.LLM_API_KEY || 'tongyiAiKey',
  tongyiModal: _Env.LLM_MODEL || 'gpt-4o',
  tongyiBaseurl: _Env.LLM_BASEURL || 'https://dashscope.aliyuncs.com/compatible-mode/v1',
  agentDeploymentUrl: _Env.AGENT_DEPLOYMENT_URL || 'http://localhost:8123',
  langsmithAiKey: _Env.LANGSMITH_API_KEY || 'langsmithAiKey',
} as const;

export type Config = typeof config;
