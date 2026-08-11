import { Form, type FormInstance, type FormProps } from 'antd';
import type { FC } from 'react';
import { v4 as uuidv4 } from 'uuid';
import ProSearch from '../pro-search';
import RenderFormItem, { type RenderFormItemProps } from '../operation-form/render-form-item.tsx';
import { getIntl } from '@/stores/i18n-store.ts';
import './index.scss';

export interface ComSearchProps {
  form: FormInstance;
  formConfig?: FormProps;
  onSearch?: () => void;
  formItemOptions: RenderFormItemProps[];
  loading?: boolean;
}

const isSearchInputField = (item: RenderFormItemProps) =>
  Boolean(item.name && !item.hidden && !item.render && !item.component && (item.type === 'Input' || !item.type));

const ComSearch: FC<ComSearchProps> = ({ form, formConfig, formItemOptions, onSearch }) => {
  const { onFinish: formOnFinish, disabled, ...restFormConfig } = formConfig || {};

  const triggerSearch = () => {
    form.submit();
  };

  const handleFinish = (values: Record<string, unknown>) => {
    if (formOnFinish) {
      formOnFinish(values);
      return;
    }
    onSearch?.();
  };

  return (
    <Form
      className="com-search"
      labelAlign={'left'}
      colon={false}
      form={form}
      layout="inline"
      disabled={disabled}
      {...restFormConfig}
      style={{ flexWrap: 'nowrap', ...restFormConfig?.style }}
      onFinish={handleFinish}
    >
      {formItemOptions?.map((item: RenderFormItemProps) => {
        const key = String(item.name || uuidv4());

        if (item.hidden) {
          return <RenderFormItem key={key} {...item} />;
        }

        if (isSearchInputField(item)) {
          const width = item.properties?.style?.width ?? 300;
          return (
            <Form.Item key={key} name={item.name} style={{ marginInlineEnd: 0, ...item.style }}>
              <ProSearch
                size="sm"
                className="custom-search-page"
                style={{ width }}
                placeholder={item.properties?.placeholder}
                closeButtonLabelText={getIntl('common.clearSearchInput')}
                disabled={disabled || item.properties?.disabled}
                onSearch={triggerSearch}
              />
            </Form.Item>
          );
        }

        return <RenderFormItem key={key} {...item} />;
      })}
    </Form>
  );
};

export default ComSearch;
