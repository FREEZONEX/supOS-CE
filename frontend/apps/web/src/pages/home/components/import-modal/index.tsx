import { type FC, useState, useEffect, useRef } from 'react';
import { FolderAdd } from '@carbon/icons-react';
import { Upload, Button, App, Flex } from 'antd';
import { useTranslate } from '@/hooks';

import { useWebSocket } from 'ahooks';

import type { Dispatch, SetStateAction } from 'react';
import type { UploadFile } from 'antd';

import './index.scss';
import InlineLoading from '@/components/inline-loading';
import ProModal from '@/components/pro-modal';
import { importGlobal } from '@/apis/inter-api/global.ts';
import { getToken } from '@/utils/auth';

const { Dragger } = Upload;

export interface ImportModalProps {
  importModal: boolean;
  setImportModal: Dispatch<SetStateAction<boolean>>;
}

interface SocketDataType {
  code?: number;
  finished?: boolean;
  msg?: string;
  progress?: number;
  task?: string;
  errTipFile?: string;
  module?: string;
  runningStatusList?: SocketDataType[];
  totalCount?: number;
  errorCount?: number;
  successCount?: number;
}

const Module: FC<ImportModalProps> = (props) => {
  const { importModal, setImportModal } = props;
  const formatMessage = useTranslate();
  const { message, modal } = App.useApp();
  const timer = useRef<number>(undefined);
  const uploadRef = useRef<any>(null);

  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const [socketUrl, setSocketUrl] = useState('');

  // 创建 WebSocket 连接
  const { readyState, disconnect, sendMessage } = useWebSocket(
    socketUrl, // 初始 URL 为 null，表示不立即连接
    {
      reconnectLimit: 0,
      onMessage: (event) => {
        if (event.data === 'pong') return;
        const data = JSON.parse(event.data);
        setSocketData(data);
      },
      onError: (error) => console.error('WebSocket error:', error),
    }
  );

  const beforeUpload = (file: any) => {
    const fileType = file.name.split('.').pop();
    if (['zip'].includes(fileType.toLowerCase())) {
      setFileList([file]);
    } else {
      message.warning(formatMessage('uns.theFileFormatOnlySupportsZip'));
    }
    return false;
  };

  const save = () => {
    if (fileList.length) {
      setLoading(true);
      importGlobal({
        value: fileList[0],
        name: 'file',
        fileName: fileList[0].name,
      })
        .then((data) => {
          if (data) {
            const protocol = location.protocol.includes('https') ? 'wss' : 'ws';
            // 创建 WebSocket 连接
            setSocketUrl(
              `${protocol}://${location.host}/inter-api/supos/uns/ws?file=${encodeURIComponent(data)}&global=${String(Date.now())}&token=${getToken()}`
            );
          }
        })
        .catch(() => {
          resetUploadStatus();
        });
    } else {
      message.warning(formatMessage('uns.pleaseUploadTheFile'));
    }
  };

  const close = () => {
    setImportModal(false);
    setFileList([]);
    if (socketData.finished) {
      resetUploadStatus();
    }
  };

  const resetUploadStatus = () => {
    setLoading(false);
    setSocketUrl('');
    setSocketData({});
    setFileList([]);
  };

  const Reupload = () => {
    resetUploadStatus();
    setTimeout(() => {
      if (uploadRef.current) uploadRef?.current?.nativeElement?.querySelector('input').click();
    });
  };

  useEffect(() => {
    if (socketData.finished && disconnect) {
      disconnect();
      clearInterval(timer.current);
      if (socketData.code === 200) {
        // message.success(formatMessage('uns.importFinished'));
        setTimeout(() => {
          close();
        }, 3000);
      }
    }
    if (socketData.code === 206) {
      modal.confirm({
        title: formatMessage('uns.PartialDataImportFailed'),
        onOk() {
          window.open(`/inter-api/supos/global/file/download?path=${socketData.errTipFile}`, '_self');
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

  useEffect(() => {
    if (readyState === 1) {
      timer.current = window.setInterval(() => {
        if (sendMessage && readyState === 1) {
          sendMessage('ping');
        } else {
          clearInterval(timer.current);
        }
      }, 30000);
    }
    return () => {
      clearInterval(timer.current);
    };
  }, [readyState]);

  useEffect(() => {
    return () => {
      if (disconnect) {
        disconnect();
      }
      clearInterval(timer.current);
    };
  }, []);

  const { code, finished, msg, task, runningStatusList } = socketData;

  const reimport = finished && code !== 200;

  return (
    <ProModal
      className="importModalGlobalWrap"
      open={importModal}
      onCancel={close}
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span>{formatMessage('common.import')}</span>
        </div>
      }
      width={460}
      maskClosable={false}
      keyboard={false}
    >
      {loading ? (
        <Flex vertical className="loadingGlobalWrap" justify="center" align="center">
          <div style={{ maxWidth: '90%' }}>
            <InlineLoading
              status={finished ? (code === 200 ? 'finished' : 'error') : 'active'}
              description={`${finished ? msg : task || ''}`}
            />
          </div>
          {runningStatusList?.length && (
            <div className="loadingGlobalContent">
              {runningStatusList?.map((m) => {
                const { code, finished, msg, task, module, totalCount, successCount } = m;
                const title = `${formatMessage('home.' + module)}${module === 'uns' && finished ? '-' + task : ''}：${finished ? msg : task || ''}`;
                return (
                  <InlineLoading
                    title={title}
                    style={{ width: '100%' }}
                    status={finished ? (code === 200 ? 'finished' : 'error') : 'active'}
                    description={
                      <Flex justify="space-between">
                        <span>{title}</span>
                        <span>
                          {code === 200 && !totalCount ? null : (
                            <>
                              <span style={{ color: '#6FDC8C' }}>{successCount ?? 0}</span>
                              <span>/{totalCount ?? 0}</span>
                            </>
                          )}
                        </span>
                      </Flex>
                    }
                  />
                );
              })}
            </div>
          )}
        </Flex>
      ) : (
        <>
          <Dragger
            ref={uploadRef}
            className="uploadGlobalWrap"
            action=""
            accept=".zip"
            maxCount={1}
            beforeUpload={beforeUpload}
            fileList={fileList}
            onRemove={() => {
              setFileList([]);
            }}
          >
            <FolderAdd size={100} style={{ color: '#E0E0E0' }} />
          </Dragger>
        </>
      )}
      <Button
        color="primary"
        variant="solid"
        onClick={reimport ? Reupload : save}
        block
        style={{ marginTop: '10px' }}
        loading={reimport ? false : loading}
        disabled={reimport ? false : loading}
      >
        {formatMessage(reimport ? 'uns.reimport' : 'common.save')}
      </Button>
    </ProModal>
  );
};
export default Module;
