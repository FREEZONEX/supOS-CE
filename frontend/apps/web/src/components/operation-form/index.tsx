import type { CSSProperties, FC, ReactNode } from 'react';
import { Button, Divider, Flex, Form, type FormInstance, type FormProps, Typography } from 'antd';
import RenderFormItem, { type RenderFormItemProps } from '../operation-form/render-form-item';
import { useTranslate } from '@/hooks';
import { v4 as uuidv4 } from 'uuid';
import './index.scss';

const { Title } = Typography;

export interface OperationFormProps {
  form: FormInstance;
  formConfig?: FormProps;
  onSave?: () => void;
  onCancel?: () => void;
  formItemOptions: RenderFormItemProps[];
  title?: ReactNode;
  loading?: boolean;
  style?: CSSProperties;
  footer?: ReactNode;
  buttonConfig?: {
    block?: boolean;
  };
}

const OperationForm: FC<OperationFormProps> = ({
  form,
  formConfig,
  onSave,
  onCancel,
  formItemOptions,
  title,
  loading,
  style,
  footer,
  buttonConfig,
}) => {
  const formatMessage = useTranslate();
  const buttonStyle = buttonConfig?.block
    ? {
        height: 40,
      }
    : {
        height: 30,
      };
  return (
    <Form
      labelAlign={'left'}
      className={'operation-form'}
      style={{ padding: '20px 40px', overflow: 'hidden', ...style }}
      colon={false}
      labelCol={{ span: 11 }}
      wrapperCol={{ span: 13 }}
      {...formConfig}
      form={form}
      labelWrap
    >
      {title && (
        <Typography style={{ marginBottom: 40 }}>
          <Title level={4}>{title}</Title>
        </Typography>
      )}
      {formItemOptions?.map((item: any) => {
        if (item.type === 'divider') {
          return <Divider key={uuidv4()} style={{ background: '#c6c6c6' }}></Divider>;
        }
        return <RenderFormItem key={item.name || uuidv4()} {...item} />;
      })}
      {footer ? (
        footer
      ) : (
        <Flex gap="10px" justify="end">
          <Button
            style={{
              ...buttonStyle,
              backgroundColor: 'var(--supos-uns-button-color)',
              color: 'var(--supos-text-color)',
              flex: 1,
              minWidth: 0,
            }}
            color="default"
            variant="filled"
            onClick={onCancel}
            block={buttonConfig?.block}
            title={formatMessage('common.cancel')}
          >
            {formatMessage('common.cancel')}
          </Button>
          <Button
            style={{
              ...buttonStyle,
              flex: 1,
              minWidth: 0,
            }}
            type="primary"
            variant="solid"
            onClick={onSave}
            loading={loading}
            block={buttonConfig?.block}
            title={formatMessage('common.save')}
          >
            {formatMessage('common.save')}
          </Button>
        </Flex>
      )}
    </Form>
  );
};

export default OperationForm;
