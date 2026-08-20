import { ButtonPermission } from '@/common-types/button-permission';
import { ProSearch } from '@/components';
import ComEmpty from '@/components/com-empty';
import HighlightText from '@/components/pro-tree/HighlightText';
import { useTranslate } from '@/hooks';
import type { FolderNode, NotebookTabType, NotebookTreeItem } from '@/pages/notebook/types';
import { hasPermission } from '@/utils/auth';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { AddLarge, Catalog, Edit, Folder, Renew, TrashCan } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import toolbarStyles from '@/components/lucide-icon/toolbar-icon.module.scss';
import { App, Button, Dropdown, Tooltip, Tree } from 'antd';
import type { DataNode } from 'antd/es/tree';
import type { MenuProps } from 'antd';
import { type CSSProperties, type FC, type ReactNode, useEffect, useMemo, useState } from 'react';
import sidebarStyles from './NotebookSidebar.module.scss';

const FOLDER_KEY_PREFIX = 'folder-';
const NOTEBOOK_KEY_PREFIX = 'notebook-';

const folderKey = (id: number) => `${FOLDER_KEY_PREFIX}${id}`;
const notebookKey = (id: number) => `${NOTEBOOK_KEY_PREFIX}${id}`;

const notebookMatches = (item: NotebookTreeItem, keyword: string) =>
  item.name.toLowerCase().includes(keyword) || item.path.toLowerCase().includes(keyword);

const treeTitleTextStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
};

const treeTitleRowStyle: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: 6,
  width: '100%',
  minWidth: 0,
  whiteSpace: 'nowrap',
  overflow: 'hidden',
};

interface NotebookSidebarProps {
  folders: FolderNode[];
  rootNotebooks: NotebookTreeItem[];
  runningNotebookIds?: number[];
  selectedFolderId: number;
  tab: NotebookTabType;
  searchValue: string;
  onSearchChange: (value: string) => void;
  onSelect: (folderId: number) => void;
  onOpenNotebook: (notebookId: number) => void;
  onNewFolder: () => void;
  onNewNotebook: () => void;
  onRefresh: () => void;
  onEditFolder: (folderId: number) => void;
  onEditNotebook: (notebookId: number) => void;
  onDeleteFolder: (folderId: number) => void | Promise<void>;
  onDeleteNotebook: (notebookId: number) => void | Promise<void>;
}

const filterTree = (nodes: FolderNode[], keyword: string): FolderNode[] => {
  const normalized = keyword.trim().toLowerCase();
  if (!normalized) return nodes;
  return nodes
    .map((node) => {
      const children = filterTree(node.children || [], keyword);
      const notebooks = (node.notebooks || []).filter((item) => notebookMatches(item, normalized));
      if (node.name.toLowerCase().includes(normalized) || children.length > 0 || notebooks.length > 0) {
        return { ...node, children, notebooks };
      }
      return null;
    })
    .filter(Boolean) as FolderNode[];
};

const filterRunningTree = (nodes: FolderNode[], runningNotebookIdSet: Set<number>): FolderNode[] =>
  nodes
    .map((node) => {
      const children = filterRunningTree(node.children || [], runningNotebookIdSet);
      const notebooks = (node.notebooks || []).filter((item) => runningNotebookIdSet.has(item.id));
      if (children.length > 0 || notebooks.length > 0) {
        return { ...node, children, notebooks };
      }
      return null;
    })
    .filter(Boolean) as FolderNode[];

const buildFolderNotebookCountMap = (nodes: FolderNode[]): Map<number, number> => {
  const countMap = new Map<number, number>();

  const walk = (node: FolderNode): number => {
    const ownNotebookCount = node.notebooks?.length ?? node.notebookCount ?? 0;
    const childrenNotebookCount = (node.children || []).reduce((sum, child) => sum + walk(child), 0);
    const totalNotebookCount = ownNotebookCount + childrenNotebookCount;

    countMap.set(node.id, totalNotebookCount);
    return totalNotebookCount;
  };

  nodes.forEach(walk);
  return countMap;
};

