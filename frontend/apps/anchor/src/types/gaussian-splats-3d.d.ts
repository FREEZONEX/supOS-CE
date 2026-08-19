// @mkkellogg/gaussian-splats-3d 无官方类型；仅声明本项目用到的最小接口
declare module '@mkkellogg/gaussian-splats-3d' {
  export enum SceneFormat {
    Splat = 0,
    KSplat = 1,
    Ply = 2,
  }

  export interface ViewerOptions {
    rootElement?: HTMLElement;
    /** 外部 three.js Scene（背景/网格等由调用方注入） */
    threeScene?: unknown;
    selfDrivenMode?: boolean;
    useBuiltInControls?: boolean;
    sharedMemoryForWorkers?: boolean;
    cameraUp?: number[];
    initialCameraPosition?: number[];
    initialCameraLookAt?: number[];
  }

  export interface AddSplatSceneOptions {
    format?: SceneFormat;
    showLoadingUI?: boolean;
    progressiveLoad?: boolean;
    /** 场景变换：四元数 [x,y,z,w] / 平移 [x,y,z] / 缩放 [x,y,z] */
    rotation?: number[];
    position?: number[];
    scale?: number[];
  }

  export class Viewer {
    constructor(options?: ViewerOptions);
    addSplatScene(url: string, options?: AddSplatSceneOptions): Promise<void>;
    start(): void;
    stop(): void;
    dispose(): void;
  }
}
