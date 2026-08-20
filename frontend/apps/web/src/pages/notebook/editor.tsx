import { createNotebookSnapshot, getNotebook, getNotebookSnapshotList } from '@/apis/core-api/notebook';
import { ButtonPermission } from '@/common-types/button-permission';
import ComBackButton from '@/components/com-back-button';
import ComContent from '@/components/com-layout/ComContent';
import { useTranslate } from '@/hooks';
import useTabName from '@/hooks/useTabName';
import type { NotebookDetail } from '@/pages/notebook/types';
import { hasPermission } from '@/utils/auth';
import { getDevProxyBaseUrl } from '@/utils/url-util';
import { ChevronDown, Renew, Save } from '@/components/lucide-icon/carbon';
import { App, Button, Popover, Space, Spin, Tooltip } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router';
import HistoryPanel from './components/HistoryPanel';
import HistoryPopover from './components/HistoryPopover';
import SaveVersionModal from './components/SaveVersionModal';
import {
  buildNotebookIframeDarkThemeCss,
  clearNotebookIframeDarkDocument,
  isNotebookDarkThemeActive,
  syncNotebookIframeDarkDocument,
} from './marimo-dark-theme';
import styles from './editor.module.scss';

const notebookIframeNameLimitStyleId = 'tier0-notebook-name-limit-style';
const notebookIframeDarkBootstrapStyleId = 'tier0-notebook-dark-bootstrap-style';
const notebookIframeParentThemeObserverKey = '__tier0NotebookParentThemeObserver';
const notebookIframeNameLimitObserverKey = '__tier0NotebookNameLimitObserver';
const notebookIframeNameLimitSignatureKey = '__tier0NotebookNameLimitSignature';
const notebookIframeNameLimitMaxWidth = 'min(340px, 26vw)';

const normalizeText = (value?: string | null) => (value || '').replace(/\s+/g, ' ').trim();

const getFileNameFromPath = (value?: string | null) => {
  const normalized = normalizeText(value);
  return normalized ? normalized.split(/[\\/]/).pop() || '' : '';
};

const isDarkThemeActive = isNotebookDarkThemeActive;

const detectNotebookIframeError = (iframe: HTMLIFrameElement | null) => {
  if (!iframe) {
    return null;
  }
  try {
    const doc = iframe.contentDocument || iframe.contentWindow?.document;
    if (!doc?.body) {
      return null;
    }
    const text = normalizeText(doc.body.innerText || doc.body.textContent || '');
    if (!text) {
      return null;
    }
    if (
      /^(dial tcp:|lookup |get "|post "|put "|patch "|delete "|502 bad gateway|503 service|504 gateway|bad gateway|invalid upstream|upstream denied)/i.test(
        text
      )
    ) {
      return text.split('\n')[0];
    }
    if (
      text.length < 400 &&
      /no such host|connection refused|i\/o timeout|connect: network is unreachable/i.test(text)
    ) {
      return text.split('\n')[0];
    }
  } catch {
    return null;
  }
  return null;
};

