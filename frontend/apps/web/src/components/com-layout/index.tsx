import Loading from '../loading';
import type { FC, ReactNode } from 'react';
import classNames from 'classnames';
import styles from './index.module.scss';

export interface ComLayoutProps {
  children?: ReactNode;
  loading?: boolean;
  className?: string;
}

const ComLayout: FC<ComLayoutProps> = ({ children, loading, className }) => {
  return (
    <Loading spinning={loading ?? false}>
      <div className={classNames(styles['com-layout'], className)}>{children}</div>
    </Loading>
  );
};

export default ComLayout;
