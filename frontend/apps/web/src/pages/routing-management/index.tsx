import { type FC, useCallback, useEffect, useMemo, useState } from 'react';
import { Add, Renew, TrashCan, Edit, Connect } from '@carbon/icons-react';
import {
  App,
  Button,
  Descriptions,
  Drawer,
  Form,
  Input,
  InputNumber,
  Segmented,
  Select,
  Space,
  Switch,
  Table,
  Tag,
  Typography,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { PageProps } from '@/common-types';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton } from '@/components/auth';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ProSearch from '@/components/pro-search';
import { useTranslate } from '@/hooks';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import {
  deleteGatewayRoute,
  listGatewayRoutes,
  saveGatewayRoute,
  testGatewayRoute,
  updateGatewayRoute,
  type GatewayRoute,
} from '@/apis/core-api/gateway';
import styles from './index.module.scss';

const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS'];

const targetTypeLabelKeys: Record<string, string> = {
  reverseProxy: 'route.target.reverseProxy',
  flowPayload: 'route.target.flowPayload',
  static: 'route.target.static',
};

const resourceOptionKeys = [
  ['', 'route.auth.publicAccess'],
  ['login', 'route.auth.loginUser'],
  ['uns.read', 'resource.uns.read'],
  ['uns.write', 'resource.uns.write'],
  ['uns.delete', 'resource.uns.delete'],
  ['uns.import', 'resource.uns.import'],
  ['uns.export', 'resource.uns.export'],
  ['uns.template.manage', 'resource.uns.template.manage'],
  ['uns.manage', 'resource.uns.manage'],
  ['flow.read', 'resource.flow.read'],
  ['flow.manage', 'resource.flow.manage'],
  ['notebook.view', 'Notebook.title'],
  ['gateway.route.manage', 'route.routingManagement'],
  ['gateway.route.create', 'resource.gateway.route.create'],
  ['gateway.route.update', 'resource.gateway.route.update'],
  ['gateway.route.delete', 'resource.gateway.route.delete'],
  ['gateway.route.test', 'resource.gateway.route.test'],
  ['cluster.edge.manage', 'mqttAuth.edgeNode'],
];

type RouteScope = 'all' | 'system' | 'custom';

const normalizeValues = (values: GatewayRoute): GatewayRoute => ({
  ...values,
  routeKey: values.routeKey?.trim(),
  name: values.name?.trim(),
  description: values.description?.trim(),
  targetType: 'reverseProxy',
  matchType: 'prefix',
  pathPattern: values.pathPattern?.trim(),
  targetUrl: values.targetUrl?.trim(),
  targetPath: values.targetPath?.trim() || '/',
  rewritePath: values.rewritePath?.trim(),
  resourceKey: values.authPolicy === 'public' ? '' : values.resourceKey || '',
  timeoutMs: Number(values.timeoutMs || 10000),
  priority: Number(values.priority || 100),
  enabled: values.enabled !== false,
});

