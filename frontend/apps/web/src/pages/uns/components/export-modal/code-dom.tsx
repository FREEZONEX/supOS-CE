import CodeMirror from '@uiw/react-codemirror';
import { Button, Flex } from 'antd';
import { useClipboard, useTranslate } from '@/hooks';
import { useSize } from 'ahooks';
import { useRef } from 'react';
import { json } from '@codemirror/lang-json';
import { Copy, Download } from '@carbon/icons-react';
import { useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { downloadFn } from '@/utils/blob';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import { exportExcel } from '@/apis/inter-api';
import ComButton from '@/components/com-button';

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
          <ComButton
            type="primary"
            onClick={() => {
              return exportExcel(params).then((jsonData) => {
                downloadFn({ data: JSON.stringify(jsonData), name: 'uns.json' });
              });
            }}
          >
            <Download />
            {formatMessage('common.download')}
          </ComButton>
        ) : (
          <div
            style={{
              height: '100%',
              borderRadius: 4,
              border: '1px solid rgb(198, 198, 198)',
              padding: 16,
              position: 'relative',
            }}
            ref={ref}
          >
            <div
              style={{
                position: 'absolute',
                right: 4,
                top: 4,
                color: 'var(--supos-text-color)',
                zIndex: 1,
                cursor: 'pointer',
              }}
              onClick={() => {
                copy(jsonData);
              }}
            >
              <Copy />
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
      <Flex justify="end" gap={8} style={{ marginTop: 16 }}>
        {jsonData && smallFile && (
          <ComButton
            type="primary"
            onClick={() => {
              return exportExcel(params).then((jsonData) => {
                downloadFn({ data: JSON.stringify(jsonData), name: 'uns.json' });
              });
            }}
          >
            <Download />
            {formatMessage('common.download')}
          </ComButton>
        )}
        <Button
          type="primary"
          onClick={() => {
            copy(jsonData);
          }}
        >
          <Copy />
          {formatMessage('common.copy')}
        </Button>
      </Flex>
    </>
  );
};
