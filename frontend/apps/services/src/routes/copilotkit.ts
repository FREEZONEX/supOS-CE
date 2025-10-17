// import { Hono } from 'hono';
// import {
//   CopilotRuntime,
//   OpenAIAdapter,
//   LangChainAdapter,
//   copilotRuntimeNodeHttpEndpoint,
//   langGraphPlatformEndpoint,
// } from '@copilotkit/runtime';
//
// // 假设这些变量在环境中定义
// const AGENT_DEPLOYMENT_URL = process.env.AGENT_DEPLOYMENT_URL || 'http://localhost:8123';
// const LANGSMITH_API_KEY = process.env.LANGSMITH_API_KEY;
//
// // 假设这些导入存在，根据实际情况调整
// const serviceAdapterByOpenai = {}; // 替换为实际的service adapter
// const llmType: any = {}; // 替换为实际的LLM类型配置
// const LLM_TYPE = 'openai'; // 替换为实际的LLM类型
//
// const copilotKitRoutes = new Hono();
//
// // CopilotKit端点
// copilotKitRoutes.use('/copilotkit', async (c, next) => {
//   try {
//     const runtime = new CopilotRuntime({
//       remoteEndpoints: [
//         langGraphPlatformEndpoint({
//           deploymentUrl: AGENT_DEPLOYMENT_URL,
//           langsmithApiKey: LANGSMITH_API_KEY,
//           agents: [
//             {
//               name: 'sample_agent',
//               description: 'A helpful LLM agent.',
//             },
//           ],
//         }),
//       ],
//     });
//     // 获取Node.js HTTP处理函数
//     const handler = copilotRuntimeNodeHttpEndpoint({
//       endpoint: '/copilotkit',
//       runtime,
//       serviceAdapter: llmType?.[LLM_TYPE] || serviceAdapterByOpenai,
//     });
//     return handler(c.req.raw, c.res);
//   } catch (e) {
//     console.log(e);
//     await next();
//   }
// });
//
// export { copilotKitRoutes };
