// 通用 3D 查看器引擎：加载 GLB、节点树、悬停/选中高亮、爆炸视图、实时绑定驱动、射线拾取、屏幕投影。
// 光照配置对齐源端 viewer-primitives/visual-setup/product-render-defaults.ts。
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
import { TransformControls } from 'three/addons/controls/TransformControls.js';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import type { MotionMappingNode } from '../api/instances';
import { buildNodeTree, findObjectByPath, type SceneNodeMeta } from './gltf-nodes';
import { applyPayloadToBindings, resetBindings, type ResolvedBinding } from './live-binding';

export type ViewerStatus = 'idle' | 'loading' | 'ready' | 'error';

export interface MarkerGizmoOptions {
  /** 模型根局部坐标（数据标签存储空间） */
  position: { x: number; y: number; z: number };
  /** 拖动过程回调（模型根局部坐标） */
  onPreview?: (position: { x: number; y: number; z: number }) => void;
  /** 拖动结束提交（仅在实际发生位移时触发） */
  onCommit?: (position: { x: number; y: number; z: number }) => void;
}

interface ExplodedEntry {
  object: THREE.Object3D;
  /** 初始世界坐标（爆炸按世界系定位后转回局部） */
  initialWorld: THREE.Vector3;
  /** 完全爆炸时相对初始位置的世界偏移 */
  finalOffset: THREE.Vector3;
  initialWorldMinY: number;
}

export interface ViewerApi {
  status: ViewerStatus;
  error: string;
  nodeTree: SceneNodeMeta[];
  modelHeight: number;
  highlight: (paths: string[]) => void;
  setExploded: (ratio: number) => void;
  resetView: () => void;
  /** 恢复绑定节点到初始姿态（实时绑定关闭时静态展示） */
  resetPose: () => void;
  setBindings: (mappings: MotionMappingNode[]) => void;
  applyPayload: (payload: Record<string, unknown>) => void;
  raycastAt: (clientX: number, clientY: number) => { path: string; local: THREE.Vector3 } | null;
  projectNodePoint: (path: string, local: { x: number; y: number; z: number }) => { x: number; y: number } | null;
  /** 模型根局部坐标 → 屏幕坐标（数据标签存储空间，对齐源端 data-tag-space） */
  projectModelPoint: (local: { x: number; y: number; z: number }) => { x: number; y: number } | null;
  /** 选中标点的平移 gizmo（对齐源端 marker TransformControls）；传 null 卸下 */
  setMarkerGizmo: (options: MarkerGizmoOptions | null) => void;
  /** 挂起/恢复渲染循环（页签隐藏或窗口后台时省资源）；恢复幂等，不会产生重复循环 */
  setSuspended: (suspended: boolean) => void;
}

const HIGHLIGHT_COLOR = new THREE.Color('#94c518');

