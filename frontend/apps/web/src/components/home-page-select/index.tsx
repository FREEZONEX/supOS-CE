import { getLaunchpadProjectsApi } from '@/apis/core-api/launchpad';
import { useTranslate } from '@/hooks';
import type { ResourceProps } from '@/stores/types';
import { type SelectProps } from 'antd';
import { type FC, useEffect, useMemo, useState } from 'react';
import ComSelect from '../com-select';

type HomePageSelectProps = Omit<SelectProps, 'options'> & {
  enabled?: boolean;
  resources?: ResourceProps[];
  targetUserId?: string | number;
};

type HomePageOption = {
  label: string;
  value: string;
};

type AppOptionState = {
  key: string;
  options: HomePageOption[];
};

const resourceHomePage = (item: ResourceProps) => {
  const url = String(item.url || '').trim();
  if (item.urlType === 1 && url.startsWith('/') && !url.startsWith('//')) {
    return url;
  }
  const code = String(item.code || '').trim();
  return code ? `/${code}` : '';
};

const HomePageSelect: FC<HomePageSelectProps> = ({ enabled = true, resources = [], targetUserId, ...props }) => {
  const formatMessage = useTranslate();
  const requestKey = enabled ? String(targetUserId ?? 'self') : '';
  const [appOptionState, setAppOptionState] = useState<AppOptionState>({ key: '', options: [] });

  useEffect(() => {
    if (!enabled) {
      return;
    }

    let active = true;
    getLaunchpadProjectsApi(undefined, targetUserId)
      .then((projects) => {
        if (!active) return;
        setAppOptionState({
          key: requestKey,
          options: projects.flatMap((project) =>
            (project.apps || []).map((app) => ({
              label: `${project.displayName || project.projectName} / ${app.displayName || app.appName}`,
              value: `/launchpad/${encodeURIComponent(String(project.projectId))}/${encodeURIComponent(String(app.appId))}`,
            }))
          ),
        });
      })
      .catch(() => {
        if (active) setAppOptionState({ key: requestKey, options: [] });
      });

    return () => {
      active = false;
    };
  }, [enabled, requestKey, targetUserId]);

  const options = useMemo(() => {
    const seen = new Set<string>();
    return [
      ...resources.map((item) => ({
        label: formatMessage(
          String(item.showName || item.code || ''),
          undefined,
          String(item.showName || item.code || '')
        ),
        value: resourceHomePage(item),
      })),
      ...(appOptionState.key === requestKey ? appOptionState.options : []),
    ].filter((option) => {
      if (!option.label || !option.value || seen.has(option.value)) return false;
      seen.add(option.value);
      return true;
    });
  }, [appOptionState, formatMessage, requestKey, resources]);

  return (
    <ComSelect
      {...props}
      options={options}
      filterOption={(input, option) =>
        String(option?.label || '')
          .toLowerCase()
          .includes(input.toLowerCase())
      }
    />
  );
};

export default HomePageSelect;
