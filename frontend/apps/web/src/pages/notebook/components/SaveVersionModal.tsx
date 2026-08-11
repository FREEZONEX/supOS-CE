import { useTranslate } from '@/hooks';
import { MAX_LENGTHS } from '@/utils/limits';
import { Form, Input, Modal } from 'antd';
import { type FC, useEffect } from 'react';

interface SaveVersionModalProps {
  open: boolean;
  confirmLoading?: boolean;
  latestVersionName?: string;
  onCancel: () => void;
  onSubmit: (versionName: string, description: string) => Promise<void> | void;
}

const generateDefaultVersionName = (latestVersionName?: string) => {
  const now = new Date();
  const pad = (value: number) => String(value).padStart(2, '0');
  const dateStr = `${now.getFullYear()}${pad(now.getMonth() + 1)}${pad(now.getDate())}`;

  // 从最新版本名中提取版本号，格式 v{number}_{date}
  let nextVersion = 1;
  if (latestVersionName) {
    const match = latestVersionName.match(/^v(\d+(?:\.\d+)?)/);
    if (match) {
      const current = parseFloat(match[1]);
      nextVersion = Math.floor(current) + 1;
    }
  }

  return `v${nextVersion}_${dateStr}`;
};

const SaveVersionModal: FC<SaveVersionModalProps> = ({
  open,
  confirmLoading,
  latestVersionName,
  onCancel,
  onSubmit,
}) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm<{ versionName: string; description?: string }>();

  useEffect(() => {
    if (open) {
      form.setFieldsValue({
        versionName: generateDefaultVersionName(latestVersionName),
        description: '',
      });
    } else {
      form.resetFields();
    }
  }, [form, open, latestVersionName]);

  return (
    <Modal
      open={open}
      title={formatMessage('Notebook.snapshot.saveVersion', {}, 'Save Version')}
      confirmLoading={confirmLoading}
      onCancel={onCancel}
      onOk={async () => {
        const values = await form.validateFields();
        await onSubmit(values.versionName, values.description || '');
        form.resetFields();
      }}
      destroyOnHidden
    >
      <Form form={form} layout="vertical">
        <Form.Item
          label={formatMessage('Notebook.snapshot.versionNameLabel', {}, 'Name')}
          name="versionName"
          rules={[
            {
              required: true,
              whitespace: true,
              message: formatMessage('Notebook.snapshot.versionNameRequired', {}, 'Please enter a version name'),
            },
            {
              max: MAX_LENGTHS.name,
              message: formatMessage(
                'Notebook.snapshot.maxChars',
                { count: MAX_LENGTHS.name },
                'Maximum {count} characters'
              ),
            },
          ]}
        >
          <Input
            placeholder={formatMessage('Notebook.snapshot.versionNamePlaceholder', {}, 'Enter version name')}
            maxLength={MAX_LENGTHS.name}
          />
        </Form.Item>
        <Form.Item label={formatMessage('Notebook.snapshot.descriptionLabel', {}, 'Description')} name="description">
          <Input.TextArea
            placeholder={formatMessage('Notebook.snapshot.descriptionPlaceholder', {}, 'Type your description here')}
            rows={3}
            maxLength={200}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default SaveVersionModal;
