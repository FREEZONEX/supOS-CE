import { type FC, useCallback, useMemo, useState, useEffect } from 'react';
import {
  Button,
  Flex,
  App,
  Tag,
  Input,
  Descriptions,
  Drawer,
  Tabs,
  Table,
  Typography,
  Modal,
  Form,
  Select,
  InputNumber,
} from 'antd';
import { Add, Renew, Edit, TrashCan } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import {
  getServices,
  createService,
  updateService,
  deleteService,
  getServiceRoutes,
  createRoute,
  updateRoute,
  deleteRoute,
} from '@/apis/inter-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';

const SVC_PROTOCOLS = ['http', 'https', 'grpc', 'grpcs', 'tcp', 'tls', 'udp'];
const ROUTE_PROTOCOLS = ['http', 'https', 'grpc', 'grpcs', 'tcp', 'tls', 'udp'];
const ROUTE_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'TRACE'];

interface ServicesTabProps {
  onViewRoute?: (route: any) => void;
}

const ServicesTab: FC<ServicesTabProps> = ({ onViewRoute }) => {
  const { modal, message } = App.useApp();
  const { data, loading, refresh } = useKongTable({ fetchApi: getServices });
  const [search, setSearch] = useState('');
  const [detailRecord, setDetailRecord] = useState<any>(null);
  const [serviceRoutes, setServiceRoutes] = useState<any[]>([]);

  const [addRouteOpen, setAddRouteOpen] = useState(false);
  const [addRouteForm] = Form.useForm();
  const [savingRoute, setSavingRoute] = useState(false);
  const [editingRoute, setEditingRoute] = useState<any>(null);

  const refreshRoutes = useCallback((serviceId: string) => {
    getServiceRoutes(serviceId)
      .then((res: any) => setServiceRoutes(res?.data ?? []))
      .catch(() => setServiceRoutes([]));
  }, []);

  useEffect(() => {
    if (!detailRecord?.id) return;
    refreshRoutes(detailRecord.id);
  }, [detailRecord?.id, refreshRoutes]);

  const handleAddRoute = useCallback(async () => {
    if (!detailRecord?.id) return;
    try {
      const values = await addRouteForm.validateFields();
      setSavingRoute(true);
      const payload: Record<string, unknown> = {
        protocols: values.protocols,
        methods: values.methods?.length ? values.methods : undefined,
        paths: values.paths
          ? values.paths
              .split(',')
              .map((p: string) => p.trim())
              .filter(Boolean)
          : undefined,
        strip_path: values.strip_path,
        preserve_host: values.preserve_host,
      };
      if (values.name) payload.name = values.name;
      if (values.hosts) {
        const hosts = values.hosts
          .split(',')
          .map((h: string) => h.trim())
          .filter(Boolean);
        if (hosts.length) payload.hosts = hosts;
      }
      if (editingRoute) {
        await updateRoute(editingRoute.id, payload);
        message.success('Route updated');
      } else {
        payload.service = { id: detailRecord.id };
        await createRoute(payload);
        message.success('Route created');
      }
      setAddRouteOpen(false);
      setEditingRoute(null);
      addRouteForm.resetFields();
      refreshRoutes(detailRecord.id);
    } catch (e: any) {
      if (!e?.errorFields) message.error(e?.message || 'Failed to save route');
    } finally {
      setSavingRoute(false);
    }
  }, [detailRecord, addRouteForm, refreshRoutes, message, editingRoute]);

  const handleEditRoute = useCallback(
    (route: any) => {
      addRouteForm.setFieldsValue({
        ...route,
        paths: Array.isArray(route.paths) ? route.paths.join(', ') : route.paths,
        hosts: Array.isArray(route.hosts) ? route.hosts.join(', ') : route.hosts,
      });
      setEditingRoute(route);
      setAddRouteOpen(true);
    },
    [addRouteForm]
  );

  const handleDeleteRoute = useCallback(
    (route: any) => {
      modal.confirm({
        title: `Delete route "${route.name ?? route.id}"?`,
        okButtonProps: { danger: true },
        onOk: async () => {
          await deleteRoute(route.id);
          if (detailRecord?.id) refreshRoutes(detailRecord.id);
        },
      });
    },
    [modal, detailRecord, refreshRoutes]
  );

  const renderForm = useCallback(
    (_form: any, editing: any) => (
      <>
        <Form.Item name="name" label="Name" rules={[{ required: true, message: 'Service name is required' }]}>
          <Input placeholder="my-service" />
        </Form.Item>
        <Flex gap={12}>
          <Form.Item name="protocol" label="Protocol" initialValue="http" style={{ flex: 1 }}>
            <Select options={SVC_PROTOCOLS.map((p) => ({ label: p, value: p }))} />
          </Form.Item>
          <Form.Item
            name="host"
            label="Host"
            rules={[{ required: true, message: 'Host is required' }]}
            style={{ flex: 2 }}
          >
            <Input placeholder="example.com" />
          </Form.Item>
          <Form.Item name="port" label="Port" initialValue={80} style={{ flex: 1 }}>
            <InputNumber min={1} max={65535} style={{ width: '100%' }} />
          </Form.Item>
        </Flex>
        <Form.Item name="path" label="Path">
          <Input placeholder="/" />
        </Form.Item>
        <Flex gap={12}>
          <Form.Item name="retries" label="Retries" initialValue={5} style={{ flex: 1 }}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="connect_timeout" label="Connect Timeout (ms)" initialValue={60000} style={{ flex: 1 }}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="write_timeout" label="Write Timeout (ms)" initialValue={60000} style={{ flex: 1 }}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="read_timeout" label="Read Timeout (ms)" initialValue={60000} style={{ flex: 1 }}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Flex>
        {editing?.id && (
          <Form.Item name="enabled" label="Enabled" initialValue={true}>
            <Select
              options={[
                { label: 'Enabled', value: true },
                { label: 'Disabled', value: false },
              ]}
            />
          </Form.Item>
        )}
      </>
    ),
    []
  );

  const { ModalDom, open } = useKongModal({
    title: 'Service',
    createApi: createService,
    updateApi: updateService,
    onSuccess: refresh,
    renderForm,
    width: 720,
  });

  const handleDelete = useCallback(
    (record: any) => {
      modal.confirm({
        title: `Delete service "${record.name ?? record.id}"?`,
        content: 'This will also remove associated routes and plugins.',
        okButtonProps: { danger: true },
        onOk: async () => {
          await deleteService(record.id);
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
      (s: any) =>
        s.name?.toLowerCase().includes(q) || s.host?.toLowerCase().includes(q) || s.id?.toLowerCase().includes(q)
    );
  }, [data, search]);

  const routeColumns = useMemo(
    () => [
      {
        title: 'Name',
        dataIndex: 'name',
        ellipsis: true,
        render: (v: string, record: any) => (
          <Typography.Link
            onClick={() => {
              setDetailRecord(null);
              onViewRoute?.(record);
            }}
          >
            {v || record.id}
          </Typography.Link>
        ),
      },
      {
        title: 'Protocols',
        dataIndex: 'protocols',
        width: 130,
        render: (v: string[]) => (
          <Flex gap={4} wrap="wrap">
            {v?.map((p) => (
              <Tag key={p}>{p}</Tag>
            ))}
          </Flex>
        ),
      },
      {
        title: 'Methods',
        dataIndex: 'methods',
        width: 200,
        render: (v: string[]) => (
          <Flex gap={4} wrap="wrap">
            {v?.map((m) => (
              <Tag key={m} color="blue">
                {m}
              </Tag>
            ))}
          </Flex>
        ),
      },
      {
        title: 'Paths',
        dataIndex: 'paths',
        ellipsis: true,
        render: (v: string[]) => v?.join(', ') || '-',
      },
      {
        title: '',
        width: 72,
        render: (_: any, record: any) => (
          <Flex gap={4}>
            <Button type="text" size="small" icon={<Edit size={14} />} onClick={() => handleEditRoute(record)} />
            <Button
              type="text"
              size="small"
              danger
              icon={<TrashCan size={14} />}
              onClick={() => handleDeleteRoute(record)}
            />
          </Flex>
        ),
      },
    ],
    [onViewRoute, handleEditRoute, handleDeleteRoute]
  );

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: 'Protocol',
      dataIndex: 'protocol',
      width: 100,
      render: (v: string) => <Tag>{v}</Tag>,
    },
    {
      title: 'Host',
      dataIndex: 'host',
      width: 200,
      ellipsis: true,
    },
    {
      title: 'Port',
      dataIndex: 'port',
      width: 80,
    },
    {
      title: 'Path',
      dataIndex: 'path',
      width: 150,
      ellipsis: true,
      render: (v: string) => v || '-',
    },
    {
      title: 'Enabled',
      dataIndex: 'enabled',
      width: 90,
      render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? 'Yes' : 'No'}</Tag>,
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
        key: 'detail',
        label: 'Details',
        children: detailRecord && (
          <Descriptions column={1} bordered size="small" className="detail-descriptions">
            <Descriptions.Item label="ID">{detailRecord.id}</Descriptions.Item>
            <Descriptions.Item label="Name">{detailRecord.name || '-'}</Descriptions.Item>
            <Descriptions.Item label="Protocol">{detailRecord.protocol}</Descriptions.Item>
            <Descriptions.Item label="Host">{detailRecord.host}</Descriptions.Item>
            <Descriptions.Item label="Port">{detailRecord.port}</Descriptions.Item>
            <Descriptions.Item label="Path">{detailRecord.path || '-'}</Descriptions.Item>
            <Descriptions.Item label="Retries">{detailRecord.retries}</Descriptions.Item>
            <Descriptions.Item label="Connect Timeout">{detailRecord.connect_timeout}ms</Descriptions.Item>
            <Descriptions.Item label="Write Timeout">{detailRecord.write_timeout}ms</Descriptions.Item>
            <Descriptions.Item label="Read Timeout">{detailRecord.read_timeout}ms</Descriptions.Item>
            <Descriptions.Item label="Enabled">
              <Tag color={detailRecord.enabled ? 'green' : 'red'}>{detailRecord.enabled ? 'Yes' : 'No'}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Created">
              {new Date(detailRecord.created_at * 1000).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="Updated">
              {new Date(detailRecord.updated_at * 1000).toLocaleString()}
            </Descriptions.Item>
          </Descriptions>
        ),
      },
      {
        key: 'routes',
        label: `Routes${serviceRoutes.length ? ` (${serviceRoutes.length})` : ''}`,
        children: (
          <>
            <Flex justify="flex-end" style={{ marginBottom: 12 }}>
              <Button
                type="primary"
                size="small"
                icon={<Add size={14} />}
                onClick={() => {
                  addRouteForm.resetFields();
                  setEditingRoute(null);
                  setAddRouteOpen(true);
                }}
              >
                Add Route
              </Button>
            </Flex>
            <Table
              rowKey="id"
              size="small"
              dataSource={serviceRoutes}
              columns={routeColumns}
              pagination={false}
              scroll={{ x: 'max-content' }}
              locale={{ emptyText: 'No routes associated with this service' }}
            />
          </>
        ),
      },
    ],
    [detailRecord, serviceRoutes, routeColumns, addRouteForm]
  );

  return (
    <>
      <div className="toolbar">
        <div className="toolbar-left">
          <Button type="primary" icon={<Add size={16} />} onClick={() => open()}>
            Add Service
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            Refresh
          </Button>
        </div>
        <div className="toolbar-right">
          <Input.Search
            placeholder="Search by name / host"
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
          pagination={{ pageSize: 20, showTotal: (t) => `Total ${t}`, showQuickJumper: true }}
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
                onClick: () => open(record),
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
      <Drawer title="Service Detail" open={!!detailRecord} onClose={() => setDetailRecord(null)} width={860}>
        <Tabs items={drawerTabs} />
      </Drawer>
      <Modal
        title={editingRoute ? 'Edit Route' : 'Add Route'}
        open={addRouteOpen}
        onCancel={() => {
          setAddRouteOpen(false);
          setEditingRoute(null);
          addRouteForm.resetFields();
        }}
        onOk={handleAddRoute}
        confirmLoading={savingRoute}
        width={560}
      >
        <Form form={addRouteForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="Name">
            <Input placeholder="my-route (optional)" />
          </Form.Item>
          <Form.Item name="protocols" label="Protocols" initialValue={['http', 'https']}>
            <Select mode="multiple" options={ROUTE_PROTOCOLS.map((p) => ({ label: p, value: p }))} />
          </Form.Item>
          <Form.Item name="methods" label="Methods">
            <Select mode="multiple" options={ROUTE_METHODS.map((m) => ({ label: m, value: m }))} />
          </Form.Item>
          <Form.Item name="paths" label="Paths" rules={[{ required: true, message: 'At least one path is required' }]}>
            <Input placeholder="/api/v1, /api/v2 (comma separated)" />
          </Form.Item>
          <Form.Item name="hosts" label="Hosts">
            <Input placeholder="example.com (comma separated, optional)" />
          </Form.Item>
          <Flex gap={12}>
            <Form.Item name="strip_path" label="Strip Path" initialValue={true} style={{ flex: 1 }}>
              <Select
                options={[
                  { label: 'Yes', value: true },
                  { label: 'No', value: false },
                ]}
              />
            </Form.Item>
            <Form.Item name="preserve_host" label="Preserve Host" initialValue={false} style={{ flex: 1 }}>
              <Select
                options={[
                  { label: 'Yes', value: true },
                  { label: 'No', value: false },
                ]}
              />
            </Form.Item>
          </Flex>
        </Form>
      </Modal>
    </>
  );
};

export default ServicesTab;
