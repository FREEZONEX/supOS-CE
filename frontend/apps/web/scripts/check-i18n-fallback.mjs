import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { fallbackKeyGroups, fallbackKeys } from './i18n-fallback-keys.mjs';

const webRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const fallbackKeyLimit = 160;
const files = {
  zh: 'src/locale/zh-CN.json',
  en: 'src/locale/en-US.json',
};

const readJson = (relativePath) => JSON.parse(fs.readFileSync(path.join(webRoot, relativePath), 'utf8'));
const sortedKeys = (value) => Object.keys(value).sort();
const missingKeys = (expected, actual) => expected.filter((key) => !actual.includes(key));
const declaredKeys = Object.values(fallbackKeyGroups).flat();
const duplicateKeys = declaredKeys.filter((key, index) => declaredKeys.indexOf(key) !== index);
const expectedKeys = [...fallbackKeys].sort();
const errors = [];

if (duplicateKeys.length > 0) {
  errors.push(`fallback allowlist contains duplicate keys: ${[...new Set(duplicateKeys)].join(', ')}`);
}
if (expectedKeys.length > fallbackKeyLimit) {
  errors.push(
    `frontend fallback has ${expectedKeys.length} allowlisted keys; hard limit is ${fallbackKeyLimit}. ` +
      'Move business translations to backend/etc/i18n instead of expanding the frontend bundle.'
  );
}

const checkKeySet = (label, actual) => {
  const missing = missingKeys(expectedKeys, actual);
  const extra = missingKeys(actual, expectedKeys);
  if (missing.length > 0) errors.push(`${label} missing keys: ${missing.join(', ')}`);
  if (extra.length > 0) errors.push(`${label} unexpected keys: ${extra.join(', ')}`);
};

checkKeySet(files.zh, sortedKeys(readJson(files.zh)));
checkKeySet(files.en, sortedKeys(readJson(files.en)));

if (errors.length > 0) {
  console.error('frontend i18n fallback check failed:');
  for (const error of errors) console.error(`- ${error}`);
  console.error('Run from repository root: node scripts/generate-i18n-fallback.mjs');
  process.exit(1);
}

console.log(`frontend i18n fallback check passed: ${expectedKeys.length} allowlisted keys.`);
