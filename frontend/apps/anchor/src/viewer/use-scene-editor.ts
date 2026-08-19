// 场景编辑器引擎：多实例摆放、TransformControls 变换 gizmo、点选、SceneConfigV4 应用、画布快照。
// 功能对齐源端 scene-editor-shell / scene-r3f-canvas-runtime（此处用纯 three.js 实现）。
import { useCallback, useEffect, useRef, useState } from 'react';
import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
import { TransformControls } from 'three/addons/controls/TransformControls.js';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { RoomEnvironment } from 'three/addons/environments/RoomEnvironment.js';
import { Reflector } from 'three/addons/objects/Reflector.js';
import type { MotionMappingNode } from '../api/instances';
import type { SceneConfigV4, ScenePlacement } from '../api/scenes';
import { findObjectByPath } from './gltf-nodes';
import { applyPayloadToBindings, type ResolvedBinding } from './live-binding';

export type GizmoMode = 'translate' | 'rotate' | 'scale';

export type ScenePivotMode = 'center' | 'bottomCenter';

export interface EditorItem {
  key: string;
  instanceId: number;
  root: THREE.Group;
  /** 模型容器（pivot 切换时移动它的局部偏移，root 原点即 gizmo 锚点） */
  inner: THREE.Object3D;
  pivot: ScenePivotMode;
  baseScale: number;
  bindings: ResolvedBinding[];
}

// 背景预设色值对齐源端 viewport-panel BACKGROUND_OPTIONS
const BG_PRESETS: Record<string, string | null> = {
  white: '#fcfcfc',
  lightGray: '#b2c1d7',
  black: '#1a1a1a',
  gradient: null,
};

// 光照预设（对齐源端 RENDER_EFFECTS_PRESET_DEFAULTS）
export const LIGHT_PRESETS: Record<string, { key: number; fill: number; back: number; ambient: number }> = {
  balanced: { key: 4, fill: 0.8, back: 0.2, ambient: 0.1 },
  soft: { key: 6.2, fill: 0.1, back: 0.06, ambient: 0.15 },
  contrast: { key: 9.5, fill: 0.08, back: 0.12, ambient: 0.09 },
};

function gradientTexture(): THREE.Texture {
  const canvas = document.createElement('canvas');
  canvas.width = 2;
  canvas.height = 256;
  const context = canvas.getContext('2d')!;
  // 对齐源端 gradient 预设：to top #1a1a1a → #404040 → #808080（canvas 从上往下画，故反向）
  const gradient = context.createLinearGradient(0, 0, 0, 256);
  gradient.addColorStop(0, '#808080');
  gradient.addColorStop(0.5, '#404040');
  gradient.addColorStop(1, '#1a1a1a');
  context.fillStyle = gradient;
  context.fillRect(0, 0, 2, 256);
  const texture = new THREE.CanvasTexture(canvas);
  texture.colorSpace = THREE.SRGBColorSpace;
  return texture;
}

// pivot 切换（对齐源端 scene-transform-runtime-math）：root 原点即 gizmo 锚点，
// 把 inner 平移使原点落到几何中心 / 底部中心；compensate=true 时补偿 root 位置保持世界位置不变
function applyPivot(item: EditorItem, mode: ScenePivotMode, compensate: boolean) {
  if (item.pivot === mode) return;
  const { root, inner } = item;
  const prevPos = root.position.clone();
  const prevQuat = root.quaternion.clone();
  const prevScale = root.scale.clone();
  root.position.set(0, 0, 0);
  root.quaternion.identity();
  root.scale.set(1, 1, 1);
  root.updateWorldMatrix(true, true);
  const box = new THREE.Box3().setFromObject(inner);
  const restore = () => {
    root.position.copy(prevPos);
    root.quaternion.copy(prevQuat);
    root.scale.copy(prevScale);
    root.updateWorldMatrix(true, true);
  };
  if (box.isEmpty()) {
    restore();
    return;
  }
  const center = box.getCenter(new THREE.Vector3());
  const target = mode === 'center' ? center : new THREE.Vector3(center.x, box.min.y, center.z);
  inner.position.sub(target);
  restore();
  if (compensate) {
    root.position.add(target.clone().multiply(prevScale).applyQuaternion(prevQuat));
    root.updateWorldMatrix(true, true);
  }
  item.pivot = mode;
}

