import { Form, InputNumber, Space } from 'antd';
import { useTranslate } from '@supos_host/hooks';
import { ComSelect } from '@supos_host/components';
import { REMOTE_NAME } from '../../variables';

const Condition = () => {
  const formatMessage = useTranslate(REMOTE_NAME);
  const commonFormatMessage = useTranslate();

  const options = [
    {
      label: formatMessage('greaterThan'),
      value: '>',
    },
    {
      label: formatMessage('greaterEqualThan'),
      value: '>=',
    },
    {
      label: formatMessage('lessThan'),
      value: '<',
    },
    {
      label: formatMessage('lessEqualThan'),
      value: '<=',
    },
    {
      label: formatMessage('equal'),
      value: '=',
    },
    {
      label: formatMessage('noEqual'),
      value: '!=',
    },
    // {
    //   label: formatMessage('alert.range'),
    //   value: 'range',
    // },
  ];

  return (
    <Form.Item label={formatMessage('condition')} required>
      <Space
        style={{ width: '100%' }}
        styles={{
          item: { width: '50%', overflow: 'hidden' },
        }}
        align="start"
      >
        <Form.Item
          name={['protocol', 'condition']}
          style={{ marginBottom: 0 }}
          rules={[{ required: true, message: commonFormatMessage('rule.required') }]}
        >
          <ComSelect allowClear options={options} style={{ width: '100%' }} placeholder={formatMessage('condition')} />
        </Form.Item>
        <Form.Item noStyle shouldUpdate={(pre, cur) => pre?.protocol?.condition !== cur?.protocol?.condition}>
          {({ getFieldValue }) => {
            const conditionType = getFieldValue(['protocol', 'condition']);
            return conditionType === 'range' ? (
              <Space
                style={{ width: '100%' }}
                styles={{
                  item: { width: '50%' },
                }}
              >
                <Form.Item name="num1" style={{ marginBottom: 0 }}>
                  <InputNumber style={{ width: '100%' }} placeholder={formatMessage('min')} />
                </Form.Item>
                <Form.Item name="num2" style={{ marginBottom: 0 }}>
                  <InputNumber style={{ width: '100%' }} placeholder={formatMessage('max')} />
                </Form.Item>
              </Space>
            ) : (
              <Form.Item
                name={['protocol', 'limitValue']}
                style={{ marginBottom: 0 }}
                rules={[{ required: true, message: commonFormatMessage('rule.required') }]}
              >
                <InputNumber style={{ width: '100%' }} placeholder={formatMessage('regularValue')} />
              </Form.Item>
            );
          }}
        </Form.Item>
      </Space>
    </Form.Item>
  );
};

export default Condition;
