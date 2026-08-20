import classNames from 'classnames';
import type { FC } from 'react';
import { Segmented } from 'antd';
import { Grid, List } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from './icon-props';
import styles from './toolbar-icon.module.scss';

export interface ViewModeSegmentedProps {
  value?: string;
  onChange?: (value: string) => void;
  cardTitle: string;
  listTitle: string;
  className?: string;
}

const ViewModeSegmented: FC<ViewModeSegmentedProps> = ({ value, onChange, cardTitle, listTitle, className }) => {
  return (
    <Segmented
      className={classNames(styles['view-mode-segmented'], className)}
      size="small"
      value={value}
      onChange={(next) => onChange?.(String(next))}
      options={[
        {
          value: 'list',
          icon: (
            <span className={styles['view-mode-icon']} title={listTitle}>
              <List {...toolbarIconProps} />
            </span>
          ),
        },
        {
          value: 'card',
          icon: (
            <span className={styles['view-mode-icon']} title={cardTitle}>
              <Grid {...toolbarIconProps} />
            </span>
          ),
        },
      ]}
    />
  );
};

export default ViewModeSegmented;
