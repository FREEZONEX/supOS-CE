import classNames from 'classnames';
import type { FC, MouseEvent } from 'react';
import { Segmented } from 'antd';
import { Code, TableSplit } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import styles from '@/components/lucide-icon/toolbar-icon.module.scss';

export type PayloadViewMode = 'table' | 'code';

export interface PayloadViewSegmentedProps {
  value: PayloadViewMode;
  onChange: (value: PayloadViewMode) => void;
  tableTitle: string;
  codeTitle: string;
  className?: string;
}

const PayloadViewSegmented: FC<PayloadViewSegmentedProps> = ({ value, onChange, tableTitle, codeTitle, className }) => {
  const stopPropagation = (event: MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();
  };

  return (
    <div className={classNames('payload-view-segmented-wrap', className)} onClick={stopPropagation}>
      <Segmented
        className={styles['view-mode-segmented']}
        size="small"
        value={value}
        onChange={(next) => onChange(next as PayloadViewMode)}
        options={[
          {
            value: 'table',
            icon: (
              <span className={styles['view-mode-icon']} title={tableTitle}>
                <TableSplit {...toolbarIconProps} />
              </span>
            ),
          },
          {
            value: 'code',
            icon: (
              <span className={styles['view-mode-icon']} title={codeTitle}>
                <Code {...toolbarIconProps} />
              </span>
            ),
          },
        ]}
      />
    </div>
  );
};

export default PayloadViewSegmented;
