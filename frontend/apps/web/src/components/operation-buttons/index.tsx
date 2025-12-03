import type { CSSProperties, FC } from 'react';
import classNames from 'classnames';
import { AuthWrapper } from '../auth';
import { Button, type ButtonProps } from 'antd';
import './index.scss';

const colorType: any = {
  outlined: {
    color: 'primary',
    variant: 'outlined',
    style: {
      background: 'var(--supos-bg-color)',
    },
  },
  primary: {
    color: 'primary',
    variant: 'solid',
  },
  dark: {
    color: 'default',
    variant: 'solid',
    style: {
      background: 'var(--supos-description-card-color)',
    },
  },
  link: {
    color: 'default',
    variant: 'link',
  },
};

export interface OperationButtonsProps {
  options?: {
    label: string;
    onClick: (item: any) => void;
    type: 'outlined' | 'primary' | 'dark' | 'link';
    btnProps?: ButtonProps;
    auth?: string | string[];
    disabled?: (item: any) => boolean;
  }[];
  record?: any;
  className?: string;
  style?: CSSProperties;
}

const OperationButtons: FC<OperationButtonsProps> = ({ options, record, className, style }) => {
  return (
    <div className={classNames('operation-buittons', className)} style={style}>
      {options?.map((item: any, i) => (
        <AuthWrapper auth={item.auth} key={item.label || `button_${i}`}>
          <Button
            {...(colorType[item.type] || colorType.outlined)}
            {...item.btnProps}
            disabled={item.disabled ? item.disabled(record) : false}
            className="button-item"
            onClick={() => !(item?.disabled ? item.disabled(record) : false) && item?.onClick?.(record, item)}
          >
            {item.label}
          </Button>
        </AuthWrapper>
      ))}
    </div>
  );
};

export default OperationButtons;
