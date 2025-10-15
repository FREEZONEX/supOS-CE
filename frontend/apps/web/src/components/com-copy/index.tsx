import { Copy } from '@carbon/icons-react';
import { useClipboard } from '@/hooks';
import { type FC, useRef } from 'react';
import { Flex, Tooltip } from 'antd';

const ComCopy: FC<{ textToCopy: string | number; title?: string }> = ({ textToCopy, title }) => {
  const buttonRef = useRef<any>(null);
  useClipboard(buttonRef, typeof textToCopy === 'string' ? textToCopy : textToCopy + '');
  return (
    <Flex ref={buttonRef} align="center">
      <Tooltip title={title}>
        <Copy style={{ cursor: 'pointer' }} />
      </Tooltip>
    </Flex>
  );
};

export default ComCopy;
