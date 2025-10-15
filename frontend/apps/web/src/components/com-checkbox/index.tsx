import { Checkbox, type CheckboxProps } from 'antd';
import type { CSSProperties, FC, ReactNode } from 'react';
import classNames from 'classnames';
import './index.scss';

export interface ComCheckboxProps extends CheckboxProps {
  label?: ReactNode;
  disabled?: boolean;
  readonly?: boolean;
  rootStyle?: CSSProperties;
  rootClassname?: CSSProperties;
}

const ComCheckbox: FC<ComCheckboxProps> = ({
  readonly,
  rootClassname,
  rootStyle,
  label,
  disabled,
  children,
  ...restProps
}) => {
  return (
    <div className={classNames('custom-checkbox-wrapper', rootClassname)} style={rootStyle}>
      {readonly ? (
        (label ?? children)
      ) : (
        <Checkbox {...restProps} disabled={disabled} className={classNames('com-checkbox', restProps.className)}>
          {label ?? children}
        </Checkbox>
      )}
    </div>
  );
};

export default ComCheckbox;
