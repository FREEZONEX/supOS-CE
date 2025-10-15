import { useState, useRef, useEffect, memo, type FC, useMemo } from 'react';
import { Button, Dropdown, Table } from 'antd';
import type { TableColumnsType } from 'antd';
import classNames from 'classnames';
import useTranslate from '@/hooks/useTranslate.ts';
import { useThemeStore } from '@/stores/theme-store.ts';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { EllipsisOutlined } from '@ant-design/icons';
import { hasPermission } from '@/utils/auth.ts';
import { commonLabelRender } from '../operation-buttons/utils.tsx';
import expandIcon from '@/assets/icons/expand.svg';
import collapseIcon from '@/assets/icons/collapse.svg';
import { ResizableTitle } from './ResizableTitle.tsx';
import type { ATableProps } from './types.ts';
import './index.scss';
import { Pin, PinFilled } from '@carbon/icons-react';
import ComButton from '../com-button';

const colorObj: any = {
  blue: {
    light: '#E8F1FF',
    dark: '#061833',
  },
  chartreuse: {
    light: '#F0FBD2',
    dark: '#242F06',
  },
};
const ProTable: FC<ATableProps> = ({
  resizeable,
  columns,
  components,
  scroll,
  className,
  hiddenEmpty,
  pagination,
  fixedPosition,
  showExpand,
  wrapperStyle,
  ...restProps
}) => {
  const formatMessage = useTranslate();
  const { theme, primaryColor } = useThemeStore((state) => ({
    theme: state.theme,
    primaryColor: state.primaryColor,
  }));
  const selectBgColor = colorObj?.[primaryColor]?.[theme];
  const [resizeColumns, setResizeColumns] = useState<TableColumnsType>(columns);
  const tableWrapRef = useRef<HTMLDivElement>(null);
  const containerWidthRef = useRef(0);

  const [isExpanded, setIsExpanded] = useState<boolean>(false);
  const [showAllColumns, setShowAllColumns] = useState<boolean>(true);

  const newPagination = pagination
    ? {
        total: pagination?.total,
        showTotal: (total: number) => `${formatMessage('common.total')}  ${total}  ${formatMessage('common.items')}`,
        style: { display: 'flex', justifyContent: 'flex-end', padding: '10px 0' },
        pageSize: 20,
        showQuickJumper: true,
        ...pagination,
      }
    : pagination;

  // 计算有效列宽总和（排除固定列的潜在样式误差）
  const calculateEffectiveWidth = (cols: TableColumnsType) => {
    return cols.reduce((sum, col) => {
      let width = typeof col.width === 'number' ? col.width : 0;
      // 修正固定列的边框误差
      if (col.fixed) width += 2; // 补偿Ant Design的固定列边框
      return sum + width;
    }, 0);
  };

  // 智能选择目标列
  const findFlexColumn = (cols: TableColumnsType) => {
    // 从右向左查找第一个非固定列
    for (let i = cols.length - 1; i >= 0; i--) {
      if (!cols[i].fixed) return i;
    }
    // 全固定列时选择倒数第二列
    return Math.max(0, cols.length - 2);
  };

  // 动态平衡列宽
  const balanceColumns = (cols: TableColumnsType, changedIndex?: number) => {
    const containerWidth = containerWidthRef.current;
    if (!containerWidth) return cols;
    // 新增：预留滚动条宽度
    const SCROLLBAR_WIDTH = 17;
    const totalWidth = calculateEffectiveWidth(cols);
    const delta = containerWidth - totalWidth - SCROLLBAR_WIDTH - (restProps?.rowSelection ? 35 : 0);

    // 当列宽总和不足时自动扩展
    if (delta > 0) {
      const targetIndex = findFlexColumn(cols);
      const newColumns = [...cols];
      newColumns[targetIndex] = {
        ...newColumns[targetIndex],
        width: (newColumns[targetIndex].width as number) + delta,
      };
      return newColumns;
    }

    // 处理主动缩小时同步调整固定列
    if (changedIndex !== undefined && cols[changedIndex]?.fixed) {
      const newColumns = [...cols];
      const nextIndex = changedIndex + 1;
      if (nextIndex < cols.length && !cols[nextIndex].fixed) {
        newColumns[nextIndex] = {
          ...newColumns[nextIndex],
          width: (newColumns[nextIndex].width as number) - delta,
        };
      }
      return newColumns;
    }
    return cols;
  };

  const handleResize = (index: number) => (width?: number) => {
    if (!width || !tableWrapRef.current) return;

    // 1：更新当前列宽
    const newColumns = [...resizeColumns];
    newColumns[index] = { ...newColumns[index], width };

    // 2：智能平衡列宽
    const balancedColumns = balanceColumns(newColumns, index);

    // 3：同步固定列样式
    balancedColumns.forEach((col) => {
      if (col.fixed) {
        const selector = col.fixed === 'right' ? '.ant-table-cell-fix-right' : '.ant-table-cell-fix-left';
        document.querySelectorAll(selector).forEach((el) => {
          const cell = el as HTMLElement;
          if (cell.textContent === col.title?.toString()) {
            cell.style.width = `${col.width}px`;
          }
        });
      }
    });

    setResizeColumns(balancedColumns);
  };

  useEffect(() => {
    const observer = new ResizeObserver((entries) => {
      if (entries[0]) {
        containerWidthRef.current = entries[0].contentRect.width;
        setResizeColumns((prev) => balanceColumns(prev));
      }
    });

    if (tableWrapRef.current) {
      observer.observe(tableWrapRef.current);
      containerWidthRef.current = tableWrapRef.current.clientWidth;
    }
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!containerWidthRef.current) return;
    const newColumns = columns.map((col) => {
      return typeof col.width === 'string'
        ? {
            ...col,
            width: col.width.includes('%') ? containerWidthRef.current * (parseFloat(col.width) / 100) : col.width,
          }
        : col;
    });
    setResizeColumns(() => balanceColumns(newColumns));
  }, [columns]);

  const mergedColumns = resizeColumns.map<TableColumnsType[number]>((col, index) => ({
    ...col,
    onHeaderCell: (column: TableColumnsType[number]) => ({
      width: column.width,
      minWidth: column.minWidth,
      changeWidth: handleResize(index),
    }),
  }));

  const _classNames = classNames(className, 'pro-table', {
    'resizable-table': resizeable,
    'hidden-empty': hiddenEmpty && restProps?.dataSource?.length === 0,
    'fixed-pagination-bottom': fixedPosition,
  });

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };

  const handleExpandClick = () => {
    setShowAllColumns(!showAllColumns);
  };

  const changeShowColumns = (oldColumns: TableColumnsType) => {
    return !showExpand
      ? oldColumns
      : oldColumns.map((col) => {
          if ([true, false].includes(col?.hidden as boolean)) {
            col.hidden = showAllColumns;
          }
          return col;
        });
  };

  return (
    <div
      ref={tableWrapRef}
      className="pro-table-container"
      style={{
        '--supos-table-select-bg-color': selectBgColor,
        width: '100%',
        ...wrapperStyle,
      }}
      onMouseEnter={toggleExpanded}
      onMouseLeave={toggleExpanded}
    >
      {resizeable ? (
        <Table
          rowKey="id"
          size={'small'}
          {...restProps}
          className={_classNames}
          columns={changeShowColumns(mergedColumns)}
          pagination={newPagination}
          scroll={{
            x: 'max-content',
            ...scroll,
          }}
          components={{ ...components, header: { cell: ResizableTitle } }}
          tableLayout="fixed"
        />
      ) : (
        <Table
          rowKey="id"
          size={'small'}
          {...restProps}
          className={_classNames}
          components={components}
          scroll={scroll}
          columns={changeShowColumns(columns)}
          pagination={newPagination}
        />
      )}
      {showExpand && (
        <div
          className={`pro-table-expand-button ${isExpanded ? 'pro-table-expanded' : ''}`}
          onClick={handleExpandClick}
        >
          <img src={showAllColumns ? expandIcon : collapseIcon} alt="" />
        </div>
      )}
    </div>
  );
};

