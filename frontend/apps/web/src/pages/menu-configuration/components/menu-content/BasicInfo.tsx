import { useTranslate } from '@/hooks';
import { Divider, Form, type FormItemProps, Input } from 'antd';
import ComCheckbox from '@/components/com-checkbox';
import ComRadio from '../../../../components/com-radio';
import SourceSelect from '@/pages/menu-configuration/components/menu-content/SourceSelect.tsx';
import { Fragment } from 'react';
import UploadPicture from '@/pages/menu-configuration/components/menu-content/UploadPicture.tsx';
import CodeInput from '@/pages/menu-configuration/components/menu-content/CodeInput.tsx';

export interface FormItemType {
  formType: string;
  formProps: FormItemProps;
  childProps?: { [key: string]: any };
}

const { TextArea } = Input;

const render = (item: FormItemType) => {
  const { formType, formProps, childProps } = item;
  switch (formType) {
    case 'checkbox':
      return (
        <Form.Item {...formProps}>
          <ComCheckbox {...childProps} />
        </Form.Item>
      );
    case 'codeInput':
      return (
        <Form.Item {...formProps}>
          <CodeInput {...childProps} />
        </Form.Item>
      );
    case 'radioGroup':
      return (
        <Form.Item {...formProps}>
          <ComRadio {...childProps} />
        </Form.Item>
      );
    case 'sourceSelect':
      return (
        <Form.Item {...formProps}>
          <SourceSelect {...childProps} />
        </Form.Item>
      );
    case 'uploadPicture':
      return (
        <Form.Item {...formProps}>
          <UploadPicture {...childProps} />
        </Form.Item>
      );
    case 'textArea':
      return (
        <Form.Item {...formProps}>
          <TextArea {...childProps} />
        </Form.Item>
      );
    default:
      return (
        <Form.Item {...formProps}>
          <Input {...childProps} />
        </Form.Item>
      );
  }
};

const BasicInfo = ({ configs }: { configs: FormItemType[] }) => {
  const formatMessage = useTranslate();
  return (
    <>
      <div style={{ fontSize: 20, fontWeight: 500 }}>{formatMessage('MenuConfiguration.basicInfo')}</div>
      <Divider style={{ backgroundColor: '#C6C6C6', margin: '16px 0' }} />
      {[
        {
          formType: 'input',
          formProps: {
            name: 'sort',
            hidden: true,
          },
        },
        {
          formType: 'input',
          formProps: {
            name: 'type',
            hidden: true,
          },
        },
        {
          formType: 'input',
          formProps: {
            name: 'id',
            hidden: true,
          },
        },
        {
          formType: 'input',
          formProps: {
            name: 'parentId',
            hidden: true,
          },
        },
        {
          formType: 'input',
          formProps: {
            name: 'icon',
            hidden: true,
          },
        },
        ...configs,
      ]?.map((item) => {
        return <Fragment key={item.formProps.name}>{render(item)}</Fragment>;
      })}
    </>
  );
};

export default BasicInfo;
