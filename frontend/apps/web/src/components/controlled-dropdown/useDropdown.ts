import {
  useState,
  useRef,
  useEffect,
  type SetStateAction,
  type Dispatch,
  type RefObject,
  type MouseEvent,
} from 'react';
import type { ItemType } from 'antd/es/menu/interface';

/**
 * 位置接口定义
 */
interface Position {
  x: number;
  y: number;
}

/**
 * 自定义下拉菜单Hook返回值接口
 */
export interface UseDropdownResult {
  open: boolean;
  menuItems: ItemType[];
  triggerRef: RefObject<HTMLSpanElement>;
  showDropdown: (e: MouseEvent, items: ItemType[]) => void;
  hideDropdown: () => void;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const defaultPosition = { x: 0, y: 0 };

/**
 * 自定义下拉菜单Hook
 * 实现在点击或右键位置显示下拉菜单
 * @returns {UseDropdownResult} 下拉菜单相关的状态和方法
 */
const useDropdown = (): UseDropdownResult => {
  // 下拉菜单显示状态
  const [open, setOpen] = useState<boolean>(false);
  // 下拉菜单位置
  const [position, setPosition] = useState<Position>(defaultPosition);
  // 下拉菜单项
  const [menuItems, setMenuItems] = useState<ItemType[]>([]);
  // 虚拟触发元素引用
  const triggerRef = useRef<HTMLSpanElement>(null);

  /**
   * 同步虚拟触发元素位置
   * 当位置或显示状态变化时更新
   */
  useEffect(() => {
    if (triggerRef.current && open) {
      // 设置绝对位置，考虑页面滚动
      triggerRef.current.style.position = 'fixed';
      // 相对于视口的，需要减去滚动偏移
      const scrollX = window.scrollX || document.documentElement.scrollLeft;
      const scrollY = window.scrollY || document.documentElement.scrollTop;
      triggerRef.current.style.top = `${position.y - scrollY}px`;
      triggerRef.current.style.left = `${position.x - scrollX}px`;
      triggerRef.current.style.width = '1px';
      triggerRef.current.style.height = '1px';
    }
  }, [position, open]);

  /**
   * 左键点击触发下拉菜单
   * @param {React.MouseEvent} e - 鼠标事件
   * @param {ItemType[]} items - 菜单项
   */
  const showDropdown = (e: MouseEvent, items: ItemType[]): void => {
    e.preventDefault();
    // 使用 clientX/Y 获取相对于视口的位置，因为使用了fixed定位，所以不需要考虑滚动偏移
    setPosition({ x: e.clientX, y: e.clientY });
    setMenuItems(items);
    setOpen(true);
  };

  /**
   * 关闭下拉菜单
   */
  const hideDropdown = (): void => {
    setOpen(false);
  };

  return {
    open,
    menuItems,
    triggerRef,
    showDropdown,
    hideDropdown,
    setOpen,
  };
};

export default useDropdown;
