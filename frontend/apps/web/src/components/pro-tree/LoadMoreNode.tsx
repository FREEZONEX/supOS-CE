import { type FC, useEffect, useRef } from 'react';
import type { DataNodeProps } from './types';

// 在ProTree组件外部定义LoadMoreNode组件
interface LoadMoreNodeProps {
  node: DataNodeProps;
  loadMoreData: (node: DataNodeProps) => void;
}
const LoadMoreNode: FC<LoadMoreNodeProps> = ({ node, loadMoreData }) => {
  const loadMoreRef = useRef<HTMLSpanElement>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);
  const loadedRef = useRef<boolean>(false);

  useEffect(() => {
    loadedRef.current = false;
    // 清理旧的observer
    if (observerRef.current) {
      observerRef.current.disconnect();
    }
    // 创建新的observer
    observerRef.current = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !loadedRef.current) {
          console.log(`加载更多节点 ${node.key} 可见，触发加载`);
          loadMoreData(node);
          loadedRef.current = true;
          observerRef.current?.disconnect();
        }
      },
      {
        threshold: 0.1, // 降低阈值，只需要10%可见就触发
        rootMargin: '100px 0px', // 增加前后100px的检测区域，提前触发
      }
    );

    if (loadMoreRef.current) {
      observerRef.current.observe(loadMoreRef.current);
    }

    // 添加一个备用方案：如果节点存在但没有被观察到，在短暂延迟后也触发加载
    const timer = setTimeout(() => {
      if (!loadedRef.current && loadMoreRef.current) {
        console.log(`加载更多节点 ${node.key} 备用触发`);
        loadMoreData?.(node);
        loadedRef.current = true;
        observerRef.current?.disconnect();
      }
    }, 300);

    return () => {
      clearTimeout(timer);
      if (observerRef.current) {
        observerRef.current?.disconnect();
        observerRef.current = null;
      }
    };
  }, [node.key]); // 只依赖key，减少不必要的effect执行

  return <span ref={loadMoreRef}>{node.title as string}</span>;
};

export default LoadMoreNode;
