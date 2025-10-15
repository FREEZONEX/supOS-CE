import type { CSSProperties, FC, ReactNode } from 'react';
import { Flex } from 'antd';
import classNames from 'classnames';
import './index.scss';

interface ComDotProps {
  children?: ReactNode;
  style?: CSSProperties;
  className?: string;
  color?: string;
  breathing?: boolean;
}

const ComDot: FC<ComDotProps> = ({
  children,
  style,
  className,
  color = 'var(--supos-theme-color)',
  breathing = false,
}) => {
  const dotClassName = classNames('com-dot-dot', {
    breathing: breathing,
  });
  return (
    <Flex align="center" gap={8} style={style} className={classNames('com-dot', className)}>
      <span className={dotClassName} style={{ backgroundColor: color }} />
      {children && <span>{children}</span>}
    </Flex>
  );
};

export default ComDot;
