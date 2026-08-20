// 移植自 tier0-frontend anchor/viewer-primitives/live-binding-core.ts 的核心算法：
// delta = (payload[value] - payload[init]) * factor * sign(axis)，叠加到节点原始位姿上。
import type * as THREE from 'three';
import type { MotionMappingNode, SignedAxis } from '../api/instances';
import { payloadNumber } from '../mqtt/use-mqtt';

export interface ResolvedBinding {
  binding: MotionMappingNode;
  node: THREE.Object3D;
  originalPosition: THREE.Vector3;
  originalRotation: THREE.Euler;
}

// 源端约定：值键 "J1_Cur" 对应初始键 "J1_Init"；无 _Cur 后缀时初始值取 0
export function initKeyFor(valueKey: string): string | null {
  if (/_cur$/i.test(valueKey)) return valueKey.replace(/_cur$/i, '_Init');
  return null;
}

export function resolveBindingValue(
  payload: Record<string, unknown> | undefined,
  valueKey: string
): number | undefined {
  const current = payloadNumber(payload, valueKey);
  if (current === undefined) return undefined;
  const initKey = initKeyFor(valueKey);
  const init = initKey ? (payloadNumber(payload, initKey) ?? 0) : 0;
  return current - init;
}

const axisLetter = (axis: SignedAxis): 'x' | 'y' | 'z' => axis[1] as 'x' | 'y' | 'z';
const axisSign = (axis: SignedAxis): number => (axis.startsWith('-') ? -1 : 1);

export function applyPayloadToBindings(payload: Record<string, unknown>, bindings: ResolvedBinding[]): void {
  for (const { binding, node, originalPosition, originalRotation } of bindings) {
    const delta = resolveBindingValue(payload, binding.value);
    if (delta === undefined) continue;
    const scaled = delta * (binding.factor || 1) * axisSign(binding.axis);
    const letter = axisLetter(binding.axis);
    if (binding.type === 'position') {
      node.position[letter] = originalPosition[letter] + scaled;
    } else {
      // 数值默认为角度制（源端机器人关节 J*_Cur 为 deg），转弧度叠加
      node.rotation[letter] = originalRotation[letter] + (scaled * Math.PI) / 180;
    }
  }
}

export function resetBindings(bindings: ResolvedBinding[]): void {
  for (const { node, originalPosition, originalRotation } of bindings) {
    node.position.copy(originalPosition);
    node.rotation.copy(originalRotation);
  }
}
