import type { FC } from 'react';
import { Form, Select } from 'antd';
import { useTranslate } from '@/hooks';
import FieldsFormList from '@/pages/uns/components/use-create-modal/components/FieldsFormList';

export interface ModelFieldsFormProps {
  types?: string[];
  options?: OptionType[];
}

interface OptionType {
  label: string;
  value: string;
}

const ModelFieldsForm: FC<ModelFieldsFormProps> = ({ types, options }) => {
  const formatMessage = useTranslate();
  const form = Form.useFormInstance();
  const modelId = Form.useWatch('modelId', form) || form.getFieldValue('modelId');

  return (
    <>
      <Form.Item name="modelId" label={formatMessage('common.template')} rules={[{ required: true }]}>
        <Select
          showSearch
          optionFilterProp="path"
          options={options}
          onChange={() => {
            form.setFieldValue('mainKey', undefined);
          }}
        />
      </Form.Item>
      {modelId && <FieldsFormList types={types} disabled showWrap={false} />}
    </>
  );
};
export default ModelFieldsForm;
