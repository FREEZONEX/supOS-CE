import { type FC, useEffect, useRef, useState } from 'react';
import type { ResourceProps } from '@/stores/types';
import { ConfigProvider, Flex, Menu, type MenuProps } from 'antd';
import { useMenuNavigate } from '@/hooks';
import styles from './index.module.scss';
import IconImage from '@/components/icon-image';
import { useThemeStore } from '@/stores/theme-store.ts';

type MenuItem = Required<MenuProps>['items'][number];
const SideMenuList: FC<{
  navList: ResourceProps[];
  openHoverNav: boolean;
  setOpenHoverNav: any;
  selectedKeys: string[];
}> = ({ navList, openHoverNav, setOpenHoverNav, selectedKeys }) => {
  const primaryColor = useThemeStore((state) => state.primaryColor);
  const handleNavigate = useMenuNavigate();
  const [items, setItems] = useState<MenuItem[]>([]);
  const [menuSelectedKeys, setSelectedKeys] = useState<string[]>([]);
  const menuRef = useRef<any>(null);
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
                  <IconImage theme={primaryColor} iconName={parent.icon} width={24} height={24} />
                  {parent.showName}
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
                    <IconImage theme={primaryColor} iconName={child.code} width={24} height={24} />
                    {child.showName}
                  </Flex>
                ),
              })),
            };
          } else {
            return {
              key: parent.code!,
              label: (
                <Flex align="center" gap={4} className={styles['side-menu-list-item']}>
                  <IconImage theme={primaryColor} iconName={parent.icon} width={24} height={24} />
                  {parent.showName}
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
  }, [openHoverNav]);

  return openHoverNav ? (
    <div ref={menuRef}>
      <ConfigProvider
        theme={{
          components: {
            Menu: {
              itemSelectedColor: 'var(--supos-theme-color)',
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
