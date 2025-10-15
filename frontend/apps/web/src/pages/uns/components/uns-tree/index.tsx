import { type Key, memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { App, Button, Divider, Flex, Popover, Radio } from 'antd';
import { getTreeStoreSnapshot, ROOT_NODE_ID, useTreeStore, useTreeStoreRef } from '../../store/treeStore';
import type { TTreeType } from '../../store/types';
import useTranslate from '@/hooks/useTranslate';
import {
  DocumentAdd,
  Filter,
  Folder,
  FolderAdd,
  FolderOpen,
  GroupObjectsNew,
  Renew,
  RepoArtifact,
  Document,
  TrashCan,
} from '@carbon/icons-react';
import Icon from '@ant-design/icons';
import { ButtonPermission } from '@/common-types/button-permission';
import { debounce } from 'lodash-es';
import type { ItemType } from 'antd/es/menu/interface';
import { getInstanceInfo, getModelInfo, pasteUns } from '@/apis/inter-api/uns';
import { useViewLabelModal } from '@/pages/uns/components';
import ReverseModal from '@/pages/uns/components/reverse-modal';
import type { UnsTreeNode } from '@/pages/uns/types';
import { filterPermissionToList } from '@/utils/auth';
import ComClickTrigger from '@/components/com-click-trigger';
import ProSearch from '@/components/pro-search';
import ProTree, { type ProTreeProps } from '@/components/pro-tree';
import TagAdd from '@/components/svg-components/TagAdd';
import { usePageIsShow } from '@/contexts/tabs-lifecycle-context.ts';
import { useUnsContext } from '@/pages/uns/UnsContext';
import StatusDot from './StatusDot';

const renderOperationDom = (type: string) => {
  switch (type) {
    case 'DocumentAdd':
      return <DocumentAdd />;
    case 'FolderAdd':
      return <FolderAdd />;
    case 'RepoArtifact':
      return <RepoArtifact />;
    case 'GroupObjectsNew':
      return <GroupObjectsNew />;
    case 'TagAdd':
      return <Icon component={TagAdd} />;
    case 'Renew':
      return <Renew />;
    default:
      return null;
  }
};

// 操作
const Operation = () => {
  const formatMessage = useTranslate();
  const { treeType, operationFns, setSelectedNode, selectedNode, loadData } = useTreeStore((state) => ({
    treeType: state.treeType,
    operationFns: state.operationFns,
    setSelectedNode: state.setSelectedNode,
    selectedNode: state.selectedNode,
    loadData: state.loadData,
  }));
  const { message } = App.useApp();
  const [reverserOpen, setReverserOpen] = useState<boolean>(false); //复制的topic
  const options = useMemo(() => {
    return filterPermissionToList<{
      onClick: () => void;
      buttonType: string;
      showTreeType?: TTreeType;
      key: string;
      title?: string;
      id?: string;
      disabled?: boolean;
    }>([
      {
        title: formatMessage('uns.newFile'),
        auth: ButtonPermission['uns.fileAdd'],
        onClick: () => {
          operationFns?.setOptionOpen?.('addFile', selectedNode);
        },
        buttonType: 'DocumentAdd',
        showTreeType: 'uns',
        key: 'addFile',
        disabled: !!selectedNode?.mount,
      },
      {
        title: formatMessage('uns.newFolder'),
        auth: ButtonPermission['uns.folderAdd'],
        onClick: () => {
          operationFns?.setOptionOpen?.('addFolder', selectedNode);
        },
        buttonType: 'FolderAdd',
        showTreeType: 'uns',
        key: 'addFolder',
        disabled: !!selectedNode?.mount,
      },
      {
        title: formatMessage('uns.batchGeneration'),
        auth: ButtonPermission['uns.batchGeneration'],
        onClick: () => {
          setReverserOpen(true);
        },
        buttonType: 'RepoArtifact',
        showTreeType: 'uns',
        key: 'batchGeneration',
        disabled: !!selectedNode?.mount,
      },
      {
        title: formatMessage('uns.addTemplate'),
        auth: ButtonPermission['uns.templateAdd'],
        onClick: () => {
          operationFns?.openTemplateModal?.('addTemplate', setSelectedNode);
        },
        buttonType: 'GroupObjectsNew',
        showTreeType: 'template',
        key: 'addTemplate',
      },
      {
        title: formatMessage('uns.newLabel'),
        auth: ButtonPermission['uns.labelAdd'],
        onClick: () => {
          operationFns?.setLabelOpen?.();
        },
        buttonType: 'TagAdd',
        showTreeType: 'label',
        key: 'addLabel',
      },
      {
        title: formatMessage('common.refresh'),
        onClick: () => {
          loadData({ reset: true, clearSelect: true }, () => {
            message.success(formatMessage('common.refreshSuccessful'));
          });
        },
        buttonType: 'Renew',
        key: 'reNew',
      },
    ])?.filter((item) => !item.showTreeType || item.showTreeType === treeType);
  }, [treeType, operationFns, loadData, selectedNode]);
  return (
    <>
      {options?.map((item) => (
        <span
          style={{
            height: 32,
            width: 32,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            cursor: item.disabled ? 'not-allowed' : 'pointer',
            color: item.disabled ? 'var(--supos-select-d-color)' : 'var(--supos-text-color)',
          }}
          id={item.id}
          onClick={item.disabled ? undefined : item.onClick}
          key={item.key}
          title={item.title}
        >
          {renderOperationDom(item.buttonType)}
        </span>
      ))}
      {reverserOpen && (
        <ReverseModal
          reverserOpen={reverserOpen}
          setReverserOpen={setReverserOpen}
          currentNode={selectedNode as any}
          initTreeData={loadData}
        />
      )}
    </>
  );
};

// 树类型
const TreeTab = () => {
  const formatMessage = useTranslate();

  const { treeType } = useTreeStore((state) => ({
    treeType: state.treeType,
  }));
  const stateRef = useTreeStoreRef();
  const { setTreeType, setSearchValue, setLazyTree, loadData } = getTreeStoreSnapshot(stateRef, (state) => ({
    setTreeType: state.setTreeType,
    setSearchValue: state.setSearchValue,
    setLazyTree: state.setLazyTree,
    loadData: state.loadData,
  }));

  return (
    <Flex justify="space-between" align="center">
      <Radio.Group
        onChange={(e) => {
          setSearchValue('');
          setTreeType(e.target.value);
          loadData({ reset: true, clearSelect: true });
        }}
        optionType="button"
        value={treeType}
        style={{ padding: '16px 0' }}
        size="small"
        options={[
          {
            label: formatMessage('uns.tree'),
            value: 'uns',
            title: formatMessage('uns.tree'),
          },
          {
            label: formatMessage('common.template'),
            value: 'template',
            title: formatMessage('common.template'),
          },
          {
            label: formatMessage('common.label'),
            value: 'label',
            title: formatMessage('common.label'),
          },
        ]}
      />
      {treeType === 'uns' && (
        <ComClickTrigger
          style={{ flex: 1, height: 24 }}
          onTrigger={() => {
            setLazyTree((pre) => !pre);
            loadData({ reset: true, clearSelect: true });
          }}
        />
      )}
      {
        <ComClickTrigger
          triggerCount={2}
          style={{ flex: 1, height: 24 }}
          onTrigger={() => {
            console.warn(getTreeStoreSnapshot(stateRef));
          }}
        />
      }
    </Flex>
  );
};

// 搜索
const Search = () => {
  const formatMessage = useTranslate();
  const { searchValue, setSearchValue, loadData, searchType, treeType } = useTreeStore((state) => ({
    searchValue: state.searchValue,
    setSearchValue: state.setSearchValue,
    loadData: state.loadData,
    searchType: state.searchType,
    treeType: state.treeType,
  }));

  const debouncedInitTreeData = useCallback(
    // eslint-disable-next-line react-hooks/use-memo
    debounce(() => {
      loadData({ reset: true, clearSelect: true, queryType: 'search' });
    }, 500),
    [loadData, searchType] // 仅当 initTreeData 或 searchType 发生变化时更新 debounce 函数
  );

  const onSearchChange = () => debouncedInitTreeData();
  const selectRef = useRef<any>(null);
  const isComposingRef = useRef(false); // 拼音输入中..
  useEffect(() => {
    //处理拼音输入法
    const inputElement = selectRef.current;

    if (inputElement) {
      const handleCompositionStart = () => {
        isComposingRef.current = true;
      };

      const handleCompositionEnd = (e: any) => {
        isComposingRef.current = false;
        const value = e.target.value;
        if (value) onSearchChange();
      };

      inputElement.addEventListener('compositionstart', handleCompositionStart);
      inputElement.addEventListener('compositionend', handleCompositionEnd);

      return () => {
        inputElement.removeEventListener('compositionstart', handleCompositionStart);
        inputElement.removeEventListener('compositionend', handleCompositionEnd);
      };
    }
  }, []);

  const placeholderMap = {
    uns: formatMessage('common.searchPlaceholderUns'),
    template: formatMessage('common.searchPlaceholderTem'),
    label: formatMessage('common.searchPlaceholderLabel'),
  };

  return (
    <ProSearch
      ref={selectRef}
      closeButtonLabelText={formatMessage('common.clearSearchInput')}
      placeholder={placeholderMap[treeType]}
      size="sm"
      value={searchValue ?? ''}
      onChange={(e) => {
        const val = e.target.value || '';
        setSearchValue(val);
        if (isComposingRef.current) return;
        onSearchChange();
      }}
      style={{ borderRadius: '3px', flex: 1 }}
      onKeyDown={(e) => {
        if (e.key === 'Enter') {
          loadData({ reset: true });
        }
      }}
      title={searchValue || placeholderMap[treeType]}
    />
  );
};

// uns搜索类型
const UnsTypeSearch = () => {
  const formatMessage = useTranslate();
  const { treeType, searchType, setSearchType, loadData } = useTreeStore((state) => ({
    searchType: state.searchType,
    treeType: state.treeType,
    setSearchType: state.setSearchType,
    loadData: state.loadData,
    setSelectedNode: state.setSelectedNode,
  }));

  const popoverContent = (
    <Radio.Group
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: 8,
      }}
      value={searchType}
      onChange={(e) => {
        setSearchType(e.target.value);
        loadData({ reset: true, clearSelect: true });
      }}
      options={[
        { value: 1, label: formatMessage('fieldTypeSTRING'), title: formatMessage('fieldTypeSTRING') },
        { value: 3, label: formatMessage('uns.hasTemplate'), title: formatMessage('uns.hasTemplate') },
        { value: 2, label: formatMessage('uns.hasLabel'), title: formatMessage('uns.hasLabel') },
      ]}
    />
  );
  return (
    treeType === 'uns' && (
      <Popover placement="bottomLeft" title="" content={popoverContent} trigger="hover">
        <Button
          icon={<Filter />}
          style={{ flexShrink: 0, background: 'var(--supos-switchwrap-bg-color)' }}
          color="default"
          variant="filled"
        />
      </Popover>
    )
  );
};

