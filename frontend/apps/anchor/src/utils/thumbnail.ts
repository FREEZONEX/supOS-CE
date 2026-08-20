// 移植自 tier0-frontend anchor/utils/thumbnail-utils.ts（去掉主题色依赖与 splat WebGL 渲染器，
// splat 直接使用投影降级方案；主题色固定为源端默认浅色）
import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';

const THUMBNAIL_WIDTH = 640;
const THUMBNAIL_HEIGHT = 320;
const THUMBNAIL_BG = '#eef2f7';
const SPLAT_VERTEX_SIZE = 32;
const SPLAT_POSITION_OFFSET = 0;
const SPLAT_COLOR_OFFSET = 24;

const getBaseName = (fileName: string) => fileName.replace(/\.[^/.]+$/, '');
const sanitizeFileName = (fileName: string) => fileName.replace(/[^a-zA-Z0-9._-]+/g, '-');

export function getModelFormat(fileName: string): 'glb' | 'gltf' | 'splat' | 'model' {
  const ext = fileName.split('.').pop()?.toLowerCase();
  if (ext === 'glb' || ext === 'gltf' || ext === 'splat') return ext;
  return 'model';
}

const blobToFile = (blob: Blob, fileName: string) =>
  new File([blob], fileName, { type: blob.type || 'application/octet-stream' });

const canvasToBlob = (canvas: HTMLCanvasElement, type: string) =>
  new Promise<Blob>((resolve, reject) => {
    canvas.toBlob((blob) => (blob ? resolve(blob) : reject(new Error('Failed to export thumbnail canvas.'))), type);
  });

const waitForNextFrame = () => new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));

function createPlaceholderSvg(fileName: string) {
  const baseName = getBaseName(fileName);
  const format = getModelFormat(fileName).toUpperCase();
  const escapedName = baseName
    .slice(0, 32)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;');
  return `<?xml version="1.0" encoding="UTF-8"?>
<svg width="${THUMBNAIL_WIDTH}" height="${THUMBNAIL_HEIGHT}" viewBox="0 0 ${THUMBNAIL_WIDTH} ${THUMBNAIL_HEIGHT}" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="640" height="320" rx="24" fill="${THUMBNAIL_BG}"/>
  <path d="M246 124L320 84L394 124V196L320 236L246 196V124Z" fill="rgba(255,255,255,0.9)"/>
  <path d="M320 84V236" stroke="rgba(71,85,105,0.35)" stroke-width="6"/>
  <path d="M246 124L320 166L394 124" stroke="rgba(71,85,105,0.35)" stroke-width="6"/>
  <rect x="40" y="36" width="76" height="30" rx="15" fill="#0F172A"/>
  <text x="78" y="56" text-anchor="middle" font-family="Arial, sans-serif" font-size="14" font-weight="700" fill="#FFFFFF">${format}</text>
  <text x="48" y="280" font-family="Arial, sans-serif" font-size="24" font-weight="700" fill="#0F172A">${escapedName}</text>
</svg>`;
}

function generatePlaceholderThumbnailFile(fileName: string) {
  const svg = createPlaceholderSvg(fileName);
  return new File([svg], `${sanitizeFileName(getBaseName(fileName))}-thumbnail.svg`, { type: 'image/svg+xml' });
}

