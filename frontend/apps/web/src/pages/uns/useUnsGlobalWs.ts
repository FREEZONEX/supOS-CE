import { useState } from 'react';
import { getToken } from '@/utils/auth';
import useSSE from '@/hooks/useSSE.ts';

export type IcmpStatesType = { topic: string; status: 0 | 1 }[];
interface WsResponseDataProps {
  icmpStates?: IcmpStatesType;
  [key: string]: any;
}

const useUnsGlobalWs = () => {
  const [data, setData] = useState<WsResponseDataProps>({});

  useSSE('/inter-api/supos/uns/newMsg?globalTopology=true&token=' + getToken(), {
    onMessage: (event) => {
      try {
        const data = JSON.parse(event.data);
        setData(data);
      } catch (e) {
        console.log(e);
        setData({});
      }
    },
    onError: (error) => console.error('WebSocket error:', error),
  });
  const { icmpStates, mountStatus, ...topologyData } = data;
  return {
    topologyData,
    icmpStates: icmpStates || [],
    mountStatus: mountStatus || {},
  };
};

export default useUnsGlobalWs;
