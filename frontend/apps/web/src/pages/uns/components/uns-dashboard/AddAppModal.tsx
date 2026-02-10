import { forwardRef, useImperativeHandle, useState } from 'react';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import { App, Col, Divider, Flex, Form, Input, Row, Segmented, type UploadFile } from 'antd';
import ComButton from '@/components/com-button';
import ComUploadPicture from '@/components/com-upload-picture';
import ComDraggerUpload from '@/components/com-dragger-upload';
import usePropsValue from '@/hooks/usePropsValue.ts';
import { ContainerSoftware, FolderAdd, Wikis } from '@carbon/icons-react';
import styles from './index.module.scss';
import CodeMirror from '@uiw/react-codemirror';
import { codemirrorTheme } from '@/theme/codemirror-theme.tsx';
import codeStyles from '@/theme/codemirror.module.scss';
import { yaml } from '@codemirror/lang-yaml';
import ComEllipsis from '@/components/com-ellipsis';
import ComCheckbox from '@/components/com-checkbox';
import { installApp } from '@/apis/inter-api/app.ts';
import { fetchBaseStore } from '@/stores/base';
import jsYaml from 'js-yaml';
import { type Diagnostic, linter, lintGutter } from '@codemirror/lint';

export interface AddAppModalRef {
  onOpen: (type: number, props?: any) => void;
  onClose: (e?: any) => void;
}

export interface AddAppModalProps {
  [key: string]: any;
}

const placeholder = `services:
  app-service:
    image: username/repository:tag
    container_name: my-web-app
    ports:
      - "4000:4000"
    environment:
      - TZ=UTC
    command: serve -s /app/web-dist -l 3000
    volumes:
      - /etc/docker/certs:/certs
    networks:
      - edge_network
    restart: always`;

const AddAppModal = forwardRef<AddAppModalRef, AddAppModalProps>((_, ref) => {
  const [visible, setVisible] = useState(false);
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message } = App.useApp();

  const onOpen = () => {
    setVisible(true);
  };

  const onClose = () => {
    setVisible(false);
    form.resetFields();
  };

  const onSave = async () => {
    const value = await form.validateFields();
    if (value?.imageConfig?.image?.[0]?.status === 'uploading') {
      return message.warning(formatMessage('uns.uploading'));
    }
    if (value?.imageConfig?.image?.[0]?.status === 'error') {
      return message.warning(formatMessage('uns.uploadError'));
    }
    return installApp({
      ...value,
      imagePath: value?.imageConfig?.type === 'image' ? value?.imageConfig?.imagePath : undefined,
      imageUrl: value?.imageConfig?.type === 'registry' ? value?.imageConfig?.imageUrl : undefined,
      iconFile: undefined,
      imageConfig: undefined,
    }).then(() => {
      message.success(formatMessage('common.optsuccess'));
      onClose?.();
      fetchBaseStore?.();
    });
  };

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  return (
    <ProModal
      // open
      open={visible}
      onCancel={onClose}
      title={formatMessage('uns.deployNewApp')}
      width={600}
      styles={{
        body: {
          paddingBlockStart: 0,
        },
      }}
      footer={
        <>
          <Divider style={{ background: '#e0e0e0', margin: '16px 0' }} />
          <Flex gap="10px" justify="end">
            <ComButton
              style={{
                backgroundColor: 'var(--supos-uns-button-color)',
                color: 'var(--supos-text-color)',
              }}
              color="default"
              variant="filled"
              onClick={onClose}
              title={formatMessage('common.cancel')}
            >
              {formatMessage('common.cancel')}
            </ComButton>
            <ComButton type="primary" variant="solid" onClick={onSave} title={formatMessage('uns.deployNow')}>
              {formatMessage('uns.deployNow')}
            </ComButton>
          </Flex>
        </>
      }
    >
      {(isFullscreen) => {
        return (
          <Form
            form={form}
            colon={false}
            layout="vertical"
            style={{ height: isFullscreen ? 'inherit' : 650, overflow: 'auto', padding: 6 }}
          >
            <Form.Item name="name" label={formatMessage('common.name')} rules={[{ required: true }]}>
              <Input
                count={{
                  show: true,
                  max: 10,
                }}
                placeholder={formatMessage('rule.pleaseInput', { label: formatMessage('common.name') })}
              />
            </Form.Item>
            <Form.Item name="description" label={formatMessage('common.description')}>
              <Input placeholder={formatMessage('rule.pleaseInput', { label: formatMessage('common.description') })} />
            </Form.Item>
            <Form.Item name="iconFile" label={formatMessage('uns.appIcon')}>
              <ComUploadPicture
                className={styles['icon-upload']}
                action="/inter-api/supos/attachment/upload"
                maxCount={1}
                withCredentials={true}
                onActionChange={(image) => {
                  if (image?.file?.status === 'done') {
                    form.setFieldValue('iconPath', image?.file?.response?.data?.storagePath);
                  }
                }}
              />
            </Form.Item>
            <Form.Item hidden name="iconPath">
              <Input />
            </Form.Item>
            <Divider style={{ background: '#e0e0e0', margin: '16px 0' }} />
            <Form.Item name="imageConfig">
              <ImageCom />
            </Form.Item>
            <Form.Item name="menuUrl" label={formatMessage('uns.menuRouting')} rules={[{ required: true }]}>
              <Input
                // addonBefore={window.location.origin + '/'}
                placeholder={formatMessage('rule.pleaseInput', { label: formatMessage('uns.menuRouting') })}
              />
            </Form.Item>
            {/*<Divider style={{ margin: '16px 0' }} className={styles['add-app-modal-divider']}>*/}
            {/*  <Flex align="center" gap={8} className={styles['add-app-modal-divider-center']} onClick={onOpenAdvanced}>*/}
            {/*    <Settings />*/}
            {/*    <span>{formatMessage('uns.showAdvanced')}</span>*/}
            {/*    {showAdvanced ? <ChevronUp /> : <ChevronDown />}*/}
            {/*  </Flex>*/}
            {/*</Divider>*/}
            <ComEllipsis>{formatMessage('uns.optionalBehaviors')}</ComEllipsis>
            <Row>
              <Col span={12}>
                <Form.Item name="routerTrim" valuePropName="checked">
                  <ComCheckbox
                    tooltip={{
                      title: formatMessage('uns.stripPathPlaceholder'),
                    }}
                  >
                    {formatMessage('uns.stripPath')}
                  </ComCheckbox>
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item name="keepHost" valuePropName="checked">
                  <ComCheckbox
                    tooltip={{
                      title: formatMessage('uns.preserveHostPlaceholder'),
                    }}
                  >
                    {formatMessage('uns.preserveHost')}
                  </ComCheckbox>
                </Form.Item>
              </Col>
            </Row>
            {/*<div style={{ display: showAdvanced ? 'inherit' : 'none' }}>*/}
            <div>
              <Form.Item name="composeYaml">
                <CodeCom />
              </Form.Item>
            </div>
          </Form>
        );
      }}
    </ProModal>
  );
});

