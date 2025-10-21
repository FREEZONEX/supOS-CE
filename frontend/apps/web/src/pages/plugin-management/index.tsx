import { type FC, useEffect, useRef, useState } from 'react';
import type { PageProps } from '@/common-types';
import { SoftwareResourceCluster as _App, InformationFilled, FolderAdd } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import { Flex, Tag, App, Popover, Empty, Upload, Button } from 'antd';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { getPluginListApi, installPluginApi, unInstallPluginApi, upgradePluginApi } from '@/apis/inter-api/plugin.ts';
import useSimpleRequest from '../../hooks/useSimpleRequest.ts';
import IconImage from '../../components/icon-image';
import { useThemeStore } from '@/stores/theme-store.ts';
import { formatTimestamp } from '@/utils/format.ts';
import { preloadPluginLang } from '@/utils/plugin.ts';
import { connectI18nMessage, useI18nStore } from '@/stores/i18n-store.ts';
import { fetchBaseStore, setPluginList } from '@/stores/base';
import { useActivate } from '@/contexts/tabs-lifecycle-context.ts';
import { useTabsContext } from '@/contexts/tabs-context.ts';
import { useLocalStorageState } from 'ahooks';
import styles from './index.module.scss';
import ProTable from '@/components/pro-table';
import ProModal from '@/components/pro-modal/index.tsx';
import ProCard from '@/components/pro-card/ProCard.tsx';
import SecondaryList from '@/components/pro-card/SecondaryList.tsx';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import ComSegmented from '@/components/com-segmented';
const { Dragger } = Upload;

const StatusOptions = [
  {
    value: 'notInstall',
    label: 'plugin.notInstall',
    color: '#E0E0E0',
  },
  {
    value: 'installFail',
    label: 'plugin.installFail',
    color: 'red',
  },
  {
    value: 'installed',
    label: 'plugin.installed',
    color: 'green',
  },
];

