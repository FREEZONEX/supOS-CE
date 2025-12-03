import { useTranslate } from '@/hooks';
import styles from './index.module.scss';
import { useState } from 'react';
import defaultIconUrl from '@/assets/home-icons/default.svg';
import { Divider, Flex, Form, Image, Input, Typography } from 'antd';
import { Launch } from '@carbon/icons-react';
import { Original_TARGET_PATH, STORAGE_PATH } from '@/common-types/constans.ts';
import { AuthButton } from '@/components/auth';
import ComCopy from '@/components/com-copy';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ProModal from '@/components/pro-modal';
import { useBaseStore } from '@/stores/base';
import ProCardContainer from '@/components/pro-card/ProCardContainer.tsx';
import { ProCard } from '@/components';

const { Title, Paragraph } = Typography;

const Index = () => {
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);
  const [openInfo, setOpenInfo] = useState<any>();
  const formatMessage = useTranslate();
  const containerList = useBaseStore((state) => state.containerList);

  return (
    <ComLayout>
      <ComContent title={<div></div>} hasBack={false} mustShowTitle={false}>
        <div className={styles['home-title']} style={{ borderBottom: '1px solid var(--supos-home-border-color)' }}>
          <Title style={{ fontWeight: 400, marginBottom: 5 }} type="secondary" level={2}>
            {formatMessage('advancedUse.advancedUse')}
          </Title>
          <Paragraph style={{ marginBottom: 0 }}>{formatMessage('advancedUse.overview')}</Paragraph>
        </div>
        <div className={styles['content-section']}>
          <ProCardContainer>
            {(containerList?.advancedUse || []).map((item, index) => {
              return (
                <ProCard
                  key={index}
                  header={{
                    iconSrc: `${STORAGE_PATH}${Original_TARGET_PATH}/${item?.envMap?.service_logo}`,
                    title: item?.name,
                    titleDescription: `${formatMessage('common.version')} ${item.version || ''}`,
                  }}
                  item={item}
                  onClick={(item) => {
                    if (item.envMap.service_password) {
                      const _item = {
                        ...item,
                        account: item.envMap.service_account,
                        password: item.envMap.service_password,
                      };
                      form.setFieldsValue(_item);
                      setOpenInfo(_item);
                      setOpen(true);
                    } else {
                      window.open(`http://${window.location.host}${item.envMap.service_redirect_url}`);
                    }
                  }}
                  description={item?.description}
                />
              );
            })}
          </ProCardContainer>
        </div>
      </ComContent>
      <ProModal
        title={
          <Flex align="center" style={{ width: '100%' }} gap={8}>
            <Image
              preview={false}
              src={`${STORAGE_PATH}${Original_TARGET_PATH}/${openInfo?.envMap?.service_logo}`}
              height={20}
              fallback={defaultIconUrl}
            />
            {openInfo?.name}
          </Flex>
        }
        open={open}
        size="xxs"
        onCancel={() => {
          setOpen(false);
          setOpenInfo(null);
          form.resetFields();
        }}
      >
        <Flex vertical align="center" justify="center" gap={8}>
          <Divider style={{ margin: 0 }} />
          <Form form={form} colon={false}>
            <Form.Item label={formatMessage('account.account')} name="account" style={{ marginBottom: 8 }}>
              <Input
                className="no-foucs"
                variant="borderless"
                readOnly
                addonAfter={<ComCopy textToCopy={openInfo?.account} />}
              />
            </Form.Item>
            <Form.Item label={formatMessage('appGui.password')} name="password" style={{ marginBottom: 8 }}>
              <Input.Password variant="borderless" readOnly addonAfter={<ComCopy textToCopy={openInfo?.password} />} />
            </Form.Item>
          </Form>
          <div>
            <AuthButton
              type="primary"
              onClick={() => {
                window.open(`http://${window.location.host}${openInfo.envMap.service_redirect_url}`);
              }}
            >
              {`${formatMessage('common.open')} ${openInfo?.name || ''}`} <Launch />
            </AuthButton>
          </div>
        </Flex>
      </ProModal>
    </ComLayout>
  );
};

export default Index;
