import { type FC, type ReactNode, useRef } from 'react';
import { Flex } from 'antd';
import { ChevronDown } from '@carbon/icons-react';
import './index.scss';

const HMenuLabel: FC<{ label: ReactNode; iconUrl?: string; expand?: boolean }> = ({ label, expand }) => {
  const ref = useRef<HTMLDivElement>(null);

  return (
    <Flex ref={ref} gap={8} align={'center'} justify="center" className="menu-label">
      {label}
      {expand && <ChevronDown />}
    </Flex>
  );
};

export default HMenuLabel;
