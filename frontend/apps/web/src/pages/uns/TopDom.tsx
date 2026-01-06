import { Button, Flex } from 'antd';
import { Copy, Rss, Workspace } from '@carbon/icons-react';
import { ButtonPermission } from '@/common-types/button-permission';
import { getTreeStoreSnapshot, useTreeStore, useTreeStoreRef } from './store/treeStore';
import { useClipboard, useTranslate } from '@/hooks';
import { type FC, type ReactNode, useCallback, useRef } from 'react';
import { ExportModal, ImportModal } from '@/pages/uns/components';
import { AuthButton, AuthWrapper } from '@/components/auth';
import ComBreadcrumb from '@/components/com-breadcrumb';
import ComText from '@/components/com-text';
import { useBaseStore } from '@/stores/base';

interface TopDomProps {
  setCurrentUnusedTopicNode: any;
  unusedTopicBreadcrumbList: any;
  currentUnusedTopicNode: any;
  changeCurrentPath: any;
}
const TopDom: FC<TopDomProps> = ({
  setCurrentUnusedTopicNode,
  unusedTopicBreadcrumbList,
  currentUnusedTopicNode,
  changeCurrentPath,
}) => {
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const formatMessage = useTranslate();
  const exportRef = useRef<any>(null);
  const importRef = useRef<any>(null);
  const copyPathRef = useRef(null);
  const { treeType, currentTreeMapType, breadcrumbList, selectedNode, setSelectedNode, treeMap } = useTreeStore(
    (state) => ({
      treeType: state.treeType,
      currentTreeMapType: state.currentTreeMapType,
      breadcrumbList: state.breadcrumbList,
      selectedNode: state.selectedNode,
      setSelectedNode: state.setSelectedNode,
      treeMap: state.treeMap,
    })
  );

  useClipboard(
    copyPathRef as any,
    currentTreeMapType === 'all' ? breadcrumbList.slice(-1)?.[0]?.path : currentUnusedTopicNode.path
  );

  const getTopicBreadcrumb = useCallback(
    (pArr: any[], addonAfter?: ReactNode | false) => (
      <ComBreadcrumb
        style={{ fontWeight: 700 }}
        items={pArr?.map((e: any, idx: number) => {
          const name = currentTreeMapType === 'all' ? e.name : e.pathName || e.name;
          if (idx + 1 === pArr?.length) {
            return {
              title: name,
            };
          }
          return {
            title: <ComText>{name}</ComText>,
            onClick: () => {
              if (currentTreeMapType === 'all') {
                setSelectedNode(e);
              } else {
                setCurrentUnusedTopicNode(e);
              }
            },
          };
        })}
        addonAfter={
          addonAfter ? (
            addonAfter
          ) : addonAfter === false ? null : (
            <div className="copyBox" ref={copyPathRef} title={formatMessage('common.copy')}>
              <Copy />
            </div>
          )
        }
      />
    ),
    [setCurrentUnusedTopicNode, setSelectedNode, currentTreeMapType]
  );

  const stateRef = useTreeStoreRef();
  const { loadData, setTreeMap } = getTreeStoreSnapshot(stateRef, (state) => ({
    loadData: state.loadData,
    setTreeMap: state.setTreeMap,
  }));

  return (
    <div className="chartTop">
      {treeMap ? (
        <div className="treemapTitle" style={{ padding: 0 }}></div>
      ) : treeType === 'uns' ? (
        <div className="chartTopL">
          {currentTreeMapType === 'all' && selectedNode?.id
            ? getTopicBreadcrumb(
                breadcrumbList,
                selectedNode.pathType === 0 ? (
                  false
                ) : selectedNode.pathType === 2 && systemInfo.useAliasPathAsTopic ? (
                  <Flex
                    align="center"
                    style={{ cursor: 'pointer' }}
                    onClick={() => {
                      const scrollWrap = document.querySelector('.topicDetailContent');
                      const targetNode = document.getElementById('sqlQuery');
                      if (scrollWrap && targetNode) {
                        const diffY =
                          scrollWrap.scrollTop +
                          targetNode.getBoundingClientRect().top -
                          scrollWrap.getBoundingClientRect().top;
                        scrollWrap.scrollTo({
                          top: diffY,
                          behavior: 'smooth',
                        });
                      }
                    }}
                    title={formatMessage('common.subscribe')}
                  >
                    <Rss />
                  </Flex>
                ) : null
              )
            : null}
          {currentTreeMapType === 'unusedTopic' && currentUnusedTopicNode.id
            ? getTopicBreadcrumb(unusedTopicBreadcrumbList)
            : null}
        </div>
      ) : (
        <span />
      )}
      <div className="chartTopR">
        <AuthButton
          auth={ButtonPermission['uns.unsImport']}
          type="primary"
          onClick={() => importRef?.current?.setOpen(true)}
        >
          {formatMessage('common.import')}
        </AuthButton>
        <AuthWrapper auth={ButtonPermission['uns.unsExport']}>
          <Button
            onClick={() => {
              exportRef.current?.setOpen(true);
            }}
          >
            {formatMessage('uns.export')}
          </Button>
        </AuthWrapper>
        <Button
          title={formatMessage('uns.backOverview')}
          style={{ padding: 8 }}
          onClick={() => {
            setTreeMap(true);
            changeCurrentPath();
          }}
        >
          <Workspace size={16} />
        </Button>
      </div>
      <ImportModal importRef={importRef} initTreeData={loadData} />
      <ExportModal exportRef={exportRef} />
    </div>
  );
};

export default TopDom;
