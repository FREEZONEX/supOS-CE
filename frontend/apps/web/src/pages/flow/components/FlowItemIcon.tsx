import { Folder } from '@/components/lucide-icon/carbon';
import { FlowKindIcon } from '@/components/lucide-icon';
import cx from 'classnames';
import type { FC } from 'react';
import './FlowItemIcon.scss';

export type FlowItemIconProps = {
  category?: string;
  flowKind?: string;
  size?: 'sm' | 'md';
};

const FlowItemIcon: FC<FlowItemIconProps> = ({ category, flowKind, size = 'md' }) => {
  const isGroup = category === 'group';
  const iconSize = size === 'sm' ? 16 : 20;

  return (
    <span className={cx('flow-item-icon', size === 'sm' && 'flow-item-icon-sm', isGroup && 'flow-item-icon-group')}>
      {isGroup ? (
        <Folder size={iconSize} strokeWidth={1.75} aria-hidden />
      ) : (
        <FlowKindIcon flowKind={flowKind} size={iconSize} strokeWidth={1.75} aria-hidden />
      )}
    </span>
  );
};

export default FlowItemIcon;
