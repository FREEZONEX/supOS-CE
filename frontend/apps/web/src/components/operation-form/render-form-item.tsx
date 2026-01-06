import { Input, Form, InputNumber, type FormItemProps, Upload, Button, Switch } from 'antd';
import ComSelect from '../com-select';
import { isValidElement, type ReactNode } from 'react';
import TagSelect from '@/pages/uns/components/use-create-modal/components/TagSelect.tsx';
import ComCheckbox from '../com-checkbox';
import { getIntl } from '@/stores/i18n-store.ts';

export const render = (type: any, properties: any) => {
  switch (type) {
    case 'Input':
      return <Input {...properties} />;
    case 'Password':
      return <Input.Password {...properties} />;
    case 'TextArea':
      return <Input.TextArea rows={2} {...properties} />;
    case 'Select':
      return <ComSelect allowClear {...properties} />;
    case 'Number':
      return <InputNumber style={{ width: '100%' }} allowClear {...properties} />;
    case 'TagSelect':
      return <TagSelect {...properties} />;
    case 'Checkbox':
      return <ComCheckbox {...properties} />;
    case 'Upload':
      return (
        <Upload {...properties}>
          <Button>
            <Upload />
            {getIntl('common.upload')}
          </Button>
        </Upload>
      );
    case 'Switch':
      return <Switch {...properties} />;
    default:
      return null;
  }
};

export interface RenderFormItemProps extends FormItemProps {
  // 完全自定义渲染
  render?: ((item?: RenderFormItemProps) => ReactNode) | ReactNode;
  // 组件类型，通过properties设置参数
  type?: string;
  // 组件参数
  properties?: { [key: string]: any };
  // 自定义form下的组件,properties、type、render无效
  component?: ReactNode;
  style?: { [key: string]: any };
}

const RenderFormItem = (item: RenderFormItemProps) => {
  const { render: renderComponent, type = 'Input', properties, component, ...restFormProps } = item;
  // 组件渲染
  if (renderComponent) {
    if (isValidElement(renderComponent)) {
      return renderComponent;
    }
    if (typeof renderComponent === 'function') {
      return renderComponent(item);
    }
  }
  return (
    <Form.Item {...restFormProps} valuePropName={['Switch', 'Checkbox'].includes(type) ? 'checked' : 'value'}>
      {component ? component : !restFormProps?.name ? null : render(type, properties)}
    </Form.Item>
  );
};

export default RenderFormItem;
