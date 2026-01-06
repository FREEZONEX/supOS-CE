import { useEffect, useState } from 'react';
import styles from './index.module.scss';
import { Button, Form } from 'antd';
import FrequencyForm from '@/pages/uns/components/use-create-modal/components/file/FrequencyForm';
import Icon, { EnterOutlined } from '@ant-design/icons';
import FileEdit from '../svg-components/FileEdit';

interface PropsTypes {
  className?: string;
  value?: string;
  onChange: (val: any) => void;
}

const CustomParagraph = (props: PropsTypes) => {
  const { className, value, onChange } = props;
  const [form] = Form.useForm();

  const [active, setActive] = useState<boolean>(false);

  useEffect(() => {
    if (!value) return;
    const match = value.match(/^(\d+)([smhd])$/);
    if (match) {
      form.setFieldValue('frequency', { value: match[1], unit: match[2] });
    }
  }, [value]);

  const handleChangeActive = () => {
    setActive(!active);
  };

  const handleSave = async () => {
    const { frequency } = await form.validateFields();
    onChange({ frequency: `${frequency.value}${frequency.unit}`, enable: true });
    setActive(false);
  };

  return (
    <div className={`${styles.container} ${className || ''}`}>
      {!active ? (
        <>
          <div className={styles.label}>{value}</div>
          <div className={styles.editIcon} onClick={handleChangeActive}>
            <Icon
              component={FileEdit}
              style={{
                fontSize: 17,
                color: 'var(--supos-text-color)',
              }}
            />
          </div>
        </>
      ) : (
        <div className={styles.editComponent}>
          <Form form={form}>
            <Form.Item className={styles.frequencyForm}>
              <FrequencyForm unitList={['s', 'm', 'h', 'd']} maxValue={99} />
              <Button icon={<EnterOutlined />} onClick={handleSave}></Button>
            </Form.Item>
          </Form>
        </div>
      )}
    </div>
  );
};

export default CustomParagraph;
