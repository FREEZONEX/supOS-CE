import { useEffect, useState } from 'react';
import { App, message, Space, Tag } from 'antd';
import { Add, IbmCloudDirectLink_2Dedicated } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import { createSecretKey, deleteSecretKey, querySecretKeyList, updateSecretKey } from '@/apis/inter-api/open-data';
import { ButtonPermission } from '@/common-types/button-permission';
import { ComLayout, ComContent, ComCopy, ProTable, ComBtnTabs, AuthWrapper, AuthButton } from '@/components';
import ServerTabs from './ServerTabs';
import styles from './index.module.scss';
import { formatTimestamp } from '@/utils';
const I18N_NAME = 'OpenData';

interface dataItemTypes {
  id: number;
  appSecretKey: string;
  createTime: number;
  status: number;
}

const OpenData = () => {
  const formatMessage = useTranslate(I18N_NAME);
  const commonFormatMessage = useTranslate();

  const [dataSource, setDataSource] = useState<dataItemTypes[]>([]);
  const { modal } = App.useApp();

  const columns = [
    {
      title: () => formatMessage('appSceretKey'),
      dataIndex: 'appSecretKey',
      render: (text: string) => {
        return (
          <Space>
            <span>{text}</span>
            <ComCopy textToCopy={text} />
          </Space>
        );
      },
    },
    {
      title: () => formatMessage('creationTime'),
      dataIndex: 'createTime',
      width: '20%',
      render: (txt: number) => {
        return txt ? <span>{formatTimestamp(txt)}</span> : null;
      },
    },
    {
      title: () => formatMessage('status'),
      dataIndex: 'status',
      width: '10%',
      render: (txt: any) => {
        return txt === 1 ? (
          <Tag bordered={false} style={{ borderRadius: 15 }} color={'green'}>
            {formatMessage('start')}
          </Tag>
        ) : (
          <Tag bordered={false} style={{ borderRadius: 15 }} color={'magenta'}>
            {formatMessage('disable')}
          </Tag>
        );
      },
    },
    {
      title: () => formatMessage('operation'),
      dataIndex: 'operation',
      width: '10%',
      render: (_: any, record: any) => {
        const disabled = dataSource?.length === 1;
        return (
          <Space>
            {/* 1 启用 0 禁用 */}
            <AuthWrapper
              auth={record.status === 1 ? ButtonPermission['OpenData.disable'] : ButtonPermission['OpenData.enable']}
            >
              <a
                onClick={() => {
                  if (disabled) return;
                  handleToggleDisabled(record);
                }}
                style={{
                  color: 'var(--supos-theme-color)',
                  opacity: disabled ? 0.5 : 1,
                  cursor: disabled ? 'not-allowed' : undefined,
                }}
              >
                {formatMessage(record.status === 1 ? 'disable' : 'start')}
              </a>
            </AuthWrapper>
            <AuthWrapper auth={ButtonPermission['OpenData.delete']}>
              {record.status === 0 && dataSource.length > 1 && (
                <a
                  style={{
                    color: 'var(--supos-theme-color)',
                  }}
                  onClick={() => handleDelete(record)}
                >
                  {formatMessage('delete')}
                </a>
              )}
            </AuthWrapper>
          </Space>
        );
      },
    },
  ];

  useEffect(() => {
    getDataSource();
  }, []);

  const getDataSource = () => {
    querySecretKeyList().then((data = []) => {
      setDataSource(data);
    });
  };

  // 新增key
  const handleAdd = () => {
    createSecretKey().then(() => {
      message.success(formatMessage('newSuccessfullyAdded'));
      getDataSource();
    });
  };

  // 禁用/启用
  const handleToggleDisabled = (record: dataItemTypes) => {
    updateSecretKey({ ...record, status: record.status === 1 ? 0 : 1 }).then(() => {
      message.success(formatMessage(record.status === 1 ? 'disabledSuccessfully' : 'startSuccessfully'));
      getDataSource();
    });
  };

  // 删除
  const handleDelete = (record: dataItemTypes) => {
    modal.confirm({
      title: formatMessage('deleteConfirm'),
      onOk: () => {
        deleteSecretKey(record.id).then(() => {
          message.success(formatMessage('deleteSuccessfully'));
          getDataSource();
        });
      },
      okButtonProps: {
        title: commonFormatMessage('common.confirm'),
      },
      cancelButtonProps: {
        title: commonFormatMessage('common.cancel'),
      },
    });
  };

  const handleSelect = (item: any) => {
    if (item.isNewWindow) {
      window.open(item.value, '_blank');
    }
  };

  return (
    <ComLayout loading={false}>
      <ComContent
        title={
          <div className={styles.titleBox}>
            <Space>
              <IbmCloudDirectLink_2Dedicated />
              <span>{formatMessage('openData')}</span>
            </Space>
            <AuthButton
              auth={ButtonPermission['OpenData.addKey']}
              type="primary"
              icon={<Add />}
              disabled={dataSource.length >= 3}
              onClick={handleAdd}
            >
              {formatMessage('newKey')}
            </AuthButton>
          </div>
        }
        mustHasBack={false}
        style={{
          display: 'flex',
          flexDirection: 'column',
          padding: '30px 36px',
        }}
      >
        <ProTable columns={columns} dataSource={dataSource} pagination={false} />
        <ComBtnTabs
          className={styles.buttonTabs}
          activeKey="builtIn"
          options={[{ label: formatMessage('builtIn'), value: 'builtIn' }]}
          onSelect={handleSelect}
        />
        <ServerTabs appSceretKey={dataSource[0]?.appSecretKey} />
      </ComContent>
    </ComLayout>
  );
};

export default OpenData;
