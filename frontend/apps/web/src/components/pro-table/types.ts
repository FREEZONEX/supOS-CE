import type { TableColumnsType, TableProps } from 'antd';
import type { CSSProperties } from 'react';
import type { OperationProps } from '@/components/operation-buttons/utils.tsx';

export interface ATableProps extends Omit<TableProps, 'columns'> {
  // titleIntlId 国际化key
  columns: TableColumnsType & { titleIntlId?: string; [key: string]: any };
  resizeable?: boolean;
  // 是否隐藏空白
  hiddenEmpty?: boolean;
  fixedPosition?: boolean; // 是否固定页码在底部
  showExpand?: boolean; // 是否显示展开按钮
  wrapperStyle?: CSSProperties;
  // 操作项配置
  operationOptions?: {
    title?: string | (() => string);
    width?: number | string;
    render: (record: any, index: number) => (OperationProps | null)[];
    disabled?: boolean;
  };
  // 置顶配置
  pinOptions?: {
    title?: string | (() => string);
    width?: number | string;
    disabled?: boolean;
    onClick?: (record: any) => Promise<any>;
    renderPinIcon?: (record: any) => boolean;
    auth?: string | string[];
  };
}

export interface TitlePropsType {
  width?: number;
  minWidth?: number;
  changeWidth: (width: number) => void;
}
