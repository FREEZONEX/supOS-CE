import '@ant-design/v5-patch-for-react-19';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router';
import App from './App';
import { ThemeBridge } from './theme/ThemeBridge';
import './index.css';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter basename="/anchor">
      <ThemeBridge>
        <App />
      </ThemeBridge>
    </BrowserRouter>
  </StrictMode>
);
