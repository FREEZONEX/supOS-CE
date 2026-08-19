import CollectionFlow from '@/pages/collection-flow';
import EventFlow from '@/pages/event-flow';
import ComLayout from '@/components/com-layout';
import { PageTitleIcon, SourceFlowIcon, EventFlowIcon } from '@/components/lucide-icon';
import { Add } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import ViewModeSegmented from '@/components/lucide-icon/ViewModeSegmented';
import { AuthButton } from '@/components';
import { ButtonPermission } from '@/common-types/button-permission';
import { useTranslate, useViewModeStorage, VIEW_MODE_STORAGE_KEYS } from '@/hooks';
import classNames from 'classnames';
import { Monitor } from 'lucide-react';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Navigate, useSearchParams } from 'react-router';
import AddFlowModal, { type AddFlowModalRef } from './components/AddFlowModal';
import type { FlowListPanelHandle, FlowTabType } from './types';
import styles from './index.module.scss';

const normalizeTab = (value: string | null): FlowTabType => {
  if (value === 'source' || value === 'event' || value === 'all') {
    return value;
  }
  return 'all';
};

const FlowPage: FC = () => {
  const formatMessage = useTranslate();
  const [searchParams, setSearchParams] = useSearchParams();
  const tab = normalizeTab(searchParams.get('tab'));
  const [toolbarHost, setToolbarHost] = useState<HTMLElement | null>(null);
  const addFlowModalRef = useRef<AddFlowModalRef>(null);
  const panelRef = useRef<FlowListPanelHandle>(null);

  useEffect(() => {
    setToolbarHost(document.getElementById('flow-subheader-actions'));
  }, [tab]);

  const openAddFlow = useCallback(() => {
    const context = panelRef.current?.getCreateContext();
    addFlowModalRef.current?.onOpen({
      flowKind: tab === 'event' ? 'event' : 'source',
      groupId: context?.groupId,
      lockFlowKind: tab !== 'all',
    });
  }, [tab]);

  const onFlowCreated = useCallback(() => {
    panelRef.current?.refreshRequest();
  }, []);

  const [viewMode, setViewMode] = useViewModeStorage(VIEW_MODE_STORAGE_KEYS.flow);

  useEffect(() => {
    const readStoredViewMode = (key: string) => {
      const raw = localStorage.getItem(key);
      if (!raw) return null;
      try {
        const parsed = JSON.parse(raw);
        return typeof parsed === 'string' ? parsed : null;
      } catch {
        return raw;
      }
    };

    if (readStoredViewMode(VIEW_MODE_STORAGE_KEYS.flow)) return;

    const legacyViewMode =
      readStoredViewMode(VIEW_MODE_STORAGE_KEYS.flowAll) ??
      readStoredViewMode(VIEW_MODE_STORAGE_KEYS.collection) ??
      readStoredViewMode(VIEW_MODE_STORAGE_KEYS.eventFlow);
    if (legacyViewMode) {
      setViewMode(legacyViewMode);
    }
  }, [setViewMode]);

  const setTab = useCallback(
    (nextTab: FlowTabType) => {
      const next = new URLSearchParams(searchParams);
      next.set('tab', nextTab);
      setSearchParams(next, { replace: true });
    },
    [searchParams, setSearchParams]
  );

  const tabs = useMemo(
    () => [
      {
        key: 'all' as const,
        label: formatMessage('common.all'),
        icon: <Monitor size={16} strokeWidth={1.75} />,
      },
      {
        key: 'source' as const,
        label: formatMessage('home.sourceFlow'),
        icon: <SourceFlowIcon />,
      },
      {
        key: 'event' as const,
        label: formatMessage('home.eventFlow'),
        icon: <EventFlowIcon />,
      },
    ],
    [formatMessage]
  );

  return (
    <ComLayout className={styles.flowPage}>
      <div className={styles.pageColumn}>
        <div className={styles.subHeader}>
          <div className={styles.pageTitleWrap}>
            <PageTitleIcon resourceKey="flow.collection.page" />
            <span className={styles.pageTitle}>{formatMessage('common.flow', {}, 'Flow')}</span>
          </div>
          <div className={styles.subHeaderActions}>
            <div id="flow-subheader-actions" className={styles.subHeaderPortal} />
            <AuthButton
              auth={[ButtonPermission['SourceFlow.add'], ButtonPermission['EventFlow.add']]}
              type="primary"
              icon={<Add {...toolbarIconProps} />}
              onClick={openAddFlow}
            >
              {formatMessage('common.newFlow')}
            </AuthButton>
          </div>
        </div>

        <div className={styles.globalHeader}>
          <div className={styles.globalTabs}>
            {tabs.map((item) => (
              <button
                key={item.key}
                type="button"
                className={classNames(
                  styles.globalTab,
                  tab === item.key && styles.globalTabActive,
                  tab !== item.key && styles.globalTabInactive
                )}
                onClick={() => setTab(item.key)}
              >
                <span className={styles.globalTabInner}>
                  {item.icon}
                  <span>{item.label}</span>
                </span>
                {tab === item.key ? <span className={styles.globalTabInk} aria-hidden /> : null}
              </button>
            ))}
          </div>
          <div className={styles.viewModeWrap}>
            <ViewModeSegmented
              value={viewMode}
              onChange={setViewMode}
              cardTitle={formatMessage('common.cardMode')}
              listTitle={formatMessage('common.listMode')}
            />
          </div>
        </div>

        <div className={styles.pageBody}>
          <div className={styles.main}>
            <div className={styles.contentFrame}>
              <div className={styles.panelHost} key={tab}>
                {tab === 'all' ? (
                  <CollectionFlow
                    ref={panelRef}
                    embedded
                    hideNewFlow
                    flowScope="all"
                    viewMode={viewMode}
                    onViewModeChange={setViewMode}
                    hideViewMode
                    toolbarPortalHost={toolbarHost}
                  />
                ) : null}
                {tab === 'source' ? (
                  <CollectionFlow
                    ref={panelRef}
                    embedded
                    hideNewFlow
                    flowScope="source"
                    viewMode={viewMode}
                    onViewModeChange={setViewMode}
                    hideViewMode
                    toolbarPortalHost={toolbarHost}
                  />
                ) : null}
                {tab === 'event' ? (
                  <EventFlow
                    ref={panelRef}
                    embedded
                    hideNewFlow
                    viewMode={viewMode}
                    onViewModeChange={setViewMode}
                    hideViewMode
                    toolbarPortalHost={toolbarHost}
                  />
                ) : null}
              </div>
            </div>
          </div>
        </div>
      </div>
      <AddFlowModal ref={addFlowModalRef} onSuccess={onFlowCreated} />
    </ComLayout>
  );
};

export const SourceFlowRedirect: FC = () => <Navigate to="/flow?tab=source" replace />;

export const EventFlowRedirect: FC = () => <Navigate to="/flow?tab=event" replace />;

export default FlowPage;
