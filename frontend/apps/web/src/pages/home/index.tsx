import { type ReactElement, useCallback, useEffect, useRef, useState } from 'react';
import { Button, message, Modal, Space, Tabs, type TabsProps, Typography } from 'antd';
import OverviewList from '@/pages/home/components/OverviewList.tsx';
import type { ResourceProps } from '@/stores/types';
import styles from './index.module.scss';
import StickyBox from 'react-sticky-box';
import { useGuideSteps, useTranslate } from '@/hooks';
import { useNavigate } from 'react-router';
import { guideSteps } from './guide-steps';
import { queryExamples, installExample, unInstallExample } from '@/apis/inter-api/example';
import { Code, Pin, TemperatureWater } from '@carbon/icons-react';
import { useActivate } from '@/contexts/tabs-lifecycle-context.ts';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { fetchBaseStore, useBaseStore } from '@/stores/base';
import { useThemeStore } from '@/stores/theme-store.ts';
// import { ButtonPermission } from '@/common-types/button-permission';
// import { AuthButton, AuthWrapper } from '@/components';
// import ImportModal from './components/import-modal';
// import ExportModal from './components/export-modal';
// import { getGlobalExportRecords } from '@/apis/inter-api/global.ts';
import ComClickTrigger from '@/components/com-click-trigger';
import IframeWrapper from '../../components/iframe-wrapper';
import { FullscreenOutlined } from '@ant-design/icons';
import screenfull from 'screenfull';
import { useI18nStore } from '@/stores/i18n-store.ts';

const { Title, Paragraph } = Typography;

// example 返回的数据结构
interface exampleItemTypes {
  id: number | string;
  name: string;
  description: string;
  status: number;
  type: number;
  dashboardType?: number;
  dashboardId?: string;
  dashboardName?: string;
}

const exampleTypes: { [x: string | number]: string } = {
  1: 'OTDataConnections',
  2: 'ITDataConnections',
};

export interface ExampleProps extends ResourceProps {
  iconComp?: ReactElement;
  status?: number | string;
}

