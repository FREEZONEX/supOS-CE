import type { ModalFuncProps } from 'antd/es/modal/interface';
import classNames from 'classnames';
import styles from './index.module.scss';

export type ConfirmModalController = {
  confirm: (props: ModalFuncProps) => unknown;
};

export type ConfirmModalOptions = Omit<ModalFuncProps, 'icon'> & {
  /** Only destructive actions should use the danger button style. */
  danger?: boolean;
};

export const confirmModalClassName = styles['confirm-modal'];

export const createConfirmModalOptions = ({
  danger = false,
  width = 420,
  centered = true,
  className,
  okButtonProps,
  ...props
}: ConfirmModalOptions): ModalFuncProps => ({
  ...props,
  width,
  centered,
  icon: null,
  className: classNames(confirmModalClassName, className),
  okButtonProps: {
    ...okButtonProps,
    danger: danger || okButtonProps?.danger,
  },
});

export const openConfirmModal = (modal: ConfirmModalController, options: ConfirmModalOptions) =>
  modal.confirm(createConfirmModalOptions(options));
