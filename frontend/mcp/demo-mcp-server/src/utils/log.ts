/**
 * A simple logger that writes to stderr to avoid interfering with stdout.
 */

const isPiped = !process.stdout.isTTY;

export const logger = {
  log: (...args: any[]) => {
    if (isPiped) {
      console.error('正常日志：', ...args);
    } else {
      console.log(...args);
    }
  },
  info: (...args: any[]) => {
    console.info(...args);
  },
  warn: (...args: any[]) => {
    console.warn(...args);
  },
  error: (...args: any[]) => {
    console.error(...args);
  },
};
