import CodeMirror from '@uiw/react-codemirror';
import { Button, Flex } from 'antd';
import { useClipboard, useTranslate } from '@/hooks';
import { useSize } from 'ahooks';
import { useEffect, useRef, useState } from 'react';
import { json } from '@codemirror/lang-json';
import { Copy, Download } from '@carbon/icons-react';
import { useTreeStore } from '@/pages/uns/components/export-modal/treeStore.tsx';
import { blobToJsonUsingTextMethod, downloadFn } from '@/utils/blob';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';

export const CodeDom = () => {
  const formatMessage = useTranslate();
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);
  const [jsonValue, setJsonValue] = useState<any>();
  const [isDownload, setIsDownload] = useState<boolean>(false);

  const { jsonData } = useTreeStore((state) => ({
    jsonData: state.jsonData,
  }));

  useEffect(() => {
    if (jsonData) {
      setIsDownload(jsonData?.size / 1024 / 1024 > 5);
      blobToJsonUsingTextMethod(jsonData).then((data) => {
        setJsonValue(data);
      });
    }
  }, [jsonData]);
  const { copy } = useClipboard();

  return (
    <>
      <div style={{ flex: 1, overflow: 'hidden' }}>
        {isDownload ? (
          <Button
            type="primary"
            onClick={() => {
              downloadFn({ data: jsonData, name: 'uns.json' });
            }}
          >
            <Download />
            {formatMessage('common.download')}
          </Button>
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
                copy(jsonValue);
              }}
            >
              <Copy />
            </div>
            <CodeMirror
              // onChange={setJsonValue}
              theme={codemirrorTheme}
              value={jsonValue}
              editable={false}
              height={(size?.height || 32) - 32 + 'px'}
              extensions={[json()]}
              placeholder={formatMessage('uns.pleaseSelectForExport')}
            />
          </div>
        )}
      </div>
      <Flex justify="end" gap={8} style={{ marginTop: 16 }}>
        <Button
          type="primary"
          onClick={() => {
            copy(jsonValue);
          }}
        >
          <Copy />
          {formatMessage('common.copy')}
        </Button>
      </Flex>
    </>
  );
};