const CardTag = ({ status, latestFailMsg }: any) => {
  const formatMessage = useTranslate();
  const info = StatusOptions?.find((f) => f.value === status) ?? {
    value: status,
    label: 'plugin.' + status,
    color: 'blue',
  };
  return (
    <Flex align="center" gap={4}>
      {['installFail'].includes(status) && latestFailMsg && (
        <Popover
          content={<div style={{ maxWidth: 400, maxHeight: 400, overflow: 'auto' }}>{latestFailMsg}</div>}
          title={formatMessage('common.errorInfo')}
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
    </Flex>
  );
};

const IconImageWrapper = ({ record }: any) => {
  const primaryColor = useThemeStore((state) => state.primaryColor);
  const icon =
    record?.plugInfoYml?.resources?.find((f: any) => f.type === 2)?.icon ||
    record?.plugInfoYml?.route?.name ||
    record?.plugInfoYml?.name;

  return (
    <>
      <IconImage theme={primaryColor} width={20} height={20} wrapperStyle={{ marginRight: 8 }} iconName={icon} />
      {record?.plugInfoYml?.showName}
    </>
  );
};
const Index: FC<PageProps> = ({ title }) => {
  const onSuccessCallback = (data: any[]) => {
    setPluginList(data);
  };
  const commonFormatMessage = useTranslate();
  const lang = useI18nStore((state) => state.lang);
  const [open, setOpen] = useState(false);
  const [fileList, setFileList] = useState<any[]>([]);
  const selectName = useRef('');
  const [uploadTitle, setUploadTitle] = useState('save');
  const isFirstRender = useRef(true);
  const primaryColor = useThemeStore((state) => state.primaryColor);

  const [mode, setMode] = useLocalStorageState<string>('SUPOS_PLUGIN_MODE', {
    defaultValue: 'card',
  });
  const { modal } = App.useApp();
  const { loading, data, refreshRequest } = useSimpleRequest({
    fetchApi: getPluginListApi,
    onSuccessCallback,
    // autoRefresh: true,
  });
  useActivate(() => {
    getPluginListApi();
  });

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
    } else {
      refreshRequest();
    }
  }, [lang]);

  const columns: any = [
    {
      dataIndex: 'showName',
      ellipsis: true,
      fixed: 'left',
      title: () => commonFormatMessage('common.pluginName'),
      width: '20%',
      render: (_: string, record: any) => {
        return <IconImageWrapper record={record} />;
      },
    },
    {
      dataIndex: 'vendorName',
      ellipsis: true,
      title: () => commonFormatMessage('common.dev'),
      width: '10%',
      render: (_: string, record: any) => {
        return record?.plugInfoYml?.vendorName;
      },
    },
    {
      dataIndex: 'version',
      ellipsis: true,
      title: () => commonFormatMessage('common.version'),
      width: '10%',
      render: (_: string, record: any) => {
        return record?.plugInfoYml?.version;
      },
    },
    {
      dataIndex: 'name',
      ellipsis: true,
      title: () => commonFormatMessage('common.name'),
      width: '10%',
    },
    {
      dataIndex: 'description',
      ellipsis: true,
      title: () => commonFormatMessage('common.description'),
      width: '23%',
      render: (_: string, record: any) => {
        return record?.plugInfoYml?.description;
      },
    },
    {
      dataIndex: 'installTime',
      ellipsis: true,
      title: () => commonFormatMessage('common.installTime'),
      width: '10%',
      render: (t: number) => {
        return formatTimestamp(t);
      },
    },
    {
      dataIndex: 'installStatus',
      ellipsis: true,
      title: () => commonFormatMessage('common.states'),
      width: '5%',
      render: (status: string) => {
        const info = StatusOptions?.find((f) => f.value === status) ?? {
          value: status,
          label: 'plugin.' + status,
          color: 'blue',
        };
        return (
          <Tag bordered={false} color={info?.color} style={{ borderRadius: 9, height: 16, lineHeight: '16px' }}>
            {commonFormatMessage(info?.label)}
          </Tag>
        );
      },
    },
  ];
  const onClose = () => {
    setFileList([]);
    setUploadTitle('save');
    setOpen(false);
    setButtonLoading(false);
  };
  const openModal = (info: any) => {
    selectName.current = info.name;
    setOpen(true);
  };
  const { message } = App.useApp();

  const [buttonLoading, setButtonLoading] = useState(false);
  const onSave = () => {
    if (fileList.length) {
      setButtonLoading(true);
      setUploadTitle('uploading');
      const item = fileList[0];
      upgradePluginApi([
        { value: item, name: 'file', fileName: item.name },
        { value: String(selectName.current), name: 'name' },
      ])
        .then(() => {
          setUploadTitle('unZiping');
          setTimeout(() => {
            setButtonLoading(false);
            setUploadTitle('save');
            refreshRequest?.();
            onClose();
            message.success(commonFormatMessage('common.optsuccess'));
          }, 1000);
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
    if (file.size <= 1024 * 1024 * 1024 * 2) {
      if (['gz'].includes(fileType.toLowerCase())) {
        setFileList([file]);
      } else {
        message.warning(commonFormatMessage('common.theFileFormatType', { fileType: '.gz' }));
      }
    } else {
      message.warning(commonFormatMessage('common.theFileSizeMax', { size: '2GB' }));
    }
    return false;
  };
  const { TabsContext } = useTabsContext();
  const onOptHandle = (apiStr: string, d: any) => {
    const api: any = {
      installPluginApi,
      unInstallPluginApi,
    };
    return api?.[apiStr]?.({ name: d?.name })
      .then(async () => {
        if (apiStr === 'installPluginApi') {
          // 安装成功，预先加载国际化
          try {
            const langMessages = await preloadPluginLang(
              [{ name: `/${d?.plugInfoYml?.route?.name}`, backendName: d?.name }],
              lang
            );
            connectI18nMessage(langMessages);
          } catch (e) {
            console.error('插件国际化', e);
          }
        } else {
          // 移除多页签
          TabsContext?.current?.onCloseTab?.(`/${d?.plugInfoYml?.route?.name}`);
        }
        refreshRequest?.();
        // 成功后刷新下菜单
        fetchBaseStore?.();
        message.success(commonFormatMessage('common.optsuccess'));
      })
      .finally(() => {});
  };
  const actions = (record: any) => {
    const btns = [
      {
        type: 'Loading',
        key: 'loading',
        label: commonFormatMessage('plugin.' + record.installStatus),
      },
      {
        key: 'install',
        label: commonFormatMessage('plugin.install'),
        auth: ButtonPermission['PluginManagement.install'],
        button: {
          type: 'primary',
        },
        onClick: () => onOptHandle('installPluginApi', record),
      },
      {
        key: 'unInstall',
        label: commonFormatMessage('common.unInstall'),
        auth: ButtonPermission['PluginManagement.unInstall'],
        disabled: record?.plugInfoYml?.removable === false,
        button: {},
        onClick: () => {
          modal.confirm({
            title: commonFormatMessage('common.uninstallConfirm'),
            content: commonFormatMessage('common.clearData'),
            onOk: () => {
              return onOptHandle('unInstallPluginApi', record);
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
        key: 'update',
        label: commonFormatMessage('plugin.update'),
        auth: ButtonPermission['PluginManagement.update'],
        button: {
          type: 'primary',
        },
        onClick: () => openModal(record),
      },
    ];
    const obj: any = {
      notInstall: ['install', 'update'],
      installFail: ['install', 'update'],
      installed: ['unInstall', 'update'],
    };
    return obj?.[record.installStatus]
      ? obj?.[record.installStatus]?.map((i: string) => btns?.find((f) => f.key === i))
      : [btns[0]];
  };
  return (
    <ComLayout loading={loading}>
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
        className={styles['plugin-management']}
      >
        <ComSegmented value={mode} onChange={setMode} defaultValue="card" />
        <div style={{ flex: 1, padding: '0 16px 16px', overflow: 'auto', alignItems: 'center' }}>
          {mode === 'card' ? (
            data.length > 0 ? (
              <ProCardContainer>
                {data?.map((d: any) => {
                  const { plugInfoYml = {}, name } = d || {};
                  const icon =
                    plugInfoYml?.resources?.find((f: any) => f.type === 2)?.icon ||
                    plugInfoYml?.route?.name ||
                    plugInfoYml?.name;
                  return (
                    <ProCard
                      key={plugInfoYml?.name}
                      statusHeader={{
                        statusTag: <CardTag status={d?.installStatus} latestFailMsg={d.latestFailMsg} />,
                      }}
                      header={{
                        customIcon: <IconImage theme={primaryColor} iconName={icon} />,
                        title: plugInfoYml.showName,
                        titleDescription: formatTimestamp(d?.installTime),
                      }}
                      description={plugInfoYml?.description}
                      secondaryDescription={
                        <SecondaryList
                          options={[
                            {
                              label: commonFormatMessage('common.dev'),
                              content: plugInfoYml?.vendorName,
                              span: 24,
                              key: 'dev',
                            },
                            {
                              label: commonFormatMessage('common.version'),
                              content: plugInfoYml?.version,
                              span: 24,
                              key: 'version',
                            },
                            {
                              label: commonFormatMessage('common.name'),
                              content: name,
                              span: 24,
                              key: 'name',
                            },
                          ]}
                        />
                      }
                      actions={actions}
                      item={d}
                    />
                  );
                  // return <LoadingCard key={d.name} openModal={openModal} d={d} refreshRequest={refreshRequest} />;
                })}
              </ProCardContainer>
            ) : (
              <Empty />
            )
          ) : (
            <ProTable
              resizeable
              style={{ height: '100%' }}
              scroll={{ y: 'calc(100vh  - 285px)', x: 'max-content' }}
              dataSource={data as any}
              columns={columns}
              pagination={false}
              operationOptions={{
                render: actions,
              }}
            />
          )}
        </div>
      </ComContent>
      <ProModal
        aria-label=""
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>{commonFormatMessage('plugin.upload')}</span>
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
          accept=".tar.gz"
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
            <span style={{ fontSize: 12 }}>{commonFormatMessage('common.theFileFormatType', { fileType: '.gz' })}</span>
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
