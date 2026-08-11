import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import { WarningFilled } from '@/components/lucide-icon/carbon';
import { Button, Checkbox, Flex } from 'antd';
import { type FC, useEffect, useState } from 'react';

interface RecycleDeleteModalProps {
  open: boolean;
  title: string;
  description: string;
  onCancel: () => void;
  onSubmit: (deleteWithFlows: boolean) => void;
}

const RecycleDeleteModal: FC<RecycleDeleteModalProps> = ({ open, title, description, onCancel, onSubmit }) => {
  const formatMessage = useTranslate();
  const [deleteWithFlows, setDeleteWithFlows] = useState(false);

  useEffect(() => {
    if (open) {
      setDeleteWithFlows(false);
    }
  }, [open]);

  return (
    <ProModal
      className="recycle-delete-modal"
      open={open}
      title={title}
      onCancel={onCancel}
      size="xs"
      draggable={false}
      fullScreenable={false}
      destroyOnHidden
      footer={
        <Flex justify="flex-end" gap={12}>
          <Button onClick={onCancel}>{formatMessage('common.cancel')}</Button>
          <Button danger type="primary" onClick={() => onSubmit(deleteWithFlows)}>
            {formatMessage('common.delete')}
          </Button>
        </Flex>
      }
    >
      <div className="recycle-delete-modal-content">
        <div className="recycle-delete-description">{description}</div>
        <div className="recycle-delete-warning">
          <WarningFilled size={20} className="recycle-delete-warning-icon" />
          <span>{formatMessage('uns.recycleDeleteWarning')}</span>
        </div>
        <Checkbox checked={deleteWithFlows} onChange={(event) => setDeleteWithFlows(event.target.checked)}>
          {formatMessage('uns.recycleDeleteFlows')}
        </Checkbox>
      </div>
    </ProModal>
  );
};

export default RecycleDeleteModal;