// 地面反射（对齐源端 reflection-config 语义：clarity → 采样分辨率，depth → 不透明度）
function reflectionResolution(clarity: number): number {
  if (clarity >= 75) return 1024;
  if (clarity >= 40) return 512;
  return 256;
}

export function useSceneEditor(options?: { interactive?: boolean }) {
  // 只读模式（扫码分享 viewer）：不创建变换 gizmo、不挂点选，仅浏览 + 实时驱动
  const interactive = options?.interactive !== false;
  const containerRef = useRef<HTMLDivElement>(null);
  const [ready, setReady] = useState(false);
  const [selectedKey, setSelectedKey] = useState<string | null>(null);
  const [transformTick, setTransformTick] = useState(0);

  const stateRef = useRef<{
    renderer?: THREE.WebGLRenderer;
    scene?: THREE.Scene;
    camera?: THREE.PerspectiveCamera;
    orbit?: OrbitControls;
    gizmo?: TransformControls;
    grid?: THREE.GridHelper;
    axes?: THREE.AxesHelper;
    reflector?: Reflector;
    envMap?: THREE.Texture;
    lights?: {
      ambient: THREE.AmbientLight;
      key: THREE.DirectionalLight;
      fill: THREE.DirectionalLight;
      back: THREE.DirectionalLight;
    };
    items: Map<string, EditorItem>;
    suspended?: boolean;
    resumeLoop?: () => void;
    loader: GLTFLoader;
    modelCache: Map<string, Promise<THREE.Group>>;
  }>({ items: new Map(), loader: new GLTFLoader(), modelCache: new Map() });

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;
    const state = stateRef.current;

    const renderer = new THREE.WebGLRenderer({ antialias: true, preserveDrawingBuffer: true });
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.toneMapping = THREE.ACESFilmicToneMapping;
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFSoftShadowMap;
    container.appendChild(renderer.domElement);

    const scene = new THREE.Scene();
    scene.background = new THREE.Color(BG_PRESETS.lightGray!);

    const camera = new THREE.PerspectiveCamera(45, container.clientWidth / container.clientHeight, 0.01, 2000);
    camera.position.set(6, 4, 6);

    const ambient = new THREE.AmbientLight('#ffffff', 0.6);
    const key = new THREE.DirectionalLight('#ffffff', 4);
    key.position.set(-4, 6, 5);
    key.castShadow = true;
    key.shadow.mapSize.set(2048, 2048);
    key.shadow.bias = -0.0005;
    const fill = new THREE.DirectionalLight('#ffffff', 0.8);
    fill.position.set(5, 3.5, 4);
    const back = new THREE.DirectionalLight('#ffffff', 0.2);
    back.position.set(0, 5, -6);
    scene.add(ambient, key, fill, back);

    const ground = new THREE.Mesh(new THREE.PlaneGeometry(200, 200), new THREE.ShadowMaterial({ opacity: 0.25 }));
    ground.rotation.x = -Math.PI / 2;
    ground.receiveShadow = true;
    scene.add(ground);

    const orbit = new OrbitControls(camera, renderer.domElement);
    orbit.enableDamping = true;

    let gizmo: TransformControls | undefined;
    if (interactive) {
      gizmo = new TransformControls(camera, renderer.domElement);
      gizmo.addEventListener('dragging-changed', (event) => {
        orbit.enabled = !event.value;
        if (!event.value) setTransformTick((tick) => tick + 1);
      });
      gizmo.addEventListener('objectChange', () => setTransformTick((tick) => tick + 1));
      scene.add(gizmo.getHelper());
    }

    // 点选（与拖拽区分：按下/抬起位移小才算点击）
    let downAt: { x: number; y: number } | null = null;
    const onPointerDown = (event: PointerEvent) => {
      downAt = { x: event.clientX, y: event.clientY };
    };
    const onPointerUp = (event: PointerEvent) => {
      if (!interactive || !downAt) return;
      const moved = Math.hypot(event.clientX - downAt.x, event.clientY - downAt.y);
      downAt = null;
      if (moved > 4 || (gizmo as unknown as { dragging: boolean } | undefined)?.dragging) return;
      const rect = container.getBoundingClientRect();
      const pointer = new THREE.Vector2(
        ((event.clientX - rect.left) / rect.width) * 2 - 1,
        -((event.clientY - rect.top) / rect.height) * 2 + 1
      );
      const raycaster = new THREE.Raycaster();
      raycaster.setFromCamera(pointer, camera);
      let hitKey: string | null = null;
      for (const [itemKey, item] of state.items) {
        if (raycaster.intersectObject(item.root, true).length > 0) {
          hitKey = itemKey;
          break;
        }
      }
      selectInternal(hitKey);
    };
    renderer.domElement.addEventListener('pointerdown', onPointerDown);
    renderer.domElement.addEventListener('pointerup', onPointerUp);

    let frameId = 0;
    // 渲染循环：挂起时不再排下一帧（loopRunning 防止恢复时产生重复循环）
    let loopRunning = false;
    const animate = () => {
      if (state.suspended) {
        loopRunning = false;
        return;
      }
      frameId = requestAnimationFrame(animate);
      orbit.update();
      renderer.render(scene, camera);
    };
    const startLoop = () => {
      if (loopRunning) return;
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

    Object.assign(state, { renderer, scene, camera, orbit, gizmo, lights: { ambient, key, fill, back } });
    setReady(true);

    return () => {
      cancelAnimationFrame(frameId);
      observer.disconnect();
      renderer.domElement.removeEventListener('pointerdown', onPointerDown);
      renderer.domElement.removeEventListener('pointerup', onPointerUp);
      gizmo?.dispose();
      orbit.dispose();
      renderer.dispose();
      if (renderer.domElement.parentNode === container) container.removeChild(renderer.domElement);
      state.items.clear();
      state.modelCache.clear();
      setReady(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const selectInternal = useCallback((key: string | null) => {
    const state = stateRef.current;
    if (!state.gizmo) return;
    if (key && state.items.has(key)) {
      state.gizmo.attach(state.items.get(key)!.root);
    } else {
      state.gizmo.detach();
    }
    setSelectedKey(key);
  }, []);

  const loadModel = useCallback((fileUrl: string) => {
    const state = stateRef.current;
    if (!state.modelCache.has(fileUrl)) {
      state.modelCache.set(
        fileUrl,
        state.loader.loadAsync(fileUrl).then((gltf) => {
          gltf.scene.traverse((object) => {
            if ((object as THREE.Mesh).isMesh) {
              object.castShadow = true;
              object.receiveShadow = true;
            }
          });
          return gltf.scene;
        })
      );
    }
    return state.modelCache.get(fileUrl)!;
  }, []);

  const addItem = useCallback(
    async (options: {
      key: string;
      instanceId: number;
      fileUrl: string;
      instanceHeight?: number;
      placement: ScenePlacement;
      motionMappings?: MotionMappingNode[];
    }) => {
      const state = stateRef.current;
      if (!state.scene || state.items.has(options.key)) return;
      const template = await loadModel(options.fileUrl);
      const root = new THREE.Group();
      const clone = template.clone(true);
      // 实例实际高度 → 基础等比缩放（AR 语义），用户 placement.scale 在此之上叠加
      const bounds = new THREE.Box3().setFromObject(clone);
      const naturalHeight = bounds.getSize(new THREE.Vector3()).y || 1;
      const baseScale =
        options.instanceHeight && options.instanceHeight > 0 ? options.instanceHeight / naturalHeight : 1;
      clone.scale.setScalar(baseScale);
      root.add(clone);
      root.position.fromArray(options.placement.position);
      root.rotation.set(...options.placement.rotation);
      root.scale.fromArray(options.placement.scale);
      state.scene.add(root);
      const item: EditorItem = {
        key: options.key,
        instanceId: options.instanceId,
        root,
        inner: clone,
        pivot: 'bottomCenter',
        baseScale,
        bindings: [],
      };
      // 恢复保存的 pivot（保存的 position 已是 pivot 点位置，只需恢复 inner 偏移，不做位置补偿）
      if (options.placement.activePivot === 'center') applyPivot(item, 'center', false);
      // 场景运行时：把实例的 motionMappings 解析到克隆体节点上，供 MQTT 消息驱动
      if (options.motionMappings?.length) {
        for (const mapping of options.motionMappings) {
          if (!mapping.value) continue;
          const node = findObjectByPath(clone, mapping.path);
          if (!node) continue;
          item.bindings.push({
            binding: mapping,
            node,
            originalPosition: node.position.clone(),
            originalRotation: node.rotation.clone(),
          });
        }
      }
      state.items.set(options.key, item);
      setTransformTick((tick) => tick + 1);
    },
    [loadModel]
  );

  // 场景运行时：把 MQTT payload 应用到某个实例条目的绑定节点
  const applyPayloadToItem = useCallback((key: string, payload: Record<string, unknown>) => {
    const item = stateRef.current.items.get(key);
    if (item?.bindings.length) applyPayloadToBindings(payload, item.bindings);
  }, []);

  const removeItem = useCallback(
    (key: string) => {
      const state = stateRef.current;
      const item = state.items.get(key);
      if (!item || !state.scene) return;
      if (selectedKey === key) selectInternal(null);
      state.scene.remove(item.root);
      state.items.delete(key);
      setTransformTick((tick) => tick + 1);
    },
    [selectedKey, selectInternal]
  );

  const getPlacement = useCallback((key: string): ScenePlacement | null => {
    const item = stateRef.current.items.get(key);
    if (!item) return null;
    return {
      position: item.root.position.toArray() as [number, number, number],
      rotation: [item.root.rotation.x, item.root.rotation.y, item.root.rotation.z],
      scale: item.root.scale.toArray() as [number, number, number],
      activePivot: item.pivot,
    };
  }, []);

  const setPlacement = useCallback((key: string, placement: ScenePlacement) => {
    const item = stateRef.current.items.get(key);
    if (!item) return;
    item.root.position.fromArray(placement.position);
    item.root.rotation.set(...placement.rotation);
    item.root.scale.fromArray(placement.scale);
    setTransformTick((tick) => tick + 1);
  }, []);

  const setMode = useCallback((mode: GizmoMode) => {
    stateRef.current.gizmo?.setMode(mode);
  }, []);

  const setPivot = useCallback((key: string, mode: ScenePivotMode) => {
    const item = stateRef.current.items.get(key);
    if (!item || item.pivot === mode) return;
    applyPivot(item, mode, true);
    setTransformTick((tick) => tick + 1);
  }, []);

  const getPivot = useCallback(
    (key: string): ScenePivotMode | null => stateRef.current.items.get(key)?.pivot ?? null,
    []
  );

  const setSpace = useCallback((space: 'local' | 'world') => {
    stateRef.current.gizmo?.setSpace(space);
  }, []);

  // 重置视角：按所有条目的整体包围盒重新取景
  const frameAll = useCallback(() => {
    const state = stateRef.current;
    if (!state.camera || !state.orbit) return;
    const bounds = new THREE.Box3();
    let hasContent = false;
    state.items.forEach((item) => {
      bounds.expandByObject(item.root);
      hasContent = true;
    });
    const camera = state.camera;
    if (!hasContent || bounds.isEmpty()) {
      camera.position.set(6, 4, 6);
      state.orbit.target.set(0, 0, 0);
      state.orbit.update();
      return;
    }
    const size = bounds.getSize(new THREE.Vector3());
    const center = bounds.getCenter(new THREE.Vector3());
    const maxDim = Math.max(size.x, size.y, size.z) || 1;
    const distance = (maxDim / (2 * Math.tan(THREE.MathUtils.degToRad(camera.fov / 2)))) * 1.6;
    camera.position.set(center.x + distance * 0.9, center.y + distance * 0.55, center.z + distance * 0.9);
    state.orbit.target.copy(center);
    state.orbit.update();
  }, []);

  // 挂起/恢复渲染循环（页签隐藏或窗口后台时省资源）；恢复幂等
  const setSuspended = useCallback((suspended: boolean) => {
    const state = stateRef.current;
    if (state.suspended === suspended) return;
    state.suspended = suspended;
    if (!suspended) state.resumeLoop?.();
  }, []);

  // 取景单个条目（扫码分享 focusInstanceId 语义：相机拉近到该实例）
  const frameItem = useCallback((key: string) => {
    const state = stateRef.current;
    const item = state.items.get(key);
    if (!state.camera || !state.orbit || !item) return;
    const bounds = new THREE.Box3().setFromObject(item.root);
    if (bounds.isEmpty()) return;
    const camera = state.camera;
    const size = bounds.getSize(new THREE.Vector3());
    const center = bounds.getCenter(new THREE.Vector3());
    const maxDim = Math.max(size.x, size.y, size.z) || 1;
    const distance = (maxDim / (2 * Math.tan(THREE.MathUtils.degToRad(camera.fov / 2)))) * 1.6;
    camera.position.set(center.x + distance * 0.9, center.y + distance * 0.55, center.z + distance * 0.9);
    state.orbit.target.copy(center);
    state.orbit.update();
  }, []);

  const applyConfig = useCallback((config: SceneConfigV4) => {
    const state = stateRef.current;
    if (!state.scene || !state.lights || !state.camera) return;
    const bg = BG_PRESETS[config.viewport.backgroundPreset];
    state.scene.background = bg ? new THREE.Color(bg) : gradientTexture();

    // 灯光强度直接应用（对齐源端 render-effects：environment=IBL 环境贴图，不再缩放各灯）
    const preset = LIGHT_PRESETS[config.light.preset] ?? LIGHT_PRESETS.balanced;
    state.lights.ambient.intensity = config.light.ambientIntensity ?? preset.ambient;
    state.lights.key.intensity = config.light.keyLightIntensity ?? preset.key;
    state.lights.fill.intensity = config.light.fillLightIntensity ?? preset.fill;
    state.lights.back.intensity = config.light.backLightIntensity ?? preset.back;
    if (config.light.environmentEnabled && state.renderer) {
      if (!state.envMap) {
        const pmrem = new THREE.PMREMGenerator(state.renderer);
        state.envMap = pmrem.fromScene(new RoomEnvironment(), 0.04).texture;
        pmrem.dispose();
      }
      state.scene.environment = state.envMap;
      state.scene.environmentIntensity = config.light.environmentIntensity ?? 1;
    } else {
      state.scene.environment = null;
    }

    // 相机：焦距（35mm 胶片等效，与源端 camera.setFocalLength 一致）+ Shift X/Y（视口偏移）+ 裁剪
    state.camera.setFocalLength(config.viewport.camera.focalLength || 50);
    const canvasEl = state.renderer?.domElement;
    const viewWidth = canvasEl?.clientWidth || 1;
    const viewHeight = canvasEl?.clientHeight || 1;
    const offsetX = Math.round((config.viewport.camera.shiftX || 0) * viewWidth);
    const offsetY = Math.round((config.viewport.camera.shiftY || 0) * viewHeight);
    if (offsetX !== 0 || offsetY !== 0) {
      state.camera.setViewOffset(viewWidth, viewHeight, offsetX, offsetY, viewWidth, viewHeight);
    } else {
      state.camera.clearViewOffset();
    }
    if (config.viewport.camera.clipMode === 'manual') {
      state.camera.near = config.viewport.camera.clipStart || 0.01;
      state.camera.far = config.viewport.camera.clipEnd || 2000;
    } else {
      // auto：按场景包围盒推导（对齐源端 applyCameraAutoClip 的量级）
      const bounds = new THREE.Box3();
      state.items.forEach((item) => bounds.expandByObject(item.root));
      const maxDim = bounds.isEmpty() ? 10 : Math.max(...bounds.getSize(new THREE.Vector3()).toArray(), 1);
      state.camera.near = Math.max(maxDim / 1000, 0.001);
      state.camera.far = Math.max(maxDim * 50, 100);
    }
    state.camera.updateProjectionMatrix();

    state.grid?.removeFromParent();
    state.axes?.removeFromParent();
    if (config.viewport.grid.show) {
      const size = config.viewport.grid.infinite ? 400 : config.viewport.grid.areaSize || 20;
      const divisions = Math.max(2, Math.round(size / Math.max(config.viewport.grid.minCellSize || 1, 0.1)));
      state.grid = new THREE.GridHelper(size, divisions, 0xc6c6c6, 0xe0e0e0);
      state.scene.add(state.grid);
    }
    if (config.viewport.grid.showAxis) {
      state.axes = new THREE.AxesHelper((config.viewport.grid.areaSize || 20) / 2);
      state.scene.add(state.axes);
    }

    // 地面反射：clarity → 反射贴图分辨率（低清晰度自然发虚），depth → 混合不透明度（对齐源端语义）
    if (state.reflector) {
      state.reflector.removeFromParent();
      state.reflector.dispose();
      state.reflector = undefined;
    }
    const reflection = config.viewport.reflection;
    if (reflection?.enabled && reflection.depth > 0) {
      const planeSize = config.viewport.grid.infinite ? 400 : (config.viewport.grid.areaSize || 20) * 2;
      const reflector = new Reflector(new THREE.PlaneGeometry(planeSize, planeSize), {
        textureWidth: reflectionResolution(reflection.clarity),
        textureHeight: reflectionResolution(reflection.clarity),
        color: 0x889999,
        clipBias: 0.003,
      });
      // Reflector 自带 shader 不支持透明度，注入 opacity uniform 实现 depth 混合
      const material = reflector.material as THREE.ShaderMaterial;
      material.transparent = true;
      material.uniforms.reflectOpacity = { value: (reflection.depth / 100) * 0.55 };
      material.fragmentShader = material.fragmentShader
        .replace('uniform vec3 color;', 'uniform vec3 color;\nuniform float reflectOpacity;')
        .replace(
          /gl_FragColor = vec4\( blendOverlay\( base.rgb, color \), 1.0 \);/,
          'gl_FragColor = vec4( blendOverlay( base.rgb, color ), reflectOpacity );'
        )
        .replace(
          /vec4\( blendOverlay\( base.rgb, color \), 1.0 \)/,
          'vec4( blendOverlay( base.rgb, color ), reflectOpacity )'
        );
      reflector.rotation.x = -Math.PI / 2;
      reflector.position.y = -0.002;
      state.reflector = reflector;
      state.scene.add(reflector);
    }
  }, []);

  const snapshot = useCallback((): Promise<File | null> => {
    const state = stateRef.current;
    if (!state.renderer) return Promise.resolve(null);
    return new Promise((resolve) => {
      state.renderer!.domElement.toBlob((blob) => {
        resolve(blob ? new File([blob], 'scene-thumbnail.png', { type: 'image/png' }) : null);
      }, 'image/png');
    });
  }, []);

  return {
    containerRef,
    ready,
    selectedKey,
    transformTick,
    select: selectInternal,
    addItem,
    removeItem,
    getPlacement,
    setPlacement,
    setMode,
    setSpace,
    setPivot,
    getPivot,
    frameAll,
    frameItem,
    setSuspended,
    applyConfig,
    applyPayloadToItem,
    snapshot,
  };
}
