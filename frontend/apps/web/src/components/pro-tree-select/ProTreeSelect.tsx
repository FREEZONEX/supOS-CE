import { Checkbox, Divider, Flex, TreeSelect, Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';
import { forwardRef, type ReactNode, useEffect, useMemo, useRef } from 'react';
import useTranslate from '../../hooks/useTranslate.ts';
import cx from 'classnames';
import './index.scss';
import { useImmer } from 'use-immer';
import usePropsValue from '../../hooks/usePropsValue.ts';
import HighlightText from '@/components/pro-tree/HighlightText.tsx';
import { debounce } from 'lodash-es';
import { produce } from 'immer';
import { Renew } from '@carbon/icons-react';
import type {
  TreeDataNodeProps,
  LoadDataOptionsType,
  PageParamsType,
  ProTreeSelectProps,
  QueryParamsType,
  TreeSelectRefType,
} from './interface.ts';
import {
  appendTreeData,
  createLoadMoreNode,
  findMissingValuesOptimized,
  findTreeNode,
  hasMoreData,
  RootId,
  SelectAllId,
} from './utils.ts';

const ProTreeSelect = forwardRef<TreeSelectRefType, ProTreeSelectProps>(
  (
    {
      placeholder,
      classNames,
      showSwitcherIcon,
      treeBlockNode = true,
      selectAllAble = true,
      labelInValue,
      popupRender,
      onChange,
      value,
      treeData,
      onSelect,
      lazy,
      treeNodeIcon,
      treeNodeExtra,
      treeNodeCount,
      treeTitleRender,
      api,
      defaultPageSize = 100,
      debounceTimeout = 500,
      onSearch,
      loadDataEnable,
      treeExpandedKeys,
      onTreeExpand,
      // treeLoadedKeys,
      ...restProps
    },
    ref
  ) => {
    const _labelInValue = selectAllAble ? true : labelInValue;
    const searchValueRef = useRef<string | undefined>(undefined);
    // 翻译，如果组件提出去，翻译要外面传入
    const formatMessage = useTranslate();
    // 是否开启了多选功能
    const multiple = restProps.treeCheckable || restProps?.multiple;
    // 全选
    const [allChecked, setAllChecked] = useImmer(false);
    // select的值
    const [_value, setValue] = usePropsValue({
      value,
    });
    // expandedKeys loadedKeys
    const [_treeExpandedKeys, setTreeExpandedKeys] = usePropsValue({
      value: treeExpandedKeys,
      onChange: onTreeExpand,
    });
    // tree数据，代理一层
    const [_treeData, setTreeData] = usePropsValue<TreeDataNodeProps[] | undefined>({
      value: treeData as any,
    });

    // 正在loading的key
    const [loadingKeys, setLoadingKeys] = useImmer(new Set());
    // abortController实例
    const [abortControllers, setAbortControllers] = useImmer(new Map());
    // 分页
    const [nodePaginationState, setNodePaginationState] = useImmer<{
      [key: string]: {
        currentPage: number;
        hasMore: boolean;
        isLoading: boolean;
      };
    }>({});
    // 触底请求逻辑
    const loadMoreData = useMemo(() => {
      return debounce((moreNodeData: any) => {
        const { parentKey: nodeKey, currentPage } = moreNodeData;
        if (loadingKeys.has(nodeKey)) {
          console.warn(`正在loading ${nodeKey}`);
          return;
        }
        const state = nodePaginationState[nodeKey];
        if (state && state.hasMore && !state.isLoading) {
          // 处理根节点和子节点的加载更多
          if (currentPage === state.currentPage) {
            loadData({
              key: nodeKey,
              currentNode: moreNodeData,
              pageNo: state.currentPage + 1,
              searchValue: searchValueRef?.current,
            });
          }
        }
      }, 300);
    }, [nodePaginationState, loadingKeys]);

    const _treeTitleRender = (node: any) => {
      const title = node[restProps?.fieldNames?.label || 'title'] as ReactNode;
      if (node.isLoadMoreNode && loadMoreData && lazy) {
        loadMoreData?.(node);
        return node.title;
      }
      const _title = treeTitleRender ? treeTitleRender?.(node) : title;
      const Icon = typeof treeNodeIcon === 'function' ? treeNodeIcon(node, _treeExpandedKeys) : treeNodeIcon;
      const Extra = typeof treeNodeExtra === 'function' ? treeNodeExtra(node) : treeNodeExtra;
      const Count = typeof treeNodeCount === 'function' ? treeNodeCount(node) : treeNodeCount;
      return (
        <Flex align="center" gap={8}>
          {loadDataEnable && [...loadingKeys]?.includes(node.key) && (
            <Spin indicator={<LoadingOutlined spin />} size="small" />
          )}
          {Icon && <div className="pro-tree-select-node-icon">{Icon}</div>}
          <Flex style={{ flex: 1, overflow: 'hidden' }} align="center" gap={8}>
            <div style={{ flex: 1, overflow: 'hidden' }}>
              <div className="pro-tree-select-node-title">
                <HighlightText needle={searchValueRef.current} haystack={_title} />
                {Count}
              </div>
            </div>
            {Extra && (
              <div className="pro-tree-select-node-extra" style={{ flexShrink: 0 }}>
                {Extra}
              </div>
            )}
          </Flex>
        </Flex>
      );
    };

    // 搜索
    const debounceFetcher = useMemo(() => {
      const loadOptions = (value: string) => {
        searchValueRef.current = value;
        loadData(
          {
            type: 'search',
            key: RootId,
            pageNo: 1,
            searchValue: value,
            currentNode: {
              key: RootId,
            },
          },
          {
            reset: true,
          }
        );
        onSearch?.(value);
      };
      return debounce(loadOptions, debounceTimeout);
    }, [debounceTimeout]);

    // 首次请求
    useEffect(() => {
      loadData({
        type: 'root',
        key: RootId,
        pageNo: 1,
        currentNode: {
          key: RootId,
        },
      });
    }, []);

    // 请求方法集合
    const loadData = async (params: QueryParamsType, options?: LoadDataOptionsType) => {
      if (!api) return;
      const { key = RootId, pageNo = 1 } = params || {};
      const { reset } = options || {};
      const pageSize = defaultPageSize;
      if (lazy) {
        if (loadingKeys.has(key)) {
          console.warn(`节点正在请求：${key};  params: ${params}`);
          return;
        }
        if (reset) {
          // 清除请求
          abortControllers.forEach((controller) => {
            controller.abort();
          });
          setAbortControllers((pre) => {
            pre.clear();
          });
          // 清空展开
          setTreeExpandedKeys([]);
        }
        setLoadingKeys((pre) => {
          pre.add(key);
        });
        setNodePaginationState((pre) => {
          pre[key as string] = { ...pre[key as string], isLoading: true };
        });
        // 懒加载请求逻辑
        // 创建AbortController用于取消请求
        const controller = new AbortController();
        const keyStr = key.toString();
        setAbortControllers((pre) => {
          pre.set(keyStr, controller);
        });
        try {
          const { data, cb, ...restResponse } = await api(
            {
              ...params,
              pageSize,
            },
            {
              signal: controller.signal,
            }
          );
          const hasMore = hasMoreData(restResponse as PageParamsType);
          if (key === RootId) {
            if (pageNo === 1) {
              // 跟节点首次请求或搜索
              if (hasMore) {
                setTreeData([
                  ...data,
                  createLoadMoreNode(key, 1, formatMessage('common.loadMore'), restProps?.fieldNames),
                ]);
                setNodePaginationState((pre) => {
                  pre[key] = { currentPage: pageNo, hasMore, isLoading: false };
                });
              } else {
                setTreeData(data);
              }
            } else {
              // 跟节点触底加载
              setTreeData((pre: TreeDataNodeProps[]) => {
                const filteredData = pre.filter((node) => !node.isLoadMoreNode);
                return hasMore
                  ? [
                      ...filteredData,
                      ...data,
                      createLoadMoreNode(key, pageNo, formatMessage('common.loadMore'), restProps?.fieldNames),
                    ]
                  : [...filteredData, ...data];
              });
              setNodePaginationState((pre) => {
                pre[key] = { currentPage: pageNo, hasMore, isLoading: false };
              });
            }
          } else {
            // 子孙节点加载逻辑
            const newData: any = hasMore
              ? [...data, createLoadMoreNode(key, pageNo, formatMessage('common.loadMore'), restProps?.fieldNames)]
              : data;
            setTreeData(
              produce((pre: TreeDataNodeProps[]) => {
                appendTreeData(pre, key, newData);
              })
            );
            setNodePaginationState((pre) => {
              pre[key as string] = { currentPage: pageNo, hasMore, isLoading: false };
            });
          }
          cb?.();
        } catch (e) {
          console.log(e);
          setNodePaginationState((pre) => {
            pre[key as string] = { ...pre[key as string], isLoading: false };
          });
        } finally {
          setLoadingKeys((pre) => {
            pre.delete(key);
          });
        }
      } else {
        // 首次请求
        api?.({ ...params, pageSize })?.then(({ data, cb }) => {
          setTreeData(data);
          cb?.();
        });
      }
    };

    const onLoadData = async (node: any) => {
      const _node = { ...node };
      return loadData(
        {
          // 必然会有，读取的是value值
          key: _node.key,
          // 当前点击的节点信息
          currentNode: _node,
          pageNo: 1,
        },
        {
          reset: true,
        }
      );
    };

    const nodeMap = useMemo(() => {
      if (_labelInValue) {
        const map = new Map();
        const buildMap = (nodes: TreeDataNodeProps) => {
          nodes.forEach((node: TreeDataNodeProps) => {
            map.set(node.key || node.id, node); // 使用 Map 的 set 方法
            if (node.children) buildMap(node.children);
          });
        };
        buildMap(_treeData || []);
        return map;
      }
      return new Map();
    }, [_treeData, _labelInValue]);

    return (
      <TreeSelect
        maxCount={allChecked ? 0 : undefined}
        classNames={{
          ...classNames,
          popup: {
            ...classNames?.popup,
            root: cx(
              classNames?.popup?.root,
              {
                'pro-tree-select-expend': !showSwitcherIcon,
                'pro-tree-select-node-block': treeBlockNode,
              },
              'pro-tree-select'
            ),
          },
        }}
        placeholder={placeholder || formatMessage('common.select')}
        popupRender={(originNode) => {
          return (
            <div style={{ background: 'var(--supos-bg-color)', padding: '8px 8px 0' }}>
              {selectAllAble && multiple && (
                <>
                  <Flex justify="space-between" align="center">
                    <Checkbox
                      checked={allChecked}
                      onChange={(e) => {
                        const selectAllOption = {
                          value: SelectAllId,
                          label: formatMessage('common.selectAll'),
                        };

                        let selectedValue;

                        if (e.target.checked) {
                          if (_labelInValue) {
                            selectedValue = multiple ? [selectAllOption] : selectAllOption;
                          } else {
                            selectedValue = multiple ? [SelectAllId] : SelectAllId;
                          }
                        } else {
                          selectedValue = multiple ? [] : undefined;
                        }
                        onChange?.(selectedValue, [], {} as any);
                        setValue(selectedValue);
                        setAllChecked(e.target.checked);
                      }}
                    >
                      {formatMessage('common.selectAll')}
                    </Checkbox>
                    <Renew
                      style={{ cursor: 'pointer' }}
                      onClick={() => {
                        loadData(
                          {
                            type: 'root',
                            key: RootId,
                            pageNo: 1,
                            currentNode: {
                              key: RootId,
                            },
                          },
                          {
                            reset: true,
                          }
                        );
                      }}
                    />
                  </Flex>
                  <Divider style={{ borderColor: '#c6c6c6', margin: '8px 0' }} />
                </>
              )}
              {popupRender ? popupRender(originNode) : originNode}
            </div>
          );
        }}
        value={_value}
        treeData={_treeData}
        onChange={(v, labelList, extra) => {
          if ((extra?.checked === false && extra?.triggerValue === SelectAllId) || !extra?.triggerValue) {
            setAllChecked(false);
          }
          const _v = selectAllAble
            ? v?.length > 0
              ? v?.filter((f: any) => (_labelInValue ? f?.value !== SelectAllId : f !== SelectAllId))
              : v
            : v;
          const nodes = _labelInValue
            ? _v
                .map((value: any) => {
                  const node = nodeMap.get(_labelInValue ? value?.value : value) || {};
                  return {
                    ...node,
                    ...value,
                  };
                })
                .filter(Boolean)
            : _v;
          console.log(nodeMap, _v);
          onChange?.(
            nodes,
            labelList?.filter((f) => f),
            extra
          );
          setValue(nodes);
        }}
        ref={ref}
        onSelect={(value, option) => {
          setAllChecked(false);
          onSelect?.(value, option);
        }}
        treeTitleRender={_treeTitleRender}
        labelInValue={_labelInValue}
        onSearch={debounceFetcher}
        autoClearSearchValue={false}
        allowClear
        filterTreeNode={false}
        treeNodeFilterProp=""
        treeLoadedKeys={_treeExpandedKeys}
        treeExpandedKeys={_treeExpandedKeys}
        onTreeExpand={(keys) => {
          setTreeExpandedKeys((pre: string[]) => {
            if (lazy && loadDataEnable) {
              const key = findMissingValuesOptimized(keys, pre)?.[0];
              if (key) {
                const node = findTreeNode(_treeData, key);
                if (!node?.children?.length) {
                  // 模拟loadData，解决搜索不触发loadData
                  onLoadData(node);
                }
              }
            }
            // 找出当前点击的值
            return keys;
          });
        }}
        {...restProps}
        loadData={undefined}
      />
    );
  }
);

/*
 * filterTreeNode={false} treeNodeFilterProp=""  不提供自带的搜索，走的是接口搜索
 * loadData=undefined 通过expand来模拟触发loadData，因为搜索情况无法触发
 * */

export default ProTreeSelect;
