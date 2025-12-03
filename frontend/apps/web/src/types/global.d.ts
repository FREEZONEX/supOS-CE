// eslint-disable-next-line @typescript-eslint/no-unused-vars
import * as React from 'react';

declare module 'react' {
  interface CSSProperties {
    [key: `--${string}`]: string | number | undefined; // 允许所有以 `--` 开头的属性
  }
  interface Attributes {
    auth?: string | string[];
  }
}