const withIntlTable = (WrappedTable: FC<ATableProps>) => {
  const IntlTableWrapper = memo(({ columns, operationOptions, pinOptions, ...restProps }: ATableProps) => {
    const lang = useI18nStore((state) => state.lang);
    const formatMessage = useTranslate();
    const _columns = useMemo(() => {
      if (pinOptions) {
        const { disabled, onClick, renderPinIcon, auth, ...restProps } = pinOptions;
        columns.unshift({
          title: ' ',
          dataIndex: 'pin',
          align: 'center',
          fixed: 'left',
          width: 40,
          render: (_: any, record: any) => {
            const isPin = renderPinIcon?.(record) ?? false;
            return (
              <ComButton
                title={isPin ? formatMessage('common.pin') : formatMessage('common.unPin')}
                auth={auth}
                disabled={disabled}
                onClick={() => onClick?.(record)}
                className={classNames('custom-pin', !isPin && 'custom-pin-fixed')}
                icon={isPin ? <Pin /> : <PinFilled />}
                size="small"
                type={'text'}
              />
            );
          },
          ...restProps,
        });
      }
      if (operationOptions) {
        // 通用操作项
        columns.push({
          title: () => formatMessage('common.operation'),
          width: 120,
          dataIndex: 'operation',
          align: 'left',
          fixed: 'right',
          ...operationOptions,
          render: (_: any, record: any, index: number) => {
            const contentRaw: any =
              operationOptions.render?.(record, index).filter((item: any) => {
                return item && (!item.auth || hasPermission(item.auth));
              }) || [];
            if (contentRaw?.length === 0) return null;
            return (
              <div className="custom-operation">
                <Dropdown
                  disabled={operationOptions.disabled}
                  menu={{
                    items: contentRaw.map((record: any) => {
                      const { key, label, icon, title, onClick, disabled, type, extra } = record;
                      return {
                        key,
                        label: commonLabelRender(record),
                        icon,
                        title: title ? title : typeof label === 'string' ? label : '',
                        onClick: type !== 'Popconfirm' && onClick,
                        disabled,
                        extra,
                      };
                    }),
                  }}
                >
                  <Button type="text" icon={<EllipsisOutlined />} style={{ height: 21 }} size="small" />
                </Dropdown>
              </div>
            );
          },
        });
      }
      return columns.map((i: any) => {
        const type = typeof i.title;
        if (i.titleIntlId) {
          return {
            ...i,
            title: () => formatMessage(i.titleIntlId),
          };
        } else if (type === 'function') {
          const originalTitleFn: any = i.title;
          return {
            ...i,
            title: (params: any) => originalTitleFn({ ...params, formatMessage }),
          };
        } else {
          return i;
        }
      });
    }, [lang, columns, operationOptions]);
    return <WrappedTable {...restProps} columns={_columns} />;
  });
  return IntlTableWrapper;
};

export default withIntlTable(ProTable);
