import type { FC } from 'react';
import { Form } from 'antd';
import { useTranslate } from '@/hooks';
import ComRadio from '@/components/com-radio';

interface BooleanOptionProps {
  name: any;
  label: string;
  isRequired?: boolean;
  onChange?: () => void;
}

const BooleanOption: FC<BooleanOptionProps> = ({ name, label, isRequired = false, onChange }) => {
  const formatMessage = useTranslate();

  return (
    <Form.Item
      label={label}
      name={name}
      rules={isRequired ? [{ required: true, message: formatMessage('rule.required') }] : []}
    >
      <ComRadio
        options={[
          { label: formatMessage('streams.yes'), value: true },
          { label: formatMessage('streams.no'), value: false },
        ]}
        onChange={onChange}
      />
    </Form.Item>
  );
};

export default BooleanOption;
