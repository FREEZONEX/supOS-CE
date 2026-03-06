import { forwardRef, useImperativeHandle, useState } from 'react';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import OperationForm from '@/components/operation-form';
import { App, Flex, Form } from 'antd';
import ComButton from '@/components/com-button';
import { addGroup, editGroup } from '@/apis/inter-api/group.ts';

export interface AddGroupModalRef {
  onOpen: (type: number, props?: any) => void;
  onClose: () => void;
}

export interface AddGroupModalProps {
  [key: string]: any;
}

const AddGroupModal = forwardRef<AddGroupModalRef, AddGroupModalProps>(({ refreshRequest }, ref) => {
  const [visible, setVisible] = useState(false);
  const formatMessage = useTranslate();
  const [title, setTitle] = useState('');
  const [form] = Form.useForm();
  const { message } = App.useApp();

  const onOpen = (type: number, props: any) => {
    form.setFieldsValue({
      ...(props || {}),
      type,
    });
    setTitle(formatMessage(props ? 'uns.editGroup' : 'uns.newGroup'));
    setVisible(true);
  };

  const onClose = () => {
    form.resetFields();
    setVisible(false);
  };

  const onSave = async () => {
    const value = await form.validateFields();
    const api = value?.id ? editGroup : addGroup;
    return api(value).then(() => {
      onClose?.();
      refreshRequest?.();
      message.success(formatMessage('common.optsuccess'));
    });
  };

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  return (
    <ProModal
      open={visible}
      onCancel={onClose}
      title={title}
      width={500}
      styles={{
        body: {
          paddingBlockStart: 0,
        },
      }}
    >
      {() => {
        return (
          <OperationForm
            form={form}
            onCancel={onClose}
            onSave={onSave}
            formConfig={{
              layout: 'vertical',
              labelCol: { span: 24 },
              wrapperCol: { span: 124 },
            }}
            formItemOptions={[
              {
                name: 'id',
                hidden: true,
              },
              {
                name: 'type',
                hidden: true,
              },
              {
                name: 'name',
                label: formatMessage('common.name'),
                rules: [{ required: true }],
                properties: {
                  placeholder: formatMessage('common.commonPlaceholder'),
                  allowClear: true,
                },
              },
              {
                name: 'description',
                label: formatMessage('uns.description'),
                properties: {
                  placeholder: formatMessage('common.commonPlaceholder'),
                  allowClear: true,
                },
              },
            ]}
            style={{ padding: 0 }}
            footer={
              <Flex gap="10px" justify="end">
                <ComButton
                  style={{
                    backgroundColor: 'var(--supos-uns-button-color)',
                    color: 'var(--supos-text-color)',
                  }}
                  color="default"
                  variant="filled"
                  onClick={onClose}
                  title={formatMessage('common.cancel')}
                >
                  {formatMessage('common.cancel')}
                </ComButton>
                <ComButton type="primary" variant="solid" onClick={onSave} title={formatMessage('common.save')}>
                  {formatMessage('common.save')}
                </ComButton>
              </Flex>
            }
          />
        );
      }}
    </ProModal>
  );
});

export default AddGroupModal;
