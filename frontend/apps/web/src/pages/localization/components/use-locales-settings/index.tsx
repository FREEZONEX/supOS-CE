import { useEffect, useImperativeHandle, useRef, useState } from 'react';
import ProModal from '@/components/pro-modal';
import useTranslate from '@/hooks/useTranslate.ts';
import {
  App,
  Badge,
  Button,
  ConfigProvider,
  Divider,
  Dropdown,
  Flex,
  Form,
  Input,
  Modal,
  Popconfirm,
  Space,
  Switch,
  Tabs,
  type UploadFile,
} from 'antd';
import ProSearch from '@/components/pro-search';
import ComDot from '@/components/com-dot';
import './index.scss';
import { AuthButton } from '@/components';
import { Attachment, Download, TrashCan } from '@carbon/icons-react';
import ComDraggerUpload from '@/components/com-dragger-upload';
import {
  deleteLangApi,
  downloadLanguageFileApi,
  downloadTemplateApi,
  exportLanguageApi,
  getLanguageRecordsApi,
  importLanguageFileApi,
  langEnableApi,
  languageRecordConfirmApi,
} from '@/apis/inter-api/i18n.ts';
import { downloadFn, getToken } from '@/utils';
import { DownOutlined } from '@ant-design/icons';
import { useWebSocket } from 'ahooks';
import InlineLoading from '../../../../components/inline-loading';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { getLangList } from '@/stores/base';

const TabLabel = ({ id, label, color }: { id: string; label?: string; color?: string }) => {
  const formatMessage = useTranslate();
  if (id === 'add')
    return (
      <Flex gap={8} style={{ overflow: 'hidden', opacity: 0.7, width: '100%' }}>
        <span>+</span>
        <div
          style={{
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            width: 'calc(100% - 20px)',
            textAlign: 'left',
          }}
        >
          {formatMessage('Localization.addLanguage')}
        </div>
      </Flex>
    );
  return <ComDot color={color}>{label}</ComDot>;
};

const TabContent = ({ info }: { info: any }) => {
  const [modal, contextHolder] = Modal.useModal();
  const [loading, setLoading] = useState(false);

  const [form] = Form.useForm();
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const configs = [
    {
      name: 'languageCode',
      label: formatMessage('Localization.languageCode'),
      type: 'readInput',
    },
    {
      name: 'languageName',
      label: formatMessage('Localization.languageName'),
      type: 'readInput',
    },
    {
      name: 'hasUsed',
      label: formatMessage('Localization.enabled'),
    },
  ];

  useEffect(() => {
    form.setFieldsValue(info);
  }, [info]);

  return (
    <Flex style={{ height: 370 }} vertical justify="space-between">
      <div>
        <ConfigProvider
          theme={{
            components: {
              Form: {
                itemMarginBottom: 8,
              },
            },
          }}
        >
          <Form
            labelAlign="left"
            colon={false}
            form={form}
            labelCol={{ span: 10 }}
            wrapperCol={{ span: 14 }}
            initialValues={{
              language: 'en',
            }}
            labelWrap
          >
            {configs.map((item) => {
              const { type, ...itemConfig } = item;
              return (
                <Form.Item key={item.name} {...itemConfig}>
                  {type === 'readInput' ? (
                    <Input readOnly variant="borderless" />
                  ) : (
                    <Switch
                      disabled={info?.languageType === 1}
                      style={{ marginLeft: 11 }}
                      onChange={(e) => {
                        langEnableApi({
                          languageCode: info.languageCode,
                          enable: e,
                        }).then(() => {
                          // 获取语言包
                          getLangList();
                        });
                      }}
                    />
                  )}
                </Form.Item>
              );
            })}
          </Form>
        </ConfigProvider>
        <Divider style={{ backgroundColor: '#BBB' }} />
        <div style={{ fontWeight: 500, fontSize: 14, marginBottom: 8 }}>
          {formatMessage('Localization.languagePack')}
        </div>
        <AuthButton
          loading={loading}
          size="small"
          type="text"
          block
          style={{ marginBottom: 60 }}
          onClick={() => {
            return exportLanguageApi({
              languageCode: info.languageCode,
            })
              .then(() => {
                let secondsToGo = 5;
                const instance = modal.success({
                  title: formatMessage('home.exportSuccess'),
                  okText: `${formatMessage('common.ok')}(${secondsToGo})`,
                });
                const timer = setInterval(() => {
                  secondsToGo -= 1;
                  instance.update({ okText: `${formatMessage('common.ok')}(${secondsToGo})` });
                }, 1000);
                setTimeout(() => {
                  clearInterval(timer);
                  instance.destroy();
                }, 5 * 1000);
              })
              .finally(() => {
                setLoading(false);
              });
          }}
        >
          <Flex align="center" justify="space-between" gap={8} style={{ width: '100%' }}>
            <Flex align="center" style={{ flex: 1, opacity: 0.8 }} gap={8}>
              <Attachment />
              <span>{info.languageName}</span>
            </Flex>
            <Download style={{ cursor: 'pointer' }} />
          </Flex>
        </AuthButton>
      </div>
      <Flex justify="flex-end">
        <Popconfirm
          title={formatMessage('common.deleteConfirm')}
          onConfirm={() => {
            deleteLangApi(info.languageCode).then(() => {
              message.success(formatMessage('common.deleteSuccessfully'));
              // 更新语言包
              getLangList();
            });
          }}
        >
          <AuthButton
            disabled={info?.languageType === 1}
            icon={
              <Flex align="center">
                <TrashCan />
              </Flex>
            }
          />
        </Popconfirm>
      </Flex>
      {contextHolder}
    </Flex>
  );
};

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

