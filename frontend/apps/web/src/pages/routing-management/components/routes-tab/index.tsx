import { type FC, useCallback, useMemo, useState, useEffect } from 'react';
import { Button, App, Tag, Input, Descriptions, Drawer, Form, Select, Typography } from 'antd';
import { Add, Renew, Edit, TrashCan } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import { getRoutes, createRoute, updateRoute, deleteRoute, getServices } from '@/apis/inter-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';

const PROTOCOLS = ['http', 'https', 'grpc', 'grpcs', 'tcp', 'tls', 'udp'];
const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'TRACE'];

interface RoutesTabProps {
  initialDetail?: any;
}

const RoutesTab: FC<RoutesTabProps> = ({ initialDetail }) => {
  const { modal } = App.useApp();
  const { data, loading, refresh } = useKongTable({ fetchApi: getRoutes });
  const [search, setSearch] = useState('');
  const [detailRecord, setDetailRecord] = useState<any>(initialDetail ?? null);
  const [serviceOptions, setServiceOptions] = useState<{ label: string; value: string }[]>([]);

  useEffect(() => {
    getServices({ size: 1000 }).then((res: any) => {
      setServiceOptions(
        (res?.data ?? []).map((s: any) => ({
          label: s.name || s.id,
          value: s.id,
        }))
      );
    });
  }, []);

  const transformValues = useCallback((values: any) => {
    const payload: Record<string, unknown> = { ...values };
    if (values.service_id) {
      payload.service = { id: values.service_id };
      delete payload.service_id;
    }
    if (typeof values.paths === 'string') {
      payload.paths = values.paths
        .split(',')
        .map((p: string) => p.trim())
        .filter(Boolean);
    }
    if (typeof values.hosts === 'string') {
      payload.hosts = values.hosts
        .split(',')
        .map((h: string) => h.trim())
        .filter(Boolean);
      if (payload.hosts && (payload.hosts as string[]).length === 0) delete payload.hosts;
    }
    return payload;
  }, []);

  const renderForm = useCallback(
    () => (
      <>
        <Form.Item name="name" label="Name">
          <Input placeholder="my-route" />
        </Form.Item>
        <Form.Item name="service_id" label="Service" rules={[{ required: true, message: 'Service is required' }]}>
          <Select showSearch placeholder="Select a service" options={serviceOptions} optionFilterProp="label" />
        </Form.Item>
        <Form.Item name="protocols" label="Protocols" initialValue={['http', 'https']}>
          <Select mode="multiple" options={PROTOCOLS.map((p) => ({ label: p, value: p }))} />
        </Form.Item>
        <Form.Item name="methods" label="Methods" initialValue={['GET']}>
          <Select mode="multiple" options={METHODS.map((m) => ({ label: m, value: m }))} />
        </Form.Item>
        <Form.Item name="paths" label="Paths" rules={[{ required: true, message: 'At least one path is required' }]}>
          <Input placeholder="/api/v1, /api/v2 (comma separated)" />
        </Form.Item>
        <Form.Item name="hosts" label="Hosts">
          <Input placeholder="example.com, api.example.com (comma separated)" />
        </Form.Item>
        <Form.Item name="strip_path" label="Strip Path" initialValue={true}>
          <Select
            options={[
              { label: 'Yes', value: true },
              { label: 'No', value: false },
            ]}
          />
        </Form.Item>
        <Form.Item name="preserve_host" label="Preserve Host" initialValue={false}>
          <Select
            options={[
              { label: 'Yes', value: true },
              { label: 'No', value: false },
            ]}
          />
        </Form.Item>
      </>
    ),
    [serviceOptions]
  );

  const { ModalDom, open } = useKongModal({
    title: 'Route',
    createApi: createRoute,
    updateApi: updateRoute,
    onSuccess: refresh,
    renderForm,
    transformValues,
  });

  const handleEdit = useCallback(
    (record: any) => {
      const formValues: any = { ...record };
      formValues.service_id = record.service?.id;
      if (Array.isArray(record.paths)) {
        formValues.paths = record.paths.join(', ');
      }
      if (Array.isArray(record.hosts)) {
        formValues.hosts = record.hosts.join(', ');
      }
      open(formValues);
    },
    [open]
  );

  const handleDelete = useCallback(
    (record: any) => {
      modal.confirm({
        title: `Delete route "${record.name ?? record.id}"?`,
        okButtonProps: { danger: true },
        onOk: async () => {
          await deleteRoute(record.id);
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
      (r: any) =>
        r.name?.toLowerCase().includes(q) ||
        r.paths?.some((p: string) => p.toLowerCase().includes(q)) ||
        r.id?.toLowerCase().includes(q)
    );
  }, [data, search]);

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 180,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: 'Protocols',
      dataIndex: 'protocols',
      width: 140,
      render: (v: string[]) => (
        <div className="tag-list">
          {v?.map((p) => (
            <Tag key={p}>{p}</Tag>
          ))}
        </div>
      ),
    },
    {
      title: 'Methods',
      dataIndex: 'methods',
      width: 200,
      render: (v: string[]) => (
        <div className="tag-list">
          {v?.map((m) => (
            <Tag key={m} color="blue">
              {m}
            </Tag>
          ))}
        </div>
      ),
    },
    {
      title: 'Paths',
      dataIndex: 'paths',
      width: 220,
      ellipsis: true,
      render: (v: string[]) => v?.join(', ') || '-',
    },
    {
      title: 'Service',
      dataIndex: 'service',
      width: 180,
      ellipsis: true,
      render: (v: any) => {
        const name = serviceOptions.find((s) => s.value === v?.id)?.label;
        return name || v?.id?.slice(0, 8) || '-';
      },
    },
    {
      title: 'Strip Path',
      dataIndex: 'strip_path',
      width: 100,
      render: (v: boolean) => (v ? 'Yes' : 'No'),
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
            Add Route
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            Refresh
          </Button>
        </div>
        <div className="toolbar-right">
          <Input.Search
            placeholder="Search by name / path"
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
      <Drawer title="Route Detail" open={!!detailRecord} onClose={() => setDetailRecord(null)} width={560}>
        {detailRecord && (
          <Descriptions column={1} bordered size="small" className="detail-descriptions">
            <Descriptions.Item label="ID">{detailRecord.id}</Descriptions.Item>
            <Descriptions.Item label="Name">{detailRecord.name || '-'}</Descriptions.Item>
            <Descriptions.Item label="Protocols">{detailRecord.protocols?.join(', ')}</Descriptions.Item>
            <Descriptions.Item label="Methods">{detailRecord.methods?.join(', ')}</Descriptions.Item>
            <Descriptions.Item label="Paths">{detailRecord.paths?.join(', ') || '-'}</Descriptions.Item>
            <Descriptions.Item label="Hosts">{detailRecord.hosts?.join(', ') || '-'}</Descriptions.Item>
            <Descriptions.Item label="Service">{detailRecord.service?.id || '-'}</Descriptions.Item>
            <Descriptions.Item label="Strip Path">{String(detailRecord.strip_path)}</Descriptions.Item>
            <Descriptions.Item label="Preserve Host">{String(detailRecord.preserve_host)}</Descriptions.Item>
            <Descriptions.Item label="Regex Priority">{detailRecord.regex_priority}</Descriptions.Item>
            <Descriptions.Item label="Created">
              {new Date(detailRecord.created_at * 1000).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label="Updated">
              {new Date(detailRecord.updated_at * 1000).toLocaleString()}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>
    </>
  );
};

export default RoutesTab;
