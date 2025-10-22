import { Form, Space } from 'antd';
import { useTranslate } from '@supos_host/hooks';
import DebounceSelect from './SearchSelect';
import { useEffect, useState } from 'react';
import { getInstanceInfo } from '@/apis';
import { ComSelect } from '@supos_host/components';
import { useBaseStore } from '@supos_host/baseStore';

import { REMOTE_NAME } from '../../variables';

const NameSpace = ({ isEdit }: any) => {
  const form = Form.useFormInstance();
  const [options, setOptions] = useState<any[]>();
  const formatMessage = useTranslate(REMOTE_NAME);
  const commonFormatMessage = useTranslate();
  const { qualityName = 'quality', timestampName = 'timeStamp' } = useBaseStore((state: any) => state.systemInfo);

  const onChange = (e: any) => {
    const _options = e?.option?.fields?.map((item: any) => {
      if ([qualityName, timestampName].includes(item.name)) {
        return { ...item, disabled: true };
      }
      return item;
    });
    setOptions(_options || []);
    form.setFieldValue(['refers', 0, 'field'], undefined);
  };

  useEffect(() => {
    if (isEdit) {
      getInstanceInfo({ id: form.getFieldValue(['refers', 0, 'refer', 'value']) }).then((res: any) => {
        const _options = res?.fields?.map((item: any) => {
          if ([qualityName, timestampName].includes(item.name)) {
            return { ...item, disabled: true };
          }
          return item;
        });
        setOptions(_options || []);
        form.setFieldValue(['refers', 0, 'refer', 'label'], res?.path);
      });
    }
  }, [isEdit]);

  return (
    <Form.Item label={formatMessage('key')} required>
      <Space
        style={{ width: '100%' }}
        styles={{
          item: { width: '50%', overflow: 'hidden' },
        }}
        align="start"
      >
        <Form.Item
          name={['refers', 0, 'refer']}
          style={{ marginBottom: 0 }}
          rules={[
            {
              required: true,
              message: commonFormatMessage('uns.pleaseInputNamespace'),
            },
          ]}
        >
          <DebounceSelect
            style={{ width: '100%' }}
            placeholder={commonFormatMessage('uns.namespace')}
            onChange={onChange}
            popupMatchSelectWidth={400}
            apiParams={{ type: 4 }}
            labelInValue
          />
        </Form.Item>
        <Form.Item
          name={['refers', 0, 'field']}
          style={{ marginBottom: 0 }}
          rules={[
            {
              required: true,
              message: commonFormatMessage('uns.pleaseSelectKeyType'),
            },
          ]}
        >
          <ComSelect
            fieldNames={{ label: 'name', value: 'name' }}
            placeholder={commonFormatMessage('uns.key')}
            options={options}
            style={{ width: '100%' }}
            allowClear
          />
        </Form.Item>
      </Space>
    </Form.Item>
  );
};

export default NameSpace;
