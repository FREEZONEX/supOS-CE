import { type FC, useCallback, useMemo, useState } from 'react';
import {
  Button,
  Flex,
  App,
  Input,
  Descriptions,
  Drawer,
  Form,
  Tag,
  Tabs,
  Table,
  Menu,
  Modal,
  Typography,
  Spin,
  Select,
} from 'antd';
import { Add, Renew, TrashCan, Edit } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import {
  getConsumers,
  createConsumer,
  updateConsumer,
  deleteConsumer,
  getConsumerBasicAuth,
  createConsumerBasicAuth,
  deleteConsumerBasicAuth,
  getConsumerKeyAuth,
  createConsumerKeyAuth,
  deleteConsumerKeyAuth,
  getConsumerHmacAuth,
  createConsumerHmacAuth,
  deleteConsumerHmacAuth,
  getConsumerOAuth2,
  createConsumerOAuth2,
  deleteConsumerOAuth2,
  getConsumerJwt,
  createConsumerJwt,
  deleteConsumerJwt,
} from '@/apis/inter-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';

const CRED_TYPES = [
  { key: 'basic-auth', label: 'BASIC' },
  { key: 'key-auth', label: 'API KEYS' },
  { key: 'hmac-auth', label: 'HMAC' },
  { key: 'oauth2', label: 'OAUTH2' },
  { key: 'jwt', label: 'JWT' },
];

const CRED_LABELS: Record<string, string> = {
  'basic-auth': 'Basic Auth',
  'key-auth': 'API Keys',
  'hmac-auth': 'HMAC Auth',
  oauth2: 'OAuth 2.0',
  jwt: 'JWT',
};

const JWT_ALGORITHMS = ['HS256', 'HS384', 'HS512', 'RS256', 'RS384', 'RS512', 'ES256', 'ES384', 'ES512'];

