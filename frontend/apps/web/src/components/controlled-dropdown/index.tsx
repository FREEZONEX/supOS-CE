import { Dropdown, type DropdownProps } from 'antd';
import useDropdown, { type UseDropdownResult } from './useDropdown.ts';
import { type CSSProperties, forwardRef, useImperativeHandle } from 'react';
import cx from 'classnames';
import './index.scss';

// 组件属性类型
export interface ControlledDropdownProps extends Omit<DropdownProps, 'open' | 'menu' | 'trigger'> {
  triggerStyle?: CSSProperties;
}

// 组件引用类型
export interface ControlledDropdownRef {
  showDropdown: UseDropdownResult['showDropdown'];
}

const ControlledDropdown = forwardRef<ControlledDropdownRef, ControlledDropdownProps>((props, ref) => {
  const { triggerStyle, overlayClassName, onOpenChange, ...restProps } = props;
  const { menuItems, triggerRef, showDropdown, open, hideDropdown } = useDropdown();

  useImperativeHandle(ref, () => ({
    showDropdown,
  }));

  return (
    <Dropdown
      placement="bottomLeft"
      {...restProps}
      open={open}
      onOpenChange={(open, info) => {
        onOpenChange?.(open, info);
        return !open && hideDropdown();
      }}
      menu={{ items: menuItems, onClick: hideDropdown }}
      trigger={['click']}
      overlayClassName={cx('controlled-dropdown', overlayClassName)}
    >
      <span
        ref={triggerRef}
        style={{
          opacity: 0,
          pointerEvents: 'none',
          ...triggerStyle,
        }}
      />
    </Dropdown>
  );
});

export default ControlledDropdown;
