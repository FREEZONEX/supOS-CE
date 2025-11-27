import { Flex } from 'antd';
import ProCard from '@/components/pro-card/ProCard.tsx';
import useTranslate from '@/hooks/useTranslate.ts';
import styles from './Overview.module.scss';
import type { OverviewProps } from './type.ts';
import type { FC } from 'react';
import { ComEllipsis } from '@/components';

const Overview: FC<OverviewProps> = ({ overviewList }) => {
  const formatMessage = useTranslate();

  return (
    <Flex vertical gap={24} style={{ flexShrink: 0 }}>
      <div className={styles['title']}>{formatMessage('common.overview')}</div>
      <div className={styles['overview']}>
        {overviewList?.map((d: any) => {
          return (
            <ProCard
              border
              key={d.key}
              header={{
                title: formatMessage(d.label),
                customIcon: (
                  <Flex style={{ height: '100%' }} align="center" justify="center">
                    {d.icon}
                  </Flex>
                ),
              }}
              iconBg={false}
              allowHover={false}
              classNames={{
                root: styles['item'],
              }}
              description={false}
              styles={{
                secondaryDescription: {
                  lineHeight: 1,
                },
                card: {
                  background: 'transparent',
                },
                headerTitle: {
                  fontSize: 16,
                  fontWeight: 'bold',
                  color: 'var(--supos-text-color)',
                },
              }}
              secondaryDescription={
                <Flex
                  style={{
                    fontSize: 50,
                    color: 'var(--supos-theme-color)',
                    fontWeight: 300,
                    marginTop: 12,
                    lineHeight: 1,
                  }}
                  align="flex-end"
                  gap={4}
                >
                  <ComEllipsis>{d.value}</ComEllipsis>
                  {d.unit && (
                    <span
                      style={{
                        flex: 1,
                        fontWeight: 400,
                        fontSize: 14,
                        lineHeight: '28px',
                        color: 'var(--supos-description-card-color)',
                      }}
                    >
                      {formatMessage(d.unit)}
                    </span>
                  )}
                </Flex>
              }
            />
          );
        })}
      </div>
    </Flex>
  );
};

export default Overview;