const toNotebookTreeNode = (
  item: NotebookTreeItem,
  actions?: { onEdit?: (id: number) => void; onDelete?: (id: number) => void; modal?: any; formatMessage?: any },
  searchValue = ''
): DataNode => ({
  key: notebookKey(item.id),
  title: (
    <span style={treeTitleRowStyle}>
      <Catalog size={16} style={{ flexShrink: 0 }} />
      <Tooltip title={item.name}>
        <span style={treeTitleTextStyle} title={item.name}>
          <HighlightText needle={searchValue} haystack={item.name} />
        </span>
      </Tooltip>
      {actions && (
        <span
          className="nodeExtra"
          style={{ display: 'none', alignItems: 'center', gap: 4, flexShrink: 0, marginLeft: 'auto' }}
          onClick={(e) => e.stopPropagation()}
        >
          <Edit size={16} style={{ cursor: 'pointer' }} onClick={() => actions.onEdit?.(item.id)} />
          <TrashCan
            size={16}
            style={{ cursor: 'pointer' }}
            onClick={() => {
              actions.modal?.confirm(
                mergeDeleteConfirmProps(
                  {
                    title: actions.formatMessage?.('Notebook.deleteNotebookConfirm', {}, 'Delete this notebook?'),
                    onOk: () => actions.onDelete?.(item.id),
                  },
                  actions.formatMessage
                )
              );
            }}
          />
        </span>
      )}
    </span>
  ),
  isLeaf: true,
});

const toTreeData = (
  nodes: FolderNode[],
  folderNotebookCountMap: Map<number, number>,
  searchValue: string,
  notebookActions?: { onDelete?: (id: number) => void; modal?: any; formatMessage?: any },
  folderActions?: { onEdit?: (id: number) => void; onDelete?: (id: number) => void; modal?: any; formatMessage?: any }
): DataNode[] =>
  nodes.map((node) => {
    const childFolders = toTreeData(
      node.children || [],
      folderNotebookCountMap,
      searchValue,
      notebookActions,
      folderActions
    );
    const childNotebooks = (node.notebooks || []).map((item) => toNotebookTreeNode(item, notebookActions, searchValue));

    return {
      key: folderKey(node.id),
      title: (
        <span style={treeTitleRowStyle}>
          <Folder size={16} style={{ flexShrink: 0 }} />
          <Tooltip title={node.name}>
            <span style={treeTitleTextStyle} title={node.name}>
              <HighlightText needle={searchValue} haystack={node.name} />
            </span>
          </Tooltip>
          <span style={{ color: 'var(--ui-icon-color)', fontSize: 12, flexShrink: 0 }}>
            ({folderNotebookCountMap.get(node.id) || 0})
          </span>
          {folderActions && (
            <span
              className="nodeExtra"
              style={{ display: 'none', alignItems: 'center', gap: 4, flexShrink: 0, marginLeft: 'auto' }}
              onClick={(e) => e.stopPropagation()}
            >
              <Edit size={16} style={{ cursor: 'pointer' }} onClick={() => folderActions.onEdit?.(node.id)} />
              <TrashCan
                size={16}
                style={{ cursor: 'pointer' }}
                onClick={() => {
                  folderActions.modal?.confirm(
                    mergeDeleteConfirmProps(
                      {
                        title: folderActions.formatMessage?.('Notebook.deleteFolderConfirm', {}, 'Delete this folder?'),
                        onOk: () => folderActions.onDelete?.(node.id),
                      },
                      folderActions.formatMessage
                    )
                  );
                }}
              />
            </span>
          )}
        </span>
      ),
      children: [...childFolders, ...childNotebooks],
    };
  });

const collectExpandedKeys = (nodes: FolderNode[]): string[] => {
  const keys: string[] = [];
  const walk = (items: FolderNode[]) => {
    items.forEach((node) => {
      keys.push(folderKey(node.id));
      if (node.children?.length) {
        walk(node.children);
      }
    });
  };
  walk(nodes);
  return keys;
};

const findAncestorKeys = (nodes: FolderNode[], targetId: number, parents: string[] = []): string[] => {
  for (const node of nodes) {
    if (node.id === targetId) {
      return parents;
    }
    const next = findAncestorKeys(node.children || [], targetId, [...parents, folderKey(node.id)]);
    if (next.length) {
      return next;
    }
  }
  return [];
};

const findFolderNode = (nodes: FolderNode[], targetId: number): FolderNode | null => {
  for (const node of nodes) {
    if (node.id === targetId) {
      return node;
    }
    const found = findFolderNode(node.children || [], targetId);
    if (found) {
      return found;
    }
  }
  return null;
};

