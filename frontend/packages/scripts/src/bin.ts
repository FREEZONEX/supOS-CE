#!/usr/bin/env node

import { runScript } from './index.ts';

const scriptName = process.argv[2] || '';
runScript(scriptName);
