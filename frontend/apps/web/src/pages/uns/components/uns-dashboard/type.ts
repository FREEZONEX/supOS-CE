import type { ReactNode } from 'react';

export interface OverviewProps {
  overviewList: OverviewListProps[];
}

export interface OverviewListProps {
  key: string;
  label: string;
  icon: ReactNode;
  value: number;
  unit?: string;
}
