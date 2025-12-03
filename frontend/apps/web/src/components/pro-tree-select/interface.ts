import { type GetRef, TreeSelect, type TreeSelectProps } from 'antd';
import type { AxiosRequestConfig } from 'axios';
import type { DataNode } from 'antd/es/tree';
import type { ReactNode } from 'react';

export type TreeSelectRefType = GetRef<typeof TreeSelect>; // BaseSelectRef

// 传递的页面请求参数
export interface PageParamsType {
  pageNo: number;
  pageSize: number;
  // 总页码
  total: number;
  // 搜索
  searchValue?: string;
}

export interface QueryParamsType extends Partial<PageParamsType> {
  // root: 根 对应首次请求; node: 节点 对于loadData; page: 分页 对应分页加载;
  type?: 'root' | 'page' | 'node' | 'search';
  // 当前请求的key值 如果是跟既为0
  key?: string;
  // 当前节点数据
  currentNode?: TreeDataNodeProps;
}

export interface LoadDataOptionsType {
  reset?: boolean;
}

export interface ProTreeSelectProps extends TreeSelectProps {
  // 是否开启分页加载
  lazy?: boolean;
  // 是否启用点击加载
  loadDataEnable?: boolean;
  // 默认分页
  defaultPageSize?: number;
  // 是否显示打开节点图标，比如单层树情况，传false
  showSwitcherIcon?: boolean;
  // 分页加载逻辑
  loadMoreData?: any;
  // 请求逻辑
  api?: (
    params?: Partial<QueryParamsType>,
    config?: AxiosRequestConfig
  ) => Promise<{
    data: DataNode[];
    pageNo?: number;
    pageSize?: number;
    total?: number;
    // 回调函数
    cb?: () => void;
  }>;
  // 搜索debounce
  debounceTimeout?: number;
  // 是否开启全选功能
  selectAllAble?: boolean;
  // 左侧icon
  treeNodeIcon?: ReactNode | ((dataNode: TreeDataNodeProps, treeExpends?: any[]) => ReactNode);
  // 右侧操作项
  treeNodeExtra?: ReactNode | ((dataNode: TreeDataNodeProps) => ReactNode);
  // 数量
  treeNodeCount?: ReactNode | ((dataNode: TreeDataNodeProps) => ReactNode);
  // 是否占据整行
  treeBlockNode?: boolean;
}

export interface TreeDataNodeProps extends Omit<DataNode, 'key'> {
  // 是否有更多节点
  isLoadMoreNode?: boolean;
  // 父级key
  parentKey?: string;
  key?: string;
  children?: any[];
  [key: string]: any;
}
