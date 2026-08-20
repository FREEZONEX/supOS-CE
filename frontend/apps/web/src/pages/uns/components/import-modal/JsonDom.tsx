import CodeMirror from '@uiw/react-codemirror';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { json, jsonParseLinter } from '@codemirror/lang-json';
import { linter, lintGutter } from '@codemirror/lint';
import { placeholder } from '@/pages/uns/components/import-modal/data.ts';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useSize } from 'ahooks';
import styles from '@/theme/codemirror.module.scss';
import { CheckmarkFilled, Copy, ErrorFilled } from '@/components/lucide-icon/carbon';
import { useClipboard, useTranslate } from '@/hooks';
import { App, Divider, Flex, Progress } from 'antd';
import ComButton from '../../../../components/com-button';
import ComEllipsis from '../../../../components/com-ellipsis';
import type { SocketDataType } from './type.ts';
import { readerSSE } from './utils.ts';
import { getToken } from '@/utils/auth.ts';

const JsonDom = ({ initTreeData, onCloseModal }: any) => {
  const [jsonValue, setJsonValue] = useState<any>();
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const { copy } = useClipboard();
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const [loading, setLoading] = useState(false);
  const timer = useRef<number>();

  const showImportError = (error?: unknown) => {
    const msg = error instanceof Error && error.message ? error.message : formatMessage('uns.importFailed');
    setSocketData({
      code: 500,
      finished: true,
      progress: 100,
      module: 'uns',
      msg,
    });
    setLoading(false);
  };

  const onSave = async () => {
    try {
      const fd = new FormData();
      if (jsonValue) {
        try {
          JSON.parse(jsonValue);
        } catch {
          message.error(formatMessage('uns.errorInTheSyntaxOfTheJSON'));
          return;
        }
        fd.append('file', new Blob([jsonValue], { type: 'application/json' }), 'uns.json');
      } else {
        message.warning(formatMessage('uns.pleaseJSON'));
        return;
      }
      setLoading(true);
      const response = await fetch('/api/core/uns/import', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${getToken() || ''}`,
        },
        body: fd,
      });
      if (!response.ok) {
        const msg = await response.text();
        throw new Error(msg || response.statusText || formatMessage('uns.importFailed'));
      }
      readerSSE(
        response,
        (data: any) => {
          const finished = data?.finished || data?.progress >= 100;
          setSocketData({
            code: data?.code,
            finished,
            progress: data?.progress,
            errTipFile: data?.errTipFile,
            module: data?.module || 'uns',
            msg: data?.msg,
          });
          if (finished && [200, 206].includes(data?.code)) initTreeData({ reset: true });
        },
        () => {
          showImportError();
        }
      );
    } catch (error) {
      showImportError(error);
    }
  };

  const onReupload = () => {
    setLoading(false);
    setSocketData({});
  };

  const { code, finished, progress } = socketData;
  const reimport = finished;
  const onClose = useCallback(() => {
    onCloseModal?.();
  }, [onCloseModal]);

  useEffect(() => {
    if (socketData.finished) {
      clearInterval(timer.current);
      if (socketData.code === 200) {
        message.success(formatMessage('uns.importFinished'));
        onClose();
      } else if (socketData.code !== 206) {
        message.error(socketData.msg || formatMessage('uns.importFailed'));
      }
    }
    if (socketData.code === 206) {
      message.warning(socketData.msg || formatMessage('uns.PartialDataImportFailed'));
    }
  }, [formatMessage, message, onClose, socketData]);

  return (
    <Flex vertical style={{ height: '100%', overflow: 'hidden' }}>
      <div className="import-json-layout">
        <div ref={ref} className={`${styles['custom-theme']} import-json-editor`}>
          <div className={`import-json-quick-tip${jsonValue ? ' is-copy' : ''}`}>
            {jsonValue ? (
              <Copy
                style={{
                  cursor: 'pointer',
                  marginTop: 4,
                }}
                onClick={() => {
                  copy(jsonValue || JSON.stringify(JSON.parse(placeholder), null, 2));
                }}
              />
            ) : (
              <span>{formatMessage('uns.ctrlPQuickApplyExample')}</span>
            )}
          </div>
          <CodeMirror
            theme={codemirrorTheme}
            placeholder={placeholder}
            onChange={setJsonValue}
            value={jsonValue}
            height={(size?.height || 32) - 32 + 'px'}
            extensions={[json(), linter(jsonParseLinter()), lintGutter()]}
            onKeyDownCapture={(e) => {
              if (e.ctrlKey && e.key === 'Enter') {
                e.preventDefault();
                e.stopPropagation();
                setJsonValue(placeholder);
              }
            }}
            onKeyDown={(e) => {
              if (e.ctrlKey && e.key === 'Enter') {
                e.preventDefault();
              }
            }}
          />
        </div>
      </div>
      {loading && (
        <div style={{ flexShrink: 0, paddingTop: 16 }}>
          <ComEllipsis style={{ color: '#525252' }}>{formatMessage('uns.overallProgress')}</ComEllipsis>
          <Flex align="center" gap={8}>
            <Progress percent={progress} showInfo={false} />
            <div>{`${progress || 0}%`}</div>
            {finished ? code === 200 ? <CheckmarkFilled fill={'#24a148'} /> : <ErrorFilled fill={'#da1e28'} /> : null}
          </Flex>
        </div>
      )}
      <div style={{ flexShrink: 0 }}>
        <Divider style={{ backgroundColor: 'rgb(198, 198, 198)', margin: '16px 0' }} />
        <Flex align="center" gap={8} justify="flex-end">
          <ComButton onClick={onClose}>{formatMessage('common.cancel')}</ComButton>
          <ComButton
            type="primary"
            onClick={reimport ? onReupload : onSave}
            loading={reimport ? false : loading}
            disabled={reimport ? false : loading}
          >
            {formatMessage(reimport ? 'uns.reimport' : 'common.save')}
          </ComButton>
        </Flex>
      </div>
    </Flex>
  );
};

export default JsonDom;
