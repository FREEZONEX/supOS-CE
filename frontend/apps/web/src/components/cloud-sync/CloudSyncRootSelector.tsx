import { getTreeData } from '@/apis/core-api/uns';
import HelpTooltip from '@/components/help-tooltip';
import { useTranslate } from '@/hooks';
import { Flex, Select, Tag } from 'antd';
import type { CustomTagProps } from 'rc-select/lib/BaseSelect';
import { Route } from 'lucide-react';
import { useEffect, useState, type MouseEvent, type ReactNode } from 'react';
import styles from '@/pages/uns/components/uns-dashboard/cloudsync.module.scss';

type RootOption = {
  label: string;
  value: string;
};

type SelectRootOption = {
  label: ReactNode;
  value: string;
  searchLabel: string;
};

interface CloudSyncRootSelectorProps {
  value: string[];
  validationVisible: boolean;
  required?: boolean;
  onChange: (value: string[]) => void;
}

const CloudSyncRootSelector = ({ value, validationVisible, required = true, onChange }: CloudSyncRootSelectorProps) => {
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(true);
  const [options, setOptions] = useState<RootOption[]>([]);

  useEffect(() => {
    let active = true;
    void getTreeData()
      .then((nodes) => {
        if (!active) return;
        const nextOptions = (Array.isArray(nodes) ? nodes : [])
          .filter((node) => {
            const parentID = String(node?.parentId || '').trim();
            return Number(node?.pathType) === 0 && (!parentID || parentID === '0');
          })
          .map((node) => ({
            label: String(node?.name || node?.title || node?.path || node?.id),
            value: String(node?.id),
          }));
        setOptions(nextOptions);
      })
      .catch(() => {
        if (active) setOptions([]);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, []);

  const selectOptions: SelectRootOption[] = options.map((option) => ({
    value: option.value,
    searchLabel: option.label,
    label: (
      <span className={styles['cloudsync-path-option']}>
        <Route size={14} strokeWidth={1.75} aria-hidden />
        <span>{option.label}</span>
      </span>
    ),
  }));

  const error = required && validationVisible && value.length === 0;

  const tagRender = ({ value: tagValue, closable, onClose }: CustomTagProps) => {
    const onPreventMouseDown = (event: MouseEvent<HTMLSpanElement>) => {
      event.preventDefault();
      event.stopPropagation();
    };
    const text = options.find((option) => option.value === String(tagValue))?.label || String(tagValue || '');
    return (
      <Tag
        className={styles['cloudsync-path-tag']}
        onMouseDown={onPreventMouseDown}
        closable={closable}
        onClose={onClose}
        icon={<Route size={12} strokeWidth={1.75} aria-hidden />}
      >
        {text}
      </Tag>
    );
  };

  return (
    <Flex vertical gap={6}>
      <div className={styles['cloudsync-required-label']}>
        {required ? <span className={styles['cloudsync-required-mark']}>*</span> : null}
        <span>{formatMessage('uns.cloudSyncSelectedPaths')}</span>
        <HelpTooltip title={formatMessage('uns.cloudSyncSelectHint')} />
      </div>
      <Select
        mode="multiple"
        className={styles['cloudsync-path-select']}
        value={value}
        options={selectOptions}
        loading={loading}
        showSearch
        optionFilterProp="searchLabel"
        placeholder={formatMessage(required ? 'uns.cloudSyncSelectTopicRequired' : 'common.select')}
        tagRender={tagRender}
        onChange={onChange}
      />
      {error ? (
        <div className={styles['cloudsync-error']}>{formatMessage('uns.cloudSyncSelectTopicRequired')}</div>
      ) : null}
    </Flex>
  );
};

export default CloudSyncRootSelector;
