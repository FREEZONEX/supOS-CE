import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { fallbackKeys } from "../frontend/apps/web/scripts/i18n-fallback-keys.mjs";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const languageFiles = [
  {
    source: "backend/etc/i18n/zh_CN.json",
    target: "frontend/apps/web/src/locale/zh-CN.json",
  },
  {
    source: "backend/etc/i18n/en_US.json",
    target: "frontend/apps/web/src/locale/en-US.json",
  },
];

for (const { source, target } of languageFiles) {
  const sourcePath = path.join(root, source);
  const targetPath = path.join(root, target);
  const messages = JSON.parse(fs.readFileSync(sourcePath, "utf8"));
  const missing = fallbackKeys.filter((key) => !(key in messages));
  if (missing.length > 0) {
    throw new Error(
      `${source} is missing fallback keys: ${missing.join(", ")}`,
    );
  }

  const fallbackMessages = Object.fromEntries(
    fallbackKeys.map((key) => [key, messages[key]]),
  );
  fs.writeFileSync(
    targetPath,
    `${JSON.stringify(fallbackMessages, null, 2)}\n`,
    "utf8",
  );
  console.log(`generated ${target} with ${fallbackKeys.length} keys`);
}
