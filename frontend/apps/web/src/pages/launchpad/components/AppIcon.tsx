import { Box } from '@/components/lucide-icon/carbon';
import { type FC, useState } from 'react';

interface AppIconProps {
  /** 应用图标地址 */
  iconUrl?: string;
  /** 图片 alt / 兜底可读文案 */
  alt?: string;
  /** 兜底 Box 图标尺寸 */
  size?: number;
  /** 兜底 Box 图标线宽 */
  strokeWidth?: number;
  /** 图片自定义 class（容器样式由父级负责） */
  imgClassName?: string;
}

/**
 * 应用图标：优先渲染 iconUrl，加载失败或缺省时兜底为 Box 图标。
 * 容器样式交给父级，本组件只负责「图片 or 兜底」的切换。
 */
const AppIcon: FC<AppIconProps> = ({ iconUrl, alt = '', size = 24, strokeWidth = 1.75, imgClassName }) => {
  // 记录加载失败的具体 url：iconUrl 变化时天然重置失败态，无需 effect
  const [failedUrl, setFailedUrl] = useState<string>();

  if (iconUrl && failedUrl !== iconUrl) {
    return <img src={iconUrl} alt={alt} className={imgClassName} onError={() => setFailedUrl(iconUrl)} />;
  }

  return <Box size={size} strokeWidth={strokeWidth} />;
};

export default AppIcon;
