import { Copy, Subtract } from '@carbon/icons-react';
import type { FC } from 'react';
import { useTreeStore } from '../../store/treeStore';
import { useTranslate } from '@/hooks';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthWrapper } from '@/components/auth';

const deletePerMap: { [key: string]: string } = {
  uns0: ButtonPermission['uns.folderDelete'],
  uns2: ButtonPermission['uns.fileDelete'],
  template: ButtonPermission['uns.templateDelete'],
  label: ButtonPermission['uns.labelDelete'],
};

const copyPerMap: { [key: string]: string } = {
  uns0: ButtonPermission['uns.folderCopy'],
  uns1: ButtonPermission['uns.fileCopy'],
  template: ButtonPermission['uns.templateCopy'],
};

const TreeNodeExtra: FC<{
  handleDelete: () => void;
  handleCopy: () => void;
  type?: number;
  deleteDisabled?: boolean;
}> = ({ handleDelete, handleCopy, type, deleteDisabled }) => {
  const formatMessage = useTranslate();
  const { treeType } = useTreeStore((state) => ({
    treeType: state.treeType,
  }));

  return (
    <>
      {treeType !== 'label' && (
        <AuthWrapper auth={copyPerMap[treeType === 'uns' ? treeType + type : treeType]}>
          <span title={formatMessage('common.copy')} style={{ lineHeight: 1 }}>
            <Copy
              onClick={(e) => {
                e.stopPropagation();
                handleCopy?.();
              }}
              style={{ cursor: 'pointer' }}
            />
          </span>
        </AuthWrapper>
      )}
      <AuthWrapper auth={deletePerMap[treeType === 'uns' ? treeType + type : treeType]}>
        <span
          title={formatMessage('common.delete')}
          style={{
            lineHeight: 1,
            cursor: deleteDisabled ? 'not-allowed' : 'pointer',
            color: deleteDisabled ? 'var(--supos-select-d-color)' : 'var(--supos-text-color)',
          }}
        >
          <Subtract
            style={{ cursor: 'pointer' }}
            onClick={(e) => {
              e?.stopPropagation();
              if (deleteDisabled) return;
              handleDelete?.();
            }}
          />
        </span>
      </AuthWrapper>
    </>
  );
};

export default TreeNodeExtra;