const NotebookEditorPage: FC = () => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const canManage = hasPermission(ButtonPermission['Notebook.manage']);
  const params = useParams();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const notebookId = params.id || '';
  const filePathFromSearch = searchParams.get('filePath') || '';
  const fromPath = searchParams.get('from') || '';
  const iframeRef = useRef<HTMLIFrameElement | null>(null);
  const [loading, setLoading] = useState(true);
  const [iframeLoading, setIframeLoading] = useState(true);
  const [iframeError, setIframeError] = useState('');
  const [iframeKey, setIframeKey] = useState(() => Date.now());
  const [detail, setDetail] = useState<NotebookDetail | null>(null);

  // 版本管理状态
  const [saveVersionOpen, setSaveVersionOpen] = useState(false);
  const [saveVersionLoading, setSaveVersionLoading] = useState(false);
  const [historyOpen, setHistoryOpen] = useState(false);
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [selectedSnapshotId, setSelectedSnapshotId] = useState<number | null>(null);
  const [latestVersionName, setLatestVersionName] = useState<string>('');
  const [historyRefreshKey, setHistoryRefreshKey] = useState(0);
  useTabName(detail?.name || formatMessage('Notebook.editor', {}, 'Notebook Editor'));
  const notebookDisplayNames = useMemo(() => {
    const name = normalizeText(detail?.name);
    const fileName = getFileNameFromPath(filePathFromSearch || detail?.filePath);
    return Array.from(new Set([name, name ? `${name}.py` : '', fileName].filter(Boolean)));
  }, [detail?.filePath, detail?.name, filePathFromSearch]);

  const loadLatestVersion = useCallback(async () => {
    if (!notebookId) return;
    try {
      const resp = await getNotebookSnapshotList(notebookId, { pageNo: 1, pageSize: 1 });
      const latest = resp?.list?.[0];
      if (latest?.versionName) {
        setLatestVersionName(latest.versionName);
      }
    } catch {
      // silent
    }
  }, [notebookId]);

  const loadDetail = useCallback(async () => {
    if (!notebookId) {
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const notebookResp = await getNotebook(notebookId);
      setDetail(notebookResp);
    } catch {
      message.error(formatMessage('Notebook.loadFailed', {}, 'Failed to load notebook'));
    } finally {
      setLoading(false);
    }
  }, [formatMessage, message, notebookId]);

  useEffect(() => {
    void loadDetail();
    void loadLatestVersion();
  }, [loadDetail, loadLatestVersion]);

  const iframeUrl = useMemo(() => {
    const filePath = filePathFromSearch || detail?.filePath;
    if (!filePath) {
      return '';
    }

    const iframePath = `${getDevProxyBaseUrl()}/marimo/home/`;
    const url = new URL(iframePath, window.location.origin);
    url.searchParams.set('file', filePath);
    return `${url.pathname}${url.search}`;
  }, [detail?.filePath, filePathFromSearch]);

  useEffect(() => {
    if (!iframeUrl) {
      setIframeLoading(false);
      setIframeError('');
    } else {
      setIframeLoading(true);
      setIframeError('');
    }
  }, [iframeKey, iframeUrl]);

  const handleSaveVersion = async (versionName: string, description: string) => {
    setSaveVersionLoading(true);
    try {
      await createNotebookSnapshot(notebookId, { versionName, description });
      message.success(formatMessage('Notebook.snapshot.saveSuccess', {}, 'Version saved'));
      setSaveVersionOpen(false);
      return true;
    } catch {
      // error handled by interceptor
      return false;
    } finally {
      setSaveVersionLoading(false);
    }
  };

  const handleReverted = async () => {
    await loadDetail();
    setIframeError('');
    setIframeLoading(true);
    setIframeKey(Date.now());
  };

  const handlePopoverSelectItem = (snapshotId: number) => {
    setSelectedSnapshotId(snapshotId);
    setPopoverOpen(false);
    setHistoryOpen(true);
  };

  const handlePopoverViewAll = () => {
    setSelectedSnapshotId(null);
    setPopoverOpen(false);
    setHistoryOpen(true);
  };

  const handleBack = () => {
    if (fromPath) {
      navigate(fromPath);
    } else {
      navigate(-1);
    }
  };

  const clearNotebookIframeDarkBootstrap = useCallback(() => {
    const iframe = iframeRef.current;
    const doc = iframe?.contentDocument || iframe?.contentWindow?.document;
    if (!doc) {
      return;
    }
    doc.getElementById(notebookIframeDarkBootstrapStyleId)?.remove();
    clearNotebookIframeDarkDocument(doc);
  }, []);

  const applyNotebookIframeDarkBootstrap = useCallback(() => {
    if (!isDarkThemeActive()) {
      clearNotebookIframeDarkBootstrap();
      return;
    }
    const iframe = iframeRef.current;
    const doc = iframe?.contentDocument || iframe?.contentWindow?.document;
    if (!iframe || !doc) {
      return;
    }

    let style = doc.getElementById(notebookIframeDarkBootstrapStyleId) as HTMLStyleElement | null;
    if (!style) {
      style = doc.createElement('style');
      style.id = notebookIframeDarkBootstrapStyleId;
      (doc.head || doc.documentElement).appendChild(style);
    }
    const themeCss = buildNotebookIframeDarkThemeCss();
    if (style.textContent !== themeCss) {
      style.textContent = themeCss;
    }
    syncNotebookIframeDarkDocument(doc);
  }, [clearNotebookIframeDarkBootstrap]);

  const applyNotebookIframeNameLimit = useCallback(() => {
    const iframe = iframeRef.current;
    const doc = iframe?.contentDocument || iframe?.contentWindow?.document;
    if (!iframe || !doc) {
      return;
    }

    let style = doc.getElementById(notebookIframeNameLimitStyleId) as HTMLStyleElement | null;
    if (!style) {
      style = doc.createElement('style');
      style.id = notebookIframeNameLimitStyleId;
      doc.head?.appendChild(style);
    }
    style.textContent = `
      .tier0-notebook-truncated-title {
        display: inline-block !important;
        max-width: ${notebookIframeNameLimitMaxWidth} !important;
        min-width: 0 !important;
        overflow: hidden !important;
        text-overflow: ellipsis !important;
        vertical-align: bottom !important;
        white-space: nowrap !important;
      }
    `;

    const matchesDisplayName = (text: string) =>
      notebookDisplayNames.some((name) => {
        return (
          text === name ||
          text.endsWith(`/${name}`) ||
          text.endsWith(`\\${name}`) ||
          (text.startsWith(name) && text.length <= name.length + 3)
        );
      });

    const isNonContentElement = (node: HTMLElement) => {
      const tagName = node.tagName.toLowerCase();
      return (
        ['marimo-filename', 'title', 'script', 'style', 'link', 'meta', 'svg'].includes(tagName) ||
        !!node.closest('head, marimo-filename, title, script, style, link, meta, svg')
      );
    };

    const applyTextOverflowStyle = (node: HTMLElement, forceWidth = false) => {
      node.style.setProperty('max-width', notebookIframeNameLimitMaxWidth, 'important');
      node.style.setProperty('min-width', '0', 'important');
      node.style.setProperty('overflow', 'hidden', 'important');
      node.style.setProperty('text-overflow', 'ellipsis', 'important');
      node.style.setProperty('white-space', 'nowrap', 'important');
      if (forceWidth) {
        node.style.setProperty('width', notebookIframeNameLimitMaxWidth, 'important');
      } else {
        node.style.setProperty('display', 'inline-block', 'important');
      }
    };

    const applyNameLimit = (node: HTMLElement, text: string) => {
      if (isNonContentElement(node)) {
        return;
      }
      if (node.closest('pre, code, textarea, button, [contenteditable="true"], .cm-editor, .cm-line')) {
        return;
      }
      node.classList.add('tier0-notebook-truncated-title');
      applyTextOverflowStyle(node);
      node.title = text;
    };

    const applyTitleAttributes = () => {
      if (!notebookDisplayNames.length) {
        return;
      }
      if (!doc.body) {
        return;
      }
      const filter =
        (iframe.contentWindow as unknown as { NodeFilter?: typeof NodeFilter } | null)?.NodeFilter || NodeFilter;
      const applyRoot = (root: Document | ShadowRoot) => {
        const nodes = Array.from(
          root.querySelectorAll<HTMLElement>(
            [
              'h1',
              'h2',
              'h3',
              'header span',
              'header div',
              '[role="heading"]',
              '[class*="file" i]',
              '[class*="filename" i]',
              '[class*="file-name" i]',
              '[class*="title" i]',
              '[class*="notebook-name" i]',
              '[class*="notebook-title" i]',
              '[aria-label*="file" i]',
              '[data-testid*="filename" i]',
              '[data-testid*="notebook" i]',
            ].join(',')
          )
        );
        nodes.forEach((node) => {
          const text = normalizeText(node.textContent);
          if (!text || text.length > 512) {
            return;
          }
          if (matchesDisplayName(text)) {
            applyNameLimit(node, text);
          }
        });

        const inputs = Array.from(root.querySelectorAll<HTMLInputElement>('input, textarea'));
        inputs.forEach((node) => {
          const text = normalizeText(node.value || node.placeholder);
          if (!text || text.length > 512 || !matchesDisplayName(text)) {
            return;
          }
          if (isNonContentElement(node)) {
            return;
          }
          applyTextOverflowStyle(node, true);
          node.title = text;
        });

        const textNodes: Text[] = [];
        const walker = doc.createTreeWalker(root, filter.SHOW_TEXT, {
          acceptNode(node) {
            const parent = node.parentElement;
            const text = normalizeText(node.textContent);
            if (!parent || isNonContentElement(parent) || !text || text.length > 512 || !matchesDisplayName(text)) {
              return filter.FILTER_REJECT;
            }
            return filter.FILTER_ACCEPT;
          },
        });
        let currentNode = walker.nextNode();
        while (currentNode) {
          textNodes.push(currentNode as Text);
          currentNode = walker.nextNode();
        }
        textNodes.forEach((textNode) => {
          const node = textNode.parentElement;
          const text = normalizeText(textNode.textContent);
          if (node && text) {
            applyNameLimit(node, text);
          }
        });

        Array.from(root.querySelectorAll<HTMLElement>('*')).forEach((node) => {
          if (node.shadowRoot) {
            applyRoot(node.shadowRoot);
          }
        });
      };

      applyRoot(doc);
    };

    applyTitleAttributes();

    const win = iframe.contentWindow as
      | (Window & {
          __tier0NotebookNameLimitObserver?: MutationObserver;
          __tier0NotebookNameLimitSignature?: string;
        })
      | null;
    const ObserverCtor =
      (win as unknown as { MutationObserver?: typeof MutationObserver } | null)?.MutationObserver || MutationObserver;
    const signature = notebookDisplayNames.join('\n');
    if (
      win &&
      ObserverCtor &&
      doc.body &&
      notebookDisplayNames.length &&
      win[notebookIframeNameLimitSignatureKey] !== signature
    ) {
      win[notebookIframeNameLimitObserverKey]?.disconnect();
      let scanFrame = 0;
      const observer = new ObserverCtor(() => {
        if (scanFrame) {
          return;
        }
        scanFrame = win.requestAnimationFrame(() => {
          scanFrame = 0;
          applyTitleAttributes();
        });
      });
      observer.observe(doc.body, { childList: true, subtree: true, characterData: true });
      win[notebookIframeNameLimitObserverKey] = observer;
      win[notebookIframeNameLimitSignatureKey] = signature;
    }
  }, [notebookDisplayNames]);

  useEffect(() => {
    if (!iframeLoading) {
      applyNotebookIframeNameLimit();
    }
  }, [applyNotebookIframeNameLimit, iframeLoading]);

  const handleIframeLoad = useCallback(() => {
    const iframe = iframeRef.current;
    const proxyError = detectNotebookIframeError(iframe);
    if (proxyError) {
      setIframeError(proxyError);
      setIframeLoading(false);
      return;
    }

    applyNotebookIframeDarkBootstrap();
    applyNotebookIframeNameLimit();
    window.requestAnimationFrame(() => {
      setIframeLoading(false);
    });
  }, [applyNotebookIframeDarkBootstrap, applyNotebookIframeNameLimit]);

  useEffect(() => {
    if (!iframeUrl) {
      return;
    }

    applyNotebookIframeDarkBootstrap();
    applyNotebookIframeNameLimit();

    const parentObserver = new MutationObserver(() => {
      applyNotebookIframeDarkBootstrap();
    });
    parentObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class'],
    });
    (window as Window & { __tier0NotebookParentThemeObserver?: MutationObserver })[
      notebookIframeParentThemeObserverKey
    ] = parentObserver;

    return () => {
      parentObserver.disconnect();
      delete (window as Window & { __tier0NotebookParentThemeObserver?: MutationObserver })[
        notebookIframeParentThemeObserverKey
      ];
      clearNotebookIframeDarkBootstrap();
    };
  }, [applyNotebookIframeDarkBootstrap, applyNotebookIframeNameLimit, clearNotebookIframeDarkBootstrap, iframeUrl]);

  return (
    <div className={styles.editorPage}>
      <ComContent
        hasBack={false}
        style={{ overflow: 'hidden' }}
        titleStyle={{ overflow: 'hidden' }}
        title={
          <div className={styles.editorTitleBar}>
            <ComBackButton onClick={handleBack} />
            <Tooltip title={detail?.name || formatMessage('Notebook.title', {}, 'Notebook')}>
              <span
                className={styles.editorTitle}
                title={detail?.name || formatMessage('Notebook.title', {}, 'Notebook')}
              >
                {detail?.name || formatMessage('Notebook.title', {}, 'Notebook')}
              </span>
            </Tooltip>
          </div>
        }
        extra={
          <Space>
            <Button
              icon={<Renew size={16} style={{ marginTop: 2 }} />}
              onClick={() => {
                setIframeError('');
                setIframeLoading(true);
                setIframeKey(Date.now());
              }}
            />
            <Button.Group>
              {canManage && (
                <Button icon={<Save size={16} style={{ marginTop: 2 }} />} onClick={() => setSaveVersionOpen(true)}>
                  {formatMessage('Notebook.snapshot.saveVersion', {}, 'Save Version')}
                </Button>
              )}
              <Popover
                open={popoverOpen}
                onOpenChange={setPopoverOpen}
                trigger="click"
                placement="bottomRight"
                content={
                  <HistoryPopover
                    notebookId={notebookId}
                    visible={popoverOpen}
                    refreshKey={historyRefreshKey}
                    onSelectItem={handlePopoverSelectItem}
                    onViewAll={handlePopoverViewAll}
                  />
                }
              >
                <Button icon={<ChevronDown size={16} />} />
              </Popover>
            </Button.Group>
          </Space>
        }
      >
        <div className={styles.editorShell}>
          <div className={styles.editorMain}>
            {(loading || iframeLoading) && (
              <div className={styles.loadingWrap}>
                <Spin />
              </div>
            )}
            {!loading && detail && (
              <div className={styles.editorBody}>
                {iframeError ? (
                  <div className={styles.iframeErrorPanel}>
                    <p className={styles.iframeErrorTitle}>
                      {formatMessage('Notebook.editorUnavailable', {}, 'Notebook editor is unavailable')}
                    </p>
                    <p className={styles.iframeErrorHint}>
                      {formatMessage(
                        'Notebook.editorUnavailableHint',
                        {},
                        'Please check that the Marimo service is running and reachable from the gateway.'
                      )}
                    </p>
                    <pre className={styles.iframeErrorDetail}>{iframeError}</pre>
                  </div>
                ) : iframeUrl ? (
                  <iframe
                    key={iframeKey}
                    ref={iframeRef}
                    className={`${styles.editorIframe} ${iframeLoading ? '' : styles.editorIframeReady}`}
                    title={detail.name || 'Notebook'}
                    src={iframeUrl}
                    onLoad={handleIframeLoad}
                  />
                ) : null}
              </div>
            )}
          </div>
          {historyOpen && (
            <div className={styles.historyPanel}>
              <HistoryPanel
                key={historyRefreshKey}
                notebookId={notebookId}
                selectedSnapshotId={selectedSnapshotId}
                onClose={() => {
                  setHistoryOpen(false);
                  setSelectedSnapshotId(null);
                }}
                onReverted={handleReverted}
              />
            </div>
          )}
        </div>
      </ComContent>

      <SaveVersionModal
        open={saveVersionOpen}
        confirmLoading={saveVersionLoading}
        latestVersionName={latestVersionName}
        onCancel={() => setSaveVersionOpen(false)}
        onSubmit={async (versionName, description) => {
          const saved = await handleSaveVersion(versionName, description);
          if (saved) {
            void loadLatestVersion();
            setHistoryRefreshKey((prev) => prev + 1);
          }
        }}
      />
    </div>
  );
};

export default NotebookEditorPage;
