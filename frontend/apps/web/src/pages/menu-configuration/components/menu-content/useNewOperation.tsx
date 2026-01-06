import useTranslate from '../../../../hooks/useTranslate.ts';
import { useState } from 'react';
import { Form } from 'antd';
import ProModal from '@/components/pro-modal';
import OperationForm from '../../../../components/operation-form';
import { passwordRegex } from '@/utils';

const useNewOperation = ({ onSuccessBack }: { onSuccessBack?: any } = {}) => {
  const formatMessage = useTranslate();
  const [open, setOpen] = useState(false);
  const [type, setType] = useState('add');
  const [form] = Form.useForm();

  const onClose = () => {
    setOpen(false);
    form.resetFields();
  };

  const onNewOperationOpen = (info?: any) => {
    setType(!info ? 'add' : 'edit');
    form.setFieldsValue({
      ...info,
    });
    setOpen(true);
  };

  const onSave = async () => {
    const data: any = await form.validateFields();
    onSuccessBack?.(data, type);
    setOpen(false);
  };

  const NewOperation = (
    <ProModal
      size="xxs"
      title={formatMessage(type === 'add' ? 'MenuConfiguration.addOperation' : 'MenuConfiguration.editOperation')}
      open={open}
      onCancel={onClose}
    >
      <OperationForm
        formConfig={{
          layout: 'vertical',
          labelCol: { span: 24 },
          wrapperCol: { span: 24 },
        }}
        buttonConfig={{
          block: true,
        }}
        style={{ padding: 0 }}
        // loading={loading}
        form={form}
        onCancel={onClose}
        onSave={onSave}
        formItemOptions={[
          {
            name: 'id',
            hidden: true,
          },
          {
            label: formatMessage('MenuConfiguration.operationCode'),
            name: 'code',
            properties: {
              disabled: type === 'edit',
            },
            rules: [
              { required: true, message: formatMessage('rule.required') },
              { pattern: passwordRegex, message: formatMessage('rule.password') },
            ],
          },
          {
            label: formatMessage('MenuConfiguration.operationName'),
            name: 'name',
            properties: {},
          },
          {
            label: formatMessage('MenuConfiguration.operationDescription'),
            name: 'description',
            properties: {},
          },
          {
            type: 'divider',
          },
        ]}
      />
    </ProModal>
  );
  return {
    NewOperation,
    onNewOperationOpen,
  };
};

export default useNewOperation;
