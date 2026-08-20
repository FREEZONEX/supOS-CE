import { useBaseStore } from '@/stores/base';
import { CloseOutlined } from '@ant-design/icons';
import { ConfigProvider, Form, type FormInstance, Tag, TreeSelect } from 'antd';
import type { CustomTagProps } from 'rc-select/lib/BaseSelect';
import ProTreeSelect from '@/components/pro-tree-select/ProTreeSelect.tsx';
import { getFlowAndGroupList } from '@/apis/core-api/flow.ts';
import { getEventFlowAndGroupList } from '@/apis/core-api/event-flow.ts';
import useTranslate from '@/hooks/useTranslate.ts';
import FlowItemIcon from '@/pages/flow/components/FlowItemIcon';
import { type MutableRefObject, useCallback, useRef } from 'react';
import styles from './index.module.scss';
const { SHOW_PARENT } = TreeSelect;

type FlowNode = {
  id: number | string;
  name?: string;
  title?: string;
  category?: string;
  flowKind?: 'source' | 'event';
  hasChildren?: boolean;
};

type FlowNodeRegistry = Map<string, FlowNode>;

const normalizeFlowNodeKey = (value: unknown) => String(value ?? '');

const registerFlowNodes = (registry: FlowNodeRegistry, items: FlowNode[] | undefined, flowKind: 'source' | 'event') => {
  items?.forEach((item) => {
    if (item?.id == null) {
      return;
    }
    registry.set(normalizeFlowNodeKey(item.id), {
      ...item,
      flowKind: item.flowKind ?? flowKind,
    });
  });
};

const mapFlowTreeNode = (item: FlowNode, flowKind: 'source' | 'event', registry: FlowNodeRegistry) => {
  const node = {
    ...item,
    isLeaf: !item.hasChildren,
    value: item.id,
    title: item.name,
    key: item.id,
    name: item.name,
    flowKind,
  };
  registry.set(normalizeFlowNodeKey(item.id), node);
  return node;
};

const resolveFlowTagValue = (value: CustomTagProps['value']) => {
  if (value && typeof value === 'object' && 'value' in value) {
    return (value as { value?: number | string }).value;
  }
  return value;
};

const resolveFlowTagLabel = (
  label: CustomTagProps['label'],
  value: CustomTagProps['value'],
  registry: FlowNodeRegistry
) => {
  const rawValue = resolveFlowTagValue(value);
  const node = registry.get(normalizeFlowNodeKey(rawValue));
  if (node?.name) {
    return node.name;
  }
  if (node?.title) {
    return String(node.title);
  }
  if (typeof label === 'string' || typeof label === 'number') {
    return String(label);
  }
  if (label && typeof label === 'object' && 'name' in label) {
    return String((label as { name?: string }).name ?? rawValue ?? '');
  }
  return String(rawValue ?? '');
};

const EXPORT_FLOW_TAG_MAX_CHARS = 4;

const truncateFlowTagLabel = (text: string, maxChars = EXPORT_FLOW_TAG_MAX_CHARS) => {
  const normalized = text.trim();
  if (!normalized) {
    return '';
  }
  const chars = Array.from(normalized);
  if (chars.length <= maxChars) {
    return normalized;
  }
  return `${chars.slice(0, maxChars).join('')}...`;
};

const createRenderExportFlowTag =
  (registry: MutableRefObject<FlowNodeRegistry>, flowKind: 'source' | 'event') =>
  (props: CustomTagProps) => {
    const { label, value, closable, onClose } = props;
    const rawValue = resolveFlowTagValue(value);
    const node = registry.current.get(normalizeFlowNodeKey(rawValue));
    const fullName = resolveFlowTagLabel(label, value, registry.current);

    return (
      <Tag
        className={styles.exportFlowTag}
        closable={closable}
        closeIcon={<CloseOutlined className={styles.exportFlowTagCloseIcon} />}
        onClose={onClose}
        title={fullName}
      >
        <span className={styles.exportFlowTagContent}>
          <FlowItemIcon
            category={node?.category}
            flowKind={node?.flowKind ?? flowKind}
            size="sm"
          />
          <span className={styles.exportFlowTagText} title={fullName}>
            {truncateFlowTagLabel(fullName)}
          </span>
        </span>
      </Tag>
    );
  };

