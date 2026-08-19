import { type FC, useEffect, useMemo, useState } from 'react';
import { Add, Edit, ListDropdown, TrashCan } from '@carbon/icons-react';
import ComEmpty from '@/components/com-empty';
import { App, Button, Flex, Form, Input, Modal, Select, Space, Spin, Switch, Table, Tag, Tree } from 'antd';
import type { DataNode } from 'antd/es/tree';
import type { ColumnsType } from 'antd/es/table';
import type { PageProps } from '@/common-types';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComLeft from '@/components/com-layout/ComLeft.tsx';
import ProSearch from '@/components/pro-search';
import {
  deleteResourceActionApi,
  getResourceActionsApi,
  postResourceActionApi,
  putResourceActionApi,
} from '@/apis/core-api/resource.ts';
import { iamApi, resourceNameMessageKey, type CoreResource } from '@/apis/core-api/core-adapter.ts';
import useTranslate from '@/hooks/useTranslate.ts';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import styles from './index.module.scss';

type FlatResource = CoreResource & {
  key: string;
};

type ResourceAction = {
  id: string | number;
  resourceId?: string | number;
  actionType: string;
  actionValue: string;
  methods?: string;
  enabled?: number;
  systemGenerated?: number;
};

const actionTypeOptions = ['ui', 'api', 'openapi', 'gateway'].map((value) => ({ label: value, value }));

const methodOptions = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'].map((value) => ({
  label: value,
  value,
}));

const flattenResources = (items: CoreResource[] = []) => {
  const out: FlatResource[] = [];
  const walk = (nodes: CoreResource[]) => {
    nodes.forEach((node) => {
      out.push({ ...node, key: String(node.id) });
      if (node.children?.length) {
        walk(node.children);
      }
    });
  };
  walk(items);
  return out;
};

const normalizeMethods = (methods?: string | string[]) => {
  const list = Array.isArray(methods) ? methods : String(methods || '').split(',');
  return list
    .map((item) => String(item).trim().toUpperCase())
    .filter(Boolean)
    .join(',');
};

const methodList = (methods?: string) =>
  String(methods || '')
    .split(',')
    .map((item) => item.trim().toUpperCase())
    .filter(Boolean);

const canMaintainActions = (type?: number) => type === 2 || type === 3 || type === 4 || type === 5;

const nodeTypeLabel = (type: number | undefined, formatMessage: (key: string) => string) => {
  if (type === 1) return formatMessage('PermissionManagement.nodeType.group');
  if (type === 2) return formatMessage('PermissionManagement.nodeType.page');
  if (type === 3) return formatMessage('PermissionManagement.nodeType.permission');
  if (type === 4) return formatMessage('PermissionManagement.nodeType.homeTab');
  if (type === 5) return formatMessage('PermissionManagement.nodeType.subMenu');
  return String(type || '-');
};

const nodeTypeColor = (type?: number) => {
  if (type === 1) return 'blue';
  if (type === 2 || type === 4 || type === 5) return 'cyan';
  if (type === 3) return 'green';
  return 'default';
};

const resourceDisplayName = (
  item: CoreResource,
  formatMessage: (key: string, opt?: any, defaultMessage?: string) => string
) => {
  const messageKey = resourceNameMessageKey(item.resourceKey);
  if (messageKey) {
    return formatMessage(messageKey, undefined, item.name || item.resourceKey);
  }
  return item.name || item.resourceKey;
};

const filterResources = (items: CoreResource[] = [], keyword: string): CoreResource[] => {
  const normalized = keyword.trim().toLowerCase();
  if (!normalized) return items;

  const walk = (nodes: CoreResource[]): CoreResource[] =>
    nodes.reduce<CoreResource[]>((result, node) => {
      const children = walk(node.children || []);
      const hit = [node.name, node.resourceKey, node.routePath]
        .filter(Boolean)
        .some((value) => String(value).toLowerCase().includes(normalized));
      if (hit || children.length > 0) {
        result.push({ ...node, children });
      }
      return result;
    }, []);

  return walk(items);
};

const toTreeData = (
  items: CoreResource[] = [],
  formatMessage: (key: string, opt?: any, defaultMessage?: string) => string
): DataNode[] =>
  items.map((item) => ({
    key: String(item.id),
    selectable: canMaintainActions(item.type),
    title: (
      <div
        className={`${styles.treeNode} ${canMaintainActions(item.type) ? '' : styles.treeNodeReadonly}`}
        title={`${resourceDisplayName(item, formatMessage)} (${item.resourceKey})`}
      >
        <span className={styles.treeNodeName}>{resourceDisplayName(item, formatMessage)}</span>
      </div>
    ),
    children: item.children?.length ? toTreeData(item.children, formatMessage) : undefined,
  }));

