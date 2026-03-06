import { type FC, useState } from 'react';
import { Apps as _App, InformationFilled, FolderAdd, Document, Grid, List } from '@carbon/icons-react';
import { usePagination, useTranslate, useMediaSize } from '@/hooks';
import { formatTimestamp } from '@/utils';
import {
  installAppApi,
  pageListApi,
  startAppApi,
  stopAppApi,
  uninstallAppApi,
  refreshAppStatus,
  uploadApp,
  deleteApp,
  batchStopApp,
  batchStartApp,
  existApp,
} from '@/apis/inter-api/third-apps';
import {
  Flex,
  Pagination,
  Space,
  Tag,
  Button,
  App,
  Popover,
  Empty,
  Upload,
  Checkbox,
  Segmented,
  type PaginationProps,
  Image,
} from 'antd';
import usePropertiesModal from './components/usePropertiesModal.tsx';
import { ButtonPermission } from '@/common-types/button-permission';
import {
  AuthButton,
  ComLayout,
  ComContent,
  ProModal,
  ProTable,
  ProCard,
  SecondaryList,
  ProCardContainer,
} from '@/components';
import styles from './index.module.scss';
const preUrl = '/files/system/resource/supos';
const { Dragger } = Upload;
import defaultIconUrl from '@/assets/home-icons/default.svg';
import { useLocalStorageState } from 'ahooks';

const I18N_NAME = 'AppManagement';

const StatusOptions = [
  {
    value: 'NOT_INSTALL',
    label: 'NOT_INSTALL',
    color: '#E0E0E0',
  },
  {
    value: 'START_WAITING',
    label: 'START_WAITING',
    color: 'orange',
  },
  {
    value: 'STOPPED',
    label: 'STOPPED',
    color: 'orange',
  },
  {
    value: 'RUNNING',
    label: 'RUNNING',
    color: 'green',
  },
  {
    value: 'INSTALL_FAIL',
    label: 'INSTALL_FAIL',
    color: 'red',
  },
  {
    value: 'START_FAIL',
    label: 'START_FAIL',
    color: 'red',
  },
];

export interface PageProps {
  location?: Partial<Location>;
  // 路由title
  title?: string;
}

const CardTag = ({ status, doc, latestFailMsg }: any) => {
  const formatMessage = useTranslate(I18N_NAME);
  const commonFormatMessage = useTranslate();
  const info = StatusOptions?.find((f) => f.value === status) ?? {
    value: status,
    label: status,
    color: 'blue',
  };
  return (
    <Flex align="center" justify="flex-end" gap={4}>
      {['INSTALL_FAIL', 'START_FAIL'].includes(status) && latestFailMsg && (
        <Popover
          content={<div style={{ maxWidth: 400, maxHeight: 400, overflow: 'auto' }}>{latestFailMsg}</div>}
          title={commonFormatMessage('common.errorInfo')}
        >
          <InformationFilled color="red" />
        </Popover>
      )}
      <Tag
        bordered={false}
        color={info?.color}
        title={formatMessage(info?.label)}
        style={{
          borderRadius: 9,
          height: 16,
          lineHeight: '16px',
          maxWidth: 120,
          overflow: 'hidden',
          whiteSpace: 'nowrap',
          textOverflow: 'ellipsis',
        }}
      >
        {formatMessage(info?.label)}
      </Tag>
      {doc && (
        <Button
          size="small"
          type="text"
          onClick={() => {
            if (doc) {
              window.open(`${preUrl}${doc}`);
            }
          }}
          style={{ color: 'var(--supos-text-color)' }}
          icon={<Document size={16} />}
        >
          {formatMessage('document')}
        </Button>
      )}
    </Flex>
  );
};

const appHealthStatusObj: any = {
  gray: {
    color: '#A8A8A8',
    key: 'unRun',
  },
  yellow: {
    color: '#F1C21B',
    key: 'partRun',
  },
  green: {
    color: '#6FDC8C',
    key: 'run',
  },
};

