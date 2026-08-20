import { GitMerge } from 'lucide-react';
import type { FC } from 'react';
import type { LucideProps } from 'lucide-react';
import { TOOLBAR_ICON_SIZE, TOOLBAR_ICON_STROKE } from './icon-props';

export type SourceFlowIconProps = LucideProps;

/** Tier0 UI — Source Flow icon (Lucide `git-merge`, Figma node 13011:84817). */
const SourceFlowIcon: FC<SourceFlowIconProps> = ({
  size = TOOLBAR_ICON_SIZE,
  strokeWidth = TOOLBAR_ICON_STROKE,
  ...props
}) => <GitMerge size={size} strokeWidth={strokeWidth} aria-hidden {...props} />;

export default SourceFlowIcon;
