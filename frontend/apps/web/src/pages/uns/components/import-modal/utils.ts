export function processChunk(rawChunk: any, callback: any) {
  // 分割多个事件
  const events = rawChunk.split('\n\n');
  events.forEach((event: any) => {
    const lines = event.split('\n');
    lines.forEach((line: string) => {
      if (line.includes('code')) {
        try {
          const data = JSON.parse(line);
          callback?.(data);
        } catch (e) {
          console.error('Error parsing JSON:', e);
        }
      }
    });
  });
}

export function readerSSE(response: Response, successHandle: any, errorHandle: any) {
  const reader = response.body?.getReader();
  if (reader) {
    const decoder = new TextDecoder('utf-8');
    let hasFinalStatus = false;
    // 递归读取流数据
    function readStream() {
      reader!
        .read()
        .then(({ done, value }) => {
          if (done) {
            if (!hasFinalStatus) errorHandle?.();
            return;
          }
          // 处理流数据块
          const chunk = decoder.decode(value, { stream: true });
          processChunk(chunk, (data: any) => {
            if (data?.finished || data?.progress >= 100) {
              hasFinalStatus = true;
            }
            successHandle?.(data);
          });
          // 继续读取下一个数据块
          readStream();
        })
        .catch((error) => {
          console.error(error);
          errorHandle?.();
        });
    }

    readStream();
  } else {
    errorHandle?.();
  }
}
