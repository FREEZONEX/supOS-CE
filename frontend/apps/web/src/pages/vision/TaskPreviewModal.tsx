import { useEffect, useRef, useState } from 'react';
import { Button, Empty, Tabs, Tag } from 'antd';
import { previewVideoCamera, stopVideoCameraPreview } from '@/apis/core-api/video';
import type { VisionCameraRegion, VisionTask } from '@/apis/core-api/task';
import { getTaskLatestResult, type VisionLatestResult } from '@/apis/core-api/vision-results';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import styles from './components/TaskPreviewModal.module.scss';

type PreviewPhase = 'connecting' | 'playing' | 'failed';

type DetectionObject = {
  label?: string;
  confidence?: number;
  bbox?: { x: number; y: number; width: number; height: number };
  trackId?: number;
};

// 识别结果超过该时长未更新就不再叠框(任务停了或链路断了,避免残影)。
const DETECTION_STALE_MS = 5000;
const DETECTION_POLL_MS = 1000;

// 视频用 object-fit: contain 渲染,把源分辨率像素坐标换算到显示坐标(含黑边偏移)。
const fitTransform = (videoW: number, videoH: number, boxW: number, boxH: number) => {
  if (!videoW || !videoH || !boxW || !boxH) return null;
  const scale = Math.min(boxW / videoW, boxH / videoH);
  return {
    scale,
    offsetX: (boxW - videoW * scale) / 2,
    offsetY: (boxH - videoH * scale) / 2,
  };
};

// 等待 ICE 收集完成后再发送 offer（ZLM WHEP 使用非 trickle 模式）。
const waitForIceGathering = (pc: RTCPeerConnection) =>
  new Promise<void>((resolve) => {
    if (pc.iceGatheringState === 'complete') {
      resolve();
      return;
    }
    const check = () => {
      if (pc.iceGatheringState === 'complete') {
        pc.removeEventListener('icegatheringstatechange', check);
        resolve();
      }
    };
    pc.addEventListener('icegatheringstatechange', check);
    // 兜底超时，避免个别浏览器长时间不触发 complete。
    setTimeout(resolve, 2000);
  });

