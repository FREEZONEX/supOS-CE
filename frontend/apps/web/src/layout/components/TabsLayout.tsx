import { type FC, useEffect, useImperativeHandle, useMemo, useState } from 'react';
import { Tabs } from 'antd';
import { useTabs } from '@/layout/useTabs';
import styles from './index.module.scss';
import { createPortal } from 'react-dom';
import ComTags from './com-tags';
import { useLocationNavigate } from '@/routers';
import { TabWrapper } from '@/layout/components/TabWrapper.tsx';
import { useMemoizedFn } from 'ahooks';
import { injectPropsToRouteNode } from '@/utils/node-utils';
import { MenuTypeEnum, type MenuTypeProps } from '@/stores/theme-store.ts';

const TabsLayout: FC<{
  menuType: MenuTypeProps;
  tabsContextRef: any;
}> = ({ menuType, tabsContextRef }) => {
  const { activeTabRoutePath, tabs, onCloseTab, onCloseOtherTab, onRefreshTab, setTabs } = useTabs();
  const [container, setContainer] = useState<HTMLElement | null>(null);
  const navigate = useLocationNavigate();

  const tabItems = useMemo(() => {
    return tabs.map((tab) => ({
      label: tab.title,
      key: tab.routePath,
      children: (
        <div key={tab.key} style={{ height: '100%', overflow: 'hidden' }}>
          <TabWrapper isActive={activeTabRoutePath === tab.routePath}>
            {injectPropsToRouteNode(tab.children, { location: tab.location, title: tab?.title })}
          </TabWrapper>
        </div>
      ),
      closable: tabs.length > 1, // 剩最后一个就不能删除了
    }));
  }, [tabs, activeTabRoutePath]);

  const onTabsChange = useMemoizedFn((tabRoutePath: string) => {
    const { location } = tabs.find((o) => o.routePath === tabRoutePath) || {};
    if (location) {
      navigate(location);
    }
  });

  useEffect(() => {
    if (menuType !== MenuTypeEnum.Top) return;
    // 使用 MutationObserver 检测 DOM 是否挂载
    const observer = new MutationObserver((mutationsList) => {
      for (const mutation of mutationsList) {
        if (mutation.type === 'childList') {
          const target = document.getElementById('custom-header-container');
          if (target) {
            setContainer(target); // 更新 container 状态
            observer.disconnect(); // 找到目标后停止监听
            break;
          }
        }
      }
    });

    // 监听根节点的子节点变化
    const rootNode = document.getElementById('root'); // 假设根节点为 #root
    if (rootNode) {
      observer.observe(rootNode, { childList: true, subtree: true });
    }

    // 清理函数，组件卸载时停止监听
    return () => {
      observer.disconnect();
    };
  }, [menuType]);

  const TabsContextValue = useMemo(
    () => ({
      onCloseTab,
      onCloseOtherTab,
      onRefreshTab,
    }),
    [onCloseTab, onCloseOtherTab, onRefreshTab]
  );

  useImperativeHandle(tabsContextRef, () => {
    return TabsContextValue;
  });
  return (
    <Tabs
      destroyOnHidden={false}
      animated={false}
      style={{ color: 'var(--supos-text-color)' }}
      renderTabBar={() => {
        const TabBar = (
          <ComTags
            setTabs={setTabs}
            onClose={onCloseTab}
            onCloseOther={onCloseOtherTab}
            onRefresh={onRefreshTab}
            activeTag={activeTabRoutePath}
            options={tabItems?.map?.((pan) => {
              return {
                children: pan.label,
                onClick: onTabsChange,
                onClose: onCloseTab,
                key: pan.key,
                closeIcon: tabItems?.length > 1,
              };
            })}
          />
        );
        return container ? createPortal(TabBar, container) : <div style={{ display: 'none' }} />;
      }}
      className={styles['tabs-layout']}
      activeKey={activeTabRoutePath}
      items={tabItems}
      type="card"
    />
  );
};

export default TabsLayout;
