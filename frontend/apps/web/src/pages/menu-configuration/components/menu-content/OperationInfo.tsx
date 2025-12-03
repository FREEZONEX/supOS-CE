import { useTranslate } from '@/hooks';
import { Divider, Flex, Form, Input } from 'antd';
import { AuthButton } from '@/components/auth';
import DndTable from '@/components/dnd-table';
import useNewOperation from './useNewOperation.tsx';
import { useMenuStore } from '@/pages/menu-configuration/store/menuStore.tsx';
import { ButtonPermission } from '@/common-types/button-permission.ts';

const OperationInfo = () => {
  const formatMessage = useTranslate();
  const form = Form.useFormInstance();
  const { selectNode } = useMenuStore((state) => ({
    selectNode: state.selectNode,
  }));
  const operationChildren = Form.useWatch('operationChildren');
  const { NewOperation, onNewOperationOpen } = useNewOperation({
    onSuccessBack: async (data: any, type: 'add' | 'edit') => {
      if (type === 'add') {
        form.setFieldValue('operationChildren', [
          ...operationChildren,
          { ...data, id: 'add_' + Math.random(), type: 3 },
        ]);
      } else if (type === 'edit') {
        form.setFieldValue(
          'operationChildren',
          operationChildren?.map((item: any) => {
            if (item.id === data.id) {
              return {
                ...item,
                ...data,
              };
            }
            return item;
          })
        );
      }
    },
  });
  const handleDragEnd = (newData: any[]) => {
    form.setFieldValue('operationChildren', newData);
  };
  return (
    <>
      {NewOperation}
      <Flex gap={8} align="center" justify="space-between">
        <span style={{ fontSize: 20, fontWeight: 500 }}>{formatMessage('MenuConfiguration.operationInfo')}</span>
        <AuthButton
          size="small"
          variant="outlined"
          color="primary"
          onClick={() => onNewOperationOpen()}
          auth={ButtonPermission['MenuConfiguration.addOperation']}
        >
          + {formatMessage('MenuConfiguration.addOperation')}
        </AuthButton>
      </Flex>
      <Divider style={{ backgroundColor: '#C6C6C6', margin: '16px 0' }} />
      <DndTable
        disabled={selectNode?.editEnable === false}
        tableConfig={{
          operationOptions: {
            disabled: selectNode?.editEnable === false,
            render: (record: any) => {
              return [
                {
                  type: 'Button',
                  key: 'edit',
                  auth: ButtonPermission['MenuConfiguration.editOperation'],
                  label: formatMessage('common.edit'),
                  disabled: selectNode?.editEnable === false,
                  onClick: () => {
                    onNewOperationOpen(record);
                  },
                },
                {
                  type: 'Popconfirm',
                  key: 'delete',
                  auth: ButtonPermission['MenuConfiguration.deleteOperation'],
                  label: formatMessage('common.delete'),
                  disabled: selectNode?.editEnable === false,
                  onClick: () => {
                    const operationChildren = form.getFieldValue('operationChildren');
                    const delOperationChildren = form.getFieldValue('delOperationChildren') || [];
                    const index = operationChildren.findIndex((item: any) => item.id === record?.id);
                    if (index !== -1) {
                      const [removedItem] = operationChildren.splice(index, 1);
                      form.setFieldsValue({
                        delOperationChildren: [...delOperationChildren, removedItem],
                        operationChildren,
                      });
                    }
                  },
                  popconfirm: {
                    title: formatMessage('common.deleteConfirm'),
                  },
                },
              ];
            },
          },
          wrapperStyle: { paddingRight: 16 },
          fixedPosition: true,
          rowKey: 'id',
          columns: [
            {
              title: () => formatMessage('MenuConfiguration.operationCode'),
              dataIndex: 'code',
              width: '30%',
              ellipsis: true,
            },
            {
              title: () => formatMessage('MenuConfiguration.operationName'),
              dataIndex: 'name',
              width: '20%',
              ellipsis: true,
            },
            {
              title: () => formatMessage('MenuConfiguration.operationDescription'),
              dataIndex: 'description',
              width: '30%',
              ellipsis: true,
            },
          ],
          dataSource: operationChildren,
          scroll: { y: 500 },
          resizeable: true,
          pagination: false,
          expandable: {
            childrenColumnName: '-children',
          },
        }}
        onDragEnd={handleDragEnd}
      />
      <Form.Item name="operationChildren" hidden>
        <Input />
      </Form.Item>
      <Form.Item name="delOperationChildren" hidden>
        <Input />
      </Form.Item>
    </>
  );
};

export default OperationInfo;
