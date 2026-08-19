import { useEffect, useState } from 'react';
import { App, Select, Tag, Tooltip } from 'antd';
import { cameraSnapshotUrl, type VisionCameraRegion } from '@/apis/core-api/task';
import type { CameraLiveStatus } from '@/apis/core-api/video';
import HelpTooltip from '@/components/help-tooltip';
import { Checkmark, ChevronDown, CopyFile, Help } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import RegionEditor from './RegionEditor';
import { CAMERA_LIVE_META } from './task-meta';
import styles from './index.module.scss';

type RegionCamera = {
  id: number;
  name: string;
  liveStatus?: CameraLiveStatus;
};

type TaskRegionSectionProps = {
  cameras: RegionCamera[];
  // 空间配置类型由上层根据算法 definition 计算好(spatialModeFor)。
  mode: 'zones' | 'line';
  value: Record<number, VisionCameraRegion>;
  onChange: (regions: Record<number, VisionCameraRegion>) => void;
};

// 创建表单里内联的检测区域/计数线配置:画布左上角切换摄像头,一次只画一台,
// 配好后可一键复制到其余已选摄像头(坐标是归一化的,相对位置保留)。
const TaskRegionSection = ({ cameras, mode, value, onChange }: TaskRegionSectionProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [activeId, setActiveId] = useState<number | undefined>(cameras[0]?.id);

  // 已选摄像头变化后,当前选中的那台可能已被移除。
  useEffect(() => {
    if (cameras.length === 0) return;
    if (!cameras.some((c) => c.id === activeId)) setActiveId(cameras[0].id);
  }, [cameras, activeId]);

  if (cameras.length === 0) return null;
  const activeCamera = cameras.find((c) => c.id === activeId) || cameras[0];

  const setCameraRegion = (camId: number, region: VisionCameraRegion) => onChange({ ...value, [camId]: region });

  const copyToOthers = (camId: number) => {
    const src = value[camId];
    if (!src) return;
    const next = { ...value };
    cameras.forEach((c) => {
      if (c.id !== camId) next[c.id] = JSON.parse(JSON.stringify(src));
    });
    onChange(next);
    message.success(formatMessage('Vision.region.copied'));
  };

  const statusLabel = (status?: CameraLiveStatus) =>
    formatMessage(CAMERA_LIVE_META[status || 'unknown']?.labelKey || 'Vision.task.camUnknown');

  return (
    <div className={styles.regionSection}>
      <div className={styles.regionSectionTitle}>
        {formatMessage(mode === 'line' ? 'Vision.region.configLine' : 'Vision.region.configZone')}
        <HelpTooltip
          title={formatMessage(mode === 'line' ? 'Vision.region.lineHint' : 'Vision.region.zoneHint')}
          icon={<Help size={12} className={styles.regionSectionTip} />}
        />
      </div>
      <RegionEditor
        snapshotUrl={cameraSnapshotUrl(activeCamera.id)}
        mode={mode}
        value={value[activeCamera.id] || {}}
        onChange={(r) => setCameraRegion(activeCamera.id, r)}
        cameraPicker={
          <Select
            size="small"
            className={styles.regionCameraPicker}
            popupClassName={styles.regionCameraDropdown}
            value={activeCamera.id}
            onChange={setActiveId}
            suffixIcon={null}
            // 下拉可比触发框更宽,避免 Online/Offline 标签被裁切。
            popupMatchSelectWidth={false}
            optionLabelProp="title"
            labelRender={({ label }) => (
              <span className={styles.regionCameraPickerValue}>
                <ChevronDown size={16} />
                <Tooltip title={label}>
                  <span className={styles.regionCameraPickerTitle}>{label}</span>
                </Tooltip>
              </span>
            )}
            options={cameras.map((cam) => ({
              value: cam.id,
              title: cam.name,
              label: cam.name,
              liveStatus: cam.liveStatus,
            }))}
            optionRender={(option) => {
              const selected = option.value === activeCamera.id;
              const name = String(option.data.title || '');
              const liveStatus = (option.data as { liveStatus?: CameraLiveStatus }).liveStatus;
              return (
                <span className={styles.regionCameraOption}>
                  <span className={styles.regionCameraOptionCheck} aria-hidden={!selected}>
                    {selected ? <Checkmark size={16} /> : null}
                  </span>
                  <Tooltip title={name}>
                    <span className={styles.regionCameraOptionName}>{name}</span>
                  </Tooltip>
                  <Tag
                    className={styles.regionCameraStatusTag}
                    color={CAMERA_LIVE_META[liveStatus || 'unknown']?.color}
                    bordered={false}
                  >
                    {statusLabel(liveStatus)}
                  </Tag>
                </span>
              );
            }}
          />
        }
        extraActions={
          cameras.length > 1 ? (
            <button
              type="button"
              className={styles.regionStageBtn}
              title={formatMessage('Vision.region.copyToOthers')}
              aria-label={formatMessage('Vision.region.copyToOthers')}
              onClick={() => copyToOthers(activeCamera.id)}
            >
              <CopyFile size={14} />
            </button>
          ) : null
        }
      />
    </div>
  );
};

export default TaskRegionSection;
