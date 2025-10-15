import { Children, cloneElement, isValidElement, type ReactNode } from 'react';

export function injectPropsToRouteNode(node: ReactNode, propsToInject: object): ReactNode {
  if (!isValidElement(node)) {
    return node;
  }
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  if (node.props && node.props.match) {
    // 如果有子元素，则对每个子元素递归注入 props
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    if (node.props.children) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-expect-error
      const children = Children.map(node.props.children, (child) => {
        // 递归注入 props 到子元素
        return cloneElement(child, { ...child.props, ...propsToInject });
      });
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-expect-error
      return cloneElement(node, { ...node.props, children });
    }
  }
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  if (node.props && node.props.children) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const children = Children.map(node.props.children, (child) => injectPropsToRouteNode(child, propsToInject));
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    return cloneElement(node, { ...node.props, children });
  }

  return node;
}
