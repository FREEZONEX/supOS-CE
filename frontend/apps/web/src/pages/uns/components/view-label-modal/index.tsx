import { useState } from 'react';
import { Button, Flex } from 'antd';
import { Launch } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';

interface LabelItem {
  id: string | number;
  labelName: string;
}

const Module = ({ toTargetNode }: any) => {
  const formatMessage = useTranslate();
  const [open, setOpen] = useState(false);
  const [labelList, setLabelList] = useState<LabelItem[]>([]);

  const setModalOpen = (labelList: LabelItem[]) => {
    setLabelList(labelList);
    setOpen(true);
  };
  const close = () => {
    setOpen(false);
  };

  const Dom = (
    <ProModal open={open} onCancel={close} title={formatMessage('common.label')} size="xxs">
      <Flex gap={10} wrap>
        {labelList.map((label: LabelItem) => (
          <Button
            key={label.id}
            color="default"
            variant="filled"
            icon={<Launch size={12} style={{ color: '#8d8d8d' }} />}
            iconPosition="end"
            size="small"
            style={{
              height: 'max-content',
              minHeight: '24px',
              maxWidth: '100%',
              border: '1px solid #CBD5E1',
              color: 'var(--supos-text-color)',
              backgroundColor: 'var(--supos-uns-button-color)',
            }}
            onClick={() => {
              toTargetNode('label', { pathType: 7, id: label.id });
              close();
            }}
          >
            <span style={{ maxWidth: '95%', whiteSpace: 'pre-wrap', textAlign: 'left' }}>{label.labelName}</span>
          </Button>
        ))}
      </Flex>
      <Button
        className="viewLabelConfirm"
        color="primary"
        variant="solid"
        onClick={close}
        block
        style={{ marginTop: '20px' }}
        size="large"
      >
        {formatMessage('common.confirm')}
      </Button>
    </ProModal>
  );
  return {
    ViewLabelModal: Dom,
    setLabelOpen: setModalOpen,
  };
};
export default Module;
