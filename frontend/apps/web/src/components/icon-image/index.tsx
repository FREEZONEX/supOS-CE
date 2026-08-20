import { Image as AntImage, type ImageProps } from 'antd';
import { createElement, type FC, useEffect, useState } from 'react';
import defaultIconUrl from '@/assets/home-icons/default.svg';
import { CUSTOM_APP_ICON_PRE, CUSTOM_MENU_ICON_PRE, CUSTOM_MENU_ICON_PRE1 } from '@/common-types/constans';
import { resolveLegacyLucideIcon, resolveMenuLucideIcon } from '@/components/lucide-icon';
import type { ResourceProps } from '@/stores/types';
import { checkImageExists, getImageSrcByTheme } from '@/utils/url-util';

interface IconImageProps extends Partial<ImageProps> {
  theme: string;
  iconName?: string;
  item?: ResourceProps;
}

const isCustomUploadedIcon = (iconName?: string) => {
  const value = String(iconName || '');
  return (
    value.includes(CUSTOM_MENU_ICON_PRE) || value.includes(CUSTOM_MENU_ICON_PRE1) || value.includes(CUSTOM_APP_ICON_PRE)
  );
};

const CustomUploadedIcon: FC<IconImageProps> = ({ theme, iconName, ...props }) => {
  const [imageSrc, setImageSrc] = useState(defaultIconUrl);

  useEffect(() => {
    const { themeImageUrl, defaultImageUrl, fallbackImageUrl } = getImageSrcByTheme(theme, iconName);
    if (themeImageUrl && defaultImageUrl) {
      const loadImage = async () => {
        const themeExists = await checkImageExists(themeImageUrl);
        if (themeExists) {
          setImageSrc(themeImageUrl);
        } else {
          const defaultExists = await checkImageExists(defaultImageUrl);
          if (defaultExists) {
            setImageSrc(defaultImageUrl);
          } else {
            setImageSrc(fallbackImageUrl);
          }
        }
      };

      loadImage();
    } else {
      setImageSrc(fallbackImageUrl);
    }
  }, [theme, iconName]);

  return <AntImage src={imageSrc} preview={false} fallback={defaultIconUrl} {...props} />;
};

const IconImage: FC<IconImageProps> = ({ theme, iconName, item, width, height, style, className, ...props }) => {
  if (!isCustomUploadedIcon(iconName)) {
    const Icon = item ? resolveMenuLucideIcon(item) : resolveLegacyLucideIcon(iconName);
    const numericWidth = typeof width === 'number' ? width : undefined;
    const numericHeight = typeof height === 'number' ? height : undefined;
    const size = numericWidth ?? numericHeight ?? 16;

    return createElement(Icon, { size, strokeWidth: 1.75, className, style, 'aria-hidden': true });
  }

  return (
    <CustomUploadedIcon
      theme={theme}
      iconName={iconName}
      width={width}
      height={height}
      style={style}
      className={className}
      {...props}
    />
  );
};

export default IconImage;
