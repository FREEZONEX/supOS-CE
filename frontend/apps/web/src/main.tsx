// import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import { CopilotKit } from '@copilotkit/react-core';
// import '@ant-design/v5-patch-for-react-19';

import './index.scss';

console.info(
  `%csupOS Frontend Version: ${import.meta.env.VITE_APP_VERSION}_${import.meta.env.VITE_APP_BUILD_TIMESTAMP}`,
  'color: #4CAF50; font-size: 16px; font-weight: bold;'
);

createRoot(document.getElementById('root')!).render(
  <CopilotKit
    // 也可以直接换成 mcpclient 的agent服务地址，不再使用.docker-nodejs 代理
    // runtimeUrl="/mcpclient/home/api/copilotkit"
    runtimeUrl="/copilotkit"
    // mcpEndpoints={[
    //   {
    //     endpoint: 'your_mcp_sse_url',
    //   },
    // ]}
    // agent="sample_agent"
    showDevConsole={false}
  >
    <App />
  </CopilotKit>
);
