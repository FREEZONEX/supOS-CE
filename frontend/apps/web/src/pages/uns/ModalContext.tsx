/**
 * @description modal相关的统一放这里,该组件渲染不影响父级，进而不影响父级下面的所有子集
 */
import { type FC, useEffect } from 'react';
import { useCreateModal, useDeleteModal } from '@/pages/uns/components';
import { getTreeStoreSnapshot, useTreeStore, useTreeStoreRef } from './store/treeStore';

interface ModalContextProps {
  addNamespaceForAi: any;
  setAddNamespaceForAi: any;
  changeCurrentPath: any;
}
const ModalContext: FC<ModalContextProps> = ({ addNamespaceForAi, setAddNamespaceForAi, changeCurrentPath }) => {
  const { selectedNode, lazyTree } = useTreeStore((state) => ({
    selectedNode: state.selectedNode,
    lazyTree: state.lazyTree,
  }));
  const stateRef = useTreeStoreRef();
  const { loadData, setTreeMap, setOperationFns, setSelectedNode } = getTreeStoreSnapshot(stateRef, (state) => ({
    loadData: state.loadData,
    setTreeMap: state.setTreeMap,
    setOperationFns: state.setOperationFns,
    setSelectedNode: state.setSelectedNode,
  }));

  const { OptionModal, setOptionOpen } = useCreateModal({
    successCallBack: loadData,
    addNamespaceForAi,
    setAddNamespaceForAi,
    changeCurrentPath,
    setTreeMap,
  });

  const { DeleteModal, setDeleteOpen } = useDeleteModal({
    successCallBack: loadData,
    currentNode: selectedNode,
    setSelectedNode,
    lazyTree,
  });

  useEffect(() => {
    setOperationFns({
      // uns创建弹框
      setOptionOpen,
    });
  }, [setOptionOpen]);
  useEffect(() => {
    setOperationFns({
      // uns删除弹框
      setDeleteOpen,
    });
  }, [setDeleteOpen]);

  return (
    <>
      {OptionModal}
      {DeleteModal}
    </>
  );
};

export default ModalContext;
