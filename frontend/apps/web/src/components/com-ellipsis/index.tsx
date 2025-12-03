import type { CSSProperties, FC } from 'react';

const ComEllipsis: FC<{ children?: string; style?: CSSProperties; title?: string; className?: string }> = ({
  children,
  style,
  title,
  className,
}) => {
  return (
    <div
      className={className}
      title={title ?? children ?? ''}
      style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', ...style }}
    >
      {children}
    </div>
  );
};

export default ComEllipsis;