const Index = () => {
  const selectedIdRef = useRef<string | number | null>(null);
  const { systemInfo, homeTree, homeTabGroup } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
    homeTree: state.homeTree,
    homeTabGroup: state.homeTabGroup,
  }));
  const lang = useI18nStore((state) => state.lang);
  const [tabKey, setTabKey] = useState('common.overview');
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const [pathname, setPathname] = useState('');
  const [exampleDataSource, setExampleDataSource] = useState<ExampleProps[]>([]);
  const [loadingViews, setLoadingViews] = useState<string[]>([]);
  // const [importModal, setImportModal] = useState(false);
  // const exportRef = useRef<any>(null);
  // const [exportRecords, setExportRecords] = useState([]);

  useEffect(() => {
    fetchBaseStore?.();
    // getRecords?.();
  }, []);
  useEffect(() => {
    getExamples();
  }, [lang]);
  const primaryColor = useThemeStore((state) => state.primaryColor);

  // 解决routesStore?.fetchRoutes导致跳转路由方法失效的问题：通过state改变再触发navigate跳转
  const handleNavigate = useCallback((path: string) => {
    setPathname(path);
  }, []);
  useEffect(() => {
    if (pathname) {
      navigate(pathname);
    }
  }, [pathname]);

  // 注意：对某个页面添加steps时，请务必在 stores -> base -> index.tsx -> GuidePagePaths 中添加该页面路由
  useGuideSteps(guideSteps(handleNavigate, { appTitle: systemInfo.appTitle }, primaryColor));

  useActivate(() => {
    // 每次进home页刷新下，保持页面完整
    fetchBaseStore?.();
    // getRecords?.();
  });

  // 获取example列表
  const getExamples = (open = true) => {
    queryExamples().then((data: any) => {
      const newExamplesMap = new Map();
      const newExampleDataSource: ExampleProps[] = [];

      if (!data?.length) return;

      // 安装后打开页面
      if (selectedIdRef.current && open) {
        const { dashboardId, dashboardName, dashboardType } = data.find(
          (item: exampleItemTypes) => item.id === selectedIdRef.current
        );
        selectedIdRef.current = null;
        if (dashboardId) {
          window.open(
            `/dashboards/preview?id=${dashboardId}&type=${dashboardType}&status=preview&name=${dashboardName}`
          );
        }
      }

      data.forEach((item: exampleItemTypes) => {
        const type = exampleTypes[item.type];

        if (!newExamplesMap.has(type)) {
          newExamplesMap.set(type, {
            showName: item.name,
            id: type,
            children: [],
          });
        }

        newExamplesMap.set(type, {
          ...newExamplesMap.get(type),
          iconComp: <Pin />,
          children: [
            ...newExamplesMap.get(type).children,
            {
              showName: item.name,
              parentId: type,
              id: item.id,
              status: item.status,
              showDescription: item.description,
              iconComp: type === 'OTDataConnections' ? <TemperatureWater color="#1D77FE" /> : <Code color="#1D77FE" />,
            },
          ],
        });
      });

      newExamplesMap.forEach((value) => {
        newExampleDataSource.push(value);
      });

      setExampleDataSource(newExampleDataSource);
    });
  };

  // 安装example
  const handleInstall = (params: ExampleProps) => {
    Modal.confirm({
      title: formatMessage('common.confirmInstall'),
      onOk: () => {
        if (!params.id) return;
        setLoadingViews([...loadingViews, params.id]);
        selectedIdRef.current = params.id;
        installExample(params.id)
          .then(() => {
            message.success(formatMessage('common.installedSuccess'));
            getExamples();
          })
          .finally(() => {
            loadingViews.splice(loadingViews.indexOf(params.id as string), 1);
            setLoadingViews(loadingViews);
          });
      },
      okButtonProps: {
        title: formatMessage('common.confirm'),
      },
      cancelButtonProps: {
        title: formatMessage('common.cancel'),
      },
    });
  };

  // 卸载example
  const handleUnInstall = (params: ExampleProps) => {
    Modal.confirm({
      title: formatMessage('common.confirmUnInstall'),
      onOk: () => {
        if (!params.id) return;
        setLoadingViews([...loadingViews, params.id]);
        selectedIdRef.current = params.id;
        unInstallExample(params.id)
          .then(() => {
            message.success(formatMessage('common.unInstalledSuccess'));
            getExamples(false);
          })
          .finally(() => {
            loadingViews.splice(loadingViews.indexOf(params.id as string), 1);
            setLoadingViews(loadingViews);
          });
      },
      okButtonProps: {
        title: formatMessage('common.confirm'),
      },
      cancelButtonProps: {
        title: formatMessage('common.cancel'),
      },
    });
  };

  // 切换tab
  const handleChangeTab = (key: string) => {
    setTabKey(key);
    if (key !== 'example') return;
    getExamples();
  };

  const renderTabBar: TabsProps['renderTabBar'] = (props, DefaultTabBar) => (
    <StickyBox offsetTop={0} offsetBottom={20} style={{ zIndex: 1 }}>
      <DefaultTabBar {...props} />
    </StickyBox>
  );

  const renderExampleOpt = (params: ExampleProps) => {
    return params.status === 1
      ? {
          type: 'button',
          key: 'install',
          button: { type: 'primary' },
          label: formatMessage('common.install'),
          onClick: () => {
            handleInstall(params);
          },
        }
      : {
          type: 'button',
          key: 'unInstall',
          label: formatMessage('common.unInstall'),
          onClick: () => {
            handleUnInstall(params);
          },
        };
  };

  // const getRecords = () => {
  //   return getGlobalExportRecords().then((data) => {
  //     setExportRecords(data);
  //   });
  // };
  const isHidden = !['common.overview', 'common.example'].includes(tabKey);
  return (
    <ComLayout>
      <ComContent title={<div></div>} hasBack={false} mustShowTitle={false}>
        <div className={styles['home-title']} style={{ display: isHidden ? 'none' : 'block' }}>
          <Title style={{ fontWeight: 400, marginBottom: 5 }} type="secondary" level={2}>
            {formatMessage('common.welcome', { appTitle: systemInfo?.appTitle })}
          </Title>

          <Paragraph style={{ marginBottom: 0 }}>{formatMessage('common.excellence')}</Paragraph>
        </div>
        <div className={styles['home-tabs']} style={{ height: isHidden ? '100%' : undefined }}>
          <Tabs
            renderTabBar={renderTabBar}
            defaultActiveKey="common.overview"
            activeKey={tabKey}
            onChange={handleChangeTab}
            tabBarExtraContent={
              <Space
                style={{
                  marginRight: 36,
                  background: 'var(--supos-bg-color)',
                }}
              >
                <ComClickTrigger
                  triggerCount={2}
                  style={{ width: 50, height: 40 }}
                  onTrigger={() => {
                    console.warn(useBaseStore.getState());
                  }}
                />
                {!isHidden ? (
                  <>
                    {/*<AuthButton*/}
                    {/*  auth={ButtonPermission['Home.import']}*/}
                    {/*  type="primary"*/}
                    {/*  onClick={() => setImportModal(true)}*/}
                    {/*>*/}
                    {/*  <Flex gap={8}>*/}
                    {/*    <Download />*/}
                    {/*    {formatMessage('common.import')}*/}
                    {/*  </Flex>*/}
                    {/*</AuthButton>*/}
                    {/*<AuthWrapper auth={ButtonPermission['Home.export']}>*/}
                    {/*  <Badge dot={exportRecords?.some((s: any) => !s.confirm)}>*/}
                    {/*    <Button*/}
                    {/*      color="default"*/}
                    {/*      variant="filled"*/}
                    {/*      style={{ background: '#c6c6c6', color: '#161616' }}*/}
                    {/*      onClick={() => {*/}
                    {/*        exportRef.current?.setOpen(true);*/}
                    {/*      }}*/}
                    {/*    >*/}
                    {/*      <Flex gap={8}>*/}
                    {/*        <Export />*/}
                    {/*        {formatMessage('common.export')}*/}
                    {/*      </Flex>*/}
                    {/*    </Button>*/}
                    {/*  </Badge>*/}
                    {/*</AuthWrapper>*/}
                  </>
                ) : (
                  <Button
                    title={formatMessage('common.fullScreen')}
                    icon={<FullscreenOutlined />}
                    onClick={() => {
                      if (screenfull.isEnabled) {
                        const el = document.getElementById(tabKey);
                        if (el) {
                          screenfull.request(el);
                        } else {
                          message.error('未找到全屏元素');
                        }
                      } else {
                        message.error('该浏览器,不支持全屏功能');
                      }
                    }}
                  />
                )}
              </Space>
            }
            items={[
              ...(homeTabGroup?.map((item) => {
                if (item.code === 'common.overview') {
                  return {
                    label: formatMessage('common.overview'),
                    key: 'common.overview',
                    children: (
                      <OverviewList
                        list={homeTree}
                        style={{
                          '--supos-line-height': 2,
                          '--supos-card-height': '125px',
                        }}
                      />
                    ),
                  };
                } else if (item.code === 'common.example') {
                  return {
                    label: formatMessage('common.example'),
                    key: 'common.example',
                    children: (
                      <OverviewList
                        list={exampleDataSource}
                        loadingViews={loadingViews}
                        type="example"
                        customOptRender={renderExampleOpt}
                      />
                    ),
                  };
                }
                return {
                  label: item.showName || item.code,
                  key: item.code + '_tab',
                  children: (
                    <IframeWrapper id={item.code + '_tab'} title={item.showName || item.code} src={item?.url || ''} />
                  ),
                };
              }) || []),
            ]?.filter((f) => f.key !== 'common.example')}
          />
        </div>
        {/*<ImportModal importModal={importModal} setImportModal={setImportModal} />*/}
        {/*<ExportModal setButtonExportRecords={setExportRecords} exportRef={exportRef} />*/}
      </ComContent>
    </ComLayout>
  );
};

export default Index;
