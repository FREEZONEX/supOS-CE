import type { ReactNode } from 'react';
import { confirmModalClassName, createConfirmModalOptions } from '@/components/confirm-modal';

type DeleteConfirmOptions = {
  title: ReactNode;
  content?: ReactNode;
  name?: ReactNode;
  okText?: string;
  cancelText?: string;
  formatMessage?: (id: string, values?: Record<string, unknown>) => string;
};

export const deleteConfirmClassName = confirmModalClassName;

const resolveDeleteConfirmContent = ({ title, content, name, formatMessage }: DeleteConfirmOptions) => {
  if (content) return content;
  if (name && formatMessage) {
    // 名称加粗,让用户在确认前一眼看清删的是哪一条。
    return formatMessage('common.deleteConfirmTarget', { name: <strong>{name}</strong> });
  }
  if (formatMessage) {
    return formatMessage('common.deleteConfirmDescription');
  }
  return title;
};

export const createDeleteConfirmOptions = (options: DeleteConfirmOptions) => {
  const okText = options.okText ?? options.formatMessage?.('common.delete') ?? 'Delete';
  const cancelText = options.cancelText ?? options.formatMessage?.('common.cancel') ?? 'Cancel';

  return createConfirmModalOptions({
    title: options.title,
    content: resolveDeleteConfirmContent(options),
    okText,
    cancelText,
    danger: true,
    okButtonProps: {
      title: okText,
    },
    cancelButtonProps: {
      title: cancelText,
    },
  });
};
