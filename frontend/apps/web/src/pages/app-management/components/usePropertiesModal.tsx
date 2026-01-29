import { useCallback, useRef, useState } from 'react';
import { App, Button, Flex, Input } from 'antd';
import ProModal from '@/components/pro-modal';
import useTranslate from '@/hooks/useTranslate.ts';
import { updateConfigApi } from '@/apis/inter-api/third-apps.ts';
const I18N_NAME = 'AppManagement';

const usePropertiesModal = ({ successBackFn }: any) => {
  const formatMessage = useTranslate(I18N_NAME);
  const commonFormatMessage = useTranslate();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();
  const infoRef = useRef<any>();
  const [value, setValue] = useState('');
  const setModalOpen = useCallback(
    (info: any) => {
      setOpen(true);
      setValue(info?.appProperties);
      infoRef.current = info;
    },
    [setOpen]
  );

  const onClose = () => {
    setOpen(false);
    setLoading(false);
  };

  const onSave = () => {
    updateConfigApi({
      appId: infoRef.current?.appId,
      properties: value,
    }).then(() => {
      onClose();
      successBackFn?.();
      message.success(formatMessage('updateConfigSuccess'));
    });
  };

  const Dom = (
    <ProModal className="labelModalWrap" open={open} onCancel={onClose} title={formatMessage('updateConfig')} size="sm">
      <Input.TextArea
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder={commonFormatMessage('common.commonPlaceholder')}
        rows={15}
        allowClear
        style={{ marginBottom: 20 }}
      />
      <Flex gap="10px">
        <Button
          loading={loading}
          style={{
            backgroundColor: 'var(--supos-button-def-10)',
            color: 'var(--supos-text-color)',
          }}
          color="default"
          variant="filled"
          block
          onClick={onClose}
        >
          {commonFormatMessage('common.cancel')}
        </Button>
        <Button loading={loading} type="primary" variant="solid" onClick={onSave} block>
          {commonFormatMessage('common.save')}
        </Button>
      </Flex>
    </ProModal>
  );
  return {
    PropertiesModal: Dom,
    setPropertiesOpen: setModalOpen,
  };
};

export default usePropertiesModal;
