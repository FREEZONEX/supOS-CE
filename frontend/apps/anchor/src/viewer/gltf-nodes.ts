// 移植自 tier0-frontend anchor/utils/gltf-node-metadata.ts 的路径规则：
// 节点路径 = "Scene/父/子" 层级拼接；tier0Id 优先取 GLB extras 注入的 userData.tier0Id，退化用路径本身。
import type * as THREE from 'three';

export interface SceneNodeMeta {
  key: string;
  title: string;
  path: string;
  tier0Id: string;
  type: string;
  isLeaf: boolean;
  children: SceneNodeMeta[];
}

export function buildNodeTree(root: THREE.Object3D): SceneNodeMeta[] {
  const build = (object: THREE.Object3D, parentPath: string): SceneNodeMeta => {
    const name = object.name || object.type || 'Node';
    const path = parentPath ? `${parentPath}/${name}` : name;
    const tier0Id = (object.userData?.tier0Id as string) || path;
    const children = object.children.map((child) => build(child, path));
    return {
      key: path,
      title: name,
      path,
      tier0Id,
      type: object.type,
      isLeaf: children.length === 0,
      children,
    };
  };
  const rootName = root.name || 'Scene';
  return [build(root, '')].map((node) => ({ ...node, title: rootName }));
}

export function flattenNodeTree(nodes: SceneNodeMeta[]): SceneNodeMeta[] {
  const out: SceneNodeMeta[] = [];
  const walk = (items: SceneNodeMeta[]) => {
    items.forEach((item) => {
      out.push(item);
      walk(item.children);
    });
  };
  walk(nodes);
  return out;
}

export function findObjectByPath(root: THREE.Object3D, path: string): THREE.Object3D | null {
  const rootName = root.name || root.type || 'Node';
  if (path === rootName) return root;
  const prefix = `${rootName}/`;
  if (!path.startsWith(prefix)) return null;
  const segments = path.slice(prefix.length).split('/');
  let current: THREE.Object3D = root;
  for (const segment of segments) {
    const next = current.children.find((child) => (child.name || child.type || 'Node') === segment);
    if (!next) return null;
    current = next;
  }
  return current;
}
