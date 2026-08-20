import { useState } from 'react';
import useSSE from '@/hooks/useSSE.ts';

export type IcmpStatesType = { topic: string; status: 0 | 1 }[];

interface WsResponseDataProps {
  icmpStates?: IcmpStatesType;
  mountStatus?: Record<string, string>;
  [key: string]: any;
}

const useUnsGlobalWs = () => {
  const [data, setData] = useState<WsResponseDataProps>({});
  const url = '/api/core/uns/newMsg?globalTopology=true';

  useSSE(url, {
    onMessage: (event) => {
      if (event.data === 'Connected') {
        return;
      }
      try {
        setData(JSON.parse(event.data));
      } catch {
        setData({});
      }
    },
  });

  const { icmpStates, mountStatus, ...topologyData } = data;
  return {
    topologyData,
    icmpStates: icmpStates || [],
    mountStatus: mountStatus || {},
  };
};

export default useUnsGlobalWs;
