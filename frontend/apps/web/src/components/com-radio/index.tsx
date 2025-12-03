import { Radio, type RadioGroupProps, type CheckboxOptionType } from 'antd';
import type { CSSProperties, FC } from 'react';
import classNames from 'classnames';
import './index.scss';

export interface RadioOptionProps extends CheckboxOptionType {
  description?: string;
}

export interface ComRadioProps extends Omit<RadioGroupProps, 'options'> {
  options?: RadioOptionProps[];
  direction?: 'horizontal' | 'vertical';
  rootStyle?: CSSProperties;
  rootClassname?: CSSProperties;
  onClick?: (e: any) => void;
}

const ComRadio: FC<ComRadioProps> = ({
  direction = 'horizontal',
  rootClassname,
  rootStyle,
  options = [],
  onClick,
  ...restProps
}) => {
  return (
    <div className={classNames('custom-radio-container', direction, rootClassname)} style={rootStyle}>
      <Radio.Group {...restProps}>
        {options.map((option) => (
          <Radio key={option.value} value={option.value} disabled={option.disabled} onClick={onClick}>
            {option.label}
            {option.description && <span className="radio-description">{option.description}</span>}
          </Radio>
        ))}
      </Radio.Group>
    </div>
  );
};

export default ComRadio;
