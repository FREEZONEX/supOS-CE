import { createElement, type FC } from 'react';
import type { LucideProps } from 'lucide-react';
import type { ResourceProps } from '@/stores/types';
import { resolveMenuLucideIcon } from './menu-icon-map';

export interface MenuLucideIconProps extends Omit<LucideProps, 'ref'> {
  item: ResourceProps;
}

const isUploadedIconPath = (icon?: string) => {
  const value = String(icon || '');
  return (
    value.startsWith('/api/core/assets/') ||
    value.startsWith('/inter-api/supos/attachment/download') ||
    value.startsWith('http://') ||
    value.startsWith('https://')
  );
};

const MenuLucideIcon: FC<MenuLucideIconProps> = ({
  item,
  size = 16,
  strokeWidth = 1.75,
  className,
  style,
  ...props
}) => {
  if (isUploadedIconPath(item.icon)) {
    return (
      <img
        src={item.icon}
        alt=""
        aria-hidden={props['aria-hidden'] ?? true}
        className={className}
        style={{ width: size, height: size, objectFit: 'contain', ...style }}
      />
    );
  }

  const Icon = resolveMenuLucideIcon(item);
  return createElement(Icon, { size, strokeWidth, className, 'aria-hidden': true, ...props, style });
};

export default MenuLucideIcon;