export function useModelViewer(fileUrl: string | undefined) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [status, setStatus] = useState<ViewerStatus>('idle');
  const [error, setError] = useState('');
  const [nodeTree, setNodeTree] = useState<SceneNodeMeta[]>([]);
  const [modelHeight, setModelHeight] = useState(0);

  const stateRef = useRef<{
    renderer?: THREE.WebGLRenderer;
    scene?: THREE.Scene;
    camera?: THREE.PerspectiveCamera;
    model?: THREE.Object3D;
    exploded: ExplodedEntry[];
    explodeBaselineY?: number;
    resetCamera?: () => void;
    bindings: ResolvedBinding[];
    highlighted: Map<THREE.Mesh, THREE.Material | THREE.Material[]>;
    markerGizmo?: TransformControls;
    markerProxy?: THREE.Object3D;
    markerGizmoOptions?: MarkerGizmoOptions | null;
    suspended?: boolean;
    resumeLoop?: () => void;
  }>({ exploded: [], bindings: [], highlighted: new Map() });

  useEffect(() => {
    if (!fileUrl || !containerRef.current) return;
    const container = containerRef.current;
    const state = stateRef.current;
    const isDark = document.documentElement.classList.contains('dark');
    setStatus('loading');
    setError('');

    const renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.toneMapping = THREE.ACESFilmicToneMapping;
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFSoftShadowMap;
    container.appendChild(renderer.domElement);

    const scene = new THREE.Scene();
    scene.background = new THREE.Color(isDark ? '#1f1f1f' : '#f0f1f3');

    const camera = new THREE.PerspectiveCamera(45, container.clientWidth / container.clientHeight, 0.01, 2000);

    // 源端 product-render-defaults：key/fill/back 三灯 + 弱环境光 + 软阴影
    const ambient = new THREE.AmbientLight('#ffffff', 0.6);
    const keyLight = new THREE.DirectionalLight('#ffffff', 4);
    keyLight.position.set(-4, 6, 5);
    keyLight.castShadow = true;
    keyLight.shadow.mapSize.set(2048, 2048);
    keyLight.shadow.bias = -0.0005;
    const fillLight = new THREE.DirectionalLight('#ffffff', 0.8);
    fillLight.position.set(5, 3.5, 4);
    const backLight = new THREE.DirectionalLight('#ffffff', 0.2);
    backLight.position.set(0, 5, -6);
    scene.add(ambient, keyLight, fillLight, backLight);

    const controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = true;

    // 选中标点的平移 gizmo（对齐源端 marker TransformControls：proxy 对象承接拖动，模型根空间存取）
    const markerProxy = new THREE.Object3D();
    scene.add(markerProxy);
    const markerGizmo = new TransformControls(camera, renderer.domElement);
    markerGizmo.setMode('translate');
    markerGizmo.setSpace('local');
    markerGizmo.setSize(0.95);
    const proxyModelLocal = () => {
      if (!state.model) return null;
      const world = markerProxy.getWorldPosition(new THREE.Vector3());
      const local = state.model.worldToLocal(world);
      return { x: local.x, y: local.y, z: local.z };
    };
    let markerDragMoved = false;
    markerGizmo.addEventListener('dragging-changed', (event) => {
      controls.enabled = !event.value;
      if (!event.value && markerDragMoved) {
        markerDragMoved = false;
        const local = proxyModelLocal();
        if (local) state.markerGizmoOptions?.onCommit?.(local);
      }
    });
    markerGizmo.addEventListener('objectChange', () => {
      if (!(markerGizmo as unknown as { dragging: boolean }).dragging) return;
      markerDragMoved = true;
      const local = proxyModelLocal();
      if (local) state.markerGizmoOptions?.onPreview?.(local);
    });
    scene.add(markerGizmo.getHelper());
    state.markerGizmo = markerGizmo;
    state.markerProxy = markerProxy;

    let disposed = false;
    let frameId = 0;
    let ground: THREE.Mesh | null = null;

    new GLTFLoader().loadAsync(fileUrl).then(
      (gltf) => {
        if (disposed) return;
        const model = gltf.scene;
        model.traverse((object) => {
          if ((object as THREE.Mesh).isMesh) {
            object.castShadow = true;
            object.receiveShadow = true;
          }
        });
        scene.add(model);

        const bounds = new THREE.Box3().setFromObject(model);
        const size = bounds.getSize(new THREE.Vector3());
        const center = bounds.getCenter(new THREE.Vector3());
        const maxDim = Math.max(size.x, size.y, size.z) || 1;
        const distance = (maxDim / (2 * Math.tan(THREE.MathUtils.degToRad(camera.fov / 2)))) * 1.4;
        camera.position.set(center.x + distance * 0.9, center.y + distance * 0.5, center.z + distance * 0.9);
        camera.near = Math.max(distance / 100, 0.01);
        camera.far = Math.max(distance * 10, 100);
        camera.updateProjectionMatrix();
        controls.target.copy(center);
        controls.update();

        // 地面（承接软阴影）
        const groundMaterial = new THREE.ShadowMaterial({ opacity: isDark ? 0.4 : 0.25 });
        ground = new THREE.Mesh(new THREE.PlaneGeometry(maxDim * 10, maxDim * 10), groundMaterial);
        ground.rotation.x = -Math.PI / 2;
        ground.position.y = bounds.min.y;
        ground.receiveShadow = true;
        scene.add(ground);
        const grid = isDark
          ? new THREE.GridHelper(maxDim * 6, 24, 0x525252, 0x393939)
          : new THREE.GridHelper(maxDim * 6, 24, 0xc6c6c6, 0xe0e0e0);
        grid.position.y = bounds.min.y;
        scene.add(grid);

        // 爆炸视图（对齐源端 TargetCubeScalingStrategy）：按 mesh 收集单元，
        // 终态 = 包围盒中心按 1.75×最长边的目标立方体逐轴等比外推；带防穿地抬升
        const exploded: ExplodedEntry[] = [];
        const targetCubeSide = maxDim * 1.75;
        const axisScale = new THREE.Vector3(
          size.x > 1e-6 ? targetCubeSide / size.x : 1,
          size.y > 1e-6 ? targetCubeSide / size.y : 1,
          size.z > 1e-6 ? targetCubeSide / size.z : 1
        );
        model.updateMatrixWorld(true);
        model.traverse((object) => {
          if (!(object as THREE.Mesh).isMesh) return;
          const meshBounds = new THREE.Box3().setFromObject(object);
          if (meshBounds.isEmpty()) return;
          const meshCenter = meshBounds.getCenter(new THREE.Vector3());
          const finalOffset = meshCenter.clone().sub(center).multiply(axisScale).add(center).sub(meshCenter);
          exploded.push({
            object,
            initialWorld: object.getWorldPosition(new THREE.Vector3()),
            finalOffset,
            initialWorldMinY: meshBounds.min.y,
          });
        });
        const explodeBaselineY = center.y - size.y * 0.5;

        state.renderer = renderer;
        state.scene = scene;
        state.camera = camera;
        state.model = model;
        state.exploded = exploded;
        state.explodeBaselineY = explodeBaselineY;
        state.bindings = [];
        state.highlighted = new Map();
        // 重置视角 = 回到初始取景
        const homePosition = camera.position.clone();
        const homeTarget = center.clone();
        state.resetCamera = () => {
          camera.position.copy(homePosition);
          controls.target.copy(homeTarget);
          controls.update();
        };

        setNodeTree(buildNodeTree(model));
        setModelHeight(size.y);
        setStatus('ready');
      },
      (e) => {
        if (disposed) return;
        setError((e as Error).message || String(e));
        setStatus('error');
      }
    );

    // 渲染循环：挂起时不再排下一帧（loopRunning 防止恢复时产生重复循环）
    let loopRunning = false;
    const animate = () => {
      if (state.suspended) {
        loopRunning = false;
        return;
      }
      frameId = requestAnimationFrame(animate);
      controls.update();
      renderer.render(scene, camera);
    };
    const startLoop = () => {
      if (loopRunning || disposed) return;
      loopRunning = true;
      animate();
    };
    state.resumeLoop = startLoop;
    startLoop();

    const handleResize = () => {
      if (!container.clientWidth || !container.clientHeight) return;
      camera.aspect = container.clientWidth / container.clientHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(container.clientWidth, container.clientHeight);
    };
    const observer = new ResizeObserver(handleResize);
    observer.observe(container);

    return () => {
      disposed = true;
      cancelAnimationFrame(frameId);
      observer.disconnect();
      markerGizmo.detach();
      markerGizmo.dispose();
      controls.dispose();
      renderer.dispose();
      if (renderer.domElement.parentNode === container) container.removeChild(renderer.domElement);
      stateRef.current = { exploded: [], bindings: [], highlighted: new Map() };
      setNodeTree([]);
      setStatus('idle');
    };
  }, [fileUrl]);

  const highlight = useCallback((paths: string[]) => {
    const state = stateRef.current;
    if (!state.model) return;
    state.highlighted.forEach((original, mesh) => {
      mesh.material = original;
    });
    state.highlighted.clear();
    const wanted = new Set(paths);
    if (wanted.size === 0) return;
    for (const path of wanted) {
      const object = findObjectByPath(state.model, path);
      if (!object) continue;
      object.traverse((child) => {
        const mesh = child as THREE.Mesh;
        if (!mesh.isMesh || state.highlighted.has(mesh)) return;
        state.highlighted.set(mesh, mesh.material);
        const cloned = (Array.isArray(mesh.material) ? mesh.material[0] : mesh.material).clone();
        if ('emissive' in cloned) {
          (cloned as THREE.MeshStandardMaterial).emissive = HIGHLIGHT_COLOR;
          (cloned as THREE.MeshStandardMaterial).emissiveIntensity = 0.55;
        }
        mesh.material = cloned;
      });
    }
  }, []);

  const setExploded = useCallback((ratio: number) => {
    const state = stateRef.current;
    const t = Math.min(Math.max(ratio, 0), 1);
    // 防穿地：爆炸后最低点低于地面基线时整体抬升（对齐源端 ensureLift）
    let minY = Infinity;
    state.exploded.forEach(({ initialWorldMinY, finalOffset }) => {
      const candidate = initialWorldMinY + finalOffset.y * t;
      if (candidate < minY) minY = candidate;
    });
    const lift = Number.isFinite(minY) ? Math.max(0, (state.explodeBaselineY ?? 0) - minY) : 0;
    const targetWorld = new THREE.Vector3();
    state.exploded.forEach(({ object, initialWorld, finalOffset }) => {
      targetWorld
        .copy(initialWorld)
        .addScaledVector(finalOffset, t)
        .add(new THREE.Vector3(0, lift, 0));
      object.parent?.updateWorldMatrix(true, false);
      object.position.copy(object.parent ? object.parent.worldToLocal(targetWorld.clone()) : targetWorld);
    });
  }, []);

  const setBindings = useCallback((mappings: MotionMappingNode[]) => {
    const state = stateRef.current;
    resetBindings(state.bindings);
    state.bindings = [];
    if (!state.model) return;
    for (const mapping of mappings) {
      if (!mapping.value) continue;
      const node = findObjectByPath(state.model, mapping.path);
      if (!node) continue;
      state.bindings.push({
        binding: mapping,
        node,
        originalPosition: node.position.clone(),
        originalRotation: node.rotation.clone(),
      });
    }
  }, []);

  const applyPayload = useCallback((payload: Record<string, unknown>) => {
    applyPayloadToBindings(payload, stateRef.current.bindings);
  }, []);

  const raycastAt = useCallback((clientX: number, clientY: number) => {
    const state = stateRef.current;
    const container = containerRef.current;
    if (!state.model || !state.camera || !container) return null;
    const rect = container.getBoundingClientRect();
    const pointer = new THREE.Vector2(
      ((clientX - rect.left) / rect.width) * 2 - 1,
      -((clientY - rect.top) / rect.height) * 2 + 1
    );
    const raycaster = new THREE.Raycaster();
    raycaster.setFromCamera(pointer, state.camera);
    const hits = raycaster.intersectObject(state.model, true);
    const hit = hits[0];
    if (!hit) return null;
    // 找到有名字的最近祖先作为标签挂载节点
    let node: THREE.Object3D | null = hit.object;
    while (node && !node.name && node.parent) node = node.parent;
    if (!node) return null;
    const segments: string[] = [];
    let cursor: THREE.Object3D | null = node;
    while (cursor) {
      segments.unshift(cursor.name || cursor.type || 'Node');
      if (cursor === state.model) break;
      cursor = cursor.parent;
    }
    // 坐标存模型根局部空间（对齐源端：worldToLocal(modelRoot)），path 仅作元数据
    const local = state.model.worldToLocal(hit.point.clone());
    return { path: segments.join('/'), local };
  }, []);

  const projectWorldPoint = useCallback((world: THREE.Vector3) => {
    const state = stateRef.current;
    const container = containerRef.current;
    if (!state.camera || !container) return null;
    const projected = world.project(state.camera);
    if (projected.z > 1) return null;
    return {
      x: ((projected.x + 1) / 2) * container.clientWidth,
      y: ((1 - projected.y) / 2) * container.clientHeight,
    };
  }, []);

  const projectNodePoint = useCallback(
    (path: string, local: { x: number; y: number; z: number }) => {
      const state = stateRef.current;
      if (!state.model) return null;
      const node = findObjectByPath(state.model, path);
      if (!node) return null;
      return projectWorldPoint(node.localToWorld(new THREE.Vector3(local.x ?? 0, local.y ?? 0, local.z ?? 0)));
    },
    [projectWorldPoint]
  );

  const projectModelPoint = useCallback(
    (local: { x: number; y: number; z: number }) => {
      const state = stateRef.current;
      if (!state.model) return null;
      return projectWorldPoint(state.model.localToWorld(new THREE.Vector3(local.x ?? 0, local.y ?? 0, local.z ?? 0)));
    },
    [projectWorldPoint]
  );

  const setSuspended = useCallback((suspended: boolean) => {
    const state = stateRef.current;
    if (state.suspended === suspended) return;
    state.suspended = suspended;
    if (!suspended) state.resumeLoop?.();
  }, []);

  const setMarkerGizmo = useCallback((options: MarkerGizmoOptions | null) => {
    const state = stateRef.current;
    state.markerGizmoOptions = options;
    const { markerGizmo, markerProxy, model } = state;
    if (!markerGizmo || !markerProxy) return;
    if (!options || !model) {
      markerGizmo.detach();
      return;
    }
    // 拖动中不回置 proxy 位置（避免预览回写抖动）
    if (!(markerGizmo as unknown as { dragging: boolean }).dragging) {
      const world = model.localToWorld(new THREE.Vector3(options.position.x, options.position.y, options.position.z));
      markerProxy.position.copy(world);
    }
    markerGizmo.attach(markerProxy);
  }, []);

  const resetView = useCallback(() => stateRef.current.resetCamera?.(), []);
  const resetPose = useCallback(() => resetBindings(stateRef.current.bindings), []);

  // api 必须引用稳定：页面高频重渲染（MQTT 实时消息）时，依赖 viewer 的 effect 不应反复重跑
  const api: ViewerApi = useMemo(
    () => ({
      status,
      error,
      nodeTree,
      modelHeight,
      highlight,
      setExploded,
      resetView,
      resetPose,
      setBindings,
      applyPayload,
      raycastAt,
      projectNodePoint,
      projectModelPoint,
      setMarkerGizmo,
      setSuspended,
    }),
    [
      status,
      error,
      nodeTree,
      modelHeight,
      highlight,
      setExploded,
      resetView,
      resetPose,
      setBindings,
      applyPayload,
      raycastAt,
      projectNodePoint,
      projectModelPoint,
      setMarkerGizmo,
      setSuspended,
    ]
  );

  return { containerRef, viewer: api };
}
