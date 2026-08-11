import { useEffect, useMemo, useState, type FC } from 'react';
import { Typography } from 'antd';
import type { Key } from 'react';
import type { FleetUNSScopeState } from '@/apis/core-api/fleet';
import ComEmpty from '@/components/com-empty';
import ProTree from '@/components/pro-tree';
import { useTranslate } from '@/hooks';
import { ModelDetail, TopicDetail } from '@/pages/uns/components';
import {
  getFolderTopicType,
  getNodeTopicType,
  getUnsTopicFolderTitleStyle,
  isTopicTypeFolder,
  UnsTreeNodeIcon,
} from '@/pages/uns/components/uns-tree/tree-icons';
import { UnsContextProvider } from '@/pages/uns/UnsContext';
import type { InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import { useBaseStore } from '@/stores/base';
import '@/pages/uns/index.scss';
import styles from './FleetUNSSyncBrowser.module.scss';

interface FleetUNSSyncBrowserProps {
  scopes: FleetUNSScopeState[];
  treeData: UnsTreeNode[];
}

const filterSyncedTree = (nodes: UnsTreeNode[], syncedNodeIDs: Set<string>): UnsTreeNode[] =>
  nodes.flatMap((node) => {
    if (!syncedNodeIDs.has(String(node.id ?? node.key))) {
      return [];
    }
    const children = filterSyncedTree(node.children || [], syncedNodeIDs);
    return [
      {
        ...node,
        children,
        countChildren: children.length,
        hasChildren: children.length > 0,
        isLeaf: node.pathType === 2 || children.length === 0,
      },
    ];
  });

const findNode = (nodes: UnsTreeNode[], key?: Key): UnsTreeNode | undefined => {
  if (key === undefined) return undefined;
  for (const node of nodes) {
    if (String(node.key) === String(key)) return node;
    const child = findNode(node.children || [], key);
    if (child) return child;
  }
  return undefined;
};

const noopInitTreeData: InitTreeDataFnType = () => {};
const noopHandleDelete = () => {};

const FleetUNSSyncBrowser: FC<FleetUNSSyncBrowserProps> = ({ scopes, treeData }) => {
  const formatMessage = useTranslate();
  const enableAutoCategorization = useBaseStore((state) => state.systemInfo.enableAutoCategorization);
  const [selectedKey, setSelectedKey] = useState<Key>();
  const [expandedKeys, setExpandedKeys] = useState<Key[]>([]);

  const syncedTree = useMemo(() => {
    const syncedNodeIDs = new Set(
      scopes
        .filter((scope) => scope.status === 'active')
        .flatMap((scope) => scope.targetNodeIDs || [])
        .map(String)
    );
    return filterSyncedTree(treeData, syncedNodeIDs);
  }, [scopes, treeData]);

  useEffect(() => {
    setSelectedKey((currentKey) => findNode(syncedTree, currentKey)?.key ?? syncedTree[0]?.key);
    setExpandedKeys(syncedTree.map((node) => node.key));
  }, [syncedTree]);

  const selectedNode = useMemo(() => findNode(syncedTree, selectedKey), [selectedKey, syncedTree]);

  return (
    <div className={styles.browser}>
      <aside className={styles['tree-pane']}>
        <Typography.Title level={5} className={styles['pane-title']}>
          {formatMessage('uns.treeList')}
        </Typography.Title>
        <div className={styles['tree-body']}>
          <ProTree
            treeData={syncedTree}
            selectedKeys={selectedKey === undefined ? [] : [selectedKey]}
            expandedKeys={expandedKeys}
            onExpand={setExpandedKeys}
            onSelect={(_, info) => setSelectedKey((info.node as UnsTreeNode).key)}
            showSwitcherIcon
            empty={<ComEmpty description={formatMessage('uns.noData')} />}
            wrapperStyle={{ padding: '0 8px 12px' }}
            treeNodeIcon={(dataNode) => (
              <UnsTreeNodeIcon
                dataNode={dataNode as UnsTreeNode}
                topicType={enableAutoCategorization ? getNodeTopicType(dataNode as UnsTreeNode) : 0}
                enableAutoCategorization={enableAutoCategorization}
                isTopicTypeFolder={isTopicTypeFolder}
              />
            )}
            treeNodeCount={(dataNode) =>
              dataNode.pathType === 0 ? (
                <span className={styles['node-count']}>({dataNode.countChildren || 0})</span>
              ) : null
            }
            renderTitleStyle={(dataNode) => {
              const topicType = getFolderTopicType(dataNode as UnsTreeNode);
              return topicType && enableAutoCategorization ? getUnsTopicFolderTitleStyle(topicType) : {};
            }}
          />
        </div>
      </aside>

      <section className={`${styles['detail-pane']} unsContainer`}>
        {!selectedNode ? (
          <div className={styles['empty-detail']}>
            <ComEmpty description={formatMessage('uns.noData')} />
          </div>
        ) : (
          <UnsContextProvider>
            <div className="chartWrap">
              {selectedNode.pathType === 0 ? (
                <ModelDetail currentNode={selectedNode} initTreeData={noopInitTreeData} readOnly />
              ) : selectedNode.pathType === 2 ? (
                <TopicDetail
                  currentNode={selectedNode}
                  initTreeData={noopInitTreeData}
                  handleDelete={noopHandleDelete}
                  readOnly
                />
              ) : (
                <div className={styles['empty-detail']}>
                  <ComEmpty description={formatMessage('uns.noData')} />
                </div>
              )}
            </div>
          </UnsContextProvider>
        )}
      </section>
    </div>
  );
};

export default FleetUNSSyncBrowser;