const RoutingManagement: FC<PageProps> = ({ title }) => {
  const formatMessage = useTranslate();
  const { message, modal } = App.useApp();
  const [form] = Form.useForm<GatewayRoute>();
  const [routes, setRoutes] = useState<GatewayRoute[]>([]);
  const [loading, setLoading] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<GatewayRoute | null>(null);
  const [detail, setDetail] = useState<GatewayRoute | null>(null);
  const [keyword, setKeyword] = useState('');
  const [routeScope, setRouteScope] = useState<RouteScope>('all');

  const resourceOptions = useMemo(
    () => resourceOptionKeys.map(([value, label]) => ({ value, label: formatMessage(label) })),
    [formatMessage]
  );

  const loadRoutes = useCallback(async () => {
    setLoading(true);
    try {
      const resp = await listGatewayRoutes();
      setRoutes(resp.list);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadRoutes();
  }, [loadRoutes]);

  const openCreate = () => {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({
      targetType: 'reverseProxy',
      matchType: 'prefix',
      methods: ['GET'],
      targetPath: '/',
      stripPrefix: true,
      authPolicy: 'login',
      timeoutMs: 10000,
      priority: 100,
      enabled: true,
    } as GatewayRoute);
    setDrawerOpen(true);
  };

  const openEdit = (record: GatewayRoute) => {
    if (record.systemBuiltin) {
      setDetail(record);
      return;
    }
    setEditing(record);
    form.setFieldsValue({
      ...record,
      methods: record.methods?.length ? record.methods : ['GET'],
      enabled: record.enabled !== false,
    });
    setDrawerOpen(true);
  };

  const submit = async () => {
    if (editing?.systemBuiltin) {
      message.warning(formatMessage('route.systemReadonly'));
      return;
    }
    const values = normalizeValues(await form.validateFields());
    if (editing) {
      await updateGatewayRoute(editing.routeKey, values);
      message.success(formatMessage('route.updated'));
    } else {
      await saveGatewayRoute(values);
      message.success(formatMessage('route.created'));
    }
    setDrawerOpen(false);
    await loadRoutes();
  };

  const remove = (record: GatewayRoute) => {
    if (record.systemBuiltin) {
      message.warning(formatMessage('route.systemReadonly'));
      return;
    }
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          title: formatMessage('route.deleteTitle', { name: record.name || record.routeKey }),
          content: formatMessage('route.deleteContent'),
          onOk: async () => {
            await deleteGatewayRoute(record.routeKey);
            message.success(formatMessage('route.deleted'));
            await loadRoutes();
          },
        },
        formatMessage
      )
    );
  };

  const testRoute = async (record: GatewayRoute) => {
    const resp = await testGatewayRoute(record.routeKey);
    message.success(resp?.reachable ? formatMessage('route.testPassed') : formatMessage('route.testReturned'));
  };

  const filteredRoutes = useMemo(() => {
    const q = keyword.trim().toLowerCase();
    if (!q) return routes;
    return routes.filter((item) =>
      [item.routeKey, item.name, item.pathPattern, item.targetUrl, item.resourceKey]
        .filter(Boolean)
        .some((value) => String(value).toLowerCase().includes(q))
    );
  }, [keyword, routes]);

  const systemRoutes = useMemo(() => filteredRoutes.filter((item) => item.systemBuiltin), [filteredRoutes]);
  const customRoutes = useMemo(() => filteredRoutes.filter((item) => !item.systemBuiltin), [filteredRoutes]);
  const visibleRoutes = useMemo(() => {
    if (routeScope === 'system') return systemRoutes;
    if (routeScope === 'custom') return customRoutes;
    return filteredRoutes;
  }, [customRoutes, filteredRoutes, routeScope, systemRoutes]);

  const columns = (): ColumnsType<GatewayRoute> => [
    {
      title: formatMessage('common.name'),
      dataIndex: 'name',
      width: 220,
      render: (value, record) => (
        <Space size={6}>
          <Typography.Link onClick={() => setDetail(record)}>{value || record.routeKey}</Typography.Link>
          {record.systemBuiltin && <Tag>{formatMessage('OpenData.builtIn')}</Tag>}
        </Space>
      ),
    },
    {
      title: formatMessage('common.path'),
      dataIndex: 'pathPattern',
      width: 220,
      ellipsis: true,
    },
    {
      title: formatMessage('OpenData.type'),
      dataIndex: 'targetType',
      width: 120,
      render: (value: string) => (
        <Tag>{formatMessage(targetTypeLabelKeys[value] || value, undefined, value || '-')}</Tag>
      ),
    },
    {
      title: formatMessage('route.methods'),
      dataIndex: 'methods',
      width: 180,
      render: (methods?: string[]) => (
        <Space size={[4, 4]} wrap>
          {(methods?.length ? methods : ['*']).map((method) => (
            <Tag key={method} color={method === 'GET' ? 'blue' : 'default'}>
              {method}
            </Tag>
          ))}
        </Space>
      ),
    },
    {
      title: formatMessage('route.upstream'),
      dataIndex: 'targetUrl',
      width: 260,
      ellipsis: true,
      render: (value, record) => {
        if (record.targetType === 'flowPayload') return formatMessage('route.generatedByBackend');
        if (record.targetType === 'static') return formatMessage('route.localStaticAssets');
        return value || '-';
      },
    },
    {
      title: formatMessage('route.auth'),
      dataIndex: 'authPolicy',
      width: 130,
      render: (value, record) => {
        if (value === 'public') return <Tag color="green">{formatMessage('route.public')}</Tag>;
        if (record.resourceKey) return <Tag color="orange">{record.resourceKey}</Tag>;
        return <Tag color="blue">{formatMessage('route.login')}</Tag>;
      },
    },
    {
      title: formatMessage('route.priority'),
      dataIndex: 'priority',
      width: 90,
    },
    {
      title: formatMessage('common.status'),
      dataIndex: 'enabled',
      width: 90,
      render: (enabled) =>
        enabled ? (
          <Tag color="green">{formatMessage('common.enable')}</Tag>
        ) : (
          <Tag>{formatMessage('common.disable')}</Tag>
        ),
    },
    {
      title: formatMessage('common.operation'),
      fixed: 'right',
      width: 240,
      render: (_, record) => {
        const editable = !record.systemBuiltin;
        return (
          <Space>
            <Button size="small" onClick={() => setDetail(record)}>
              {formatMessage('common.view')}
            </Button>
            <AuthButton
              auth={ButtonPermission['RoutingManagement.test']}
              size="small"
              icon={<Connect size={14} />}
              onClick={() => testRoute(record)}
            >
              {formatMessage('route.check')}
            </AuthButton>
            {editable ? (
              <>
                <AuthButton
                  auth={ButtonPermission['RoutingManagement.edit']}
                  size="small"
                  icon={<Edit size={14} />}
                  onClick={() => openEdit(record)}
                />
                <AuthButton
                  auth={ButtonPermission['RoutingManagement.delete']}
                  size="small"
                  danger
                  icon={<TrashCan size={14} />}
                  onClick={() => remove(record)}
                />
              </>
            ) : null}
          </Space>
        );
      },
    },
  ];

  return (
    <ComLayout>
      <ComContent hasBack={false} title={title || formatMessage('route.routingManagement')}>
        <div className={styles['routing-page']}>
          <div className="toolbar">
            <Space>
              <AuthButton
                auth={ButtonPermission['RoutingManagement.add']}
                type="primary"
                icon={<Add size={16} />}
                onClick={openCreate}
              >
                {formatMessage('route.newRoute')}
              </AuthButton>
              <Button icon={<Renew size={16} />} onClick={loadRoutes}>
                {formatMessage('common.refresh')}
              </Button>
            </Space>
            <Space>
              <Segmented
                value={routeScope}
                onChange={(value) => setRouteScope(value as RouteScope)}
                options={[
                  { label: `${formatMessage('common.all')} (${filteredRoutes.length})`, value: 'all' },
                  { label: `${formatMessage('route.systemRoutes')} (${systemRoutes.length})`, value: 'system' },
                  { label: `${formatMessage('route.customRoutes')} (${customRoutes.length})`, value: 'custom' },
                ]}
              />
              <ProSearch
                size="sm"
                placeholder={formatMessage('route.searchPlaceholder')}
                style={{ width: 320 }}
                value={keyword}
                onChange={(event) => setKeyword(event.target.value)}
                onSearch={setKeyword}
              />
            </Space>
          </div>
          <div className="table-area">
            <div className="section-header">
              <div>
                <div className="section-title">{formatMessage('route.routeList')}</div>
                <div className="section-meta">{formatMessage('route.routeListHint')}</div>
              </div>
              <Tag>{formatMessage('route.count', { count: visibleRoutes.length })}</Tag>
            </div>
            <Table
              rowKey="routeKey"
              size="small"
              loading={loading}
              dataSource={visibleRoutes}
              columns={columns()}
              pagination={{ pageSize: 10, showSizeChanger: true }}
              scroll={{ x: 1300 }}
            />
          </div>
        </div>

        <Drawer
          title={
            editing
              ? formatMessage('route.editRouteTitle', { name: editing.name || editing.routeKey })
              : formatMessage('route.newGatewayRoute')
          }
          open={drawerOpen}
          width={620}
          onClose={() => setDrawerOpen(false)}
          extra={
            <Space>
              <Button onClick={() => setDrawerOpen(false)}>{formatMessage('common.cancel')}</Button>
              <AuthButton
                auth={editing ? ButtonPermission['RoutingManagement.edit'] : ButtonPermission['RoutingManagement.add']}
                type="primary"
                onClick={submit}
              >
                {formatMessage('common.save')}
              </AuthButton>
            </Space>
          }
          destroyOnClose
        >
          <Form form={form} layout="vertical" preserve={false}>
            <Form.Item
              name="routeKey"
              label={formatMessage('route.routeKey')}
              rules={[{ required: true, message: formatMessage('route.routeKeyRequired') }]}
              tooltip={formatMessage('route.routeKeyTooltip')}
            >
              <Input disabled={!!editing} placeholder="proxy.custom-app" />
            </Form.Item>
            <Form.Item
              name="name"
              label={formatMessage('common.name')}
              rules={[{ required: true, message: formatMessage('route.nameRequired') }]}
            >
              <Input placeholder={formatMessage('route.namePlaceholder')} />
            </Form.Item>
            <Form.Item name="description" label={formatMessage('common.description')}>
              <Input.TextArea rows={2} placeholder={formatMessage('route.optional')} />
            </Form.Item>
            <Form.Item
              name="pathPattern"
              label={formatMessage('route.entryPath')}
              rules={[{ required: true, message: formatMessage('route.entryPathRequired') }]}
            >
              <Input placeholder="/custom/**" />
            </Form.Item>
            <Form.Item
              name="methods"
              label={formatMessage('route.allowedMethods')}
              rules={[{ required: true, message: formatMessage('route.methodsRequired') }]}
            >
              <Select mode="multiple" options={METHODS.map((method) => ({ label: method, value: method }))} />
            </Form.Item>
            <Form.Item
              name="targetUrl"
              label={formatMessage('route.upstreamUrl')}
              rules={[{ required: true, message: formatMessage('route.upstreamRequired') }]}
            >
              <Input placeholder="http://service:8080" />
            </Form.Item>
            <Form.Item name="targetPath" label={formatMessage('route.upstreamPath')}>
              <Input placeholder="/" />
            </Form.Item>
            <Form.Item name="stripPrefix" label={formatMessage('route.stripPrefix')} valuePropName="checked">
              <Switch />
            </Form.Item>
            <Form.Item name="rewritePath" label={formatMessage('route.rewritePath')}>
              <Input placeholder={formatMessage('route.rewritePathPlaceholder')} />
            </Form.Item>
            <Form.Item
              name="authPolicy"
              label={formatMessage('route.authPolicy')}
              rules={[{ required: true, message: formatMessage('route.authPolicyRequired') }]}
            >
              <Select
                options={[
                  { label: formatMessage('route.public'), value: 'public' },
                  { label: formatMessage('route.login'), value: 'login' },
                ]}
              />
            </Form.Item>
            <Form.Item noStyle shouldUpdate={(prev, next) => prev.authPolicy !== next.authPolicy}>
              {({ getFieldValue }) =>
                getFieldValue('authPolicy') === 'login' ? (
                  <Form.Item name="resourceKey" label={formatMessage('route.resourceKey')}>
                    <Select showSearch options={resourceOptions.filter((item) => item.value !== 'login')} />
                  </Form.Item>
                ) : null
              }
            </Form.Item>
            <Form.Item
              name="priority"
              label={formatMessage('route.priority')}
              rules={[{ required: true, message: formatMessage('route.priorityRequired') }]}
            >
              <InputNumber min={1} max={9999} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="timeoutMs" label={formatMessage('route.timeoutMs')}>
              <InputNumber min={1000} max={600000} step={1000} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="enabled" label={formatMessage('common.enable')} valuePropName="checked">
              <Switch />
            </Form.Item>
          </Form>
        </Drawer>

        <Drawer title={formatMessage('route.detail')} open={!!detail} width={620} onClose={() => setDetail(null)}>
          {detail && (
            <Descriptions column={1} bordered size="small">
              <Descriptions.Item label="Key">{detail.routeKey}</Descriptions.Item>
              <Descriptions.Item label={formatMessage('common.name')}>{detail.name}</Descriptions.Item>
              <Descriptions.Item label={formatMessage('common.path')}>{detail.pathPattern}</Descriptions.Item>
              <Descriptions.Item label={formatMessage('OpenData.type')}>
                {formatMessage(
                  targetTypeLabelKeys[detail.targetType] || detail.targetType,
                  undefined,
                  detail.targetType || '-'
                )}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('route.methods')}>
                {detail.methods?.join(', ') || '*'}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('route.upstream')}>
                {detail.targetType === 'flowPayload'
                  ? formatMessage('route.generatedByBackend')
                  : detail.targetType === 'static'
                    ? formatMessage('route.localStaticAssets')
                    : detail.targetUrl || '-'}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('route.stripPrefix')}>
                {detail.stripPrefix ? formatMessage('common.yes') : formatMessage('common.no')}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('route.auth')}>{detail.authPolicy}</Descriptions.Item>
              <Descriptions.Item label={formatMessage('route.resourceKey')}>
                {detail.resourceKey || '-'}
              </Descriptions.Item>
              <Descriptions.Item label={formatMessage('route.priority')}>{detail.priority}</Descriptions.Item>
              <Descriptions.Item label={formatMessage('common.description')}>
                {detail.description || '-'}
              </Descriptions.Item>
            </Descriptions>
          )}
        </Drawer>
      </ComContent>
    </ComLayout>
  );
};

export default RoutingManagement;
