import { Tag } from 'antd';
import classNames from 'classnames';
import type { FC, ReactNode } from 'react';
import styles from './FlowStatusTag.module.scss';

export type FlowStatusValue = 'RUNNING' | 'PENDING' | 'STOPPED' | 'DISABLED' | 'DRAFT';

const STATUS_CLASS_MAP: Record<string, string> = {
  RUNNING: styles.running,
  PENDING: styles.pending,
  STOPPED: styles.stopped,
  DISABLED: styles.disabled,
  DRAFT: styles.draft,
};

type FlowStatusTagProps = {
  status?: string | null;
  children: ReactNode;
  title?: string;
  ellipsis?: boolean;
  className?: string;
};

const FlowStatusTag: FC<FlowStatusTagProps> = ({ status, children, title, ellipsis, className }) => {
  const variantClass = (status && STATUS_CLASS_MAP[status]) || styles.fallback;

  return (
    <Tag className={classNames(styles.tag, variantClass, ellipsis && styles.ellipsis, className)} title={title}>
      {children}
    </Tag>
  );
};

export default FlowStatusTag;
