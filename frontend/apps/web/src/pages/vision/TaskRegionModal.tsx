import { useEffect, useMemo, useState } from 'react';
import { App, Button, Segmented, Space } from 'antd';
import { cameraSnapshotUrl, saveTaskRegions, type VisionCameraRegion, type VisionTask } from '@/apis/core-api/task';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import RegionEditor from './RegionEditor';
import { spatialModeFor } from './task-meta';
import styles from './index.module.scss';

type TaskRegionModalProps = {
  task: VisionTask | null;
  onClose: () => void;
  onSaved: () => void;
};

const TaskRegionModal = ({ task, onClose, onSaved }: TaskRegionModalProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [cameraId, setCameraId] = useState<number | null>(null);
  const [region, setRegion] = useState<VisionCameraRegion>({});
  const [snapshotUrl, setSnapshotUrl] = useState('');
  const [saving, setSaving] = useState(false);

  // 空间配置类型按算法 definition 计算;'none' 的算法不会打开本弹窗,兜底按 zones。
  const regionMode = spatialModeFor(task?.algorithm);
  const mode: 'zones' | 'line' = regionMode === 'line' ? 'line' : 'zones';
  const cameras = useMemo(() => task?.cameras || [], [task]);

  useEffect(() => {
    // 打开弹窗时默认选中第一路相机(外部触发的初始化,非渲染派生)。
    if (!task || cameras.length === 0) return;
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setCameraId(cameras[0].id);
  }, [task, cameras]);

  useEffect(() => {
    // 切换相机时载入其已存区域与快照(与外部选择同步)。
    if (cameraId == null || !task) return;
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setRegion(task.regions?.[String(cameraId)] || {});
    setSnapshotUrl(cameraSnapshotUrl(cameraId));
  }, [cameraId, task]);

  const save = async () => {
    if (cameraId == null || !task) return;
    setSaving(true);
    try {
      await saveTaskRegions(task.id, {
        cameraId,
        zones: region.zones,
        countingLine: region.countingLine ?? null,
      });
      message.success(formatMessage('common.optsuccess'));
      onSaved();
    } finally {
      setSaving(false);
    }
  };

  return (
    <ProModal
      open={Boolean(task)}
      title={formatMessage(mode === 'line' ? 'Vision.region.titleLine' : 'Vision.region.titleZone')}
      width={720}
      fullScreenable={false}
      onCancel={onClose}
      destroyOnHidden
      maskClosable={false}
    >
      {cameras.length > 1 && (
        <div className={styles.regionCamBar}>
          <Segmented
            value={cameraId ?? undefined}
            onChange={(v) => setCameraId(v as number)}
            options={cameras.map((c) => ({ label: c.name, value: c.id }))}
          />
        </div>
      )}
      {cameraId != null && <RegionEditor snapshotUrl={snapshotUrl} mode={mode} value={region} onChange={setRegion} />}
      <div className={styles.footer}>
        <span />
        <Space>
          <Button onClick={onClose} disabled={saving}>
            {formatMessage('common.cancel')}
          </Button>
          <Button type="primary" loading={saving} onClick={() => void save()}>
            {formatMessage('common.save')}
          </Button>
        </Space>
      </div>
    </ProModal>
  );
};

export default TaskRegionModal;
