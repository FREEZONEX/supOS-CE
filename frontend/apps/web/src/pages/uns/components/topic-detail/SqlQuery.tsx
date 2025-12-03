import { type FC, useEffect, useState } from 'react';
import MqttDemo from '@/components/server-demo/MqttDemo.tsx';
import { useTranslate } from '@/hooks';
import FetchData from '@/pages/uns/components/topic-detail/FetchData';
import demoData from '@/components/server-demo/data';
import ServerDemo from '@/components/server-demo';

import type { FieldItem } from '@/pages/uns/types';
import CodeSnippet from '@/components/code-snippet';
import ComBtnTabs, { type OptionTypes } from '@/components/com-btn-tabs';
import { useBaseStore } from '@/stores/base';

interface SqlQueryProps {
  instanceInfo: { [key: string]: any };
  id: string;
}

const SqlQuery: FC<SqlQueryProps> = ({ instanceInfo, id }) => {
  const formatMessage = useTranslate();
  const dataBaseType = useBaseStore((state) => state.dataBaseType);
  const [activeTab, setActiveTab] = useState('');
  const [list, setList] = useState<OptionTypes[]>([]);
  const getSQL = () => {
    let code = '';
    if (instanceInfo.fields && activeTab === 'Grafana') {
      if (instanceInfo.dataType === 2) {
        code = `select * from "public"."${instanceInfo.alias}" limit 10`;
      } else {
        code = `select _ct,${instanceInfo.fields.map((e: FieldItem) => e.name)} from \`public\`.\`${instanceInfo.alias}\` where \`_ct\` > NOW - 2h;`;
      }
    }
    return code;
  };
  const reset = () => {
    document.querySelector('.topicDetailContent')?.scrollTo(0, 0);
  };
  useEffect(() => {
    if (id) {
      reset();
      getShowDiv(instanceInfo.dataType);
    }
  }, [id, instanceInfo]);

  const getShowDiv = (type: number) => {
    const list: OptionTypes[] = [{ label: formatMessage('uns.upload'), value: 'upload' }];

    if (type === 1 && dataBaseType?.includes('tdengine')) {
      list.push({ label: formatMessage('uns.dbInfo'), value: 'dbInfo' });
    }

    setList(list);
    setActiveTab(list?.[0].value ?? '');
  };

  const handleSelectTab = (item: OptionTypes) => {
    setActiveTab(item.value);
  };

  const renderActiveTab = () => {
    switch (activeTab) {
      case 'Grafana':
        return (
          <CodeSnippet
            className="codeViewWrap"
            type="multi"
            minCollapsedNumberOfRows={1}
            align="top-right"
            showLessText={formatMessage('uns.showLess')}
            showMoreText={formatMessage('uns.showMore')}
            aria-label={formatMessage('uns.copyToClipboard')}
            copyButtonDescription={formatMessage('uns.copyToClipboard')}
          >
            {getSQL()}
          </CodeSnippet>
        );
      // broker信息
      case 'upload':
        return <MqttDemo instanceInfo={instanceInfo} />;
      case 'fetch':
        return <FetchData instanceInfo={instanceInfo} />;
      // 数据库信息
      case 'dbInfo':
        return <ServerDemo {...(typeof demoData.dbInfo === 'function' ? demoData.dbInfo?.({}) : {})} />;
      default:
        return null;
    }
  };

  return (
    <>
      <ComBtnTabs options={list} activeKey={activeTab} onSelect={handleSelectTab} style={{ marginBottom: 20 }} />
      {renderActiveTab()}
    </>
  );
};
export default SqlQuery;
