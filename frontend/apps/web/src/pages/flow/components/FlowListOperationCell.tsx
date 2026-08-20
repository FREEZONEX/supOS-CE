import { AuthWrapper } from '@/components/auth';
import { Edit, TrashCan } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import type { FC } from 'react';
import styles from './FlowListOperation.module.scss';

type FlowListOperationCellProps = {
  editAuth?: string | string[];
  deleteAuth?: string | string[];
  onEdit: () => void;
  onDelete: () => void;
};

const FlowListOperationCell: FC<FlowListOperationCellProps> = ({ editAuth, deleteAuth, onEdit, onDelete }) => {
  const formatMessage = useTranslate();

  return (
    <span className={styles.operation}>
      <AuthWrapper auth={editAuth}>
        <i title={formatMessage('common.edit')}>
          <Edit className="custom-operation" style={{ cursor: 'pointer' }} onClick={onEdit} />
        </i>
      </AuthWrapper>
      <AuthWrapper auth={deleteAuth}>
        <i title={formatMessage('common.delete')}>
          <TrashCan className="custom-operation" style={{ cursor: 'pointer' }} onClick={onDelete} />
        </i>
      </AuthWrapper>
    </span>
  );
};

export default FlowListOperationCell;
