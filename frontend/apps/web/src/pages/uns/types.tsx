import type { DataNodeProps } from '@/components/pro-tree';
import type { Key } from 'react';

export interface UnsTreeNode extends Omit<DataNodeProps, 'children'> {
  key: Key;
  id?: Key;
  parentId?: Key;
  path?: string;
  parentPath?: string;
  /*
   * @description type => pathType
   * 0 文件夹 2 文件 1 模板 7 标签
   * */
  // type?: number | null;
  pathType?: number | null;
  name?: string;
  alias?: string;
  parentAlias?: string;
  children?: UnsTreeNode[];
  // 子孙文件的个数
  countChildren?: number;
  // 是否有子集
  hasChildren?: boolean;
  dataType?: number;
  parentDataType?: number;
  [key: string]: any;
}

export interface FieldItem {
  name: string;
  type: string;
  displayName?: string;
  remark?: string;
  unique?: boolean;
  index?: number | string;
  systemField?: boolean;
  maxLen?: number;
  unit?: string;
  upperLimit?: number;
  lowerLimit?: number;
  decimal?: number;
}

interface InitTreeDataParamsType {
  reset?: boolean;
  query?: string;
  type?: number;
  [key: string]: any;
}

export type InitTreeDataFnType = (params: InitTreeDataParamsType, cb?: () => void) => void;

export interface SelectTreeNode extends Omit<UnsTreeNode, 'children'> {
  fields: FieldItem[];
  children: SelectTreeNode[];
}
