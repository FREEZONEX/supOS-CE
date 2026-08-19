import { exportNotebook } from '@/apis/core-api/notebook';
import { ButtonPermission } from '@/common-types/button-permission';
import ProTable from '@/components/pro-table';
import { useTranslate } from '@/hooks';
import type { NotebookItem } from '@/pages/notebook/types';
import { useBaseStore } from '@/stores/base';
import { hasPermission } from '@/utils/auth';
import { formatTimestamp } from '@/utils/format';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import {
  Copy,
  Download,
  Edit,
  Catalog,
  Folder,
  FlashOff,
  OverflowMenuHorizontal,
  TrashCan,
} from '@/components/lucide-icon/carbon';
import { App, Button, Dropdown, Tag, Tooltip } from 'antd';
import NotebookStatusTag, { isNotebookRunning } from './NotebookStatusTag';
import type { ColumnsType } from 'antd/es/table';
import type { MenuProps } from 'antd';
import { type FC, useMemo } from 'react';

interface NotebookTableProps {
  items: NotebookItem[];
  loading?: boolean;
  onOpenFolder: (item: NotebookItem) => void;
  onOpenNotebook: (item: NotebookItem) => void;
  onEditFolder: (item: NotebookItem) => void;
  onEditNotebook: (item: NotebookItem) => void;
  onDeleteNotebook: (item: NotebookItem) => void;
  onCloneNotebook: (item: NotebookItem) => void;
  onShutdownNotebook: (item: NotebookItem) => void;
  onDeleteFolder: (item: NotebookItem) => void;
}

