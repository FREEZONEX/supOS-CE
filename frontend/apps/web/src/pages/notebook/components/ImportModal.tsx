import { importNotebook } from '@/apis/core-api/notebook';
import { useTranslate } from '@/hooks';
import type { FolderNode } from '@/pages/notebook/types';
import { Upload as UploadIcon } from '@/components/lucide-icon/carbon';
import { App, Divider, Modal, TreeSelect, Upload, type UploadFile } from 'antd';
import type { FC } from 'react';
import { useState } from 'react';

const { Dragger } = Upload;

interface ImportModalProps {
  open: boolean;
  folders: FolderNode[];
  selectedFolderId: number;
  onCancel: () => void;
  onSuccess: () => void;
}

const ACCEPTED_EXTENSIONS = ['py'];
const ACCEPT_STRING = '.py';

const buildTreeSelectData = (nodes: FolderNode[]): any[] =>
  nodes.map((node) => ({
    title: node.name,
    value: node.id,
    key: node.id,
    children: buildTreeSelectData(node.children || []),
  }));

const ImportModal: FC<ImportModalProps> = ({ open, folders, selectedFolderId, onCancel, onSuccess }) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [folderId, setFolderId] = useState<number>(selectedFolderId);
  const [loading, setLoading] = useState(false);

  const getImportErrorMessage = (error: any) => {
    if (error?.msg === 'notebook.duplicate' || error?.msg === 'common.duplicateResource') {
      return formatMessage('common.duplicateResource', {}, 'Resource already exists. Please use another name.');
    }

    if (error?.msg === 'file required') {
      return formatMessage('Notebook.import.fileRequired', {}, 'Please select a Notebook file to import');
    }

    if (error?.msg === 'invalid notebook') {
      return formatMessage('Notebook.import.invalidContent', {}, 'The Notebook file content is invalid');
    }

    if (error?.msg === 'Notebook.notebookNameTooLong') {
      return formatMessage('Notebook.import.nameTooLong', {}, 'The Notebook file name is too long');
    }

    if (error?.msg === 'notebook.import.unsupportedFormat') {
      return formatMessage(
        'Notebook.import.unsupportedFormatDetailed',
        {},
        'Import failed. Only .py notebook files are supported.'
      );
    }

    return formatMessage('Notebook.import.failed', {}, 'Import failed');
  };

  const handleCancel = () => {
    setFileList([]);
    setLoading(false);
    onCancel();
  };

  const handleSubmit = async () => {
    if (fileList.length === 0) {
      message.warning(formatMessage('Notebook.import.selectFile', {}, 'Please select a file'));
      return;
    }
    setLoading(true);
    try {
      const formData = new FormData();
      formData.append('file', fileList[0] as any);
      formData.append('folderId', String(folderId || 0));
      await importNotebook(formData);
      message.success(formatMessage('Notebook.import.success', {}, 'Notebook imported successfully'));
      setFileList([]);
      onSuccess();
    } catch (error) {
      message.error(getImportErrorMessage(error));
    } finally {
      setLoading(false);
    }
  };

  const treeData = [
    { title: formatMessage('Notebook.rootFolder', {}, 'Root'), value: 0, key: 0 },
    ...buildTreeSelectData(folders),
  ];

  return (
    <Modal
      title={formatMessage('common.import', {}, 'Import')}
      open={open}
      onCancel={handleCancel}
      onOk={handleSubmit}
      confirmLoading={loading}
      okText={formatMessage('common.import', {}, 'Import')}
      destroyOnClose
      afterOpenChange={(visible) => {
        if (visible) {
          setFolderId(selectedFolderId);
          setFileList([]);
        }
      }}
    >
      <Dragger
        fileList={fileList}
        beforeUpload={(file) => {
          const ext = file.name.split('.').pop()?.toLowerCase();
          if (!ext || !ACCEPTED_EXTENSIONS.includes(ext)) {
            message.warning(formatMessage('Notebook.import.unsupportedFormat', {}, 'Supported formats: .py'));
            return false;
          }
          setFileList([file]);
          return false;
        }}
        onRemove={() => {
          setFileList([]);
        }}
        maxCount={1}
        accept={ACCEPT_STRING}
      >
        <div style={{ padding: '20px 0' }}>
          <UploadIcon size={32} style={{ color: 'var(--ui-text-color)', opacity: 0.45 }} />
          <p style={{ marginTop: 8, color: 'var(--ui-text-color)' }}>
            {formatMessage('Notebook.import.dragHint', {}, 'Click or drag file to this area to upload')}
          </p>
          <p style={{ color: 'var(--ui-text-color)', opacity: 0.45, fontSize: 12 }}>
            {formatMessage('Notebook.import.supportHint', {}, 'Supported formats: .py (Max size: 100MB)')}
          </p>
        </div>
      </Dragger>

      <Divider style={{ margin: '16px 0' }} />
      <div>
        <div style={{ marginBottom: 8, fontWeight: 500 }}>
          {formatMessage('Notebook.import.targetFolder', {}, 'Target Folder')}
        </div>
        <TreeSelect
          value={folderId}
          onChange={setFolderId}
          treeData={treeData}
          style={{ width: '100%' }}
          placeholder={formatMessage('Notebook.import.selectFolder', {}, 'Select folder')}
          treeDefaultExpandAll
        />
      </div>
    </Modal>
  );
};

export default ImportModal;
