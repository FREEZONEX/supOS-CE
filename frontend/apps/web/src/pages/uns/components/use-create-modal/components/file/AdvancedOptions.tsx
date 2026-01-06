import { Form, Input, Select, DatePicker } from 'antd';
import BooleanOption from '../BooleanOption';
import { useTranslate } from '@/hooks';

const AdvancedOptions = () => {
  const formatMessage = useTranslate();
  const form = Form.useFormInstance();
  const trigger = Form.useWatch(['_advancedOptions', 'trigger']);
  const isFillHistory = Form.useWatch(['_advancedOptions', 'fillHistory']) || false;
  const windowType = form.getFieldValue(['streamOptions', 'window', 'windowType']);

  return (
    <div className="formBox">
      <Form.Item label={formatMessage('streams.trigger')} name={['_advancedOptions', 'trigger']}>
        <Select
          onChange={() => {
            form.setFieldValue(['_advancedOptions', 'delayTime'], undefined);
          }}
        >
          <Select.Option value="AT_ONCE">AT_ONCE</Select.Option>
          <Select.Option value="WINDOW_CLOSE">WINDOW_CLOSE</Select.Option>
          <Select.Option value="MAX_DELAY">MAX_DELAY</Select.Option>
        </Select>
      </Form.Item>
      {trigger === 'MAX_DELAY' && (
        <Form.Item label={formatMessage('streams.delayTime')} name={['_advancedOptions', 'delayTime']}>
          <Input placeholder={`${formatMessage('uns.eg')}: 120s`} />
        </Form.Item>
      )}
      <Form.Item
        label={formatMessage('streams.watermark')}
        name={['_advancedOptions', 'waterMark']}
        rules={[{ required: windowType === 'COUNT_WINDOW' }]}
      >
        <Input placeholder={`${formatMessage('uns.eg')}: 100s`} />
      </Form.Item>
      <Form.Item label={formatMessage('streams.deleteMark')} name={['_advancedOptions', 'deleteMark']}>
        <Input placeholder={`${formatMessage('uns.eg')}: 5m`} />
      </Form.Item>

      <BooleanOption
        name={['_advancedOptions', 'fillHistory']}
        label={formatMessage('streams.fillHistory')}
        onChange={() => {
          form.setFieldValue(['_advancedOptions', 'startTime'], undefined);
          form.setFieldValue(['_advancedOptions', 'endTime'], undefined);
        }}
      />
      <BooleanOption name={['_advancedOptions', 'ignoreUpdate']} label={formatMessage('streams.ignoreUpdate')} />
      <BooleanOption name={['_advancedOptions', 'ignoreExpired']} label={formatMessage('streams.ignoreExpired')} />
      {isFillHistory && (
        <>
          <Form.Item label={formatMessage('streams.startTime')} name={['_advancedOptions', 'startTime']}>
            <DatePicker />
          </Form.Item>
          <Form.Item label={formatMessage('streams.endTime')} name={['_advancedOptions', 'endTime']}>
            <DatePicker />
          </Form.Item>
        </>
      )}
    </div>
  );
};

export default AdvancedOptions;