async function generateGltfThumbnailFile(file: File) {
  const objectUrl = URL.createObjectURL(file);
  const canvas = document.createElement('canvas');
  canvas.width = THUMBNAIL_WIDTH;
  canvas.height = THUMBNAIL_HEIGHT;
  let renderer: THREE.WebGLRenderer | null = null;
  try {
    renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true, preserveDrawingBuffer: true });
    renderer.setSize(THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT, false);
    renderer.setPixelRatio(1);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.toneMapping = THREE.ACESFilmicToneMapping;

    const scene = new THREE.Scene();
    scene.background = new THREE.Color(THUMBNAIL_BG);

    const camera = new THREE.PerspectiveCamera(32, THUMBNAIL_WIDTH / THUMBNAIL_HEIGHT, 0.01, 2000);
    const ambientLight = new THREE.AmbientLight('#ffffff', 2.4);
    const keyLight = new THREE.DirectionalLight('#ffffff', 2.2);
    const fillLight = new THREE.DirectionalLight('#dbeafe', 1.4);
    keyLight.position.set(3, 4, 5);
    fillLight.position.set(-4, 2, -3);
    scene.add(ambientLight, keyLight, fillLight);

    const gltf = await new GLTFLoader().loadAsync(objectUrl);
    scene.add(gltf.scene);

    const bounds = new THREE.Box3().setFromObject(gltf.scene);
    if (bounds.isEmpty()) throw new Error('Model bounds are empty.');
    const size = bounds.getSize(new THREE.Vector3());
    const center = bounds.getCenter(new THREE.Vector3());
    const maxDim = Math.max(size.x, size.y, size.z) || 1;
    const fitHeightDistance = maxDim / (2 * Math.tan(THREE.MathUtils.degToRad(camera.fov / 2)));
    const distance = fitHeightDistance * 1.35;
    camera.position.set(center.x + distance * 0.95, center.y + distance * 0.45, center.z + distance * 0.95);
    camera.lookAt(center);
    camera.near = Math.max(distance / 100, 0.01);
    camera.far = Math.max(distance * 10, 100);
    camera.updateProjectionMatrix();

    await waitForNextFrame();
    renderer.render(scene, camera);

    const blob = await canvasToBlob(canvas, 'image/png');
    return blobToFile(blob, `${sanitizeFileName(getBaseName(file.name))}-thumbnail.png`);
  } finally {
    URL.revokeObjectURL(objectUrl);
    renderer?.dispose();
  }
}

type SplatPoint = { x: number; y: number; z: number; color: string; alpha: number; depth: number };

