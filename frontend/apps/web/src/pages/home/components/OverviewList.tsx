import type { CSSProperties, FC } from 'react';
import { Divider, Flex, Typography } from 'antd';
import type { ResourceProps } from '@/stores/types';
import { useMenuNavigate } from '@/hooks';
import IconImage from '@/components/icon-image';
import { useThemeStore } from '@/stores/theme-store.ts';
import type { ExampleProps } from '@/pages/home';
import ProCardContainer from '@/components/pro-card/ProCardContainer';
import ProCard from '@/components/pro-card/ProCard.tsx';
const { Paragraph } = Typography;

interface MenuListProps {
  list: ExampleProps[];
  type?: string;
  customOptRender?: (params: any) => any;
  loadingViews?: string[];
  style?: CSSProperties;
}

const OverviewList: FC<MenuListProps> = ({ list, type, customOptRender, loadingViews, style }) => {
  const handleNavigate = useMenuNavigate();
  const primaryColor = useThemeStore((state) => state.primaryColor);

  const handleClickItem = (item: ResourceProps) => {
    if (type === 'example') {
      return;
    }
    handleNavigate(item);
  };

  return (
    <div style={{ padding: '0 36px', ...style }}>
      {list.map((item, index) => {
        return (
          // 新手导航 id
          <div key={item.id} id={`home_section_step${index + 1}`}>
            <Paragraph style={{ margin: '30px 0 20px' }}>
              <Flex align="center" gap={4}>
                {item.iconComp || (
                  <IconImage
                    theme={primaryColor}
                    iconName={item.icon}
                    width={24}
                    height={24}
                    style={{ paddingRight: 4, verticalAlign: 'middle' }}
                  />
                )}
                <span style={{ fontSize: 20 }}>{item.showName}</span>
              </Flex>
            </Paragraph>
            <ProCardContainer minWidth={350}>
              {(item?.children?.length ? item?.children : [item])?.map?.((c: any) => {
                // 新手导航 id
                let unsMenuId;
                if (c?.url === '/uns') {
                  unsMenuId = 'home_route_uns';
                }
                return (
                  <div id={unsMenuId} key={c.id}>
                    <ProCard
                      loading={(loadingViews || []).includes(c.id as string)}
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
                        rows: type === 'example' ? 2 : 2,
                      }}
                      actions={type === 'example' ? ([customOptRender?.(c)] as any) : undefined}
                    />
                  </div>
                );
              })}
            </ProCardContainer>
            <Divider />
          </div>
        );
      })}
    </div>
  );
};

export default OverviewList;
