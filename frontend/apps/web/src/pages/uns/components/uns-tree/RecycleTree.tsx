import { clearRecycleTree, forceDeleteRecycleNode, getRecycleTreeData, restoreRecycleNode } from '@/apis/core-api/uns';
import ProTree from '@/components/pro-tree';
import ComEmpty from '@/components/com-empty';
import { ChevronDown, ChevronUp, TrashCan, Undo, WarningFilled } from '@/components/lucide-icon/carbon';
import { treeIconProps } from '@/components/lucide-icon/icon-props';
import { useTranslate } from '@/hooks';
import type { UnsTreeNode } from '@/pages/uns/types';
import { useBaseStore } from '@/stores/base';
import { App, Button, Flex } from 'antd';
import type { ItemType } from 'antd/es/menu/interface';
import { type FC, type Key, type KeyboardEvent, useCallback, useEffect, useMemo, useState } from 'react';
import RecycleDeleteModal from './RecycleDeleteModal';
import { getNodeTopicType, isTopicTypeFolder, UnsTreeNodeIcon } from './tree-icons';

type DeleteTarget = {
  type: 'one' | 'all';
  node?: UnsTreeNode;
};

interface RecycleTreeProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onRecycleChanged?: () => void;
}

const nodeLabel = (node?: UnsTreeNode) =>
  String(node?.title || node?.pathName || node?.name || node?.displayName || node?.alias || node?.path || '').trim();

const isRecycleActionableNode = (node?: UnsTreeNode) =>
  Number(node?.deletedTime || 0) > 0 && Number(node?.recycleIsDel ?? 2) === 2;

const isRecycleContextNode = (node?: UnsTreeNode) => !!node?.id && !isRecycleActionableNode(node);

const countRecycleFiles = (node?: UnsTreeNode): number =>
  (node?.children || []).reduce((total, child) => {
    const childCount = child.pathType === 2 && isRecycleActionableNode(child) ? 1 : 0;
    return total + childCount + countRecycleFiles(child);
  }, 0);

const RecycleNodeIcon: FC<{ dataNode: UnsTreeNode; expandedKeys: Key[] }> = ({ dataNode }) => {
  const {
    systemInfo: { enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));

  const mutedClass = isRecycleContextNode(dataNode) ? 'recycle-node-icon-muted' : '';

  return (
    <span className={mutedClass}>
      <UnsTreeNodeIcon
        dataNode={dataNode}
        topicType={enableAutoCategorization ? getNodeTopicType(dataNode) : 0}
        enableAutoCategorization={enableAutoCategorization}
        isTopicTypeFolder={isTopicTypeFolder}
      />
    </span>
  );
};

