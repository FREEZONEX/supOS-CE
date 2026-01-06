import { type ReactNode, useRef, useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { Tree, type TreeProps } from 'antd';
import './index.scss';
import { loadIbmFont } from '@/utils/uns';

export interface ComTreeProps extends TreeProps {
  notDataContent?: ReactNode;
  height?: number;
  ibmStyle?: boolean;
}
export interface ComTreeRef {
  scrollTo?: (options: { key: string | number; align?: 'top' | 'bottom' | 'auto'; offset?: number }) => void;
}

const ComTree = forwardRef<ComTreeRef, ComTreeProps>((props, ref) => {
  const { notDataContent, height, ibmStyle = true, ...rest } = props;

  const [treeHeight, setTreeHeight] = useState(height); //树的高度，开启虚拟滚动

  // 创建一个 ref 来引用 treeWrap 元素
  const treeWrapRef = useRef<HTMLDivElement | null>(null);
  const treeRef: any = useRef(null);

  useImperativeHandle(ref, () => ({
    scrollTo: treeRef.current?.scrollTo,
  }));

  useEffect(() => {
    loadIbmFont();
    // 如果没有找到目标节点，则停止执行
    if (!treeWrapRef.current) {
      console.error('未找到类名为 com-tree-wrap 的元素');
      return;
    }

    // 使用 ResizeObserver 监控元素尺寸变化
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        if (entry.contentRect.height > 0) {
          setTreeHeight(entry.contentRect.height);
        }
      }
    });

    // 开始观察目标节点
    resizeObserver.observe(treeWrapRef.current);

    // 清理函数：当组件卸载或依赖项改变时取消观察
    return () => {
      if (treeWrapRef.current) {
        resizeObserver.unobserve(treeWrapRef.current);
      }
    };
  }, []);

  return (
    <div className="com-tree-wrap" ref={treeWrapRef}>
      <div className={`com-tree-box ${ibmStyle ? 'ibm-style' : ''}`}>
        {notDataContent}
        <Tree {...rest} ref={treeRef} height={treeHeight} />
      </div>
    </div>
  );
});
export default ComTree;
