import { type FC, useEffect, useRef } from 'react';
import type { PageProps } from '@/common-types';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { App } from 'antd';
import { Add, Document, Folder, Renew } from '@carbon/icons-react';
import useTranslate from '@/hooks/useTranslate';
import ComLeft from '@/components/com-layout/ComLeft.tsx';
import { SortableTree } from './components/menu-tree';
import { MenuStoreProvider, useMenuStore } from './store/menuStore.tsx';
import MenuContent from './components/menu-content/MenuContent.tsx';
import EmptyDetail from './components/empty-detail';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { useTabsContext } from '@/contexts/tabs-context.ts';
import { AuthButton } from '@/components/auth';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import styles from './index.module.scss';

const Module: FC<PageProps> = ({ title }) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const isFirstRender = useRef(true);
  const { requestMenu, menuList, menuTree, setContentType, contentType, setSelectNode, selectNode, loading } =
    useMenuStore((state) => ({
      requestMenu: state.requestMenu,
      menuList: state.menuList,
      menuTree: state.menuTree,
      setContentType: state.setContentType,
      contentType: state.contentType,
      setSelectNode: state.setSelectNode,
      selectNode: state.selectNode,
      loading: state.loading,
    }));
  const { TabsContext } = useTabsContext();

  const lang = useI18nStore((state) => state.lang);
  const menuLabel = (node: any) => formatMessage(node?.showName, undefined, node?.showName || node?.label || '');
  const isSystemMenu = (node?: any) => Boolean(node?.coreResourceId && node?.editEnable === false);

  const onAddExternalMenu = () => {
    const parentId =
      selectNode && !isSystemMenu(selectNode)
        ? selectNode.type === 1
          ? selectNode.id
          : selectNode.parentId
        : undefined;
    setSelectNode({
      id: '',
      parentId,
      type: 2,
      code: '',
      showName: '',
      sort: (menuList?.length || 0) * 10 + 100,
      url: '',
      urlType: 2,
      openType: 0,
      enable: true,
      children: [],
    } as any);
    setContentType('addMenu');
  };

  const onRefreshMenu = () => {
    setContentType(null);
    setSelectNode(null);
    requestMenu().then(() => {
      message.success(formatMessage('common.refreshSuccessful'));
    });
  };

  useEffect(() => {
    requestMenu();
  }, [requestMenu]);

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
    } else {
      TabsContext?.current?.onRefreshTab?.('/MenuConfiguration');
    }
  }, [TabsContext, lang]);

  return (
    <ComLayout className={styles.splitLayout}>
      <ComContent className={styles.menuConfigurationContent} hasBack={false} title={title}>
        <ComLayout className={styles.splitLayout}>
          <ComLeft
            resize
            defaultWidth={256}
            className={styles.sidebar}
            style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden', minHeight: 0 }}
          >
            <div className={styles.sidebarHeader}>
              <span className={styles.sidebarTitle}>{formatMessage('MenuConfiguration.menuList')}</span>
              <div className={styles.sidebarActions}>
                <AuthButton
                  auth={ButtonPermission['MenuConfiguration.addMenu']}
                  type="text"
                  className={styles.iconAction}
                  icon={<Add size={16} />}
                  title={formatMessage('MenuConfiguration.addMenu')}
                  onClick={onAddExternalMenu}
                />
                <button
                  type="button"
                  className={styles.iconAction}
                  title={formatMessage('common.refresh')}
                  onClick={onRefreshMenu}
                >
                  <Renew size={16} />
                </button>
              </div>
            </div>
            <div className={styles.tree}>
              <SortableTree
                loading={loading}
                treeData={menuTree as any}
                renderLabel={menuLabel}
                style={{ height: '100%' }}
                indicator
                indentationWidth={24}
                selectedKey={selectNode ? selectNode.id : null}
                onSelect={(key: any, node: any) => {
                  setSelectNode(node);
                  if (!key) {
                    setContentType(null);
                  } else {
                    setContentType(node.type === 1 ? 'editGroup' : 'editMenu');
                  }
                }}
                leftExtra={(node: any) => {
                  if (node.type === 1) {
                    return <Folder size={16} style={{ flexShrink: 0 }} />;
                  }
                  if (node.type === 2) {
                    return <Document size={16} style={{ flexShrink: 0 }} />;
                  }
                  return null;
                }}
                disabledSelected={(node: any) => node.type == 4 && !node.url}
                disabledDraggable={isSystemMenu}
                allowDrop={() => false}
              />
            </div>
          </ComLeft>
          <ComContent className={styles.detailPanel} hasBack={false} mustShowTitle={false} border={false}>
            {contentType ? <MenuContent /> : <EmptyDetail />}
          </ComContent>
        </ComLayout>
      </ComContent>
    </ComLayout>
  );
};

const WrapperModule: FC<PageProps> = (props) => {
  return (
    <MenuStoreProvider>
      <Module {...props} />
    </MenuStoreProvider>
  );
};

export default WrapperModule;
