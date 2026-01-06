import { useState, useRef, type ReactNode, type FC, useEffect } from 'react';
import { Breadcrumb, Dropdown } from 'antd';
import type { BreadcrumbProps } from 'antd';
import { useDeepCompareEffect, useSize } from 'ahooks';
import type { ItemType } from 'antd/es/menu/interface';
import './index.scss';

export interface ComBreadcrumbProps extends BreadcrumbProps {
  // 尾部标签
  addonAfter?: ReactNode;
  // 是否开启响应式
  responsive?: boolean;
}

interface WidthData {
  [index: number]: number;
}

const ComBreadcrumb: FC<ComBreadcrumbProps> = ({ items = [], addonAfter, responsive = true, ...restProps }) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const shadowRef = useRef<HTMLDivElement>(null);
  const addonAfterRef = useRef<HTMLDivElement>(null);
  // 拿到原始的items长度
  const [widthData, setWidthData] = useState<WidthData>({});
  // 真实渲染
  const [renderItems, setRenderItems] = useState<typeof items>([]);
  // 是否在计算原始的长度
  const [measuring, setMeasuring] = useState(false);

  const containerSize = useSize(containerRef);
  const addonAfterSize = useSize(addonAfterRef);

  // 拿到原始的items长度
  useDeepCompareEffect(() => {
    if (!responsive) {
      setRenderItems(items);
      return;
    }
    setMeasuring(true);
    // 设置容器宽度监听
    const updateWidth = () => {
      if (shadowRef.current) {
        const olElement = shadowRef.current?.querySelector('ol');
        if (!olElement) return;
        const breadcrumbChildren = olElement?.children || [];
        if (!breadcrumbChildren?.length) return;
        const _w: WidthData = {};
        for (let i = 0; i < (olElement?.childElementCount ?? 1); i++) {
          const child = breadcrumbChildren[i];
          const { width } = child.getBoundingClientRect();
          _w[i] = width;
        }
        setWidthData(_w);
        setMeasuring(false);
      }
    };
    const resizeObserver = new ResizeObserver(updateWidth);
    if (shadowRef.current) {
      resizeObserver.observe(shadowRef.current);
    }
    return () => {
      resizeObserver.disconnect();
    };
  }, [items, responsive]);

  // 计算真实items
  useEffect(() => {
    if (!responsive) return;
    if (!containerSize?.width || Object.keys(widthData).length === 0) return;

    // 如果items只有一个，直接显示
    if (items.length <= 1) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setRenderItems(items);
      return;
    }

    // 分离项目和分隔符宽度
    const itemWidths: number[] = [];
    const separatorWidths: number[] = [];

    Object.entries(widthData).forEach(([indexStr, width]) => {
      const index = Number(indexStr);
      // 偶数索引是项目，奇数索引是分隔符
      if (index % 2 === 0) {
        itemWidths[index / 2] = width;
      } else {
        separatorWidths[Math.floor(index / 2)] = width;
      }
    });

    const _containerWidth = containerSize.width - (addonAfterSize?.width ?? 0);

    const lastIndex = items.length - 1;
    const firstItem = items[0];
    const lastItem = items[lastIndex];

    // 计算所有分隔符宽度总和（除了最后一个分隔符）
    const totalSeparatorWidth = separatorWidths.slice(0, -1).reduce((sum, width) => sum + (width || 8), 0);

    let availableWidth = _containerWidth;
    const newItems = [];

    availableWidth -= itemWidths[0];
    newItems.push(firstItem);
    availableWidth -= itemWidths[lastIndex];
    availableWidth -= totalSeparatorWidth + 40; // 预留省略号按钮宽度

    for (let i = 1; i < lastIndex; i++) {
      // 计算当前项和其前一个分隔符的宽度
      const prevSeparatorWidth = separatorWidths[i - 1] || 8;
      const neededWidth = itemWidths[i] + prevSeparatorWidth;

      if (availableWidth >= neededWidth) {
        newItems.push(items[i]);
        availableWidth -= neededWidth;
      } else {
        const _items = items.slice(i, lastIndex).map((item, idx) => ({
          key: idx,
          label: item.title,
          onClick: item.onClick,
        })) as ItemType[];
        // 添加折叠菜单
        newItems.push({
          title: (
            <Dropdown
              menu={{
                items: _items,
              }}
            >
              <span
                style={{ cursor: 'pointer', display: 'inline-block', maxWidth: 40, color: 'var(--supos-theme-color)' }}
              >
                ...
              </span>
            </Dropdown>
          ),
        });
        break;
      }
    }

    newItems.push(lastItem);
    setRenderItems(newItems);
  }, [widthData, containerSize?.width, items, addonAfterSize?.width]);

  return (
    <div ref={containerRef} className={responsive ? 'com-breadcrumb' : ''}>
      <div ref={shadowRef} style={{ position: 'absolute', zIndex: -99999, whiteSpace: 'nowrap', visibility: 'hidden' }}>
        {measuring && <Breadcrumb items={items} {...restProps} />}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', width: '100%' }}>
        <Breadcrumb items={renderItems} {...restProps} />
        {addonAfter && (
          <div style={{ paddingLeft: 8, flexShrink: 0 }} ref={addonAfterRef}>
            {addonAfter}
          </div>
        )}
      </div>
    </div>
  );
};

export default ComBreadcrumb;
