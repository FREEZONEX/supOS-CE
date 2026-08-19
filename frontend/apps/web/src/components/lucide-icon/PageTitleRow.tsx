import { Flex } from 'antd';
import type { FC, ReactNode } from 'react';
import PageTitleIcon, { type PageTitleIconProps } from './PageTitleIcon';

export interface PageTitleRowProps extends PageTitleIconProps {
  children: ReactNode;
  className?: string;
}

const PageTitleRow: FC<PageTitleRowProps> = ({ children, className, ...iconProps }) => (
  <Flex align="center" gap={8} className={className}>
    <PageTitleIcon {...iconProps} />
    {children}
  </Flex>
);

export default PageTitleRow;
