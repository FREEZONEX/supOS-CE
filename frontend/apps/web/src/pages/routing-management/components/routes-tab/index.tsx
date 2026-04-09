import { type FC, useCallback, useMemo, useState, useEffect } from 'react';
import { Button, App, Tag, Input, Descriptions, Drawer, Form, Select, Typography } from 'antd';
import { Add, Renew, Edit, TrashCan } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import { getRoutes, createRoute, updateRoute, deleteRoute, getServices } from '@/apis/inter-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';
import useTranslate from '@/hooks/useTranslate';

const PROTOCOLS = ['http', 'https', 'grpc', 'grpcs', 'tcp', 'tls', 'udp'];
const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', 'TRACE'];

interface RoutesTabProps {
  initialDetail?: any;
}

const RoutesTab: FC<RoutesTabProps> = ({ initialDetail }) => {
  const { modal } = App.useApp();
  const formatMessage = useTranslate();
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
        <Form.Item name="name" label={formatMessage('kong.labelName')}>
          <Input placeholder={formatMessage('kong.phRouteName')} />
        </Form.Item>
        <Form.Item
          name="service_id"
          label={formatMessage('kong.colService')}
          rules={[{ required: true, message: formatMessage('kong.ruleServiceRequired') }]}
        >
          <Select
            showSearch
            placeholder={formatMessage('common.select')}
            options={serviceOptions}
            optionFilterProp="label"
          />
        </Form.Item>
        <Form.Item name="protocols" label={formatMessage('kong.labelProtocols')} initialValue={['http', 'https']}>
          <Select mode="multiple" options={PROTOCOLS.map((p) => ({ label: p, value: p }))} />
        </Form.Item>
        <Form.Item name="methods" label={formatMessage('kong.labelMethods')} initialValue={['GET']}>
          <Select mode="multiple" options={METHODS.map((m) => ({ label: m, value: m }))} />
        </Form.Item>
        <Form.Item
          name="paths"
          label={formatMessage('kong.labelPaths')}
          rules={[{ required: true, message: formatMessage('kong.rulePathRequired') }]}
        >
          <Input placeholder={formatMessage('kong.phPaths')} />
        </Form.Item>
        <Form.Item name="hosts" label={formatMessage('kong.labelHosts')}>
          <Input placeholder={formatMessage('kong.phHosts')} />
        </Form.Item>
        <Form.Item name="strip_path" label={formatMessage('kong.labelStripPath')} initialValue={true}>
          <Select
            options={[
              { label: formatMessage('kong.valYes'), value: true },
              { label: formatMessage('kong.valNo'), value: false },
            ]}
          />
        </Form.Item>
        <Form.Item name="preserve_host" label={formatMessage('kong.labelPreserveHost')} initialValue={false}>
          <Select
            options={[
              { label: formatMessage('kong.valYes'), value: true },
              { label: formatMessage('kong.valNo'), value: false },
            ]}
          />
        </Form.Item>
      </>
    ),
    [serviceOptions, formatMessage]
  );

  const { ModalDom, open } = useKongModal({
    title: formatMessage('kong.routes'),
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
        title: formatMessage('kong.deleteRoute', { name: record.name ?? record.id }),
        okButtonProps: { danger: true },
        onOk: async () => {
          await deleteRoute(record.id);
          refresh();
        },
      });
    },
    [modal, refresh, formatMessage]
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
      title: formatMessage('kong.colName'),
      dataIndex: 'name',
      width: 180,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v || record.id}</Typography.Link>
      ),
    },
    {
      title: formatMessage('kong.colProtocols'),
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
      title: formatMessage('kong.colMethods'),
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
      title: formatMessage('kong.colPaths'),
      dataIndex: 'paths',
      width: 220,
      ellipsis: true,
      render: (v: string[]) => v?.join(', ') || '-',
    },
    {
      title: formatMessage('kong.colService'),
      dataIndex: 'service',
      width: 180,
      ellipsis: true,
      render: (v: any) => {
        const name = serviceOptions.find((s) => s.value === v?.id)?.label;
        return name || v?.id?.slice(0, 8) || '-';
      },
    },
    {
      title: formatMessage('kong.colStripPath'),
      dataIndex: 'strip_path',
      width: 100,
      render: (v: boolean) => (v ? formatMessage('kong.valYes') : formatMessage('kong.valNo')),
    },
    {
      title: formatMessage('kong.colCreated'),
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
            {formatMessage('kong.addRoute')}
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            {formatMessage('common.refresh')}
          </Button>
        </div>
        <div className="toolbar-right">
          <Input.Search
            placeholder={formatMessage('kong.searchRoute')}
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
        title={formatMessage('kong.routeDetail')}
        open={!!detailRecord}
        onClose={() => setDetailRecord(null)}
        width={560}
      >
        {detailRecord && (
          <Descriptions column={1} bordered size="small" className="detail-descriptions">
            <Descriptions.Item label={formatMessage('kong.colID')}>{detailRecord.id}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colName')}>{detailRecord.name || '-'}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colProtocols')}>
              {detailRecord.protocols?.join(', ')}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colMethods')}>
              {detailRecord.methods?.join(', ')}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colPaths')}>
              {detailRecord.paths?.join(', ') || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.labelHosts')}>
              {detailRecord.hosts?.join(', ') || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colService')}>
              {detailRecord.service?.id || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colStripPath')}>
              {String(detailRecord.strip_path)}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.labelPreserveHost')}>
              {String(detailRecord.preserve_host)}
            </Descriptions.Item>
            <Descriptions.Item label="Regex Priority">{detailRecord.regex_priority}</Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colCreated')}>
              {new Date(detailRecord.created_at * 1000).toLocaleString()}
            </Descriptions.Item>
            <Descriptions.Item label={formatMessage('kong.colUpdated')}>
              {new Date(detailRecord.updated_at * 1000).toLocaleString()}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>
    </>
  );
};

export default RoutesTab;
