import { ModelDetail, RealTimeData, TopicDetail, EmptyDetail } from '@/pages/uns/components';
import { getTreeStoreSnapshot, useTreeStore, useTreeStoreRef } from './store/treeStore';
import type { FC } from 'react';
import type { UnsTreeNode } from '@/pages/uns/types';

interface DetailDomProps {
  handleDelete: (item: UnsTreeNode) => void;
  currentUnusedTopicNode?: any;
}

const DetailDom: FC<DetailDomProps> = ({ handleDelete, currentUnusedTopicNode }) => {
  const { currentTreeMapType, selectedNode } = useTreeStore((state) => ({
    currentTreeMapType: state.currentTreeMapType,
    selectedNode: state.selectedNode,
  }));

  const stateRef = useTreeStoreRef();
  const { loadData } = getTreeStoreSnapshot(stateRef, (state) => ({
    loadData: state.loadData,
  }));

  const getDetailDom = (selectedNode?: UnsTreeNode) => {
    if (!selectedNode) return <EmptyDetail />;
    switch (selectedNode.pathType) {
      case 0:
        return <ModelDetail currentNode={selectedNode} initTreeData={loadData} />;
      case 2:
        return <TopicDetail currentNode={selectedNode} handleDelete={handleDelete} initTreeData={loadData} />;
      default:
        return <EmptyDetail />;
    }
  };
  return currentTreeMapType === 'all' ? (
    getDetailDom(selectedNode)
  ) : (
    <RealTimeData currentNode={currentUnusedTopicNode} />
  );
};

export default DetailDom;
