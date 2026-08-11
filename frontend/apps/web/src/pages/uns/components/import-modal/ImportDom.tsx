import cx from 'classnames';
import { CheckmarkFilled, ErrorFilled, Upload as UploadIcon } from '@/components/lucide-icon/carbon';
import ComEllipsis from '../../../../components/com-ellipsis';
import { App, Divider, Flex, Progress, Upload, type UploadFile } from 'antd';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useTranslate } from '@/hooks';
import ComButton from '../../../../components/com-button';
import type { SocketDataType } from './type.ts';
import InlineLoading from '@/components/inline-loading/index.tsx';
import { readerSSE } from './utils.ts';
import { getToken } from '@/utils/auth.ts';
const { Dragger } = Upload;

const ImportDom = ({ initTreeData, onCloseModal, fillHeight = false }: any) => {
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const uploadRef = useRef<any>(null);
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const [moduleMap, setModuleMap] = useState(new Map());
  const timer = useRef<number>();

  const showImportError = (error?: unknown) => {
    const msg = error instanceof Error && error.message ? error.message : formatMessage('uns.importFailed');
    const item = {
      code: 500,
      finished: true,
      progress: 100,
      module: 'uns',
      msg,
    };
    setModuleMap((prevMap) => {
      const newMap = new Map(prevMap);
      newMap.set(item.module, item);
      return newMap;
    });
    setSocketData(item);
    setLoading(false);
  };

  const beforeUpload = (file: any) => {
    const fileType = (file.name.split('.').pop() || '').toLowerCase();
    if (fileType === 'json') {
      setFileList([file]);
    } else {
      message.warning(formatMessage('common.theFileFormatType', { fileType: '.json' }));
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
          const item = {
            module: data?.module || 'uns',
            ...data,
            finished: data?.finished || data?.progress >= 100,
          };
          setModuleMap((prevMap) => {
            const newMap = new Map(prevMap);
            newMap.set(item.module, item);
            return newMap;
          });
          setSocketData({
            code: data?.code,
            finished: item.finished,
            progress: data?.progress,
            errTipFile: data?.errTipFile,
            module: item.module,
            msg: data?.msg,
          });
          if (item.finished && [200, 206].includes(item.code)) initTreeData({ reset: true });
        },
        () => {
          showImportError();
        }
      );
    } catch (error) {
      console.error(error);
      showImportError(error);
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
    <Flex vertical style={{ height: fillHeight ? '100%' : undefined, overflow: 'hidden' }}>
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
                <span style={{ color: 'var(--ui-text-color)' }}>{formatMessage('uns.uploadSuccess')}</span>
              </Flex>
              {!socketData?.finished && formatMessage('uns.waitingFormParsing')}
            </Flex>
          ) : (
            <div className="upload-drag-content">
              <UploadIcon size={32} className="upload-drag-icon" />
              <p className="upload-hint-primary">{formatMessage('common.clickOrDragForUpload')}</p>
            </div>
          )}
        </Dragger>
      </div>
      {loading ? (
        <div style={{ flex: fillHeight ? 1 : undefined, minHeight: fillHeight ? 0 : undefined, overflow: 'auto', paddingTop: 16 }}>
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
                      <span>{item.msg || item.module}</span>
                    </Flex>
                  }
                />
              );
            })}
          </div>
        </div>
      ) : null}
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