const NotebookTable: FC<NotebookTableProps> = ({
  items,
  loading,
  onOpenFolder,
  onOpenNotebook,
  onEditFolder,
  onEditNotebook,
  onDeleteNotebook,
  onCloneNotebook,
  onShutdownNotebook,
  onDeleteFolder,
}) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const canManage = hasPermission(ButtonPermission['Notebook.manage']);

  const handleDownload = async (record: NotebookItem, format: string, includeCode?: boolean) => {
    try {
      const resp = await exportNotebook(record.id, {
        exportFormat: format,
        ...(format === 'html' ? { includeCode } : {}),
      });
      const extMap: Record<string, string> = { py: '.py', html: '.html', markdown: '.md' };
      const ext = extMap[format] || '.py';
      const blob = resp instanceof Blob ? resp : new Blob([resp as any]);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = record.name + ext;
      link.click();
      window.URL.revokeObjectURL(url);
      message.success(formatMessage('Notebook.export.success', {}, 'Downloaded successfully'));
    } catch {
      message.error(formatMessage('Notebook.export.failed', {}, 'Download failed'));
    }
  };
  const { modal } = App.useApp();
  const currentUserInfo = useBaseStore((state) => state.currentUserInfo);
  const currentUserDisplayName = currentUserInfo.firstName || currentUserInfo.preferredUsername || '';

  const getOwnerDisplayName = (owner?: string) => {
    if (!owner) return '-';
    if (owner === currentUserInfo.sub && currentUserDisplayName) {
      return currentUserDisplayName;
    }
    return owner;
  };

  const columns = useMemo<ColumnsType<NotebookItem>>(
    () => [
      {
        title: formatMessage('common.name', {}, 'Name'),
        dataIndex: 'name',
        key: 'name',
        width: 200,
        className: 'notebook-col-name',
        ellipsis: true,
        render: (_, record) => {
          const icon = record.type === 'folder' ? <Folder size={16} /> : <Catalog size={16} />;
          const onClick = () => (record.type === 'folder' ? onOpenFolder(record) : onOpenNotebook(record));

          return (
            <Tooltip title={record.name}>
              <button type="button" className="notebook-name-cell" onClick={onClick}>
                <span className="notebook-name-icon">{icon}</span>
                <span className="notebook-name-text" title={record.name}>
                  {record.name}
                </span>
              </button>
            </Tooltip>
          );
        },
      },
      {
        title: formatMessage('Notebook.colType', {}, 'Type'),
        dataIndex: 'type',
        key: 'type',
        width: 100,
        ellipsis: true,
        render: (value) => (
          <Tag className="notebook-type-tag">
            {value === 'folder'
              ? formatMessage('Notebook.typeFolder', {}, 'Folder')
              : formatMessage('Notebook.typeNotebook', {}, 'Notebook')}
          </Tag>
        ),
      },
      {
        title: formatMessage('Notebook.colOwner', {}, 'Owner'),
        dataIndex: 'owner',
        key: 'owner',
        width: 100,
        ellipsis: true,
        render: (value) => {
          const ownerName = getOwnerDisplayName(value);
          return ownerName === '-' ? (
            '-'
          ) : (
            <Tooltip title={ownerName}>
              <span className="notebook-cell-ellipsis">{ownerName}</span>
            </Tooltip>
          );
        },
      },
      {
        title: formatMessage('common.status', {}, 'Status'),
        dataIndex: 'status',
        key: 'status',
        width: 100,
        className: 'notebook-col-status',
        ellipsis: true,
        render: (value, record) => (record.type === 'notebook' ? <NotebookStatusTag status={value} /> : '-'),
      },
      {
        title: formatMessage('Notebook.colCreatedAt', {}, 'Time'),
        dataIndex: 'updatedAt',
        key: 'updatedAt',
        width: 180,
        ellipsis: true,
        render: (value) => {
          const text = value ? formatTimestamp(value) : '-';
          return text === '-' ? (
            '-'
          ) : (
            <Tooltip title={text}>
              <span className="notebook-cell-ellipsis">{text}</span>
            </Tooltip>
          );
        },
      },
      {
        title: formatMessage('common.description', {}, 'Description'),
        dataIndex: 'description',
        key: 'description',
        width: 200,
        className: 'notebook-col-description',
        ellipsis: true,
        render: (value) =>
          value ? (
            <Tooltip title={value}>
              <span className="notebook-cell-ellipsis">{value}</span>
            </Tooltip>
          ) : (
            '-'
          ),
      },
      {
        title: formatMessage('common.operation'),
        key: 'actions',
        width: 120,
        fixed: 'right',
        render: (_, record) => {
          const menuItems: MenuProps['items'] =
            record.type === 'folder'
              ? canManage
                ? [
                    {
                      key: 'edit',
                      label: formatMessage('common.edit', {}, 'Edit'),
                      icon: <Edit size={16} />,
                      onClick: () => onEditFolder(record),
                    },
                    {
                      key: 'delete',
                      label: formatMessage('common.delete', {}, 'Delete'),
                      icon: <TrashCan size={16} />,
                      danger: true,
                      onClick: () => {
                        modal.confirm(
                          mergeDeleteConfirmProps(
                            {
                              title: formatMessage('Notebook.deleteFolderConfirm', {}, 'Delete this folder?'),
                              onOk: () => onDeleteFolder(record),
                            },
                            formatMessage
                          )
                        );
                      },
                    },
                  ]
                : []
              : canManage
                ? [
                    {
                      key: 'edit',
                      label: formatMessage('common.edit', {}, 'Edit'),
                      icon: <Edit size={16} />,
                      onClick: () => onEditNotebook(record),
                    },
                    ...(isNotebookRunning(record.status)
                      ? [
                          {
                            key: 'shutdown',
                            label: formatMessage('Notebook.shutdown', {}, 'Shutdown'),
                            icon: <FlashOff size={16} style={{ color: '#da1e28' }} />,
                            onClick: () => onShutdownNotebook(record),
                          },
                        ]
                      : []),
                    {
                      key: 'clone',
                      label: formatMessage('Notebook.clone', {}, 'Clone'),
                      icon: <Copy size={16} />,
                      onClick: () => onCloneNotebook(record),
                    },
                    {
                      key: 'download',
                      label: formatMessage('Notebook.download', {}, 'Download'),
                      icon: <Download size={16} />,
                      children: [
                        ...(isNotebookRunning(record.status)
                          ? [
                              {
                                key: 'download-html',
                                label: formatMessage('Notebook.downloadHtml', {}, 'Download as HTML'),
                                onClick: () => handleDownload(record, 'html', true),
                              },
                              {
                                key: 'download-html-no-code',
                                label: formatMessage(
                                  'Notebook.downloadHtmlNoCode',
                                  {},
                                  'Download as HTML(exclude code)'
                                ),
                                onClick: () => handleDownload(record, 'html', false),
                              },
                              {
                                key: 'download-markdown',
                                label: formatMessage('Notebook.downloadMarkdown', {}, 'Download as Markdown'),
                                onClick: () => handleDownload(record, 'markdown'),
                              },
                            ]
                          : []),
                        {
                          key: 'download-python',
                          label: formatMessage('Notebook.downloadPython', {}, 'Download Python code'),
                          onClick: () => handleDownload(record, 'py'),
                        },
                      ],
                    },
                    { type: 'divider' as const },
                    {
                      key: 'delete',
                      label: formatMessage('common.delete', {}, 'Delete'),
                      icon: <TrashCan size={16} />,
                      danger: true,
                      onClick: () => {
                        modal.confirm(
                          mergeDeleteConfirmProps(
                            {
                              title: formatMessage('Notebook.deleteNotebookConfirm', {}, 'Delete this notebook?'),
                              onOk: () => onDeleteNotebook(record),
                            },
                            formatMessage
                          )
                        );
                      },
                    },
                  ]
                : [];

          if (!menuItems?.length) {
            return null;
          }
          return (
            <Dropdown
              overlayClassName="pro-table-operation-menu"
              menu={{ items: menuItems }}
              trigger={['click']}
              placement="bottomRight"
            >
              <Button type="text" icon={<OverflowMenuHorizontal size={16} />} />
            </Dropdown>
          );
        },
      },
    ],
    [
      formatMessage,
      currentUserDisplayName,
      currentUserInfo.sub,
      canManage,
      modal,
      onCloneNotebook,
      onDeleteFolder,
      onDeleteNotebook,
      onEditFolder,
      onEditNotebook,
      handleDownload,
      onOpenFolder,
      onOpenNotebook,
      onShutdownNotebook,
    ]
  );

  return (
    <ProTable
      className="notebook-table"
      resizeable
      columnFit
      rowKey={(record) => `${record.type}-${record.id}`}
      columns={columns}
      dataSource={items}
      loading={loading}
      sticky
      fixedPosition
      pagination={{
        pageSize: 20,
        showSizeChanger: true,
        showQuickJumper: true,
        pageSizeOptions: ['10', '20', '50', '100'],
      }}
    />
  );
};

export default NotebookTable;
