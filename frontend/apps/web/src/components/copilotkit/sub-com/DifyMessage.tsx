import { type FC, useState } from 'react';
import CommonTextMessage from './CommonTextMessage';
import { useDeepCompareEffect } from 'ahooks';
import InlineLoading from '@/components/inline-loading';
import { useBaseStore } from '@/stores/base';

const DifyMessage: FC<any> = (props) => {
  const { args, handler, status } = props;
  const currentUserInfo = useBaseStore((state) => state.currentUserInfo);

  const [isLoading, setIsLoading] = useState(true);
  const [message, setMessage] = useState('');
  // 处理数据块的函数

  function processChunk(rawChunk: any) {
    // 分割多个事件
    const events = rawChunk.split('\n\n');
    events.forEach((event: any) => {
      const lines = event.split('\n');
      lines.forEach((line: any) => {
        if (line.startsWith('data: ')) {
          try {
            const jsonStr = line.replace('data: ', '');
            const data = JSON.parse(jsonStr);

            // 根据事件类型处理不同数据
            switch (data.event) {
              case 'message':
                {
                  if (isLoading) {
                    setIsLoading(false);
                  }
                  setMessage((prev) => {
                    return prev + data.answer;
                  });
                }
                break;
              default:
                break;
            }
          } catch (e) {
            console.error('Error parsing JSON:', e);
          }
        }
      });
    });
  }

  const sendStreamingRequest = async (query: string) => {
    setIsLoading(true);
    try {
      const response = await fetch('http://office.unibutton.com:7580/v1/chat-messages', {
        method: 'POST',
        headers: {
          Authorization: `Bearer app-EHfiOdzTu5hFVRoIrYoUnlXT`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          inputs: {},
          query,
          response_mode: 'streaming',
          conversation_id: '',
          user: currentUserInfo.preferredUsername,
          files: [],
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      if (reader) {
        const decoder = new TextDecoder('utf-8');
        // 递归读取流数据
        function readStream() {
          reader!
            .read()
            .then(({ done, value }) => {
              if (done) {
                handler?.('回答完成');
                return;
              }
              // 处理流数据块
              const chunk = decoder.decode(value, { stream: true });
              processChunk(chunk);
              // 继续读取下一个数据块
              readStream();
            })
            .catch((error) => {
              console.error('Stream reading error:', error);
            });
        }

        readStream();
      }
    } catch (error) {
      console.log(error);
      setMessage('请求失败');
      handler?.('失败');
      setIsLoading(false);
    }
  };

  useDeepCompareEffect(() => {
    // 拿到用户问询内容，发起问询
    if (args.prompt && status === 'executing') {
      sendStreamingRequest(args.prompt);
    }
  }, [args, status]);

  return isLoading ? (
    <InlineLoading status="active" description="Loading data..." />
  ) : (
    <CommonTextMessage>{message ?? 'Loading data...'}</CommonTextMessage>
  );
};

export default DifyMessage;
