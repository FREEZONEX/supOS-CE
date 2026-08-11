import CodeMirror from '@uiw/react-codemirror';
import { Alert, Flex } from 'antd';
import { useClipboard, useTranslate } from '@/hooks';
import { useSize } from 'ahooks';
import { useRef } from 'react';
import { json } from '@codemirror/lang-json';
import { Copy, Download } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { downloadFn } from '@/utils/blob';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { exportExcel } from '@/apis/core-api';
import ComButton from '@/components/com-button';
import styles from './index.module.scss';

export const CodeDom = () => {
  const formatMessage = useTranslate();
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);

  const { smallFile, jsonData, params } = useTreeStore((state) => ({
    smallFile: state.smallFile,
    jsonData: state.jsonData,
    params: state.params,
  }));

  const { copy } = useClipboard();

  return (
    <>
      <div style={{ flex: 1, overflow: 'hidden' }}>
        {!smallFile ? (
          <>
            <div>
              <ComButton
                type="primary"
                icon={<Download size={16} />}
                onClick={() => {
                  return exportExcel(params).then((jsonData) => {
                    downloadFn({ data: JSON.stringify(jsonData), name: 'uns.json' });
                  });
                }}
              >
                {formatMessage('common.download')}
              </ComButton>
            </div>
            <Alert
              style={{ marginTop: 16, padding: '8px 16px' }}
              description={formatMessage('uns.exportJsonWarning')}
              type="warning"
              showIcon
            />
          </>
        ) : (
          <div className={styles.exportCodePanel} ref={ref}>
            <div
              style={{
                position: 'absolute',
                right: 8,
                top: 8,
                color: 'var(--ui-text-color)',
                zIndex: 1,
                cursor: 'pointer',
                display: 'flex',
                alignItems: 'center',
              }}
              onClick={() => {
                copy(jsonData);
              }}
            >
              <Copy {...toolbarIconProps} />
            </div>
            <CodeMirror
              // onChange={setJsonValue}
              theme={codemirrorTheme}
              value={jsonData}
              editable={false}
              height={(size?.height || 32) - 32 + 'px'}
              extensions={[json()]}
              placeholder={formatMessage('uns.pleaseSelectForExport')}
            />
          </div>
        )}
      </div>
      <Flex className={styles.exportPanelActions} justify="end" gap={8} style={{ marginTop: 16 }}>
        {
          <ComButton
            type="primary"
            icon={<Download size={16} />}
            onClick={() => {
              return exportExcel(params).then((jsonData) => {
                downloadFn({ data: JSON.stringify(jsonData), name: 'uns.json' });
              });
            }}
            disabled={!(jsonData && smallFile)}
          >
            {formatMessage('common.download')}
          </ComButton>
        }
        <ComButton
          icon={<Copy size={16} />}
          onClick={() => {
            copy(jsonData);
          }}
          disabled={!(jsonData && smallFile)}
        >
          {formatMessage('common.copy')}
        </ComButton>
      </Flex>
    </>
  );
};
