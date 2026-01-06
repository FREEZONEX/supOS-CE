import fs from 'fs';
import { REMOTE_NAME } from './variables';

interface ConfigTypes {
  assetPrefix?: string;
  hostOrigin?: string;
  name: string;
}

let development: any = {};
let production: any = {};
try {
  development = await import('./supos-env.dev.ts');
} catch {
  development = {
    assetPrefix: 'http://localhost:5174',
    hostOrigin: 'http://localhost:5173',
  };
  fs.writeFileSync('./supos-env.dev.ts', `export default ${JSON.stringify(development, null, 2)}`);
}
try {
  production = await import('./supos-env.pro.ts');
} catch {
  production = {
    assetPrefix: '',
    hostOrigin: '',
  };
  fs.writeFileSync('./supos-env.pro.ts', `export default ${JSON.stringify(production, null, 2)}`);
}

const config: ConfigTypes = {
  ...(process.env.NODE_ENV === 'production' ? production?.default || production : development?.default || development),
  name: REMOTE_NAME,
};

export default config;
