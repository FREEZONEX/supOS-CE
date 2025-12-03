export interface MenuProps {
  id: string; //id不传走新增，id传了走更新
  // 资源编码, 唯一
  code: string;
  // 父级ID null 为顶级
  parentId?: string | null;
  // 资源类型1-组 2-菜单 3-按钮（操作权限） 4 tab（home页）
  type: 1 | 2 | 3 | 4;
  // 地址
  url?: string;
  // 类型 1-内部地址 2-外部链接
  urlType: 1 | 2;
  // 打开方式：1-当前页面跳转 2-新窗口打开
  openType: 1 | 2;
  // 图标  传文件上传的附件地址
  icon: string;
  // 描述国际化Key
  description: string;
  // 排序
  sort: number;
  //是否可编辑
  editEnable?: boolean;
  //是否显示在首页
  homeEnable?: boolean;
  //是否启用
  enable?: boolean;
  // 子级
  children?: MenuProps[];
  // home页tab
  tabChildren?: MenuProps[];
  // 操作页
  operationChildren?: MenuProps[];
  isLeaf?: boolean;
  [key: string]: any;
}

type ContentTypeProps = 'addMenu' | 'addGroup' | 'editMenu' | 'editGroup' | null;

export type MenuStoreState = {
  menuList?: MenuProps[];
  menuTree?: MenuProps[];
  contentType?: ContentTypeProps;
  selectNode?: MenuProps | null;
  loading?: boolean;
};

export type MenuStoreActions = {
  requestMenu: () => Promise<void>;
  setContentType: (type: MenuStoreState['contentType']) => void;
  setSelectNode: (type: MenuStoreState['selectNode']) => void;
  setMenuInfo: (menuTree: MenuStoreState['menuTree'], menuList: MenuProps[]) => void;
};

export type MenuStoreProps = MenuStoreState & MenuStoreActions;
