import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import {
  fallbackKeyGroups,
  fallbackKeys,
} from "../frontend/apps/web/scripts/i18n-fallback-keys.mjs";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");

const files = {
  backendZh: "backend/etc/i18n/zh_CN.json",
  backendEn: "backend/etc/i18n/en_US.json",
  fallbackZh: "frontend/apps/web/src/locale/zh-CN.json",
  fallbackEn: "frontend/apps/web/src/locale/en-US.json",
};

const fallbackKeyLimit = 160;

const readJson = (relativePath) => {
  const absolutePath = path.join(root, relativePath);
  return JSON.parse(fs.readFileSync(absolutePath, "utf8"));
};

const sortedKeys = (value) => Object.keys(value).sort();
const missingKeys = (expected, actual) =>
  expected.filter((key) => !actual.includes(key));

const backendZh = readJson(files.backendZh);
const backendEn = readJson(files.backendEn);
const fallbackZh = readJson(files.fallbackZh);
const fallbackEn = readJson(files.fallbackEn);
const errors = [];

const backendZhKeys = sortedKeys(backendZh);
const backendEnKeys = sortedKeys(backendEn);
const fallbackZhKeys = sortedKeys(fallbackZh);
const fallbackEnKeys = sortedKeys(fallbackEn);
const expectedFallbackKeys = [...fallbackKeys].sort();
const declaredFallbackKeys = Object.values(fallbackKeyGroups).flat();
const duplicateFallbackKeys = declaredFallbackKeys.filter(
  (key, index) => declaredFallbackKeys.indexOf(key) !== index,
);

if (duplicateFallbackKeys.length > 0) {
  errors.push(
    `fallback allowlist contains duplicate keys: ${[...new Set(duplicateFallbackKeys)].join(", ")}`,
  );
}
if (fallbackKeys.length > fallbackKeyLimit) {
  errors.push(
    `frontend fallback has ${fallbackKeys.length} allowlisted keys; hard limit is ${fallbackKeyLimit}`,
  );
}

const checkKeySet = (label, expected, actual) => {
  const missing = missingKeys(expected, actual);
  const extra = missingKeys(actual, expected);

  if (missing.length > 0) {
    errors.push(`${label} missing keys: ${missing.join(", ")}`);
  }
  if (extra.length > 0) {
    errors.push(`${label} unexpected keys: ${extra.join(", ")}`);
  }
};

checkKeySet("backend en_US.json", backendZhKeys, backendEnKeys);
checkKeySet("backend zh_CN.json", backendEnKeys, backendZhKeys);
checkKeySet("frontend zh-CN.json", expectedFallbackKeys, fallbackZhKeys);
checkKeySet("frontend en-US.json", expectedFallbackKeys, fallbackEnKeys);

for (const key of fallbackKeys) {
  if (!(key in backendZh) || !(key in backendEn)) {
    errors.push(`fallback key is missing from backend language packs: ${key}`);
    continue;
  }

  if (fallbackZh[key] !== backendZh[key]) {
    errors.push(`frontend zh-CN value differs from backend zh_CN: ${key}`);
  }
  if (fallbackEn[key] !== backendEn[key]) {
    errors.push(`frontend en-US value differs from backend en_US: ${key}`);
  }
}

if (errors.length > 0) {
  console.error("i18n check failed:");
  for (const error of errors) {
    console.error(`- ${error}`);
  }
  process.exit(1);
}

console.log(
  `i18n check passed: ${backendZhKeys.length} backend keys, ${fallbackKeys.length} frontend fallback keys.`,
);
