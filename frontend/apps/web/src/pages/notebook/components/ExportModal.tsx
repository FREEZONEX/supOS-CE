import { exportNotebook } from '@/apis/core-api/notebook';
import { useTranslate } from '@/hooks';
import { App, Checkbox, Modal, Radio, Space } from 'antd';
import type { FC } from 'react';
import { useState } from 'react';

interface ExportModalProps {
  open: boolean;
  notebookId: number | null;
  notebookName: string;
  onCancel: () => void;
}

const ExportModal: FC<ExportModalProps> = ({ open, notebookId, notebookName, onCancel }) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [exportFormat, setExportFormat] = useState('py');
  const [includeCode, setIncludeCode] = useState(true);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!notebookId) return;
    setLoading(true);
    try {
      const resp = await exportNotebook(notebookId, {
        exportFormat,
        ...(exportFormat === 'html' ? { includeCode } : {}),
      });

      // 触发浏览器下载
      const extMap: Record<string, string> = { py: '.py', html: '.html', markdown: '.md', md: '.md' };
      const ext = extMap[exportFormat] || '.py';
      const fileName = notebookName + ext;

      const blob = resp instanceof Blob ? resp : new Blob([resp as any]);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      link.click();
      window.URL.revokeObjectURL(url);

      message.success(formatMessage('Notebook.export.success', {}, 'Notebook exported successfully'));
      onCancel();
    } catch {
      message.error(
        formatMessage(
          'Notebook.export.failed',
          {},
          'Export failed. Make sure the notebook is running for HTML/Markdown export.'
        )
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={formatMessage('Notebook.export.title', {}, 'Export Notebook')}
      open={open}
      onCancel={onCancel}
      onOk={handleSubmit}
      confirmLoading={loading}
      okText={formatMessage('common.export', {}, 'Export')}
      destroyOnClose
      afterOpenChange={(visible) => {
        if (visible) {
          setExportFormat('py');
          setIncludeCode(true);
        }
      }}
    >
      <div style={{ marginBottom: 16 }}>
        <div style={{ marginBottom: 8, fontWeight: 500 }}>
          {formatMessage('Notebook.export.selectFormat', {}, 'Select format')}
        </div>
        <Radio.Group value={exportFormat} onChange={(e) => setExportFormat(e.target.value)}>
          <Space direction="vertical">
            <Radio value="py">Python (.py)</Radio>
            <Radio value="html">HTML (.html)</Radio>
            <Radio value="markdown">Markdown (.md)</Radio>
          </Space>
        </Radio.Group>
      </div>

      {exportFormat === 'html' && (
        <Checkbox checked={includeCode} onChange={(e) => setIncludeCode(e.target.checked)}>
          {formatMessage('Notebook.export.includeCode', {}, 'Include code')}
        </Checkbox>
      )}

      {(exportFormat === 'html' || exportFormat === 'markdown') && (
        <div
          style={{
            marginTop: 12,
            padding: '8px 12px',
            background: 'var(--ui-craft-bg-color, #f4f4f4)',
            borderRadius: 4,
            fontSize: 12,
            color: 'var(--ui-text-color)',
            opacity: 0.65,
          }}
        >
          {formatMessage(
            'Notebook.export.runningHint',
            {},
            'Note: HTML and Markdown export require the notebook to be running.'
          )}
        </div>
      )}
    </Modal>
  );
};

export default ExportModal;
