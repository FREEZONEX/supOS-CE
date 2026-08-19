import { type FC, useEffect, useRef, useState } from 'react';
import type { ResourceProps } from '@/stores/types';
import { ConfigProvider, Flex, Menu, type MenuProps } from 'antd';
import { useMenuNavigate, useTranslate } from '@/hooks';
import styles from './index.module.scss';
import { MenuLucideIcon } from '@/components/lucide-icon';

type MenuItem = Required<MenuProps>['items'][number];
const SideMenuList: FC<{
  navList: ResourceProps[];
  openHoverNav: boolean;
  setOpenHoverNav: any;
  selectedKeys: string[];
}> = ({ navList, openHoverNav, setOpenHoverNav, selectedKeys }) => {
  const handleNavigate = useMenuNavigate();
  const formatMessage = useTranslate();
  const [items, setItems] = useState<MenuItem[]>([]);
  const [menuSelectedKeys, setSelectedKeys] = useState<string[]>([]);
  const menuRef = useRef<any>(null);
  const menuLabel = (showName?: string) => formatMessage(showName || '', undefined, showName || '');
  const handleClickOutside = (event: any) => {
    if (menuRef.current) {
      if (event.target.closest('.imgWrap')) return;
      if (event.target.closest('.ant-menu-submenu-popup')) return;
      if (!menuRef.current?.contains?.(event.target)) {
        setOpenHoverNav(false);
      }
    }
  };
  useEffect(() => {
    // 当 menu 打开时，监听点击事件
    if (openHoverNav) {
      setItems(
        navList?.map?.((parent) => {
          if (parent.children?.length && parent.type !== 2) {
            return {
              key: parent.code!,
              label: (
                <Flex align="center" gap={4} className={styles['side-menu-list-item']}>
                  <MenuLucideIcon item={parent} size={24} />
                  {menuLabel(parent.showName)}
                </Flex>
              ),
              children: parent?.children?.map((child) => ({
                key: child.code!,
                onClick: () => {
                  handleNavigate(child);
                  setOpenHoverNav?.(false);
                },
                label: (
                  <Flex align="center" gap={4} className={styles['side-menu-list-item']}>
                    <MenuLucideIcon item={child} size={24} />
                    {menuLabel(child.showName)}
                  </Flex>
                ),
              })),
            };
          } else {
            return {
              key: parent.code!,
              label: (
                <Flex align="center" gap={4} className={styles['side-menu-list-item']}>
                  <MenuLucideIcon item={parent} size={24} />
                  {menuLabel(parent.showName)}
                </Flex>
              ),
              onClick: () => {
                handleNavigate(parent);
                setOpenHoverNav?.(false);
              },
            };
          }
        })
      );
      setTimeout(() => {
        setSelectedKeys(selectedKeys);
      });
      document.addEventListener('mousedown', handleClickOutside);
    } else {
      setItems([]);
      document.removeEventListener('mousedown', handleClickOutside);
    }

    // 组件卸载时清除事件监听器
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [openHoverNav, navList, selectedKeys, formatMessage]);

  return openHoverNav ? (
    <div ref={menuRef}>
      <ConfigProvider
        theme={{
          components: {
            Menu: {
              itemSelectedColor: 'var(--ui-theme-color)',
            },
          },
        }}
      >
        <Menu
          key={selectedKeys.join(',')}
          style={{ width: 174, maxHeight: 500 }}
          selectedKeys={menuSelectedKeys}
          items={items}
        />
      </ConfigProvider>
    </div>
  ) : null;
};

export default SideMenuList;
