import { type FC, useCallback, useEffect, useState } from 'react';
import { useDeepCompareEffect } from 'ahooks';
import { getTreeStoreSnapshot, TreeStoreProvider, useTreeStore, useTreeStoreRef } from './store/treeStore';
import { findNodeInfoById, getParentNodes, handlerTreeData } from './store/utils';
import ModalContext from './ModalContext';
import TopDom from './TopDom';
import DetailDom from './DetailDom';
import LeftDom from './LeftDom';
import type { UnsTreeNode } from './types';
// import { useLocation } from 'react-router';
// import { guideSteps } from './guide-steps';

import './index.scss';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { setAiResult, useAiStore } from '@/stores/ai-store.ts';
import { UnsContextProvider } from './UnsContext';
import type { PageProps } from '@/common-types.ts';
import { useLocation } from 'react-router';
import { getModelInfo, getTreeData } from '@/apis/core-api/uns';

const initNode = {
  key: '',
  id: '',
  pathType: null,
};

const normalizeNodeKey = (value: unknown) => {
  const key = String(value ?? '').trim();
  return key && key !== 'undefined' && key !== 'null' ? key : '';
};

const findFirstMatchingNode = (treeData: UnsTreeNode[], rawValue: string) => {
  const target = normalizeNodeKey(rawValue).toLowerCase();
  if (!target) return undefined;
  let fallback: UnsTreeNode | undefined;

  const walk = (nodes: UnsTreeNode[]): UnsTreeNode | undefined => {
    for (const node of nodes) {
      const candidates = [node.id, node.key, node.path, node.alias, node.name, node.title, node.pathName]
        .map((value) => normalizeNodeKey(value).toLowerCase())
        .filter(Boolean);
      if (candidates.some((value) => value === target)) {
        return node;
      }
      if (!fallback && candidates.some((value) => value.includes(target))) {
        fallback = node;
      }
      const child = walk(node.children || []);
      if (child) {
        return child;
      }
    }
    return undefined;
  };

  return walk(treeData) || fallback;
};


