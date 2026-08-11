import { Flex } from 'antd';
import ProCard from '@/components/pro-card/ProCard.tsx';
import useTranslate from '@/hooks/useTranslate.ts';
import ComEllipsis from '@/components/com-ellipsis';
import type { OverviewProps } from './type.ts';
import type { FC } from 'react';
import styles from './index.module.scss';

const Overview: FC<OverviewProps> = ({ overviewList }) => {
  const formatMessage = useTranslate();

  return (
    <Flex vertical gap={18} className={styles['overview-section']}>
      <ComEllipsis className={styles['title']}>{formatMessage('common.overview')}</ComEllipsis>
      <div className={styles['overview']}>
        {overviewList?.map((d: any) => {
          return (
            <ProCard
              border
              key={d.key}
              header={{
                title: formatMessage(d.label),
                customIcon: (
                  <Flex className={styles['overview-icon']} align="center" justify="center">
                    {d.icon}
                  </Flex>
                ),
              }}
              iconBg={false}
              allowHover={false}
              classNames={{
                root: styles['overview-item'],
                card: styles['overview-card'],
                header: styles['overview-card-header'],
                secondaryDescription: styles['overview-card-value-wrap'],
              }}
              description={false}
              styles={{
                secondaryDescription: {
                  lineHeight: 1,
                },
                card: {
                  background: 'var(--ui-promodal-bg-pg-color)',
                  border: '1px solid var(--ui-select-card-color)',
                },
                headerTitle: {
                  fontSize: 16,
                  fontWeight: 'bold',
                  color: 'var(--ui-text-color)',
                },
              }}
              secondaryDescription={
                <Flex className={styles['overview-value']} align="flex-end" gap={4}>
                  <ComEllipsis>{d.value}</ComEllipsis>
                  {d.unit && <span className={styles['overview-unit']}>{formatMessage(d.unit)}</span>}
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
