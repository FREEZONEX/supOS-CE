import { type FC, useCallback, useMemo, useState } from 'react';
import { Button, App, Tag, Input, Descriptions, Drawer, Form, Typography } from 'antd';
import { Add, Renew, Edit, TrashCan } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import { getCertificates, createCertificate, updateCertificate, deleteCertificate } from '@/apis/inter-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';

const CertificatesTab: FC = () => {
  const { modal } = App.useApp();
  const { data, loading, refresh } = useKongTable({ fetchApi: getCertificates });
  const [search, setSearch] = useState('');
  const [detailRecord, setDetailRecord] = useState<any>(null);

  const transformValues = useCallback((values: any) => {
    const payload: Record<string, unknown> = { ...values };
    if (typeof values.snis === 'string') {
      payload.snis = values.snis
        .split(',')
        .map((s: string) => s.trim())
        .filter(Boolean);
      if ((payload.snis as string[]).length === 0) delete payload.snis;
    }
    if (typeof values.tags === 'string') {
      payload.tags = values.tags
        .split(',')
        .map((t: string) => t.trim())
        .filter(Boolean);
      if ((payload.tags as string[]).length === 0) delete payload.tags;
    }
    return payload;
  }, []);

  const renderForm = useCallback(
    () => (
      <>
        <Form.Item
          name="cert"
          label="Certificate (PEM)"
          rules={[{ required: true, message: 'Certificate is required' }]}
        >
          <Input.TextArea
            rows={6}
            placeholder="-----BEGIN CERTIFICATE-----&#10;...&#10;-----END CERTIFICATE-----"
          />
        </Form.Item>
        <Form.Item
          name="key"
          label="Private Key (PEM)"
          rules={[{ required: true, message: 'Private key is required' }]}
        >
          <Input.TextArea
            rows={6}
            placeholder="-----BEGIN PRIVATE KEY-----&#10;...&#10;-----END PRIVATE KEY-----"
          />
        </Form.Item>
        <Form.Item name="snis" label="SNIs">
          <Input placeholder="example.com, *.example.com (comma separated)" />
        </Form.Item>
        <Form.Item name="tags" label="Tags">
          <Input placeholder="tag1, tag2 (comma separated)" />
        </Form.Item>
      </>
    ),
    []
  );

  const { ModalDom, open } = useKongModal({
    title: 'Certificate',
    createApi: createCertificate,
    updateApi: updateCertificate,
    onSuccess: refresh,
    renderForm,
    transformValues,
    width: 680,
  });

  const handleEdit = useCallback(
    (record: any) => {
      const formValues = {
        ...record,
        snis: record.snis?.map?.((s: any) => (typeof s === 'string' ? s : s.name))?.join(', ') ?? '',
        tags: record.tags?.join(', ') ?? '',
      };
      open(formValues);
    },
    [open]
  );

  const handleDelete = useCallback(
    (record: any) => {
      modal.confirm({
        title: `Delete certificate ${record.id.slice(0, 8)}...?`,
        okButtonProps: { danger: true },
        onOk: async () => {
          await deleteCertificate(record.id);
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
        c.id?.toLowerCase().includes(q) ||
        c.snis?.some?.((s: any) => {
          const name = typeof s === 'string' ? s : s.name;
          return name?.toLowerCase().includes(q);
        }) ||
        c.tags?.some?.((t: string) => t.toLowerCase().includes(q))
    );
  }, [data, search]);

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 280,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v}</Typography.Link>
      ),
    },
    {
      title: 'SNIs',
      dataIndex: 'snis',
      width: 260,
      render: (v: any[]) => {
        if (!v?.length) return '-';
        return (
          <div className="tag-list">
            {v.map((s: any) => {
              const name = typeof s === 'string' ? s : s.name;
              return <Tag key={name}>{name}</Tag>;
            })}
          </div>
        );
      },
    },
    {
      title: 'Tags',
      dataIndex: 'tags',
      width: 180,
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

  return (
    <>
      <div className="toolbar">
        <div className="toolbar-left">
          <Button type="primary" icon={<Add size={16} />} onClick={() => open()}>
            Add Certificate
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            Refresh
          </Button>
        </div>
        <div className="toolbar-right">
          <Input.Search
            placeholder="Search by SNI / tag"
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
      <Drawer title="Certificate Detail" open={!!detailRecord} onClose={() => setDetailRecord(null)} width={640}>
        {detailRecord && (
          <>
            <Descriptions column={1} bordered size="small" className="detail-descriptions">
              <Descriptions.Item label="ID">{detailRecord.id}</Descriptions.Item>
              <Descriptions.Item label="SNIs">
                {detailRecord.snis?.map((s: any) => (typeof s === 'string' ? s : s.name))?.join(', ') || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="Tags">{detailRecord.tags?.join(', ') || '-'}</Descriptions.Item>
              <Descriptions.Item label="Created">
                {new Date(detailRecord.created_at * 1000).toLocaleString()}
              </Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16, marginBottom: 8 }}>Certificate</h4>
            <div className="json-preview">{detailRecord.cert}</div>
            <h4 style={{ marginTop: 16, marginBottom: 8 }}>Private Key</h4>
            <div className="json-preview">{detailRecord.key}</div>
          </>
        )}
      </Drawer>
    </>
  );
};

export default CertificatesTab;
