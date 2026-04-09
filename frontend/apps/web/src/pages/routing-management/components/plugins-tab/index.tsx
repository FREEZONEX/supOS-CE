import { type FC, useCallback, useMemo, useState, useEffect } from 'react';
import { Button, App, Tag, Input, Drawer, Form, Select, Switch, Descriptions, Typography } from 'antd';
import { Add, Renew, Edit, TrashCan } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import ProCodemirror from '@/components/pro-codemirror';
import { json } from '@codemirror/lang-json';
import {
  getPlugins,
  getEnabledPlugins,
  createPlugin,
  updatePlugin,
  deletePlugin,
  getServices,
  getRoutes,
  getConsumers,
} from '@/apis/inter-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';

const PluginsTab: FC = () => {
  const { modal } = App.useApp();
  const { data, loading, refresh } = useKongTable({ fetchApi: getPlugins });
  const [search, setSearch] = useState('');
  const [detailRecord, setDetailRecord] = useState<any>(null);

  const [availablePlugins, setAvailablePlugins] = useState<string[]>([]);
  const [serviceOpts, setServiceOpts] = useState<{ label: string; value: string }[]>([]);
  const [routeOpts, setRouteOpts] = useState<{ label: string; value: string }[]>([]);
  const [consumerOpts, setConsumerOpts] = useState<{ label: string; value: string }[]>([]);

  useEffect(() => {
    getEnabledPlugins()
      .then((res: any) => setAvailablePlugins(res?.enabled_plugins ?? []))
      .catch(() => {});
    getServices({ size: 1000 }).then((res: any) =>
      setServiceOpts((res?.data ?? []).map((s: any) => ({ label: s.name || s.id, value: s.id })))
    );
    getRoutes({ size: 1000 }).then((res: any) =>
      setRouteOpts((res?.data ?? []).map((r: any) => ({ label: r.name || r.id, value: r.id })))
    );
    getConsumers({ size: 1000 }).then((res: any) =>
      setConsumerOpts((res?.data ?? []).map((c: any) => ({ label: c.username || c.id, value: c.id })))
    );
  }, []);

  const transformValues = useCallback((values: any) => {
    const payload: Record<string, unknown> = {
      name: values.name,
      enabled: values.enabled ?? true,
    };
    if (values.service_id) payload.service = { id: values.service_id };
    if (values.route_id) payload.route = { id: values.route_id };
    if (values.consumer_id) payload.consumer = { id: values.consumer_id };
    if (values.config) {
      try {
        payload.config = JSON.parse(values.config);
      } catch {
        payload.config = {};
      }
    }
    return payload;
  }, []);

  const renderForm = useCallback(
    (_form: any, editing: any) => (
      <>
        <Form.Item name="name" label="Plugin" rules={[{ required: true, message: 'Plugin name is required' }]}>
          <Select
            showSearch
            placeholder="Select a plugin"
            disabled={!!editing?.id}
            options={availablePlugins.map((p) => ({ label: p, value: p }))}
            optionFilterProp="label"
          />
        </Form.Item>
        <Form.Item name="service_id" label="Service (optional)">
          <Select showSearch allowClear placeholder="Global if empty" options={serviceOpts} optionFilterProp="label" />
        </Form.Item>
        <Form.Item name="route_id" label="Route (optional)">
          <Select showSearch allowClear placeholder="Global if empty" options={routeOpts} optionFilterProp="label" />
        </Form.Item>
        <Form.Item name="consumer_id" label="Consumer (optional)">
          <Select showSearch allowClear placeholder="Global if empty" options={consumerOpts} optionFilterProp="label" />
        </Form.Item>
        <Form.Item name="enabled" label="Enabled" valuePropName="checked" initialValue={true}>
          <Switch />
        </Form.Item>
        <Form.Item name="config" label="Config (JSON)">
          <ProCodemirror
            extensions={[json()]}
            height="160px"
            basicSetup={{ lineNumbers: false, foldGutter: false }}
            showHint={false}
          />
        </Form.Item>
      </>
    ),
    [availablePlugins, serviceOpts, routeOpts, consumerOpts]
  );

  const { ModalDom, open } = useKongModal({
    title: 'Plugin',
    createApi: createPlugin,
    updateApi: updatePlugin,
    onSuccess: refresh,
    renderForm,
    transformValues,
    width: 680,
  });

  const handleEdit = useCallback(
    (record: any) => {
      const formValues: any = {
        ...record,
        service_id: record.service?.id,
        route_id: record.route?.id,
        consumer_id: record.consumer?.id,
        config: record.config ? JSON.stringify(record.config, null, 2) : '',
      };
      open(formValues);
    },
    [open]
  );

  const handleDelete = useCallback(
    (record: any) => {
      modal.confirm({
        title: `Delete plugin "${record.name}" (${record.id.slice(0, 8)}...)?`,
        okButtonProps: { danger: true },
        onOk: async () => {
          await deletePlugin(record.id);
          refresh();
        },
      });
    },
    [modal, refresh]
  );

  const filteredData = useMemo(() => {
    if (!search) return data;
    const q = search.toLowerCase();
    return data.filter((p: any) => p.name?.toLowerCase().includes(q) || p.id?.toLowerCase().includes(q));
  }, [data, search]);

  const columns = [
    {
      title: 'Plugin',
      dataIndex: 'name',
      width: 180,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: 'Scope',
      width: 160,
      render: (_: any, record: any) => {
        if (record.service?.id) return <Tag color="blue">Service</Tag>;
        if (record.route?.id) return <Tag color="cyan">Route</Tag>;
        if (record.consumer?.id) return <Tag color="purple">Consumer</Tag>;
        return <Tag>Global</Tag>;
      },
    },
    {
      title: 'Applied To',
      width: 220,
      ellipsis: true,
      render: (_: any, record: any) => {
        if (record.service?.id) {
          const name = serviceOpts.find((s) => s.value === record.service.id)?.label;
          return name || record.service.id.slice(0, 8);
        }
        if (record.route?.id) {
          const name = routeOpts.find((r) => r.value === record.route.id)?.label;
          return name || record.route.id.slice(0, 8);
        }
        if (record.consumer?.id) {
          const name = consumerOpts.find((c) => c.value === record.consumer.id)?.label;
          return name || record.consumer.id.slice(0, 8);
        }
        return '-';
      },
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

  return (
    <>
      <div className="toolbar">
        <div className="toolbar-left">
          <Button type="primary" icon={<Add size={16} />} onClick={() => open()}>
            Add Plugin
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            Refresh
          </Button>
        </div>
        <div className="toolbar-right">
          <Input.Search
            placeholder="Search by plugin name"
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
      <Drawer title="Plugin Detail" open={!!detailRecord} onClose={() => setDetailRecord(null)} width={640}>
        {detailRecord && (
          <>
            <Descriptions column={1} bordered size="small" className="detail-descriptions">
              <Descriptions.Item label="ID">{detailRecord.id}</Descriptions.Item>
              <Descriptions.Item label="Plugin">{detailRecord.name}</Descriptions.Item>
              <Descriptions.Item label="Enabled">
                <Tag color={detailRecord.enabled ? 'green' : 'red'}>{detailRecord.enabled ? 'Yes' : 'No'}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Service">{detailRecord.service?.id || '-'}</Descriptions.Item>
              <Descriptions.Item label="Route">{detailRecord.route?.id || '-'}</Descriptions.Item>
              <Descriptions.Item label="Consumer">{detailRecord.consumer?.id || '-'}</Descriptions.Item>
              <Descriptions.Item label="Created">
                {new Date(detailRecord.created_at * 1000).toLocaleString()}
              </Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16, marginBottom: 8 }}>Config</h4>
            <ProCodemirror
              value={JSON.stringify(detailRecord.config, null, 2)}
              extensions={[json()]}
              readOnly
              editable={false}
              basicSetup={{ lineNumbers: false, foldGutter: false }}
            />
          </>
        )}
      </Drawer>
    </>
  );
};

export default PluginsTab;
