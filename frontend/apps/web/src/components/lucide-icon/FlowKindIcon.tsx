import type { FC } from 'react';
import type { LucideProps } from 'lucide-react';
import EventFlowIcon from './EventFlowIcon';
import SourceFlowIcon from './SourceFlowIcon';

export type FlowKind = 'source' | 'event';

export type FlowKindIconProps = LucideProps & {
  flowKind?: FlowKind | string;
};

const FlowKindIcon: FC<FlowKindIconProps> = ({ flowKind, ...props }) =>
  flowKind === 'event' ? <EventFlowIcon {...props} /> : <SourceFlowIcon {...props} />;

export default FlowKindIcon;
