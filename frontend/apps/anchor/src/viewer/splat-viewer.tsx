// SPLAT（3D Gaussian Splatting）查看器：基于 @mkkellogg/gaussian-splats-3d 独立 Viewer。
// 对齐源端 utils/3dgs/gaussian-splat.ts 的呈现约定：
//   .splat 原始数据是 y 向下（COLMAP 系），渲染时绕 X 轴转 180°（y/z 取反）；
//   解析包围盒后用 [-center.x, -min.y, -center.z] 偏移让内容落地居中，相机按包围盒取景。
// 数据标签存储空间对齐源端（storageRoot=splatPickGroup）：渲染世界坐标 - contentOffset，
//   path/nodeID 固定 '__3dgs__'；拾取用库内置 Raycaster.intersectSplatMesh。
// SPLAT 实例只支持静态数据标点（无节点树/运动映射），场景摆放暂不支持。
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import * as THREE from 'three';
import { t } from '../i18n';

const SPLAT_VERTEX_SIZE = 32;

export const SPLAT_TAG_PATH = '__3dgs__';

interface SplatBounds {
  min: THREE.Vector3;
  max: THREE.Vector3;
  center: THREE.Vector3;
  size: THREE.Vector3;
}

// 与源端 computeSplatBoundingBox 一致：按渲染空间 (x, -y, -z) 统计包围盒
function computeSplatBounds(buffer: ArrayBuffer): SplatBounds | null {
  const vertexCount = Math.floor(buffer.byteLength / SPLAT_VERTEX_SIZE);
  if (vertexCount === 0) return null;
  const view = new DataView(buffer);
  const min = new THREE.Vector3(Infinity, Infinity, Infinity);
  const max = new THREE.Vector3(-Infinity, -Infinity, -Infinity);
  for (let index = 0; index < vertexCount; index += 1) {
    const offset = index * SPLAT_VERTEX_SIZE;
    const sx = view.getFloat32(offset, true);
    const sy = view.getFloat32(offset + 4, true);
    const sz = view.getFloat32(offset + 8, true);
    if (!Number.isFinite(sx) || !Number.isFinite(sy) || !Number.isFinite(sz)) continue;
    const x = sx;
    const y = -sy;
    const z = -sz;
    min.x = Math.min(min.x, x);
    min.y = Math.min(min.y, y);
    min.z = Math.min(min.z, z);
    max.x = Math.max(max.x, x);
    max.y = Math.max(max.y, y);
    max.z = Math.max(max.z, z);
  }
  if (min.x === Infinity || max.x === -Infinity) return null;
  const center = min.clone().add(max).multiplyScalar(0.5);
  const size = max.clone().sub(min);
  return { min, max, center, size };
}

export interface SplatViewerApi {
  status: 'idle' | 'loading' | 'ready' | 'error';
  error: string;
  /** 内容包围盒高度（真实高度展示/默认实例高度） */
  modelHeight: number;
  /** 点击拾取（返回存储空间坐标 + 固定 path '__3dgs__'），未命中返回 null */
  raycastAt: (clientX: number, clientY: number) => { path: string; local: THREE.Vector3 } | null;
  /** 存储空间坐标 → 屏幕坐标（供 TagOverlay 投影） */
  projectModelPoint: (local: { x: number; y: number; z: number }) => { x: number; y: number } | null;
  /** SPLAT 无节点树，恒为 null（兼容 TagOverlay 的运动标签投影签名） */
  projectNodePoint: (path: string, local: { x: number; y: number; z: number }) => { x: number; y: number } | null;
  /** 挂起/恢复渲染（页签隐藏或窗口后台时省资源）；恢复幂等 */
  setSuspended: (suspended: boolean) => void;
}

interface SplatRuntime {
  camera?: THREE.PerspectiveCamera;
  raycaster?: {
    setFromCameraAndScreenPosition: (
      camera: THREE.Camera,
      position: { x: number; y: number },
      dimensions: { x: number; y: number }
    ) => void;
    intersectSplatMesh: (mesh: unknown, outHits: { origin: THREE.Vector3 }[]) => void;
  };
  splatMesh?: unknown;
  contentOffset: THREE.Vector3;
  viewer?: { start?: () => void; stop?: () => void };
  suspended?: boolean;
}