const TreeHeader = () => {
  return (
    <div>
      <TreeTab />
      <Flex gap={8} align="center">
        <UnsTypeSearch />
        <Search />
        <Operation />
      </Flex>

      <Divider
        style={{
          background: '#c6c6c6',
          margin: '16px 0 10px',
        }}
      />
    </div>
  );
};

// uns树的icon展示
const TreeNodeIcon = memo(({ dataNode }: { dataNode: UnsTreeNode }) => {
  const { mountStatus } = useUnsContext();
  const { expandedKeys } = useTreeStore((state) => ({
    expandedKeys: state.expandedKeys,
  }));
  if (dataNode.type === 0) {
    return (
      <Flex align="center">
        <div style={{ width: 10, display: 'flex', alignItems: 'center' }}>
          {dataNode.alias && mountStatus[dataNode.alias] && <StatusDot status={mountStatus[dataNode.alias]} />}
        </div>
        {expandedKeys.includes(dataNode.key) && !dataNode.hasChildren ? (
          <FolderOpen style={{ flexShrink: 0, marginRight: '5px' }} />
        ) : (
          <Folder style={{ flexShrink: 0, marginRight: '5px' }} />
        )}
      </Flex>
    );
  } else if (dataNode.type === 2) {
    return (
      <Flex align="center">
        <div style={{ width: 10 }} />
        <Document style={{ flexShrink: 0, marginRight: '5px' }} />
      </Flex>
    );
  }
  return null;
});

