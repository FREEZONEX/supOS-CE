import { type FC, useCallback, useMemo, useState } from 'react';
import { Button, App, Tag, Input, Descriptions, Drawer, Form, Typography } from 'antd';
import { Add, Renew, Edit, TrashCan } from '@carbon/icons-react';
import ProTable from '@/components/pro-table';
import ProSearch from '@/components/pro-search';
import { getCertificates, createCertificate, updateCertificate, deleteCertificate } from '@/apis/core-api/kong';
import useKongTable from '../../hooks/useKongTable';
import useKongModal from '../../hooks/useKongModal';
import useTranslate from '@/hooks/useTranslate';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';

const CertificatesTab: FC = () => {
  const { modal } = App.useApp();
  const formatMessage = useTranslate();
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
          label={formatMessage('kong.labelCertificatePEM')}
          rules={[{ required: true, message: formatMessage('kong.ruleCertRequired') }]}
        >
          <Input.TextArea
            rows={6}
            placeholder="-----BEGIN CERTIFICATE-----&#10;...&#10;-----END CERTIFICATE-----"
          />
        </Form.Item>
        <Form.Item
          name="key"
          label={formatMessage('kong.labelPrivateKeyPEM')}
          rules={[{ required: true, message: formatMessage('kong.ruleKeyRequired') }]}
        >
          <Input.TextArea
            rows={6}
            placeholder="-----BEGIN PRIVATE KEY-----&#10;...&#10;-----END PRIVATE KEY-----"
          />
        </Form.Item>
        <Form.Item name="snis" label={formatMessage('kong.labelSNIs')}>
          <Input placeholder={formatMessage('kong.phSNIs')} />
        </Form.Item>
        <Form.Item name="tags" label={formatMessage('kong.labelTags')}>
          <Input placeholder={formatMessage('kong.phTags')} />
        </Form.Item>
      </>
    ),
    [formatMessage]
  );

  const { ModalDom, open } = useKongModal({
    title: formatMessage('kong.certificates'),
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
      modal.confirm(
        mergeDeleteConfirmProps(
          {
            title: formatMessage('kong.deleteCertificate', { id: record.id.slice(0, 8) }),
            onOk: async () => {
              await deleteCertificate(record.id);
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
      title: formatMessage('kong.colID'),
      dataIndex: 'id',
      width: 280,
      ellipsis: true,
      render: (v: string, record: any) => (
        <Typography.Link onClick={() => setDetailRecord(record)}>{v}</Typography.Link>
      ),
    },
    {
      title: formatMessage('kong.colSNIs'),
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
      title: formatMessage('kong.colTags'),
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
            {formatMessage('kong.addCertificate')}
          </Button>
          <Button icon={<Renew size={16} />} onClick={refresh}>
            {formatMessage('common.refresh')}
          </Button>
        </div>
        <div className="toolbar-right">
          <ProSearch
            size="sm"
            placeholder={formatMessage('kong.searchCertificate')}
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
        title={formatMessage('kong.certificateDetail')}
        open={!!detailRecord}
        onClose={() => setDetailRecord(null)}
        width={640}
      >
        {detailRecord && (
          <>
            <Descriptions column={1} bordered size="small" className="detail-descriptions">
              <Descriptions.Item label={formatMessage('kong.colID')}>{detailRecord.id}</Descriptions.Item>
              <Descriptions.Item label={formatMessage('kong.colSNIs')}>
                {detailRecord.snis?.map((s: any) => (typeof s === 'string' ? s : s.name))?.join(', ') || '-'}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('kong.colTags')}>
                {detailRecord.tags?.join(', ') || '-'}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('kong.colCreated')}>
                {new Date(detailRecord.created_at * 1000).toLocaleString()}
              </Descriptions.Item>
            </Descriptions>
            <h4 style={{ marginTop: 16, marginBottom: 8 }}>{formatMessage('kong.colCertificate')}</h4>
            <div className="json-preview">{detailRecord.cert}</div>
            <h4 style={{ marginTop: 16, marginBottom: 8 }}>{formatMessage('kong.colPrivateKey')}</h4>
            <div className="json-preview">{detailRecord.key}</div>
          </>
        )}
      </Drawer>
    </>
  );
};

export default CertificatesTab;
