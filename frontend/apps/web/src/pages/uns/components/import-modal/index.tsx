import { type FC, useEffect, useImperativeHandle, useRef, useState } from 'react';
import { Upload, App, Flex, Button, type UploadFile } from 'antd';
import { useClipboard, useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import ComRadio from '@/components/com-radio';
import ComEllipsis from '@/components/com-ellipsis';
import ComButton from '@/components/com-button';
import { importExcel } from '@/apis/inter-api';
import CodeMirror from '@uiw/react-codemirror';
import { json } from '@codemirror/lang-json';
import { useSize, useWebSocket } from 'ahooks';
import { Copy, Download, FolderAdd } from '@carbon/icons-react';
import cx from 'classnames';
import { getToken } from '@/utils';
import InlineLoading from '@/components/inline-loading';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import styles from '@/theme/codemirror.module.scss';
import './index.scss';

const { Dragger } = Upload;

export interface ImportModalProps {
  initTreeData: any;
  importRef: any;
}

interface SocketDataType {
  code?: number;
  finished?: boolean;
  msg?: string;
  progress?: number;
  task?: string;
  errTipFile?: string;
}
const placeholder = `{
  "Template": [],
  "Label": [],
  "UNS": [
    {
      "name": "v1",
      "type": "folder",
      "children": [
        {
          "name": "Plant_Name",
          "type": "folder",
          "children": [
            {
              "name": "SMT-Area-1",
              "type": "folder",
              "children": [
                {
                  "name": "SMT-Line-1",
                  "type": "folder",
                  "children": [
                    {
                      "name": "Printer-Cell",
                      "type": "folder",
                      "children": [
                        {
                          "name": "Printer01",
                          "type": "folder",
                          "children": [
                            {
                              "name": "State",
                              "type": "folder",
                              "topicType": "STATE",
                              "children": [
                                {
                                  "name": "current_job",
                                  "type": "file",
                                  "topicType": "STATE",
                                  "dataType": "RELATION_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "job_id",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "product_id",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "planned_quantity",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "completed_quantity",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "status",
                                      "type": "LONG"
                                    }
                                  ]
                                },
                                {
                                  "name": "alarm_status",
                                  "type": "file",
                                  "topicType": "STATE",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE"
                                }
                              ]
                            },
                            {
                              "name": "Action",
                              "type": "folder",
                              "topicType": "ACTION",
                              "children": [
                                {
                                  "name": "start_job",
                                  "type": "file",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                },
                                {
                                  "name": "stop_job",
                                  "type": "file",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                }
                              ]
                            },
                            {
                              "name": "Metric",
                              "type": "folder",
                              "topicType": "METRIC",
                              "children": [
                                {
                                  "name": "board_cycle_time",
                                  "type": "file",
                                  "topicType": "METRIC",
                                  "dataType": "TIME_SEQUENCE_TYPE",
                                  "generateDashboard": "TRUE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "cycle_time_ms",
                                      "type": "LONG"
                                    }
                                  ]
                                },
                                {
                                  "name": "boards_count",
                                  "type": "file",
                                  "topicType": "METRIC",
                                  "dataType": "TIME_SEQUENCE_TYPE",
                                  "generateDashboard": "TRUE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "good_count",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "ng_count",
                                      "type": "LONG"
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}`;
const Module: FC<ImportModalProps> = (props) => {
  const { importRef, initTreeData } = props;
  const [open, setOpen] = useState(false);
  const formatMessage = useTranslate();
  const { message, modal } = App.useApp();
  const [type, setType] = useState('json');
  const ref = useRef<HTMLDivElement>(null);
  const size = useSize(ref);
  const [jsonValue, setJsonValue] = useState<any>();
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const uploadRef = useRef<any>(null);
  const timer = useRef<NodeJS.Timeout>();
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const [socketUrl, setSocketUrl] = useState('');
  const [loading, setLoading] = useState(false);

  useImperativeHandle(importRef, () => ({
    setOpen: setOpen,
  }));

  // 创建 WebSocket 连接
  const { readyState, disconnect, sendMessage } = useWebSocket(
    socketUrl, // 初始 URL 为 null，表示不立即连接
    {
      reconnectLimit: 0,
      onMessage: (event) => {
        if (event.data === 'pong') return;
        const data = JSON.parse(event.data);
        setSocketData(data);
        if (data.finished) initTreeData({ reset: true });
      },
      onError: (error) => console.error('WebSocket error:', error),
    }
  );

  const save = () => {
    if (type === 'json') {
      if (jsonValue) {
        setLoading(true);
        importExcel({
          value: new Blob([jsonValue], { type: 'application/json' }),
          name: 'file',
          fileName: 'uns.json',
        })
          .then((data) => {
            if (data) {
              const protocol = location.protocol.includes('https') ? 'wss' : 'ws';
              // 创建 WebSocket 连接
              setSocketUrl(
                `${protocol}://${location.host}/inter-api/supos/uns/ws?file=${encodeURIComponent(data)}&token=${getToken()}`
              );
            }
          })
          .catch(() => {
            resetUploadStatus();
          });
      } else {
        message.warning(formatMessage('uns.pleaseJSON'));
      }
    } else {
      if (fileList.length) {
        setLoading(true);
        importExcel({
          value: fileList[0],
          name: 'file',
          fileName: fileList[0].name,
        })
          .then((data) => {
            if (data) {
              const protocol = location.protocol.includes('https') ? 'wss' : 'ws';
              // 创建 WebSocket 连接
              setSocketUrl(
                `${protocol}://${location.host}/inter-api/supos/uns/ws?file=${encodeURIComponent(data)}&token=${getToken()}`
              );
            }
          })
          .catch(() => {
            resetUploadStatus();
          });
      } else {
        message.warning(formatMessage('uns.pleaseUploadTheFile'));
      }
    }
  };

  const close = () => {
    setOpen(false);
    setFileList([]);
    setType('json');
    setJsonValue(undefined);
    if (socketData.finished) {
      resetUploadStatus();
    }
  };
  const { code, finished, msg, task, progress } = socketData;
  const reimport = finished && code !== 200;

  const beforeUpload = (file: any) => {
    const fileType = file.name.split('.').pop();
    if (['json'].includes(fileType.toLowerCase())) {
      setFileList([file]);
    } else {
      message.warning(formatMessage('common.theFileFormatType', { fileType: '.json' }));
    }
    return false;
  };

  const resetUploadStatus = () => {
    setLoading(false);
    setSocketUrl('');
    setSocketData({});
    setFileList([]);
    setJsonValue(undefined);
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
          window.open(`/inter-api/supos/uns/excel/download?path=${socketData.errTipFile}`, '_self');
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
      timer.current = setInterval(() => {
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
  const { copy } = useClipboard();

  if (!open) return null;
  return (
    <ProModal
      className="importModalWrap"
      open={open}
      onCancel={close}
      title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span>{formatMessage('common.import')}</span>
        </div>
      }
      width={460}
      maskClosable={false}
      keyboard={false}
      destroyOnHidden
    >
      {(isFullscreen) => {
        return (
          <Flex vertical style={{ height: isFullscreen ? '100%' : 400 }}>
            <ComRadio
              style={{ margin: '8px 0' }}
              value={type}
              onChange={(e) => {
                setType(e.target.value);
              }}
              options={[
                { label: 'JSON', value: 'json' },
                { label: formatMessage('common.uploadFile'), value: 'document' },
              ]}
            />
            <div style={{ flex: 1, overflow: 'hidden' }}>
              {loading ? (
                <div className="loadingWrap">
                  <InlineLoading
                    status={finished ? (code === 200 ? 'finished' : 'error') : 'active'}
                    description={`${formatMessage('common.importProgress')}：${progress || 0}%${msg || task ? '，' : ''}${finished ? msg : task || ''}`}
                  />
                </div>
              ) : type === 'json' ? (
                <div
                  ref={ref}
                  style={{
                    height: '100%',
                    borderRadius: 4,
                    border: '1px solid rgb(198, 198, 198)',
                    padding: 16,
                    position: 'relative',
                  }}
                  className={styles['custom-theme']}
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
                      copy(jsonValue || placeholder);
                    }}
                  >
                    <Copy />
                  </div>
                  <CodeMirror
                    theme={codemirrorTheme}
                    placeholder={placeholder}
                    onChange={setJsonValue}
                    value={jsonValue}
                    height={(size?.height || 32) - 32 + 'px'}
                    extensions={[json()]}
                  />
                </div>
              ) : (
                <Dragger
                  ref={uploadRef}
                  className={cx('uploadWrap', fileList?.length > 0 && 'uploadWrapFile')}
                  action=""
                  accept=".json"
                  maxCount={1}
                  beforeUpload={beforeUpload}
                  fileList={fileList}
                  onRemove={() => {
                    setFileList([]);
                  }}
                >
                  <FolderAdd size={48} style={{ color: '#E0E0E0' }} />
                  <ComEllipsis style={{ padding: '16px 0' }}>
                    {formatMessage('common.clickOrDragForUpload')}
                  </ComEllipsis>
                  <Button
                    onClick={(e) => {
                      e.stopPropagation();
                      window.open(`/inter-api/supos/uns/excel/template/download?fileType=json`, '_self');
                    }}
                  >
                    <Download />
                    {formatMessage('common.downloadTemplate')}
                  </Button>
                </Dragger>
              )}
            </div>

            <Flex justify="end" gap={8} style={{ marginTop: 16 }}>
              <ComButton onClick={close}>{formatMessage('common.cancel')}</ComButton>
              <ComButton
                loading={reimport ? false : loading}
                disabled={reimport ? false : loading}
                type="primary"
                onClick={reimport ? Reupload : save}
              >
                {formatMessage(reimport ? 'uns.reimport' : 'common.save')}
              </ComButton>
            </Flex>
          </Flex>
        );
      }}
    </ProModal>
  );
};
export default Module;
