const PATH_SEGMENT_REG = /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/;
const RESERVED_FOLDER_NAMES = new Set(['label', 'template']);
const MAX_PATH_SEGMENT_LENGTH = 63;

const CATEGORY_SEGMENT_TO_PARENT_DATA_TYPE: Record<string, number> = {
  State: 1,
  Action: 2,
  Metric: 3,
};

const ALLOWED_DATA_TYPES_BY_PARENT: Record<number, number[]> = {
  1: [2, 6, 7, 8],
  2: [8],
  3: [1, 3, 4, 6, 7],
};

export type NamespaceValidationResult = {
  key: string;
  values?: Record<string, number>;
} | null;

export const splitNamespaceSegments = (name?: string) => (name || '').split('/');

export const isMultiSegmentName = (name?: string) => !!name && name.includes('/');

export const getLeafSegment = (name?: string) => {
  const segments = splitNamespaceSegments(name).filter(Boolean);
  return segments[segments.length - 1] || '';
};

export const buildNamespacePreview = (sourcePath?: string, name?: string) =>
  isMultiSegmentName(name) ? name || '' : `${sourcePath || ''}${name || ''}`;

export const deriveParentDataTypeFromName = (name?: string) => {
  if (!isMultiSegmentName(name)) return undefined;
  const segments = splitNamespaceSegments(name).filter(Boolean);
  if (segments.length < 2) return undefined;
  return CATEGORY_SEGMENT_TO_PARENT_DATA_TYPE[segments[segments.length - 2]];
};

export const getDefaultDataTypeByParent = (parentDataType?: number) => {
  switch (parentDataType) {
    case 1:
    case 2:
      return 8;
    case 3:
      return 1;
    default:
      return undefined;
  }
};

export const isDataTypeMatched = (parentDataType?: number, dataType?: number) => {
  if (!parentDataType || !dataType) return false;
  return ALLOWED_DATA_TYPES_BY_PARENT[parentDataType]?.includes(dataType) || false;
};

export const validateNamespaceName = (
  name: string,
  options: {
    isCreateFolder: boolean;
    dataType?: number;
    parentDataType?: number;
  }
): NamespaceValidationResult => {
  if (!name) return null;
  if (!options.isCreateFolder && name.endsWith('/')) {
    return { key: 'uns.namespaceCategoryRequired' };
  }
  if (name.startsWith('/') || name.endsWith('/') || name.includes('//')) {
    return { key: 'uns.namespaceNameFormat' };
  }

  const segments = splitNamespaceSegments(name);
  const lastDirIndex = segments.length - 2;

  for (const [index, segment] of segments.entries()) {
    if (!segment || !PATH_SEGMENT_REG.test(segment)) {
      return { key: 'uns.namespaceNameFormat' };
    }
    if (segment.length > MAX_PATH_SEGMENT_LENGTH) {
      return { key: 'uns.namespaceSegmentTooLong', values: { length: MAX_PATH_SEGMENT_LENGTH } };
    }
    if (options.isCreateFolder || index <= lastDirIndex) {
      if (RESERVED_FOLDER_NAMES.has(segment.toLowerCase())) {
        return { key: 'uns.prohibitKeywords' };
      }
    }
    if (options.isCreateFolder && index < segments.length - 1) {
      if (CATEGORY_SEGMENT_TO_PARENT_DATA_TYPE[segment]) {
        return { key: 'uns.namespaceCategoryFolderChildForbidden' };
      }
    }
  }

  if (options.isCreateFolder && !isMultiSegmentName(name) && (options.parentDataType || 0) > 0) {
    return { key: 'uns.namespaceCategoryFolderChildForbidden' };
  }

  if (!options.isCreateFolder && segments.length > 1) {
    for (const segment of segments.slice(0, -2)) {
      if (CATEGORY_SEGMENT_TO_PARENT_DATA_TYPE[segment]) {
        return { key: 'uns.namespaceCategoryFolderChildForbidden' };
      }
    }
    const parentDataType = deriveParentDataTypeFromName(name);
    if (!parentDataType) {
      return { key: 'uns.namespaceCategoryRequired' };
    }
    if (options.dataType && !isDataTypeMatched(parentDataType, options.dataType)) {
      return { key: 'uns.namespaceCategoryMismatch' };
    }
  }

  return null;
};
