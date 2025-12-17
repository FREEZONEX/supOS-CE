import { useEffect, useState } from 'react';
import { Flex, Segmented } from 'antd';
import ComEllipsis from '@/components/com-ellipsis';
import useTranslate from '@/hooks/useTranslate.ts';
import { fetchBaseStore, useBaseStore } from '@/stores/base';
import ProCard from '@/components/pro-card/ProCard.tsx';
import IconImage from '@/components/icon-image';
import type { ResourceProps } from '@/stores/types.ts';
import { useMenuNavigate } from '@/hooks';
import { useThemeStore } from '@/stores/theme-store.ts';
import styles from './index.module.scss';
import { useActivate } from '@/contexts/tabs-lifecycle-context.ts';

const Functions = () => {
  useActivate(() => {
    fetchBaseStore?.();
  });
  useEffect(() => {
    fetchBaseStore?.();
  }, []);
  const formatMessage = useTranslate();
  const { homeTree } = useBaseStore((state) => ({
    homeTree: state.homeTree,
  }));
  const list = homeTree?.map?.((item) => {
    if (item.children && item.children.length) {
      return item;
    } else {
      return {
        ...item,
        children: [item],
      };
    }
  });
  const primaryColor = useThemeStore((state) => state.primaryColor);
  const handleNavigate = useMenuNavigate();
  const [groupId, setGroupId] = useState<string | null>(list?.[0]?.id || '');

  const handleClickItem = (item: ResourceProps) => {
    handleNavigate(item);
  };
  return (
    <Flex vertical gap={24} style={{ marginBottom: 24 }}>
      <Flex justify="space-between" align="flex-start" gap={16}>
        <ComEllipsis className={styles['title']}>{formatMessage('uns.functions')}</ComEllipsis>
        <Segmented<string>
          style={{ overflowX: 'auto' }}
          options={homeTree?.map((item) => ({
            label: item.showName,
            value: item.id,
            title: item.showName,
          }))}
          onChange={(value) => {
            setGroupId(value);
          }}
        />
      </Flex>
      <Flex gap={16} wrap={true}>
        {list
          ?.find((f) => f.id === groupId)
          ?.children?.map?.((c: any) => {
            // 新手导航 id
            let unsMenuId;
            if (c?.url === '/uns') {
              unsMenuId = 'home_route_uns';
            }
            return (
              <div id={unsMenuId} key={c.id} className={styles['functions-item']}>
                <ProCard
                  header={{
                    title: c.showName,
                    customIcon: c.iconComp ? (
                      <div style={{ justifyContent: 'center', alignItems: 'center', display: 'flex', width: 28 }}>
                        {c.iconComp}
                      </div>
                    ) : (
                      <IconImage theme={primaryColor} iconName={c.icon} width={28} height={28} />
                    ),
                  }}
                  onClick={() => handleClickItem(c)}
                  description={{
                    content: c.showDescription,
                  }}
                />
              </div>
            );
          })}
      </Flex>
    </Flex>
  );
};

export default Functions;
