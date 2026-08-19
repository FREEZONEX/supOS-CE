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
import ProSearch from '@/components/pro-search';
import {
  getServices,
  createService,
  updateService,
  deleteService,
  getServiceRoutes,
  createRoute,
  updateRoute,
  deleteRoute,
} from '@/apis/core-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';
import useTranslate from '@/hooks/useTranslate';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';

const SVC_PROTOCOLS = ['http', 'https', 'grpc', 'grpcs', 'tcp', 'tls', 'udp'];
const ROUTE_PROTOCOLS = ['http', 'https', 'grpc', 'grpcs', 'tcp', 'tls', 'udp'];
const ROUTE_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'TRACE'];

interface ServicesTabProps {
  onViewRoute?: (route: any) => void;
}

const ServicesTab: FC<ServicesTabProps> = ({ onViewRoute }) => {
  const { modal, message } = App.useApp();
  const formatMessage = useTranslate();
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
        message.success(formatMessage('kong.routeUpdated'));
      } else {
        payload.service = { id: detailRecord.id };
        await createRoute(payload);
        message.success(formatMessage('kong.routeCreated'));
      }
      setAddRouteOpen(false);
      setEditingRoute(null);
      addRouteForm.resetFields();
      refreshRoutes(detailRecord.id);
    } catch (e: any) {
      if (!e?.errorFields) message.error(e?.message || formatMessage('kong.routeSaveFailed'));
    } finally {
      setSavingRoute(false);
    }
  }, [detailRecord, addRouteForm, refreshRoutes, message, editingRoute, formatMessage]);

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
      modal.confirm(
        mergeDeleteConfirmProps(
          {
            title: formatMessage('kong.deleteRoute', { name: route.name ?? route.id }),
            onOk: async () => {
              await deleteRoute(route.id);
              if (detailRecord?.id) refreshRoutes(detailRecord.id);
            },
          },
          formatMessage
        )
      );
    },
    [modal, detailRecord, refreshRoutes, formatMessage]
  );

  const renderForm = useCallback(
    (_form: any, editing: any) => (
      <>
        <Form.Item
          name="name"
          label={formatMessage('kong.labelName')}
          rules={[{ required: true, message: formatMessage('kong.ruleServiceNameRequired') }]}
        >
          <Input placeholder={formatMessage('kong.phServiceName')} />
        </Form.Item>
        <Flex gap={12}>
          <Form.Item
            name="protocol"
            label={formatMessage('kong.labelProtocol')}
            initialValue="http"
            style={{ flex: 1 }}
          >
            <Select options={SVC_PROTOCOLS.map((p) => ({ label: p, value: p }))} />
          </Form.Item>
          <Form.Item
            name="host"
            label={formatMessage('kong.labelHost')}
            rules={[{ required: true, message: formatMessage('kong.ruleHostRequired') }]}
            style={{ flex: 2 }}
          >
            <Input placeholder={formatMessage('kong.phHost')} />
          </Form.Item>
          <Form.Item name="port" label={formatMessage('kong.labelPort')} initialValue={80} style={{ flex: 1 }}>
            <InputNumber min={1} max={65535} style={{ width: '100%' }} />
          </Form.Item>
        </Flex>
        <Form.Item name="path" label={formatMessage('kong.labelPath')}>
          <Input placeholder={formatMessage('kong.phPath')} />
        </Form.Item>
        <Flex gap={12}>
          <Form.Item name="retries" label={formatMessage('kong.labelRetries')} initialValue={5} style={{ flex: 1 }}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="connect_timeout"
            label={formatMessage('kong.labelConnectTimeout')}
            initialValue={60000}
            style={{ flex: 1 }}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="write_timeout"
            label={formatMessage('kong.labelWriteTimeout')}
            initialValue={60000}
            style={{ flex: 1 }}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="read_timeout"
            label={formatMessage('kong.labelReadTimeout')}
            initialValue={60000}
            style={{ flex: 1 }}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Flex>
        {editing?.id && (
          <Form.Item name="enabled" label={formatMessage('kong.labelEnabled')} initialValue={true}>
            <Select
              options={[
                { label: formatMessage('kong.valEnabled'), value: true },
                { label: formatMessage('kong.valDisabled'), value: false },
              ]}
            />
          </Form.Item>
        )}
      </>
    ),
    [formatMessage]
  );

  const { ModalDom, open } = useKongModal({
    title: formatMessage('kong.services'),
    createApi: createService,
    updateApi: updateService,
    onSuccess: refresh,
    renderForm,
    width: 720,
  });

  const handleDelete = useCallback(
    (record: any) => {
      modal.confirm(
        mergeDeleteConfirmProps(
          {
            title: formatMessage('kong.deleteService', { name: record.name ?? record.id }),
            content: formatMessage('kong.deleteServiceContent'),
            onOk: async () => {
              await deleteService(record.id);
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
      (s: any) =>
        s.name?.toLowerCase().includes(q) || s.host?.toLowerCase().includes(q) || s.id?.toLowerCase().includes(q)
    );
  }, [data, search]);

  const routeColumns = useMemo(
    () => [
      {
        title: formatMessage('kong.colName'),
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
        title: formatMessage('kong.colProtocols'),
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
        title: formatMessage('kong.colMethods'),
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
        title: formatMessage('kong.colPaths'),
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
    [onViewRoute, handleEditRoute, handleDeleteRoute, formatMessage]
  );

  const columns = [
    {
      title: formatMessage('kong.colName'),
      dataIndex: 'name',
      width: 200,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: formatMessage('kong.colProtocol'),
      dataIndex: 'protocol',
      width: 100,
      render: (v: string) => <Tag>{v}</Tag>,
    },
    {
      title: formatMessage('kong.colHost'),
      dataIndex: 'host',
      width: 200,
      ellipsis: true,
    },
    {
      title: formatMessage('kong.colPort'),
      dataIndex: 'port',
      width: 80,
    },
    {
      title: formatMessage('kong.colPath'),
      dataIndex: 'path',
      width: 150,
      ellipsis: true,
      render: (v: string) => v || '-',
    },
    {
      title: formatMessage('kong.colEnabled'),
      dataIndex: 'enabled',
      width: 90,
      render: (v: boolean) => (
        <Tag color={v ? 'green' : 'red'}>{v ? formatMessage('kong.valYes') : formatMessage('kong.valNo')}</Tag>
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
        key: 'detail',
        label: formatMessage('kong.details'),
        children: detailRecord && (
          <Descriptions column={1} bordered size="small" className="detail-descriptions">
            <Descriptions.Item label={formatMessage('kong.colID')}>{detailRecord.id}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colName')}>{detailRecord.name || '-'}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colProtocol')}>{detailRecord.protocol}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colHost')}>{detailRecord.host}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colPort')}>{detailRecord.port}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colPath')}>{detailRecord.path || '-'}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.labelRetries')}>{detailRecord.retries}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.labelConnectTimeout')}>
              {detailRecord.connect_timeout}ms
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.labelWriteTimeout')}>
              {detailRecord.write_timeout}ms
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.labelReadTimeout')}>
              {detailRecord.read_timeout}ms
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colEnabled')}>
              <Tag color={detailRecord.enabled ? 'green' : 'red'}>
                {detailRecord.enabled ? formatMessage('kong.valYes') : formatMessage('kong.valNo')}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colCreated')}>
              {new Date(detailRecord.created_at * 1000).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colUpdated')}>
              {new Date(detailRecord.updated_at * 1000).toLocaleString()}
            </Descriptions.Item>
          </Descriptions>
        ),
      },
      {
        key: 'routes',
        label: `${formatMessage('kong.routes')}${serviceRoutes.length ? ` (${serviceRoutes.length})` : ''}`,
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
                {formatMessage('kong.addRoute')}
              </Button>
            </Flex>
            <Table
              rowKey="id"
              size="small"
              dataSource={serviceRoutes}
              columns={routeColumns}
              pagination={false}
              scroll={{ x: 'max-content' }}
              locale={{ emptyText: formatMessage('kong.noRoutes') }}
            />
          </>
        ),
      },
    ],
    [detailRecord, serviceRoutes, routeColumns, addRouteForm, formatMessage]
  );

  return (
    <>
      <div className="toolbar">
        <div className="toolbar-left">
          <Button type="primary" icon={<Add size={16} />} onClick={() => open()}>
            {formatMessage('kong.addService')}
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            {formatMessage('common.refresh')}
          </Button>
        </div>
        <div className="toolbar-right">
          <ProSearch
            size="sm"
            placeholder={formatMessage('kong.searchService')}
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
          pagination={{ pageSize: 20, showTotal: (t) => `Total ${t}`, showQuickJumper: true }}
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
                onClick: () => open(record),
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
        title={formatMessage('kong.serviceDetail')}
        open={!!detailRecord}
        onClose={() => setDetailRecord(null)}
        width={860}
      >
        <Tabs items={drawerTabs} />
      </Drawer>
      <Modal
        title={editingRoute ? formatMessage('kong.editRoute') : formatMessage('kong.addRoute')}
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
          <Form.Item name="name" label={formatMessage('kong.labelName')}>
            <Input placeholder={formatMessage('kong.phRouteNameOptional')} />
          </Form.Item>
          <Form.Item name="protocols" label={formatMessage('kong.labelProtocols')} initialValue={['http', 'https']}>
            <Select mode="multiple" options={ROUTE_PROTOCOLS.map((p) => ({ label: p, value: p }))} />
          </Form.Item>
          <Form.Item name="methods" label={formatMessage('kong.labelMethods')}>
            <Select mode="multiple" options={ROUTE_METHODS.map((m) => ({ label: m, value: m }))} />
          </Form.Item>
          <Form.Item
            name="paths"
            label={formatMessage('kong.labelPaths')}
            rules={[{ required: true, message: formatMessage('kong.rulePathRequired') }]}
          >
            <Input placeholder={formatMessage('kong.phPaths')} />
          </Form.Item>
          <Form.Item name="hosts" label={formatMessage('kong.labelHosts')}>
            <Input placeholder={formatMessage('kong.phHostsOptional')} />
          </Form.Item>
          <Flex gap={12}>
            <Form.Item
              name="strip_path"
              label={formatMessage('kong.labelStripPath')}
              initialValue={true}
              style={{ flex: 1 }}
            >
              <Select
                options={[
                  { label: formatMessage('kong.valYes'), value: true },
                  { label: formatMessage('kong.valNo'), value: false },
                ]}
              />
            </Form.Item>
            <Form.Item
              name="preserve_host"
              label={formatMessage('kong.labelPreserveHost')}
              initialValue={false}
              style={{ flex: 1 }}
            >
              <Select
                options={[
                  { label: formatMessage('kong.valYes'), value: true },
                  { label: formatMessage('kong.valNo'), value: false },
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
