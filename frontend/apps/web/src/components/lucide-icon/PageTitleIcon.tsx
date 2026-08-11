import { createElement, type FC } from 'react';
import type { LucideProps } from 'lucide-react';
import type { ResourceProps } from '@/stores/types';
import { useBaseStore } from '@/stores/base';
import { resolveMenuLucideIcon } from './menu-icon-map';
import { titleIconProps } from './icon-props';

export interface PageTitleIconProps extends Omit<LucideProps, 'ref'> {
  item?: Partial<ResourceProps>;
  resourceKey?: string;
  url?: string;
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

const PageTitleIcon: FC<PageTitleIconProps> = ({ item, resourceKey, url, className, style, ...props }) => {
  const currentMenuInfo = useBaseStore((state) => state.currentMenuInfo);
  const menuItem = {
    ...currentMenuInfo,
    ...item,
    resourceKey: resourceKey ?? item?.resourceKey ?? currentMenuInfo?.resourceKey,
    url: url ?? item?.url ?? currentMenuInfo?.url,
    icon: item?.icon ?? currentMenuInfo?.icon,
    code: item?.code ?? currentMenuInfo?.code,
  } as ResourceProps;

  if (isUploadedIconPath(menuItem.icon)) {
    return (
      <img
        src={menuItem.icon}
        alt=""
        aria-hidden={props['aria-hidden'] ?? true}
        className={className}
        style={{ width: titleIconProps.size, height: titleIconProps.size, objectFit: 'contain', ...style }}
      />
    );
  }

  const Icon = resolveMenuLucideIcon(menuItem);

  return createElement(Icon, { ...titleIconProps, className, ...props, style });
};

export default PageTitleIcon;
