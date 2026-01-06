import { AddLarge, Apps } from '@carbon/icons-react';
import { message, Spin, Typography, Button, Flex } from 'antd';
import AppSpaceList from '@/pages/app-management/components/AppSpaceList';
import useAddModal from '@/pages/app-management/components/useAddModal';
import { useNavigate } from 'react-router';
import AppUrlPreview from '@/pages/app-management/components/AppUrlPreview';
import AppEmpty from '@/pages/app-management/components/AppEmpty';
import { destroyApp, destroyHtml, getApps, getSingleApp, setHomepage } from '@/apis/inter-api/apps';
import { type FC, useEffect, useState } from 'react';
import { useDebounce, useUpdateEffect } from 'ahooks';
import type { PageProps } from '@/common-types';
import { useTranslate } from '@/hooks';
import { ButtonPermission } from '@/common-types/button-permission';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { AuthWrapper, AuthButton } from '@/components/auth';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ComLeft from '@/components/com-layout/ComLeft';
import ProSearch from '@/components/pro-search';
const { Title } = Typography;

const Module: FC<PageProps> = ({ location }) => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const { state } = location || {};
  const [list, setList] = useState([]);
  const [htmlList, setHtmlList] = useState<any>([]);
  const [htmlLoading, setHtmlLoading] = useState(false);
  const [searchValue, setSearchValue] = useState<string>('');
  const [clickItemId, setClickItemIdItemId] = useState(null);
  const debouncedSearchValue = useDebounce(searchValue, { wait: 300 });
  const [loading, setLoading] = useState(false);
  useActivate(() => {
    getAppsFn(debouncedSearchValue);
  });

  const getAppsFn = (k?: any) => {
    setLoading(true);
    getApps({ k })
      .then((data: any) => {
        setList(data || []);
      })
      .catch(() => {})
      .finally(() => {
        setLoading(false);
      });
  };
  useUpdateEffect(() => {
    getAppsFn(debouncedSearchValue);
  }, [debouncedSearchValue]);
  const { ModalDom, setModalOpen } = useAddModal({ successCallBack: getAppsFn });
  useEffect(() => {
    getAppsFn();
  }, []);
  const getHtmlList = (name?: any) => {
    setHtmlLoading(true);
    getSingleApp(name)
      .then((data: any) => {
        const urls = Object.entries(data?.urls)?.map((item) => ({
          id: item[0],
          url: item[1],
          isHomePage: data?.homepage === item[1],
        }));
        setHtmlList(urls);
      })
      .finally(() => {
        setHtmlLoading(false);
      });
  };

  const onSearchChange = (e: any) => {
    const value = e?.target?.value || '';
    setSearchValue(value);
  };
  const deleteHandle = (appName: string) => {
    destroyApp(appName).then(() => {
      if (appName === clickItemId) {
        setClickItemIdItemId(null);
        setHtmlList([]);
      }
      getAppsFn?.();
    });
  };
  useEffect(() => {
    if (state?.name) {
      setClickItemIdItemId(state?.name);
      getHtmlList?.(state?.name);
    }
  }, [state?.name]);
  return (
    <ComLayout loading={loading}>
      <ComLeft
        style={{ padding: '20px 10px 10px', backgroundColor: 'var(--supos-bg-color)' }}
        title={
          <Flex
            justify="flex-start"
            align="center"
            gap={5}
            style={{ height: '100%', color: 'var(--supos-text-color)' }}
          >
            <Apps size={19} style={{ color: 'var(--supos-text-color)' }} />
            {formatMessage('appSpace.tree')}
          </Flex>
        }
      >
        <Flex gap={4}>
          <ProSearch
            closeButtonLabelText={formatMessage('common.clearSearchInput')}
            id="search"
            placeholder={formatMessage('uns.inputText')}
            role="searchbox"
            size="md"
            type="text"
            style={{ border: 'none' }}
            onChange={onSearchChange}
            value={searchValue}
          />
          <AuthWrapper auth={ButtonPermission['appSpace.add']}>
            <Button
              style={{ borderRadius: 0, height: 40, padding: '0 20px' }}
              title={formatMessage('appSpace.add')}
              type="primary"
              onClick={() => setModalOpen(true)}
              icon={<AddLarge />}
            />
          </AuthWrapper>
        </Flex>
        <AppSpaceList
          list={list}
          getHtmlList={getHtmlList}
          clickItemId={clickItemId}
          setClickItemIdItemId={setClickItemIdItemId}
          deleteHandle={deleteHandle}
        />
        {ModalDom}
      </ComLeft>
      <ComContent
        title={
          <Flex
            style={{ height: '100%', backgroundColor: 'var(--supos-user-color)' }}
            align={'center'}
            justify={'space-between'}
          >
            <span style={{ fontSize: 18, fontWeight: 700 }}></span>
            <AuthButton
              auth={ButtonPermission['appSpace.newPage']}
              style={{ width: 102 }}
              type="primary"
              onClick={() => {
                navigate('/app-gui');
              }}
            >
              {formatMessage('appSpace.newPage')}
            </AuthButton>
          </Flex>
        }
      >
        <div style={{ height: '100%', padding: '40px 20px' }}>
          <Title style={{ fontWeight: 400 }} type="secondary" level={2}>
            {clickItemId}
          </Title>
          {htmlLoading ? (
            <Spin>loading...</Spin>
          ) : htmlList?.length > 0 ? (
            htmlList?.map((item: any) => (
              <AppUrlPreview
                item={item}
                key={item.id}
                editHandle={() =>
                  navigate('/app-preview', {
                    state: { htmlName: item.url, appName: clickItemId },
                  })
                }
                setHomepage={() => {
                  setHomepage({
                    htmlId: item.id,
                    appName: clickItemId,
                  }).then(() => {
                    message.success(formatMessage('common.settingSuccess'));
                    getHtmlList(clickItemId);
                  });
                }}
                deleteHandle={() => {
                  destroyHtml(clickItemId!, item.id).then(() => {
                    message.success(formatMessage('common.deleteHtmlSuccess'));
                    getHtmlList(clickItemId);
                  });
                }}
              />
            ))
          ) : (
            <AppEmpty appName={clickItemId} />
          )}
        </div>
      </ComContent>
    </ComLayout>
  );
};

export default Module;
