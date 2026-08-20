import type { LucideProps } from 'lucide-react';

/** Figma toolbar / segmented icon size (Tier0 Edge Enterprise UI). */
export const TOOLBAR_ICON_SIZE = 16;
export const TOOLBAR_ICON_STROKE = 1.75;
export const TITLE_ICON_SIZE = 20;
export const TREE_ICON_SIZE = 14;

export const toolbarIconProps: Pick<LucideProps, 'size' | 'strokeWidth' | 'aria-hidden'> = {
  size: TOOLBAR_ICON_SIZE,
  strokeWidth: TOOLBAR_ICON_STROKE,
  'aria-hidden': true,
};

export const titleIconProps: Pick<LucideProps, 'size' | 'strokeWidth' | 'aria-hidden'> = {
  size: TITLE_ICON_SIZE,
  strokeWidth: TOOLBAR_ICON_STROKE,
  'aria-hidden': true,
};

export const treeIconProps: Pick<LucideProps, 'size' | 'strokeWidth' | 'aria-hidden'> = {
  size: TREE_ICON_SIZE,
  strokeWidth: TOOLBAR_ICON_STROKE,
  'aria-hidden': true,
};