async function generateSplatThumbnailProjectionFile(file: File) {
  const buffer = await file.arrayBuffer();
  const view = new DataView(buffer);
  const vertexCount = Math.floor(buffer.byteLength / SPLAT_VERTEX_SIZE);
  if (vertexCount === 0) throw new Error('Splat file is empty.');

  let minX = Infinity,
    minY = Infinity,
    minZ = Infinity,
    maxX = -Infinity,
    maxY = -Infinity,
    maxZ = -Infinity;
  for (let index = 0; index < vertexCount; index += 1) {
    const offset = index * SPLAT_VERTEX_SIZE + SPLAT_POSITION_OFFSET;
    const x = view.getFloat32(offset, true);
    const y = -view.getFloat32(offset + 4, true);
    const z = -view.getFloat32(offset + 8, true);
    if (!Number.isFinite(x) || !Number.isFinite(y) || !Number.isFinite(z)) continue;
    minX = Math.min(minX, x);
    minY = Math.min(minY, y);
    minZ = Math.min(minZ, z);
    maxX = Math.max(maxX, x);
    maxY = Math.max(maxY, y);
    maxZ = Math.max(maxZ, z);
  }
  if (!Number.isFinite(minX) || !Number.isFinite(maxX)) throw new Error('Failed to parse splat positions.');

  const centerX = (minX + maxX) / 2;
  const centerY = (minY + maxY) / 2;
  const centerZ = (minZ + maxZ) / 2;
  const maxDim = Math.max(maxX - minX, maxY - minY, maxZ - minZ, 1e-6);
  const sampleStep = Math.max(1, Math.ceil(vertexCount / 18000));
  const points: SplatPoint[] = [];

  for (let index = 0; index < vertexCount; index += sampleStep) {
    const offset = index * SPLAT_VERTEX_SIZE + SPLAT_POSITION_OFFSET;
    const x = view.getFloat32(offset, true);
    const y = -view.getFloat32(offset + 4, true);
    const z = -view.getFloat32(offset + 8, true);
    if (!Number.isFinite(x) || !Number.isFinite(y) || !Number.isFinite(z)) continue;
    const normalizedX = (x - centerX) / maxDim;
    const normalizedY = (y - centerY) / maxDim;
    const normalizedZ = (z - centerZ) / maxDim;
    const colorOffset = index * SPLAT_VERTEX_SIZE + SPLAT_COLOR_OFFSET;
    const r = view.getUint8(colorOffset);
    const g = view.getUint8(colorOffset + 1);
    const b = view.getUint8(colorOffset + 2);
    const alpha = Math.max(view.getUint8(colorOffset + 3) / 255, 0.12);
    points.push({
      x: (normalizedX - normalizedZ) * 0.92,
      y: -(normalizedY * 1.2 + (normalizedX + normalizedZ) * 0.34),
      z,
      depth: normalizedX + normalizedY - normalizedZ,
      color: `rgb(${r}, ${g}, ${b})`,
      alpha,
    });
  }
  if (points.length === 0) throw new Error('Failed to sample splat points.');

  let projectedMinX = Infinity,
    projectedMaxX = -Infinity,
    projectedMinY = Infinity,
    projectedMaxY = -Infinity;
  for (const point of points) {
    projectedMinX = Math.min(projectedMinX, point.x);
    projectedMaxX = Math.max(projectedMaxX, point.x);
    projectedMinY = Math.min(projectedMinY, point.y);
    projectedMaxY = Math.max(projectedMaxY, point.y);
  }

  const canvas = document.createElement('canvas');
  canvas.width = THUMBNAIL_WIDTH;
  canvas.height = THUMBNAIL_HEIGHT;
  const context = canvas.getContext('2d');
  if (!context) throw new Error('Failed to create thumbnail context.');

  context.fillStyle = THUMBNAIL_BG;
  context.fillRect(0, 0, THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT);

  const contentWidth = Math.max(projectedMaxX - projectedMinX, 1e-6);
  const contentHeight = Math.max(projectedMaxY - projectedMinY, 1e-6);
  const scale = Math.min((THUMBNAIL_WIDTH * 0.66) / contentWidth, (THUMBNAIL_HEIGHT * 0.66) / contentHeight);
  const offsetX = THUMBNAIL_WIDTH / 2 - ((projectedMinX + projectedMaxX) / 2) * scale;
  const offsetY = THUMBNAIL_HEIGHT / 2 - ((projectedMinY + projectedMaxY) / 2) * scale + 6;

  points.sort((a, b) => a.depth - b.depth);
  context.save();
  context.shadowBlur = 22;
  context.shadowColor = 'rgba(15, 23, 42, 0.12)';
  for (const point of points) {
    const x = point.x * scale + offsetX;
    const y = point.y * scale + offsetY;
    const radius = Math.max(0.9, 2.1 - point.depth * 0.22);
    context.globalAlpha = Math.min(0.9, point.alpha * 0.9);
    context.fillStyle = point.color;
    context.beginPath();
    context.arc(x, y, radius, 0, Math.PI * 2);
    context.fill();
  }
  context.restore();

  const blob = await canvasToBlob(canvas, 'image/png');
  return blobToFile(blob, `${sanitizeFileName(getBaseName(file.name))}-thumbnail.png`);
}

export async function generateModelThumbnailFile(file: File): Promise<File> {
  const fileName = file.name || 'model';
  const format = getModelFormat(fileName);
  if (format === 'splat') {
    try {
      return await generateSplatThumbnailProjectionFile(file);
    } catch {
      return generatePlaceholderThumbnailFile(fileName);
    }
  }
  if (format !== 'glb' && format !== 'gltf') {
    return generatePlaceholderThumbnailFile(fileName);
  }
  try {
    return await generateGltfThumbnailFile(file);
  } catch {
    return generatePlaceholderThumbnailFile(fileName);
  }
}
