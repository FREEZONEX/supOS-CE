import { App, Empty, Flex, Image, Tooltip } from 'antd';
import { forwardRef, useEffect, useImperativeHandle, useState } from 'react';
import { deleteAttachments, getAttachment, getAttachmentsList } from '@/apis/inter-api/attachments.ts';
import { Close, Download, View } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import ComCopy from '@/components/com-copy';
import ComItem from '@/components/com-item';
import { downloadFn } from '@/utils/blob';
import { validPicRegex } from '@/utils/pattern';
export interface DocumentListRef {
  refresh: () => any;
}
const DocumentList = forwardRef<DocumentListRef | undefined, { alias: string }>(function ({ alias }, ref) {
  const formatMessage = useTranslate();
  const [visible, setVisible] = useState(false);
  const [imageUrl, setImageUrl] = useState('');
  const { message } = App.useApp();
  const [data, setData] = useState<any[]>([]);
  const request = () => {
    if (!alias) return;
    getAttachmentsList({ alias }).then((data: any) => {
      setData(
        data?.list?.map((item: any) => {
          const lastDotIndex = item.originalName.lastIndexOf('.');
          const _label = item.originalName.slice(0, lastDotIndex);
          const _type = item.originalName.slice(lastDotIndex + 1);
          return {
            ...item,
            _label,
            _type: _type.toUpperCase(),
            isPic: validPicRegex.test(_type),
          };
        })
      );
    });
  };
  useImperativeHandle(ref, () => {
    return {
      refresh: () => request(),
    };
  });
  useEffect(() => {
    request();
  }, [alias]);

  const onDelete = (item: any) => {
    deleteAttachments({ objectName: item?.attachmentPath }).then(() => {
      message.success(formatMessage('common.deleteSuccessfully'));
      request();
    });
  };

  const onDownload = (item: any) => {
    getAttachment({ objectName: item?.attachmentPath }).then((data: any) => {
      downloadFn({ data, name: item.originalName });
    });
  };

  const onPreview = (item: any) => {
    getAttachment({ objectName: item?.attachmentPath }).then((data: any) => {
      const url = window.URL.createObjectURL(new Blob([data]));
      setImageUrl(url);
      setVisible(true);
    });
  };
  return data?.length > 0 ? (
    <>
      <Flex gap={6} wrap>
        {data?.map((item) => (
          <ComItem
            key={item.id}
            addonBefore={<span style={{ color: 'var(--supos-theme-color)' }}>{item._type}</span>}
            label={item._label}
            style={{ flex: '0 0 calc(50% - 6px)' }}
            extra={
              <Flex justify="end" align="center" gap={10} style={{ cursor: 'pointer' }}>
                {item.isPic && (
                  <Tooltip title={formatMessage('dashboards.preview')}>
                    <View onClick={() => onPreview(item)} />
                  </Tooltip>
                )}
                {item.isPic && (
                  <ComCopy
                    textToCopy={`${location.origin}/inter-api/supos/uns/attachment?objectName=${item?.attachmentPath}`}
                    title={formatMessage('common.copy')}
                  />
                )}
                <Tooltip title={formatMessage('common.download')}>
                  <Download size={16} onClick={() => onDownload(item)} />
                </Tooltip>
                <Tooltip title={formatMessage('common.delete')}>
                  <Close size={16} onClick={() => onDelete(item)} />
                </Tooltip>
              </Flex>
            }
          />
        ))}
      </Flex>
      <Image
        width={200}
        style={{ display: 'none' }}
        src={imageUrl}
        preview={{
          visible,
          src: imageUrl,
          onVisibleChange: (value) => {
            if (!value) {
              setImageUrl('');
            }
            setVisible(value);
          },
        }}
      />
    </>
  ) : (
    <Empty />
  );
});

export default DocumentList;