const findNodeWithParent = (
  tree: any[],
  id: string,
  parent: any | null = null
): { node: any; parent: any | null; index: number } | null => {
  for (let i = 0; i < tree.length; i++) {
    const node = tree[i];
    if (node.id === id) return { node, parent, index: i };
    if (node.children) {
      const result = findNodeWithParent(node.children, id, node);
      if (result) return result;
    }
  }
  return null;
};

const TopTreeCom = ({
  header,
  treeNodeExtra,
}: {
  header: ProTreeProps['header'];
  treeNodeExtra?: ProTreeProps['treeNodeExtra'];
}) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  // 创建一个 ref 来引用 tree 元素
  const treeRef = useRef<any>(null);

  const {
    lazyTree,
    loadData,
    treeData,
    expandedKeys,
    setExpandedKeys,
    loadedKeys,
    setLoadedKeys,
    nodePaginationState,
    loadingKeys,
    loading,
    treeType,
    selectedNode,
    setSelectedNode,
    setTreeMap,
    breadcrumbList,
    searchValue,
    operationFns,
    setPasteNode,
    pasteNode,
    onRefresh,
    setCurrentTreeMapType,
    setScrollTreeNode,
    setTreeType,
    handleExpandNode,
    setTreeData,
    setLoading,
  } = useTreeStore((state) => ({
    lazyTree: state.lazyTree,
    loadData: state.loadData,
    treeData: state.treeData,
    expandedKeys: state.expandedKeys,
    setExpandedKeys: state.setExpandedKeys,
    loadedKeys: state.loadedKeys,
    setLoadedKeys: state.setLoadedKeys,
    nodePaginationState: state.nodePaginationState,
    setNodePaginationState: state.setNodePaginationState,
    loadingKeys: state.loadingKeys,
    loading: state.loading,
    treeType: state.treeType,
    selectedNode: state.selectedNode,
    setSelectedNode: state.setSelectedNode,
    breadcrumbList: state.breadcrumbList,
    searchValue: state.searchValue,
    operationFns: state.operationFns,
    setPasteNode: state.setPasteNode,
    pasteNode: state.pasteNode,
    onRefresh: state.onRefresh,
    setTreeMap: state.setTreeMap,
    setCurrentTreeMapType: state.setCurrentTreeMapType,
    setScrollTreeNode: state.setScrollTreeNode,
    setTreeType: state.setTreeType,
    handleExpandNode: state.handleExpandNode,
    setTreeData: state.setTreeData,
    setLoading: state.setLoading,
  }));

  const { mountStatus } = useUnsContext();

  //滚动到目标树节点
  const scrollTreeNode = useCallback((id: Key) => {
    setTimeout(() => {
      if (treeRef.current) treeRef.current.scrollTo?.({ key: id, align: 'top' });
    }, 500);
  }, []);

  useEffect(() => {
    setScrollTreeNode(scrollTreeNode);
  }, [scrollTreeNode]);

  const onLoadData = async (node: any) => {
    const _node = { ...node };
    return loadData({
      key: _node.key,
      parentInfo: _node,
    });
  };
  const toTargetNode = (type: TTreeType, node: any) => {
    setTreeMap(false);
    setTreeType(type);
    scrollTreeNode(node.id);
    loadData({ queryType: 'viewTemplate', newNodeKey: node.id }, (data) => {
      setSelectedNode(data?.find((f) => f.id === node.id));
      scrollTreeNode(node.id);
      setCurrentTreeMapType('all');
    });
  };
  const { ViewLabelModal, setLabelOpen } = useViewLabelModal({ toTargetNode: toTargetNode });

  const handleRenderLoadMoreNode = (moreNodeData: any) => {
    const { parentKey: nodeKey, currentPage } = moreNodeData;
    if (loadingKeys.has(nodeKey)) {
      console.log(`正在loading ${nodeKey}`);
      return;
    }
    const state = nodePaginationState[nodeKey];
    if (state && state.hasMore && !state.isLoading) {
      // 处理根节点和子节点的加载更多
      if (currentPage === state.currentPage) {
        // 以防重复请求2次
        loadData({
          key: nodeKey,
          page: state.currentPage + 1,
          parentInfo: moreNodeData?.parentInfo,
        });
      }
    }
  };
  const isUns = treeType === 'uns';
  const isShow = usePageIsShow();
  const pasteHandle = (source: any, target: any) => {
    if (source) {
      setLoading(true);
      const targetId = (target?.type === 0 ? target?.id : target?.parentId) || '';
      pasteUns({
        sourceId: source?.id,
        targetId,
      })
        .then((data) => {
          loadData({
            queryType: source?.type === 0 ? 'addFolder' : 'addFile',
            key: targetId ? targetId : ROOT_NODE_ID,
            newNodeKey: data,
            reset: !targetId,
          });
        })
        .catch(() => {
          setLoading(false);
        });
    } else {
      message.warning(formatMessage('uns.copyTip'));
    }
  };
  return (
    <>
      {ViewLabelModal}
      <ProTree
        onDndDragStart={(info) => {
          const { active } = info;
          if (active.type == 0 && expandedKeys.includes(active?.id)) {
            // 关闭文件夹
            setExpandedKeys((expandedKeys) => {
              return expandedKeys.filter((k) => active?.id !== k);
            });
          }
        }}
        onDndDragEnd={(info) => {
          const { active, over, isInset } = info;
          if (active?.id && over?.id && active.id !== over.id) {
            let _isInset = isInset;
            if (over?.type === 2) {
              _isInset = false;
            }
            // active是我移动的节点，over是我放置的节点位置，_isInset是我放置在over这个节点里面还是over节点的下面
            console.log(active, over, _isInset);
            setTreeData((draft) => {
              const activeInfo = findNodeWithParent(draft, active.id);
              const overInfo = findNodeWithParent(draft, over.id);
              if (activeInfo && overInfo) {
                // 移除原节点
                if (activeInfo.parent) {
                  activeInfo.parent.children.splice(activeInfo.index, 1);
                } else {
                  draft.splice(activeInfo.index, 1);
                }

                // 处理插入位置
                const isRootOperation = !overInfo.parent;
                const targetArray = isRootOperation ? draft : overInfo.parent?.children || [];

                let insertIndex = _isInset ? 0 : overInfo.index + 1;

                // 边界检查
                insertIndex = Math.min(insertIndex, targetArray.length);

                if (_isInset) {
                  overInfo.node.children = [activeInfo.node, ...(overInfo.node.children || [])];
                } else {
                  targetArray.splice(insertIndex, 0, activeInfo.node);
                }
              }
            });
          }
        }}
        // dndDraggable={isUns}
        isShow={isShow}
        ref={treeRef}
        rightClickMenuItems={
          isUns
            ? ({ node }) => {
                if (!node) {
                  // 空白点击
                  return [
                    {
                      auth: [ButtonPermission['uns.folderPaste'], ButtonPermission['uns.filePaste']],
                      key: 'paste',
                      label: formatMessage('common.paste'),
                      onClick: () => {
                        pasteHandle(pasteNode, null);
                      },
                    },
                    {
                      auth: [ButtonPermission['uns.folderPaste'], ButtonPermission['uns.filePaste']],
                      key: 'pasteAndEdit',
                      label: formatMessage('common.pasteAndEdit'),
                      onClick: () => {
                        if (pasteNode) {
                          //纯前端方案
                          operationFns?.setOptionOpen?.('paste', null, pasteNode);
                        } else {
                          message.warning(formatMessage('uns.copyTip'));
                        }
                      },
                    },
                    {
                      auth: ButtonPermission['uns.folderAdd'],
                      key: 'addFolder',
                      label: formatMessage('common.createNewFolder'),
                      onClick: () => {
                        operationFns?.setOptionOpen?.('addFolder');
                      },
                    },
                    {
                      auth: ButtonPermission['uns.fileAdd'],
                      key: 'addFile',
                      label: formatMessage('common.createNewFile'),
                      onClick: () => {
                        operationFns?.setOptionOpen?.('addFile');
                      },
                    },
                  ];
                }
                const _node = { ...node };
                const baseItems = ['viewTemplate', 'copy', 'paste', 'pasteAndEdit', 'addFolder', 'addFile', 'delete'];
                const folderItems = lazyTree
                  ? ['refresh', ...baseItems, 'collapseFolder']
                  : [...baseItems, 'expandFolder', 'collapseFolder'];

                const mapItem = _node.type === 0 ? folderItems : ['viewLabels', ...baseItems];
                const isMountTarget = !!node.mount && mountStatus[node.alias];
                const isMountFile = !!node.mount;

                return filterPermissionToList<ItemType>(
                  [
                    {
                      key: 'refresh',
                      label: formatMessage('common.refresh'),
                      onClick: () => {
                        onRefresh(_node);
                      },
                    },
                    {
                      key: 'viewLabels',
                      label: formatMessage('common.viewLabels'),
                      onClick: async () => {
                        const getInfo = _node.type === 2 ? getInstanceInfo : getModelInfo;
                        const detail: any = await getInfo({ id: _node.id });
                        if (detail?.labelList?.length > 0) {
                          setLabelOpen(detail.labelList);
                        } else {
                          message.warning(formatMessage('uns.noLabel'));
                        }
                      },
                    },
                    {
                      key: 'viewTemplate',
                      label: formatMessage('common.viewTemplate'),
                      onClick: async () => {
                        const getInfo = _node.type === 2 ? getInstanceInfo : getModelInfo;
                        const detail: any = await getInfo({ id: _node.id });
                        if (detail.modelId) {
                          toTargetNode('template', { key: detail.modelId, type: 1, id: detail.modelId });
                        } else {
                          message.warning(formatMessage('uns.noTemplate'));
                        }
                      },
                    },
                    {
                      auth: _node.type === 0 ? ButtonPermission['uns.folderCopy'] : ButtonPermission['uns.fileCopy'],
                      key: 'copy',
                      label: formatMessage('common.copy'),
                      onClick: () => {
                        setPasteNode(_node);
                        message.success(formatMessage('common.copySuccess'));
                      },
                    },
                    {
                      auth: _node.type === 0 ? ButtonPermission['uns.folderPaste'] : ButtonPermission['uns.filePaste'],
                      key: 'paste',
                      label: formatMessage('common.paste'),
                      onClick: () => {
                        pasteHandle(pasteNode, _node);
                      },
                      disabled: isMountFile,
                    },
                    {
                      auth: _node.type === 0 ? ButtonPermission['uns.folderPaste'] : ButtonPermission['uns.filePaste'],
                      key: 'pasteAndEdit',
                      label: formatMessage('common.pasteAndEdit'),
                      onClick: () => {
                        if (pasteNode) {
                          //纯前端方案
                          operationFns?.setOptionOpen?.('paste', _node, pasteNode);
                        } else {
                          message.warning(formatMessage('uns.copyTip'));
                        }
                      },
                      disabled: isMountFile,
                    },
                    {
                      auth: ButtonPermission['uns.folderAdd'],
                      key: 'addFolder',
                      label: formatMessage('common.createNewFolder'),
                      onClick: () => {
                        operationFns?.setOptionOpen?.('addFolder', _node);
                      },
                      disabled: isMountFile,
                    },
                    {
                      auth: ButtonPermission['uns.fileAdd'],
                      key: 'addFile',
                      label: formatMessage('common.createNewFile'),
                      onClick: () => {
                        operationFns?.setOptionOpen?.('addFile', _node);
                      },
                      disabled: isMountFile,
                    },
                    {
                      key: 'expandFolder',
                      label: formatMessage('common.expandFolder'),
                      onClick: () => {
                        handleExpandNode(true, _node);
                      },
                    },
                    {
                      key: 'collapseFolder',
                      label: formatMessage('common.collapseFolder'),
                      onClick: () => {
                        handleExpandNode(false, _node);
                      },
                    },
                    {
                      type: 'divider',
                    },
                    {
                      auth:
                        _node.type === 0 ? ButtonPermission['uns.folderDelete'] : ButtonPermission['uns.fileDelete'],
                      key: 'delete',
                      label: formatMessage('common.delete'),
                      onClick: () => {
                        operationFns?.setDeleteOpen?.(_node);
                      },
                      extra: (
                        <div style={{ display: 'flex' }}>
                          <TrashCan />
                        </div>
                      ),
                      disabled: isMountFile && !isMountTarget,
                    },
                  ]?.filter((f) => !f.key || mapItem.includes(f.key)) as any
                );
              }
            : undefined
        }
        matchHighlightValue={searchValue}
        showSwitcherIcon={isUns}
        selectedKeys={selectedNode ? [selectedNode.key] : []}
        onSelect={(_, { node, selected }) => {
          const selectedNode = selected ? { ...node } : undefined;
          setSelectedNode(selectedNode);
          setTreeMap(false);
          setCurrentTreeMapType('all');
        }}
        loading={loading}
        treeData={treeData}
        loadData={onLoadData}
        loadMoreData={handleRenderLoadMoreNode}
        loadedKeys={loadedKeys}
        // onLoad={(newLoadedKeys) => setLoadedKeys(newLoadedKeys)}
        expandedKeys={expandedKeys}
        onExpand={(expandedKeys) => {
          setExpandedKeys(expandedKeys);
          setLoadedKeys(expandedKeys);
        }}
        lazy={treeType === 'template' ? true : lazyTree}
        header={header}
        wrapperStyle={{ padding: '0 14px' }}
        height={0}
        treeNodeIcon={isUns ? (dataNode) => <TreeNodeIcon dataNode={dataNode} /> : undefined}
        filterTreeNode={(node) => {
          // 高亮字段
          return (
            breadcrumbList
              ?.slice(0, -1)
              .map((e) => e.key)
              .includes(node.key) ?? false
          );
        }}
        treeNodeCount={(dataNode) => {
          return (
            dataNode.type === 0 && (
              <span style={{ color: 'var(--supos-text-color)', fontSize: '12px', opacity: 0.5 }}>
                ({dataNode.countChildren})
              </span>
            )
          );
        }}
        treeNodeExtra={treeNodeExtra}
        overlayChildren={(dataNode: any) => {
          return (
            <Flex
              justify="flex-start"
              align="center"
              className="overlay-dom"
              style={{
                color: 'var(--supos-text-color)',
              }}
            >
              {isUns ? <TreeNodeIcon dataNode={dataNode} /> : undefined}
              <span style={{ fontWeight: 'bold', fontSize: 14 }}>{dataNode.title}</span>
              {dataNode.type === 0 && (
                <span style={{ fontSize: '12px', opacity: 0.5 }}>({dataNode.countChildren})</span>
              )}
            </Flex>
          );
        }}
        drapOverChanges={(info) => {
          const { node, classNames, isInset } = info;
          if (node?.type === 2) {
            return classNames.out;
          } else if (node?.type === 0) {
            return isInset ? classNames.in : classNames.out;
          } else {
            return '';
          }
        }}
      />
    </>
  );
};

const UnsTree = ({ treeNodeExtra }: { treeNodeExtra?: ProTreeProps['treeNodeExtra'] }) => {
  return <TopTreeCom treeNodeExtra={treeNodeExtra} header={<TreeHeader />} />;
};
export default UnsTree;
