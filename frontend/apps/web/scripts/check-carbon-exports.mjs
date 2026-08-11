import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const carbonPath = path.join(root, 'src/components/lucide-icon/carbon.ts');
const srcRoot = path.join(root, 'src');

const carbonSource = fs.readFileSync(carbonPath, 'utf8');
const exportBlock = carbonSource.match(/export\s*\{([\s\S]*?)\}\s*from\s*['"]lucide-react['"]/);
if (!exportBlock) {
  console.error('check-carbon-exports: cannot parse carbon.ts exports');
  process.exit(1);
}

const exportedNames = new Set();
for (const line of exportBlock[1].split(',')) {
  const trimmed = line.trim();
  if (!trimmed) continue;
  const aliasMatch = trimmed.match(/\bas\s+(\w+)\s*$/);
  exportedNames.add(aliasMatch ? aliasMatch[1] : trimmed.split(/\s+/)[0]);
}

const importPattern = /import\s*\{([^}]+)\}\s*from\s*['"]@\/components\/lucide-icon\/carbon['"]/g;
const missing = new Map();

const walk = (dir) => {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name === 'node_modules') continue;
      walk(fullPath);
      continue;
    }
    if (!/\.(tsx?|jsx?)$/.test(entry.name)) continue;
    const content = fs.readFileSync(fullPath, 'utf8');
    let match;
    while ((match = importPattern.exec(content)) !== null) {
      const importBlock = match[1].replace(/\/\/.*$/gm, '');
      for (const part of importBlock.split(',')) {
        const cleaned = part.trim();
        if (!cleaned) continue;
        const sourceName = cleaned.includes(' as ') ? cleaned.split(/\s+as\s+/)[0].trim() : cleaned;
        if (!sourceName || exportedNames.has(sourceName)) continue;
        const rel = path.relative(root, fullPath);
        if (!missing.has(sourceName)) missing.set(sourceName, new Set());
        missing.get(sourceName).add(rel);
      }
    }
  }
};

walk(srcRoot);

if (missing.size > 0) {
  console.error('check-carbon-exports: missing exports from carbon.ts:');
  for (const [name, files] of missing.entries()) {
    console.error(`  - ${name}: ${[...files].join(', ')}`);
  }
  process.exit(1);
}

console.log('check-carbon-exports: ok');
