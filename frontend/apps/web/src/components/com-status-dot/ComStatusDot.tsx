import React from 'react';
import './ComStatusDot.scss';

type StatusType = 'stop' | 'active' | 'breathing' | 'unactive';

interface StatusDotProps {
  status: StatusType;
  size?: number;
  className?: string;
  ringMinSize?: number;
  ringMaxSize?: number; // 默认 18
  animationDuration?: number;
  ringOpacityMin?: number;
  ringOpacityMax?: number;
}

const ComStatusDot: React.FC<StatusDotProps> = ({
  status,
  size = 8,
  className = '',
  ringMinSize = 0,
  ringMaxSize = 4,
  animationDuration = 2.2,
  ringOpacityMin = 0.15,
  ringOpacityMax = 0.5,
}) => {
  const isBreathing = status === 'breathing';

  const style = isBreathing
    ? ({
        '--ring-min': `${ringMinSize}px`,
        '--ring-max': `${ringMaxSize}px`,
        '--duration': `${animationDuration}s`,
        '--opacity-min': ringOpacityMin,
        '--opacity-max': ringOpacityMax,
      } as React.CSSProperties)
    : {};

  return (
    <span
      className={`status-dot ${status} ${className}`}
      style={{
        width: size,
        height: size,
        minWidth: size,
        ...style,
      }}
      aria-hidden="true"
    />
  );
};

export default ComStatusDot;
