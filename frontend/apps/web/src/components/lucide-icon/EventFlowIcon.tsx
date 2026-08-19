import { Zap } from 'lucide-react';
import type { FC } from 'react';
import type { LucideProps } from 'lucide-react';
import { TOOLBAR_ICON_SIZE, TOOLBAR_ICON_STROKE } from './icon-props';

export type EventFlowIconProps = LucideProps;

/** Tier0 UI — Event Flow icon (Lucide `zap`, matches Flow page event tab). */
const EventFlowIcon: FC<EventFlowIconProps> = ({
  size = TOOLBAR_ICON_SIZE,
  strokeWidth = TOOLBAR_ICON_STROKE,
  ...props
}) => <Zap size={size} strokeWidth={strokeWidth} aria-hidden {...props} />;

export default EventFlowIcon;
