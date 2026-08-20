import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { useEffect } from 'react';

const DevPage = () => {
  useEffect(() => {}, []);
  console.log('ces');
  return (
    <ComLayout>
      <ComContent title="test" hasBack={false} />
    </ComLayout>
  );
};

export default DevPage;
