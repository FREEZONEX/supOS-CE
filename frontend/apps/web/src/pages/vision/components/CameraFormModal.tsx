import { useEffect, useState } from 'react';
import { App, Alert, Button, Form, Input, Select, Space } from 'antd';
import {
  createVideoCamera,
  probeVideoCamera,
  updateVideoCamera,
  type CameraProbeResult,
  type CameraRtpType,
  type CameraStatus,
  type CameraTransport,
  type VideoCamera,
  type VideoCameraPayload,
} from '@/apis/core-api/video';
import { CheckmarkFilled, ErrorFilled, View, WarningFilled } from '@/components/lucide-icon/carbon';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import styles from './CameraFormModal.module.scss';

type CameraFormValues = VideoCameraPayload;

type CameraFormModalProps = {
  open: boolean;
  camera: VideoCamera | null;
  onCancel: () => void;
  onSaved: () => void;
};

const cameraCodePattern = /^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$/;

const CameraFormModal = ({ open, camera, onCancel, onSaved }: CameraFormModalProps) => {
  const [form] = Form.useForm<CameraFormValues>();
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<CameraProbeResult | null>(null);
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const transport = Form.useWatch('transport', form);
  const isEdit = Boolean(camera);

  useEffect(() => {
    if (!open) return;
    form.setFieldsValue({
      cameraCode: camera?.cameraCode || '',
      name: camera?.name || '',
      location: camera?.location || '',
      description: camera?.description || '',
      status: camera?.status || 1,
      sourceUrl: camera?.sourceUrl || '',
      transport: camera?.transport || 'rtsp',
      rtpType: camera?.rtpType || 'tcp',
    });
  }, [camera, form, open]);

  const submit = async () => {
    const values = await form.validateFields();
    const payload: VideoCameraPayload = {
      ...values,
      cameraCode: values.cameraCode.trim(),
      name: values.name.trim(),
      location: values.location?.trim(),
      description: values.description?.trim(),
      sourceUrl: values.sourceUrl.trim(),
    };
    setSaving(true);
    try {
      if (camera) {
        await updateVideoCamera(camera.id, payload);
      } else {
        await createVideoCamera(payload);
      }
      message.success(formatMessage('common.optsuccess'));
      onSaved();
    } finally {
      setSaving(false);
    }
  };

  const testConnection = async () => {
    const values = await form.validateFields(['transport', 'rtpType', 'sourceUrl']);
    const sourceUrl = values.sourceUrl.trim();
    // 编辑态且地址未改、且已保存凭据时，按摄像头复用已存凭据；否则按填写地址探测。
    const reuseSaved = isEdit && camera?.hasCredentials && sourceUrl === (camera?.sourceUrl || '');
    setTesting(true);
    setTestResult(null);
    try {
      const result = await probeVideoCamera(
        reuseSaved ? { cameraId: camera?.id } : { sourceUrl, transport: values.transport, rtpType: values.rtpType }
      );
      setTestResult(result);
      // 三态提示:可达且出画 / 可达但无预览 / 不可达。
      if (!result.reachable) {
        message.error(formatMessage('Vision.camera.testFailed'));
      } else if (result.resolution) {
        message.success(formatMessage('Vision.camera.testConnected'));
      } else {
        message.warning(formatMessage('Vision.camera.testNoPreviewToast'));
      }
    } catch (error) {
      // 后端对地址格式非法直接返回参数错误,单独给出可操作的提示。
      const code = (error as { code?: number })?.code;
      setTestResult({ reachable: false, resolution: '' });
      if (code === 400) message.error(formatMessage('Vision.camera.testInvalidUrl'));
    } finally {
      setTesting(false);
    }
  };

  const resetTestResult = () => setTestResult(null);

  return (
    <ProModal
      open={open}
      title={formatMessage(isEdit ? 'Vision.camera.editTitle' : 'Vision.camera.addTitle')}
      width={640}
      fullScreenable={false}
      onCancel={onCancel}
      destroyOnHidden
      maskClosable={false}
    >
      <Form
        form={form}
        layout="vertical"
        className={styles.form}
        preserve={false}
        onValuesChange={(changed) => {
          // 修改地址或传输协议后，原测试结果失效（文档 4.5）。
          if ('sourceUrl' in changed || 'transport' in changed || 'rtpType' in changed) {
            resetTestResult();
          }
        }}
      >
        <div className={styles.grid}>
          <Form.Item
            name="name"
            label={formatMessage('Vision.camera.name')}
            rules={[{ required: true, whitespace: true, message: formatMessage('Vision.camera.nameRequired') }]}
          >
            <Input maxLength={128} placeholder={formatMessage('Vision.camera.namePlaceholder')} />
          </Form.Item>
          <Form.Item
            name="cameraCode"
            label={formatMessage('Vision.camera.code')}
            extra={isEdit ? formatMessage('Vision.camera.codeImmutableHint') : undefined}
            rules={[
              { required: true, whitespace: true, message: formatMessage('Vision.camera.codeRequired') },
              { pattern: cameraCodePattern, message: formatMessage('Vision.camera.codeInvalid') },
            ]}
          >
            <Input maxLength={128} disabled={isEdit} placeholder={formatMessage('Vision.camera.codePlaceholder')} />
          </Form.Item>
          <Form.Item name="transport" label={formatMessage('Vision.camera.inputType')} rules={[{ required: true }]}>
            <Select
              options={[
                { value: 'rtsp' satisfies CameraTransport, label: 'RTSP' },
                { value: 'rtmp_pull' satisfies CameraTransport, label: formatMessage('Vision.camera.rtmpPull') },
              ]}
            />
          </Form.Item>
          <Form.Item name="rtpType" label={formatMessage('Vision.camera.protocol')} rules={[{ required: true }]}>
            <Select
              disabled={transport === 'rtmp_pull'}
              options={[
                { value: 'tcp' satisfies CameraRtpType, label: 'TCP' },
                { value: 'udp' satisfies CameraRtpType, label: 'UDP' },
              ]}
            />
          </Form.Item>
        </div>
        <Form.Item
          name="sourceUrl"
          label={formatMessage('Vision.camera.streamUrl')}
          extra={isEdit && camera?.hasCredentials ? formatMessage('Vision.camera.credentialsRetained') : undefined}
          dependencies={['transport']}
          rules={[
            { required: true, whitespace: true, message: formatMessage('Vision.camera.streamUrlRequired') },
            {
              validator: (_, value: string) => {
                const scheme = String(value || '')
                  .trim()
                  .split(':')[0]
                  ?.toLowerCase();
                const valid =
                  transport === 'rtmp_pull' ? ['rtmp', 'rtmps'].includes(scheme) : ['rtsp', 'rtsps'].includes(scheme);
                return valid
                  ? Promise.resolve()
                  : Promise.reject(new Error(formatMessage('Vision.camera.streamUrlInvalid')));
              },
            },
          ]}
        >
          <Input
            placeholder={transport === 'rtmp_pull' ? 'rtmp://host/app/stream' : 'rtsp://user:password@host/path'}
          />
        </Form.Item>
        <div className={styles.grid}>
          <Form.Item name="location" label={formatMessage('Vision.camera.location')}>
            <Input maxLength={256} placeholder={formatMessage('Vision.camera.locationPlaceholder')} />
          </Form.Item>
          <Form.Item name="status" label={formatMessage('Vision.camera.status')} rules={[{ required: true }]}>
            <Select
              options={[
                { value: 1 satisfies CameraStatus, label: formatMessage('common.enable') },
                { value: 2 satisfies CameraStatus, label: formatMessage('common.disable') },
              ]}
            />
          </Form.Item>
        </div>
        <Form.Item name="description" label={formatMessage('common.description')}>
          <Input.TextArea rows={3} maxLength={1000} showCount />
        </Form.Item>
      </Form>
      {/* 失败/无预览按设计给出可操作说明;成功则不用横幅,直接展示预览。 */}
      {testResult && !testResult.reachable && (
        <Alert
          className={styles.testResult}
          type="error"
          showIcon
          message={formatMessage('Vision.camera.testFailedDetail')}
        />
      )}
      {testResult?.reachable && !testResult.resolution && (
        <Alert
          className={styles.testResult}
          type="warning"
          showIcon
          message={formatMessage('Vision.camera.testNoPreviewDetail')}
        />
      )}
      {testResult?.reachable && testResult.resolution && (
        <div className={styles.previewBlock}>
          {testResult.snapshot ? (
            <img
              className={styles.previewImg}
              src={testResult.snapshot}
              alt={formatMessage('Vision.camera.previewBlock')}
            />
          ) : (
            <div className={styles.previewPlaceholder}>
              <View size={20} />
              <span>{formatMessage('Vision.camera.previewBlock')}</span>
            </div>
          )}
        </div>
      )}
      <div className={styles.footer}>
        {testResult ? (
          <span
            className={`${styles.testState} ${
              !testResult.reachable
                ? styles.testStateFailed
                : testResult.resolution
                  ? styles.testStateOk
                  : styles.testStateWarn
            }`}
          >
            {!testResult.reachable ? (
              <ErrorFilled size={16} />
            ) : testResult.resolution ? (
              <CheckmarkFilled size={16} />
            ) : (
              <WarningFilled size={16} />
            )}
            {formatMessage(
              !testResult.reachable
                ? 'Vision.camera.stateFailed'
                : testResult.resolution
                  ? 'Vision.camera.stateConnected'
                  : 'Vision.camera.stateVerified'
            )}
          </span>
        ) : (
          <Button loading={testing} disabled={saving} onClick={() => void testConnection()}>
            {formatMessage('Vision.camera.testConnection')}
          </Button>
        )}
        <Space>
          <Button onClick={onCancel} disabled={saving}>
            {formatMessage('common.cancel')}
          </Button>
          <Button type="primary" loading={saving} onClick={() => void submit()}>
            {formatMessage('common.confirm')}
          </Button>
        </Space>
      </div>
    </ProModal>
  );
};

export default CameraFormModal;
