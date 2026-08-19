import { CloseOutlined } from '@ant-design/icons';
import type { ModalFuncProps } from 'antd/es/modal/interface';

export const DELETE_CONFIRM_MODAL_CLASS = 'delete-confirm-modal';

type FormatMessage = (id: string, values?: Record<string, unknown>) => string;

export function mergeDeleteConfirmProps(props: ModalFuncProps, formatMessage?: FormatMessage): ModalFuncProps {
  const { className, okButtonProps, okText, cancelText, closeIcon, ...rest } = props;

  return {
    ...rest,
    icon: null,
    closable: true,
    centered: true,
    closeIcon: closeIcon ?? <CloseOutlined className="delete-confirm-close-icon" />,
    okText: okText ?? formatMessage?.('common.delete'),
    cancelText: cancelText ?? formatMessage?.('common.cancel'),
    okButtonProps: { danger: true, ...okButtonProps },
    className: [DELETE_CONFIRM_MODAL_CLASS, className].filter(Boolean).join(' ') || undefined,
  };
}

export function openDeleteConfirm(
  modal: { confirm: (props: ModalFuncProps) => void },
  props: ModalFuncProps,
  formatMessage?: FormatMessage
) {
  return modal.confirm(mergeDeleteConfirmProps(props, formatMessage));
}
