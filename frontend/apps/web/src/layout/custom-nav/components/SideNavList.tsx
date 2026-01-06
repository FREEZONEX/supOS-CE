import type { FC } from 'react';
import { Flex, Menu } from 'antd';
import type { ResourceProps } from '@/stores/types';
import { useMenuNavigate } from '@/hooks';
import styles from './index.module.scss';
import IconImage from '@/components/icon-image';
import { useThemeStore } from '@/stores/theme-store.ts';

const SideNavList: FC<{ navList: ResourceProps[]; selectedKeys: string[] }> = ({ navList, selectedKeys }) => {
  const handleNavigate = useMenuNavigate();
  const primaryColor = useThemeStore((state) => state.primaryColor);
  const theme = useThemeStore((state) => state.theme);

  const createMenuItems = (): any[] => {
    return navList?.map((parent) => {
      if (parent.children?.length && parent.type !== 2) {
        return {
          key: parent.code!,
          label: parent.showName,
          icon: <IconImage theme={primaryColor} iconName={parent.icon} width={'0.875rem'} height={'0.875rem'} />,
          children: parent.children?.map((child) => ({
            key: child.code!,
            label: (
              <div>
                <Flex align="center" gap={4}>
                  <IconImage theme={primaryColor} iconName={child.icon} width={'0.875rem'} height={'0.875rem'} />
                  {child.showName}
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
              <IconImage theme={primaryColor} iconName={parent.icon} width={'0.875rem'} height={'0.875rem'} />
              {parent.showName}
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
