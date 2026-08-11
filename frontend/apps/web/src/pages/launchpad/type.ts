/**
 * Launchpad 相关类型定义
 */

import type { ProjectItem } from '../project/types';

/**
 * 应用对象（对应 LaunchpadAppVO）
 */
export interface App {
  /** 应用 ID */
  appId: number;
  /** 应用名称（系统名） */
  appName: string;
  /** 应用展示名 */
  displayName?: string;
  /** 应用图标地址 */
  iconUrl?: string;
  /** 应用封面地址 */
  coverUrl?: string;
  /** 应用访问地址 */
  url: string;
  /** 应用类型 */
  siteType?: string;
  /** 应用来源类型 */
  appType?: string;
  /** 是否手工添加 */
  manual?: boolean;
  /** 是否在平台内打开 */
  openInPlatform?: boolean;
  /** 是否已收藏（后端就绪后回填；收藏功能后端未上线前前端本地维护） */
  isFavorite?: boolean;
  /** 所属项目名（收藏/最近区展示所属项目用） */
  projectName?: string;
  /** 所属项目展示名 */
  projectDisplayName?: string;
  /** 创建人 */
  createdBy: string;
  /** 创建时间 */
  createdAt: string | number;
  /** 更新时间 */
  updatedAt: string | number;
}

/**
 * 项目对象（前端统一使用）
 */
export interface Project {
  /** 项目 ID */
  projectId: number;
  /** 项目名称（系统名） */
  projectName: string;
  /** 项目展示名 */
  displayName?: string;
  /** 项目描述 */
  description: string;
  /** 是否已收藏（后端就绪后回填） */
  isFavorite?: boolean;
  /** 是否可编辑 */
  canEdit: boolean;
  /** 项目下应用列表 */
  apps: App[];
  /** 创建人 */
  createdBy: string;
  /** 创建时间 */
  createdAt: string | number;
  /** 更新时间 */
  updatedAt: string | number;
}

export interface ProjectDetail {
  /** 是否可编辑 */
  canEdit: boolean;
  /** 项目下应用列表 */
  apps: App[];
  project: ProjectItem;
}
