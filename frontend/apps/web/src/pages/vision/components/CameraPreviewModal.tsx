import { useEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { Button, Tag } from 'antd';
import { previewVideoCamera, stopVideoCameraPreview, type VideoCamera } from '@/apis/core-api/video';
import { Expand, Shrink } from '@/components/lucide-icon/carbon';
import { formatTimestamp } from '@/utils/format';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import styles from './CameraPreviewModal.module.scss';

type PreviewPhase = 'connecting' | 'playing' | 'failed';

type CameraPreviewModalProps = {
  open: boolean;
  camera: VideoCamera | null;
  onClose: () => void;
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

const CameraPreviewModal = ({ open, camera, onClose }: CameraPreviewModalProps) => {
  const formatMessage = useTranslate();
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const stageRef = useRef<HTMLDivElement | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const [expanded, setExpanded] = useState(false);
  const [nativeFullscreen, setNativeFullscreen] = useState(false);
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const [phase, setPhase] = useState<PreviewPhase>('connecting');
  const [resolution, setResolution] = useState('');
  const [failReason, setFailReason] = useState('');
  const [retryToken, setRetryToken] = useState(0);

  useEffect(() => {
    if (!open || !camera) return;
    let disposed = false;
    const videoEl = videoRef.current;

    const pc = new RTCPeerConnection();
    pcRef.current = pc;
    pc.addTransceiver('video', { direction: 'recvonly' });
    pc.addTransceiver('audio', { direction: 'recvonly' });
    pc.ontrack = (event) => {
      if (!event.streams[0]) return;
      streamRef.current = event.streams[0];
      if (videoRef.current) {
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
        const answer = await previewVideoCamera(camera.id, pc.localDescription?.sdp || offer.sdp || '');
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
      void stopVideoCameraPreview(camera.id).catch(() => undefined);
    };
  }, [open, camera, formatMessage, retryToken]);

  const handlePlaying = () => {
    setPhase('playing');
    const video = videoRef.current;
    if (video && video.videoWidth && video.videoHeight) {
      setResolution(`${video.videoWidth}x${video.videoHeight}`);
    }
  };

  // 全屏只放大画面本身。优先用浏览器原生全屏;产品支持 iframe 嵌入,宿主没给
  // allow="fullscreen" 时 requestFullscreen 会被拒,退回铺满视口的放大态。
  const toggleFullscreen = () => {
    if (document.fullscreenElement) {
      void document.exitFullscreen();
      return;
    }
    if (expanded) {
      setExpanded(false);
      return;
    }
    const request = stageRef.current?.requestFullscreen?.();
    if (request) {
      request.catch(() => setExpanded(true));
    } else {
      setExpanded(true);
    }
  };

  // 原生全屏是浏览器在管的,按 Esc 或系统手势退出都不经过上面的 toggle,
  // 所以图标状态得跟着 fullscreenchange 走,不能只看 expanded。
  useEffect(() => {
    const sync = () => setNativeFullscreen(document.fullscreenElement === stageRef.current);
    document.addEventListener('fullscreenchange', sync);
    return () => document.removeEventListener('fullscreenchange', sync);
  }, []);

  const enlarged = expanded || nativeFullscreen;

  // 放大态用 portal 挂到 body:弹窗被 react-draggable 的 transform 包住,
  // 留在原位的 position: fixed 只会铺满弹窗。换父节点会让 video 重新挂载,
  // 所以这里把已协商好的流重新接回去,不必重建 PeerConnection。
  useEffect(() => {
    const video = videoRef.current;
    if (video && streamRef.current && video.srcObject !== streamRef.current) {
      video.srcObject = streamRef.current;
      void video.play().catch(() => undefined);
    }
  }, [expanded]);

  // Esc 只退出放大。原生全屏由浏览器自己退,这里只负责拦掉冒泡,
  // 否则 antd Modal 会顺手把整个弹窗也关了。
  useEffect(() => {
    if (!enlarged) return undefined;
    const onKey = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return;
      event.stopPropagation();
      if (expanded) setExpanded(false);
    };
    window.addEventListener('keydown', onKey, true);
    return () => window.removeEventListener('keydown', onKey, true);
  }, [enlarged, expanded]);

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

  const stage = (
    <div className={`${styles.stage} ${expanded ? styles.stageExpanded : ''}`} ref={stageRef}>
      <video
        ref={videoRef}
        className={styles.video}
        autoPlay
        playsInline
        muted
        onLoadedMetadata={handlePlaying}
        onPlaying={handlePlaying}
      />
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
      <button
        type="button"
        className={styles.fullscreenBtn}
        aria-label={formatMessage(enlarged ? 'common.exitFullScreen' : 'Vision.camera.fullscreen')}
        onClick={toggleFullscreen}
      >
        {enlarged ? <Shrink size={12} /> : <Expand size={12} />}
      </button>
    </div>
  );

  return (
    <ProModal
      open={open}
      title={formatMessage('Vision.camera.previewTitle')}
      width={760}
      fullScreenable={false}
      onCancel={onClose}
      destroyOnHidden
      footer={
        <div className={styles.previewFooter}>
          <Button type="primary" onClick={onClose}>
            {formatMessage('Vision.camera.closePreview')}
          </Button>
        </div>
      }
    >
      {expanded ? createPortal(stage, document.body) : stage}
      <div className={styles.previewMeta}>
        <div className={styles.previewField}>
          <span className={styles.previewLabel}>{formatMessage('common.name')}</span>
          <span className={styles.previewValue}>{camera?.name || '-'}</span>
        </div>
        <div className={styles.previewField}>
          <span className={styles.previewLabel}>{formatMessage('Vision.camera.latestUpdate')}</span>
          <span className={styles.previewValue}>
            {camera?.lastTestedAt ? formatTimestamp(camera.lastTestedAt) : '-'}
          </span>
        </div>
        <div className={styles.previewField}>
          <span className={styles.previewLabel}>{formatMessage('Vision.camera.liveStatus')}</span>
          <span className={styles.previewValue}>
            {statusTag()}
            {resolution && <span className={styles.resolution}>{resolution}</span>}
          </span>
        </div>
        <div className={styles.previewField}>
          <span className={styles.previewLabel}>{formatMessage('Vision.camera.protocol')}</span>
          <span className={styles.previewValue}>{camera?.rtpType?.toUpperCase() || '-'}</span>
        </div>
      </div>
    </ProModal>
  );
};

export default CameraPreviewModal;
