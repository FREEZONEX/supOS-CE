import { useEffect, useMemo, useRef, useState } from 'react';
import { Alert, App, Button, Checkbox, Collapse, Form, Input, InputNumber, Radio, Select, Space, Tag, Tooltip } from 'antd';
import {
  createVisionTask,
  getVisionTask,
  updateVisionTask,
  type VisionCameraRegion,
  type VisionScheduleRule,
  type VisionTask,
  type VisionTaskCreatePayload,
} from '@/apis/core-api/task';
import { listVideoCameras, type VideoCamera } from '@/apis/core-api/video';
import { listVisionAlgorithms, type VisionAlgorithm } from '@/apis/core-api/algorithm';
import { ComDrawer } from '@/components';
import { useTranslate } from '@/hooks';
import ScheduleEditor from './ScheduleEditor';
import TaskRegionSection from './TaskRegionSection';
import {
  CAMERA_LIVE_META,
  DETECTION_FREQUENCY_OPTIONS,
  FREQUENCY_TO_MS,
  MODEL_STATUS_META,
  outputOptionsFor,
  RETENTION_OPTIONS,
  showsPushIntervalFor,
  spatialModeFor,
} from './task-meta';
import formStyles from './components/CameraFormModal.module.scss';
import styles from './index.module.scss';

type TaskFormValues = {
  name: string;
  note?: string;
  cameraIds: number[];
  algorithmId: number;
  outputCategories?: string[];
  scheduleMode: 'allDay' | 'custom';
  resultPushIntervalMs?: number;
  detectionFrequency: 'low' | 'standard' | 'high' | 'custom';
  samplingIntervalMs?: number;
  confThreshold: number;
  iouThreshold: number;
  inputWidth: number;
  inputHeight: number;
  resultRetention: string;
};

type TaskFormModalProps = {
  open: boolean;
  onCancel: () => void;
  onSaved: () => void;
  editTask?: VisionTask | null;
};

const DEFAULTS: Partial<TaskFormValues> = {
  scheduleMode: 'allDay',
  detectionFrequency: 'standard',
  samplingIntervalMs: 200,
  confThreshold: 0.5,
  iouThreshold: 0.45,
  inputWidth: 640,
  inputHeight: 640,
  resultRetention: 'permanent',
  resultPushIntervalMs: 0,
};

