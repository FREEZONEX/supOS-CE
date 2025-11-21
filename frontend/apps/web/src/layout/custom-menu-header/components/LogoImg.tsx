import { Image as AntImage, type ImageProps } from 'antd';
import { type FC, useEffect, useState } from 'react';
import logoBlackWhite from '@/assets/custom-nav/logo-white.svg';
import logoBlack from '@/assets/custom-nav/logo-black.svg';
import LoadingDots from '@/layout/custom-menu-header/components/LoadingDots.tsx';
// import { MENU_TARGET_PATH, STORAGE_PATH } from '@/common-types/constans';
// import { checkImageExists, getBaseUrl } from '@/utils/url-util';
// import { useBaseStore } from '@/stores/base';

interface IconImgProps extends Partial<ImageProps> {
  isDark: boolean;
}
const LogoImg: FC<IconImgProps> = ({ isDark, onClick, ...props }) => {
  const [imageSrc, setImageSrc] = useState('');
  // const { systemInfo } = useBaseStore((state) => ({
  //   systemInfo: state.systemInfo,
  // }));

  useEffect(() => {
    setImageSrc('');
    const loadImage = async () => {
      // const baseUrl = `${getBaseUrl()}${STORAGE_PATH}${MENU_TARGET_PATH}`;
      // const themeLogoUrl = systemInfo?.themeConfig?.navigationIcon
      //   ? `${getBaseUrl()}${systemInfo.themeConfig.navigationIcon}?t=${new Date().getTime()}`
      //   : `${baseUrl}/logo-${isDark ? 'dark' : 'light'}.svg`;
      // const themeExists = await checkImageExists(themeLogoUrl);
      // if (themeExists) {
      //   setImageSrc(themeLogoUrl); // 如果主题图片存在
      // } else {
      //   setImageSrc(isDark ? logoBlackWhite : logoBlack); // 如果默认图片存在
      // }
      setImageSrc(isDark ? logoBlackWhite : logoBlack); // 如果默认图片存在
    };
    loadImage();
  }, [isDark]);
  return (
    <div
      onClick={onClick}
      style={{
        cursor: 'pointer',
        minWidth: 50,
        overflow: 'hidden',
        marginRight: 8,
      }}
    >
      {!imageSrc ? (
        <LoadingDots color={isDark ? 'white' : '#333'} />
      ) : (
        <AntImage src={imageSrc} preview={false} fallback={isDark ? logoBlackWhite : logoBlack} {...props} />
      )}
    </div>
  );
};

export default LogoImg;
