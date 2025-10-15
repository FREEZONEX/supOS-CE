import { Image as AntImage, type ImageProps } from 'antd';
import { type FC, useEffect, useState } from 'react';
import defaultIconUrl from '@/assets/home-icons/default.svg';
import { checkImageExists, getImageSrcByTheme } from '@/utils/url-util';

interface IconImageProps extends Partial<ImageProps> {
  theme: string;
  iconName?: string;
}
const IconImage: FC<IconImageProps> = ({ theme, iconName, ...props }) => {
  const [imageSrc, setImageSrc] = useState(defaultIconUrl);
  useEffect(() => {
    const { themeImageUrl, defaultImageUrl, fallbackImageUrl } = getImageSrcByTheme(theme, iconName);
    if (themeImageUrl && defaultImageUrl) {
      // 尝试加载主题图片 -> 默认图片 -> 静态资源
      const loadImage = async () => {
        const themeExists = await checkImageExists(themeImageUrl);
        if (themeExists) {
          setImageSrc(themeImageUrl); // 如果主题图片存在
        } else {
          const defaultExists = await checkImageExists(defaultImageUrl);
          if (defaultExists) {
            setImageSrc(defaultImageUrl); // 如果默认图片存在
          } else {
            setImageSrc(fallbackImageUrl); // 如果都不存在，加载静态资源中的默认图片
          }
        }
      };

      loadImage();
    } else {
      setImageSrc(fallbackImageUrl);
    }
  }, [theme]);

  return <AntImage src={imageSrc} preview={false} fallback={defaultIconUrl} {...props} />;
};

export default IconImage;
