import { type FC, useEffect, useRef, useState } from 'react';
import { TreeView as TreeViewIcon } from '@carbon/icons-react';
import { App, Drawer, Splitter } from 'antd';
import { useUnsTreeMapContext } from '@/UnsTreeMapContext.tsx';
import { useMediaSize, useTranslate } from '@/hooks';
import { Tree, UnusedTopicTree } from '@/pages/uns/components';
import TreeNodeExtra from '@/pages/uns/components/uns-tree/TreeNodeExtra.tsx';
import { ChevronDown, ChevronUp } from '@carbon/icons-react';
import { getTreeStoreSnapshot, useTreeStore, useTreeStoreRef } from '@/pages/uns/store/treeStore.tsx';
import { getParentNodes } from '@/pages/uns/store/utils.ts';
import type { InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import { getExternalTreeData } from '@/apis/inter-api/external.ts';
import ComLeft from '@/components/com-layout/ComLeft';
import Loading from '@/components/loading';
import { useWebSocket } from 'ahooks';
import { useBaseStore } from '@/stores/base';

const panelCloseSize = 48;
const panelOpenSize = 500;
const initNode = {
  key: '',
  id: '',
  pathType: null,
};

const LeftDom: FC<{
  changeCurrentPath: any;
  handleDelete: any;
  currentUnusedTopicNode: any;
  setCurrentUnusedTopicNode: any;
  unusedTopicBreadcrumbList: any;
  setUnusedTopicBreadcrumbList: any;
  changeCurrentUnusedTopicPath: any;
}> = ({
  changeCurrentPath,
  handleDelete,
  currentUnusedTopicNode,
  setCurrentUnusedTopicNode,
  unusedTopicBreadcrumbList,
  setUnusedTopicBreadcrumbList,
  changeCurrentUnusedTopicPath,
}) => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  // 针对Unused订阅的topic进行实时数据展示
  const [unusedTopicTreeData, setUnusedTopicTreeData] = useState<UnsTreeNode[]>([]); // unusedTopic订阅
  const {
    systemInfo: { enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));

  const { treeType, operationFns, selectedNode } = useTreeStore((state) => ({
    treeType: state.treeType,
    operationFns: state.operationFns,
    selectedNode: state.selectedNode,
  }));
  const stateRef = useTreeStoreRef();
  const { setTreeMap, setPasteNode, loadData, setCurrentTreeMapType } = getTreeStoreSnapshot(stateRef, (state) => ({
    setTreeMap: state.setTreeMap,
    setPasteNode: state.setPasteNode,
    loadData: state.loadData,
    setCurrentTreeMapType: state.setCurrentTreeMapType,
  }));
  const [unusedTopicLoading, setUnusedTopicLoading] = useState(false);
  const [unusedTopicPanelSize, setUnusedTopicPanelSize] = useState([0, panelCloseSize]);

  const splitterWrapRef = useRef<any>(null);
  const unusedTopicTreeRef = useRef<any>(null);
  const handleTreeDataAddId = (treeData: UnsTreeNode[]) => {
    treeData.forEach((item) => {
      item.id = item.path || '';
      if (item.children && item.children.length) {
        handleTreeDataAddId(item.children);
      }
    });
  };

  useEffect(() => {
    loadData({ reset: true });
  }, []);

  useEffect(() => {
    if (selectedNode?.id) {
      // 如果有选中node，要清空 其他topic
      setCurrentUnusedTopicNode(initNode);
    }
  }, [selectedNode?.id]);

  const { connect, disconnect } = useWebSocket(
    `${location.protocol.includes('https') ? 'wss' : 'ws'}://${location.host}/inter-api/supos/uns/external/topic`,
    {
      reconnectLimit: 0,
      onMessage: (event) => {
        if (event.data === 'refresh_notify') {
          initUnusedTopicTreeData({ query: unusedTopicTreeRef?.current?.searchQuery });
        }
      },
      manual: true,
      onError: (error) => console.error('WebSocket error:', error),
    }
  );

  //滚动到目标树节点
  const scrollUnusedTopicTreeNode = (id: string) => {
    setTimeout(() => {
      if (unusedTopicTreeRef.current) unusedTopicTreeRef.current.scrollTo(id);
    }, 500);
  };
  useEffect(() => {
    const updatePanelSize = () => {
      const isShow = !!splitterWrapRef?.current?.offsetHeight;
      if (splitterWrapRef.current && isShow && !unusedTopicPanelSize[0]) {
        setUnusedTopicPanelSize([splitterWrapRef.current.offsetHeight - panelCloseSize, panelCloseSize]);
      }
    };
    updatePanelSize();
    const resizeObserver = new ResizeObserver(updatePanelSize);
    if (splitterWrapRef.current) {
      resizeObserver.observe(splitterWrapRef.current);
    }
    return () => {
      resizeObserver.disconnect();
    };
  }, [unusedTopicPanelSize[0]]);

  useEffect(() => {
    //监听选中节点获取面包屑数据
    if (!!currentUnusedTopicNode.id && treeType === 'uns') {
      const nodeParents = getParentNodes(unusedTopicTreeData as any, currentUnusedTopicNode.id);
      setUnusedTopicBreadcrumbList(nodeParents as any);
    }
  }, [currentUnusedTopicNode.id, unusedTopicTreeData, treeType]);
  // const abortControllerRef = useRef<AbortController>();

  //初始化UnusedTopic树数据
  const initUnusedTopicTreeData: InitTreeDataFnType = ({ reset, query }, cb) => {
    if (treeType === 'uns') {
      // 取消之前的请求
      // const prevController = abortControllerRef.current;
      // abortControllerRef.current = new AbortController();
      // prevController?.abort();
      setUnusedTopicLoading(true);
      // getExternalTreeData({ key: query }, { signal: abortControllerRef.current.signal })
      getExternalTreeData({ key: query })
        .then((res: any) => {
          // 组树
          if (res?.length) {
            handleTreeDataAddId(res);
            setUnusedTopicTreeData(res);
            if (reset) {
              changeCurrentUnusedTopicPath();
              unusedTopicTreeRef.current?.setExpandedArr([]);
              scrollUnusedTopicTreeNode(res[0]?.path);
            }
          } else {
            setUnusedTopicTreeData([]);
            setCurrentUnusedTopicNode(initNode);
            unusedTopicTreeRef.current?.setExpandedArr([]);
          }
          cb?.();
        })
        .catch((err) => {
          setUnusedTopicTreeData([]);
          setCurrentUnusedTopicNode(initNode);
          console.log(err);
        })
        .finally(() => {
          setUnusedTopicLoading(false);
        });
    }
  };

  useEffect(() => {
    if (enableAutoCategorization) return;
    if (treeType === 'uns') {
      connect?.();
    }
    // 初始化 其他topic
    initUnusedTopicTreeData({ reset: true });
    setCurrentTreeMapType('all');
    return () => {
      // abortControllerRef.current?.abort();
      disconnect?.();
    };
  }, [treeType, enableAutoCategorization]);

  const { isTreeMapVisible, setTreeMapVisible } = useUnsTreeMapContext();
  const { isH5 } = useMediaSize();
  const treeMapHtml = (
    <div ref={splitterWrapRef} style={{ height: 'calc(100% - 48px)' }}>
      <Splitter layout="vertical" onResize={setUnusedTopicPanelSize} className="unusedTopicTree-Splitter">
        <Splitter.Panel min={120} size={unusedTopicPanelSize[0]}>
          <Tree
            changeCurrentPath={changeCurrentPath}
            treeNodeExtra={(dataNode) => {
              return (
                <TreeNodeExtra
                  type={dataNode.pathType}
                  handleDelete={() => handleDelete(dataNode)}
                  handleCopy={() => handleCopy(dataNode)}
                />
              );
            }}
          />
        </Splitter.Panel>
        {/*暂时隐藏*/}
        {treeType === 'uns' && !enableAutoCategorization && (
          <Splitter.Panel size={unusedTopicPanelSize[1]} min={panelCloseSize} style={{ overflow: 'hidden' }}>
            <div
              className="unusedTopicTree-collapsible"
              onClick={() => {
                if (unusedTopicPanelSize[1] === panelCloseSize) {
                  setUnusedTopicPanelSize([splitterWrapRef.current?.offsetHeight - panelOpenSize, panelOpenSize]);
                } else {
                  setUnusedTopicPanelSize([splitterWrapRef.current?.offsetHeight - panelCloseSize, panelCloseSize]);
                }
              }}
            >
              {formatMessage('uns.otherTopic', 'Raw Topics')}
              {unusedTopicPanelSize[1] === panelCloseSize ? <ChevronUp /> : <ChevronDown />}
            </div>
            <div style={{ height: 'calc(100% - 48px)', padding: '8px 14px 0' }}>
              <Loading spinning={unusedTopicLoading}>
                <UnusedTopicTree
                  unsTreeRef={unusedTopicTreeRef}
                  treeData={unusedTopicTreeData}
                  currentNode={currentUnusedTopicNode}
                  initTreeData={initUnusedTopicTreeData}
                  setTreeMap={setTreeMap}
                  changeCurrentPath={changeCurrentUnusedTopicPath}
                  treeType={treeType}
                  unusedTopicBreadcrumbList={unusedTopicBreadcrumbList}
                />
              </Loading>
            </div>
          </Splitter.Panel>
        )}
      </Splitter>
    </div>
  );

  const handleCopy = (item: UnsTreeNode) => {
    const { id } = item;
    switch (treeType) {
      case 'uns':
        setPasteNode(item);
        message.success(formatMessage('common.copySuccess'));
        break;
      case 'template':
        operationFns?.openTemplateModal('copyTemplate', id as string);
        break;
      default:
        break;
    }
  };

  return !isH5 ? (
    <ComLeft style={{ overflow: 'hidden' }} resize defaultWidth={360}>
      <div
        className="treemapTitle"
        onClick={() => {
          setTreeMap(true);
          changeCurrentPath();
        }}
      >
        <TreeViewIcon />
        <span>{formatMessage('uns.treeList')}</span>
      </div>
      {treeMapHtml}
    </ComLeft>
  ) : (
    isTreeMapVisible && (
      <Drawer
        className="treemap-drawer"
        rootClassName="treemap-drawer-root"
        placement="right"
        onClose={() => setTreeMapVisible(false)}
        open={isTreeMapVisible}
        getContainer={false}
      >
        {treeMapHtml}
      </Drawer>
    )
  );
};

export default LeftDom;
