import { useEffect, useState } from 'react';
import { Button, Flex, Form, Typography } from 'antd';
import { useTranslate } from '@/hooks';
import ProModal from '@/components/pro-modal';
import styles from './index.module.scss';
import FrequencyForm from '../use-create-modal/components/file/FrequencyForm';
import { BookmarkAdd, BookmarkFilled } from '@carbon/icons-react';

interface PropsTypes {
  value: boolean; // 是否选中
  showModal?: boolean; // 是否显示Modal
  topic?: string; // topic名称
  fileCount?: number; // 文件数量
  speed?: string; // 频率
  unit?: string; // 单位
  onChange: (checked: boolean, config?: any) => Promise<void>;
  subscribeFrequency?: string;
  hidden?: boolean;
}

const Subscribe = (props: PropsTypes) => {
  const { onChange, value, showModal, topic, fileCount, subscribeFrequency, hidden } = props;
  const [form] = Form.useForm();
  const [open, setOpen] = useState<boolean>(false);
  const [checked, setChecked] = useState<boolean>(false);
  const formatMessage = useTranslate();

  useEffect(() => {
    setChecked(value);
  }, [value]);

  const handleChange = () => {
    setChecked((state) => !state);
    if (!showModal || !checked === false) {
      onChange?.(!checked);
      return;
    }

    form.resetFields();

    setOpen(true);
  };

  const handleCancel = () => {
    setOpen(false);
    setChecked(value || false);
  };

  const handleOk = async () => {
    const { frequency } = await form.validateFields();
    onChange(true, `${frequency.value}${frequency.unit}`).then(() => {
      setOpen(false);
    });
  };

  if (hidden) return null;
  return (
    <Flex className={styles.subscribe}>
      <Button
        color="default"
        variant="filled"
        onClick={handleChange}
        style={{ background: '#F4F4F4', color: '#585C62' }}
      >
        {!checked ? <BookmarkAdd /> : <BookmarkFilled style={{ color: '#F1C21B' }} />}
        {checked ? formatMessage('common.subscribed') : formatMessage('common.subscribe')}
        {subscribeFrequency ? `: ${subscribeFrequency}` : ''}
      </Button>
      {/*<span className={styles.label}>{formatMessage('common.subscribe')}</span>*/}
      {/*<Switch checked={checked} onChange={handleChange} />*/}
      <ProModal width={500} destroyOnHidden title={formatMessage('common.config')} open={open} onCancel={handleCancel}>
        <div>
          <div className={styles.configBox}>
            <span className={styles.configLabel}>Topic</span>
            <Typography.Text title={topic} style={{ flex: 1 }} ellipsis>
              {topic}
            </Typography.Text>
            {/*<span className={styles.configValue}>{topic}</span>*/}
          </div>
          {fileCount ? (
            <div className={styles.configBox}>
              <span className={styles.configLabel}>{formatMessage('uns.fileCount')}</span>
              <span className={`${styles.configValue} ${styles.fileCount}`}>{fileCount}</span>
            </div>
          ) : null}
          <div className={styles.configBox}>
            <span className={styles.configLabel}>{formatMessage('uns.subscriptionFrequency')}</span>
            <span className={styles.configValue} style={{ flex: 1 }}>
              <Form form={form} style={{ width: '100%' }}>
                <Form.Item style={{ marginBottom: 0 }}>
                  <FrequencyForm unitList={['s', 'm', 'h', 'd']} maxValue={99} />
                </Form.Item>
              </Form>
            </span>
          </div>
        </div>
        <Flex justify="flex-end">
          <Button onClick={handleOk} style={{ marginTop: 8 }} type="primary">
            {formatMessage('common.save')}
          </Button>
        </Flex>
      </ProModal>
    </Flex>
  );
};

export default Subscribe;
