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
import ProSearch from '@/components/pro-search';
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
} from '@/apis/core-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';
import useTranslate from '@/hooks/useTranslate';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';

const JWT_ALGORITHMS = ['HS256', 'HS384', 'HS512', 'RS256', 'RS384', 'RS512', 'ES256', 'ES384', 'ES512'];

const ConsumersTab: FC = () => {
  const { modal, message } = App.useApp();
  const formatMessage = useTranslate();
  const { data, loading, refresh } = useKongTable({ fetchApi: getConsumers });
  const [search, setSearch] = useState('');
  const [detailRecord, setDetailRecord] = useState<any>(null);

  const CRED_TYPES = [
    { key: 'basic-auth', label: formatMessage('kong.credTypeBasic') },
    { key: 'key-auth', label: formatMessage('kong.credTypeApiKeys') },
    { key: 'hmac-auth', label: formatMessage('kong.credTypeHmac') },
    { key: 'oauth2', label: formatMessage('kong.credTypeOAuth2') },
    { key: 'jwt', label: formatMessage('kong.credTypeJwt') },
  ];

  const CRED_LABELS: Record<string, string> = {
    'basic-auth': formatMessage('kong.credLabelBasic'),
    'key-auth': formatMessage('kong.credLabelApiKeys'),
    'hmac-auth': formatMessage('kong.credLabelHmac'),
    oauth2: formatMessage('kong.credLabelOAuth2'),
    jwt: formatMessage('kong.credLabelJwt'),
  };

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
      modal.confirm(
        mergeDeleteConfirmProps(
          {
            title: formatMessage('kong.deleteCredential'),
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
          },
          formatMessage
        )
      );
    },
    [modal, detailRecord, credType, fetchCredentials, formatMessage]
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
      message.success(formatMessage('kong.credentialCreated'));
      setCredModalOpen(false);
      credForm.resetFields();
      fetchCredentials(detailRecord.id, credType);
    } catch (e: any) {
      if (!e?.errorFields) {
        message.error(e?.message || formatMessage('kong.credentialFailed'));
      }
    } finally {
      setSavingCred(false);
    }
  }, [credForm, credType, detailRecord, fetchCredentials, message, formatMessage]);

  const credColumns = useMemo(() => {
    const actionsCol = {
      title: '',
      width: 60,
      render: (_: any, r: any) => (
        <Button type="text" size="small" danger icon={<TrashCan size={14} />} onClick={() => handleDeleteCred(r.id)} />
      ),
    };
    const idCol = {
      title: formatMessage('kong.colID'),
      dataIndex: 'id',
      width: 110,
      ellipsis: true,
      render: (v: string) => <Typography.Text copyable={{ text: v }}>{v?.slice(0, 8)}…</Typography.Text>,
    };

    switch (credType) {
      case 'basic-auth':
        return [idCol, { title: formatMessage('kong.colUsername'), dataIndex: 'username', ellipsis: true }, actionsCol];
      case 'key-auth':
        return [
          idCol,
          {
            title: formatMessage('kong.colKey'),
            dataIndex: 'key',
            ellipsis: true,
            render: (v: string) => <Typography.Text copyable>{v}</Typography.Text>,
          },
          actionsCol,
        ];
      case 'hmac-auth':
        return [
          idCol,
          { title: formatMessage('kong.colUsername'), dataIndex: 'username', ellipsis: true },
          { title: formatMessage('kong.colSecret'), dataIndex: 'secret', ellipsis: true },
          actionsCol,
        ];
      case 'oauth2':
        return [
          idCol,
          { title: formatMessage('kong.colName'), dataIndex: 'name', ellipsis: true },
          { title: formatMessage('kong.colClientID'), dataIndex: 'client_id', ellipsis: true },
          {
            title: formatMessage('kong.colRedirectURIs'),
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
            title: formatMessage('kong.colKey'),
            dataIndex: 'key',
            ellipsis: true,
            render: (v: string) => <Typography.Text copyable>{v}</Typography.Text>,
          },
          { title: formatMessage('kong.colAlgorithm'), dataIndex: 'algorithm', width: 100 },
          actionsCol,
        ];
      default:
        return [];
    }
  }, [credType, handleDeleteCred, formatMessage]);

  const renderCredForm = useCallback(() => {
    switch (credType) {
      case 'basic-auth':
        return (
          <>
            <Form.Item name="username" label={formatMessage('kong.labelUsername')} rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="password" label={formatMessage('kong.labelPassword')} rules={[{ required: true }]}>
              <Input.Password />
            </Form.Item>
          </>
        );
      case 'key-auth':
        return (
          <Form.Item name="key" label={formatMessage('kong.labelKey')} extra={formatMessage('kong.phAutoGenerate')}>
            <Input placeholder={formatMessage('kong.phAutoGenerateInput')} />
          </Form.Item>
        );
      case 'hmac-auth':
        return (
          <>
            <Form.Item name="username" label={formatMessage('kong.labelUsername')} rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item
              name="secret"
              label={formatMessage('kong.labelSecret')}
              extra={formatMessage('kong.phAutoGenerate')}
            >
              <Input placeholder={formatMessage('kong.phAutoGenerateInput')} />
            </Form.Item>
          </>
        );
      case 'oauth2':
        return (
          <>
            <Form.Item name="name" label={formatMessage('kong.colName')} rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item
              name="client_id"
              label={formatMessage('kong.labelClientID')}
              extra={formatMessage('kong.phAutoGenerate')}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="client_secret"
              label={formatMessage('kong.labelClientSecret')}
              extra={formatMessage('kong.phAutoGenerate')}
            >
              <Input />
            </Form.Item>
            <Form.Item name="redirect_uris" label={formatMessage('kong.labelRedirectURIs')}>
              <Input placeholder={formatMessage('kong.phOAuth2Redirect')} />
            </Form.Item>
          </>
        );
      case 'jwt':
        return (
          <>
            <Form.Item name="algorithm" label={formatMessage('kong.labelAlgorithm')} initialValue="HS256">
              <Select
                options={JWT_ALGORITHMS.map((a) => ({ label: a, value: a }))}
                onChange={(v) => setJwtAlgorithm(v)}
              />
            </Form.Item>
            <Form.Item name="key" label={formatMessage('kong.labelKey')} extra={formatMessage('kong.phAutoGenerate')}>
              <Input />
            </Form.Item>
            {jwtAlgorithm.startsWith('HS') && (
              <Form.Item
                name="secret"
                label={formatMessage('kong.labelSecret')}
                extra={formatMessage('kong.phAutoGenerate')}
              >
                <Input />
              </Form.Item>
            )}
            {(jwtAlgorithm.startsWith('RS') || jwtAlgorithm.startsWith('ES')) && (
              <Form.Item name="rsa_public_key" label={formatMessage('kong.labelPublicKey')}>
                <Input.TextArea rows={4} />
              </Form.Item>
            )}
          </>
        );
      default:
        return null;
    }
  }, [credType, jwtAlgorithm, formatMessage]);

  const renderForm = useCallback(
    () => (
      <>
        <Form.Item
          name="username"
          label={formatMessage('kong.labelUsername')}
          rules={[{ required: true, message: formatMessage('kong.ruleUsernameRequired') }]}
        >
          <Input placeholder={formatMessage('kong.phConsumerName')} />
        </Form.Item>
        <Form.Item name="custom_id" label={formatMessage('kong.labelCustomID')}>
          <Input placeholder={formatMessage('kong.phCustomID')} />
        </Form.Item>
        <Form.Item name="tags" label={formatMessage('kong.labelTags')}>
          <Input placeholder={formatMessage('kong.phTags')} />
        </Form.Item>
      </>
    ),
    [formatMessage]
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
    title: formatMessage('kong.consumers'),
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
      modal.confirm(
        mergeDeleteConfirmProps(
          {
            title: formatMessage('kong.deleteConsumer', { name: record.username ?? record.id }),
            onOk: async () => {
              await deleteConsumer(record.id);
              refresh();
            },
          },
          formatMessage
        )
      );
    },
    [modal, refresh, formatMessage]
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
      title: formatMessage('kong.colUsername'),
      dataIndex: 'username',
      width: 200,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => handleOpenDetail(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: formatMessage('kong.colCustomID'),
      dataIndex: 'custom_id',
      width: 200,
      ellipsis: true,
      render: (v: string) => v || '-',
    },
    {
      title: formatMessage('kong.colTags'),
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
      title: formatMessage('kong.colCreated'),
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
        label: formatMessage('kong.details'),
        children: detailRecord && (
          <Descriptions column={1} bordered size="small" className="detail-descriptions">
            <Descriptions.Item label={formatMessage('kong.colID')}>{detailRecord.id}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colUsername')}>
              {detailRecord.username || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colCustomID')}>
              {detailRecord.custom_id || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colTags')}>
              {detailRecord.tags?.join(', ') || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colCreated')}>
              {detailRecord.created_at ? new Date(detailRecord.created_at * 1000).toLocaleString() : '-'}
            </Descriptions.Item>
          </Descriptions>
        ),
      },
      {
        key: 'credentials',
        label: formatMessage('kong.credentials'),
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
                  {formatMessage('kong.createCredentials')}
                </Button>
              </Flex>
              {credLoading ? (
                <Flex justify="center" style={{ padding: 32 }}>
                  <Spin />
                </Flex>
              ) : credentials.length === 0 ? (
                <Typography.Text type="secondary">
                  {formatMessage('kong.noCredentials', { label: CRED_LABELS[credType] })}
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
    [
      detailRecord,
      credType,
      handleCredTypeChange,
      credForm,
      credLoading,
      credentials,
      credColumns,
      formatMessage,
      CRED_TYPES,
      CRED_LABELS,
    ]
  );

  return (
    <>
      <div className="toolbar">
        <div className="toolbar-left">
          <Button type="primary" icon={<Add size={16} />} onClick={() => open()}>
            {formatMessage('kong.addConsumer')}
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            {formatMessage('common.refresh')}
          </Button>
        </div>
        <div className="toolbar-right">
          <ProSearch
            size="sm"
            placeholder={formatMessage('kong.searchConsumer')}
            style={{ width: 280 }}
            onChange={(e) => {
              if (!e.target.value) {
                setSearch('');
              }
            }}
            onSearch={setSearch}
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
                    {formatMessage('common.edit')} <Edit size={14} />
                  </span>
                ),
                onClick: () => handleEdit(record),
              },
              {
                key: 'delete',
                label: (
                  <span style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 24 }}>
                    {formatMessage('common.delete')} <TrashCan size={14} />
                  </span>
                ),
                onClick: () => handleDelete(record),
              },
            ],
          }}
        />
      </div>
      {ModalDom}
      <Drawer
        title={formatMessage('kong.consumerDetail')}
        open={!!detailRecord}
        onClose={handleCloseDetail}
        width={1000}
      >
        <Tabs items={drawerTabs} />
      </Drawer>
      <Modal
        title={formatMessage('kong.createCredentialTitle', { label: CRED_LABELS[credType] })}
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
