import { Flex, type FlexProps } from 'antd';
import type { FC, ReactNode } from 'react';
import classNames from 'classnames';
import './index.scss';

interface ComItemProps extends Omit<FlexProps, 'children'> {
  label?: string;
  extra?: ReactNode;
  addonBefore?: ReactNode;
}

const ComItem: FC<ComItemProps> = ({ rootClassName, label, extra, addonBefore, ...restProps }) => {
  return (
    <Flex align="center" className={classNames('com-item', rootClassName)} {...restProps}>
      {addonBefore && <div>{addonBefore}</div>}
      <div className="label" title={label}>
        {label}
      </div>
      {extra && <div>{extra}</div>}
    </Flex>
  );
};

export default ComItem;
