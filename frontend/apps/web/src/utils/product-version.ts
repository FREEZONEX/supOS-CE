type ProductVersionSource = {
  productVersion?: string;
  version?: string;
};

export const resolveProductVersion = (source?: ProductVersionSource) =>
  String(source?.productVersion || source?.version || 'dev').trim() || 'dev';
