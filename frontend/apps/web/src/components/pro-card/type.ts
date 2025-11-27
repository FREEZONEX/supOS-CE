export type SemanticDOM =
  | 'root'
  | 'card'
  | 'statusHeader'
  | 'statusInfo'
  | 'statusTag'
  | 'header'
  | 'headerTitle'
  | 'secondaryDescription'
  | 'actions';

import type { CSSProperties, ReactNode } from 'react';
import type { OperationProps } from '@/components/operation-buttons/utils.tsx';

export interface ProCardProps {
  loading?: boolean;
  // 影响hover的样式
  allowHover?: boolean;
  styles?: Partial<Record<SemanticDOM, CSSProperties>>;
  classNames?: Partial<Record<SemanticDOM, string>>;
  value?: boolean;
  // 是否有icon背景色
  iconBg?: boolean;
  onChange?: (e: any) => void;
  statusHeader?: {
    allowCheck?: boolean;
    statusInfo?: {
      label: string;
      color: string;
      title: string;
    };
    statusTag?: ReactNode;
    pinOptions?: {
      disabled?: boolean;
      onClick?: (record: any) => Promise<any>;
      renderPinIcon?: (record: any) => boolean;
      auth?: string | string[];
    };
  };
  item?: any;
  onClick?: (item?: any) => void;
  header?: {
    customIcon?: ReactNode;
    iconSrc?: string;
    defaultIconUrl?: string;
    title?: ReactNode;
    titleDescription?: ReactNode;
    onClick?: (item?: any) => void;
  };
  description?: false | string | { content?: string; rows?: number; empty?: string };
  secondaryDescription?: ReactNode;
  actions?: ((item: any) => OperationProps[]) | OperationProps[];
  actionConfig?: {
    num?: number;
  };
  border?: boolean;
}
