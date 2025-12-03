import type { CSSProperties, FC, ReactNode } from 'react';
import './index.scss';

interface ComTextProps {
  children?: ReactNode;
  style?: CSSProperties;
}
const ComText: FC<ComTextProps> = ({ children, style }) => {
  return (
    <span style={style} className={'com-text'}>
      {children}
    </span>
  );
};

export default ComText;
