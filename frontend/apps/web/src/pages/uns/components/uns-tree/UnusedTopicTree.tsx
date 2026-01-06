import { type FC, useState, useEffect, useRef, useImperativeHandle, useCallback } from 'react';
import { Folder, FolderOpen, Document, DocumentImport } from '@carbon/icons-react';
import { debounce } from 'lodash-es';
import { useTranslate } from '@/hooks';
import { App, Flex, Divider } from 'antd';
import HighlightText from './HighlightText';
import './index.scss';

import type { UnsTreeNode } from '@/pages/uns/types';
import ComTree from '@/components/com-tree';
import ProSearch from '@/components/pro-search';
import { collectChildrenIds, findParentIds } from '@/utils/uns';
import { useTreeStore } from '../../store/treeStore';
import { useBaseStore } from '@/stores/base';
import { ClearOutlined } from '@ant-design/icons';
import { clearExternalTreeData } from '@/apis/inter-api/external.ts';

export interface UnusedTopicTreeProps {
  treeData: any[];
  currentNode: UnsTreeNode;
  initTreeData: (data: any, cb?: () => void) => void;
  setTreeMap?: (params: any) => void;
  changeCurrentPath: (node?: any) => void;
  treeHeight?: number;
  unsTreeRef?: any;
  treeType: string;
  unusedTopicBreadcrumbList: UnsTreeNode[];
}

export interface RightKeyMenuItem {
  auth: string;
  menuTitle: string;
  click: () => void;
  height: number;
}

