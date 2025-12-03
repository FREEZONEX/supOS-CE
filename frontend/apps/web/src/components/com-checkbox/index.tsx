import { Checkbox, type CheckboxProps, Flex, Tooltip, type TooltipProps } from 'antd';
import type { CSSProperties, FC, ReactNode } from 'react';
import classNames from 'classnames';
import './index.scss';
import { QuestionCircleOutlined } from '@ant-design/icons';
import ComEllipsis from '../com-ellipsis';

export interface ComCheckboxProps extends CheckboxProps {
  label?: ReactNode;
  disabled?: boolean;
  readonly?: boolean;
  rootStyle?: CSSProperties;
  rootClassname?: CSSProperties;
  tooltip?: TooltipProps;
}

const ComCheckbox: FC<ComCheckboxProps> = ({
  readonly,
  rootClassname,
  rootStyle,
  label,
  disabled,
  children,
  tooltip,
  ...restProps
}) => {
  return (
    <div className={classNames('custom-checkbox-wrapper', rootClassname)} style={rootStyle}>
      {readonly ? (
        (label ?? children)
      ) : (
        <Checkbox {...restProps} disabled={disabled} className={classNames('com-checkbox', restProps.className)}>
          {tooltip ? (
            <Flex gap={8} style={{ overflow: 'hidden' }}>
              <ComEllipsis>{(label ?? children) as string}</ComEllipsis>
              {tooltip && (
                <Tooltip title={tooltip?.title}>
                  <QuestionCircleOutlined />
                </Tooltip>
              )}
            </Flex>
          ) : (
            (label ?? children)
          )}
        </Checkbox>
      )}
    </div>
  );
};

export default ComCheckbox;