interface ImageComProps {
  type: string;
  image?: UploadFile[];
  registry?: string;
}

const CodeCom = ({ value, onChange }: { value?: string; onChange?: (v: string) => void }) => {
  const formatMessage = useTranslate();

  return (
    <div>
      <Flex align="center" gap={8} style={{ marginBottom: 8 }}>
        <ContainerSoftware size={16} />
        <ComEllipsis>{formatMessage('uns.containerConfiguration')}</ComEllipsis>
      </Flex>
      <div
        style={{
          borderRadius: 4,
          border: '1px solid rgb(198, 198, 198)',
          padding: 16,
          position: 'relative',
          overflow: 'hidden',
        }}
        className={codeStyles['custom-theme']}
      >
        <CodeMirror
          theme={codemirrorTheme}
          onChange={onChange}
          value={value}
          height={'200px'}
          extensions={[
            yaml(),
            linter((view) => {
              const diagnostics: Diagnostic[] = [];
              const doc = view.state.doc.toString();
              try {
                jsYaml.load(doc);
              } catch (e: any) {
                if (e.mark) {
                  diagnostics.push({
                    from: e.mark.position,
                    to: e.mark.position + 1,
                    severity: 'error',
                    message: e.reason,
                  });
                }
              }
              return diagnostics;
            }),
            lintGutter(),
          ]}
          placeholder={placeholder}
        />
      </div>
    </div>
  );
};

const ImageCom = ({
  value,
  onChange,
  defaultValue = { type: 'image' },
}: {
  value?: ImageComProps;
  onChange?: (v: ImageComProps) => void;
  defaultValue?: ImageComProps;
}) => {
  const formatMessage = useTranslate();
  const [v, setV] = usePropsValue<ImageComProps>({
    value,
    onChange,
    defaultValue,
  });

  return (
    <>
      <Segmented<string>
        value={v.type}
        onChange={(type) => setV((pre: ImageComProps) => ({ ...pre, type }))}
        options={[
          {
            value: 'image',
            label: formatMessage('uns.uploadImage'),
          },
          {
            value: 'registry',
            label: formatMessage('uns.dockerRegistry'),
          },
        ]}
        style={{ marginBottom: 16 }}
      />
      <div style={{ display: v.type === 'image' ? 'inherit' : 'none' }}>
        <ComDraggerUpload
          acceptList={['tar']}
          className={styles['image-upload']}
          value={v?.image || []}
          size={1024 * 1024 * 1024 * 2}
          action="/inter-api/supos/attachment/upload"
          maxCount={1}
          withCredentials={true}
          onActionChange={(image) => {
            setV((pre: ImageComProps) => ({
              ...pre,
              image: image?.fileList,
              imagePath: image?.fileList?.[0]?.response?.data?.storagePath,
            }));
          }}
        >
          <Flex align="center" justify="center" vertical gap={8}>
            <FolderAdd size={40} style={{ color: '#E0E0E0' }} />
            <ComEllipsis>{formatMessage('uns.clickUploadImage')}</ComEllipsis>
            <ComEllipsis>{formatMessage('common.theFileSizeMax', { size: '2GB' })}</ComEllipsis>
          </Flex>
        </ComDraggerUpload>
      </div>
      <div style={{ display: v.type === 'registry' ? 'inherit' : 'none' }}>
        <Input
          placeholder={formatMessage('uns.dockerExample', { example: 'grafana/grafana:latest' })}
          prefix={<Wikis />}
          value={v.imageUrl}
          onChange={(e) => setV((pre: ImageComProps) => ({ ...pre, imageUrl: e?.target?.value }))}
        />
      </div>
    </>
  );
};

export default AddAppModal;
