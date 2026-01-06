import type { CSSProperties, FC } from 'react';

const Module: FC<{ style?: CSSProperties }> = ({ style }) => {
  return (
    <div
      id="iframeMask"
      style={{
        display: 'none',
        width: '100%',
        height: '100%',
        position: 'fixed',
        top: 0,
        left: 0,
        ...style,
      }}
    />
  );
};
export default Module;
