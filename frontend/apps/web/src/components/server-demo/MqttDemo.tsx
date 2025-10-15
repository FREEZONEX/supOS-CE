import type { FC } from 'react';
import { fromPairs, map } from 'lodash-es';
import ServerDemo from './index';
import demoData from './data';
import { getExampleForJavaType } from '@/utils/example';
import { useBaseStore } from '@/stores/base';

const MqttDemo: FC<any> = ({ instanceInfo, ...other }) => {
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const wsPort = systemInfo?.mqttWebsocketTslPort ?? window.location.port;
  const tcpPort = systemInfo?.mqttTcpPort ?? window.location.port;
  const topic = instanceInfo.topic;
  const hostName = window.location.hostname;
  const fieldExampleList = instanceInfo?.fields?.map((item: any) => {
    return {
      key: item.name,
      value: getExampleForJavaType(item.type, item.name),
      type: item.type,
    };
  });
  const jsObj = fromPairs(map(fieldExampleList, (item) => [item.key, item.value]));

  const data =
    typeof demoData.mqtt === 'function' && demoData.mqtt
      ? demoData.mqtt({
          dataType: instanceInfo.dataType,
          hostName,
          wsPort,
          topic,
          jsObj,
          tcpPort,
          fieldExampleList,
        })
      : {};

  return <ServerDemo {...data} {...(other || {})} />;
};

export default MqttDemo;
