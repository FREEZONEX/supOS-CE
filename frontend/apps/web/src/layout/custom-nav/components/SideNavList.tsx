import type { FC } from 'react';
import { Flex, Menu } from 'antd';
import type { ResourceProps } from '@/stores/types';
import { useMenuNavigate, useTranslate } from '@/hooks';
import styles from './index.module.scss';
import { MenuLucideIcon } from '@/components/lucide-icon';
import { useThemeStore } from '@/stores/theme-store.ts';

const SideNavList: FC<{ navList: ResourceProps[]; selectedKeys: string[] }> = ({ navList, selectedKeys }) => {
  const handleNavigate = useMenuNavigate();
  const formatMessage = useTranslate();
  const theme = useThemeStore((state) => state.theme);
  const menuLabel = (showName?: string) => formatMessage(showName || '', undefined, showName || '');

  const createMenuItems = (): any[] => {
    return navList?.map((parent) => {
      if (parent.children?.length && parent.type !== 2) {
        return {
          key: parent.code!,
          label: menuLabel(parent.showName),
          icon: <MenuLucideIcon item={parent} size={14} />,
          children: parent.children?.map((child) => ({
            key: child.code!,
            label: (
              <div>
                <Flex align="center" gap={4}>
                  <MenuLucideIcon item={child} size={14} />
                  {menuLabel(child.showName)}
                </Flex>
              </div>
            ),
            onClick: () => {
              handleNavigate(child);
            },
          })),
        };
      }

      return {
        key: parent.code!,
        label: (
          <div onClick={() => handleNavigate(parent)}>
            <Flex align="center" gap={4}>
              <MenuLucideIcon item={parent} size={14} />
              {menuLabel(parent.showName)}
            </Flex>
          </div>
        ),
        onClick: () => {
          handleNavigate(parent);
        },
      };
    });
  };
  return (
    <Menu
      className={styles['side-nav-list']}
      mode="inline"
      selectedKeys={selectedKeys}
      theme={theme === 'dark' ? 'dark' : 'light'}
      items={createMenuItems()}
    />
  );
};

export default SideNavList;
