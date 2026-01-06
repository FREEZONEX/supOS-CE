import type { FC } from 'react';
import { Form, Input, InputNumber, Select } from 'antd';
import { useTranslate } from '@/hooks';
import FunctionList from './FunctionList';

interface AggregateWindowFormProps {
  windowType?: string;
}

const AggForm: FC<AggregateWindowFormProps> = () => {
  const formatMessage = useTranslate();
  const form = Form.useFormInstance();
  const windowType = Form.useWatch(['streamOptions', 'window', 'windowType']);
  const whereFieldList = Form.useWatch('whereFieldList', form) || [];

  const renderWindowTypeFields = () => {
    switch (windowType) {
      case 'INTERVAL':
        return (
          <>
            <Form.Item
              name={['streamOptions', 'window', 'options', 'intervalValue']}
              label={formatMessage('streams.intervalValue')}
              rules={[{ required: true }]}
            >
              <Input placeholder={`${formatMessage('uns.eg')}: 1m, 5s, 1h`} />
            </Form.Item>
            <Form.Item
              name={['streamOptions', 'window', 'options', 'intervalOffset']}
              label={formatMessage('streams.intervalOffset')}
            >
              <Input placeholder={`${formatMessage('uns.optional')}, ${formatMessage('uns.eg')}: 30s`} />
            </Form.Item>
          </>
        );

      case 'SESSION':
        return (
          <Form.Item
            name={['streamOptions', 'window', 'options', 'tolValue']}
            label={formatMessage('streams.tolValue')}
            rules={[{ required: true }]}
          >
            <Input placeholder={`${formatMessage('uns.eg')}: 10s`} />
          </Form.Item>
        );

      case 'STATE_WINDOW':
        return (
          <Form.Item
            name={['streamOptions', 'window', 'options', 'field']}
            label={formatMessage('streams.referenceField')}
            rules={[{ required: true }]}
          >
            <Select
              options={whereFieldList.filter((field: any) =>
                ['integer', 'long', 'boolean', 'string'].includes(field.type)
              )}
            />
          </Form.Item>
        );

      case 'EVENT_WINDOW':
        return (
          <>
            <Form.Item
              name={['streamOptions', 'window', 'options', 'startWith']}
              label={formatMessage('streams.startWith')}
              rules={[{ required: true }]}
            >
              <Input placeholder={`${formatMessage('uns.eg')}: temperature > 90`} />
            </Form.Item>
            <Form.Item
              name={['streamOptions', 'window', 'options', 'endWith']}
              label={formatMessage('streams.endWith')}
              rules={[{ required: true }]}
            >
              <Input placeholder={`${formatMessage('uns.eg')}: temperature < 85`} />
            </Form.Item>
          </>
        );

      case 'COUNT_WINDOW':
        return (
          <>
            <Form.Item
              name={['streamOptions', 'window', 'options', 'countValue']}
              label={formatMessage('streams.countValue')}
              rules={[{ required: true }]}
            >
              <InputNumber
                min={2}
                step={1}
                max={Math.pow(2, 31) - 1}
                placeholder="[2,2^31-1]"
                style={{ width: '100%' }}
              />
            </Form.Item>
            <Form.Item
              name={['streamOptions', 'window', 'options', 'slidingValue']}
              label={formatMessage('streams.slidingValue')}
              rules={[
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    console.log(value);
                    if (!value || getFieldValue(['streamOptions', 'window', 'options', 'countValue']) > value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error(formatMessage('streams.countSlidingLessThanCountValue')));
                  },
                }),
              ]}
            >
              <InputNumber min={1} step={1} placeholder={formatMessage('uns.optional')} style={{ width: '100%' }} />
            </Form.Item>
          </>
        );

      default:
        return null;
    }
  };

  return (
    <>
      <FunctionList />
      <Form.Item
        label={formatMessage('streams.windowType')}
        name={['streamOptions', 'window', 'windowType']}
        style={{ marginBottom: 20 }}
        rules={[{ required: true }]}
      >
        <Select
          onChange={(val) => {
            form.setFieldValue(['streamOptions', 'window', 'options'], undefined);
            form.setFieldValue('advancedOptions', val === 'COUNT_WINDOW');
          }}
        >
          <Select.Option value="INTERVAL">INTERVAL</Select.Option>
          <Select.Option value="SESSION">SESSION</Select.Option>
          <Select.Option value="STATE_WINDOW">STATE WINDOW</Select.Option>
          <Select.Option value="EVENT_WINDOW">EVENT WINDOW</Select.Option>
          <Select.Option value="COUNT_WINDOW">COUNT WINDOW</Select.Option>
        </Select>
      </Form.Item>
      {windowType && renderWindowTypeFields()}
    </>
  );
};

export default AggForm;
