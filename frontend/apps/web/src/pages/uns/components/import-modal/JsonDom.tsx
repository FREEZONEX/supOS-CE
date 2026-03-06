import CodeMirror from '@uiw/react-codemirror';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { json, jsonParseLinter } from '@codemirror/lang-json';
import { linter, lintGutter } from '@codemirror/lint';
import { placeholder } from '@/pages/uns/components/import-modal/data.ts';
import { useEffect, useRef, useState } from 'react';
import { useSize } from 'ahooks';
import styles from '@/theme/codemirror.module.scss';
import { CheckmarkFilled, Copy, ErrorFilled } from '@carbon/icons-react';
import { useClipboard, useTranslate } from '@/hooks';
import { App, Divider, Flex, Progress } from 'antd';
import ComButton from '../../../../components/com-button';
import ComEllipsis from '../../../../components/com-ellipsis';
import type { SocketDataType } from './type.ts';
import { readerSSE } from '@/pages/uns/components/import-modal/utils.ts';

const JsonDom = ({ initTreeData, onCloseModal }: any) => {
  const [jsonValue, setJsonValue] = useState<any>();
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);
  const formatMessage = useTranslate();
  const { message, modal } = App.useApp();
  const { copy } = useClipboard();
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const [moduleMap, setModuleMap] = useState(new Map());
  console.log(moduleMap);
  const [loading, setLoading] = useState(false);
  const timer = useRef<number>();

  const onSave = async () => {
    try {
      const fd = new FormData();
      if (jsonValue) {
        try {
          JSON.parse(jsonValue);
        } catch (e) {
          console.log(e);
          message.error(formatMessage('uns.errorInTheSyntaxOfTheJSON'));
          return;
        }
        fd.append('file', new Blob([jsonValue], { type: 'application/json' }), 'uns.json');
      } else {
        message.warning(formatMessage('uns.pleaseJSON'));
        return;
      }
      setLoading(true);
      const response = await fetch('/inter-api/supos/uns/importExport/import', {
        method: 'POST',
        body: fd,
      });
      readerSSE(
        response,
        (data: any) => {
          setModuleMap((prevMap) => {
            const newMap = new Map(prevMap);
            newMap.set(data.module, data);
            return newMap;
          });
          setSocketData({
            code: data?.code,
            finished: data?.progress >= 100,
            progress: data?.progress,
            errTipFile: data?.errTipFile,
          });
          if (data?.progress >= 100) initTreeData({ reset: true });
        },
        () => {
          setLoading(false);
        }
      );
    } catch (error) {
      console.error(error);
      setLoading(false);
    }
  };

  const onReupload = () => {
    setLoading(false);
    setSocketData({});
    setModuleMap(new Map());
    setJsonValue(undefined);
  };

  const { code, finished, progress } = socketData;
  const reimport = finished;
  const onClose = () => {
    onCloseModal?.();
  };

  useEffect(() => {
    if (socketData.finished) {
      clearInterval(timer.current);
      if (socketData.code === 200) {
        message.success(formatMessage('uns.importFinished'));
        // setTimeout(() => {
        //   onClose();
        // }, 3000);
      }
    }
    if (socketData.code === 206) {
      modal.confirm({
        title: formatMessage('uns.PartialDataImportFailed'),
        onOk() {
          window.open(`/inter-api/supos/uns/importExport/file/download?path=${socketData.errTipFile}`, '_self');
        },
        okButtonProps: {
          title: formatMessage('common.confirm'),
        },
        cancelButtonProps: {
          title: formatMessage('common.cancel'),
        },
      });
    }
  }, [socketData]);

  return (
    <Flex vertical style={{ height: '100%', overflow: 'hidden' }}>
      <div
        ref={ref}
        style={{
          flex: 1,
          borderRadius: 4,
          border: '1px solid rgb(198, 198, 198)',
          padding: 16,
          position: 'relative',
          overflow: 'hidden',
        }}
        className={styles['custom-theme']}
      >
        <div
          style={{
            position: 'absolute',
            right: 4,
            top: 0,
            color: 'var(--supos-text-color)',
            zIndex: 1,
          }}
        >
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
            <span
              style={{
                marginRight: 14,
                fontSize: '12px',
                pointerEvents: 'none',
                zIndex: 10,
                color: '#c6c6c6',
              }}
            >
              {formatMessage('uns.ctrlPQuickApplyExample')}
            </span>
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