const Module: FC = () => {
  const location = useLocation();
  // 默认落在 Agent（画布）——对齐 saas unsTree store 的 `activeTab: 'canvas'` 初始值；
  // 首次选中节点后由下方 effect 自动切到 detail（同 saas onSelect→uns-info）。
  const { treeType, selectedNode, operationFns } = useTreeStore((state) => ({
    treeType: state.treeType,
    selectedNode: state.selectedNode,
    operationFns: state.operationFns,
  }));
  const stateRef = useTreeStoreRef();

  const { setSelectedNode, setCurrentTreeMapType, setTreeMap } = getTreeStoreSnapshot(stateRef, (state) => ({
    setSelectedNode: state.setSelectedNode,
    setTreeMap: state.setTreeMap,
    setCurrentTreeMapType: state.setCurrentTreeMapType,
    setPasteNode: state.setPasteNode,
  }));

  const [addNamespaceForAi, setAddNamespaceForAi] = useState<any>(null);
  const aiResult = useAiStore((state) => state.aiResult);

  const [currentUnusedTopicNode, setCurrentUnusedTopicNode] = useState<UnsTreeNode>(initNode); // 当前unusedTopic节点
  const [unusedTopicBreadcrumbList, setUnusedTopicBreadcrumbList] = useState<UnsTreeNode[]>([]); //当前文件路径Array


  const resetToNamespaceGuide = useCallback(() => {
    setTreeMap(false);
    setSelectedNode(undefined);
    setCurrentTreeMapType('all');
    setCurrentUnusedTopicNode(initNode);
  }, [setCurrentTreeMapType, setSelectedNode, setTreeMap]);

  useEffect(() => {
    resetToNamespaceGuide();
  }, [location.state?.resetUnsLanding, resetToNamespaceGuide]);

  const focusVisibleUnsNode = useCallback(
    (node: UnsTreeNode, treeData?: UnsTreeNode[]) => {
      const store = stateRef.getState();
      const sourceTree = treeData || store.treeData;
      const nodeKey = normalizeNodeKey(node.id ?? node.key);
      if (!nodeKey) return;

      const visibleNode = findNodeInfoById(sourceTree, nodeKey) || node;
      const parents = getParentNodes(sourceTree, visibleNode.key);
      const selected = parents[parents.length - 1] || visibleNode;
      const parentKeys = parents
        .slice(0, -1)
        .map((item) => normalizeNodeKey(item.id ?? item.key))
        .filter(Boolean);

      store.setTreeMap(false);
      store.setCurrentTreeMapType('all');
      store.setSearchValue('');
      setCurrentUnusedTopicNode(initNode);
      if (parentKeys.length) {
        store.setExpandedKeys((expandedKeys) => Array.from(new Set([...expandedKeys, ...parentKeys])));
      }
      store.setSelectedNode(selected, parents.length === 0);
      window.setTimeout(() => {
        stateRef.getState().scrollTreeNode?.(selected.id ?? selected.key);
      }, 0);
    },
    [stateRef]
  );

  const focusUnsNode = useCallback(
    async (namespaceId: string) => {
      const targetId = normalizeNodeKey(namespaceId);
      if (!targetId) return;

      const store = stateRef.getState();
      const currentNode = findNodeInfoById(store.treeData, targetId);
      if (currentNode) {
        focusVisibleUnsNode(currentNode, store.treeData);
        return;
      }

      store.setLoading(true);
      try {
        const detail = await getModelInfo({ id: targetId });
        const searchKey = normalizeNodeKey(
          detail?.path || detail?.alias || detail?.namespace || detail?.name || targetId
        );
        const searchedTree = searchKey
          ? handlerTreeData(await getTreeData({ type: 1, key: searchKey, keyword: searchKey }))
          : [];
        const targetNode = findNodeInfoById(searchedTree, targetId);

        if (targetNode) {
          store.setTreeData(searchedTree);
          focusVisibleUnsNode(targetNode, searchedTree);
          return;
        }

        if (detail?.id) {
          focusVisibleUnsNode(detail);
        }
      } catch (error) {
        console.error('Failed to focus UNS node:', error);
      } finally {
        stateRef.getState().setLoading(false);
      }
    },
    [focusVisibleUnsNode, stateRef]
  );

  const focusUnsSearchKey = useCallback(
    async (namespaceKey: string) => {
      const searchKey = normalizeNodeKey(namespaceKey);
      if (!searchKey) return;

      const store = stateRef.getState();
      store.setLoading(true);
      try {
        const searchedTree = handlerTreeData(await getTreeData({ type: 1, key: searchKey, keyword: searchKey }));
        const targetNode = findFirstMatchingNode(searchedTree, searchKey);
        if (targetNode) {
          store.setTreeData(searchedTree);
          focusVisibleUnsNode(targetNode, searchedTree);
        }
      } catch (error) {
        console.error('Failed to focus UNS node:', error);
      } finally {
        stateRef.getState().setLoading(false);
      }
    },
    [focusVisibleUnsNode, stateRef]
  );

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const namespaceId = params.get('namespaceId');
    const namespaceKey = params.get('namespaceKey');
    if (namespaceId) {
      void focusUnsNode(namespaceId);
    } else if (namespaceKey) {
      void focusUnsSearchKey(namespaceKey);
    }
  }, [focusUnsNode, focusUnsSearchKey, location.search]);

  useDeepCompareEffect(() => {
    if (aiResult?.uns) {
      setAddNamespaceForAi(aiResult?.uns);
      setAiResult('uns', undefined);
    }
  }, [aiResult?.uns]);

  // uns删除操作
  const handleDelete = (item: UnsTreeNode) => {
    switch (treeType) {
      case 'uns':
        operationFns?.setDeleteOpen(item as any);
        break;
      default:
        break;
    }
  };

  const changeCurrentPath = (node?: UnsTreeNode) => {
    setSelectedNode(node?.id === selectedNode?.id ? undefined : node);
    setCurrentUnusedTopicNode(initNode);
    setCurrentTreeMapType('all');
  };

  // const location = useLocation();
  // 新手导航步骤
  // useGuideSteps(guideSteps(), location?.state?.stepId);

  return (
    <ComLayout className="unsContainer">
      <LeftDom
        changeCurrentPath={changeCurrentPath}
        handleDelete={handleDelete}
        setCurrentUnusedTopicNode={setCurrentUnusedTopicNode}
        setUnusedTopicBreadcrumbList={setUnusedTopicBreadcrumbList}
      />
      <ComContent>
        <div className="chartWrap">
          <TopDom
            changeCurrentPath={changeCurrentPath}
            setCurrentUnusedTopicNode={setCurrentUnusedTopicNode}
            unusedTopicBreadcrumbList={unusedTopicBreadcrumbList}
            currentUnusedTopicNode={currentUnusedTopicNode}
          />
          <DetailDom handleDelete={handleDelete} currentUnusedTopicNode={currentUnusedTopicNode} />
        </div>
      </ComContent>
      <ModalContext
        addNamespaceForAi={addNamespaceForAi}
        setAddNamespaceForAi={setAddNamespaceForAi}
        changeCurrentPath={changeCurrentPath}
      />
    </ComLayout>
  );
};

const WrapperModule: FC<PageProps> = () => {
  return (
    <TreeStoreProvider>
      <UnsContextProvider>
        <Module />
      </UnsContextProvider>
    </TreeStoreProvider>
  );
};
export default WrapperModule;
