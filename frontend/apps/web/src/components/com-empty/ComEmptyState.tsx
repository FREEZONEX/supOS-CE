import classNames from 'classnames';
import type { FC, ReactNode } from 'react';
import ComEmpty from './index';
import ComEmptyDescription from './ComEmptyDescription';
import styles from './index.module.scss';

export interface ComEmptyStateProps {
  title?: ReactNode;
  description?: ReactNode;
  className?: string;
  variant?: 'page' | 'inline';
}

const ComEmptyState: FC<ComEmptyStateProps> = ({ title, description, className, variant = 'page' }) => {
  const emptyDescription = (() => {
    if (!title && !description) {
      return undefined;
    }
    if (!title && (typeof description === 'string' || typeof description === 'number')) {
      return description;
    }
    return <ComEmptyDescription title={title} description={description} />;
  })();

  return (
    <div className={classNames(styles.emptyWrap, variant === 'inline' && styles.emptyWrapInline, className)}>
      <ComEmpty description={emptyDescription} />
    </div>
  );
};

export default ComEmptyState;
