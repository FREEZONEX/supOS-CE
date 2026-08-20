import { TreeSelect } from 'antd';
import { useCallback, useEffect, useState } from 'react';
import { getUnsNode, isTopicNode, listUnsNodes, normalizeNodeList, type UnsNode } from '../api/uns';
import { t } from '../i18n';

interface TreeNodeData {
  id: string;
  pId: string | null;
  value: string;
  title: string;
  isLeaf: boolean;
  selectable: boolean;
  topicPath: string;
}

const toTreeNode = (node: UnsNode, parentId: string | null): TreeNodeData => {
  const topic = isTopicNode(node);
  // UNS 数据面约定（dev 起）：flow 发布到 namespace 全路径（旧版为 alias）
  const path = String(node.namespace || node.path || node.alias || node.name || '');
  return {
    id: String(node.id),
    pId: parentId,
    value: String(node.id),
    title: String(node.name || node.namespace || path),
    isLeaf: topic,
    selectable: topic,
    topicPath: path,
  };
};

// UNS 树懒加载选择器：只允许选择文件型节点（可订阅 topic）
export function TopicSelect({
  value,
  onChange,
}: {
  value?: { unsNodeId: string; topic: string };
  onChange: (next: { unsNodeId: string; topic: string }) => void;
}) {
  const [treeData, setTreeData] = useState<TreeNodeData[]>([]);
  // 树是懒加载的，已选节点大概率不在 treeData 里——labelInValue + 拉节点详情解析真实路径回显
  const [selectedLabel, setSelectedLabel] = useState<string>('');

  const loadChildren = useCallback(async (parentId: string | null) => {
    const data = await listUnsNodes(
      parentId === null ? { parentIdSet: true, parentId: 0 } : { parentIdSet: true, parentId }
    );
    return normalizeNodeList(data).map((node) => toTreeNode(node, parentId));
  }, []);

  useEffect(() => {
    loadChildren(null)
      .then(setTreeData)
      .catch(() => setTreeData([]));
  }, [loadChildren]);

  useEffect(() => {
    if (!value?.unsNodeId) {
      setSelectedLabel('');
      return;
    }
    let cancelled = false;
    getUnsNode(value.unsNodeId)
      .then((node) => {
        if (cancelled) return;
        setSelectedLabel(String(node.namespace || node.path || node.name || value.topic || value.unsNodeId));
      })
      .catch(() => {
        if (!cancelled) setSelectedLabel(value.topic || value.unsNodeId);
      });
    return () => {
      cancelled = true;
    };
  }, [value?.unsNodeId, value?.topic]);

  return (
    <TreeSelect
      style={{ width: '100%' }}
      treeDataSimpleMode
      treeData={treeData}
      labelInValue
      value={value?.unsNodeId ? { value: value.unsNodeId, label: selectedLabel || value.topic } : undefined}
      placeholder={t('instance.topicPlaceholder')}
      loadData={async (node) => {
        const children = await loadChildren(String(node.id));
        setTreeData((prev) => {
          const existing = new Set(prev.map((item) => item.id));
          return [...prev, ...children.filter((item) => !existing.has(item.id))];
        });
      }}
      onSelect={(_selected, node) => {
        const item = node as unknown as TreeNodeData;
        if (!item.selectable) return;
        onChange({ unsNodeId: item.id, topic: item.topicPath });
      }}
    />
  );
}
