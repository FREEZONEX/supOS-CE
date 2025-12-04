import { Button, Flex } from 'antd';
import ComEllipsis from '@/components/com-ellipsis';
import ComCopy from '@/components/com-copy';
import useTranslate from '@/hooks/useTranslate.ts';
import { useBaseStore } from '@/stores/base';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect.tsx';
import { type CSSProperties, useRef, useState } from 'react';
import { getInstanceInfo } from '@/apis/inter-api';
import { getExampleForJavaType } from '@/utils';
import { fromPairs, map } from 'lodash-es';
import DatabaseInfoModal, { type ModalRef } from './DatabaseInfoModal.tsx';
import { DataBase } from '@carbon/icons-react';
import styles from './index.module.scss';
import HelpTooltip from '../../../../components/help-tooltip';

const Item = ({ item, height = 32, ellipsis = true }: any) => {
  const formatMessage = useTranslate();
  const customStyle: CSSProperties = ellipsis
    ? { whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }
    : {
        overflow: 'auto',
        whiteSpace: 'pre-wrap',
        wordWrap: 'break-word',
        wordBreak: 'break-all',
      };
  return (
    <div key={item.key}>
      <Flex justify="space-between" align="center" style={{ marginBottom: 8 }}>
        <ComEllipsis style={{ fontWeight: 400, fontSize: 12, lineHeight: '20px', color: '#525252' }}>
          {formatMessage(item.label)}
        </ComEllipsis>
        {item?.extra && <div style={{ flexShrink: 0, lineHeight: 1 }}>{item?.extra}</div>}
      </Flex>
      <Flex
        title={item.text || formatMessage('uns.selectTopic')}
        align="center"
        justify="space-between"
        gap={6}
        style={item.style}
      >
        <pre
          style={{
            background: 'var(--supos-bg-color)',
            borderRadius: '3px',
            border: '1px solid #E0E0E0',
            width: '100%',
            height,
            padding: '4px 12px',
            ...customStyle,
          }}
        >
          {item.text || formatMessage('uns.selectTopic')}
        </pre>
        <ComCopy style={{ height }} bg textToCopy={item.text || formatMessage('uns.selectTopic')} />
      </Flex>
    </div>
  );
};

const MQTT = () => {
  const formatMessage = useTranslate();
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const wsPort = systemInfo?.mqttTcpPort ?? window.location.port;
  const mqttList = [
    {
      key: 'url',
      label: 'uns.MQTTUrl',
      style: { marginBottom: 8 },
      text: `mqtt://${window.location.hostname}:${wsPort}`,
    },
    {
      key: 'port',
      label: 'uns.MQTTPort',
      text: wsPort,
    },
  ];
  const [topicInfo, setTopicInfo] = useState<any>(null);
  const modalRef = useRef<ModalRef>(null);
  const [payloadInfo, setPayLoadInfo] = useState<any>(null);
  return (
    <Flex vertical className={styles['mqtt']}>
      <Flex align="center" gap={8} style={{ marginBottom: 5 }}>
        <DataBase size={24} />
        <ComEllipsis style={{ fontWeight: 600 }}>{formatMessage('uns.mqttAccess')}</ComEllipsis>
        <HelpTooltip title={formatMessage('uns.mqttDescription')} />
      </Flex>
      <Flex style={{ flex: 1, overflow: 'hidden' }} vertical>
        {mqttList?.map((item: any) => {
          return <Item item={item} key={item.key} />;
        })}
        <ComEllipsis style={{ fontWeight: 400, fontSize: 12, lineHeight: '20px', color: '#525252', margin: '8px 0' }}>
          {formatMessage('uns.topic')}
        </ComEllipsis>
        <SearchSelect
          apiParams={{
            type: 2,
          }}
          style={{
            width: '100%',
            marginBottom: 8,
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
              setPayLoadInfo(null);
            }
          }}
          labelInValue
        />
        <Item
          height={125}
          ellipsis={false}
          item={{
            key: 'payload',
            label: 'uns.payload',
            text: payloadInfo,
            extra: topicInfo?.withSave2db ? (
              <Button
                // type="link"
                title={formatMessage('uns.databaseInfo')}
                size="small"
                style={{ height: 20 }}
                onClick={() => modalRef.current?.onOpen(topicInfo)}
              >
                <DataBase />
              </Button>
            ) : null,
          }}
        />
      </Flex>
      <DatabaseInfoModal ref={modalRef} />
    </Flex>
  );
};

export default MQTT;
