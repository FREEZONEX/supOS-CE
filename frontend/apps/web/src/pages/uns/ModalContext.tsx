/**
 * @description modal相关的统一放这里,该组件渲染不影响父级，进而不影响父级下面的所有子集
 */
import { type FC, useEffect } from 'react';
import { useCreateModal, useDeleteModal, useLabelModal, useTemplateModal } from '@/pages/uns/components';
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
  const { loadData, setTreeMap, scrollTreeNode, setOperationFns, setSelectedNode } = getTreeStoreSnapshot(
    stateRef,
    (state) => ({
      loadData: state.loadData,
      setTreeMap: state.setTreeMap,
      scrollTreeNode: state.scrollTreeNode,
      setOperationFns: state.setOperationFns,
      setSelectedNode: state.setSelectedNode,
    })
  );

  const { OptionModal, setOptionOpen } = useCreateModal({
    successCallBack: loadData,
    addNamespaceForAi,
    setAddNamespaceForAi,
    changeCurrentPath,
    setTreeMap,
  });

  const { TemplateModal, openTemplateModal } = useTemplateModal({
    successCallBack: loadData,
    changeCurrentPath,
    scrollTreeNode,
    setTreeMap,
  });

  const { LabelModal, setLabelOpen } = useLabelModal({
    successCallBack: loadData,
    changeCurrentPath,
    scrollTreeNode,
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
  //
  useEffect(() => {
    setOperationFns({
      // 创建模板弹框
      openTemplateModal,
    });
  }, [openTemplateModal]);

  useEffect(() => {
    setOperationFns({
      // uns删除弹框
      setDeleteOpen,
    });
  }, [setDeleteOpen]);

  useEffect(() => {
    setOperationFns({
      // 创建标签弹框
      setLabelOpen,
    });
  }, [setLabelOpen]);

  return (
    <>
      {OptionModal}
      {TemplateModal}
      {LabelModal}
      {DeleteModal}
    </>
  );
};

export default ModalContext;