const collectFolderBranchKeys = (node: FolderNode): string[] => {
  const keys = [folderKey(node.id)];
  for (const child of node.children || []) {
    keys.push(...collectFolderBranchKeys(child));
  }
  return keys;
};

const NotebookSidebar: FC<NotebookSidebarProps> = ({
  folders,
  rootNotebooks,
  runningNotebookIds = [],
  selectedFolderId,
  tab,
  searchValue,
  onSearchChange,
  onSelect,
  onOpenNotebook,
  onNewFolder,
  onNewNotebook,
  onRefresh,
  onEditFolder,
  onEditNotebook,
  onDeleteFolder,
  onDeleteNotebook,
}) => {
  const formatMessage = useTranslate();
  const { modal } = App.useApp();
  const canManage = hasPermission(ButtonPermission['Notebook.manage']);
  const runningNotebookIdSet = useMemo(() => new Set(runningNotebookIds), [runningNotebookIds]);
  const visibleFolders = useMemo(
    () => (tab === 'running' ? filterRunningTree(folders, runningNotebookIdSet) : folders),
    [folders, runningNotebookIdSet, tab]
  );
  const visibleRootNotebooks = useMemo(
    () => (tab === 'running' ? rootNotebooks.filter((item) => runningNotebookIdSet.has(item.id)) : rootNotebooks),
    [rootNotebooks, runningNotebookIdSet, tab]
  );
  const filteredFolders = useMemo(() => filterTree(visibleFolders, searchValue), [visibleFolders, searchValue]);
  const filteredRootNotebooks = useMemo(() => {
    const normalized = searchValue.trim().toLowerCase();
    if (!normalized) return visibleRootNotebooks;
    return visibleRootNotebooks.filter((item) => notebookMatches(item, normalized));
  }, [visibleRootNotebooks, searchValue]);
  const folderNotebookCountMap = useMemo(() => buildFolderNotebookCountMap(visibleFolders), [visibleFolders]);
  const notebookActions = useMemo(
    () => (canManage ? { onEdit: onEditNotebook, onDelete: onDeleteNotebook, modal, formatMessage } : undefined),
    [canManage, onEditNotebook, onDeleteNotebook, modal, formatMessage]
  );
  const folderActions = useMemo(
    () => (canManage ? { onEdit: onEditFolder, onDelete: onDeleteFolder, modal, formatMessage } : undefined),
    [canManage, onEditFolder, onDeleteFolder, modal, formatMessage]
  );
  const treeData = useMemo(
    () => [
      ...toTreeData(filteredFolders, folderNotebookCountMap, searchValue, notebookActions, folderActions),
      ...filteredRootNotebooks.map((item) => toNotebookTreeNode(item, notebookActions, searchValue)),
    ],
    [filteredFolders, filteredRootNotebooks, folderNotebookCountMap, searchValue, notebookActions, folderActions]
  );
  const [manualExpandedKeys, setManualExpandedKeys] = useState<string[]>([]);

  useEffect(() => {
    if (searchValue.trim() || selectedFolderId <= 0) {
      return;
    }
    const ancestorKeys = findAncestorKeys(visibleFolders, selectedFolderId);
    setManualExpandedKeys((prev) => Array.from(new Set([...prev, ...ancestorKeys, folderKey(selectedFolderId)])));
  }, [visibleFolders, searchValue, selectedFolderId]);

  const expandedKeys = useMemo(() => {
    if (searchValue.trim()) {
      return collectExpandedKeys(filteredFolders);
    }
    return manualExpandedKeys;
  }, [filteredFolders, manualExpandedKeys, searchValue]);

  const getFolderMenuItems = (node: FolderNode): MenuProps['items'] => {
    const currentKey = folderKey(node.id);
    const branchKeys = collectFolderBranchKeys(node);
    const isExpanded = expandedKeys.includes(currentKey);

    return [
      {
        key: 'expand',
        label: formatMessage('Notebook.expandFolder', {}, 'Expand Folder'),
        onClick: () => {
          setManualExpandedKeys((prev) => Array.from(new Set([...prev, ...branchKeys])));
        },
      },
      {
        key: 'collapse',
        label: formatMessage('Notebook.collapseFolder', {}, 'Collapse Folder'),
        onClick: () => {
          setManualExpandedKeys((prev) => prev.filter((key) => !branchKeys.includes(key)));
        },
        disabled: !isExpanded,
      },
      ...(canManage
        ? [
            { type: 'divider' as const },
            {
              key: 'delete',
              label: formatMessage('common.delete', {}, 'Delete'),
              danger: true,
              onClick: () => {
                modal.confirm(
                  mergeDeleteConfirmProps(
                    {
                      title: formatMessage('Notebook.deleteFolderConfirm', {}, 'Delete this folder?'),
                      onOk: () => onDeleteFolder(node.id),
                    },
                    formatMessage
                  )
                );
              },
            },
          ]
        : []),
    ];
  };

  return (
    <div className={sidebarStyles.sidebarBody}>
      <div className={sidebarStyles.sidebarToolbar}>
        <ProSearch
          value={searchValue}
          closeButtonLabelText={formatMessage('common.clearSearchInput', {}, 'Clear search input')}
          placeholder={formatMessage('Notebook.searchPlaceholder', {}, 'Search')}
          onChange={(event) => onSearchChange(event.target.value)}
          onClear={() => onSearchChange('')}
          title={searchValue || formatMessage('Notebook.searchPlaceholder', {}, 'Search')}
          size="sm"
          className={sidebarStyles.sidebarSearch}
        />
        {canManage && (
          <Dropdown
            menu={{
              items: [
                {
                  key: 'newFolder',
                  label: formatMessage('Notebook.newFolder', {}, 'New Folder'),
                  onClick: onNewFolder,
                },
                {
                  key: 'newNotebook',
                  label: formatMessage('Notebook.newNotebook', {}, 'New Notebook'),
                  onClick: onNewNotebook,
                },
              ] as MenuProps['items'],
            }}
            placement="bottom"
          >
            <Button
              className={`${toolbarStyles['toolbar-icon-btn']} ${sidebarStyles.sidebarIconBtn}`}
              color="default"
              variant="filled"
            >
              <AddLarge {...toolbarIconProps} />
            </Button>
          </Dropdown>
        )}
        <Button
          className={`${toolbarStyles['toolbar-icon-btn']} ${sidebarStyles.sidebarIconBtn}`}
          color="default"
          variant="filled"
          onClick={onRefresh}
        >
          <Renew {...toolbarIconProps} />
        </Button>
      </div>
      <div className={sidebarStyles.sidebarDivider} />
      <div className={sidebarStyles.treeWrap}>
        <Tree
          blockNode
          expandedKeys={expandedKeys}
          autoExpandParent={false}
          selectedKeys={selectedFolderId > 0 ? [folderKey(selectedFolderId)] : []}
          treeData={treeData.map(function attachMenu(node: DataNode): DataNode {
            if (!String(node.key).startsWith(FOLDER_KEY_PREFIX)) {
              return node;
            }

            const folderId = Number(String(node.key).slice(FOLDER_KEY_PREFIX.length));
            const folderNode = findFolderNode(filteredFolders.length ? filteredFolders : visibleFolders, folderId);
            if (!folderNode) {
              return node;
            }

            return {
              ...node,
              title: (
                <Dropdown menu={{ items: getFolderMenuItems(folderNode) }} trigger={['contextMenu']}>
                  <span>{node.title as ReactNode}</span>
                </Dropdown>
              ),
              children: node.children?.map(attachMenu) as DataNode[] | undefined,
            };
          })}
          onExpand={(keys) => setManualExpandedKeys(keys as string[])}
          onSelect={(keys) => {
            const key = String(keys[0] ?? '');
            if (!key) {
              onSelect(0);
              return;
            }

            if (key.startsWith(NOTEBOOK_KEY_PREFIX)) {
              const notebookId = Number(key.slice(NOTEBOOK_KEY_PREFIX.length));
              if (!Number.isNaN(notebookId)) {
                onOpenNotebook(notebookId);
              }
              return;
            }

            const next = Number(key.slice(FOLDER_KEY_PREFIX.length));
            onSelect(Number.isNaN(next) ? 0 : next);
          }}
        />
        {folders.length === 0 && rootNotebooks.length === 0 && (
          <ComEmpty description={formatMessage('Notebook.noFolders', {}, 'No folders yet')} />
        )}
      </div>
    </div>
  );
};

export default NotebookSidebar;
