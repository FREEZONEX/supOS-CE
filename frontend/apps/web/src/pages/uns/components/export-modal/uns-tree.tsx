import ProTree from '@/components/pro-tree';
import { Checkbox, Divider, Flex } from 'antd';
import {
  ChartLine,
  Renew,
  SendAlt,
  Document,
  FolderOpen,
  Folder,
  WatsonHealth3DCurveAutoColon,
} from '@carbon/icons-react';
import ComButton from '@/components/com-button';
import { useTranslate } from '@/hooks';
import { useTreeStore } from './treeStore.tsx';
import { type FC, type Key, memo, useCallback, useEffect, useRef } from 'react';
import { useBaseStore } from '@/stores/base';
import type { UnsTreeNode } from '../../types.tsx';
import { exportExcel } from '@/apis/inter-api/uns';
import { processedCheckedKeys } from '@/pages/uns/store/utils.ts';
import { getParamsForArray } from '@/utils';

const TreeNodeIcon = memo(({ dataNode }: { dataNode: UnsTreeNode }) => {
  const { expandedKeys } = useTreeStore((state) => ({
    expandedKeys: state.expandedKeys,
  }));
  const {
    systemInfo: { enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));

  let Dom;
  let color;
  const parentDataType = enableAutoCategorization
    ? dataNode.type === 0
      ? dataNode.dataType
      : dataNode.parentDataType
    : 0;

  switch (parentDataType) {
    case 1:
      Dom = Document;
      color = '#D2A106';
      break;
    case 2:
      Dom = SendAlt;
      color = '#94C518';
      break;
    case 3:
      Dom = ChartLine;
      color = '#1D77FE';
      break;
    default:
      break;
  }
  const commonStyle = { flexShrink: 0, marginRight: '5px' };

  if (dataNode.pathType === 0) {
    return (
      <Flex align="center">
        {enableAutoCategorization ? (
          dataNode.dataType && Dom ? (
            <Dom style={{ ...commonStyle, color: color }} />
          ) : (
            <WatsonHealth3DCurveAutoColon style={commonStyle} />
          )
        ) : expandedKeys.includes(dataNode.key) && dataNode.hasChildren ? (
          <FolderOpen style={commonStyle} />
        ) : (
          <Folder style={commonStyle} />
        )}
      </Flex>
    );
  } else if (dataNode.pathType === 2) {
    return (
      <Flex align="center">
        {Dom ? <Dom style={{ ...commonStyle, color: color }} /> : <Document style={commonStyle} />}
      </Flex>
    );
  }
  return null;
});

