import { useMemo, useState } from 'react';
import { Button, Tabs } from 'antd';
import type { TabsProps } from 'antd';
import { Api } from '@carbon/icons-react';
import { ServerDemo, demoData, MqttDemo } from '@/components';
import { useTranslate } from '@/hooks';
import styles from './ServerTabs.module.scss';
const I18N_NAME = 'OpenData';

interface PropsTypes {
  appSceretKey?: string;
}

const ServerTabs = (props: PropsTypes) => {
  const { appSceretKey } = props;
  const formatMessage = useTranslate(I18N_NAME);
  const [activeKey, setActiveKey] = useState('RestAPI');

  const items: TabsProps['items'] = [
    {
      key: 'RestAPI',
      label: 'RestAPI',
      children: (
        <ServerDemo
          {...(typeof demoData.restApi === 'function' ? demoData.restApi({ appSceretKey }) : demoData.restApi)}
          className={styles.demoBox}
        />
      ),
    },
    // {
    //   key: 'websocket',
    //   label: 'websocket',
    //   children: <ServerDemo {...(demoData?.websocket as any)?.({})} className={styles.demoBox} />,
    // },
    {
      key: 'MQTT',
      label: 'MQTT',
      children: (
        <MqttDemo
          instanceInfo={{
            dataType: 1,
            topic: 'demo',
            fields: [
              {
                name: 'timeStamp',
                type: 'DATETIME',
                unique: true,
              },
              {
                name: 'aa',
                type: 'DOUBLE',
              },
              {
                name: 'quality',
                type: 'INT',
              },
            ],
          }}
          className={styles.demoBox}
        />
      ),
    },
    {
      key: 'mcpServer',
      label: 'MCP Server',
      children: <ServerDemo {...(demoData.mcpServer as any)?.({})} className={styles.demoBox} />,
    },
  ];

  const onChange = (key: string) => {
    console.log(key);
    setActiveKey(key);
  };

  const handleLink = () => {
    window.open('/swagger-ui/');
  };

  const operations = useMemo(() => {
    return activeKey === 'RestAPI' ? (
      <Button type="primary" icon={<Api />} onClick={handleLink}>
        {formatMessage('apiList')}
      </Button>
    ) : null;
  }, [activeKey]);

  return (
    <div className={styles.container}>
      <Tabs
        activeKey={activeKey}
        items={items}
        onChange={onChange}
        tabBarExtraContent={operations}
        tabBarStyle={{ margin: '0 14px 0' }}
      />
    </div>
  );
};

export default ServerTabs;