const UnusedTopicTree: FC<UnusedTopicTreeProps> = ({
  treeData,
  currentNode,
  initTreeData,
  changeCurrentPath,
  setTreeMap,
  treeHeight = 0,
  unsTreeRef,
  treeType,
  unusedTopicBreadcrumbList,
}) => {
  const [expandedArr, setExpandedArr] = useState<any>([]); //展开的树节点
  const [searchQuery, setSearchQuery] = useState(''); //搜索参数
  const { message } = App.useApp();
  const formatMessage = useTranslate();

  // 创建一个 ref 来引用 tree 元素
  const treeRef: any = useRef(null);
  const selectRef = useRef<any>(null);
  const isComposingRef = useRef(false); // 拼音输入中..

  const { operationFns, setOperationFns } = useTreeStore((state) => ({
    operationFns: state.operationFns,
    setOperationFns: state.setOperationFns,
  }));
  const { systemInfo } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));

  useImperativeHandle(unsTreeRef, () => ({
    expandedArr: expandedArr,
    setExpandedArr: setExpandedArr,
    scrollTo: scrollTreeNode,
    searchQuery: searchQuery,
  }));

  useEffect(() => {
    //处理拼音输入法
    const inputElement = selectRef.current;

    if (inputElement) {
      const handleCompositionStart = () => {
        isComposingRef.current = true;
      };

      const handleCompositionEnd = (e: any) => {
        isComposingRef.current = false;
        const value = e.target.value;
        if (value) onSearchChange(value);
      };

      inputElement.addEventListener('compositionstart', handleCompositionStart);
      inputElement.addEventListener('compositionend', handleCompositionEnd);

      return () => {
        inputElement.removeEventListener('compositionstart', handleCompositionStart);
        inputElement.removeEventListener('compositionend', handleCompositionEnd);
      };
    }
  }, []);

  useEffect(() => {
    setSearchQuery('');
  }, [treeType]);

  useEffect(() => {
    if (searchQuery && treeType === 'treemap') {
      if (searchQuery.includes('/')) {
        setExpandedArr(collectChildrenIds(treeData, ''));
      } else {
        setExpandedArr(findParentIds(searchQuery, treeData));
      }
    }
  }, [treeData]);

  //滚动到目标树节点
  const scrollTreeNode = (key: any) => {
    treeRef.current.scrollTo({ key, align: 'top' });
  };

  //处理树节点展开收起逻辑
  const onExpand = (expandedKeys: any, { expanded, node }: any) => {
    const newArr = expandedKeys;
    if (expanded) {
      newArr.push(node.value);
    } else {
      const index = expandedKeys.findIndex((id: any) => id === node.value);
      newArr.splice(index, 1);
    }
    setExpandedArr([...new Set(newArr)]);
  };

  const debouncedInitTreeData = useCallback(
    debounce((value) => {
      initTreeData({ reset: true, query: value });
    }, 500),
    [initTreeData] // 仅当 initTreeData 发生变化时更新 debounce 函数
  );

  //查询树节点
  const onSearchChange = (val: string) => {
    return debouncedInitTreeData(val);
  };
  const refreshUnusedTopicTree = (path: string) => {
    if (path === currentNode?.id) changeCurrentPath();
    initTreeData({ reset: false, query: searchQueryRef.current });
  };

  useEffect(() => {
    setOperationFns({ refreshUnusedTopicTree });
  }, []);

  const searchQueryRef = useRef(searchQuery);
  useEffect(() => {
    searchQueryRef.current = searchQuery;
  }, [searchQuery]);

  return (
    <div className="treeWrap">
      <Flex gap="8px" className="treeSearchBox">
        <ProSearch
          ref={selectRef}
          closeButtonLabelText={formatMessage('common.clearSearchInput')}
          id="search-playground-1"
          placeholder={formatMessage('uns.inputText')}
          role="searchbox"
          size="sm"
          type="text"
          onChange={(e) => {
            const val = e.target.value || '';
            setSearchQuery(val);
            if (isComposingRef.current) return;
            onSearchChange(val);
          }}
          value={searchQuery}
          style={{ borderRadius: '3px', flex: 1 }}
          onKeyDown={(e) => {
            if (e.keyCode === 13) {
              initTreeData({ reset: true, query: searchQuery });
            }
          }}
        />
        <div className="treeOperateIconWrap">
          <span title={formatMessage('common.clear')}>
            <ClearOutlined
              onClick={() => {
                setSearchQuery('');
                clearExternalTreeData().then(() => {
                  message.success(formatMessage('common.clearSuccessful'));
                  initTreeData({ reset: true });
                });
              }}
            />
          </span>
        </div>
      </Flex>
      <Divider style={{ borderColor: '#e0e0e0', margin: '16px 0 10px 0' }} />
      <ComTree
        ref={treeRef}
        height={treeHeight}
        selectedKeys={[currentNode.id as string]}
        treeData={treeData}
        fieldNames={{ title: 'pathName', key: 'id' }}
        blockNode
        expandedKeys={expandedArr}
        onExpand={onExpand}
        notDataContent={
          treeData.length === 0 ? <div style={{ textAlign: 'center' }}>{formatMessage('uns.noData')}</div> : null
        }
        titleRender={(item: any) => {
          return (
            <div
              className={`customTreeNode`}
              onClick={() => {
                changeCurrentPath(item);
                setTreeMap?.(false);
              }}
              style={
                treeType === 'treemap' &&
                currentNode.id &&
                unusedTopicBreadcrumbList
                  .slice(0, -1)
                  .map((e) => e.id)
                  .includes(item.id)
                  ? {
                      color: 'var(--supos-theme-color)',
                    }
                  : {}
              }
            >
              {item.type === 0 &&
                (expandedArr.includes(item.id) && item.children ? (
                  <FolderOpen style={{ flexShrink: 0, marginRight: '5px' }} />
                ) : (
                  <Folder style={{ flexShrink: 0, marginRight: '5px' }} />
                ))}
              {item.type === 2 && <Document style={{ flexShrink: 0, marginRight: '5px' }} />}
              <Flex align="center" gap={8} style={{ flex: 1, overflow: 'hidden' }}>
                <div style={{ flex: 1, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  <HighlightText needle={searchQuery} haystack={item.pathName || item.name} />
                  {item.type === 0 && (
                    <span style={{ color: 'var(--supos-text-color)', fontSize: '12px', opacity: 0.5 }}>
                      ({item.countChildren})
                    </span>
                  )}
                </div>
                {!systemInfo?.useAliasPathAsTopic && item.type === 2 && (
                  <Flex
                    className="treeNodeIconWrap"
                    align="center"
                    style={{ flexShrink: 0, padding: '0 4px' }}
                    title={formatMessage('uns.topicToFile')}
                    onClick={(e) => {
                      e.stopPropagation();
                      operationFns?.setOptionOpen?.('topicToFile', item);
                    }}
                  >
                    <DocumentImport />
                  </Flex>
                )}
              </Flex>
            </div>
          );
        }}
      />
    </div>
  );
};
export default UnusedTopicTree;