const TaskFormModal = ({ open, onCancel, onSaved, editTask }: TaskFormModalProps) => {
  const [form] = Form.useForm<TaskFormValues>();
  const [saving, setSaving] = useState(false);
  const [cameras, setCameras] = useState<VideoCamera[]>([]);
  const [algorithms, setAlgorithms] = useState<VisionAlgorithm[]>([]);
  const [scheduleRules, setScheduleRules] = useState<VisionScheduleRule[]>([]);
  const [regions, setRegions] = useState<Record<number, VisionCameraRegion>>({});
  // 编辑回填时跳过一次算法联动(避免把回填的 outputCategories 覆盖为全选)。
  const skipAutoOutputRef = useRef(false);
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const algorithmId = Form.useWatch('algorithmId', form);
  const cameraIds = Form.useWatch('cameraIds', form);
  const taskName = Form.useWatch('name', form);
  const scheduleMode = Form.useWatch('scheduleMode', form);
  const detectionFrequency = Form.useWatch('detectionFrequency', form);

  const selectedAlgorithm = useMemo(() => algorithms.find((a) => a.id === algorithmId), [algorithms, algorithmId]);
  // 输出项/空间配置/推送间隔均由算法 definition 驱动(缺省时按 algoType fallback)。
  const outputOptions = selectedAlgorithm ? outputOptionsFor(selectedAlgorithm) : [];
  const regionMode = spatialModeFor(selectedAlgorithm);
  const showPushInterval = showsPushIntervalFor(selectedAlgorithm);
  const selectedCameras = useMemo(
    () =>
      cameras
        .filter((c) => (cameraIds || []).includes(c.id))
        .map((c) => ({ id: c.id, name: c.name, cameraCode: c.cameraCode, liveStatus: c.liveStatus })),
    [cameras, cameraIds]
  );
  const hasOfflineCamera = selectedCameras.some((c) => c.liveStatus === 'offline');

  useEffect(() => {
    if (!open) return;
    form.resetFields();
    skipAutoOutputRef.current = false;
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setScheduleRules([]);

    setRegions({});
    void listVideoCameras({ pageNo: 1, pageSize: 200 }).then((res) => setCameras(res.data));
    void listVisionAlgorithms({ pageNo: 1, pageSize: 200 }).then((res) => setAlgorithms(res.data));

    if (!editTask) return;
    // 编辑模式:拉完整任务并回填表单。
    void getVisionTask(editTask.id).then((task) => {
      // 回填后跳过一次算法联动,保留已选输出类别。
      skipAutoOutputRef.current = true;
      form.setFieldsValue({
        name: task.name,
        note: task.note,
        cameraIds: task.cameraIds,
        algorithmId: task.algorithmId,
        outputCategories: task.outputCategories,
        confThreshold: task.params?.confThreshold ?? DEFAULTS.confThreshold,
        iouThreshold: task.params?.iouThreshold ?? DEFAULTS.iouThreshold,
        inputWidth: task.params?.inputWidth ?? DEFAULTS.inputWidth,
        inputHeight: task.params?.inputHeight ?? DEFAULTS.inputHeight,
        detectionFrequency: task.params?.detectionFrequency ?? DEFAULTS.detectionFrequency,
        samplingIntervalMs: task.params?.samplingIntervalMs ?? DEFAULTS.samplingIntervalMs,
        resultRetention: task.params?.resultRetention ?? DEFAULTS.resultRetention,
        resultPushIntervalMs: task.params?.resultPushIntervalMs ?? DEFAULTS.resultPushIntervalMs,
        scheduleMode: task.schedule?.mode === 'custom' ? 'custom' : 'allDay',
      });
      setScheduleRules(task.schedule?.rules || []);
      const nextRegions: Record<number, VisionCameraRegion> = {};
      Object.entries(task.regions || {}).forEach(([cid, region]) => {
        nextRegions[Number(cid)] = region;
      });
      setRegions(nextRegions);
    });
  }, [form, open, editTask]);

  // 切换算法后,输出类别默认全选该算法的可选项(编辑回填首次跳过)。
  useEffect(() => {
    if (skipAutoOutputRef.current) {
      skipAutoOutputRef.current = false;
      return;
    }
    if (selectedAlgorithm) {
      form.setFieldValue('outputCategories', outputOptionsFor(selectedAlgorithm));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAlgorithm]);

  // 非 custom 频率:采样间隔跟随预设并只读。
  useEffect(() => {
    if (detectionFrequency && detectionFrequency !== 'custom') {
      form.setFieldValue('samplingIntervalMs', FREQUENCY_TO_MS[detectionFrequency]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [detectionFrequency]);

  const submit = async () => {
    const v = await form.validateFields();
    const payload: VisionTaskCreatePayload = {
      name: v.name.trim(),
      note: v.note?.trim(),
      cameraIds: v.cameraIds,
      algorithmId: v.algorithmId,
      outputCategories: v.outputCategories,
      params: {
        confThreshold: v.confThreshold,
        iouThreshold: v.iouThreshold,
        inputWidth: v.inputWidth,
        inputHeight: v.inputHeight,
        detectionFrequency: v.detectionFrequency,
        samplingIntervalMs: v.samplingIntervalMs,
        resultRetention: v.resultRetention as NonNullable<VisionTaskCreatePayload['params']>['resultRetention'],
        resultPushIntervalMs: showPushInterval ? v.resultPushIntervalMs : 0,
      },
      schedule: v.scheduleMode === 'custom' ? { mode: 'custom', rules: scheduleRules } : { mode: 'allDay' },
      regions: v.cameraIds
        .filter((camId) => regions[camId]?.zones?.length || regions[camId]?.countingLine)
        .map((camId) => ({
          cameraId: camId,
          zones: regions[camId]?.zones,
          countingLine: regions[camId]?.countingLine ?? null,
        })),
    };
    setSaving(true);
    try {
      if (editTask) {
        await updateVisionTask(editTask.id, payload);
      } else {
        await createVisionTask(payload);
      }
      message.success(formatMessage(editTask ? 'Vision.task.updateSuccess' : 'common.optsuccess'));
      onSaved();
    } finally {
      setSaving(false);
    }
  };

  return (
    <ComDrawer
      open={open}
      title={formatMessage(editTask ? 'Vision.task.editTitle' : 'Vision.task.addTitle')}
      width={720}
      onClose={onCancel}
      destroyOnClose
      maskClosable={false}
      footer={
        <div className={styles.drawerFooter}>
          <Space>
            <Button onClick={onCancel} disabled={saving}>
              {formatMessage('common.cancel')}
            </Button>
            <Button type="primary" loading={saving} onClick={() => void submit()}>
              {formatMessage(editTask ? 'common.save' : 'common.create')}
            </Button>
          </Space>
        </div>
      }
    >
      <div className={styles.taskDrawerBody}>
        <Form
          form={form}
          layout="horizontal"
          labelAlign="left"
          colon={false}
          labelCol={{ flex: '0 0 128px' }}
          wrapperCol={{ flex: '1 1 auto' }}
          preserve={false}
          initialValues={DEFAULTS}
          // 设计稿里星号跟在标签文字后面,antd 默认放前面,这里自定义渲染。
          requiredMark={(label, { required }) => (
            <>
              {label}
              {required && <span className={styles.requiredMark}>*</span>}
            </>
          )}
        >
          <Form.Item
            name="name"
            label={formatMessage('Vision.task.name')}
            rules={[{ required: true, whitespace: true, message: formatMessage('Vision.task.nameRequired') }]}
          >
            <Input maxLength={128} placeholder={formatMessage('Vision.task.namePlaceholder')} />
          </Form.Item>
          <Form.Item name="note" label={formatMessage('Vision.task.note')}>
            <Input.TextArea rows={2} maxLength={1000} showCount />
          </Form.Item>
          <Form.Item
            name="cameraIds"
            label={formatMessage('Vision.task.cameras')}
            rules={[{ required: true, message: formatMessage('Vision.task.cameraRequired') }]}
          >
            <Select
              mode="multiple"
              className={styles.cameraMultiSelect}
              popupClassName={styles.cameraMultiDropdown}
              placeholder={formatMessage('Vision.task.cameraPlaceholder')}
              // 通用做法:按可用宽度自适应展示标签,超出用 +N,避免换行/滚动条。
              maxTagCount="responsive"
              maxTagPlaceholder={(omitted) => {
                const names = omitted
                  .map((item) => String(item.label ?? ''))
                  .filter(Boolean);
                return (
                  <Tooltip
                    title={
                      <div className={styles.cameraRestTooltip}>
                        {names.map((name) => (
                          <div key={name}>{name}</div>
                        ))}
                      </div>
                    }
                  >
                    <span>+{omitted.length}</span>
                  </Tooltip>
                );
              }}
              filterOption={(input, opt) =>
                String(opt?.title ?? '')
                  .toLowerCase()
                  .includes(input.toLowerCase())
              }
              notFoundContent={<span className={styles.selectEmpty}>{formatMessage('Vision.task.noCamera')}</span>}
              options={cameras.map((cam) => ({
                value: cam.id,
                // 选中标签只展示相机名(对齐 Figma 13652-56275)。
                label: cam.name,
                title: `${cam.name} ${cam.cameraCode}`,
                cameraCode: cam.cameraCode,
                liveStatus: cam.liveStatus,
              }))}
              // 勾选态由左侧自定义 checkbox 展示,隐藏 antd 默认右侧勾。
              menuItemSelectedIcon={null}
              optionRender={(option) => {
                const cameraCode = String((option.data as { cameraCode?: string }).cameraCode || '');
                const liveStatus = (option.data as { liveStatus?: VideoCamera['liveStatus'] }).liveStatus;
                const selected = (cameraIds || []).includes(option.value as number);
                return (
                  <span className={styles.camOption}>
                    <span className={styles.camOptionCheck} aria-hidden>
                      <Checkbox checked={selected} tabIndex={-1} />
                    </span>
                    <span className={styles.camOptionMain}>
                      <span className={styles.camOptionName}>{option.data.label}</span>
                      {cameraCode ? <span className={styles.camOptionCode}>{cameraCode}</span> : null}
                    </span>
                    <Tag color={CAMERA_LIVE_META[liveStatus || 'unknown']?.color} bordered={false}>
                      {formatMessage(CAMERA_LIVE_META[liveStatus || 'unknown']?.labelKey || 'Vision.task.camUnknown')}
                    </Tag>
                  </span>
                );
              }}
            />
          </Form.Item>
          {(hasOfflineCamera || selectedCameras.length > 0) && (
            <Form.Item label=" " colon={false} className={styles.subItem}>
              {hasOfflineCamera && (
                <Alert
                  className={styles.formAlert}
                  type="warning"
                  showIcon
                  message={formatMessage('Vision.task.cameraOfflineWarn')}
                />
              )}
              {selectedCameras.length > 0 && (
                <div className={styles.unsPreview}>
                  <div className={styles.unsTitle}>{formatMessage('Vision.task.generatedTopic')}</div>
                  <div className={styles.unsList}>
                    {selectedCameras.map((c) => (
                      <div key={c.id} className={styles.unsRow}>
                        <span className={styles.unsCam}>{c.name}</span>
                        {/* 与后端 unsTopicForTaskCamera 规则一致:Vision/<任务名>/Metric/<相机码> */}
                        <code className={styles.unsTopic}>
                          Vision/{(taskName || '').trim().replaceAll('/', '-') || '<task>'}/Metric/{c.cameraCode}
                        </code>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </Form.Item>
          )}
          <Form.Item
            name="algorithmId"
            label={formatMessage('Vision.task.algorithm')}
            rules={[{ required: true, message: formatMessage('Vision.task.algorithmRequired') }]}
          >
            <Select
              placeholder={formatMessage('Vision.task.algorithmPlaceholder')}
              filterOption={(input, opt) =>
                String(opt?.title ?? '')
                  .toLowerCase()
                  .includes(input.toLowerCase())
              }
              options={algorithms.map((algo) => ({
                value: algo.id,
                title: algo.name,
                disabled: algo.modelStatus !== 'available',
                label: (
                  <span className={styles.camOption}>
                    <span>{algo.name}</span>
                    <Tag color={MODEL_STATUS_META[algo.modelStatus]?.color} bordered={false}>
                      {formatMessage(MODEL_STATUS_META[algo.modelStatus]?.labelKey || 'Vision.task.modelMissing')}
                    </Tag>
                  </span>
                ),
              }))}
            />
          </Form.Item>

          {outputOptions.length > 0 && (
            <Form.Item
              name="outputCategories"
              label={formatMessage('Vision.task.outputCategories')}
              rules={[{ required: true, message: formatMessage('Vision.task.outputRequired') }]}
            >
              <Checkbox.Group options={outputOptions.map((value) => ({ label: value, value }))} />
            </Form.Item>
          )}

          {selectedAlgorithm && showPushInterval && (
            <Form.Item
              name="resultPushIntervalMs"
              label={formatMessage('Vision.task.pushInterval')}
              tooltip={formatMessage('Vision.task.pushIntervalHint')}
            >
              <InputNumber min={0} step={100} addonAfter="ms" style={{ width: 180 }} />
            </Form.Item>
          )}

          {/* 设计稿里 Run Schedule 整段上下各有一条分隔线 */}
          <div className={styles.formDivider} />
          <Form.Item name="scheduleMode" label={formatMessage('Vision.task.runSchedule')} required>
            <Radio.Group>
              <Radio value="allDay">{formatMessage('Vision.task.runAllDay')}</Radio>
              <Radio value="custom">{formatMessage('Vision.task.customSchedule')}</Radio>
            </Radio.Group>
          </Form.Item>
          {scheduleMode === 'custom' && <ScheduleEditor value={scheduleRules} onChange={setScheduleRules} />}
          <div className={styles.formDivider} />

          {selectedAlgorithm && regionMode !== 'none' && selectedCameras.length > 0 && (
            <TaskRegionSection cameras={selectedCameras} mode={regionMode} value={regions} onChange={setRegions} />
          )}

          <Collapse
            ghost
            className={styles.advancedCollapse}
            items={[
              {
                key: 'advanced',
                label: formatMessage('Vision.task.advancedSettings'),
                children: (
                  <div className={styles.advancedGrid}>
                    <Form.Item
                      name="detectionFrequency"
                      label={formatMessage('Vision.task.detectionFrequency')}
                      tooltip={formatMessage('Vision.task.detectionFrequencyHint')}
                    >
                      <Select
                        options={DETECTION_FREQUENCY_OPTIONS.map((o) => ({
                          value: o.value,
                          label: formatMessage(o.labelKey),
                        }))}
                      />
                    </Form.Item>
                    <Form.Item
                      name="samplingIntervalMs"
                      label={formatMessage('Vision.task.samplingInterval')}
                      tooltip={formatMessage('Vision.task.samplingIntervalHint')}
                    >
                      <InputNumber
                        min={50}
                        step={10}
                        addonAfter="ms"
                        style={{ width: '100%' }}
                        disabled={detectionFrequency !== 'custom'}
                      />
                    </Form.Item>
                    <Form.Item
                      name="confThreshold"
                      label={formatMessage('Vision.task.confThreshold')}
                      tooltip={formatMessage('Vision.task.confThresholdHint')}
                    >
                      <InputNumber min={0.01} max={1} step={0.01} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                      name="iouThreshold"
                      label={formatMessage('Vision.task.iouThreshold')}
                      tooltip={formatMessage('Vision.task.iouThresholdHint')}
                    >
                      <InputNumber min={0.01} max={1} step={0.01} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                      name="inputWidth"
                      label={formatMessage('Vision.task.inputWidth')}
                      tooltip={formatMessage('Vision.task.inputSizeHint')}
                    >
                      <InputNumber min={64} step={32} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                      name="inputHeight"
                      label={formatMessage('Vision.task.inputHeight')}
                      tooltip={formatMessage('Vision.task.inputSizeHint')}
                    >
                      <InputNumber min={64} step={32} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item
                      className={styles.advancedFull}
                      name="resultRetention"
                      label={formatMessage('Vision.task.resultRetention')}
                      tooltip={formatMessage('Vision.task.resultRetentionHint')}
                    >
                      <Select
                        options={RETENTION_OPTIONS.map((o) => ({ value: o.value, label: formatMessage(o.labelKey) }))}
                      />
                    </Form.Item>
                  </div>
                ),
              },
            ]}
          />
        </Form>

        {selectedAlgorithm && selectedAlgorithm.modelStatus !== 'available' && (
          <div className={formStyles.testResult}>
            <Tag color="warning" bordered={false}>
              {formatMessage('Vision.task.modelNotReady')}
            </Tag>
          </div>
        )}
      </div>
    </ComDrawer>
  );
};

export default TaskFormModal;