const ConsumersTab: FC = () => {
  const { modal, message } = App.useApp();
  const { data, loading, refresh } = useKongTable({ fetchApi: getConsumers });
  const [search, setSearch] = useState('');
  const [detailRecord, setDetailRecord] = useState<any>(null);

  const [credType, setCredType] = useState('basic-auth');
  const [credentials, setCredentials] = useState<any[]>([]);
  const [credLoading, setCredLoading] = useState(false);
  const [credModalOpen, setCredModalOpen] = useState(false);
  const [credForm] = Form.useForm();
  const [savingCred, setSavingCred] = useState(false);
  const [jwtAlgorithm, setJwtAlgorithm] = useState('HS256');

  const fetchCredentials = useCallback((consumerId: string, type: string) => {
    setCredLoading(true);
    let req: Promise<any>;
    switch (type) {
      case 'basic-auth':
        req = getConsumerBasicAuth(consumerId);
        break;
      case 'key-auth':
        req = getConsumerKeyAuth(consumerId);
        break;
      case 'hmac-auth':
        req = getConsumerHmacAuth(consumerId);
        break;
      case 'oauth2':
        req = getConsumerOAuth2(consumerId);
        break;
      case 'jwt':
        req = getConsumerJwt(consumerId);
        break;
      default:
        req = Promise.resolve({ data: [] });
    }
    req
      .then((res: any) => setCredentials(res?.data ?? []))
      .catch(() => setCredentials([]))
      .finally(() => setCredLoading(false));
  }, []);

  const handleOpenDetail = useCallback(
    (record: any) => {
      setDetailRecord(record);
      setCredType('basic-auth');
      setCredentials([]);
      fetchCredentials(record.id, 'basic-auth');
    },
    [fetchCredentials]
  );

  const handleCloseDetail = useCallback(() => {
    setDetailRecord(null);
    setCredentials([]);
  }, []);

  const handleCredTypeChange = useCallback(
    ({ key }: { key: string }) => {
      setCredType(key);
      if (detailRecord?.id) fetchCredentials(detailRecord.id, key);
    },
    [detailRecord, fetchCredentials]
  );

  const handleDeleteCred = useCallback(
    (credId: string) => {
      if (!detailRecord?.id) return;
      modal.confirm({
        title: 'Delete this credential?',
        okButtonProps: { danger: true },
        onOk: async () => {
          switch (credType) {
            case 'basic-auth':
              await deleteConsumerBasicAuth(detailRecord.id, credId);
              break;
            case 'key-auth':
              await deleteConsumerKeyAuth(detailRecord.id, credId);
              break;
            case 'hmac-auth':
              await deleteConsumerHmacAuth(detailRecord.id, credId);
              break;
            case 'oauth2':
              await deleteConsumerOAuth2(detailRecord.id, credId);
              break;
            case 'jwt':
              await deleteConsumerJwt(detailRecord.id, credId);
              break;
          }
          fetchCredentials(detailRecord.id, credType);
        },
      });
    },
    [modal, detailRecord, credType, fetchCredentials]
  );

  const handleCreateCred = useCallback(async () => {
    if (!detailRecord?.id) return;
    try {
      const values = await credForm.validateFields();
      setSavingCred(true);
      const payload: Record<string, unknown> = {};
      switch (credType) {
        case 'basic-auth':
          await createConsumerBasicAuth(detailRecord.id, { username: values.username, password: values.password });
          break;
        case 'key-auth':
          if (values.key) payload.key = values.key;
          await createConsumerKeyAuth(detailRecord.id, payload);
          break;
        case 'hmac-auth':
          payload.username = values.username;
          if (values.secret) payload.secret = values.secret;
          await createConsumerHmacAuth(detailRecord.id, payload);
          break;
        case 'oauth2':
          payload.name = values.name;
          if (values.client_id) payload.client_id = values.client_id;
          if (values.client_secret) payload.client_secret = values.client_secret;
          if (values.redirect_uris) {
            payload.redirect_uris = values.redirect_uris
              .split(',')
              .map((s: string) => s.trim())
              .filter(Boolean);
          }
          await createConsumerOAuth2(detailRecord.id, payload);
          break;
        case 'jwt':
          payload.algorithm = values.algorithm || 'HS256';
          if (values.key) payload.key = values.key;
          if (values.secret) payload.secret = values.secret;
          if (values.rsa_public_key) payload.rsa_public_key = values.rsa_public_key;
          await createConsumerJwt(detailRecord.id, payload);
          break;
      }
      message.success('Credential created successfully');
      setCredModalOpen(false);
      credForm.resetFields();
      fetchCredentials(detailRecord.id, credType);
    } catch (e: any) {
      if (!e?.errorFields) {
        message.error(e?.message || 'Failed to create credential');
      }
    } finally {
      setSavingCred(false);
    }
  }, [credForm, credType, detailRecord, fetchCredentials, message]);

  const credColumns = useMemo(() => {
    const actionsCol = {
      title: '',
      width: 60,
      render: (_: any, r: any) => (
        <Button type="text" size="small" danger icon={<TrashCan size={14} />} onClick={() => handleDeleteCred(r.id)} />
      ),
    };
    const idCol = {
      title: 'ID',
      dataIndex: 'id',
      width: 110,
      ellipsis: true,
      render: (v: string) => <Typography.Text copyable={{ text: v }}>{v?.slice(0, 8)}…</Typography.Text>,
    };

    switch (credType) {
      case 'basic-auth':
        return [idCol, { title: 'Username', dataIndex: 'username', ellipsis: true }, actionsCol];
      case 'key-auth':
        return [
          idCol,
          {
            title: 'Key',
            dataIndex: 'key',
            ellipsis: true,
            render: (v: string) => <Typography.Text copyable>{v}</Typography.Text>,
          },
          actionsCol,
        ];
      case 'hmac-auth':
        return [
          idCol,
          { title: 'Username', dataIndex: 'username', ellipsis: true },
          { title: 'Secret', dataIndex: 'secret', ellipsis: true },
          actionsCol,
        ];
      case 'oauth2':
        return [
          idCol,
          { title: 'Name', dataIndex: 'name', ellipsis: true },
          { title: 'Client ID', dataIndex: 'client_id', ellipsis: true },
          {
            title: 'Redirect URIs',
            dataIndex: 'redirect_uris',
            ellipsis: true,
            render: (v: string[]) => v?.join(', ') || '-',
          },
          actionsCol,
        ];
      case 'jwt':
        return [
          idCol,
          {
            title: 'Key',
            dataIndex: 'key',
            ellipsis: true,
            render: (v: string) => <Typography.Text copyable>{v}</Typography.Text>,
          },
          { title: 'Algorithm', dataIndex: 'algorithm', width: 100 },
          actionsCol,
        ];
      default:
        return [];
    }
  }, [credType, handleDeleteCred]);

  const renderCredForm = useCallback(() => {
    switch (credType) {
      case 'basic-auth':
        return (
          <>
            <Form.Item name="username" label="Username" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="password" label="Password" rules={[{ required: true }]}>
              <Input.Password />
            </Form.Item>
          </>
        );
      case 'key-auth':
        return (
          <Form.Item name="key" label="Key" extra="Leave empty to auto-generate">
            <Input placeholder="Auto-generated if empty" />
          </Form.Item>
        );
      case 'hmac-auth':
        return (
          <>
            <Form.Item name="username" label="Username" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="secret" label="Secret" extra="Leave empty to auto-generate">
              <Input placeholder="Auto-generated if empty" />
            </Form.Item>
          </>
        );
      case 'oauth2':
        return (
          <>
            <Form.Item name="name" label="Name" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="client_id" label="Client ID" extra="Leave empty to auto-generate">
              <Input />
            </Form.Item>
            <Form.Item name="client_secret" label="Client Secret" extra="Leave empty to auto-generate">
              <Input />
            </Form.Item>
            <Form.Item name="redirect_uris" label="Redirect URIs">
              <Input placeholder="https://example.com/callback (comma separated)" />
            </Form.Item>
          </>
        );
      case 'jwt':
        return (
          <>
            <Form.Item name="algorithm" label="Algorithm" initialValue="HS256">
              <Select
                options={JWT_ALGORITHMS.map((a) => ({ label: a, value: a }))}
                onChange={(v) => setJwtAlgorithm(v)}
              />
            </Form.Item>
            <Form.Item name="key" label="Key" extra="Leave empty to auto-generate">
              <Input />
            </Form.Item>
            {jwtAlgorithm.startsWith('HS') && (
              <Form.Item name="secret" label="Secret" extra="Leave empty to auto-generate">
                <Input />
              </Form.Item>
            )}
            {(jwtAlgorithm.startsWith('RS') || jwtAlgorithm.startsWith('ES')) && (
              <Form.Item name="rsa_public_key" label="Public Key">
                <Input.TextArea rows={4} />
              </Form.Item>
            )}
          </>
        );
      default:
        return null;
    }
  }, [credType, jwtAlgorithm]);

  const renderForm = useCallback(
    () => (
      <>
        <Form.Item name="username" label="Username" rules={[{ required: true, message: 'Username is required' }]}>
          <Input placeholder="my-consumer" />
        </Form.Item>
        <Form.Item name="custom_id" label="Custom ID">
          <Input placeholder="custom-id-123" />
        </Form.Item>
        <Form.Item name="tags" label="Tags">
          <Input placeholder="tag1, tag2 (comma separated)" />
        </Form.Item>
      </>
    ),
    []
  );

  const transformValues = useCallback((values: any) => {
    const payload: Record<string, unknown> = { ...values };
    if (typeof values.tags === 'string') {
      payload.tags = values.tags
        .split(',')
        .map((t: string) => t.trim())
        .filter(Boolean);
      if ((payload.tags as string[]).length === 0) delete payload.tags;
    }
    return payload;
  }, []);

  const { ModalDom, open } = useKongModal({
    title: 'Consumer',
    createApi: createConsumer,
    updateApi: updateConsumer,
    onSuccess: refresh,
    renderForm,
    transformValues,
  });

  const handleEdit = useCallback(
    (record: any) => {
      open({ ...record, tags: record.tags?.join(', ') ?? '' });
    },
    [open]
  );

  const handleDelete = useCallback(
    (record: any) => {
      modal.confirm({
        title: `Delete consumer "${record.username ?? record.id}"?`,
        okButtonProps: { danger: true },
        onOk: async () => {
          await deleteConsumer(record.id);
          refresh();
        },
      });
    },
    [modal, refresh]
  );

  const filteredData = useMemo(() => {
    if (!search) return data;
    const q = search.toLowerCase();
    return data.filter(
      (c: any) =>
        c.username?.toLowerCase().includes(q) ||
        c.custom_id?.toLowerCase().includes(q) ||
        c.id?.toLowerCase().includes(q)
    );
  }, [data, search]);

  const columns = [
    {
      title: 'Username',
      dataIndex: 'username',
      width: 200,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => handleOpenDetail(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: 'Custom ID',
      dataIndex: 'custom_id',
      width: 200,
      ellipsis: true,
      render: (v: string) => v || '-',
    },
    {
      title: 'Tags',
      dataIndex: 'tags',
      width: 200,
      render: (v: string[]) =>
        v?.length ? (
          <div className="tag-list">
            {v.map((t) => (
              <Tag key={t}>{t}</Tag>
            ))}
          </div>
        ) : (
          '-'
        ),
    },
    {
      title: 'Created',
      dataIndex: 'created_at',
      width: 180,
      sorter: (a: any, b: any) => (a.created_at ?? 0) - (b.created_at ?? 0),
      render: (v: number) => (v ? new Date(v * 1000).toLocaleString() : '-'),
    },
  ];

  const drawerTabs = useMemo(
    () => [
      {
        key: 'details',
        label: 'Details',
        children: detailRecord && (
          <Descriptions column={1} bordered size="small" className="detail-descriptions">
            <Descriptions.Item label="ID">{detailRecord.id}</Descriptions.Item>
            <Descriptions.Item label="Username">{detailRecord.username || '-'}</Descriptions.Item>
            <Descriptions.Item label="Custom ID">{detailRecord.custom_id || '-'}</Descriptions.Item>
            <Descriptions.Item label="Tags">{detailRecord.tags?.join(', ') || '-'}</Descriptions.Item>
            <Descriptions.Item label="Created">
              {detailRecord.created_at ? new Date(detailRecord.created_at * 1000).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        ),
      },
      {
        key: 'credentials',
        label: 'Credentials',
        children: (
          <Flex style={{ minHeight: 320 }}>
            <Menu
              mode="inline"
              selectedKeys={[credType]}
              onClick={handleCredTypeChange}
              items={CRED_TYPES.map((ct) => ({ key: ct.key, label: ct.label }))}
              className="routing-sidebar-menu"
              style={{ width: 160, flexShrink: 0, borderRight: '1px solid var(--ant-color-border)' }}
            />
            <div style={{ flex: 1, paddingLeft: 16, overflow: 'auto' }}>
              <Flex justify="space-between" align="center" style={{ marginBottom: 16 }}>
                <Typography.Title level={5} style={{ margin: 0 }}>
                  {CRED_LABELS[credType]}
                </Typography.Title>
                <Button
                  type="primary"
                  size="small"
                  icon={<Add size={14} />}
                  onClick={() => {
                    credForm.resetFields();
                    setJwtAlgorithm('HS256');
                    setCredModalOpen(true);
                  }}
                >
                  Create Credentials
                </Button>
              </Flex>
              {credLoading ? (
                <Flex justify="center" style={{ padding: 32 }}>
                  <Spin />
                </Flex>
              ) : credentials.length === 0 ? (
                <Typography.Text type="secondary">
                  You have not created any {CRED_LABELS[credType]} credentials for this consumer yet
                </Typography.Text>
              ) : (
                <Table
                  rowKey="id"
                  size="small"
                  dataSource={credentials}
                  columns={credColumns}
                  pagination={false}
                  scroll={{ x: 'max-content' }}
                />
              )}
            </div>
          </Flex>
        ),
      },
    ],
    [detailRecord, credType, handleCredTypeChange, credForm, credLoading, credentials, credColumns]
  );

  return (
    <>
      <div className="toolbar">
        <div className="toolbar-left">
          <Button type="primary" icon={<Add size={16} />} onClick={() => open()}>
            Add Consumer
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            Refresh
          </Button>
        </div>
        <div className="toolbar-right">
          <Input.Search
            placeholder="Search by username"
            allowClear
            style={{ width: 280 }}
            onSearch={setSearch}
            onChange={(e) => !e.target.value && setSearch('')}
          />
        </div>
      </div>
      <div className="table-area">
        <ProTable
          rowKey="id"
          size="small"
          loading={loading}
          dataSource={filteredData}
          columns={columns}
          pagination={{ pageSize: 20, showTotal: (t: number) => `Total ${t}`, showQuickJumper: true }}
          scroll={{ y: 'calc(100vh - 285px)', x: 'max-content' }}
          operationOptions={{
            render: (record) => [
              {
                key: 'edit',
                label: (
                  <span style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 24 }}>
                    Edit <Edit size={14} />
                  </span>
                ),
                onClick: () => handleEdit(record),
              },
              {
                key: 'delete',
                label: (
                  <span style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 24 }}>
                    Delete <TrashCan size={14} />
                  </span>
                ),
                onClick: () => handleDelete(record),
              },
            ],
          }}
        />
      </div>
      {ModalDom}
      <Drawer title="Consumer Detail" open={!!detailRecord} onClose={handleCloseDetail} width={1000}>
        <Tabs items={drawerTabs} />
      </Drawer>
      <Modal
        title={`Create ${CRED_LABELS[credType]} Credential`}
        open={credModalOpen}
        onCancel={() => {
          setCredModalOpen(false);
          credForm.resetFields();
        }}
        onOk={handleCreateCred}
        confirmLoading={savingCred}
      >
        <Form form={credForm} layout="vertical" style={{ marginTop: 16 }}>
          {renderCredForm()}
        </Form>
      </Modal>
    </>
  );
};

export default ConsumersTab;
