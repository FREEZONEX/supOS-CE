import { type CSSProperties, type FC, type Key, useCallback, useEffect, useRef, useState } from 'react';
import { Button, Checkbox, type CheckboxOptionType, Divider, Flex, Popover, Tag } from 'antd';
import { Filter } from '@carbon/icons-react';
import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import usePropsValue from '@/hooks/usePropsValue.ts';
import './index.scss';
import useTranslate from '@/hooks/useTranslate.ts';

export interface TagFilterProps {
  value?: Key[];
  onChange?: (value: Key[], options?: any[]) => void;
  options?: CheckboxOptionType[];
  defaultValue?: Key[];
  style?: CSSProperties;
  showTag?: boolean;
  showNumber?: boolean;
}
const Index: FC<TagFilterProps> = ({
  value,
  onChange,
  showTag = false,
  showNumber = true,
  options = [],
  defaultValue,
  style,
}) => {
  const formatMessage = useTranslate();
  const [v, setV] = usePropsValue<Key[]>({
    value,
    defaultValue,
  });

  const onValueChange = (values: any[]) => {
    setV(values);
    onChange?.(
      values,
      values.map((v) => options?.find((f) => f.value === v))
    );
  };

  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);

  const checkArrows = useCallback(() => {
    const element = scrollContainerRef.current;
    if (element) {
      const isOverflowing = element.scrollWidth > element.clientWidth;
      setCanScrollLeft(isOverflowing && element.scrollLeft > 0);
      setCanScrollRight(isOverflowing && element.scrollLeft < element.scrollWidth - element.clientWidth - 1);
    }
  }, []);

  useEffect(() => {
    if (!showTag) return;
    const element = scrollContainerRef.current;
    if (element) {
      checkArrows();
      element.addEventListener('scroll', checkArrows);
      const resizeObserver = new ResizeObserver(checkArrows);
      resizeObserver.observe(element);
      return () => {
        element.removeEventListener('scroll', checkArrows);
        resizeObserver.unobserve(element);
      };
    }
  }, [v, checkArrows, showTag]);

  useEffect(() => {
    // 处理options变化，过滤掉v中不在options里的值
    if (v && v.length > 0) {
      const validValues = v.filter((value: string) => options.some((option) => option.value === value));
      console.log(validValues);
      if (validValues.length !== v.length) {
        onValueChange?.(validValues);
      }
    }
  }, [options]);

  const handleScroll = (direction: 'left' | 'right') => {
    const element = scrollContainerRef.current;
    if (element) {
      const scrollAmount = element.clientWidth * 0.8;
      element.scrollBy({
        left: direction === 'left' ? -scrollAmount : scrollAmount,
        behavior: 'smooth',
      });
    }
  };
  const allOptionsValues = options?.map((o) => o.value) || [];
  const isAllSelected = v && v.length > 0 && v.length === allOptionsValues.length;
  const isIndeterminate = v && v.length > 0 && v.length < allOptionsValues.length;
  const popoverContent = (
    <div>
      <Checkbox
        indeterminate={isIndeterminate}
        onChange={(e) => {
          onValueChange(e.target.checked ? allOptionsValues : []);
        }}
        checked={!!isAllSelected}
      >
        {formatMessage('common.allSelect')}
      </Checkbox>
      <Divider style={{ margin: '8px 0', background: '#c6c6c6' }} />
      <Checkbox.Group
        value={v}
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 8,
        }}
        onChange={onValueChange}
        options={options}
      />
    </div>
  );
  return (
    <Flex className="com-tag-filter" align="center" gap={4} style={{ overflow: 'hidden', ...style }}>
      <Popover arrow={false} placement="bottomLeft" title="" content={popoverContent} trigger="hover">
        <Button
          icon={
            <Flex align="center">
              <Filter />
            </Flex>
          }
          size="small"
          style={{ flexShrink: 0, background: 'var(--supos-switchwrap-bg-color)' }}
          color="default"
          variant="filled"
        />
      </Popover>
      {showNumber && (
        <Tag style={{ backgroundColor: 'var(--supos-table-head-color)', flexShrink: 0 }}>{v?.length ?? 0}</Tag>
      )}
      {showTag && canScrollLeft && (
        <Button
          size="small"
          shape="circle"
          icon={<LeftOutlined />}
          onClick={() => handleScroll('left')}
          style={{ flexShrink: 0 }}
        />
      )}
      {showTag && (
        <div
          ref={scrollContainerRef}
          style={{
            flex: 1,
            display: 'flex',
            alignItems: 'center',
            overflow: 'hidden',
            whiteSpace: 'nowrap',
          }}
          onWheel={(e) => {
            if (scrollContainerRef.current) {
              scrollContainerRef.current.scrollLeft += e.deltaY; // 横向滚动
            }
          }}
        >
          {(v || [])?.map((id: Key) => {
            const item = options?.find((f) => f.value === id);
            return (
              <Tag
                closable
                style={{ backgroundColor: 'var(--supos-bg-color)', flexShrink: 0 }}
                key={id}
                onClose={() => {
                  onValueChange(v.filter((value: Key) => value !== id));
                }}
              >
                {item?.label}
              </Tag>
            );
          })}
        </div>
      )}
      {showTag && canScrollRight && (
        <Button
          size="small"
          shape="circle"
          icon={<RightOutlined />}
          onClick={() => handleScroll('right')}
          style={{ flexShrink: 0 }}
        />
      )}
    </Flex>
  );
};

export default Index;