export const UnsTree: FC<{ open: boolean }> = ({ open }) => {
  const formatMessage = useTranslate();
  const treeRef = useRef<any>(null);

  const {
    loadData,
    treeData,
    setCheckedKeys,
    checkedKeys,
    loading,
    loadedKeys,
    loadingKeys,
    nodePaginationState,
    expandedKeys,
    setExpandedKeys,
    setLoadedKeys,
    setLazyTree,
    setScrollTreeNode,
    allChecked,
    setAllChecked,
    setLoading,
    setJsonData,
    setParams,
    setSmallFile,
  } = useTreeStore((state) => ({
    loadData: state.loadData,
    treeData: state.treeData,
    setCheckedKeys: state.setCheckedKeys,
    checkedKeys: state.checkedKeys,
    loading: state.loading,
    loadedKeys: state.loadedKeys,
    setLoadedKeys: state.setLoadedKeys,
    loadingKeys: state.loadingKeys,
    nodePaginationState: state.nodePaginationState,
    expandedKeys: state.expandedKeys,
    setExpandedKeys: state.setExpandedKeys,
    setLazyTree: state.setLazyTree,
    setScrollTreeNode: state.setScrollTreeNode,
    allChecked: state.allChecked,
    setAllChecked: state.setAllChecked,
    setLoading: state.setLoading,
    setJsonData: state.setJsonData,
    setParams: state.setParams,
    setSmallFile: state.setSmallFile,
  }));

  const {
    lazyTree,
    systemInfo: { enableAutoCategorization },
  } = useBaseStore((state) => ({
    lazyTree: state.systemInfo?.lazyTree || false,
    systemInfo: state.systemInfo,
  }));

  useEffect(() => {
    setLazyTree(lazyTree || false);
    if (open) {
      loadData({ reset: true });
    }
  }, [open, lazyTree, loadData]);

  const onLoadData = async (node: any) => {
    const _node = { ...node };
    return loadData({
      key: _node.key,
      parentInfo: _node,
    });
  };
  const handleRenderLoadMoreNode = (moreNodeData: any) => {
    const { parentKey: nodeKey, currentPage } = moreNodeData;
    if (loadingKeys.has(nodeKey)) {
      console.log(`正在loading ${nodeKey}`);
      return;
    }
    const state = nodePaginationState[nodeKey];
    if (state && state.hasMore && !state.isLoading) {
      // 处理根节点和子节点的加载更多
      if (currentPage === state.currentPage) {
        // 以防重复请求2次
        loadData({
          key: nodeKey,
          page: state.currentPage + 1,
          parentInfo: moreNodeData?.parentInfo,
        });
      }
    }
  };

  //滚动到目标树节点
  const scrollTreeNode = useCallback((id: Key) => {
    setTimeout(() => {
      if (treeRef.current) treeRef.current.scrollTo?.({ key: id, align: 'top' });
    }, 500);
  }, []);

  useEffect(() => {
    setScrollTreeNode(scrollTreeNode);
  }, [scrollTreeNode]);

  return (
    <>
      {/*<Input*/}
      {/*  style={{ marginBottom: 8 }}*/}
      {/*  disabled={loading}*/}
      {/*  placeholder={formatMessage('common.search')}*/}
      {/*  prefix={<Search />}*/}
      {/*/>*/}
      <div style={{ flex: 1, overflow: 'hidden' }}>
        <ProTree
          disabled={allChecked}
          ref={treeRef}
          selectable={false}
          checkedKeys={checkedKeys}
          onCheck={(checkedKeysValue) => {
            setCheckedKeys(checkedKeysValue as Key[]);
          }}
          checkable
          wrapperStyle={{
            border: '1px solid #c6c6c6',
            borderRadius: 4,
            padding: 4,
            flex: 1,
          }}
          height={0}
          treeData={treeData}
          treeNodeIcon={(dataNode) => <TreeNodeIcon dataNode={dataNode} />}
          treeNodeCount={(dataNode) => {
            return (
              dataNode.pathType === 0 && (
                <span
                  style={{
                    color: enableAutoCategorization ? '#161616' : 'var(--supos-text-color)',
                    fontSize: '12px',
                    opacity: 0.5,
                  }}
                >
                  ({dataNode.countChildren})
                </span>
              )
            );
          }}
          renderTitleStyle={(dataNode) => {
            const bgColor =
              dataNode.type === 0 && dataNode.dataType && enableAutoCategorization
                ? dataNode.dataType === 1
                  ? '#FCF4D6'
                  : dataNode.dataType === 2
                    ? '#F0FBD2'
                    : '#E8F1FF'
                : '';
            return bgColor
              ? {
                  height: '26px',
                  backgroundColor: bgColor,
                  borderRadius: '3px',
                  paddingRight: '8px',
                  color: '#161616',
                  paddingLeft: 10,
                }
              : {};
          }}
          header={
            <>
              <Flex justify="space-between" align="center" style={{ marginTop: 4, padding: '0 4px' }}>
                <Checkbox
                  checked={allChecked}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setCheckedKeys([]);
                    }
                    setAllChecked(e.target.checked);
                  }}
                >
                  {formatMessage('common.selectAll')}
                </Checkbox>
                <Renew
                  style={{ cursor: 'pointer' }}
                  onClick={() => {
                    loadData({ reset: true });
                  }}
                />
              </Flex>
              <Divider style={{ borderColor: '#c6c6c6', margin: '8px 0' }} />
            </>
          }
          loadData={onLoadData}
          loading={loading}
          loadMoreData={handleRenderLoadMoreNode}
          loadedKeys={loadedKeys}
          expandedKeys={expandedKeys}
          onExpand={(expandedKeys) => {
            setExpandedKeys(expandedKeys);
            setLoadedKeys(expandedKeys);
          }}
          lazy={lazyTree}
        />
      </div>

      <Flex justify="end" gap={8} style={{ marginTop: 16 }}>
        <ComButton
          onClick={() => {
            setCheckedKeys([]);
            setAllChecked(false);
          }}
        >
          {formatMessage('common.reset')}
        </ComButton>
        <ComButton
          loading={loading}
          disabled={checkedKeys?.length === 0 && !allChecked}
          type="primary"
          onClick={() => {
            setLoading(true);
            let params: any = {
              fileType: 'json',
              checkSmallFile: true,
            };
            if (allChecked) {
              params['exportType'] = 'ALL';
            } else {
              // 根据checkedKeys匹配节点信息
              const matchedNodes = processedCheckedKeys({
                checkedKeys,
                strategy: 'SHOW_PARENT',
                treeData, // 添加treeData参数
              });
              params = {
                ...params,
                ...getParamsForArray(matchedNodes as any[], 'pathType', {
                  groups: {
                    0: 'folders',
                    2: 'files',
                  },
                  extract: 'id',
                }),
                checkSmallFile: true,
              };
            }
            return exportExcel(params)
              .then((info) => {
                if (info?.code) {
                  //
                  if (info.code === 200) {
                    setSmallFile(false);
                  } else {
                    setJsonData(undefined);
                    setSmallFile(true);
                  }
                } else {
                  // code不存在，就直接是数据
                  try {
                    setJsonData(JSON.stringify(info, null, 2));
                  } catch (e) {
                    console.log(e);
                    setJsonData(undefined);
                  }
                  setSmallFile(true);
                }
              })
              .finally(() => {
                setParams({
                  ...params,
                  checkSmallFile: undefined,
                });
                setLoading(false);
              });
          }}
        >
          {formatMessage('appSpace.newgenerate')}
        </ComButton>
      </Flex>
    </>
  );
};
