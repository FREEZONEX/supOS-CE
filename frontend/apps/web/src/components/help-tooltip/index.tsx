import type { FC, ReactNode } from 'react';
import { Tooltip } from 'antd';
import { Help } from '@carbon/icons-react';
import type { TooltipProps } from 'antd';

export interface HelpTooltipProps extends Omit<TooltipProps, 'children'> {
  icon?: ReactNode;
}

const HelpTooltip: FC<HelpTooltipProps> = ({ icon, ...restProps }) => {
  return <Tooltip {...restProps}>{icon ? icon : <Help style={{ cursor: 'help' }} />}</Tooltip>;
};

export default HelpTooltip;