const Index: FC<PageProps> = ({ title }) => {
  const commonFormatMessage = useTranslate();
  const formatMessage = useTranslate(I18N_NAME);
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useLocalStorageState<string>('SUPOS_APP_MODE', {
    defaultValue: 'card',
  });
  const { message, modal } = App.useApp();
  const [buttonLoading, setButtonLoading] = useState(false);
  const [fileList, setFileList] = useState<any[]>([]);
  const { isH5 } = useMediaSize();
  const [uploadTitle, setUploadTitle] = useState('save');
  const showTotal: PaginationProps['showTotal'] = (total) =>
    isH5 ? null : `${commonFormatMessage('common.total')}  ${total}  ${commonFormatMessage('common.items')}`;

  const { loading, pagination, data, refreshRequest, rowSelection, selectedRows, setLoading } = usePagination({
    fetchApi: pageListApi,
    autoRefresh: true,
    initPageSize: 18,
  });
  const onClose = () => {
    setFileList([]);
    setOpen(false);
    setUploadTitle('save');
    setButtonLoading(false);
  };
  const onSave = () => {
    if (fileList.length) {
      setButtonLoading(true);
      setUploadTitle('uploading');
      const item = fileList[0];
      uploadApp([{ value: item, name: 'file', fileName: item.name }])
        .then((res: any) => {
          setUploadTitle('unZiping');
          const fn = () => {
            existApp(res?.data)
              ?.then((data) => {
                if (data) {
                  refreshRequest?.();
                  onClose();
                  message.success(commonFormatMessage('common.optsuccess'));
                  setButtonLoading(false);
                  setUploadTitle('save');
                } else {
                  setTimeout(() => {
                    fn();
                  }, 2000);
                }
              })
              .catch(() => {
                setButtonLoading(false);
                setUploadTitle('save');
              });
          };
          setTimeout(() => {
            fn();
          }, 2000);
        })
        .catch(() => {
          setButtonLoading(false);
          setUploadTitle('save');
        });
    } else {
      message.warning(commonFormatMessage('uns.pleaseUploadTheFile'));
    }
  };
  const beforeUpload = (file: any) => {
    const fileType = file.name.split('.').pop();
    if (['zip'].includes(fileType.toLowerCase())) {
      setFileList([file]);
    } else {
      message.warning(commonFormatMessage('common.theFileFormatType', { fileType: '.zip' }));
    }
    return false;
  };
  const columns: any = [
    {
      dataIndex: 'appShowName',
      ellipsis: true,
      fixed: 'left',
      title: () => formatMessage('appName'),
      width: '10%',
      render: (text: string, record: any) => {
        return (
          <>
            <Image
              wrapperStyle={{ marginRight: 8 }}
              preview={false}
              src={`${preUrl}${record.icon}`}
              width={20}
              height={20}
              fallback={defaultIconUrl}
            />
            {text}
          </>
        );
      },
    },
    {
      dataIndex: 'vendorName',
      ellipsis: true,
      title: () => commonFormatMessage('common.dev'),
      width: '10%',
    },
    {
      dataIndex: 'version',
      ellipsis: true,
      title: () => commonFormatMessage('common.version'),
      width: '10%',
    },
    {
      dataIndex: 'dependencies',
      ellipsis: true,
      title: () => commonFormatMessage('common.dependencies'),
      width: '10%',
    },
    {
      dataIndex: 'lastInstallTime',
      ellipsis: true,
      title: () => formatMessage('installationTime'),
      width: '10%',
      render: (t: number) => {
        return formatTimestamp(t);
      },
    },
    {
      dataIndex: 'status',
      ellipsis: true,
      title: () => commonFormatMessage('common.states'),
      width: '10%',
      render: (status: string, record: any) => {
        const info = StatusOptions?.find((f) => f.value === status) ?? {
          value: status,
          label: status,
          color: 'blue',
        };
        return (
          <Flex align="center" gap={4}>
            {['INSTALL_FAIL', 'START_FAIL'].includes(status) && record?.latestFailMsg && (
              <Popover
                content={<div style={{ maxWidth: 400, maxHeight: 400, overflow: 'auto' }}>{record?.latestFailMsg}</div>}
                title={commonFormatMessage('common.errorInfo')}
              >
                <InformationFilled color="red" />
              </Popover>
            )}
            <Tag bordered={false} color={info?.color} style={{ borderRadius: 9, height: 16, lineHeight: '16px' }}>
              {formatMessage(info?.label)}
            </Tag>
            {record?.doc && (
              <Button
                size="small"
                type="text"
                onClick={() => {
                  if (record?.doc) {
                    window.open(`${preUrl}${record?.doc}`);
                  }
                }}
                style={{ color: 'var(--supos-text-color)' }}
                icon={<Document size={16} />}
              >
                {formatMessage('document')}
              </Button>
            )}
          </Flex>
        );
      },
    },
    {
      dataIndex: 'appHealthStatus',
      title: () => formatMessage('healthStates'),
      width: '5%',
      render: (text: string) => {
        const info = appHealthStatusObj[text];
        return (
          <Flex justify="flex-start" align="center" gap={4}>
            <div style={{ width: 5, height: 5, borderRadius: '50%', background: info?.color }} />
            {commonFormatMessage(`common.${info?.key}`)}
          </Flex>
        );
      },
    },
  ];
  const onBatchHandle = (str: string) => {
    setLoading(true);
    const api: any = {
      batchStopApp,
      batchStartApp,
    };
    api?.[str]?.({
      appIdList: selectedRows?.map((m: any) => m.appId),
    })
      .then((data: any) => {
        console.log(data);
        message.success(data);
        refreshRequest?.(undefined, { clearSelect: true });
      })
      .finally(() => {
        setLoading(false);
      });
  };
  const onOptHandle = (apiStr: string, d: any) => {
    const api: any = {
      uninstallAppApi,
      installAppApi,
      stopAppApi,
      startAppApi,
      deleteApp,
    };
    return api?.[apiStr]?.({ appId: d.appId }).then(() => {
      refreshRequest?.();
      message.success(commonFormatMessage('common.optsuccess'));
    });
  };
  const { PropertiesModal, setPropertiesOpen } = usePropertiesModal({
    successBackFn: refreshRequest,
  });
  const actions = (record: any) => {
    const btns = [
      {
        type: 'Loading',
        key: 'loading',
        label: formatMessage(record.status),
      },
      {
        key: 'start',
        label: formatMessage('start'),
        auth: ButtonPermission['AppManagement.start'],
        button: {
          color: 'primary',
          variant: 'outlined',
        },
        onClick: () => onOptHandle('startAppApi', record),
      },
      {
        key: 'pause',
        label: formatMessage('pause'),
        auth: ButtonPermission['AppManagement.pause'],
        button: {
          style: {
            backgroundColor: '#6F6F6F',
            color: 'white',
          },
        },
        onClick: () => onOptHandle('stopAppApi', record),
      },
      {
        key: 'install',
        label: formatMessage('install'),
        auth: ButtonPermission['AppManagement.install'],
        button: {},
        onClick: () => onOptHandle('installAppApi', record),
      },
      {
        key: 'reInstall',
        label: formatMessage('reInstall'),
        auth: ButtonPermission['AppManagement.install'],
        button: {},
        onClick: () => onOptHandle('installAppApi', record),
      },
      {
        key: 'configUpdate',
        label: formatMessage('updateConfig'),
        auth: ButtonPermission['AppManagement.configUpdate'],
        button: {
          style: { color: 'var(--supos-theme-color)' },
        },
        onClick: () => setPropertiesOpen({ appProperties: record.appProperties, appId: record.appId }),
      },
      {
        key: 'delete',
        label: formatMessage('delete'),
        auth: ButtonPermission['AppManagement.delete'],
        button: {},
        onClick: () => {
          modal.confirm({
            title: commonFormatMessage('common.deleteConfirm'),
            onOk: () => {
              onOptHandle('deleteApp', record);
            },
            okButtonProps: {
              title: commonFormatMessage('common.confirm'),
            },
            cancelButtonProps: {
              title: commonFormatMessage('common.cancel'),
            },
          });
        },
      },
      {
        key: 'unInstall',
        label: formatMessage('unInstall'),
        auth: ButtonPermission['AppManagement.unInstall'],
        button: {
          style: { color: 'var(--supos-text-color)' },
        },
        onClick: () => onOptHandle('uninstallAppApi', record),
      },
    ];
    const obj: any = {
      NOT_INSTALL: ['install', 'configUpdate', 'delete'],
      INSTALL_FAIL: ['reInstall', 'configUpdate', 'delete'],
      STOPPED: ['start', 'configUpdate', 'unInstall'],
      RUNNING: ['pause', 'configUpdate'],
      START_FAIL: ['start', 'configUpdate', 'unInstall'],
    };
    return [
      ...(obj?.[record.status]
        ? obj?.[record.status]?.map((i: string) => btns?.find((f) => f.key === i)) || []
        : [btns[0]]),
      {
        key: 'renew',
        label: commonFormatMessage('common.refresh'),
        // icon: <Renew />,
        onClick: () => {
          return refreshAppStatus(record.appId).then(() => {
            message.success(commonFormatMessage('common.refreshSuccessful'));
          });
        },
      },
    ];
  };
  return (
    <ComLayout loading={loading}>
      {PropertiesModal}
      <ComContent
        title={
          <div>
            <_App size={20} style={{ justifyContent: 'center', verticalAlign: 'middle' }} /> {title}
          </div>
        }
        hasBack={false}
        style={{
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
        }}
        extra={
          <Space>
            <AuthButton
              auth={ButtonPermission['AppManagement.start']}
              style={{ height: 28 }}
              color="primary"
              variant="outlined"
              disabled={!selectedRows?.length}
              onClick={() => {
                onBatchHandle('batchStartApp');
              }}
            >
              <Flex align="center" gap={6}>
                {formatMessage('batchStart')}
              </Flex>
            </AuthButton>
            <AuthButton
              auth={ButtonPermission['AppManagement.pause']}
              style={{ height: 28 }}
              color="primary"
              variant="outlined"
              disabled={!selectedRows?.length}
              onClick={() => {
                onBatchHandle('batchStopApp');
              }}
            >
              <Flex align="center" gap={6}>
                {formatMessage('batchPause')}
              </Flex>
            </AuthButton>
            <AuthButton
              auth={ButtonPermission['AppManagement.upload']}
              style={{ height: 28 }}
              onClick={() => {
                setOpen(true);
              }}
            >
              <Flex align="center" gap={6}>
                {commonFormatMessage('common.uploadApp')}
              </Flex>
            </AuthButton>
          </Space>
        }
        className={styles['app-management']}
      >
        <Flex justify="flex-end" align="center" style={{ marginBottom: 16, marginTop: 16, paddingRight: 16 }}>
          <Segmented
            size="small"
            value={mode}
            onChange={(v) => setMode(v)}
            options={[
              {
                value: 'card',
                icon: (
                  <span className={styles['flex']} title={commonFormatMessage('common.cardMode')}>
                    <Grid />
                  </span>
                ),
              },
              {
                value: 'list',
                icon: (
                  <span className={styles['flex']} title={commonFormatMessage('common.listMode')}>
                    <List />
                  </span>
                ),
              },
            ]}
          />
        </Flex>

        <div style={{ flex: 1, padding: '0 16px 16px', overflow: 'auto' }}>
          {mode === 'card' ? (
            data?.length > 0 ? (
              <Checkbox.Group
                style={{ width: '100%' }}
                value={rowSelection.selectedRowKeys}
                onChange={(checkedValues) => {
                  rowSelection.onChange(checkedValues);
                }}
              >
                <ProCardContainer>
                  {data?.map((d: any) => {
                    const info = appHealthStatusObj?.[d?.appHealthStatus];
                    return (
                      <ProCard
                        key={d.id}
                        loading={loading}
                        statusHeader={{
                          allowCheck: true,
                          statusInfo: {
                            label: commonFormatMessage(`common.${info?.key}`),
                            title: formatMessage(`healthStates`),
                            color: info?.color,
                          },
                          statusTag: <CardTag status={d.status} doc={d.doc} latestFailMsg={d.latestFailMsg} />,
                        }}
                        styles={{
                          // card: { width: 320, height: 285 },
                          card: { height: 280 },
                        }}
                        header={{
                          iconSrc: `${preUrl}${d.icon}`,
                          title: d.appShowName,
                          titleDescription: formatTimestamp(d?.lastInstallTime) || ' ',
                        }}
                        description={d.description}
                        secondaryDescription={
                          <SecondaryList
                            options={[
                              {
                                label: commonFormatMessage('common.dev'),
                                content: d?.vendorName,
                                span: 24,
                                key: 'dev',
                              },
                              {
                                label: commonFormatMessage('common.version'),
                                content: d?.version,
                                span: 24,
                                key: 'version',
                              },
                              {
                                label: commonFormatMessage('common.dependencies'),
                                content: d?.dependencies,
                                span: 24,
                                key: 'dependencies',
                              },
                            ]}
                          />
                        }
                        actions={actions}
                        item={d}
                        value={d.id}
                      />
                    );
                  })}
                </ProCardContainer>
              </Checkbox.Group>
            ) : (
              <Flex style={{ width: '100%' }}>
                <Empty style={{ width: '100%' }} />
              </Flex>
            )
          ) : (
            <ProTable
              resizeable
              style={{ height: '100%' }}
              scroll={{ y: 'calc(100vh  - 285px)', x: 'max-content' }}
              dataSource={data as any}
              columns={columns}
              operationOptions={{
                render: actions,
              }}
              pagination={{
                total: pagination?.total,
                showTotal: showTotal,
                style: { display: 'flex', justifyContent: 'flex-end', padding: '10px 0' },
                pageSize: pagination?.pageSize || 20,
                current: pagination?.page,
                showQuickJumper: true,
                showSizeChanger: true,
                pageSizeOptions: pagination?.pageSizes,
                onChange: pagination.onChange,
                onShowSizeChange: (current: number, size: number) => {
                  pagination.onChange({ page: current, pageSize: size });
                },
              }}
              rowSelection={rowSelection}
            />
          )}
        </div>
        {mode === 'card' && (
          <Pagination
            size="small"
            className="custom-pagination"
            align="center"
            style={{ margin: '20px 0' }}
            total={pagination?.total}
            showSizeChanger={false}
            onChange={pagination.onChange}
            pageSize={pagination?.pageSize || 20}
            current={pagination?.page}
          />
        )}
      </ComContent>
      <ProModal
        aria-label=""
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>{commonFormatMessage('common.uploadApp')}</span>
          </div>
        }
        onCancel={onClose}
        open={open}
        className="importModalWrap"
        size="xxs"
      >
        <Dragger
          className="uploadWrap"
          action=""
          accept=".zip"
          maxCount={1}
          fileList={fileList}
          disabled={buttonLoading}
          beforeUpload={beforeUpload}
          onRemove={() => {
            setFileList([]);
          }}
        >
          <Flex vertical align="center" gap={10}>
            <FolderAdd size={100} style={{ color: '#E0E0E0' }} />
            <span style={{ fontSize: 12 }}>
              {commonFormatMessage('common.theFileFormatType', { fileType: '.zip' })}
            </span>
          </Flex>
        </Dragger>
        <Button
          loading={buttonLoading}
          color="primary"
          variant="solid"
          block
          onClick={onSave}
          style={{ marginTop: 20 }}
        >
          {commonFormatMessage(`common.${uploadTitle}`)}
        </Button>
      </ProModal>
    </ComLayout>
  );
};

export default Index;
