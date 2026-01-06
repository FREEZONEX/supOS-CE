import { App, Button, Divider, Empty, Flex, Form, Popover, Radio, Space, Tag } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import ProCard from '@/components/pro-card/ProCard.tsx';
import { useEffect, useState } from 'react';
import { addMcp, deleteMcp, getMcpList } from '@/apis/copilotkit/mcp.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { TrashCan } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import ComButton from '@/components/com-button';
import { formatTimestamp } from '@/utils';
import ProModal from '@/components/pro-modal';
import McpTypeForm from './McpTypeForm.tsx';

const McpSetting = () => {
  const [mcpList, setMcpList] = useState<any[]>([]);
  const [form] = Form.useForm();
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const getMcpListFn = () => {
    getMcpList().then((data) => {
      setMcpList(data);
    });
  };

  useEffect(() => {
    getMcpListFn();
  }, []);

  const handleModalCancel = () => {
    setIsModalOpen(false);
    form.resetFields();
  };

  const handleAddServer = () => {
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    const { json, type, ...restValues } = values;
    let requestData: any = [];
    if (type === 'json') {
      Object.entries(JSON.parse(json)?.mcpServers).forEach(([key, value]: any) => {
        if (value.transportType === 'stdio') {
          requestData.push({
            transportType: value.transportType,
            name: key,
            config: {
              args: value.args || [],
              env: Object.keys(value.env || {}).map((key) => {
                return {
                  key: key,
                  value: value?.env?.[key],
                };
              }),
              command: value.command || 'npx',
            },
          });
        } else {
          requestData.push({
            transportType: value.transportType,
            name: key,
            config: {
              url: value.url,
            },
          });
        }
      });
    } else {
      // Form模式：单个配置
      requestData = restValues;
    }
    return addMcp(requestData).then(() => {
      handleModalCancel();
      message.success(formatMessage('uns.newSuccessfullyAdded'));
      getMcpListFn();
    });
  };

  return (
    <div>
      <ProModal
        title={formatMessage('copilotkit.add')}
        destroyOnHidden
        width={500}
        open={isModalOpen}
        onCancel={handleModalCancel}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="type" initialValue="form">
            <Radio.Group
              options={[
                { value: 'form', label: 'Form' },
                { value: 'json', label: 'Json' },
              ]}
            />
          </Form.Item>
          <McpTypeForm />
          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={handleModalCancel}>{formatMessage('common.cancel')}</Button>
              <ComButton onClick={handleSubmit} type="primary" htmlType="submit">
                {formatMessage('common.save')}
              </ComButton>
            </Space>
          </Form.Item>
        </Form>
      </ProModal>
      <Flex justify="flex-end" align="center" gap={8}>
        <ComButton
          type="primary"
          size="small"
          icon={<PlusOutlined />}
          onClick={() => {
            return handleAddServer();
          }}
        >
          {formatMessage('copilotkit.add')}
        </ComButton>
      </Flex>
      <Divider
        style={{
          background: '#c6c6c6',
          margin: '15px auto',
        }}
      />
      <ProCardContainer minWidth={200}>
        {mcpList?.length > 0 ? (
          mcpList?.map((d) => {
            return (
              <ProCard
                header={{
                  title: d?.clientName,
                  titleDescription: formatTimestamp(d?.lastUsed),
                }}
                key={d.clientName}
                item={d}
                description={false}
                statusHeader={{
                  statusTag: (
                    <Tag style={{ borderRadius: 15, lineHeight: '16px', margin: 0 }} bordered={false} color="success">
                      {d?.transportType || '未知'}
                    </Tag>
                  ),
                  statusInfo: {
                    label: d?.isConnected ? 'ON' : 'OFF',
                    title: formatMessage(`common.status`),
                    color: d?.isConnected ? '#6FDC8C' : '#A8A8A8',
                  },
                }}
                secondaryDescription={
                  <Popover
                    content={
                      <div style={{ maxHeight: 200, width: 200, overflow: 'auto' }}>
                        {d?.tools?.map((t: any) => {
                          return (
                            <div key={t.name}>
                              {t.name}:{t.description}
                              <Divider style={{ margin: '4px 0' }} />
                            </div>
                          );
                        })}
                      </div>
                    }
                    trigger="hover"
                  >
                    <Button type="link">tools({d?.tools?.length || 0})</Button>
                  </Popover>
                }
                actions={(record) => {
                  return [
                    {
                      key: 'delete',
                      label: formatMessage('common.delete'),
                      auth: ButtonPermission['SourceFlow.delete'],
                      extra: (
                        <Flex justify="center" align="center">
                          <TrashCan />
                        </Flex>
                      ),
                      onClick: () => {
                        deleteMcp({
                          endpoint: record.endpoint,
                        }).then(() => {
                          getMcpListFn();
                          message.success(formatMessage('common.deleteSuccessfully'));
                        });
                      },
                    },
                  ];
                }}
              />
            );
          })
        ) : (
          <Empty />
        )}
      </ProCardContainer>
    </div>
  );
};

export { McpSetting };
