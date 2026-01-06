import { useState } from 'react';
import { App } from 'antd';
import { useDeepCompareEffect } from 'ahooks';
import { useTranslate } from '@/hooks';
import { deleteLabel, deleteTemplate } from '@/apis/inter-api/uns';
import { getTreeStoreSnapshot, TreeStoreProvider, useTreeStore, useTreeStoreRef } from './store/treeStore';
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

const initNode = {
  key: '',
  id: '',
  pathType: null,
};

const Module = () => {
  const { treeType, selectedNode, operationFns } = useTreeStore((state) => ({
    treeType: state.treeType,
    selectedNode: state.selectedNode,
    operationFns: state.operationFns,
  }));
  const stateRef = useTreeStoreRef();

  const { loadData, setSelectedNode, setCurrentTreeMapType } = getTreeStoreSnapshot(stateRef, (state) => ({
    loadData: state.loadData,
    setSelectedNode: state.setSelectedNode,
    setTreeMap: state.setTreeMap,
    setCurrentTreeMapType: state.setCurrentTreeMapType,
    setPasteNode: state.setPasteNode,
  }));

  const { modal, message } = App.useApp();
  const formatMessage = useTranslate();
  const [addNamespaceForAi, setAddNamespaceForAi] = useState<any>(null);
  const aiResult = useAiStore((state) => state.aiResult);

  const [currentUnusedTopicNode, setCurrentUnusedTopicNode] = useState<UnsTreeNode>(initNode); // 当前unusedTopic节点
  const [unusedTopicBreadcrumbList, setUnusedTopicBreadcrumbList] = useState<UnsTreeNode[]>([]); //当前文件路径Array

  useDeepCompareEffect(() => {
    if (aiResult?.uns) {
      setAddNamespaceForAi(aiResult?.uns);
      setAiResult('uns', undefined);
    }
  }, [aiResult?.uns]);

  // uns、template、label删除操作
  const handleDelete = (item: UnsTreeNode) => {
    const { id } = item;
    switch (treeType) {
      case 'uns':
        operationFns?.setDeleteOpen(item as any);
        break;
      case 'label':
        modal.confirm({
          content: formatMessage('uns.areYouSureToDeleteThisLabel'),
          onOk() {
            deleteLabel(id as string).then(() => {
              loadData({ clearSelect: id === selectedNode?.id });
              message.success(formatMessage('common.deleteSuccessfully'));
            });
          },
          okButtonProps: {
            title: formatMessage('common.confirm'),
          },
          cancelButtonProps: {
            title: formatMessage('common.cancel'),
          },
        });
        break;
      case 'template':
        modal.confirm({
          content: formatMessage('common.deleteTemplateConfirm'),
          onOk() {
            deleteTemplate(id as string).then(() => {
              loadData(
                id === selectedNode?.id
                  ? { clearSelect: id === selectedNode?.id }
                  : {
                      queryType: 'deleteTemplate',
                      newNodeKey: selectedNode?.id,
                    }
              );
              message.success(formatMessage('common.deleteSuccessfully'));
            });
          },
          okButtonProps: {
            title: formatMessage('common.confirm'),
          },
          cancelButtonProps: {
            title: formatMessage('common.cancel'),
          },
        });
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
  const changeCurrentUnusedTopicPath = (node?: UnsTreeNode) => {
    setCurrentUnusedTopicNode(node?.id === currentUnusedTopicNode.id ? initNode : node || initNode);
    setSelectedNode();
    if (node?.id) {
      setCurrentTreeMapType('unusedTopic');
    }
  };

  // const location = useLocation();
  // 新手导航步骤
  // useGuideSteps(guideSteps(), location?.state?.stepId);

  return (
    <ComLayout className="unsContainer">
      <LeftDom
        changeCurrentPath={changeCurrentPath}
        handleDelete={handleDelete}
        currentUnusedTopicNode={currentUnusedTopicNode}
        setCurrentUnusedTopicNode={setCurrentUnusedTopicNode}
        unusedTopicBreadcrumbList={unusedTopicBreadcrumbList}
        changeCurrentUnusedTopicPath={changeCurrentUnusedTopicPath}
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

const WrapperModule = () => {
  return (
    <TreeStoreProvider>
      <UnsContextProvider>
        <Module />
      </UnsContextProvider>
    </TreeStoreProvider>
  );
};
export default WrapperModule;
