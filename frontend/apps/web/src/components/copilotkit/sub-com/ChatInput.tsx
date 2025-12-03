import { Button, Flex } from 'antd';
import { type FC, useRef } from 'react';
import type { InputProps } from '@copilotkit/react-ui';
import { useTranslate } from '@/hooks';

const ChatInput: FC<InputProps> = ({ onSend, inProgress }) => {
  const formatMessage = useTranslate();

  const inputRef = useRef<any>(null);
  const onSendHandle = () => {
    onSend?.(inputRef.current?.value);
    if (inputRef.current?.value) {
      inputRef.current.value = '';
    }
  };
  const handleKeyDown = (e: any) => {
    // 阻止按下空格和回车键时 Tooltip 关闭
    if (e.keyCode === 32 || e.keyCode === 13 || e.keyCode === 229) {
      e.stopPropagation();
      if (e.keyCode === 13) {
        onSendHandle();
      }
    }
  };
  return (
    <Flex className="chat-input" align="center" justify="center" gap={10}>
      <input
        onKeyDown={handleKeyDown}
        ref={inputRef}
        style={{ height: 38, borderColor: 'transparent', flex: 1, outlineColor: 'transparent', paddingLeft: 15 }}
        placeholder={formatMessage('common.question')}
      />
      <Button loading={inProgress} style={{ height: 38 }} disabled={inProgress} type="primary" onClick={onSendHandle}>
        {formatMessage('common.send')}
      </Button>
    </Flex>
  );
};

export default ChatInput;
