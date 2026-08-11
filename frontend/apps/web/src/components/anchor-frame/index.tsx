import { useEffect, useMemo, useRef, useState, type FC } from 'react';
import { Button, Spin } from 'antd';
import { useMemoizedFn } from 'ahooks';
import {
  createAnchorHost,
  type AnchorContextPayload,
  type AnchorHost,
  type AnchorRoutePayload,
} from '@internal/anchor-bridge';
import type { PageProps } from '@/common-types';
import { useTabLifecycle, useActivate, useUnActivate } from '@/contexts/tabs-lifecycle-context.ts';
import useTranslate from '@/hooks/useTranslate.ts';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { ThemeType, useThemeStore } from '@/stores/theme-store.ts';

// Anchor 3D 子应用的专用 iframe 宿主（替代通用 DynamicIframe 承载 /anchor/*）：
//   - 只允许同源 /anchor/* 地址，菜单配置不能把它变成任意 URL 入口；
//   - 通过 anchor-bridge 向子应用推送主题/主色/语言上下文与页签生命周期；
//   - 带加载态、超时提示与重试；子应用路由变化只回传标题展示，不改写宿主路由。
const LOAD_TIMEOUT_MS = 20_000;

const isSafeAnchorUrl = (url?: string) => Boolean(url && /^\/anchor(\/|$|\?)/.test(url));

interface AnchorFrameProps extends PageProps {
  url?: string;
  name?: string;
}

const AnchorFrame: FC<AnchorFrameProps> = ({ url, name, location }) => {
  const formatMessage = useTranslate();
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const hostRef = useRef<AnchorHost | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'timeout'>('loading');
  const [reloadKey, setReloadKey] = useState(0);
  const [frameTitle, setFrameTitle] = useState<string | undefined>(undefined);

  const theme = useThemeStore((state) => state.theme);
  const primaryColor = useThemeStore((state) => state.primaryColor);
  const lang = useI18nStore((state) => state.lang);
  const { isShowRef } = useTabLifecycle();

  const context = useMemo<AnchorContextPayload>(
    () => ({
      theme: theme === ThemeType.Dark ? 'dark' : 'light',
      primaryColor,
      language: lang,
    }),
    [theme, primaryColor, lang]
  );
  const contextRef = useRef(context);
  contextRef.current = context;

  const search = location?.search ?? '';
  const src = isSafeAnchorUrl(url) ? `${url}${search}` : '';

  // bridge 宿主：iframe 挂载即监听；ready 后推送初始上下文与当前激活态
  useEffect(() => {
    const iframe = iframeRef.current;
    if (!iframe || !src) return;
    const host = createAnchorHost(iframe, {
      onReady: () => {
        host.sendContext(contextRef.current);
        host.sendLifecycle({ active: isShowRef?.current !== false });
      },
      onRoute: (payload: AnchorRoutePayload) => {
        setFrameTitle(payload.title || undefined);
      },
    });
    hostRef.current = host;
    return () => {
      host.dispose();
      hostRef.current = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [src, reloadKey]);

  // 主题/主色/语言变化 → 实时推送（子应用未接 bridge 时该消息被忽略，无副作用）
  useEffect(() => {
    hostRef.current?.sendContext(context);
  }, [context]);

  // 页签激活/隐藏 → 生命周期通知（子应用据此暂停/恢复渲染与订阅）
  useActivate(
    useMemoizedFn(() => {
      hostRef.current?.sendLifecycle({ active: true });
    })
  );
  useUnActivate(
    useMemoizedFn(() => {
      hostRef.current?.sendLifecycle({ active: false });
    })
  );

  // 加载超时提示（iframe load 即认为可用，bridge ready 与否不阻塞展示）
  useEffect(() => {
    if (!src) return;
    setStatus('loading');
    const timer = window.setTimeout(() => {
      setStatus((prev) => (prev === 'loading' ? 'timeout' : prev));
    }, LOAD_TIMEOUT_MS);
    return () => window.clearTimeout(timer);
  }, [src, reloadKey]);

  if (!src) {
    return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
        {formatMessage('anchor.frameInvalidUrl', {}, 'Invalid Anchor page address')}
      </div>
    );
  }

  return (
    <div style={{ position: 'relative', width: '100%', height: '100%' }}>
      <iframe
        key={reloadKey}
        ref={iframeRef}
        src={src}
        title={frameTitle || name || 'Anchor'}
        allowFullScreen
        allow="fullscreen; clipboard-write"
        style={{ width: '100%', height: '100%', border: 'none', display: 'block' }}
        onLoad={() => setStatus('ready')}
      />
      {status !== 'ready' ? (
        <div
          style={{
            position: 'absolute',
            inset: 0,
            display: 'flex',
            flexDirection: 'column',
            gap: 12,
            alignItems: 'center',
            justifyContent: 'center',
            background: 'var(--ui-bg-color)',
          }}
        >
          {status === 'loading' ? (
            <Spin />
          ) : (
            <>
              <span style={{ color: 'var(--ui-text-color)' }}>
                {formatMessage('anchor.frameTimeout', {}, '3D page took too long to load')}
              </span>
              <Button
                type="primary"
                onClick={() => {
                  setStatus('loading');
                  setReloadKey((key) => key + 1);
                }}
              >
                {formatMessage('anchor.frameRetry', {}, 'Retry')}
              </Button>
            </>
          )}
        </div>
      ) : null}
    </div>
  );
};

export default AnchorFrame;
