import { Children, cloneElement, isValidElement, type ReactNode } from 'react';

export function injectPropsToRouteNode(node: ReactNode, propsToInject: object): ReactNode {
  if (!isValidElement(node)) {
    return node;
  }
  if (node.props && node.props.match) {
    // 如果有子元素，则对每个子元素递归注入 props
    if (node.props.children) {
      const children = Children.map(node.props.children, (child) => {
        // 递归注入 props 到子元素
        return cloneElement(child, { ...child.props, ...propsToInject });
      });
      return cloneElement(node, { ...node.props, children });
    }
  }
  if (node.props && node.props.children) {
    const children = Children.map(node.props.children, (child) => injectPropsToRouteNode(child, propsToInject));
    return cloneElement(node, { ...node.props, children });
  }

  return node;
}