const AddContent = ({ addRef }: any) => {
  const formatMessage = useTranslate();
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [socketUrl, setSocketUrl] = useState('');
  const timer = useRef<number>(undefined);
  const [socketData, setSocketData] = useState<SocketDataType>({});
  const { message, modal } = App.useApp();
  const uploadRef = useRef<any>(null);

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

  useImperativeHandle(addRef, () => ({
    resetUploadStatus,
  }));

  useEffect(() => {
    if (socketData.finished && disconnect) {
      disconnect();
      clearInterval(timer.current);
      if (socketData.code === 200) {
        message.success(formatMessage('uns.importFinished'));
        setFileList([]);
        // 获取语言包
        getLangList();
        if (socketData.finished) {
          resetUploadStatus();
        }
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
  const save = () => {
    if (fileList.length) {
      setLoading(true);
      importLanguageFileApi({
        value: fileList[0],
        name: 'file',
        fileName: fileList[0].name,
      })
        .then((data) => {
          if (data) {
            const protocol = location.protocol.includes('https') ? 'wss' : 'ws';
            // 创建 WebSocket 连接
            setSocketUrl(
              `${protocol}://${location.host}/inter-api/supos/uns/ws?file=${encodeURIComponent(data)}&i18n=${String(Date.now())}&token=${getToken()}`
            );
          }
        })
        .catch(() => {
          resetUploadStatus();
        });
    }
  };
  const { code, finished, msg, task, runningStatusList } = socketData;

  const reimport = finished && code !== 200;
  return (
    <Flex style={{ height: 370 }} vertical justify="space-between">
      <div>
        <Flex style={{ fontWeight: 500, fontSize: 14, marginBottom: 16, overflow: 'hidden' }} justify="space-between">
          <span
            style={{
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
            title={formatMessage('Localization.languagePack')}
          >
            {formatMessage('Localization.languagePack')}
          </span>
          <AuthButton
            size="small"
            type="text"
            onClick={() => {
              downloadTemplateApi({ fileType: 'excel' }).then((data) => {
                downloadFn({ data, name: 'i18n_languageCode.xlsx' });
              });
              // window.open(`/inter-api/supos/i18n/excel/template/download?fileType=excel`, '_self');
            }}
          >
            <Flex align="center" style={{ cursor: 'pointer' }} gap={8}>
              <Download />
              <span>{formatMessage('common.downloadTemplate')}</span>
            </Flex>
          </AuthButton>
        </Flex>
        {loading ? (
          <Flex vertical className="useLocalesSettingsWrap" justify="center" align="center">
            <div style={{ maxWidth: '90%' }}>
              <InlineLoading
                status={finished ? (code === 200 ? 'finished' : 'error') : 'active'}
                description={`${finished ? msg : task || ''}`}
              />
            </div>
            {runningStatusList?.length && (
              <div className="useLocalesSettingsContent">
                {runningStatusList?.map((m) => {
                  const { code, finished, msg, task, module, totalCount, successCount } = m;
                  const title = `${formatMessage('home.' + module)}：${finished ? msg : task || ''}`;
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
          <ComDraggerUpload
            acceptList={['xlsx']}
            ref={uploadRef}
            onChange={(v) => {
              setFileList(v);
            }}
          />
        )}
      </div>
      <AuthButton
        color="primary"
        variant="solid"
        block
        onClick={reimport ? Reupload : save}
        loading={reimport ? false : loading}
        disabled={reimport ? false : loading}
      >
        {formatMessage(reimport ? 'uns.reimport' : 'common.save')}
      </AuthButton>
    </Flex>
  );
};

const useLocalesSettings = ({ setButtonExportRecords }: any) => {
  const formatMessage = useTranslate();
  const [open, setOpen] = useState(false);
  const addRef = useRef<any>(null);
  const onLocalesModalOpen = () => {
    setOpen(true);
  };
  const [exportRecords, setExportRecords] = useState([]);
  const langData = useI18nStore((state) => state.langList);
  const [searchValue, setSearchValue] = useState('');
  const [filteredLangData, setFilteredLangData] = useState(langData);

  const onClose = () => {
    setOpen(false);
    addRef?.current?.resetUploadStatus();
    setSearchValue('');
    setFilteredLangData(langData);
  };

  useEffect(() => {
    if (searchValue) {
      const filtered = langData.filter(
        (lang: any) =>
          lang.languageName?.toLowerCase().includes(searchValue.toLowerCase()) ||
          lang.languageCode?.toLowerCase().includes(searchValue.toLowerCase())
      );
      setFilteredLangData(filtered);
    } else {
      setFilteredLangData(langData);
    }
  }, [searchValue, langData]);

  useEffect(() => {
    if (open) {
      getLanguageRecordsApi().then((data) => {
        setExportRecords(data);
      });
    }
  }, [open]);

  const LocalesModal = (
    <ProModal
      size="xs"
      className="use-locales-settings"
      open={open}
      title={
        <Flex justify="space-between" align="center" style={{ overflow: 'hidden' }}>
          <div
            title={formatMessage('Localization.localesSetting')}
            style={{ flex: 1, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', minWidth: 0 }}
          >
            {formatMessage('Localization.localesSetting')}
          </div>
          <Dropdown
            onOpenChange={(open) => {
              if (open) {
                getLanguageRecordsApi().then((data) => {
                  setExportRecords(data);
                  const ids = data?.filter((f: any) => !f.confirm)?.map((d: any) => d.id);
                  if (ids?.length > 0) {
                    languageRecordConfirmApi({ ids }).then(() => {
                      setButtonExportRecords((pre: any) => {
                        return pre.map((i: any) => ({
                          ...i,
                          confirm: true,
                        }));
                      });
                      setExportRecords((pre: any) => {
                        return pre.map((i: any) => ({
                          ...i,
                          confirm: true,
                        }));
                      });
                    });
                  }
                });
              }
            }}
            menu={{
              items:
                exportRecords?.length > 0
                  ? [
                      ...(exportRecords?.map((m: any) => {
                        return {
                          label: m?.fileName,
                          key: m.id,
                          extra: <Download style={{ verticalAlign: 'middle' }} />,
                          onClick: () => {
                            downloadLanguageFileApi({ path: m.filePath }).then((data) => {
                              downloadFn({ data, name: m.fileName });
                            });
                          },
                        };
                      }) || []),
                      {
                        type: 'divider',
                      },
                      {
                        key: '-2',
                        label: formatMessage('home.fiveRecord'),
                        disabled: true,
                      },
                    ]
                  : [
                      {
                        disabled: true,
                        label: formatMessage('home.noExport'),
                        key: '-1',
                      },
                    ],
            }}
          >
            <Badge dot={exportRecords?.some((s: any) => !s.confirm)}>
              <Button color="default" variant="filled" iconPosition="end" style={{ padding: '4px 10px' }}>
                <Space>
                  <div
                    title={formatMessage('common.exported')}
                    style={{ maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
                  >
                    {formatMessage('common.exported')}
                  </div>
                  <DownOutlined />
                </Space>
              </Button>
            </Badge>
          </Dropdown>
        </Flex>
      }
      onCancel={onClose}
    >
      {/*<Flex gap={8} align="center" style={{ marginBottom: 16 }}>*/}
      {/*  <span>{formatMessage('Localization.defaultLanguage')}</span>*/}
      {/*  <Select*/}
      {/*    placeholder={formatMessage('common.select')}*/}
      {/*    variant="borderless"*/}
      {/*    style={{ width: 100 }}*/}
      {/*    options={[*/}
      {/*      { value: 'jack', label: 'Jack' },*/}
      {/*      { value: 'lucy', label: 'Lucy' },*/}
      {/*      { value: 'Yiminghe', label: 'yiminghe' },*/}
      {/*    ]}*/}
      {/*  />*/}
      {/*</Flex>*/}
      <Tabs
        tabBarGutter={0}
        tabBarExtraContent={{
          left: (
            <div style={{ marginRight: 16, marginBottom: 16 }}>
              <ProSearch
                value={searchValue}
                onChange={(e) => setSearchValue(e.target.value)}
                style={{ width: 200 }}
                placeholder={formatMessage('Localization.searchLanguage')}
                onClear={() => setSearchValue('')}
                size="sm"
              />
            </div>
          ),
        }}
        className="custom-tab"
        tabPosition="left"
        items={[
          ...filteredLangData.map((lang: any) => {
            return {
              label: <TabLabel id={lang.id} label={lang.languageName} color={lang.hasUsed ? undefined : '#8D8D8D'} />,
              key: lang.id,
              children: <TabContent info={lang} />,
            };
          }),
          {
            label: <TabLabel id="add" />,
            key: 'add',
            children: <AddContent addRef={addRef} />,
          },
        ]}
      />
    </ProModal>
  );

  return {
    LocalesModal,
    onLocalesModalOpen,
  };
};

export default useLocalesSettings;
