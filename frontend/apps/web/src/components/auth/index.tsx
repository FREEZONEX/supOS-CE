import { cloneElement, type ComponentType, type FC, Children, type ReactNode, isValidElement } from 'react';
import { Button } from 'antd';
import { hasPermission } from '@/utils/auth';

const withAuth =
  <T extends { auth?: string | string[] }>(Component: ComponentType<T>) =>
  ({ auth, ...props }: T) => {
    if (auth && !hasPermission(auth)) {
      return null;
    }
    return <Component {...(props as T)} data-button-auth={auth} />;
  };

export const AuthWrapper: FC<{ children: ReactNode; auth?: string | string[] }> = ({ children, auth }) => {
  if (auth && !hasPermission(auth)) {
    return null;
  }
  return (
    <>
      {Children.map(children, (child) => {
        // 如果 child 是有效的 React 元素，则克隆并添加 'data-auth' 属性
        if (isValidElement(child)) {
          if (
            typeof child.type === 'object' &&
            child.type !== null &&
            'render' in child.type &&
            (child.type as { render?: { displayName?: string } }).render?.displayName === 'Tooltip'
          ) {
            return <div data-button-auth={auth}>{child}</div>;
          }
          return cloneElement(child, { ...child.props, 'data-button-auth': auth });
        }
        // 如果 child 不是有效的 React 元素，直接返回它
        return <span data-button-auth={auth}>{child}</span>;
      })}
    </>
  );
};

// 权限按钮
export const AuthButton = withAuth(Button);

export default withAuth;
