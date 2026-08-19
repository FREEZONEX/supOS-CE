// import Uns from '@/pages/uns/index';
import { type FC } from 'react';
import TabLayout from '../../TabLayout';
import Namespace from './Namespace';

const Uns: FC<{ projectId: string }> = ({ projectId }) => {
  return (
    <TabLayout style={{ paddingLeft: '24px' }}>
      <Namespace projectId={projectId} />
    </TabLayout>
  );
};

export default Uns;
