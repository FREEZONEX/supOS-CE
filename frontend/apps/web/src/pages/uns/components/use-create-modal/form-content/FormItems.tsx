import type { FC } from 'react';
import { Form, Input, Select, Divider } from 'antd';
import type { FormItemProps } from 'antd';
import TagSelect from '@/pages/uns/components/use-create-modal/components/TagSelect';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect';
import FieldsFormList from '@/pages/uns/components/use-create-modal/components/FieldsFormList';
import ModelFieldsForm from '@/pages/uns/components/use-create-modal/components/file/ModelFieldsForm';
import ReverseGeneration from '@/pages/uns/components/use-create-modal/components/file/ReverseGeneration';
import TopicToUnsFieldsList from '@/pages/uns/components/use-create-modal/components/file/TopicToUnsFieldsList';
import FrequencyForm from '@/pages/uns/components/use-create-modal/components/file/FrequencyForm';
import ExpressionForm from '@/pages/uns/components/use-create-modal/components/file/timeSeries/ExpressionForm';
import AggForm from '@/pages/uns/components/use-create-modal/components/file/timeSeries/AggForm';
import AdvancedOptions from '@/pages/uns/components/use-create-modal/components/file/AdvancedOptions';
import ExpandedKeyFormList from '@/pages/uns/components/ExpandedKeyFormList';
import ComCheckbox from '@/components/com-checkbox';
import ComRadio from '@/components/com-radio';
import AttributeTypeForm from '@/pages/uns/components/use-create-modal/components/file/AttributeTypeForm';

const { TextArea } = Input;

export interface FormItemType {
  formType: string;
  formProps: FormItemProps;
  childProps?: { [key: string]: any };
}

export interface FormItemsProps {
  formData: FormItemType[];
}
const FormItems: FC<FormItemsProps> = ({ formData }) => {
  return (
    <>
      {formData.map((item: FormItemType) => {
        const { formType, formProps = {}, childProps = {} } = item;
        const key = formProps.name;
        switch (formType) {
          case 'showTopic':
            return (
              <Form.Item {...formProps} key={key}>
                <div className="namespaceValue">{formProps.initialValue}</div>
              </Form.Item>
            );
          case 'divider':
            return <Divider style={{ borderColor: '#c6c6c6' }} key={key} />;
          case 'input':
            return (
              <Form.Item {...formProps} key={key}>
                <Input {...childProps} />
              </Form.Item>
            );
          case 'textArea':
            return (
              <Form.Item {...formProps} key={key}>
                <TextArea {...childProps} />
              </Form.Item>
            );
          case 'select':
            return (
              <Form.Item {...formProps} key={key}>
                <Select {...childProps} />
              </Form.Item>
            );
          case 'tagSelect':
            return (
              <Form.Item {...formProps} key={key}>
                <TagSelect {...childProps} />
              </Form.Item>
            );
          case 'searchSelect':
            return (
              <Form.Item {...formProps} key={key}>
                <SearchSelect {...childProps} />
              </Form.Item>
            );
          case 'radioGroup':
            return (
              <Form.Item {...formProps} key={key}>
                <ComRadio {...childProps} />
              </Form.Item>
            );
          case 'checkbox':
            return (
              <Form.Item {...formProps} key={key}>
                <ComCheckbox {...childProps} />
              </Form.Item>
            );
          case 'frequency':
            return (
              <Form.Item {...formProps} key={key}>
                <FrequencyForm {...childProps} />
              </Form.Item>
            );
          case 'expandFormList':
            return <ExpandedKeyFormList key={key} />;
          case 'fieldsFormList':
            return <FieldsFormList {...childProps} key={key} />;
          case 'modelFieldsForm':
            return <ModelFieldsForm {...childProps} key={key} />;
          case 'reverseGeneration':
            return <ReverseGeneration {...childProps} key={key} />;
          case 'topicToUnsFieldsList':
            return <TopicToUnsFieldsList {...childProps} key={key} />;
          case 'expressionForm':
            return <ExpressionForm key={key} {...childProps} />;
          case 'aggForm':
            return <AggForm key={key} />;
          case 'advancedOptions':
            return <AdvancedOptions key={key} />;
          case 'attributeTypeForm':
            return <AttributeTypeForm key={key} {...childProps} />;
          default:
            return null;
        }
      })}
    </>
  );
};

export default FormItems;
