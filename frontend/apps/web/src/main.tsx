// import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
// import '@ant-design/v5-patch-for-react-19';

import './index.scss';

console.info(
  `%csupOS Frontend Version: ${import.meta.env.VITE_APP_VERSION}_${import.meta.env.VITE_APP_BUILD_TIMESTAMP}`,
  'color: #4CAF50; font-size: 16px; font-weight: bold;'
);

createRoot(document.getElementById('root')!).render(<App />);
