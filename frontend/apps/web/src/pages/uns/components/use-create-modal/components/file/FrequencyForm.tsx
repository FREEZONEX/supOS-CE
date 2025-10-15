import type { FC } from 'react';
import { Form, Space, InputNumber, Select } from 'antd';
import { useTranslate } from '@/hooks';

interface FrequencyFormProps {
  unitList?: string[];
  maxValue?: number;
}

const FrequencyForm: FC<FrequencyFormProps> = ({ unitList = ['s', 'm', 'h'], maxValue }) => {
  const formatMessage = useTranslate();
  const commonRules = (message: string) => {
    return {
      rules: [
        {
          required: true,
          message,
        },
      ],
    };
  };

  const options = [
    { value: 's', label: formatMessage('uns.second') },
    { value: 'm', label: formatMessage('uns.minute') },
    { value: 'h', label: formatMessage('uns.hour') },
    { value: 'd', label: formatMessage('uns.day') },
  ];
  return (
    <Space.Compact block style={{ width: '100%' }}>
      <Form.Item name={['frequency', 'value']} {...commonRules(formatMessage('uns.pleaseInputValue'))} noStyle>
        <InputNumber
          placeholder={formatMessage('common.commonPlaceholder')}
          style={{ width: '50%' }}
          min={1}
          step="1"
          precision={0}
          max={maxValue}
        />
      </Form.Item>
      <Form.Item name={['frequency', 'unit']} {...commonRules(formatMessage('uns.pleaseSelectUnit'))} noStyle>
        <Select
          placeholder={formatMessage('common.select')}
          style={{ width: '50%' }}
          options={options.filter((e) => unitList.includes(e.value))}
        />
      </Form.Item>
    </Space.Compact>
  );
};

export default FrequencyForm;
