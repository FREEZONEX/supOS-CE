import { useState, type FC, useEffect } from 'react';
import { Select, App } from 'antd';
import { usePropsValue, useTranslate } from '@/hooks';
import type { SelectProps } from 'antd';
import { getAllLabel } from '@/apis/inter-api/uns';

interface TagSelectProps extends SelectProps {
  tagMaxLen?: number;
}

const TagSelect: FC<TagSelectProps> = ({ value, onChange, tagMaxLen, ...rest }) => {
  const [options, setOptions] = useState<any>([]);
  const [val, setVal] = usePropsValue({
    value,
    onChange,
  });

  const { message } = App.useApp();
  const formatMessage = useTranslate();

  const searchAllTags = async () => {
    const res = await getAllLabel();
    if (res && Array.isArray(res)) {
      setOptions(res);
    }
  };

  const handleChange = (newValues: any) => {
    let filteredValues = newValues.slice();

    if (tagMaxLen && newValues.some((item: any) => (item.label || item.value)?.length > tagMaxLen)) {
      filteredValues = filteredValues.filter((item: any) => (item.label || item.value)?.length <= tagMaxLen);
      message.error(formatMessage('uns.labelMaxLength', { label: formatMessage('common.label'), length: tagMaxLen }));
    }

    if (newValues.some((item: any) => (item.label || item.value)?.trim() === '')) {
      filteredValues = filteredValues.filter((item: any) => (item.label || item.value)?.trim() !== '');
      message.error(formatMessage('common.prohibitSpacesTip'));
    }

    const reg = /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/;
    if (newValues.some((item: any) => !reg.test(item.label || item.value))) {
      filteredValues = filteredValues.filter((item: any) => reg.test(item.label || item.value));
      message.error(formatMessage('uns.nameFormat'));
    }

    setVal(filteredValues);
    onChange?.(filteredValues);
  };

  useEffect(() => {
    searchAllTags();
  }, []);

  return (
    <Select
      {...rest}
      showSearch
      filterOption
      options={options}
      value={val}
      onChange={handleChange}
      mode="tags"
      allowClear
      fieldNames={{
        label: 'labelName',
        value: 'id',
      }}
      labelInValue
      optionFilterProp="labelName"
    />
  );
};

export default TagSelect;
