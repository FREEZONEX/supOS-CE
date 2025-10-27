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
import { fromPairs, map } from 'lodash-es';
import DatabaseInfoModal, { type ModalRef } from './DatabaseInfoModal.tsx';
import { DataBase } from '@carbon/icons-react';

const Item = ({ item }: any) => {
  const formatMessage = useTranslate();
  return (
    <div key={item.key}>
      <Flex justify="space-between" align="center">
        <ComEllipsis style={{ fontWeight: 400, fontSize: 12, color: '#525252' }}>
          {formatMessage(item.label)}
        </ComEllipsis>
        {item?.extra && <div style={{ flexShrink: 0 }}>{item?.extra}</div>}
      </Flex>
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
        <pre
          style={{
            overflow: 'auto',
            whiteSpace: 'pre-wrap',
            wordWrap: 'break-word',
            wordBreak: 'break-all',
            width: '100%',
            maxHeight: 150,
          }}
        >
          {item.text}
        </pre>
        <ComCopy textToCopy={item.text} />
      </Flex>
    </div>
  );
};

const MQTT = () => {
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const wsPort = systemInfo?.mqttTcpPort ?? window.location.port;
  const formatMessage = useTranslate();
  const [payloadInfo, setPayLoadInfo] = useState<any>(formatMessage('uns.selectTopic'));
  const modalRef = useRef<ModalRef>(null);
  const [topicInfo, setTopicInfo] = useState<any>(null);
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
      </Flex>
      <div style={{ flex: 1, background: 'var(--supos-card-bg)', padding: 16, overflow: 'auto' }}>
        {mqttList?.map((item: any) => {
          return <Item item={item} key={item.key} />;
        })}
        <ComEllipsis style={{ fontWeight: 400, fontSize: 12, color: '#525252' }}>
          {formatMessage('uns.topic')}
        </ComEllipsis>
        <SearchSelect
          apiParams={{
            type: 2,
          }}
          style={{
            margin: '12px 0',
            width: '100%',
          }}
          placeholder={formatMessage('common.select')}
          onChange={(e) => {
            if (e?.value) {
              getInstanceInfo({ id: e?.value })
                .then((data) => {
                  setTopicInfo(data);
                  const fieldExampleList = data?.fields?.map((item: any) => {
                    return {
                      key: item.name,
                      value: getExampleForJavaType(item.type, item.name),
                      type: item.type,
                    };
                  });
                  if (data?.dataType === 8) {
                    setPayLoadInfo(formatMessage('uns.jsonBExample'));
                  } else {
                    const jsObj = fromPairs(map(fieldExampleList, (item) => [item.key, item.value]));
                    setPayLoadInfo(JSON.stringify(jsObj, null, 2));
                  }
                })
                .catch(() => {
                  setTopicInfo(null);
                });
            } else {
              setTopicInfo(null);
              setPayLoadInfo(formatMessage('uns.selectTopic'));
            }
          }}
          labelInValue
        />
        <Item
          item={{
            key: 'payload',
            label: 'uns.payload',
            text: payloadInfo,
            extra: topicInfo?.withSave2db ? (
              <Button
                title={formatMessage('uns.databaseInfo')}
                size="small"
                onClick={() => modalRef.current?.onOpen(topicInfo)}
              >
                <DataBase />
              </Button>
            ) : null,
          }}
        />
      </div>
      <DatabaseInfoModal ref={modalRef} />
    </Flex>
  );
};

export default MQTT;
