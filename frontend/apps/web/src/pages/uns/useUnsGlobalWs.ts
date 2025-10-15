import { useWebSocket } from 'ahooks';
import { useState } from 'react';
import { getToken } from '@/utils/auth';

export type IcmpStatesType = { topic: string; status: 0 | 1 }[];
interface WsResponseDataProps {
  icmpStates?: IcmpStatesType;
  [key: string]: any;
}

const useUnsGlobalWs = () => {
  const [data, setData] = useState<WsResponseDataProps>({});

  useWebSocket(
    `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/inter-api/supos/uns/ws?globalTopology=true&token=${getToken()}`,
    {
      reconnectLimit: 0,
      onMessage: (event) => {
        const data = JSON.parse(event.data);
        setData(data);
      },
      onError: (error) => console.error('WebSocket error:', error),
    }
  );
  const { icmpStates, mountStatus, ...topologyData } = data;
  return {
    topologyData,
    icmpStates: icmpStates || [],
    mountStatus: mountStatus || {},
  };
};

export default useUnsGlobalWs;
