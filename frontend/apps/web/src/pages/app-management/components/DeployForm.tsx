import { Button, Divider, Flex, Form, Input, message, Select, Typography } from 'antd';
import { type FC, useEffect, useRef, useState } from 'react';
import { getApps, uploadHtml } from '@/apis/inter-api/apps';
import { useNavigate } from 'react-router';
import styles from './DeployForm.module.scss';
import { useTranslate } from '@/hooks';
import { getBaseFileName } from '@/utils/url-util';
const { Title } = Typography;

const DeployForm: FC<any> = ({ htmlName, appName, show, setShow, htmlContent, getHtmlContent }) => {
  const [success, setSuccess] = useState(false);
  const [options, setOptions] = useState([]);
  const [form] = Form.useForm();
  const [name, setName] = useState('');
  const navigate = useNavigate();
  const originUrl = useRef<string | null>(null);
  const formatMessage = useTranslate();

  useEffect(() => {
    if (show) {
      getApps().then((data: any) => {
        setOptions(data?.map((d: any) => d.name) || []);
        if (htmlName || appName) {
          form.setFieldsValue({
            appName,
            htmlName: getBaseFileName(htmlName),
          });
        }
      });
    }
  }, [show]);
  const onCancel = () => {
    setShow(false);
  };
  const onSave = async () => {
    if (!htmlContent && !getHtmlContent) return;
    const content = htmlContent || getHtmlContent?.();
    const values = await form.validateFields();
    const blob = new Blob([content], { type: 'text/html' });
    uploadHtml(values.appName, {
      value: blob,
      name: 'file',
      fileName: `${values.htmlName}.html`,
    }).then(() => {
      // 根据生成的html名字，把url匹配出来
      originUrl.current = `/inter-api/supos/app/${values.appName}/preview/${values.htmlName}.html`;
      setName(values.appName);
      setSuccess(true);
      form.setFieldsValue({
        appName: appName,
        htmlName: undefined,
      });
    });
  };
  return success ? (
    <Flex style={{ height: 'calc(100% - 70px)', overflow: 'hidden' }} justify="center" vertical align="center" gap={20}>
      <Title style={{ marginBottom: 0 }} type="secondary" level={5}>
        {formatMessage('appGui.saveSuccess')}
      </Title>
      <div>
        <Button
          type="link"
          style={{ textDecoration: 'underline' }}
          onClick={() => {
            navigate('/app-space', {
              state: { name },
            });
            setShow(false);
            setSuccess(false);
          }}
        >
          {formatMessage('appGui.goAppDisplay')}
        </Button>
        <Button
          type="link"
          style={{ textDecoration: 'underline' }}
          onClick={() => {
            if (originUrl?.current) {
              window.open(originUrl.current);
              setSuccess(false);
            } else {
              message.error(formatMessage('common.UrlLose'));
            }
          }}
        >
          {formatMessage('appGui.goHTML')}
        </Button>
        <Button
          type="link"
          style={{
            textDecoration: 'underline',
            backgroundColor: 'var(--supos-button-def-10)',
            color: 'var(--supos-text-color)',
          }}
          onClick={() => {
            setSuccess(false);
          }}
        >
          {formatMessage('common.back')}
        </Button>
      </div>
    </Flex>
  ) : (
    <Form
      style={{ padding: '20px 40px', overflow: 'hidden' }}
      colon={false}
      labelCol={{ span: 8 }}
      wrapperCol={{ span: 16 }}
      form={form}
      className={styles['deploy-form']}
    >
      <Form.Item label={formatMessage('appGui.deploy')}></Form.Item>
      <Form.Item label={formatMessage('appGui.targetApp')} name="appName" required>
        <Select placeholder="select app" disabled={!!appName}>
          {options?.map((item) => {
            return <Select.Option value={item}>{item}</Select.Option>;
          })}
        </Select>
      </Form.Item>
      <Form.Item label={formatMessage('appGui.targetHtml')} name="htmlName" required>
        <Input placeholder="type html name" disabled={!!htmlName} />
      </Form.Item>
      <Divider></Divider>
      <Flex justify="flex-end" gap="10px">
        <Button
          style={{
            width: 76,
            height: 30,
            backgroundColor: 'var(--supos-button-def-10)',
            color: 'var(--supos-text-color)',
          }}
          color="default"
          variant="filled"
          onClick={onCancel}
        >
          {formatMessage('common.cancel')}
        </Button>
        <Button style={{ width: 76, height: 30 }} type="primary" variant="solid" onClick={onSave}>
          {formatMessage('common.save')}
        </Button>
      </Flex>
    </Form>
  );
};

export default DeployForm;