const renderFlowTreeIcon = (dataNode: unknown, flowKind: 'source' | 'event') => {
  const node = (dataNode || {}) as { category?: string; flowKind?: string };
  return (
    <FlowItemIcon
      category={node.category}
      flowKind={node.flowKind ?? flowKind}
      size="sm"
    />
  );
};

const OtherDom = ({ form }: { form: FormInstance }) => {
  const { containerList } = useBaseStore((state) => ({
    containerList: state.containerList,
  }));
  const hasNodeRed = !!containerList?.aboutUs?.some((s) => s.name === 'nodered') || true;
  const hasEventflow = !!containerList?.aboutUs?.some((s) => s.name === 'eventflow') || true;
  const formatMessage = useTranslate();
  const sourceFlowRegistryRef = useRef<FlowNodeRegistry>(new Map());
  const eventFlowRegistryRef = useRef<FlowNodeRegistry>(new Map());

  const buildFlowTreeSelectProps = useCallback(
    (flowKind: 'source' | 'event', registry: MutableRefObject<FlowNodeRegistry>) => ({
      className: styles.exportFlowTreeSelect,
      popupClassName: 'export-flow-tree-select-popup',
      size: 'middle' as const,
      selectAllAble: true,
      labelInValue: true,
      showSearch: false,
      loadDataEnable: true,
      lazy: true,
      listHeight: 350,
      showSwitcherIcon: true,
      treeCheckable: true,
      fieldNames: { label: 'name', value: 'id' },
      showCheckedStrategy: SHOW_PARENT,
      allowClear: true,
      tagRender: createRenderExportFlowTag(registry, flowKind),
    }),
    []
  );

  return (
    <div style={{ width: '100%' }}>
      <ConfigProvider
        theme={{
          components: {
            Form: {
              itemMarginBottom: 12,
            },
          },
        }}
      >
        <Form
          layout="vertical"
          name="exportForm"
          form={form}
          colon={false}
          style={{ color: 'var(--ui-text-color)' }}
          initialValues={{
            sourceFlowExportParam: [],
            eventFlowExportParam: [],
          }}
          // disabled={loading}
        >
          {hasNodeRed && (
            <Form.Item label={formatMessage('home.sourceFlow')} name="sourceFlowExportParam">
              <ProTreeSelect
                {...buildFlowTreeSelectProps('source', sourceFlowRegistryRef)}
                api={(params, config) =>
                  getFlowAndGroupList(
                    {
                      k: params?.searchValue,
                      groupId: params?.key ? params?.key : undefined,
                      pageNo: params?.pageNo,
                      pageSize: params?.pageSize,
                    },
                    config
                  ).then((data) => {
                    const nodes = data?.data?.map((item: FlowNode) =>
                      mapFlowTreeNode(item, 'source', sourceFlowRegistryRef.current)
                    );
                    registerFlowNodes(sourceFlowRegistryRef.current, nodes, 'source');
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: nodes,
                    };
                  })
                }
                treeNodeIcon={(dataNode) => renderFlowTreeIcon(dataNode, 'source')}
              />
            </Form.Item>
          )}
          {hasEventflow && (
            <Form.Item label={formatMessage('home.eventFlow')} name="eventFlowExportParam">
              <ProTreeSelect
                {...buildFlowTreeSelectProps('event', eventFlowRegistryRef)}
                api={(params, config) =>
                  getEventFlowAndGroupList(
                    {
                      k: params?.searchValue,
                      groupId: params?.key ? params?.key : undefined,
                      pageNo: params?.pageNo,
                      pageSize: params?.pageSize,
                    },
                    config
                  ).then((data) => {
                    const nodes = data?.data?.map((item: FlowNode) =>
                      mapFlowTreeNode(item, 'event', eventFlowRegistryRef.current)
                    );
                    registerFlowNodes(eventFlowRegistryRef.current, nodes, 'event');
                    return {
                      pageNo: data?.pageNo,
                      pageSize: data?.pageSize,
                      total: data?.total,
                      data: nodes,
                    };
                  })
                }
                treeNodeIcon={(dataNode) => renderFlowTreeIcon(dataNode, 'event')}
              />
            </Form.Item>
          )}
        </Form>
      </ConfigProvider>
    </div>
  );
};

export default OtherDom;
