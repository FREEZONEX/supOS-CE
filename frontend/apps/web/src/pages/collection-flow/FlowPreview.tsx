import {
  deployFlow,
  getFlowDetail,
  saveFlow,
} from '@/apis/core-api/flow';
import type { PageProps } from '@/common-types';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { AuthButton } from '@/components/auth';
import ComBackButton from '@/components/com-back-button';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComText from '@/components/com-text';
import { postFlowIframeTheme, useFlowIframeThemeBridge, useLocalStorage, useTabName, useTranslate } from '@/hooks';
import { hasPermission } from '@/utils/auth';
import { getDevProxyBaseUrl, getSearchParamsObj, getSearchParamsString } from '@/utils/url-util';
import { ChevronLeft, OverflowMenuVertical } from '@/components/lucide-icon/carbon';
import { useUpdateEffect } from 'ahooks';
import { Breadcrumb, Button, Dropdown, Flex, message, Space } from 'antd';
import { type FC, useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router';
import './index.scss';

type NodeRedWorkspaceState = {
  activeWorkspaceId?: string;
  activeWorkspaceName?: string;
  activeWorkspaceType?: string;
  rootWorkspaceId?: string;
  rootWorkspaceName?: string;
  isSubflow?: boolean;
};

const FLOW_TAB_MAX_NAME_LENGTH = 18;
const truncateFlowTabName = (value: string) =>
  value.length > FLOW_TAB_MAX_NAME_LENGTH ? `${value.slice(0, FLOW_TAB_MAX_NAME_LENGTH)}...` : value;

const FlowPreview: FC<PageProps> = ({ location }) => {
  const state = getSearchParamsObj(location?.search) || {};
  const navigate = useNavigate();
  const formatMessage = useTranslate();
  const [flowName, setFlowName] = useState(state.name || '');
  const tabRawName = String(flowName || state.name || '').trim();
  const flowTabPrefix = `${formatMessage('common.flow')} · `;
  useTabName(tabRawName, {
    displayName: tabRawName ? `${flowTabPrefix}${truncateFlowTabName(tabRawName)}` : undefined,
    fullName: tabRawName ? `${flowTabPrefix}${tabRawName}` : undefined,
  });
  const nodeRedLang = useLocalStorage('editor-language');
  const flowIframeRef = useRef<HTMLIFrameElement | null>(null);
  const [loading, setLoading] = useState(true);
  const [nodeRedWorkspace, setNodeRedWorkspace] = useState<NodeRedWorkspaceState | null>(null);
  const observerRef = useRef<any>(null);
  // const [buttonDisabled, setDisabled] = useState(state?.status === 'RUNNING');
  const loadRef = useRef(false);
  // iframe的key

  const [key, setKey] = useState(() => Date.now());
  const iframeTheme = useFlowIframeThemeBridge(flowIframeRef);
  const iframeThemeRef = useRef(iframeTheme);
  const [initialIframeTheme] = useState(iframeTheme);
  const iframeParams = new URLSearchParams({
    flowId: String(state.id || ''),
    rootWorkspaceId: String(state.flowId || ''),
    theme: initialIframeTheme,
  });
  const runtimeFlowId = String(state.flowId || '');
  const iframeHash = runtimeFlowId ? `#flow/${encodeURIComponent(runtimeFlowId)}` : '';
  const iframeUrl = `/nodered/home/?${iframeParams.toString()}${iframeHash}`;
  const isInSubflow = !!nodeRedWorkspace?.isSubflow;
  useEffect(() => {
    iframeThemeRef.current = iframeTheme;
  }, [iframeTheme]);

  const onShowRootWorkspace = () => {
    flowIframeRef.current?.contentWindow?.postMessage(
      { type: 'nodeRedShowWorkspace', data: { id: nodeRedWorkspace?.rootWorkspaceId } },
      '*'
    );
  };
  const breadcrumbList = [
    {
      name: flowName || state.name,
      onClick: isInSubflow ? onShowRootWorkspace : undefined,
    },
    ...(isInSubflow
      ? [
          {
            name: nodeRedWorkspace?.activeWorkspaceName || formatMessage('flowEditor.subflow'),
          },
        ]
      : []),
  ];

  useUpdateEffect(() => {
    if (!loadRef.current) return;
    if (nodeRedLang) {
      setKey(Date.now());
    }
  }, [nodeRedLang]);

  useEffect(() => {
    if (!state?.id) return;
    getFlowDetail(state.id)
      .then((flow) => {
        if (flow?.flowName) {
          setFlowName(flow.flowName);
        }
      })
      .catch(() => {
        setFlowName(state.name || '');
      });
  }, [state?.id, state.name]);
  // 将 flows 数据保存到后端
const saveFlowsToBackend = async (data: any) => {
  try {
    const { flows, type } = data;
    const filteredFlows = flows?.filter(
      (item: any) => item.type !== 'tab' && item._contextOnly !== true && item._contextOnly !== 'true'
    );
    const api = type === 'save' ? saveFlow : deployFlow;
    setLoading(true);
    const result: any = await api({ flows: filteredFlows, id: state?.id });
    if (type === 'deploy') {
      const flowId = result?.flowId;
      if (!state.flowId && flowId) {
        navigate(`/collection-flow/flow-editor?${getSearchParamsString({ ...state, flowId })}`, { replace: true });
      }
      setKey(Date.now());
      message.success(formatMessage('appGui.deployOk'));
    } else {
      message.success(formatMessage('appGui.saveSuccess'));
    }
  } catch (error) {
    console.error('Error saving flows:', error);
  } finally {
    setLoading(false);
  }
}

  // 监听 iframe 加载
  useUpdateEffect(() => {
    const iframe = flowIframeRef.current;
    if (!iframe) return;

    const handleLoad = () => {
      postFlowIframeTheme(iframe, iframeThemeRef.current);
      setLoading(false);
    };

    iframe.addEventListener('load', handleLoad);
    return () => iframe.removeEventListener('load', handleLoad);
  }, [key]); // 依赖 iframeKey 的变化（重新加载时触发）

  useEffect(() => {
    setLoading(true);
    // 监听来自 Node-RED 的 flows 数据
    const handleMessage = (event: any) => {
      if (!flowIframeRef.current || event.source !== flowIframeRef.current.contentWindow) return;
      if (event.data.type === 'currentFlows') {
        saveFlowsToBackend(event.data.data);
      } else if (event.data.type === 'flowsChange') {
        // setDisabled(!event.data?.data?.contentsChanged);
      } else if (event.data.type === 'nodeRedWorkspaceState') {
        setNodeRedWorkspace(event.data.data || null);
      }
    };

    const loadFn = () => {
      loadRef.current = true;
      postFlowIframeTheme(flowIframeRef.current, iframeThemeRef.current);
      setLoading(false);
    };
    if (flowIframeRef.current) {
      flowIframeRef.current.addEventListener('load', loadFn);
    }
    window.addEventListener('message', handleMessage);

    return () => {
      window.removeEventListener('message', handleMessage);
      if (flowIframeRef.current) {
        // eslint-disable-next-line @typescript-eslint/no-unused-expressions
        flowIframeRef.current && flowIframeRef.current.removeEventListener('load', loadFn);
      }
    };
  }, [state?.id, state?.flowId]);

  useEffect(() => {
    setNodeRedWorkspace(null);
  }, [key, state?.id, state?.flowId]);

  const setPostMessage = (type: string) => {
    if (flowIframeRef.current) {
      setLoading(true);
      flowIframeRef.current.contentWindow!.postMessage({ data: type, type: 'requestFlows' }, '*');
    }
  };

  // 点击按钮请求 flows 数据
  const onDeployFlows = () => {
    setPostMessage('deploy');
  };

  const onOpenMenuHandle = (id: string) => {
    if (flowIframeRef.current) {
      flowIframeRef.current.contentWindow!.postMessage({ data: { id }, type: 'openMenu' }, '*');
    }
  };

  useEffect(() => {
    const targetId = 'red-ui-loading-progress';
    const iframe = flowIframeRef.current;
    const handleVisibilityChange = (isVisible: boolean) => {
      setLoading(isVisible);
    };
    const setupObserver = () => {
      if (!iframe) return;
      try {
        // 获取iframe的document对象
        const iframeDoc = iframe.contentDocument || iframe.contentWindow?.document;

        // 创建文档级观察器
        const docObserver = new MutationObserver((_, observer) => {
          const targetElement = iframeDoc?.getElementById(targetId);
          if (targetElement) {
            // 找到元素后立即断开文档观察器
            observer.disconnect();
            // 初始化样式观察器
            const styleObserver = new MutationObserver((mutations) => {
              mutations.forEach((mutation) => {
                if (mutation.attributeName === 'style') {
                  handleVisibilityChange((mutation.target as HTMLElement).style.display !== 'none');
                }
              });
            });

            styleObserver.observe(targetElement, {
              attributes: true,
              attributeFilter: ['style'],
            });

            // 存储观察器引用
            observerRef.current = styleObserver;

            // 初始状态检查
            handleVisibilityChange(targetElement.style.display !== 'none');
          }
        });
        // 监听整个文档的DOM变化
        docObserver.observe(iframeDoc!, {
          childList: true,
          subtree: true,
        });
        handleVisibilityChange(true);
      } catch (error) {
        console.error('访问iframe内容出错:', error);
      }
    };
    if (iframe) {
      iframe.addEventListener('load', setupObserver);
    }
    return () => {
      if (iframe) {
        iframe.removeEventListener('load', setupObserver);
      }
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [key, state?.id, state?.flowId]);

  const items: any = [
    {
      key: 'menu-item-import',
      auth: ButtonPermission['SourceFlow.import'],
      label: formatMessage('common.import'),
    },
    {
      key: 'menu-item-export',
      auth: ButtonPermission['SourceFlow.export'],
      label: formatMessage('uns.export'),
    },
    {
      type: 'divider',
    },
    {
      key: 'menu-item-search',
      // auth: ButtonPermission['SourceFlow.process'],
      label: formatMessage('flowEditor.process'),
    },
    {
      type: 'divider',
    },
    // {
    //   key: 'config-nodes',
    //   label: <span onClick={() => onOpenMenuHandle('menu-item-config-nodes')}>修改节点配置</span>,
    // },
    // {
    //   type: 'divider',
    // },
    {
      key: 'menu-item-edit-palette',
      auth: ButtonPermission['SourceFlow.nodeManagement'],
      label: formatMessage('flowEditor.nodeManagement'),
    },
    // {
    //   type: 'divider',
    // },
    // {
    //   key: 'menu-item-user-settings',
    //   label: <span>设置</span>,
    // },
  ]
    ?.filter((i) => i.type === 'divider' || i.label)
    ?.filter((f) => {
      return !f.auth || hasPermission(f.auth);
    });
  const handleBack = () => {
    const fromPath = state.from;
    if (fromPath) {
      navigate(fromPath);
    } else {
      navigate(-1);
    }
  };
  return (
    <ComLayout loading={loading}>
      <ComContent
        mustHasBack={false}
        style={{ overflow: 'hidden' }}
        hasPadding
        border={false}
        title={
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <ComBackButton onClick={handleBack} />
              <Breadcrumb
                separator=">"
                items={breadcrumbList?.map((item: any, idx: number) => {
                  if (idx + 1 === breadcrumbList?.length) {
                    return {
                      title: item.name,
                    };
                  }
                  return {
                    title: <ComText>{item.name}</ComText>,
                    onClick: () => {
                      if (item.onClick) {
                        item.onClick();
                        return;
                      }
                      if (!item.path) return;
                      navigate(item.path);
                    },
                  };
                })}
              />
              {isInSubflow && (
                <Button
                  variant="outlined"
                  color="default"
                  style={{ paddingLeft: '5.5px', gap: '3px' }}
                  onClick={onShowRootWorkspace}
                >
                  <Flex align="center" gap={8}>
                    <ChevronLeft size={16} />
                    {formatMessage('flowEditor.returnParent')}
                  </Flex>
                </Button>
              )}
            </div>
            <Space>
              <AuthButton
                auth={ButtonPermission['SourceFlow.deploy']}
                loading={loading}
                type="primary"
                onClick={onDeployFlows}
                // disabled={buttonDisabled}
              >
                {formatMessage('appGui.deploy')}
              </AuthButton>
              <Dropdown
                menu={{
                  onClick: (e) => {
                    onOpenMenuHandle(e.key);
                  },
                  items: items,
                }}
                placement="bottomRight"
              >
                <div className="flow-dropdown-more">
                  <OverflowMenuVertical size={16} />
                </div>
              </Dropdown>
            </Space>
          </div>
        }
      >
        <iframe
          key={key}
          ref={flowIframeRef}
          style={{
            width: '100%',
            height: '100%',
            border: 'none',
          }}
          title={'Node-RED'}
          src={`${getDevProxyBaseUrl()}${iframeUrl}`}
        />
      </ComContent>
    </ComLayout>
  );
};

export default FlowPreview;
