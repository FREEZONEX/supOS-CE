import { useTranslate } from '@/hooks';
import { MAX_LENGTHS } from '@/utils/limits';
import type { FolderNode, NotebookDetail } from '@/pages/notebook/types';
import { Folder } from '@/components/lucide-icon/carbon';
import { Form, Input, Modal, TreeSelect } from 'antd';
import type { DataNode } from 'antd/es/tree';
import { type FC, useEffect } from 'react';
import styles from './NotebookModal.module.scss';

interface NotebookModalProps {
  open: boolean;
  title: string;
  folders: FolderNode[];
  notebook?: NotebookDetail | null;
  selectedFolderId: number;
  confirmLoading?: boolean;
  onCancel: () => void;
  onSubmit: (values: { name: string; description?: string; folderId: number }) => Promise<void> | void;
}

const folderTitle = (name: string) => (
  <span className={styles.folderTitle} title={name}>
    <Folder size={16} className={styles.folderIcon} />
    <span className={styles.folderName}>{name}</span>
  </span>
);

const toTreeData = (folders: FolderNode[]): DataNode[] =>
  folders.map((folder) => ({
    title: folderTitle(folder.name),
    value: folder.id,
    key: folder.id,
    children: toTreeData(folder.children || []),
  }));

const NotebookModal: FC<NotebookModalProps> = ({
  open,
  title,
  folders,
  notebook,
  selectedFolderId,
  confirmLoading,
  onCancel,
  onSubmit,
}) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm<{ name: string; description?: string; folderId: number }>();
  const maxNotebookNameLength = MAX_LENGTHS.name;

  useEffect(() => {
    form.setFieldsValue({
      name: notebook?.name || '',
      description: notebook?.description || '',
      folderId: notebook?.folderId ?? selectedFolderId,
    });
  }, [form, notebook, open, selectedFolderId]);

  return (
    <Modal
      open={open}
      title={title}
      confirmLoading={confirmLoading}
      onCancel={onCancel}
      onOk={async () => {
        const values = await form.validateFields();
        await onSubmit(values);
        form.resetFields();
      }}
      destroyOnHidden
    >
      <Form form={form} layout="vertical">
        <Form.Item
          label={formatMessage('Notebook.notebookNameLabel', {}, 'Notebook Name')}
          name="name"
          rules={[
            {
              required: true,
              whitespace: true,
              message: formatMessage('Notebook.notebookNameRequired', {}, 'Notebook name is required'),
            },
            {
              max: maxNotebookNameLength,
              message: formatMessage(
                'Notebook.notebookNameTooLong',
                { length: maxNotebookNameLength },
                `Notebook name cannot exceed ${maxNotebookNameLength} characters`
              ),
            },
          ]}
        >
          <Input
            maxLength={maxNotebookNameLength}
            showCount
            placeholder={formatMessage('Notebook.notebookNamePlaceholder', {}, 'Notebook name')}
          />
        </Form.Item>
        <Form.Item label={formatMessage('Notebook.descriptionLabel', {}, 'Description')} name="description">
          <Input.TextArea
            rows={3}
            maxLength={MAX_LENGTHS.description}
            showCount
            placeholder={formatMessage('Notebook.descriptionPlaceholder', {}, 'Description')}
          />
        </Form.Item>
        <Form.Item label={formatMessage('Notebook.folderLabel', {}, 'Folder')} name="folderId" initialValue={0}>
          <TreeSelect
            className={styles.folderSelect}
            popupClassName={styles.folderTreePopup}
            treeDefaultExpandAll
            allowClear={false}
            treeData={[
              {
                title: folderTitle(formatMessage('Notebook.rootFolder', {}, 'Root')),
                value: 0,
                key: 0,
                children: toTreeData(folders),
              },
            ]}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default NotebookModal;
