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
import useTranslate from '@/hooks/useTranslate';
import styles from './index.module.scss';

const RoutingManagement: FC<PageProps> = ({ title }) => {
  const formatMessage = useTranslate();
  const [activeKey, setActiveKey] = useState('services');
  const [pendingRouteDetail, setPendingRouteDetail] = useState<any>(null);

  const menuItems = [
    { key: 'services', icon: <DataStructured size={16} />, label: formatMessage('kong.services') },
    { key: 'routes', icon: <DirectionFork size={16} />, label: formatMessage('kong.routes') },
    { key: 'consumers', icon: <UserMultiple size={16} />, label: formatMessage('kong.consumers') },
    { key: 'plugins', icon: <PlugFilled size={16} />, label: formatMessage('kong.plugins') },
    { key: 'certificates', icon: <Certificate size={16} />, label: formatMessage('kong.certificates') },
  ];

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
      <ComContent hasBack={false} title={title || formatMessage('kong.routingManagement')}>
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
