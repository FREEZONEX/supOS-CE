import { type FC, useCallback, useEffect, useState } from 'react';
import { Add, Renew, Connect } from '@/components/lucide-icon/carbon';
import { PageTitleRow } from '@/components/lucide-icon';
import { App, Button, Form, Input, Modal, Space, Table, Tabs, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { PageProps } from '@/common-types';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { useTranslate } from '@/hooks';
import {
  createClusterToken,
  forceSyncCloudSync,
  getCloudSyncConfig,
  getClusterNodes,
  getClusterOutbox,
  getClusterScopes,
  handshakeCluster,
  type ClusterNode,
  type ClusterOutbox,
  type ClusterScope,
} from '@/apis/core-api/cloudsync';
import styles from './index.module.scss';

const formatTime = (value?: number) => (value ? new Date(value).toLocaleString() : '-');

const ClusterPage: FC<PageProps> = ({ title }) => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const [tokenForm] = Form.useForm();
  const [handshakeForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [syncing, setSyncing] = useState(false);
  const [activeTab, setActiveTab] = useState('nodes');
  const [nodes, setNodes] = useState<ClusterNode[]>([]);
  const [scopes, setScopes] = useState<ClusterScope[]>([]);
  const [outbox, setOutbox] = useState<ClusterOutbox[]>([]);
  const [tokenOpen, setTokenOpen] = useState(false);
  const [handshakeOpen, setHandshakeOpen] = useState(false);
  const [createdToken, setCreatedToken] = useState('');

  const load = useCallback(async (options: { silent?: boolean } = {}) => {
    if (!options.silent) {
      setLoading(true);
    }
    try {
      const [nodeResp, scopeResp, outboxResp] = await Promise.all([
        getClusterNodes(),
        getClusterScopes(),
        getClusterOutbox(),
      ]);
      setNodes(nodeResp?.list || []);
      setScopes(scopeResp?.list || []);
      setOutbox(outboxResp?.list || []);
    } finally {
      if (!options.silent) {
        setLoading(false);
      }
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const createToken = async () => {
    const values = await tokenForm.validateFields();
    const resp = await createClusterToken({
      tokenName: values.tokenName,
      description: values.description || '',
    });
    setCreatedToken(resp?.token || '');
  };

  const handshake = async () => {
    const values = await handshakeForm.validateFields();
    const resp = await handshakeCluster({
      token: values.token,
      nodeKey: values.nodeKey,
      nodeName: values.nodeName || values.nodeKey,
      role: values.role || 'edge',
      endpoint: values.endpoint || '',
      capabilitiesJson: values.capabilitiesJson || '{}',
    });
    message.success(`握手成功：${resp?.protocolVersion || 'sync.v1'}`);
    setHandshakeOpen(false);
    await load();
  };

  const forceSync = async () => {
    setSyncing(true);
    try {
      const config = await getCloudSyncConfig();
      if (!config?.configured) {
        message.warning(
          formatMessage('cluster.forceSyncNotConfigured', {}, 'Configure Edge Connection before forcing sync.')
        );
        return;
      }
      await forceSyncCloudSync();
      message.success(formatMessage('cluster.forceSyncSuccess', {}, 'Force sync triggered and list refreshed.'));
      await load({ silent: true });
    } catch {
      message.error(formatMessage('cluster.forceSyncFailed', {}, 'Force sync failed.'));
    } finally {
      setSyncing(false);
    }
  };

  const nodeColumns: ColumnsType<ClusterNode> = [
    { title: '节点 Key', dataIndex: 'nodeKey', width: 220 },
    { title: '名称', dataIndex: 'nodeName', width: 220 },
    { title: '角色', dataIndex: 'role', width: 120, render: (value) => <Tag>{value}</Tag> },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (value) => <Tag color={value === 'online' ? 'green' : 'default'}>{value}</Tag>,
    },
    { title: 'Endpoint', dataIndex: 'endpoint', width: 260, ellipsis: true, render: (value) => value || '-' },
    { title: '最后连接', dataIndex: 'lastSeenTime', width: 180, render: formatTime },
  ];

  const scopeColumns: ColumnsType<ClusterScope> = [
    { title: '节点 ID', dataIndex: 'nodeId', width: 100 },
    { title: '范围类型', dataIndex: 'scopeType', width: 140, render: (value) => <Tag color="blue">{value}</Tag> },
    { title: '范围 Key', dataIndex: 'scopeKey', width: 180 },
    { title: '方向', dataIndex: 'direction', width: 160 },
    {
      title: '启用',
      dataIndex: 'enabled',
      width: 100,
      render: (value) => (value ? <Tag color="green">启用</Tag> : <Tag>停用</Tag>),
    },
  ];

  const outboxColumns: ColumnsType<ClusterOutbox> = [
    { title: '事件 ID', dataIndex: 'eventId', width: 240, ellipsis: true },
    { title: '事件类型', dataIndex: 'eventType', width: 200 },
    { title: '聚合', dataIndex: 'aggregateType', width: 120 },
    { title: '聚合 ID', dataIndex: 'aggregateId', width: 140 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (value) => <Tag color={value === 'dead' ? 'red' : 'orange'}>{value}</Tag>,
    },
    { title: '重试', dataIndex: 'attempts', width: 80 },
    { title: '错误', dataIndex: 'lastError', ellipsis: true, render: (value) => value || '-' },
  ];

  return (
    <ComLayout loading={loading}>
      <ComContent
        hasBack={false}
        title={
          <PageTitleRow resourceKey="cluster.edge.manage">
            <span>{title || 'Cluster 边边协同'}</span>
          </PageTitleRow>
        }
        extra={
          <Space>
            <Button icon={<Renew size={16} />} onClick={() => void load()}>
              刷新
            </Button>
            <Button icon={<Connect size={16} />} onClick={() => setHandshakeOpen(true)}>
              握手测试
            </Button>
            <Button type="primary" icon={<Add size={16} />} onClick={() => setTokenOpen(true)}>
              创建 Token
            </Button>
          </Space>
        }
      >
        <div className={styles['cluster-page']}>
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            tabBarExtraContent={
              activeTab === 'nodes' ? (
                <Button size="small" icon={<Renew size={16} />} loading={syncing} onClick={forceSync}>
                  {formatMessage('cluster.forceSyncRefresh', {}, 'Force Sync / Refresh')}
                </Button>
              ) : null
            }
            items={[
              {
                key: 'nodes',
                label: 'Edge 节点',
                children: (
                  <Table
                    rowKey="id"
                    size="small"
                    dataSource={nodes}
                    columns={nodeColumns}
                    pagination={false}
                    scroll={{ x: 1100 }}
                  />
                ),
              },
              {
                key: 'scopes',
                label: '同步范围',
                children: (
                  <Table
                    rowKey="id"
                    size="small"
                    dataSource={scopes}
                    columns={scopeColumns}
                    pagination={{ pageSize: 20 }}
                    scroll={{ x: 800 }}
                  />
                ),
              },
              {
                key: 'outbox',
                label: '待同步事件',
                children: (
                  <Table
                    rowKey="id"
                    size="small"
                    dataSource={outbox}
                    columns={outboxColumns}
                    pagination={{ pageSize: 20 }}
                    scroll={{ x: 1200 }}
                  />
                ),
              },
            ]}
          />
        </div>

        <Modal
          title="创建协同 Token"
          open={tokenOpen}
          onOk={createToken}
          onCancel={() => setTokenOpen(false)}
          width={620}
          destroyOnClose
        >
          <Form form={tokenForm} layout="vertical" preserve={false}>
            <Form.Item name="tokenName" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
              <Input placeholder="edge-a-token" />
            </Form.Item>
            <Form.Item name="description" label="描述">
              <Input.TextArea rows={4} placeholder="请输入描述" />
            </Form.Item>
          </Form>
          {createdToken && (
            <Typography.Paragraph copyable className="cluster-token">
              {createdToken}
            </Typography.Paragraph>
          )}
        </Modal>

        <Modal
          title="握手测试"
          open={handshakeOpen}
          onOk={handshake}
          onCancel={() => setHandshakeOpen(false)}
          width={620}
          destroyOnClose
        >
          <Form
            form={handshakeForm}
            layout="vertical"
            preserve={false}
            initialValues={{ role: 'edge', capabilitiesJson: '{}' }}
          >
            <Form.Item name="token" label="Token" rules={[{ required: true, message: '请输入 token' }]}>
              <Input.Password />
            </Form.Item>
            <Form.Item name="nodeKey" label="节点 Key" rules={[{ required: true, message: '请输入节点 Key' }]}>
              <Input placeholder="edge-a" />
            </Form.Item>
            <Form.Item name="nodeName" label="节点名称">
              <Input />
            </Form.Item>
            <Form.Item name="role" label="角色">
              <Input />
            </Form.Item>
            <Form.Item name="endpoint" label="Endpoint">
              <Input />
            </Form.Item>
            <Form.Item name="capabilitiesJson" label="能力 JSON">
              <Input.TextArea rows={4} />
            </Form.Item>
          </Form>
        </Modal>
      </ComContent>
    </ComLayout>
  );
};

export default ClusterPage;
