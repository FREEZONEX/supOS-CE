import { type FC, useState, useEffect, useRef } from 'react';
import { Button, Space, Breadcrumb, Flex } from 'antd';
import { useNavigate } from 'react-router';
import type { PageProps } from '@/common-types';
import { getDashboardDetail } from '@/apis/inter-api/uns';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { usePrevious } from 'ahooks';
import { useTranslate } from '@/hooks';
import ComText from '@/components/com-text';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { getSearchParamsObj } from '@/utils/url-util';
import { AuthButton } from '@/components/auth';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { ChevronLeft } from '@carbon/icons-react';

const FlowPreview: FC<PageProps> = ({ location }) => {
  const formatMessage = useTranslate();
  const [iframeUrl, setIframeUrl] = useState('');
  const state = getSearchParamsObj(location?.search) || {};
  const breadcrumbList = [
    {
      name: state.name,
    },
  ];
  const navigate = useNavigate();
  const timer: any = useRef(null);
  const { id, type, status } = state;
  const previous = usePrevious(state);

  const getIframeUrl = (isFirst = false) => {
    if (!isFirst && previous?.id === id && previous?.status === status) return;
    setIframeUrl('');
    if (id) {
      //fuxa
      if ([2, '2'].includes(type)) {
        setIframeUrl(`/fuxa/home/?id=${id}&status=${status === 'design' ? 'editor' : 'lab'}`);
        return;
      }

      // grafana
      getDashboardDetail(id).then((res: any) => {
        if (res?.meta?.url) {
          setIframeUrl(`${res?.meta?.url}${status === 'design' ? '' : '?kiosk'}`);
        }
      });
    }
  };

  useActivate(() => {
    getIframeUrl();
  });

  useEffect(() => {
    getIframeUrl(true);
  }, []);

  useEffect(() => {
    if (iframeUrl) {
      localStorage.setItem('SearchBar_Hidden', 'true');
      const iframe: any = document?.getElementById('dashboardIframe');
      if (status === 'design' && [1, '1'].includes(type) && iframe) {
        iframe.onload = function () {
          console.log('iframe加载完成');
          timer.current = setInterval(() => {
            const megaMenuToggle = iframe?.contentWindow?.document?.querySelector('#mega-menu-toggle');
            const breadcrumbs = iframe?.contentWindow?.document?.querySelector('[aria-label="Breadcrumbs"]');
            const kioskModeBtn =
              iframe?.contentWindow?.document?.querySelector('[title="Enable kiosk mode"]') ||
              iframe?.contentWindow?.document?.querySelector('[title="启用 kiosk 模式"]');
            const bar =
              iframe?.contentWindow?.document?.querySelector('[title="Toggle top search bar"]') ||
              iframe?.contentWindow?.document?.querySelector('[title="切换顶部搜索栏"]');
            if (megaMenuToggle) {
              try {
                // 隐藏元素
                megaMenuToggle.style.display = 'none';
                megaMenuToggle.style.pos = 'none';
                // 禁用事件监听器
                megaMenuToggle.addEventListener('click', function (event: any) {
                  event.stopPropagation();
                  event.preventDefault();
                });

                breadcrumbs.style.display = 'none';
                // 禁用事件监听器
                breadcrumbs.addEventListener('click', function (event: any) {
                  event.stopPropagation();
                  event.preventDefault();
                });
                kioskModeBtn.style.display = 'none';
                // 禁用事件监听器
                kioskModeBtn.addEventListener('click', function (event: any) {
                  event.stopPropagation();
                  event.preventDefault();
                });

                bar.style.display = 'none';
                // 禁用事件监听器
                bar.addEventListener('click', function (event: any) {
                  event.stopPropagation();
                  event.preventDefault();
                });

                clearInterval(timer.current);
              } catch (err) {
                console.error(err);
                clearInterval(timer.current);
              }
            }
          }, 10);
        };
      }
    }
    return () => {
      clearInterval(timer.current);
    };
  }, [iframeUrl, status]);

  const handleClick = (type: string) => {
    if (!(document.getElementById('dashboardIframe') as HTMLIFrameElement)?.contentWindow?.postMessage) return;

    (document.getElementById('dashboardIframe') as HTMLIFrameElement).contentWindow?.postMessage({
      from: 'supos',
      type,
    });
  };

  return (
    <ComLayout>
      <ComContent
        mustHasBack={false}
        hasPadding
        title={
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <Button
                variant="outlined"
                color="default"
                style={{ paddingLeft: '5.5px', gap: '3px' }}
                onClick={() => {
                  navigate('/dashboards');
                }}
              >
                <Flex align="center" gap={8}>
                  <ChevronLeft size={16} />
                  {formatMessage('common.back')}
                </Flex>
              </Button>
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
                      if (!item.path) return;
                      navigate(item.path);
                    },
                  };
                })}
              />
            </div>
            {[2, '2'].includes(type) && (
              <div>
                <Space>
                  {status === 'design' && (
                    <>
                      <AuthButton auth={ButtonPermission['Dashboards.save']} onClick={() => handleClick('save')}>
                        {formatMessage('common.save')}
                      </AuthButton>
                      <AuthButton auth={ButtonPermission['Dashboards.export']} onClick={() => handleClick('export')}>
                        {formatMessage('uns.export')}
                      </AuthButton>
                      <AuthButton auth={ButtonPermission['Dashboards.import']} onClick={() => handleClick('import')}>
                        {formatMessage('common.import')}
                      </AuthButton>
                    </>
                  )}
                  <Button onClick={() => handleClick('share')}>{formatMessage('common.share')}</Button>
                </Space>
              </div>
            )}
          </div>
        }
      >
        <iframe
          key={iframeUrl ?? '-1'}
          id="dashboardIframe"
          src={iframeUrl}
          style={{ width: '100%', height: '100%', display: 'block' }}
        />
      </ComContent>
    </ComLayout>
  );
};

export default FlowPreview;
