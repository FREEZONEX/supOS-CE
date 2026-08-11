import { useUpdateEffect } from 'ahooks';
import type { RefObject } from 'react';
import { ThemeType, useThemeStore } from '@/stores/theme-store.ts';

export type FlowIframeTheme = 'light' | 'dark';

export const FLOW_IFRAME_THEME_MESSAGE = 'tier0FlowThemeChange';

const normalizeFlowIframeTheme = (theme: string): FlowIframeTheme => (theme === ThemeType.Dark ? 'dark' : 'light');

export const postFlowIframeTheme = (iframe: HTMLIFrameElement | null, theme: FlowIframeTheme) => {
  iframe?.contentWindow?.postMessage(
    {
      type: FLOW_IFRAME_THEME_MESSAGE,
      data: { theme },
    },
    '*'
  );
};

const useFlowIframeThemeBridge = (iframeRef: RefObject<HTMLIFrameElement | null>): FlowIframeTheme => {
  const theme = useThemeStore((state) => state.theme);
  const iframeTheme = normalizeFlowIframeTheme(theme);

  useUpdateEffect(() => {
    postFlowIframeTheme(iframeRef.current, iframeTheme);
  }, [iframeTheme]);

  return iframeTheme;
};

export default useFlowIframeThemeBridge;
