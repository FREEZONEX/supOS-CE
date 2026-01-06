import type { CSSProperties, FC, ReactNode } from 'react';
import { Close } from '@carbon/icons-react';
import classNames from 'classnames';
import styles from './index.module.scss';

export interface ComRightProps {
  children?: ReactNode;
  style?: CSSProperties;
  onCancel?: () => void;
  className?: string;
}

const ComRight: FC<ComRightProps> = ({ children, style, onCancel, className }) => {
  return (
    <div className={classNames(styles['com-right'], className)} style={{ ...style, position: 'relative' }}>
      {onCancel && (
        <div className="close">
          <Close style={{ cursor: 'pointer' }} size={20} onClick={() => onCancel?.()} />
        </div>
      )}
      {children}
    </div>
  );
};

export default ComRight;