const PermissionManagement: FC<PageProps> = ({ title }) => {
  const formatMessage = useTranslate();
  const { modal, message } = App.useApp();
  const [form] = Form.useForm();
  const [resources, setResources] = useState<CoreResource[]>([]);
  const [selectedKey, setSelectedKey] = useState<string>();
  const [actions, setActions] = useState<ResourceAction[]>([]);
  const [resourceLoading, setResourceLoading] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [keyword, setKeyword] = useState('');
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<ResourceAction | null>(null);

  const flatResources = useMemo(() => flattenResources(resources), [resources]);
  const selected = flatResources.find((item) => item.key === selectedKey);
  const canEditActions = canMaintainActions(selected?.type);
  const filteredResources = useMemo(() => filterResources(resources, keyword), [keyword, resources]);

  const requestResources = async () => {
    setResourceLoading(true);
    try {
      const list = (await iamApi.get('/resources')) || [];
      setResources(list);
      const flat = flattenResources(list);
      const first = flat.find((item) => canMaintainActions(item.type));
      setSelectedKey((current) => {
        if (current && flat.some((item) => item.key === current && canMaintainActions(item.type))) {
          return current;
        }
        return first?.key;
      });
    } finally {
      setResourceLoading(false);
    }
  };

  const requestActions = async (resourceId?: string) => {
    if (!resourceId) {
      setActions([]);
      return;
    }
    setActionLoading(true);
    try {
      const data = await getResourceActionsApi(resourceId);
      setActions(data || []);
    } finally {
      setActionLoading(false);
    }
  };

  useEffect(() => {
    requestResources();
  }, []);

  useEffect(() => {
    requestActions(selected?.key);
  }, [selected?.key]);

  const openActionModal = (record?: ResourceAction) => {
    if (!canEditActions) {
      message.warning(formatMessage('PermissionManagement.actionReadonly'));
      return;
    }
    setEditing(record || null);
    form.setFieldsValue({
      actionType: record?.actionType || 'api',
      actionValue: record?.actionValue || '',
      methods: methodList(record?.methods),
      enabled: record ? record.enabled !== 0 : true,
    });
    setModalOpen(true);
  };

  const onSave = async () => {
    if (!selected?.key || !canEditActions) return;
    const values = await form.validateFields();
    const payload = {
      actionType: values.actionType,
      actionValue: String(values.actionValue || '').trim(),
      methods: normalizeMethods(values.methods),
      enabled: values.enabled ? 1 : 0,
    };
    setSaving(true);
    try {
      if (editing?.id) {
        await putResourceActionApi(selected.key, String(editing.id), payload);
      } else {
        await postResourceActionApi(selected.key, payload);
      }
      message.success(formatMessage('common.optsuccess'));
      setModalOpen(false);
      requestActions(selected.key);
    } finally {
      setSaving(false);
    }
  };

  const onDelete = (record: ResourceAction) => {
    if (!selected?.key || !canEditActions) return;
    modal.confirm({
      ...createDeleteConfirmOptions({
        title: formatMessage('common.deleteConfirm'),
        name: record?.actionValue,
        formatMessage,
      }),
      onOk: async () => {
        await deleteResourceActionApi(selected.key, String(record.id));
        message.success(formatMessage('common.deleteSuccessfully'));
        requestActions(selected.key);
      },
    });
  };

  const columns: ColumnsType<ResourceAction> = [
    {
      title: formatMessage('resourceAction.type'),
      dataIndex: 'actionType',
      width: 120,
      render: (value) => <Tag bordered={false}>{value}</Tag>,
    },
    {
      title: formatMessage('resourceAction.value'),
      dataIndex: 'actionValue',
      ellipsis: true,
    },
    {
      title: formatMessage('resourceAction.methods'),
      dataIndex: 'methods',
      width: 180,
      render: (value) => {
        const methods = methodList(value);
        if (methods.length === 0) return '*';
        return (
          <Space size={4} wrap>
            {methods.map((method) => (
              <Tag key={method} bordered={false}>
                {method}
              </Tag>
            ))}
          </Space>
        );
      },
    },
    {
      title: formatMessage('common.status'),
      dataIndex: 'enabled',
      width: 110,
      render: (value) => (
        <Tag bordered={false} color={value === 0 ? 'default' : 'green'}>
          {formatMessage(value === 0 ? 'resourceAction.disabled' : 'resourceAction.enabled')}
        </Tag>
      ),
    },
    {
      title: formatMessage('common.operation'),
      width: 120,
      fixed: 'right',
      render: (_, record) => (
        <Space size={4}>
          <Button
            type="text"
            size="small"
            disabled={!canEditActions || record.systemGenerated === 1}
            icon={<Edit size={14} />}
            title={formatMessage('common.edit')}
            onClick={() => openActionModal(record)}
          />
          <Button
            type="text"
            size="small"
            danger
            disabled={!canEditActions || record.systemGenerated === 1}
            icon={<TrashCan size={14} />}
            title={formatMessage('common.delete')}
            onClick={() => onDelete(record)}
          />
        </Space>
      ),
    },
  ];

  return (
    <ComLayout>
      <ComContent
        hasBack={false}
        title={
          <Flex align="center" gap={8} style={{ lineHeight: 1 }}>
            <ListDropdown size={20} style={{ justifyContent: 'center', verticalAlign: 'middle' }} />
            <span>{title}</span>
          </Flex>
        }
      >
        <ComLayout>
          <ComLeft
            resize
            defaultWidth={380}
            style={{ display: 'flex', flexDirection: 'column', padding: '16px 0 16px 16px' }}
          >
            <Flex justify="space-between" align="center" className={styles.leftHeader}>
              <div>
                <div className={styles.leftTitle}>{formatMessage('PermissionManagement.resourceTree')}</div>
                <div className={styles.leftSubTitle}>{formatMessage('PermissionManagement.resourceTreeHint')}</div>
              </div>
              <Button size="small" type="text" onClick={requestResources}>
                {formatMessage('common.refresh')}
              </Button>
            </Flex>
            <ProSearch
              size="sm"
              placeholder={formatMessage('PermissionManagement.searchPlaceholder')}
              value={keyword}
              onChange={(event) => setKeyword(event.target.value)}
              className={styles.search}
            />
            <Spin spinning={resourceLoading}>
              {filteredResources.length > 0 ? (
                <Tree
                  blockNode
                  showLine
                  selectedKeys={selectedKey ? [selectedKey] : []}
                  treeData={toTreeData(filteredResources, formatMessage)}
                  defaultExpandAll
                  onSelect={(keys) => setSelectedKey(String(keys[0] || selectedKey || ''))}
                />
              ) : (
                <ComEmpty />
              )}
            </Spin>
          </ComLeft>
          <ComContent>
            <div className={styles.detailPane}>
              {selected ? (
                <>
                  <Flex align="center" justify="space-between" gap={16} className={styles.detailHeader}>
                    <div className={styles.titleBlock}>
                      <div className={styles.eyebrow}>{formatMessage('PermissionManagement.selectedNode')}</div>
                      <Flex align="center" gap={8}>
                        <span className={styles.resourceName}>{resourceDisplayName(selected, formatMessage)}</span>
                        <Tag bordered={false} color={nodeTypeColor(selected.type)}>
                          {nodeTypeLabel(selected.type, formatMessage)}
                        </Tag>
                      </Flex>
                      <div className={styles.resourceKey}>{selected.resourceKey}</div>
                    </div>
                    <Button
                      type="primary"
                      icon={<Add size={16} />}
                      disabled={!canEditActions}
                      onClick={() => openActionModal()}
                    >
                      {formatMessage('resourceAction.add')}
                    </Button>
                  </Flex>

                  <div className={styles.section}>
                    <div className={styles.sectionTitle}>{formatMessage('PermissionManagement.nodeInfo')}</div>
                    <div className={styles.infoGrid}>
                      <div>
                        <span>{formatMessage('PermissionManagement.resourceKey')}</span>
                        <strong>{selected.resourceKey || '-'}</strong>
                      </div>
                      <div>
                        <span>{formatMessage('PermissionManagement.nodeName')}</span>
                        <strong>{resourceDisplayName(selected, formatMessage)}</strong>
                      </div>
                      <div>
                        <span>{formatMessage('PermissionManagement.nodeType')}</span>
                        <strong>{nodeTypeLabel(selected.type, formatMessage)}</strong>
                      </div>
                      <div>
                        <span>{formatMessage('PermissionManagement.routePath')}</span>
                        <strong>{selected.routePath || '-'}</strong>
                      </div>
                    </div>
                  </div>

                  <div className={styles.section}>
                    <Flex align="center" justify="space-between" className={styles.sectionHeader}>
                      <div>
                        <div className={styles.sectionTitle}>{formatMessage('PermissionManagement.entryBindings')}</div>
                        <div className={styles.sectionHint}>
                          {formatMessage('PermissionManagement.entryBindingsHint')}
                        </div>
                      </div>
                      <Tag bordered={false}>{actions.length}</Tag>
                    </Flex>
                    <Table
                      rowKey={(record) =>
                        String(record.id || `${record.actionType}:${record.actionValue}:${record.methods}`)
                      }
                      size="small"
                      loading={actionLoading}
                      columns={columns}
                      dataSource={actions}
                      pagination={false}
                      scroll={{ x: 860 }}
                    />
                  </div>
                </>
              ) : (
                <ComEmpty description={formatMessage('PermissionManagement.selectNode')} />
              )}
            </div>
          </ComContent>
        </ComLayout>
      </ComContent>

      <Modal
        open={modalOpen}
        title={formatMessage(editing ? 'resourceAction.edit' : 'resourceAction.add')}
        onCancel={() => setModalOpen(false)}
        onOk={onSave}
        confirmLoading={saving}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="actionType" label={formatMessage('resourceAction.type')} rules={[{ required: true }]}>
            <Select options={actionTypeOptions} />
          </Form.Item>
          <Form.Item name="actionValue" label={formatMessage('resourceAction.value')} rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="methods" label={formatMessage('resourceAction.methods')}>
            <Select mode="tags" options={methodOptions} tokenSeparators={[',']} />
          </Form.Item>
          <Form.Item name="enabled" label={formatMessage('common.status')} valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </ComLayout>
  );
};

export default PermissionManagement;
