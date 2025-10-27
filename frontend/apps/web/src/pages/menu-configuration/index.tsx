import { type FC, useEffect, useRef } from 'react';
import type { PageProps } from '@/common-types';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { App, Flex } from 'antd';
import {
  Close,
  Document,
  DocumentAdd,
  Folder,
  FolderAdd,
  ListDropdown,
  Renew,
  View,
  ViewOff,
} from '@carbon/icons-react';
import useTranslate from '@/hooks/useTranslate';
import ComLeft from '@/components/com-layout/ComLeft.tsx';
import { SortableTree } from './components/menu-tree';
import useMenuSetting from './components/menu-setting/useMenuSetting.tsx';
import { AuthButton, AuthWrapper } from '@/components/auth';
import { MenuStoreProvider, useMenuStore } from './store/menuStore.tsx';
import MenuContent from './components/menu-content/MenuContent.tsx';
import EmptyDetail from './components/empty-detail';
import { batchEditResourceApi, deleteResourceApi } from '@/apis/inter-api/resource.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { useTabsContext } from '@/contexts/tabs-context.ts';

const Module: FC<PageProps> = ({ title }) => {
  const formatMessage = useTranslate();
  const { onMenuModalOpen, MenuModal } = useMenuSetting();
  const { modal, message } = App.useApp();
  const isFirstRender = useRef(true);
  const { requestMenu, menuTree, setContentType, contentType, setSelectNode, selectNode, setMenuInfo, loading } =
    useMenuStore((state) => ({
      requestMenu: state.requestMenu,
      menuTree: state.menuTree,
      setContentType: state.setContentType,
      contentType: state.contentType,
      setSelectNode: state.setSelectNode,
      selectNode: state.selectNode,
      setMenuInfo: state.setMenuInfo,
      loading: state.loading,
    }));
  const { TabsContext } = useTabsContext();

  const lang = useI18nStore((state) => state.lang);

  useEffect(() => {
    requestMenu();
  }, []);

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
    } else {
      TabsContext?.current?.onRefreshTab?.('/MenuConfiguration');
    }
  }, [lang]);

  const onEnabledHandle = (node: any, e: any) => {
    e.stopPropagation();
    batchEditResourceApi([
      {
        id: node.id,
        enable: !node?.enable,
      },
    ]).then(() => {
      message.success(formatMessage('common.optsuccess'));
      requestMenu();
    });
  };

  return (
    <ComLayout>
      {MenuModal}
      <ComContent
        hasBack={false}
        title={
          <Flex align="center" gap={8} style={{ lineHeight: 1 }}>
            <ListDropdown
              size={20}
              style={{ justifyContent: 'center', verticalAlign: 'middle' }}
              onClick={() => onMenuModalOpen()}
            />
            <span>{title}</span>
          </Flex>
        }
      >
        <ComLayout>
          <ComLeft
            resize
            defaultWidth={360}
            style={{ display: 'flex', flexDirection: 'column', padding: '16px 0 16px 16px' }}
          >
            <Flex style={{ marginBottom: 16, marginRight: 16 }} justify="space-between" align="center">
              <span style={{ fontSize: 20, fontWeight: 500 }}>{formatMessage('MenuConfiguration.menuList')}</span>
              <Flex align="center">
                <AuthButton
                  disabled={selectNode?.editEnable === false || [2, 3].includes(selectNode?.type || 0)}
                  size="small"
                  type="text"
                  onClick={() => setContentType('addMenu')}
                  auth={ButtonPermission['MenuConfiguration.addMenu']}
                  title={formatMessage('MenuConfiguration.addMenu')}
                >
                  <DocumentAdd />
                </AuthButton>
                <AuthButton
                  disabled={selectNode?.editEnable === false || [1, 2, 3].includes(selectNode?.type || 0)}
                  size="small"
                  type="text"
                  onClick={() => setContentType('addGroup')}
                  auth={ButtonPermission['MenuConfiguration.addMenu']}
                  title={formatMessage('MenuConfiguration.addGroup')}
                >
                  <FolderAdd />
                </AuthButton>
                <AuthButton
                  size="small"
                  type="text"
                  title={formatMessage('common.refresh')}
                  onClick={() => {
                    setContentType(null);
                    setSelectNode(null);
                    requestMenu().then(() => {
                      message.success(formatMessage('common.refreshSuccessful'));
                    });
                  }}
                >
                  <Renew />
                </AuthButton>
              </Flex>
            </Flex>
            <SortableTree
              loading={loading}
              onHandleDragEnd={(newData: any, tree: any) => {
                // 先设置，后请求
                setMenuInfo(tree, newData);
                batchEditResourceApi(
                  newData
                    ?.map((item: any, index: number) => ({
                      id: item.id,
                      sort: (item?.index ?? index) + 1,
                      parentId: item.parentId,
                    }))
                    ?.filter((f: any) => !f.id?.includes('tab_container'))
                ).then(() => {
                  message.success(formatMessage('common.optsuccess'));
                  requestMenu();
                });
              }}
              treeData={menuTree as any}
              style={{ flex: 1, overflow: 'auto', scrollbarGutter: 'stable', paddingRight: 8 }}
              indicator
              indentationWidth={32}
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
                  return <Folder style={{ flexShrink: 0 }} />;
                } else if (node.type === 2) {
                  return <Document style={{ flexShrink: 0 }} />;
                }
                return null;
              }}
              rightExtra={(node: any) => {
                if (!node?.editEnable) return null;
                return (
                  <Flex gap={8}>
                    {node?.enable ? (
                      <AuthWrapper auth={ButtonPermission['MenuConfiguration.enabledMenu']}>
                        <Flex
                          title={formatMessage('MenuConfiguration.disabled')}
                          onClick={(e) => onEnabledHandle(node, e)}
                        >
                          <View style={{ cursor: 'pointer' }} />
                        </Flex>
                      </AuthWrapper>
                    ) : (
                      <AuthWrapper auth={ButtonPermission['MenuConfiguration.enabledMenu']}>
                        <Flex
                          title={formatMessage('MenuConfiguration.enable')}
                          onClick={(e) => onEnabledHandle(node, e)}
                        >
                          <ViewOff style={{ cursor: 'pointer' }} />
                        </Flex>
                      </AuthWrapper>
                    )}
                    {(![1, 2].includes(node?.type) ||
                      (node?.type === 1 && !['80', '50'].includes(node.id)) ||
                      (node?.type === 2 && node?.urlType === 2)) && (
                      <AuthWrapper auth={ButtonPermission['MenuConfiguration.deleteMenu']}>
                        <Flex
                          title={formatMessage('common.delete')}
                          onClick={(e) => {
                            e.stopPropagation();
                            modal.confirm({
                              title: formatMessage('common.deleteConfirm'),
                              onOk: async () => {
                                return deleteResourceApi(node.id).then(() => {
                                  message.success(formatMessage('common.deleteSuccessfully'));
                                  requestMenu();
                                });
                              },
                              okButtonProps: {
                                title: formatMessage('common.confirm'),
                              },
                              cancelButtonProps: {
                                title: formatMessage('common.cancel'),
                              },
                            });
                          }}
                        >
                          <Close style={{ cursor: 'pointer' }} />
                        </Flex>
                      </AuthWrapper>
                    )}
                  </Flex>
                );
              }}
              disabledSelected={(node: any) => {
                return node.type == 4 && !node.url;
              }}
              allowDrop={({ drop, drag }: any) => {
                // 跟目录都可以放置
                if (!drop) return true;
                // 固定不可放置
                if (drag?.fixed || drop.fixed) return false;
                // 不可放置 文件
                if (drop.type === 2) return false;
                // 组 不可放置 组
                if (drag?.type === 1 && drop?.type === 1) {
                  return false;
                }
                return true;
              }}
            />
          </ComLeft>
          <ComContent>{contentType ? <MenuContent /> : <EmptyDetail />}</ComContent>
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
