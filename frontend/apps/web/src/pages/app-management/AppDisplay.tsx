import { useEffect, useState } from 'react';
import AppList from '@/pages/app-management/components/AppList';
import { getApps } from '@/apis/inter-api/apps';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';

const Module = () => {
  const [loading, setLoading] = useState(false);
  const [list, setList] = useState([]);
  const getAppsFn = () => {
    setLoading(true);
    getApps()
      .then((data: any) => {
        setList(data || []);
      })
      .catch(() => {})
      .finally(() => {
        setLoading(false);
      });
  };
  useEffect(() => {
    getAppsFn();
  }, []);
  return (
    <ComLayout loading={loading}>
      <ComContent title={<div></div>} mustShowTitle={false}>
        <AppList list={list} successCallBack={getAppsFn} />
      </ComContent>
    </ComLayout>
  );
};
export default Module;
