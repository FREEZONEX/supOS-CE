import cx from 'classnames';
import { CheckmarkFilled, Download, ErrorFilled, FolderAdd } from '@carbon/icons-react';
import ComEllipsis from '../../../../components/com-ellipsis';
import { App, Button, Divider, Flex, Progress, Upload, type UploadFile } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { useTranslate } from '@/hooks';
import ComButton from '../../../../components/com-button';
import type { SocketDataType } from './type.ts';
import InlineLoading from '@/components/inline-loading/index.tsx';
import { readerSSE } from '@/pages/uns/components/import-modal/utils.ts';
const { Dragger } = Upload;

const ImportDom = ({ initTreeData, onCloseModal }: any) => {
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const uploadRef = useRef<any>(null);
  const { message, modal } = App.useApp();
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const [moduleMap, setModuleMap] = useState(new Map());
  const timer = useRef<number>();

  const beforeUpload = (file: any) => {
    const fileType = file.name.split('.').pop();
    if (['json', 'zip'].includes(fileType.toLowerCase())) {
      setFileList([file]);
    } else {
      message.warning(formatMessage('common.theFileFormatType', { fileType: '.json,.zip' }));
    }
    return false;
  };

  const { code, finished, progress } = socketData;
  const reimport = finished;

  const onSave = async () => {
    try {
      const fd = new FormData();
      if (fileList.length) {
        fd.append('file', fileList[0] as any, fileList[0].name);
      } else {
        message.warning(formatMessage('uns.pleaseUploadTheFile'));
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
    setFileList([]);
    setModuleMap(new Map());
    setTimeout(() => {
      if (uploadRef.current) uploadRef?.current?.nativeElement?.querySelector('input').click();
    });
  };
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
      <div style={{ flexShrink: 0 }}>
        <Dragger
          ref={uploadRef}
          className={cx('import-upload', fileList?.length > 0 && 'upload-file')}
          action=""
          accept=".json"
          maxCount={1}
          disabled={loading}
          beforeUpload={beforeUpload}
          fileList={fileList}
          onRemove={() => {
            setFileList([]);
          }}
        >
          {loading ? (
            <Flex vertical align="center" justify="center" gap={8}>
              <Flex align="center" gap={4}>
                <CheckmarkFilled fill={'#24a148'} />
                <span style={{ color: 'var(--supos-text-color)' }}>{formatMessage('uns.uploadSuccess')}</span>
              </Flex>
              {!socketData?.finished && formatMessage('uns.waitingFormParsing')}
            </Flex>
          ) : (
            <div style={{ height: 170 }}>
              <FolderAdd size={48} style={{ color: '#E0E0E0' }} />
              <ComEllipsis style={{ padding: '16px 0' }}>{formatMessage('common.clickOrDragForUpload')}</ComEllipsis>
              <Button
                onClick={(e) => {
                  e.stopPropagation();
                  window.open(`/inter-api/supos/uns/importExport/template/download?fileType=json`, '_self');
                }}
              >
                <Download />
                {formatMessage('common.downloadTemplate')}
              </Button>
            </div>
          )}
        </Dragger>
      </div>
      <div style={{ flex: 1, overflow: 'auto', paddingTop: 16 }}>
        {loading && (
          <div>
            <ComEllipsis style={{ color: '#525252' }}>{formatMessage('uns.overallProgress')}</ComEllipsis>
            <Flex align="center" gap={8}>
              <Progress percent={progress} showInfo={false} />
              <div>{`${progress || 0}%`}</div>
              {finished ? code === 200 ? <CheckmarkFilled fill={'#24a148'} /> : <ErrorFilled fill={'#da1e28'} /> : null}
            </Flex>
            <ComEllipsis style={{ color: '#525252' }}>{formatMessage('uns.individualProgress')}</ComEllipsis>
            {Array.from(moduleMap.values()).map((item: any) => {
              return (
                <InlineLoading
                  key={item.module}
                  title={item.module}
                  style={{ width: '100%' }}
                  status={item?.finished ? (item?.code === 200 ? 'finished' : 'error') : 'active'}
                  description={
                    <Flex justify="space-between">
                      <span>{item.module}</span>
                    </Flex>
                  }
                />
              );
            })}
          </div>
        )}
      </div>
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

export default ImportDom;
