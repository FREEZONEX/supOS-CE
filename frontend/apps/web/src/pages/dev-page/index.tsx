import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { useEffect } from 'react';
import DashboardBinding from '@/pages/uns/components/topic-detail/dashboard-binding/DashboardBinding.tsx';
import { Flex } from 'antd';

const DevPage = () => {
  useEffect(() => {}, []);
  console.log('ces');
  return (
    <ComLayout>
      <ComContent title="test" hasBack={false}>
        <Flex justify="flex-end" style={{ width: '100%' }}>
          <DashboardBinding />
        </Flex>
      </ComContent>
    </ComLayout>
  );
};

export default DevPage;