export function useSplatViewer(fileUrl: string | undefined, onHeight?: (height: number) => void) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [status, setStatus] = useState<'idle' | 'loading' | 'ready' | 'error'>('idle');
  const [error, setError] = useState('');
  const [modelHeight, setModelHeight] = useState(0);
  const runtimeRef = useRef<SplatRuntime>({ contentOffset: new THREE.Vector3() });
  const onHeightRef = useRef(onHeight);
  useEffect(() => {
    onHeightRef.current = onHeight;
  }, [onHeight]);

  useEffect(() => {
    const container = containerRef.current;
    if (!fileUrl || !container) return;
    let disposed = false;
    let viewer: { dispose?: () => void; stop?: () => void } | null = null;
    let objectUrl = '';
    setStatus('loading');
    setError('');

    (async () => {
      try {
        const [GaussianSplats3D, response] = await Promise.all([
          import('@mkkellogg/gaussian-splats-3d'),
          fetch(fileUrl, { credentials: 'include' }),
        ]);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const buffer = await response.arrayBuffer();
        if (disposed) return;

        const bounds = computeSplatBounds(buffer);
        // 落地居中偏移（源端 createSplatContentOffset）
        const contentOffset = bounds
          ? new THREE.Vector3(-bounds.center.x, -bounds.min.y, -bounds.center.z)
          : new THREE.Vector3();
        const size = bounds?.size ?? new THREE.Vector3(2, 2, 2);
        const maxDim = Math.max(size.x, size.y, size.z) || 1;
        const lookAt = new THREE.Vector3(0, size.y / 2, 0);
        // 取景与 GLB 查看器同款：按最长边+fov 推相机距离
        const distance = (maxDim / (2 * Math.tan(THREE.MathUtils.degToRad(50 / 2)))) * 1.4;
        const cameraPosition = new THREE.Vector3(distance * 0.9, lookAt.y + distance * 0.5, distance * 0.9);

        // 主题背景 + 地面网格（与 GLB 查看器一致），通过外部 scene 注入
        const isDark = document.documentElement.classList.contains('dark');
        const scene = new THREE.Scene();
        scene.background = new THREE.Color(isDark ? '#1f1f1f' : '#f0f1f3');
        const grid = isDark
          ? new THREE.GridHelper(maxDim * 6, 24, 0x525252, 0x393939)
          : new THREE.GridHelper(maxDim * 6, 24, 0xc6c6c6, 0xe0e0e0);
        scene.add(grid);

        const instance = new GaussianSplats3D.Viewer({
          rootElement: container,
          threeScene: scene,
          selfDrivenMode: true,
          useBuiltInControls: true,
          sharedMemoryForWorkers: false, // 避免 COOP/COEP 要求
          cameraUp: [0, 1, 0],
          initialCameraPosition: cameraPosition.toArray(),
          initialCameraLookAt: lookAt.toArray(),
        });
        viewer = instance;
        objectUrl = URL.createObjectURL(new Blob([buffer]));
        await instance.addSplatScene(objectUrl, {
          format: GaussianSplats3D.SceneFormat.Splat,
          showLoadingUI: false,
          // 源端渲染约定：绕 X 轴 180°（y/z 取反）再平移落地
          rotation: [1, 0, 0, 0],
          position: contentOffset.toArray(),
        });
        if (disposed) return;
        instance.start();

        const runtime = runtimeRef.current;
        runtime.contentOffset = contentOffset;
        runtime.camera = (instance as unknown as { camera?: THREE.PerspectiveCamera }).camera;
        runtime.raycaster = (instance as unknown as { raycaster?: SplatRuntime['raycaster'] }).raycaster;
        runtime.splatMesh = (instance as unknown as { splatMesh?: unknown }).splatMesh;
        runtime.viewer = instance;
        // 加载期间若已被挂起（页签在后台完成加载），立即停掉自驱动循环
        if (runtime.suspended) instance.stop?.();

        setStatus('ready');
        if (bounds) {
          setModelHeight(size.y);
          onHeightRef.current?.(size.y);
        }
      } catch (e) {
        if (!disposed) {
          setError((e as Error).message || String(e));
          setStatus('error');
        }
      }
    })();

    return () => {
      disposed = true;
      if (objectUrl) URL.revokeObjectURL(objectUrl);
      runtimeRef.current = { contentOffset: new THREE.Vector3() };
      setStatus('idle');
      setModelHeight(0);
      try {
        viewer?.stop?.();
        viewer?.dispose?.();
      } catch {
        /* dispose 竞态忽略 */
      }
    };
  }, [fileUrl]);

  const raycastAt = useCallback((clientX: number, clientY: number) => {
    const container = containerRef.current;
    const { camera, raycaster, splatMesh, contentOffset } = runtimeRef.current;
    if (!container || !camera || !raycaster || !splatMesh) return null;
    const rect = container.getBoundingClientRect();
    if (rect.width <= 0 || rect.height <= 0) return null;
    raycaster.setFromCameraAndScreenPosition(
      camera,
      { x: clientX - rect.left, y: clientY - rect.top },
      { x: rect.width, y: rect.height }
    );
    const hits: { origin: THREE.Vector3 }[] = [];
    raycaster.intersectSplatMesh(splatMesh, hits);
    const hit = hits[0];
    if (!hit) return null;
    // 存储空间 = 渲染世界坐标 - contentOffset（对齐源端 splatPickGroup storageRoot）
    const local = hit.origin.clone().sub(contentOffset);
    return { path: SPLAT_TAG_PATH, local };
  }, []);

  const projectModelPoint = useCallback((local: { x: number; y: number; z: number }) => {
    const container = containerRef.current;
    const { camera, contentOffset } = runtimeRef.current;
    if (!container || !camera) return null;
    const world = new THREE.Vector3(local.x ?? 0, local.y ?? 0, local.z ?? 0).add(contentOffset);
    const projected = world.project(camera);
    if (projected.z > 1) return null;
    return {
      x: ((projected.x + 1) / 2) * container.clientWidth,
      y: ((1 - projected.y) / 2) * container.clientHeight,
    };
  }, []);

  const projectNodePoint = useCallback(() => null, []);

  const setSuspended = useCallback((suspended: boolean) => {
    const runtime = runtimeRef.current;
    if (runtime.suspended === suspended) return;
    runtime.suspended = suspended;
    if (suspended) runtime.viewer?.stop?.();
    else runtime.viewer?.start?.();
  }, []);

  const viewer: SplatViewerApi = useMemo(
    () => ({ status, error, modelHeight, raycastAt, projectModelPoint, projectNodePoint, setSuspended }),
    [status, error, modelHeight, raycastAt, projectModelPoint, projectNodePoint, setSuspended]
  );

  return { containerRef, viewer };
}

export function SplatViewer({
  fileUrl,
  onHeight,
  suspended = false,
}: {
  fileUrl: string;
  onHeight?: (height: number) => void;
  /** 页签隐藏/窗口后台时暂停渲染 */
  suspended?: boolean;
}) {
  const { containerRef, viewer } = useSplatViewer(fileUrl, onHeight);
  useEffect(() => {
    viewer.setSuspended(suspended);
  }, [viewer, suspended]);

  return (
    <div className="relative h-full w-full">
      <div ref={containerRef} className="h-full w-full" style={{ background: 'var(--ui-bg-color)' }} />
      {viewer.status !== 'ready' ? (
        <div
          className="absolute inset-0 flex items-center justify-center text-sm"
          style={{ background: 'var(--ui-bg-color)', color: 'var(--ui-description-card-color)' }}
        >
          {viewer.status === 'error' ? `${t('model.viewer.error')}: ${viewer.error}` : t('model.viewer.loading')}
        </div>
      ) : null}
    </div>
  );
}
