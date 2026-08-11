import { ButtonPermission } from '@/common-types/button-permission';
import { AuthWrapper } from '@/components/auth';
import { Copy, TrashCan } from '@/components/lucide-icon/carbon';
import { treeIconProps } from '@/components/lucide-icon/icon-props';
import { useTranslate } from '@/hooks';
import type { UnsTreeNode } from '@/pages/uns/types';
import type { FC } from 'react';
import { getDefaultDataRootNode, hasMountedTopicInSubtree, normalizeNodeId } from '.';
import { useTreeStore } from '../../store/treeStore';

const deletePerMap: { [key: string]: string } = {
  uns0: ButtonPermission['uns.folderDelete'],
  uns2: ButtonPermission['uns.fileDelete'],
};

const copyPerMap: { [key: string]: string } = {
  uns0: ButtonPermission['uns.folderCopy'],
  uns2: ButtonPermission['uns.fileCopy'],
};

const TreeNodeExtra: FC<{
  handleDelete: () => void;
  handleCopy: () => void;
  type?: number;
  node?: UnsTreeNode;
  disableDeleteRoot?: boolean;
}> = ({ handleDelete, handleCopy, type, node, disableDeleteRoot }) => {
  const formatMessage = useTranslate();
  const { treeType, treeData } = useTreeStore((state) => ({
    treeType: state.treeType,
    treeData: state.treeData,
  }));

  const currentDataRootNode = getDefaultDataRootNode(treeData || []);
  const currentDataRootNodeId = normalizeNodeId(currentDataRootNode?.id);
  const currentNodeId = normalizeNodeId(node?.id);
  const isMountedTopic = Boolean(node?.mount && type === 2);
  const deleteDisabled =
    treeType === 'uns' &&
    (hasMountedTopicInSubtree(node) ||
      (!!disableDeleteRoot && !!currentDataRootNodeId && currentNodeId === currentDataRootNodeId));

  if (isMountedTopic) {
    return null;
  }

  return (
    <>
      <AuthWrapper auth={copyPerMap[treeType + type]}>
        <span title={formatMessage('common.copy')} style={{ lineHeight: 1 }}>
          <Copy
            {...treeIconProps}
            onClick={(e: any) => {
              e.stopPropagation();
              handleCopy?.();
            }}
            style={{ cursor: 'pointer' }}
          />
        </span>
      </AuthWrapper>
      <AuthWrapper auth={deletePerMap[treeType + type]}>
        <span
          title={formatMessage('common.delete')}
          style={{
            lineHeight: 1,
            cursor: deleteDisabled ? 'not-allowed' : 'pointer',
            opacity: deleteDisabled ? 0.45 : 1,
          }}
        >
          <TrashCan
            {...treeIconProps}
            onClick={(e: any) => {
              e?.stopPropagation();
              if (deleteDisabled) {
                return;
              }
              handleDelete?.();
            }}
          />
        </span>
      </AuthWrapper>
    </>
  );
};

export default TreeNodeExtra;
