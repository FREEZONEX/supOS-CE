import { Drawer, type DrawerProps } from 'antd';
import './index.scss';
import type { FC } from 'react';
import classNames from 'classnames';
import { Close } from '@carbon/icons-react';

export interface ComDrawerProps extends DrawerProps {
  onClose?: () => void;
}

const Index: FC<ComDrawerProps> = (props) => {
  const { rootClassName, onClose, style, ...restProps } = props;
  return (
    <Drawer
      closable={false}
      extra={<Close style={{ cursor: 'pointer' }} size={20} onClick={onClose} />}
      {...restProps}
      onClose={onClose}
      style={{ backgroundColor: 'var(--supos-bg-color) !important', color: 'var(--supos-text-color)', ...style }}
      rootClassName={classNames('com-drawer', rootClassName)}
    />
  );
};

export default Index;
