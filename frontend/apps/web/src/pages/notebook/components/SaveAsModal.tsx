import { useTranslate } from '@/hooks';
import { MAX_LENGTHS } from '@/utils/limits';
import { Form, Input, Modal } from 'antd';
import { type FC, useEffect } from 'react';
import styles from './SaveAsModal.module.scss';

const NAME_MAX_LENGTH = MAX_LENGTHS.name;
const DESCRIPTION_MAX_LENGTH = 200;

interface SaveAsModalProps {
  open: boolean;
  confirmLoading?: boolean;
  onCancel: () => void;
  onSubmit: (newName: string, description: string) => Promise<boolean> | boolean;
}

const SaveAsModal: FC<SaveAsModalProps> = ({ open, confirmLoading, onCancel, onSubmit }) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm<{ newName?: string; description?: string }>();
  const newName = Form.useWatch('newName', form) || '';
  const description = Form.useWatch('description', form) || '';

  useEffect(() => {
    if (!open) {
      form.resetFields();
    }
  }, [form, open]);

  return (
    <Modal
      open={open}
      title={formatMessage('Notebook.snapshot.saveAsTitle', {}, 'Save as New File')}
      confirmLoading={confirmLoading}
      onCancel={onCancel}
      onOk={async () => {
        const values = await form.validateFields();
        const success = await onSubmit(values.newName || '', values.description || '');
        if (success) {
          form.resetFields();
        }
      }}
      destroyOnHidden
    >
      <Form form={form} layout="vertical">
        <div className={styles.field}>
          <div className={styles.labelRow}>
            <label className={styles.label}>
              {formatMessage('Notebook.snapshot.newNameLabel', {}, 'Name')} <span className={styles.required}>*</span>
            </label>
            <span className={newName.length >= NAME_MAX_LENGTH ? styles.countError : styles.count}>
              {newName.length}/{NAME_MAX_LENGTH}
            </span>
          </div>
          <Form.Item
            name="newName"
            rules={[
              {
                required: true,
                whitespace: true,
                message: formatMessage('Notebook.notebookNameRequired', {}, 'Notebook name is required'),
              },
              {
                max: NAME_MAX_LENGTH,
                message: formatMessage(
                  'Notebook.snapshot.maxChars',
                  { count: NAME_MAX_LENGTH },
                  'Maximum {count} characters'
                ),
              },
            ]}
          >
            <Input
              placeholder={formatMessage('Notebook.snapshot.newNamePlaceholder', {}, 'Enter new file name here')}
              maxLength={NAME_MAX_LENGTH}
            />
          </Form.Item>
        </div>
        <div className={styles.field}>
          <div className={styles.labelRow}>
            <label className={styles.label}>
              {formatMessage('Notebook.snapshot.descriptionLabel', {}, 'Description')}
            </label>
            <span className={description.length >= DESCRIPTION_MAX_LENGTH ? styles.countError : styles.count}>
              {description.length}/{DESCRIPTION_MAX_LENGTH}
            </span>
          </div>
          <Form.Item
            name="description"
            rules={[
              {
                max: DESCRIPTION_MAX_LENGTH,
                message: formatMessage(
                  'Notebook.snapshot.maxChars',
                  { count: DESCRIPTION_MAX_LENGTH },
                  'Maximum {count} characters'
                ),
              },
            ]}
          >
            <Input.TextArea
              placeholder={formatMessage('Notebook.snapshot.descriptionPlaceholder', {}, 'Enter your description here')}
              rows={3}
              maxLength={DESCRIPTION_MAX_LENGTH}
            />
          </Form.Item>
        </div>
      </Form>
    </Modal>
  );
};

export default SaveAsModal;
