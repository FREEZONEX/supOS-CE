import { useEffect, useState } from 'react';
import { App, Button, Form, Input, Select, Space, Upload } from 'antd';
import {
  createVisionAlgorithm,
  updateVisionAlgorithm,
  type AlgorithmType,
  type VisionAlgorithm,
  type VisionAlgorithmCreatePayload,
} from '@/apis/core-api/algorithm';
import { listVisionModels, uploadVisionModel, type VisionModel } from '@/apis/core-api/model';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import ModelStatusTag from './ModelStatusTag';
import styles from './components/CameraFormModal.module.scss';
import pageStyles from './index.module.scss';

type AlgorithmFormValues = {
  name: string;
  algorithmCode: string;
  algoType: AlgorithmType;
  scene?: string;
  description?: string;
  labels?: string[];
  defaultOutput?: string[];
  modelId?: number;
};

type AlgorithmFormModalProps = {
  open: boolean;
  algorithm?: VisionAlgorithm | null;
  onCancel: () => void;
  onSaved: () => void;
};

const algorithmCodePattern = /^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$/;
const ALGO_TYPES: AlgorithmType[] = [
  'object_detection',
  'helmet',
  'intrusion',
  'line_counting',
  'ppe',
  'behavior',
  'interaction',
  'environment',
];

const AlgorithmFormModal = ({ open, algorithm, onCancel, onSaved }: AlgorithmFormModalProps) => {
  const [form] = Form.useForm<AlgorithmFormValues>();
  const [saving, setSaving] = useState(false);
  const [models, setModels] = useState<VisionModel[]>([]);
  const [uploadingModel, setUploadingModel] = useState(false);
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const isEdit = Boolean(algorithm);

  useEffect(() => {
    if (!open) return;
    form.setFieldsValue({
      name: algorithm?.name || '',
      algorithmCode: algorithm?.algorithmCode || '',
      algoType: algorithm?.algoType || 'object_detection',
      scene: algorithm?.scene || '',
      description: algorithm?.description || '',
      labels: algorithm?.labels || [],
      defaultOutput: algorithm?.defaultOutput || [],
      modelId: algorithm?.modelId || undefined,
    });
  }, [algorithm, form, open]);

  useEffect(() => {
    if (!open) return;
    void listVisionModels({ pageNo: 1, pageSize: 200 }).then((res) => setModels(res.data));
  }, [open]);

  // 上传新模型:成功后刷新列表并自动选中。
  const handleModelUpload = (file: File) => {
    setUploadingModel(true);
    void uploadVisionModel(file)
      .then((model) => {
        message.success(formatMessage('Vision.model.uploadSuccess'));
        setModels((prev) => [model, ...prev.filter((m) => m.id !== model.id)]);
        form.setFieldValue('modelId', model.id);
      })
      .finally(() => setUploadingModel(false));
    return false;
  };

  const submit = async () => {
    const values = await form.validateFields();
    setSaving(true);
    try {
      if (algorithm) {
        await updateVisionAlgorithm(algorithm.id, {
          name: values.name.trim(),
          scene: values.scene?.trim(),
          description: values.description?.trim(),
          labels: values.labels,
          defaultOutput: values.defaultOutput,
          modelId: values.modelId ?? 0,
        });
      } else {
        const payload: VisionAlgorithmCreatePayload = {
          algorithmCode: values.algorithmCode.trim(),
          name: values.name.trim(),
          algoType: values.algoType,
          scene: values.scene?.trim(),
          description: values.description?.trim(),
          labels: values.labels,
          defaultOutput: values.defaultOutput,
          modelId: values.modelId ?? 0,
        };
        await createVisionAlgorithm(payload);
      }
      message.success(formatMessage('common.optsuccess'));
      onSaved();
    } finally {
      setSaving(false);
    }
  };

  return (
    <ProModal
      open={open}
      title={
        isEdit
          ? formatMessage('Vision.algorithm.editTitle', { name: algorithm?.name })
          : formatMessage('Vision.algorithm.addTitle')
      }
      width={640}
      fullScreenable={false}
      onCancel={onCancel}
      destroyOnHidden
      maskClosable={false}
    >
      <Form form={form} layout="vertical" className={styles.form} preserve={false}>
        <div className={styles.grid}>
          <Form.Item
            name="algorithmCode"
            label={formatMessage('Vision.algorithm.algorithmId')}
            extra={isEdit ? formatMessage('Vision.algorithm.codeHint') : undefined}
            rules={[
              { required: true, whitespace: true, message: formatMessage('Vision.algorithm.codeRequired') },
              { pattern: algorithmCodePattern, message: formatMessage('Vision.algorithm.codeInvalid') },
            ]}
          >
            <Input maxLength={128} disabled={isEdit} placeholder={formatMessage('Vision.algorithm.codePlaceholder')} />
          </Form.Item>
          <Form.Item
            name="name"
            label={formatMessage('Vision.algorithm.name')}
            rules={[{ required: true, whitespace: true, message: formatMessage('Vision.algorithm.nameRequired') }]}
          >
            <Input maxLength={128} placeholder={formatMessage('Vision.algorithm.namePlaceholder')} />
          </Form.Item>
          <Form.Item
            name="algoType"
            label={formatMessage('Vision.algorithm.type')}
            rules={[{ required: true, message: formatMessage('Vision.algorithm.typeRequired') }]}
          >
            <Select
              disabled={isEdit}
              options={ALGO_TYPES.map((value) => ({
                value,
                label: formatMessage(`Vision.algorithm.type.${value}`),
              }))}
            />
          </Form.Item>
          <Form.Item label={formatMessage('Vision.algorithm.modelFormat')}>
            <Input value="ONNX" disabled />
          </Form.Item>
        </div>
        <Form.Item label={formatMessage('Vision.algorithm.linkedModel')}>
          <div className={styles.modelRow}>
            <Form.Item name="modelId" noStyle>
              <Select
                allowClear
                showSearch
                placeholder={formatMessage('Vision.model.selectPlaceholder')}
                filterOption={(input, opt) =>
                  String(opt?.title ?? '')
                    .toLowerCase()
                    .includes(input.toLowerCase())
                }
                options={models.map((m) => ({
                  value: m.id,
                  title: `${m.name} ${m.version}`,
                  disabled: m.status !== 'available',
                  label: (
                    <span className={pageStyles.camOption}>
                      <span>
                        {m.name}
                        {m.version ? ` ${m.version}` : ''}
                      </span>
                      <ModelStatusTag status={m.status} />
                    </span>
                  ),
                }))}
              />
            </Form.Item>
            <Upload accept=".onnx" showUploadList={false} beforeUpload={handleModelUpload}>
              <Button loading={uploadingModel}>{formatMessage('Vision.model.uploadNew')}</Button>
            </Upload>
          </div>
        </Form.Item>
        <Form.Item name="labels" label={formatMessage('Vision.algorithm.labels')}>
          <Select
            mode="tags"
            tokenSeparators={[',']}
            placeholder={formatMessage('Vision.algorithm.labelsPlaceholder')}
            open={false}
            suffixIcon={null}
          />
        </Form.Item>
        <Form.Item name="scene" label={formatMessage('Vision.algorithm.scene')}>
          <Input maxLength={256} placeholder={formatMessage('Vision.algorithm.scenePlaceholder')} />
        </Form.Item>
        <Form.Item name="description" label={formatMessage('common.description')}>
          <Input.TextArea rows={3} maxLength={1000} showCount />
        </Form.Item>
      </Form>
      <div className={styles.footer}>
        <span />
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

export default AlgorithmFormModal;
