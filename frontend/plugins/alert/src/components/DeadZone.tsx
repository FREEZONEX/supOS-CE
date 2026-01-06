import { Form, InputNumber, Space } from 'antd';
import { useTranslate } from '@supos_host/hooks';
import { ComSelect } from '@supos_host/components';
import { REMOTE_NAME } from '../../variables';

const DeadZone = () => {
  const formatMessage = useTranslate(REMOTE_NAME);

  const options = [
    {
      label: formatMessage('value'),
      value: 1,
    },
    {
      label: formatMessage('percent'),
      value: 2,
    },
  ];

  return (
    <Form.Item label={formatMessage('deadZone')}>
      <Space
        style={{ width: '100%' }}
        styles={{
          item: { width: '50%', overflow: 'hidden' },
        }}
        align="start"
      >
        <Form.Item name={['protocol', 'deadBandType']} style={{ marginBottom: 0 }}>
          <ComSelect options={options} style={{ width: '100%' }} placeholder={formatMessage('deadZone')} />
        </Form.Item>
        <Form.Item name={['protocol', 'deadBand']} style={{ marginBottom: 0 }}>
          <InputNumber style={{ width: '100%' }} placeholder={formatMessage('regularValue')} />
        </Form.Item>
      </Space>
    </Form.Item>
  );
};

export default DeadZone;
