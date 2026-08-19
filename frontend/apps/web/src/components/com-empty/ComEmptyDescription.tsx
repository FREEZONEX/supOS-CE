import type { FC, ReactNode } from 'react';
import styles from './index.module.scss';

export interface ComEmptyDescriptionProps {
  title?: ReactNode;
  description?: ReactNode;
}

const ComEmptyDescription: FC<ComEmptyDescriptionProps> = ({ title, description }) => {
  if (!title && !description) {
    return null;
  }

  return (
    <div className={styles.emptyDesc}>
      {title ? <span className={styles.emptyLine}>{title}</span> : null}
      {description ? <span className={styles.emptyLine}>{description}</span> : null}
    </div>
  );
};

export default ComEmptyDescription;