// 单相机 WHEP 实时画面：每个实例独立 pc + video ref + 失败态。
// 挂载时建连，卸载时清理（stopVideoCameraPreview + pc.close），失败可重试。
// detection 为该相机最新一帧识别结果,叠加 SVG 检测框(近实时,受结果推送间隔限制)。
const CameraLiveView = ({
  cameraId,
  name,
  detection,
  region,
}: {
  cameraId: number;
  name?: string;
  detection?: VisionLatestResult;
  region?: VisionCameraRegion;
}) => {
  const formatMessage = useTranslate();
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const stageRef = useRef<HTMLDivElement | null>(null);
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const [phase, setPhase] = useState<PreviewPhase>('connecting');
  const [resolution, setResolution] = useState('');
  const [failReason, setFailReason] = useState('');
  const [retryToken, setRetryToken] = useState(0);
  const [stageSize, setStageSize] = useState({ width: 0, height: 0 });

  // 跟踪显示区尺寸(弹窗打开动画/窗口缩放都会变)。
  useEffect(() => {
    const stage = stageRef.current;
    if (!stage || typeof ResizeObserver === 'undefined') return undefined;
    const observer = new ResizeObserver((entries) => {
      const rect = entries[0]?.contentRect;
      if (rect) setStageSize({ width: rect.width, height: rect.height });
    });
    observer.observe(stage);
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    let disposed = false;
    const videoEl = videoRef.current;

    const pc = new RTCPeerConnection();
    pcRef.current = pc;
    pc.addTransceiver('video', { direction: 'recvonly' });
    pc.addTransceiver('audio', { direction: 'recvonly' });
    pc.ontrack = (event) => {
      if (videoRef.current && event.streams[0]) {
        videoRef.current.srcObject = event.streams[0];
      }
    };
    pc.onconnectionstatechange = () => {
      if (disposed) return;
      if (pc.connectionState === 'failed' || pc.connectionState === 'disconnected') {
        setPhase('failed');
        setFailReason(formatMessage('Vision.camera.previewConnLost'));
      }
    };

    const start = async () => {
      try {
        const offer = await pc.createOffer();
        await pc.setLocalDescription(offer);
        await waitForIceGathering(pc);
        const answer = await previewVideoCamera(cameraId, pc.localDescription?.sdp || offer.sdp || '');
        if (disposed) return;
        await pc.setRemoteDescription({ type: 'answer', sdp: answer.sdp });
      } catch (error) {
        if (disposed) return;
        setPhase('failed');
        setFailReason(error instanceof Error ? error.message : formatMessage('Vision.camera.previewFailed'));
      }
    };
    void start();

    return () => {
      disposed = true;
      pc.ontrack = null;
      pc.onconnectionstatechange = null;
      pc.close();
      pcRef.current = null;
      if (videoEl) {
        videoEl.srcObject = null;
      }
      void stopVideoCameraPreview(cameraId).catch(() => undefined);
    };
  }, [cameraId, formatMessage, retryToken]);

  const handlePlaying = () => {
    setPhase('playing');
    const video = videoRef.current;
    if (video && video.videoWidth && video.videoHeight) {
      setResolution(`${video.videoWidth}x${video.videoHeight}`);
    }
  };

  // 计算当前应叠加的检测框(显示坐标系);结果过期或视频未就绪时为空。
  // 用 effect 而非渲染期计算:需要读取 video ref 与当前时间(过期判定)。
  const [detectionBoxes, setDetectionBoxes] = useState<
    { key: string; x: number; y: number; width: number; height: number; label: string }[]
  >([]);
  useEffect(() => {
    if (phase !== 'playing' || !detection || Date.now() - detection.updatedAt > DETECTION_STALE_MS) {
      setDetectionBoxes([]);
      return;
    }
    const video = videoRef.current;
    const transform = fitTransform(video?.videoWidth || 0, video?.videoHeight || 0, stageSize.width, stageSize.height);
    if (!transform) {
      setDetectionBoxes([]);
      return;
    }
    const objects = (detection.payload?.objects as DetectionObject[] | undefined) || [];
    setDetectionBoxes(
      objects
        .filter((obj) => obj.bbox && obj.bbox.width > 0 && obj.bbox.height > 0)
        .map((obj, index) => {
          const bbox = obj.bbox as NonNullable<DetectionObject['bbox']>;
          return {
            key: `${obj.trackId ?? index}-${obj.label ?? ''}`,
            x: transform.offsetX + bbox.x * transform.scale,
            y: transform.offsetY + bbox.y * transform.scale,
            width: bbox.width * transform.scale,
            height: bbox.height * transform.scale,
            label: `${obj.label ?? ''}${obj.confidence ? ` ${(obj.confidence * 100).toFixed(0)}%` : ''}`.trim(),
          };
        })
    );
  }, [phase, detection, stageSize]);

  // 用户划定的检测区域/计数线(归一化坐标 → 显示坐标),随视频尺寸变化重算。
  const [regionShapes, setRegionShapes] = useState<{
    zones: { key: string; points: string; name?: string; labelX: number; labelY: number }[];
    line: { x1: number; y1: number; x2: number; y2: number } | null;
  }>({ zones: [], line: null });
  useEffect(() => {
    if (phase !== 'playing' || !region) {
      setRegionShapes({ zones: [], line: null });
      return;
    }
    const video = videoRef.current;
    const videoW = video?.videoWidth || 0;
    const videoH = video?.videoHeight || 0;
    const transform = fitTransform(videoW, videoH, stageSize.width, stageSize.height);
    if (!transform) {
      setRegionShapes({ zones: [], line: null });
      return;
    }
    const toX = (nx: number) => transform.offsetX + nx * videoW * transform.scale;
    const toY = (ny: number) => transform.offsetY + ny * videoH * transform.scale;
    const zones = (region.zones || [])
      .filter((zone) => zone.enabled !== false && (zone.points?.length ?? 0) >= 3)
      .map((zone, index) => ({
        key: `${zone.name ?? 'zone'}-${index}`,
        points: zone.points.map((p) => `${toX(p.x)},${toY(p.y)}`).join(' '),
        name: zone.name,
        labelX: toX(Math.min(...zone.points.map((p) => p.x))),
        labelY: Math.max(toY(Math.min(...zone.points.map((p) => p.y))) - 6, 12),
      }));
    const cl = region.countingLine;
    const line =
      cl && cl.start && cl.end
        ? { x1: toX(cl.start.x), y1: toY(cl.start.y), x2: toX(cl.end.x), y2: toY(cl.end.y) }
        : null;
    setRegionShapes({ zones, line });
  }, [phase, region, stageSize]);

  const statusTag = () => {
    if (phase === 'playing') {
      return (
        <Tag color="success" bordered={false}>
          {formatMessage('Vision.camera.liveOnline')}
        </Tag>
      );
    }
    if (phase === 'failed') {
      return (
        <Tag color="error" bordered={false}>
          {formatMessage('Vision.camera.previewFailed')}
        </Tag>
      );
    }
    return (
      <Tag color="processing" bordered={false}>
        {formatMessage('Vision.camera.previewConnecting')}
      </Tag>
    );
  };

  return (
    <div className={styles.liveView}>
      <div className={styles.header}>
        {name && <span className={styles.name}>{name}</span>}
        {statusTag()}
        {resolution && <span className={styles.resolution}>{resolution}</span>}
      </div>
      <div className={styles.stage} ref={stageRef}>
        <video
          ref={videoRef}
          className={styles.video}
          autoPlay
          playsInline
          muted
          onLoadedMetadata={handlePlaying}
          onPlaying={handlePlaying}
        />
        {(detectionBoxes.length > 0 || regionShapes.zones.length > 0 || regionShapes.line) && (
          <svg className={styles.detOverlay} width={stageSize.width} height={stageSize.height}>
            {regionShapes.zones.map((zone) => (
              <g key={zone.key}>
                <polygon className={styles.zoneShape} points={zone.points} />
                {zone.name && (
                  <text className={styles.detLabel} x={zone.labelX + 4} y={zone.labelY}>
                    {zone.name}
                  </text>
                )}
              </g>
            ))}
            {regionShapes.line && (
              <line
                className={styles.lineShape}
                x1={regionShapes.line.x1}
                y1={regionShapes.line.y1}
                x2={regionShapes.line.x2}
                y2={regionShapes.line.y2}
              />
            )}
            {detectionBoxes.map((box) => (
              <g key={box.key}>
                <rect className={styles.detBox} x={box.x} y={box.y} width={box.width} height={box.height} />
                {box.label && (
                  <text className={styles.detLabel} x={box.x + 4} y={Math.max(box.y - 6, 12)}>
                    {box.label}
                  </text>
                )}
              </g>
            ))}
          </svg>
        )}
        {phase !== 'playing' && (
          <div className={styles.overlay}>
            {phase === 'connecting' && <span>{formatMessage('Vision.camera.previewConnecting')}</span>}
            {phase === 'failed' && (
              <div className={styles.failBox}>
                <span>{failReason || formatMessage('Vision.camera.previewFailed')}</span>
                <Button
                  size="small"
                  onClick={() => {
                    setPhase('connecting');
                    setResolution('');
                    setFailReason('');
                    setRetryToken((token) => token + 1);
                  }}
                >
                  {formatMessage('Vision.camera.previewRetry')}
                </Button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

type TaskPreviewModalProps = {
  task: VisionTask | null;
  onClose: () => void;
};

// 任务实时预览：单相机直接渲染，多相机用 Tabs（destroyOnHidden 保证只连当前 tab）。
// 弹窗打开期间轮询任务最新识别结果,把各相机的检测框叠加到对应画面上。
const TaskPreviewModal = ({ task, onClose }: TaskPreviewModalProps) => {
  const formatMessage = useTranslate();
  const cameras = task?.cameras ?? [];
  const taskId = task?.id;
  const [detections, setDetections] = useState<Record<number, VisionLatestResult>>({});

  useEffect(() => {
    if (!taskId) {
      setDetections({});
      return undefined;
    }
    let disposed = false;
    const poll = async () => {
      try {
        const list = await getTaskLatestResult(taskId);
        if (disposed) return;
        const next: Record<number, VisionLatestResult> = {};
        list.forEach((item) => {
          next[item.cameraId] = item;
        });
        setDetections(next);
      } catch {
        // 静默:预览叠框是增强能力,拉取失败不打扰播放
      }
    };
    void poll();
    const timer = window.setInterval(poll, DETECTION_POLL_MS);
    return () => {
      disposed = true;
      window.clearInterval(timer);
    };
  }, [taskId]);

  const renderBody = () => {
    if (cameras.length === 0) {
      return (
        <Empty
          className={styles.empty}
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={formatMessage('Vision.task.previewNoCamera')}
        />
      );
    }
    if (cameras.length === 1) {
      const camera = cameras[0];
      return (
        <CameraLiveView
          key={camera.id}
          cameraId={camera.id}
          name={camera.name || camera.code}
          detection={detections[camera.id]}
          region={task?.regions?.[String(camera.id)]}
        />
      );
    }
    return (
      <Tabs
        className={styles.tabs}
        destroyOnHidden
        items={cameras.map((camera) => ({
          key: String(camera.id),
          label: camera.name || camera.code,
          children: (
            <CameraLiveView
              cameraId={camera.id}
              detection={detections[camera.id]}
              region={task?.regions?.[String(camera.id)]}
            />
          ),
        }))}
      />
    );
  };

  return (
    <ProModal
      open={Boolean(task)}
      title={formatMessage('Vision.task.previewTitle')}
      width={720}
      fullScreenable={false}
      onCancel={onClose}
      destroyOnHidden
      footer={null}
    >
      <div className={styles.body}>{renderBody()}</div>
    </ProModal>
  );
};

export default TaskPreviewModal;
