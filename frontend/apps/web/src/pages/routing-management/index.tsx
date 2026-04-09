import { type FC, useState, useCallback } from 'react';
import { Menu } from 'antd';
import { DataStructured, DirectionFork, UserMultiple, PlugFilled, Certificate } from '@carbon/icons-react';
import type { PageProps } from '@/common-types';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComLeft from '@/components/com-layout/ComLeft';
import ServicesTab from './components/services-tab';
import RoutesTab from './components/routes-tab';
import ConsumersTab from './components/consumers-tab';
import PluginsTab from './components/plugins-tab';
import CertificatesTab from './components/certificates-tab';
import styles from './index.module.scss';

const menuItems = [
  { key: 'services', icon: <DataStructured size={16} />, label: 'Services' },
  { key: 'routes', icon: <DirectionFork size={16} />, label: 'Routes' },
  { key: 'consumers', icon: <UserMultiple size={16} />, label: 'Consumers' },
  { key: 'plugins', icon: <PlugFilled size={16} />, label: 'Plugins' },
  { key: 'certificates', icon: <Certificate size={16} />, label: 'Certificates' },
];

const RoutingManagement: FC<PageProps> = ({ title }) => {
  const [activeKey, setActiveKey] = useState('services');
  const [pendingRouteDetail, setPendingRouteDetail] = useState<any>(null);

  const handleViewRoute = useCallback((route: any) => {
    setPendingRouteDetail(route);
    setActiveKey('routes');
  }, []);

  const handleMenuClick = useCallback(({ key }: { key: string }) => {
    if (key !== 'routes') setPendingRouteDetail(null);
    setActiveKey(key);
  }, []);

  return (
    <ComLayout>
      <ComContent hasBack={false} title={title || 'Routing Management'}>
        <ComLayout className={styles['routing-inner']}>
          <ComLeft defaultWidth={360} resize style={{ padding: '8px 0' }}>
            <Menu
              mode="inline"
              selectedKeys={[activeKey]}
              onClick={handleMenuClick}
              items={menuItems}
              className="routing-sidebar-menu"
            />
          </ComLeft>
          <ComContent hasBack={false} mustShowTitle={false}>
            <div className="routing-panel">
              {activeKey === 'services' && <ServicesTab key="services" onViewRoute={handleViewRoute} />}
              {activeKey === 'routes' && (
                <RoutesTab key={`routes-${pendingRouteDetail?.id ?? ''}`} initialDetail={pendingRouteDetail} />
              )}
              {activeKey === 'consumers' && <ConsumersTab key="consumers" />}
              {activeKey === 'plugins' && <PluginsTab key="plugins" />}
              {activeKey === 'certificates' && <CertificatesTab key="certificates" />}
            </div>
          </ComContent>
        </ComLayout>
      </ComContent>
    </ComLayout>
  );
};

export default RoutingManagement;
