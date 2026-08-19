import { useTranslate } from '@/hooks';
import DetailDom from '@/pages/uns/DetailDom';
import LeftDom from '@/pages/uns/LeftDom';
import ModalContext from '@/pages/uns/ModalContext';
import { getTreeStoreSnapshot, TreeStoreProvider, useTreeStore, useTreeStoreRef } from '@/pages/uns/store/treeStore';
import type { UnsTreeNode } from '@/pages/uns/types';
import { useDeepCompareEffect } from 'ahooks';
import { type FC, useEffect, useState } from 'react';
// import { useLocation } from 'react-router';
// import { guideSteps } from './guide-steps';

import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import '@/pages/uns/index.scss';
import { UnsContextProvider } from '@/pages/uns/UnsContext';
import { setAiResult, useAiStore } from '@/stores/ai-store.ts';
import { PageTitleIcon } from '@/components/lucide-icon';
import styles from './index.module.scss';

const initNode = {
  key: '',
  id: '',
  pathType: null,
};

const Module: FC = () => {
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

  const formatMessage = useTranslate();
  const [addNamespaceForAi, setAddNamespaceForAi] = useState<any>(null);
  const aiResult = useAiStore((state) => state.aiResult);

  const [currentUnusedTopicNode, setCurrentUnusedTopicNode] = useState<UnsTreeNode>(initNode); // 当前unusedTopic节点

  useDeepCompareEffect(() => {
    if (aiResult?.uns) {
      setAddNamespaceForAi(aiResult?.uns);
      setAiResult('uns', undefined);
    }
  }, [aiResult?.uns]);

  // useEffect(() => {
  //   if (locationTreeMap) {
  //     // 点击logo过来的需要跳到overview页
  //     setTreeMap(true);
  //     changeCurrentPath();
  //     navigate(location.pathname, { replace: true, state: {} });
  //   }
  // }, [locationTreeMap]);

  useEffect(() => {
    setTreeMap(false);
  }, []);

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
        className={styles['left-dom']}
        changeCurrentPath={changeCurrentPath}
        handleDelete={handleDelete}
        setCurrentUnusedTopicNode={setCurrentUnusedTopicNode}
        customTitle={
          <div
            className="treemapTitle"
            style={{
              background: 'var(--ui-bg-color)',
            }}
          >
            <PageTitleIcon resourceKey="uns.page" />
            <span>{formatMessage('uns.treeList')}</span>
          </div>
        }
        customTreeTab={<p style={{ paddingTop: '16px' }}></p>}
        enableRightClick={false}
        enableRootAsParent
        disableDeleteRoot
      />
      <ComContent>
        <div className="chartWrap">
          {/* <TopDom
            changeCurrentPath={changeCurrentPath}
            setCurrentUnusedTopicNode={setCurrentUnusedTopicNode}
            currentUnusedTopicNode={currentUnusedTopicNode}
          /> */}
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

interface UnsProps {
  projectId: string;
}

const WrapperModule: FC<UnsProps> = ({ projectId }) => {
  return (
    <TreeStoreProvider
      initProps={{
        loadDataDefaultParams: {
          projectId,
        },
      }}
    >
      <UnsContextProvider>
        <Module />
      </UnsContextProvider>
    </TreeStoreProvider>
  );
};
export default WrapperModule;
