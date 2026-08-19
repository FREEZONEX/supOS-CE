import { Copy } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useClipboard } from '@/hooks';
import { type CSSProperties, type FC, useRef } from 'react';
import { Flex, Tooltip } from 'antd';

const ComCopy: FC<{ textToCopy: string | number; title?: string; bg?: boolean; style?: CSSProperties }> = ({
  textToCopy,
  title,
  bg,
  style,
}) => {
  const buttonRef = useRef<any>(null);
  useClipboard(buttonRef, typeof textToCopy === 'string' ? textToCopy : textToCopy + '');
  return (
    <Flex
      ref={buttonRef}
      align="center"
      style={
        bg
          ? {
              padding: 6,
              background: 'var(--ui-switchwrap-bg-color)',
              cursor: 'pointer',
              ...style,
            }
          : {
              cursor: 'pointer',
              ...style,
            }
      }
    >
      <Tooltip title={title}>
        <Copy {...toolbarIconProps} />
      </Tooltip>
    </Flex>
  );
};

export default ComCopy;
