import { Copy, Download, Upload } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useTreeStore } from './store/treeStore';
import { useClipboard, useTranslate } from '@/hooks';
import { type FC, type ReactNode, useCallback, useRef } from 'react';
import ComBreadcrumb from '@/components/com-breadcrumb';
import ComText from '@/components/com-text';
import ImportModal from './components/import-modal';
import ExportModal from './components/export-modal';
import { AuthButton } from '@/components/auth';
import { ButtonPermission } from '@/common-types/button-permission';

interface TopDomProps {
  setCurrentUnusedTopicNode: any;
  unusedTopicBreadcrumbList: any;
  currentUnusedTopicNode: any;
  changeCurrentPath: any;
}
const TopDom: FC<TopDomProps> = ({ setCurrentUnusedTopicNode, unusedTopicBreadcrumbList, currentUnusedTopicNode }) => {
  const formatMessage = useTranslate();
  const copyPathRef = useRef(null);
  const importRef = useRef<any>(null);
  const exportRef = useRef<any>(null);
  const { treeType, currentTreeMapType, breadcrumbList, selectedNode, setSelectedNode, loadData } = useTreeStore(
    (state) => ({
      treeType: state.treeType,
      currentTreeMapType: state.currentTreeMapType,
      breadcrumbList: state.breadcrumbList,
      selectedNode: state.selectedNode,
      setSelectedNode: state.setSelectedNode,
      loadData: state.loadData,
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
              title: (
                <span className="com-breadcrumb-current" title={name}>
                  {name}
                </span>
              ),
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
              <Copy {...toolbarIconProps} />
            </div>
          )
        }
      />
    ),
    [setCurrentUnusedTopicNode, setSelectedNode, currentTreeMapType, formatMessage]
  );

  return (
    <>
      <div className="chartTop">
        {treeType === 'uns' ? (
          <div className="chartTopL">
            {currentTreeMapType === 'all' && selectedNode?.id
              ? getTopicBreadcrumb(breadcrumbList, selectedNode.pathType === 0 ? false : null)
              : null}
            {currentTreeMapType === 'unusedTopic' && currentUnusedTopicNode.id
              ? getTopicBreadcrumb(unusedTopicBreadcrumbList)
              : null}
          </div>
        ) : (
          <span />
        )}
        <div className="chartTopR">
          {treeType === 'uns' && (
            <>
              <AuthButton
                auth={ButtonPermission['uns.import']}
                type="primary"
                icon={<Upload size={16} />}
                onClick={() => importRef.current?.setOpen(true)}
              >
                {formatMessage('common.import')}
              </AuthButton>
              <AuthButton
                auth={ButtonPermission['uns.export']}
                type="primary"
                icon={<Download size={16} />}
                onClick={() => exportRef.current?.setOpen(true)}
              >
                {formatMessage('common.export')}
              </AuthButton>
            </>
          )}
        </div>
      </div>
      <ImportModal importRef={importRef} initTreeData={loadData} />
      <ExportModal exportRef={exportRef} />
    </>
  );
};

export default TopDom;
