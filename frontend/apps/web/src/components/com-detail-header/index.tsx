import type { CSSProperties, FC, ReactNode } from 'react';

import ComBackButton from '@/components/com-back-button';
import './index.scss';

export interface ComDetailHeaderProps {
  rightExtra?: ReactNode;
  titleNode?: ReactNode;
  onBack?: () => void;
  title?: string;
  desc?: string;
  showBack?: boolean;
  showDesc?: boolean;
  style?: CSSProperties;
}

const ComDetailHeader: FC<ComDetailHeaderProps> = ({
  rightExtra,
  onBack,
  title,
  desc,
  showBack,
  showDesc,
  titleNode,
  style,
}) => {
  const showDescription = showDesc ?? true;
  return (
    <div className="com-detail-header" style={style}>
      <div className="left-section">
        {showBack && <ComBackButton onClick={onBack} />}
        <div className="info">
          {titleNode || (
            <h1 className="name" title={title}>
              {title}
            </h1>
          )}
          {showDescription && (
            <p className="description" title={desc}>
              {desc || '--'}
            </p>
          )}
        </div>
      </div>
      {rightExtra && <div className="right-section">{rightExtra}</div>}
    </div>
  );
};

export default ComDetailHeader;
