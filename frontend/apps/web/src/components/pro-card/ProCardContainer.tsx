import type { CSSProperties, FC, ReactNode } from 'react';
import './proCardContainer.scss';

const ProCardContainer: FC<{ children?: ReactNode; style?: CSSProperties; minWidth?: number }> = ({
  children,
  style,
  minWidth = 260,
}) => {
  return (
    <div className="proCardContainer" style={{ ...style, '--supos-pro-card-grid-min-width': minWidth + 'px' }}>
      {children}
    </div>
  );
};

export default ProCardContainer;
