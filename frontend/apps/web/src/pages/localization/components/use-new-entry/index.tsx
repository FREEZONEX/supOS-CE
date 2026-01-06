import useTranslate from '@/hooks/useTranslate.ts';
import ProModal from '@/components/pro-modal';
import { useRef, useState } from 'react';
import { App, Form } from 'antd';
import OperationForm from '@/components/operation-form';
import { addResourcesApi, editResourcesApi } from '@/apis/inter-api/i18n.ts';
import { passwordRegex } from '@/utils';
import { useI18nStore } from '@/stores/i18n-store.ts';

const useNewEntry = ({ onSuccessBack }: { onSuccessBack?: any } = {}) => {
  const formatMessage = useTranslate();
  const [open, setOpen] = useState(false);
  const [type, setType] = useState('add');
  const { message } = App.useApp();
  const info = useRef<any>(null);
  const [form] = Form.useForm();
  const [formItemOptions, setFormItemOptions] = useState<any[]>([]);
  const langData = useI18nStore((state) => state.langList);

  const onNewModalOpen = (_info: any, data?: any) => {
    info.current = _info;
    setType(data ? 'edit' : 'add');
    setFormItemOptions([
      {
        label: formatMessage('Localization.i18nMainKey'),
        name: 'key',
        properties: {
          disabled: !!data,
        },
        rules: [
          { required: true, message: formatMessage('rule.required') },
          { pattern: passwordRegex, message: formatMessage('rule.password') },
        ],
      },
      ...(langData?.map((f: any) => {
        return {
          label: f.languageName,
          name: f.languageCode,
        };
      }) || []),
      {
        type: 'divider',
      },
    ]);
    if (data) {
      form.setFieldsValue({
        key: data?.i18nKey,
        ...(data?.values || {}),
      });
    }
    setOpen(true);
  };
  const onClose = () => {
    setOpen(false);
    form.resetFields();
  };
  const onSave = async () => {
    const data: any = await form.validateFields();
    const { key, ...restData } = data;
    const api = type === 'edit' ? editResourcesApi : addResourcesApi;

    api({ values: restData, key, moduleCode: info.current.moduleCode }).then(() => {
      onClose();
      onSuccessBack?.();
      message.success(formatMessage(type === 'edit' ? 'uns.editSuccessful' : 'uns.newSuccessfullyAdded'));
    });
  };
  const NewEntryModal = (
    <ProModal
      size="xxs"
      title={formatMessage(type === 'add' ? 'Localization.newEntry' : 'Localization.editEntry')}
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
        formItemOptions={formItemOptions}
      />
    </ProModal>
  );

  return {
    NewEntryModal,
    onNewModalOpen,
  };
};

export default useNewEntry;