const RecycleTree: FC<RecycleTreeProps> = ({ open, onOpenChange, onRecycleChanged }) => {
  const formatMessage = useTranslate();
  const { message, modal } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [treeData, setTreeData] = useState<UnsTreeNode[]>([]);
  const [expandedKeys, setExpandedKeys] = useState<Key[]>([]);
  const [selectedNode, setSelectedNode] = useState<UnsTreeNode>();
  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget>();

  const loadRecycleData = useCallback(() => {
    setLoading(true);
    getRecycleTreeData()
      .then((data) => {
        setTreeData(data);
      })
      .catch(() => {
        setTreeData([]);
      })
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (open) {
      loadRecycleData();
    }
  }, [open, loadRecycleData]);

  const handleRestore = useCallback(
    async (node?: UnsTreeNode) => {
      if (!node?.id || !isRecycleActionableNode(node)) return;
      await restoreRecycleNode({ id: String(node.id), confirm: true });
      message.success(formatMessage('uns.recycleRestoreSuccess'));
      setSelectedNode(undefined);
      loadRecycleData();
      onRecycleChanged?.();
    },
    [formatMessage, loadRecycleData, message, onRecycleChanged]
  );

  const openRestoreConfirm = useCallback(
    (node?: UnsTreeNode) => {
      const okText = formatMessage('uns.recycleRestore');
      const cancelText = formatMessage('common.cancel');

      modal.confirm({
        title: formatMessage('uns.recycleRestoreTitle'),
        width: 420,
        icon: null,
        content: (
          <Flex align="flex-start" gap={8} style={{ marginTop: 4 }}>
            <WarningFilled style={{ color: '#faad14', flexShrink: 0, width: 20, height: 20 }} />
            <span style={{ overflowWrap: 'anywhere' }}>
              {formatMessage('uns.recycleRestoreDesc', { name: nodeLabel(node) })}
            </span>
          </Flex>
        ),
        okText,
        cancelText,
        okButtonProps: {
          title: okText,
        },
        cancelButtonProps: {
          title: cancelText,
        },
        onOk: () => handleRestore(node),
      });
    },
    [formatMessage, handleRestore, modal]
  );

  const openDeleteConfirm = useCallback((target: DeleteTarget) => {
    setDeleteTarget(target);
  }, []);

  const handleDelete = useCallback(
    async (deleteWithFlows: boolean) => {
      if (!deleteTarget) return;
      if (deleteTarget.type === 'all') {
        await clearRecycleTree({ deleteFlow: deleteWithFlows });
        message.success(formatMessage('uns.recycleClearSuccess'));
      } else if (deleteTarget.node?.id) {
        await forceDeleteRecycleNode({ id: String(deleteTarget.node.id), deleteFlow: deleteWithFlows });
        message.success(formatMessage('common.deleteSuccessfully'));
      }
      setDeleteTarget(undefined);
      setSelectedNode(undefined);
      loadRecycleData();
      onRecycleChanged?.();
    },
    [deleteTarget, formatMessage, loadRecycleData, message, onRecycleChanged]
  );

  const deleteTitle = deleteTarget?.type === 'all' ? 'uns.recycleClearAllTitle' : 'uns.recycleDeleteTitle';
  const deleteDescription =
    deleteTarget?.type === 'all'
      ? formatMessage('uns.recycleClearAllDesc')
      : formatMessage('uns.recycleDeleteDesc', { name: nodeLabel(deleteTarget?.node) });

  const rightClickMenuItems = useCallback(
    ({ node }: { node?: UnsTreeNode }): ItemType[] => {
      if (!node?.id || !isRecycleActionableNode(node as UnsTreeNode)) return [];
      return [
        {
          key: 'restore',
          label: formatMessage('uns.recycleRestore'),
          onClick: () => openRestoreConfirm(node),
        },
        { type: 'divider' },
        {
          key: 'delete',
          label: formatMessage('common.delete'),
          danger: true,
          extra: (
            <Flex align="center">
              <TrashCan {...treeIconProps} />
            </Flex>
          ),
          onClick: () => openDeleteConfirm({ type: 'one', node }),
        },
      ];
    },
    [formatMessage, openDeleteConfirm, openRestoreConfirm]
  );

  const selectedKeys = useMemo(() => (selectedNode?.key ? [selectedNode.key] : []), [selectedNode]);
  const toggleRecyclePanel = useCallback(() => {
    onOpenChange(!open);
  }, [onOpenChange, open]);

  const handleHeaderKeyDown = useCallback(
    (event: KeyboardEvent<HTMLDivElement>) => {
      if (event.target !== event.currentTarget) return;
      if (event.key !== 'Enter' && event.key !== ' ') return;
      event.preventDefault();
      toggleRecyclePanel();
    },
    [toggleRecyclePanel]
  );

  return (
    <div className={`recycle-panel ${open ? 'is-open' : ''}`}>
      <div
        className="recycle-panel-header"
        role="button"
        tabIndex={0}
        onClick={toggleRecyclePanel}
        onKeyDown={handleHeaderKeyDown}
      >
        <span className="recycle-panel-title">
          <TrashCan size={18} strokeWidth={1.75} aria-hidden />
          <span className="recycle-panel-title-text">{formatMessage('uns.recycleBin')}</span>
        </span>
        <span className="recycle-panel-actions">
          {open && (
            <Button
              danger
              size="small"
              className="recycle-clear-all"
              onClick={(event) => {
                event.stopPropagation();
                openDeleteConfirm({ type: 'all' });
              }}
            >
              {formatMessage('uns.recycleClearAll')}
            </Button>
          )}
          {open ? <ChevronDown size={16} /> : <ChevronUp size={16} />}
        </span>
      </div>
      {open && (
        <div className="recycle-panel-content">
          <ProTree
            loading={loading}
            empty={
              <div className="recycle-empty">
                <ComEmpty />
              </div>
            }
            treeData={treeData}
            selectedKeys={selectedKeys}
            showSwitcherIcon
            titleRender={(dataNode) => (
              <span className={isRecycleContextNode(dataNode as UnsTreeNode) ? 'recycle-node-title is-context' : ''}>
                {nodeLabel(dataNode as UnsTreeNode)}
              </span>
            )}
            onSelect={(_, { node, selected }) => {
              const nextNode = node as UnsTreeNode;
              setSelectedNode(selected && isRecycleActionableNode(nextNode) ? { ...nextNode } : undefined);
            }}
            expandedKeys={expandedKeys}
            onExpand={(keys) => setExpandedKeys(keys)}
            rightClickMenuItems={rightClickMenuItems as any}
            treeNodeIcon={(dataNode) => (
              <RecycleNodeIcon dataNode={dataNode as UnsTreeNode} expandedKeys={expandedKeys} />
            )}
            treeNodeExtra={(dataNode) =>
              isRecycleActionableNode(dataNode as UnsTreeNode) ? (
                <span title={formatMessage('uns.recycleRestore')}>
                  <Undo
                    size={15}
                    className="recycle-restore-action"
                    onClick={(event: any) => {
                      event.stopPropagation();
                      openRestoreConfirm(dataNode as UnsTreeNode);
                    }}
                  />
                </span>
              ) : null
            }
            treeNodeCount={(dataNode) =>
              dataNode.pathType === 0 ? (
                <span className="recycle-node-count">({countRecycleFiles(dataNode as UnsTreeNode)})</span>
              ) : null
            }
            wrapperStyle={{ padding: '8px 14px 0' }}
            height={0}
          />
        </div>
      )}
      <RecycleDeleteModal
        open={!!deleteTarget}
        title={formatMessage(deleteTitle)}
        description={deleteDescription}
        onCancel={() => setDeleteTarget(undefined)}
        onSubmit={handleDelete}
      />
    </div>
  );
};

export default RecycleTree;
