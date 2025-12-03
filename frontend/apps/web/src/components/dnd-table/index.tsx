import { DndContext, type DragEndEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import { arrayMove, SortableContext, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import type { ATableProps } from '../pro-table/types.ts';
import ProTable from '../pro-table';

interface PropsTypes {
  tableConfig: ATableProps;
  onDragEnd: (arg0: any, arg1: any) => void;
  disabled?: boolean;
}

const Row: React.FC<Readonly<any>> = (props) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: props['data-row-key'],
  });

  const style: React.CSSProperties = {
    ...props.style,
    transform: CSS.Translate.toString(transform),
    transition,
    cursor: !props?.disabled ? 'move' : 'inherit',
    ...(isDragging ? { position: 'relative' } : {}),
  };

  return <tr {...props} ref={setNodeRef} style={style} {...attributes} {...listeners} />;
};

const DndTable = (props: PropsTypes) => {
  const { tableConfig, onDragEnd, disabled } = props;
  const { dataSource = [] } = tableConfig;
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 1,
      },
    })
  );
  const unitKey = typeof tableConfig?.rowKey === 'string' ? tableConfig?.rowKey : 'code';

  // 拖拽
  const handleDragEnd = ({ active, over }: DragEndEvent) => {
    if (active.id === over?.id) return;

    const newDataSource = [...dataSource];
    const activeIndex = newDataSource.findIndex((i) => i[unitKey] === active.id);
    const overIndex = newDataSource.findIndex((i) => i[unitKey] === over?.id);
    const currentId = newDataSource[activeIndex]?.id;

    let prevId, nextId;

    // 向上移动
    if (activeIndex > overIndex) {
      nextId = newDataSource[overIndex]?.id;
      prevId = overIndex === 0 ? undefined : newDataSource[overIndex - 1]?.id;
    } else {
      // 向下移动
      prevId = newDataSource[overIndex]?.id;
      nextId = overIndex === newDataSource.length - 1 ? undefined : newDataSource[overIndex + 1]?.id;
    }

    onDragEnd(arrayMove(newDataSource, activeIndex, overIndex), {
      prevId,
      nextId,
      currentId,
    });
  };
  return (
    <DndContext sensors={sensors} modifiers={[restrictToVerticalAxis]} onDragEnd={handleDragEnd}>
      <SortableContext
        disabled={disabled}
        items={dataSource.map((i) => i?.[unitKey])}
        strategy={verticalListSortingStrategy}
      >
        <ProTable
          {...tableConfig}
          columns={tableConfig.columns || []}
          components={{
            body: { row: (props: any) => <Row {...props} disabled={disabled} /> },
          }}
        />
      </SortableContext>
    </DndContext>
  );
};

export default DndTable;
