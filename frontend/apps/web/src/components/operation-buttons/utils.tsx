import type { ReactNode } from 'react';
import { AuthButton } from '@/components/auth';
import InlineLoading, { type InlineLoadingProps } from '../inline-loading';
import { type ButtonProps, Popconfirm, type PopconfirmProps } from 'antd';

export interface OperationProps {
  key: string | number;
  type?: 'Button' | 'Custom' | 'Popconfirm' | 'Loading';
  auth?: string | string[];
  label?: ReactNode;
  title?: string;
  component?: ReactNode;
  icon?: ReactNode;
  extra?: ReactNode;
  disabled?: boolean;
  danger?: boolean;
  onClick?: () => void;
  status?: InlineLoadingProps['status'];
  // icon  disabled onClick 放在外部处理
  button?: Omit<ButtonProps, 'icon' | 'disabled' | 'onClick'>;
  //  disabled onConfirm（对于onClick） 放在外部处理
  popconfirm?: Omit<PopconfirmProps, 'onConfirm' | 'disabled'>;
}

// 获取menu label
export const commonLabelRender = (record: OperationProps | undefined) => {
  if (!record) return null;
  const { type, status = '', label, onClick, popconfirm, disabled } = record || {};
  switch (type) {
    case 'Loading':
      return <InlineLoading status={status || 'active'} description={label} />;
    case 'Popconfirm':
      return (
        <Popconfirm
          {...popconfirm}
          title={popconfirm?.title ?? '未填写Popconfirm的title'}
          onConfirm={onClick}
          disabled={disabled}
        >
          {label}
        </Popconfirm>
      );
    default:
      return label;
  }
};

// 获取单个操作的样式
export const commonOperationRender = (record: OperationProps | undefined) => {
  if (!record) return null;
  const { type, status = '', label, button, onClick, icon, disabled, popconfirm } = record || {};
  switch (type) {
    case 'Loading':
      return <InlineLoading status={status || 'active'} description={label} />;
    case 'Popconfirm':
      return (
        <Popconfirm {...popconfirm} title={popconfirm?.title ?? '未填写Popconfirm的title'} onConfirm={onClick}>
          <AuthButton {...button} icon={icon} disabled={disabled}>
            {label}
          </AuthButton>
        </Popconfirm>
      );
    default:
      return label;
  }
};
