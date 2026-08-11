export type NotebookTabType = 'all' | 'running';

export interface FolderNode {
  id: number;
  name: string;
  parentId: number;
  path: string;
  notebookCount?: number;
  notebooks?: NotebookTreeItem[];
  children?: FolderNode[];
}

export interface NotebookTreeItem {
  id: number;
  name: string;
  path: string;
  folderId: number;
  status?: string;
  updatedAt?: string;
}

export interface NotebookItem {
  id: number;
  type: 'folder' | 'notebook';
  name: string;
  path: string;
  isDirectory: boolean;
  isFavorite?: boolean;
  description?: string;
  folderId?: number;
  status?: string;
  owner?: string;
  fileSize?: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface NotebookDetail {
  id: number;
  name: string;
  description?: string;
  folderId: number;
  filePath: string;
  status: string;
  owner?: string;
  isFavorite?: boolean;
  fileSize: number;
  createdAt?: string;
  updatedAt?: string;
}
