import { App, Button, Divider, Empty, Flex, Form, Input, Select, Space, Tag } from 'antd';
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
import { McpTransportForm } from './McpTransportForm.tsx';

const McpSetting = () => {
  const [mcpList, setMcpList] = useState<any[]>([]);
  const [form] = Form.useForm();
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const [transportType, setTransportType] = useState<'sse' | 'streamable-http' | 'stdio'>('stdio');
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
    setTransportType('stdio');
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    const values = await form.validateFields();
    return addMcp(values).then(() => {
      handleModalCancel();
      message.success(formatMessage('uns.newSuccessfullyAdded'));
      getMcpListFn();
    });
  };

  return (
    <div>
      <ProModal title="新增" destroyOnHidden width={500} open={isModalOpen} onCancel={handleModalCancel}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="服务器名称" rules={[{ required: true, message: '请输入服务器名称' }]}>
            <Input placeholder="请输入服务器名称" />
          </Form.Item>
          <Form.Item name="transportType" label="传输模式" initialValue="stdio">
            <Select
              onChange={(value: any) => setTransportType(value)}
              options={[
                {
                  label: 'sse',
                  value: 'sse',
                },
                {
                  label: 'streamable-http',
                  value: 'streamable-http',
                },
                {
                  label: 'stdio',
                  value: 'stdio',
                },
              ]}
            />
          </Form.Item>
          <McpTransportForm transportType={transportType} form={form} />

          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={handleModalCancel}>取消</Button>
              <ComButton onClick={handleSubmit} type="primary" htmlType="submit">
                添加
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
          添加服务器
        </ComButton>
      </Flex>
      <Divider
        style={{
          background: '#c6c6c6',
          margin: '15px auto',
        }}
      />
      <ProCardContainer minWidth={200} hiddenEmpty={false}>
        {mcpList?.length > 0 ? (
          mcpList?.map((d) => {
            return (
              <ProCard
                header={{
                  title: d?.clientName,
                  titleDescription: formatTimestamp(d?.lastUsed),
                }}
                key={d.enpoint}
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
