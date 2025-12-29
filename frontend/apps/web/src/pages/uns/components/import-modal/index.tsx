import { type FC, useEffect, useImperativeHandle, useRef, useState } from 'react';
import { Upload, App, Flex, Button, type UploadFile } from 'antd';
import { useClipboard, useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import ComRadio from '@/components/com-radio';
import ComEllipsis from '@/components/com-ellipsis';
import ComButton from '@/components/com-button';
import CodeMirror from '@uiw/react-codemirror';
import { json, jsonParseLinter } from '@codemirror/lang-json';
import { linter, lintGutter } from '@codemirror/lint';
import { useSize } from 'ahooks';
import { Copy, Download, FolderAdd } from '@carbon/icons-react';
import cx from 'classnames';
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
  "notes": "type:folder|file,topicType:STATE|ACTION|METRIC,dataType:TEMPLATE_TYPE|TIME_SEQUENCE_TYPE|RELATION_TYPE|CALCULATION_REAL_TYPE|CALCULATION_HIST_TYPE|MERGE_TYPE|CITING_TYPE|JSONB_TYPE|,fields.type:INTEGER|LONG|FLOAT|DOUBLE|BOOLEAN|DATETIME|STRING",
  "Template": [],
  "Label": [],
  "UNS": [
    {
      "name": "v1",
      "type": "path",
      "children": [
        {
          "name": "Plant_Name",
          "type": "path",
          "children": [
            {
              "name": "SMT-Area-1",
              "type": "path",
              "children": [
                {
                  "name": "SMT-Line-1",
                  "type": "path",
                  "children": [
                    {
                      "name": "Printer-Cell",
                      "type": "path",
                      "children": [
                        {
                          "name": "Printer01",
                          "type": "path",
                          "children": [
                            {
                              "name": "State",
                              "type": "path",
                              "topicType": "STATE",
                              "children": [
                                {
                                  "name": "current_job",
                                  "type": "topic",
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
                                  "type": "topic",
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
                              "type": "path",
                              "topicType": "ACTION",
                              "children": [
                                {
                                  "name": "start_job",
                                  "type": "topic",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                },
                                {
                                  "name": "stop_job",
                                  "type": "topic",
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
                              "type": "path",
                              "topicType": "METRIC",
                              "children": [
                                {
                                  "name": "board_cycle_time",
                                  "type": "topic",
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
                                  "type": "topic",
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
  const [loading, setLoading] = useState(false);

  useImperativeHandle(importRef, () => ({
    setOpen: setOpen,
  }));

  function processChunk(rawChunk: any) {
    // 分割多个事件
    const events = rawChunk.split('\n\n');
    events.forEach((event: any) => {
      const lines = event.split('\n');
      lines.forEach((line: string) => {
        if (line.includes('code')) {
          try {
            const data = JSON.parse(line);
            setSocketData(data);
            if (data.finished) initTreeData({ reset: true });
          } catch (e) {
            console.error('Error parsing JSON:', e);
          }
        }
      });
    });
  }

  const save = async () => {
    try {
      const fd = new FormData();

      if (type == 'json') {
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
      } else {
        if (fileList.length) {
          fd.append('file', fileList[0] as any, fileList[0].name);
        } else {
          message.warning(formatMessage('uns.pleaseUploadTheFile'));
          return;
        }
      }
      setLoading(true);
      const response = await fetch('/inter-api/supos/uns/importExport/import', {
        method: 'POST',
        body: fd,
      });

      if (!response.ok) {
        setLoading(false);
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      if (reader) {
        const decoder = new TextDecoder('utf-8');
        // 递归读取流数据
        function readStream() {
          reader!
            .read()
            .then(({ done, value }) => {
              if (done) {
                return;
              }
              // 处理流数据块
              const chunk = decoder.decode(value, { stream: true });
              processChunk(chunk);
              // 继续读取下一个数据块
              readStream();
            })
            .catch((error) => {
              console.error(error);
              setLoading(false);
            });
        }

        readStream();
      }
    } catch (error) {
      console.error(error);
      setLoading(false);
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
    if (socketData.finished) {
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
                    style={{ width: '100%', display: 'flex', padding: '0 8px', justifyContent: 'center' }}
                    status={finished ? (code === 200 ? 'finished' : 'error') : 'active'}
                    textMode="custom-lines"
                    lineClamp={3}
                    title={`${formatMessage('common.importProgress')}：${progress || 0}%${msg || task ? '，' : ''}${finished ? msg : task || ''}`}
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
                      window.open(`/inter-api/supos/uns/importExport/template/download?fileType=json`, '_self');
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
