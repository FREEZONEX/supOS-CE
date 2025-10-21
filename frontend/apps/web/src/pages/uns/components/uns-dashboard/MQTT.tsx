import { Button, Flex } from 'antd';
import cx from 'classnames';
import styles from './Topology.module.scss';
import ComEllipsis from '@/components/com-ellipsis';
import ComCopy from '@/components/com-copy';
import useTranslate from '@/hooks/useTranslate.ts';
import { useBaseStore } from '@/stores/base';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect.tsx';
import { useRef, useState } from 'react';
import { getInstanceInfo } from '@/apis/inter-api';
import { getExampleForJavaType } from '@/utils';
import { fromPairs, map } from 'lodash';
import DatabaseInfoModal, { type ModalRef } from './DatabaseInfoModal.tsx';

const defaultPayload = JSON.stringify(
  {
    Name: 'Name',
    CurrentValue: 1,
    InitialValue: 1,
  },
  null,
  2
);

const Item = ({ item }: any) => {
  const formatMessage = useTranslate();
  return (
    <div key={item.key}>
      <ComEllipsis style={{ fontWeight: 500, fontSize: 16 }}>{formatMessage(item.label)}</ComEllipsis>
      <Flex
        title={item.text}
        style={{
          background: 'var(--supos-bg-color)',
          padding: '4px 12px',
          margin: '12px 0',
          borderRadius: '3px',
          border: '1px solid #E0E0E0',
        }}
        align="center"
        justify="space-between"
      >
        <pre>{item.text}</pre>
        <ComCopy textToCopy={item.text} />
      </Flex>
    </div>
  );
};

const MQTT = () => {
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const wsPort = systemInfo?.mqttTcpPort ?? window.location.port;
  const formatMessage = useTranslate();
  const [payloadInfo, setPayLoadInfo] = useState<any>(defaultPayload);
  const modalRef = useRef<ModalRef>(null);
  const topicInfo = useRef(null);
  const mqttList = [
    {
      key: 'url',
      label: 'uns.MQTTUrl',
      text: `mqtt://${window.location.hostname}:${wsPort}`,
    },
    {
      key: 'port',
      label: 'uns.MQTTPort',
      text: wsPort,
    },
  ];

  return (
    <Flex vertical className={cx(styles['item'], styles['item-right'])} gap={16}>
      <Flex justify="space-between" align="center">
        <ComEllipsis className={styles['title']}>{formatMessage('uns.mqttAccess')}</ComEllipsis>
        <Button size="small" onClick={() => modalRef.current?.onOpen(topicInfo.current)}>
          {formatMessage('uns.databaseInfo')}
        </Button>
      </Flex>
      <div style={{ flex: 1, background: 'var(--supos-card-bg)', padding: 16, overflow: 'auto' }}>
        {mqttList?.map((item: any) => {
          return <Item item={item} key={item.key} />;
        })}
        <ComEllipsis style={{ fontWeight: 500, fontSize: 16 }}>{formatMessage('uns.topic')}</ComEllipsis>
        <SearchSelect
          style={{
            margin: '12px 0',
            width: '100%',
          }}
          placeholder={formatMessage('uns.namespace')}
          onChange={(e) => {
            if (e?.value) {
              getInstanceInfo({ id: e?.value })
                .then((data) => {
                  topicInfo.current = data;
                  const fieldExampleList = data?.fields?.map((item: any) => {
                    return {
                      key: item.name,
                      value: getExampleForJavaType(item.type, item.name),
                      type: item.type,
                    };
                  });
                  const jsObj = fromPairs(map(fieldExampleList, (item) => [item.key, item.value]));
                  setPayLoadInfo(JSON.stringify(jsObj, null, 2));
                })
                .catch(() => {
                  topicInfo.current = null;
                });
            } else {
              topicInfo.current = null;
              setPayLoadInfo(defaultPayload);
            }
          }}
          labelInValue
        />
        <Item
          item={{
            key: 'payload',
            label: 'uns.payload',
            text: payloadInfo,
          }}
        />
      </div>
      <DatabaseInfoModal ref={modalRef} />
    </Flex>
  );
};

export default MQTT;
