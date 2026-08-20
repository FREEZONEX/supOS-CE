import { useTranslate } from '@/hooks';
import { Form, type FormItemProps, Input } from 'antd';
import type { ReactNode } from 'react';
import ComCheckbox from '@/components/com-checkbox';
import ComRadio from '../../../../components/com-radio';
import SourceSelect from '@/pages/menu-configuration/components/menu-content/SourceSelect.tsx';
import { Fragment } from 'react';
import CodeInput from '@/pages/menu-configuration/components/menu-content/CodeInput.tsx';
import MenuIconField from '@/pages/menu-configuration/components/menu-content/MenuIconField.tsx';
import type { MenuProps } from '@/pages/menu-configuration/store/types.ts';
import styles from './BasicInfo.module.scss';

export interface FormItemType {
  formType: string;
  label?: string;
  formProps: FormItemProps;
  childProps?: { [key: string]: any };
}

const { TextArea } = Input;

const FieldRow = ({ label, children, className }: { label?: string; children: ReactNode; className?: string }) => (
  <div className={`${styles.fieldRow}${className ? ` ${className}` : ''}`}>
    <div className={styles.fieldLabel}>{label}</div>
    <div className={styles.fieldValue}>{children}</div>
  </div>
);

const render = (item: FormItemType, iconMenuItem?: MenuProps | null, iconDisabled?: boolean) => {
  const { formType, label, formProps, childProps } = item;
  if (formProps.hidden) {
    switch (formType) {
      case 'checkbox':
        return (
          <Form.Item {...formProps} key={formProps.name}>
            <ComCheckbox {...childProps} />
          </Form.Item>
        );
      default:
        return (
          <Form.Item {...formProps} key={formProps.name}>
            <Input {...childProps} />
          </Form.Item>
        );
    }
  }

  switch (formType) {
    case 'checkbox':
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            <ComCheckbox {...childProps} />
          </Form.Item>
        </FieldRow>
      );
    case 'codeInput':
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            <CodeInput {...childProps} />
          </Form.Item>
        </FieldRow>
      );
    case 'radioGroup':
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            <ComRadio {...childProps} />
          </Form.Item>
        </FieldRow>
      );
    case 'custom':
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            {childProps?.children}
          </Form.Item>
        </FieldRow>
      );
    case 'sourceSelect':
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            <SourceSelect {...childProps} />
          </Form.Item>
        </FieldRow>
      );
    case 'menuIcon':
      return (
        <FieldRow key={formProps.name} label={label} className={styles.menuIconRow}>
          <Form.Item {...formProps} noStyle>
            <MenuIconField menuItem={iconMenuItem} disabled={iconDisabled} />
          </Form.Item>
        </FieldRow>
      );
    case 'textArea':
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            <TextArea {...childProps} />
          </Form.Item>
        </FieldRow>
      );
    default:
      return (
        <FieldRow key={formProps.name} label={label}>
          <Form.Item {...formProps} noStyle>
            <Input {...childProps} />
          </Form.Item>
        </FieldRow>
      );
  }
};

const BasicInfo = ({
  configs,
  iconMenuItem,
  iconDisabled,
}: {
  configs: FormItemType[];
  iconMenuItem?: MenuProps | null;
  iconDisabled?: boolean;
}) => {
  const formatMessage = useTranslate();
  return (
    <section className={styles.basicSection}>
      <h3 className={styles.sectionTitle}>{formatMessage('MenuConfiguration.basicInfo')}</h3>
      <div className={styles.fieldList}>
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
          ...configs,
        ]?.map((item) => {
          return <Fragment key={item.formProps.name}>{render(item, iconMenuItem, iconDisabled)}</Fragment>;
        })}
      </div>
    </section>
  );
};

export default BasicInfo;
