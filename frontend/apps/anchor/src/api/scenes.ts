import { api } from './client';
import type { InstanceInfo } from './instances';

// 与源端 SceneConfigV4 对齐（tier0-frontend anchor/scene/[sceneId]/types.ts）
export type SceneBackgroundPreset = 'white' | 'lightGray' | 'black' | 'gradient';
export type SceneLightPreset = 'balanced' | 'soft' | 'contrast';
export type SceneVector3 = [number, number, number];

export interface ScenePlacement {
  position: SceneVector3;
  rotation: SceneVector3;
  scale: SceneVector3;
  [key: string]: unknown;
}

export interface SceneConfigV4 {
  version: 4;
  viewport: {
    backgroundPreset: SceneBackgroundPreset;
    camera: {
      focalLength: number;
      lensUnit: 'millimeters';
      shiftX: number;
      shiftY: number;
      clipStart: number;
      clipEnd: number;
      clipMode: 'auto' | 'manual';
    };
    grid: { show: boolean; showAxis: boolean; areaSize: number; minCellSize: number; infinite: boolean };
    reflection: { enabled: boolean; clarity: number; depth: number };
  };
  light: {
    preset: SceneLightPreset;
    environmentEnabled: boolean;
    environmentIntensity: number;
    ambientIntensity: number;
    keyLightIntensity: number;
    fillLightIntensity: number;
    backLightIntensity: number;
  };
}

export const defaultSceneConfig = (): SceneConfigV4 => ({
  version: 4,
  viewport: {
    backgroundPreset: 'lightGray',
    camera: {
      focalLength: 50,
      lensUnit: 'millimeters',
      shiftX: 0,
      shiftY: 0,
      clipStart: 0.01,
      clipEnd: 2000,
      clipMode: 'auto',
    },
    grid: { show: true, showAxis: true, areaSize: 20, minCellSize: 1, infinite: false },
    // 对齐源端 DEFAULT_GROUND_REFLECTION
    reflection: { enabled: true, clarity: 60, depth: 45 },
  },
  light: {
    preset: 'balanced',
    environmentEnabled: true,
    environmentIntensity: 1,
    ambientIntensity: 0.1,
    keyLightIntensity: 4,
    fillLightIntensity: 0.8,
    backLightIntensity: 0.2,
  },
});

export const defaultPlacement = (): ScenePlacement => ({
  position: [0, 0, 0],
  rotation: [0, 0, 0],
  scale: [1, 1, 1],
});

export function parseSceneConfig(raw: string | undefined): SceneConfigV4 {
  const fallback = defaultSceneConfig();
  if (!raw) return fallback;
  try {
    const parsed = JSON.parse(raw) as Partial<SceneConfigV4>;
    if (!parsed || typeof parsed !== 'object' || !parsed.viewport) return fallback;
    return {
      ...fallback,
      ...parsed,
      viewport: {
        ...fallback.viewport,
        ...parsed.viewport,
        grid: { ...fallback.viewport.grid, ...parsed.viewport?.grid },
        camera: { ...fallback.viewport.camera, ...parsed.viewport?.camera },
        reflection: { ...fallback.viewport.reflection, ...parsed.viewport?.reflection },
      },
      light: { ...fallback.light, ...parsed.light },
    } as SceneConfigV4;
  } catch {
    return fallback;
  }
}

export function parsePlacement(raw: string | undefined): ScenePlacement {
  const fallback = defaultPlacement();
  if (!raw) return fallback;
  try {
    const parsed = JSON.parse(raw) as Partial<ScenePlacement>;
    return {
      ...fallback,
      ...parsed,
      position: (parsed.position as SceneVector3) ?? fallback.position,
      rotation: (parsed.rotation as SceneVector3) ?? fallback.rotation,
      scale: (parsed.scale as SceneVector3) ?? fallback.scale,
    };
  } catch {
    return fallback;
  }
}

export interface SceneItem {
  id: number;
  instanceId: number;
  placementJson: string;
  sort: number;
  instance?: InstanceInfo;
}

export interface SceneInfo {
  id: number;
  name: string;
  description: string;
  configJson: string;
  thumbnailAssetId: number;
  thumbnailUrl: string;
  items: SceneItem[];
  itemCount: number;
  createdTime: string;
  updatedTime: string;
}

export interface SceneListResult {
  list: SceneInfo[];
  total: number;
  page: number;
  size: number;
}

const BASE = '/api/core/anchor';

export const listScenes = (params: { page?: number; size?: number; keyword?: string }) =>
  api.get<SceneListResult>(`${BASE}/scenes`, params);

export const getScene = (id: number | string) => api.get<SceneInfo>(`${BASE}/scenes/${id}`);

export const createScene = (body: { name: string; description?: string; configJson?: string }) =>
  api.post<SceneInfo>(`${BASE}/scenes`, body);

export const updateScene = (
  id: number,
  body: {
    name?: string;
    description?: string;
    configJson?: string;
    thumbnailAssetId?: number;
    items?: { instanceId: number; placementJson: string; sort: number }[];
    itemsSet?: boolean;
  }
) => api.put<SceneInfo>(`${BASE}/scenes/${id}`, body);

export const deleteScene = (id: number) => api.delete<{ id: number }>(`${BASE}/scenes/${id}`);

// 场景二维码分享配置（免登录接口，供 /viewer 扫码页使用）
export interface SceneQrItem {
  instanceId: number;
  placementJson: string;
  sort: number;
  instance: {
    id: number;
    name: string;
    modelName: string;
    modelFileUrl: string;
    height: number;
    topic: string;
    bindingJson: string;
  };
}

export interface SceneQrConfig {
  scene: {
    id: number;
    name: string;
    description: string;
    configJson: string;
    items: SceneQrItem[];
  };
  mqtt: {
    username: string;
    password: string;
    clientId: string;
    wsPort: number;
    wssPort: number;
    path: string;
  };
}

export const sceneQrConfigUrl = (id: number | string) => `${BASE}/scenes/${id}/qr-config`;
