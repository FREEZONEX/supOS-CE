import { useTranslate } from '@/hooks';
import { Form, Input, Modal } from 'antd';
import { type FC, useEffect } from 'react';

interface FolderModalProps {
  open: boolean;
  title: string;
  initialName?: string;
  confirmLoading?: boolean;
  onCancel: () => void;
  onSubmit: (name: string) => Promise<void> | void;
}

const FolderModal: FC<FolderModalProps> = ({ open, title, initialName, confirmLoading, onCancel, onSubmit }) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm<{ name: string }>();
  const namePattern = /^[\p{L}\p{N}._ -]+$/u;
  const maxFolderNameLength = 63;

  useEffect(() => {
    form.setFieldsValue({ name: initialName || '' });
  }, [form, initialName, open]);

  return (
    <Modal
      open={open}
      title={title}
      confirmLoading={confirmLoading}
      onCancel={onCancel}
      onOk={async () => {
        const values = await form.validateFields();
        await onSubmit(values.name);
        form.resetFields();
      }}
      destroyOnHidden
    >
      <Form form={form} layout="vertical">
        <Form.Item
          label={formatMessage('Notebook.folderNameLabel', {}, 'Folder Name')}
          name="name"
          rules={[
            {
              required: true,
              whitespace: true,
              message: formatMessage('Notebook.folderNameRequired', {}, 'Folder name is required'),
            },
            {
              max: maxFolderNameLength,
              message: formatMessage(
                'Notebook.folderNameTooLong',
                { length: maxFolderNameLength },
                `Folder name cannot exceed ${maxFolderNameLength} characters`
              ),
            },
            {
              validator: async (_, value: string | undefined) => {
                if (
                  String(value || '')
                    .trim()
                    .toLowerCase()
                    .endsWith('.py')
                ) {
                  throw new Error(
                    formatMessage('Notebook.namePySuffixInvalid', {}, 'Do not include the .py suffix in the name')
                  );
                }
              },
            },
            {
              pattern: namePattern,
              message: formatMessage(
                'Notebook.nameInvalid',
                {},
                'Only Chinese characters, letters, numbers, spaces, dots (.), underscores (_), and hyphens (-) are allowed'
              ),
            },
          ]}
        >
          <Input placeholder={formatMessage('Notebook.folderNamePlaceholder', {}, 'Folder name')} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default FolderModal;
