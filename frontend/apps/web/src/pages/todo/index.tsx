import { useEffect, useMemo, useState } from 'react';
import { Task } from '@carbon/icons-react';
import { useNavigate } from 'react-router';
import { Radio, Space, Select, Button, Spin, Input, Checkbox } from 'antd';
import type { PaginationProps } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { getModuleList, todoPageList } from '@/apis/inter-api/todo';
import { useTranslate, usePagination, useMediaSize } from '@/hooks';
import { useActivate } from '@/contexts/tabs-lifecycle-context.ts';
import styles from './index.module.scss';
import { AuthWrapper } from '@/components/auth';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import ProTable from '@/components/pro-table';
import { formatTimestamp } from '@/utils/format';
const Todo = () => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const [moduleList, setModuleList] = useState<any>();
  const [moduleListName, setModuleListName] = useState<string>();
  const [radioValue, setRadioValue] = useState<any>('todo');
  const [taskValue, setTaskValue] = useState<string>();
  const [todoTotal, setTodoTotal] = useState<number>();
  const [checkboxValue, setCheckboxValue] = useState<boolean>(false);
  const { isH5 } = useMediaSize();
  const { pagination, data, loading, refreshRequest, setSearchParams } = usePagination({
    initPageSize: 20,
    defaultParams: {
      status: 0,
    },
    fetchApi: todoPageList,
    onSuccessCallback: ({ total }: any) => {
      if (radioValue === 'todo') {
        setTodoTotal(total);
      }
    },
  });
  const _moduleListName = useMemo(() => {
    const lastDashIndex = moduleListName?.lastIndexOf('-');
    return moduleListName?.slice?.(0, lastDashIndex) ?? '';
  }, [moduleListName]);
  const showTotal: PaginationProps['showTotal'] = (total) =>
    isH5 ? null : `${formatMessage('common.total')}  ${total}  ${formatMessage('common.items')}`;

  useEffect(() => {
    getModuleList().then((res: any) => {
      setModuleList(
        res?.map((item: any, index: number) => ({ ...item, moduleCodeAndName: `${item.moduleCode}-${index}` })) ?? []
      );
    });
  }, []);
  useActivate(() => {
    refreshRequest?.();
  });

  const columns: any = useMemo(() => {
    return [
      {
        title: () => formatMessage('common.origin'),
        dataIndex: 'moduleName',
        key: 'moduleName',
        width: '10%',
      },
      {
        title: () => formatMessage('common.task'),
        dataIndex: 'todoMsg',
        key: 'todoMsg',
        width: '40%',
        render: (item: any) => {
          console.log(item);
          return item?.replace(/\$\{(.*?)\}/, (_: any, p1: any) => {
            return formatTimestamp(Number(p1));
          });
        },
      },
      {
        title: () => formatMessage('common.creationTime'),
        dataIndex: 'createAt',
        width: '14%',
        key: 'createAt',
        render: (item: any) => formatTimestamp(item),
      },
      ...(radioValue === 'completed'
        ? [
            {
              title: () => formatMessage('common.processingTime'),
              dataIndex: 'handlerTime',
              width: '14%',
              key: 'handlerTime',
              render: (item: any) => formatTimestamp(item),
            },
          ]
        : []),
      ...(radioValue === 'todo'
        ? [
            {
              title: () => formatMessage('common.ToDoPerson'),
              dataIndex: 'username',
              width: '10%',
              key: 'username',
            },
          ]
        : [
            {
              title: () => formatMessage('common.completedPerson'),
              dataIndex: 'handlerUsername',
              width: '10%',
              key: 'handlerUsername',
            },
            {
              title: () => formatMessage('common.recipient'),
              dataIndex: 'username',
              width: '10%',
              key: 'username',
            },
          ]),
      ...(radioValue === 'todo'
        ? [
            {
              title: () => formatMessage('common.operation'),
              key: 'Operation',
              fixed: 'right',
              width: '10%',
              dataIndex: 'Operation',
              render: (_: any, e: any) => (
                <AuthWrapper>
                  <Space size="middle">
                    <Button
                      type="text"
                      style={{
                        color: 'var(--supos-theme-color)',
                        fontSize: 13,
                        whiteSpace: 'normal',
                        wordWrap: 'break-word',
                      }}
                      disabled={e?.status || !e.link}
                      onClick={() => {
                        if (e.moduleCode === 'system.module.alarm') {
                          navigate('/Alert', { state: { businessId: e.businessId } });
                        } else {
                          if (e.link) {
                            window.location.href = e.link;
                          }
                        }
                      }}
                    >
                      {formatMessage('common.processTask')}
                    </Button>
                  </Space>
                </AuthWrapper>
              ),
            },
          ]
        : []),
    ];
  }, [radioValue]);
  const handleSearchClick = () => {
    setSearchParams({
      myTodo: radioValue === 'completed' && checkboxValue ? true : null,
      moduleCode: _moduleListName,
      status: radioValue === 'todo' ? 0 : 1,
      todoMsg: taskValue,
    });
  };
  const handleRadioChange = (e: any) => {
    switch (e.target.value) {
      case 'todo':
        setSearchParams({ status: 0, moduleCode: _moduleListName, todoMsg: taskValue });
        break;
      case 'completed':
        setSearchParams({ status: 1, moduleCode: _moduleListName, todoMsg: taskValue });
        break;
      default:
        break;
    }
  };
  const checkboxChange = (e: any) => {
    setCheckboxValue(e.target.checked);
  };
  const handlePressEnter = (e: any) => {
    e.preventDefault();
    setSearchParams({
      myTodo: radioValue === 'completed' && checkboxValue ? true : null,
      moduleCode: _moduleListName,
      status: radioValue === 'todo' ? 0 : 1,
      todoMsg: e.target.value,
    });
  };
  return (
    <Spin spinning={loading}>
      <ComLayout className={styles['TodoTable']}>
        <ComContent
          hasBack={false}
          title={
            <>
              <Task />
              <span style={{ marginLeft: 10 }}>{formatMessage('common.taskCenter')}</span>
            </>
          }
          style={{
            padding: '0 20px',
          }}
        >
          <div className={styles['filter-container']}>
            <div className={styles['radio-group']}>
              <Radio.Group
                onChange={(e) => {
                  handleRadioChange(e);
                  setRadioValue(e.target.value);
                }}
                defaultValue="todo"
              >
                <Radio.Button value="todo">{`${formatMessage('common.todo')}(${todoTotal ?? 0})`}</Radio.Button>
                <Radio.Button value="completed">{formatMessage('common.completed')}</Radio.Button>
              </Radio.Group>
            </div>
            <div className={styles['right-filter']}>
              {radioValue === 'completed' ? (
                <Checkbox onChange={checkboxChange}>{formatMessage('common.myCompletedList')}</Checkbox>
              ) : (
                ''
              )}
              <div>
                <span className={styles['filter-label-task']}>{formatMessage('common.task')}</span>
                <Input
                  onPressEnter={handlePressEnter}
                  allowClear
                  className={styles['filter-input']}
                  onChange={(e) => setTaskValue(e.target.value)}
                />
              </div>
              <div>
                <span className={styles['filter-label-origin']}>{formatMessage('common.origin')}</span>
                <Select
                  className={styles['filter-select']}
                  onChange={(e: string) => setModuleListName(e)}
                  fieldNames={{ label: 'moduleName', value: 'moduleCodeAndName' }}
                  options={moduleList}
                  allowClear
                />
              </div>
              <Button
                className={styles['search-btn']}
                onClick={handleSearchClick}
                type="primary"
                icon={<SearchOutlined />}
              >
                {formatMessage('common.search')}
              </Button>
            </div>
          </div>
          <ProTable
            scroll={{ x: 'max-content', y: data.length > 0 ? 'calc(100vh - 280px)' : undefined }}
            columns={columns}
            key={radioValue}
            dataSource={data as any}
            resizeable
            pagination={{
              total: pagination?.total,
              showTotal: showTotal,
              style: { display: 'flex', justifyContent: 'flex-end', padding: '10px 0' },
              pageSize: pagination?.pageSize || 20,
              current: pagination?.page,
              showQuickJumper: true,
              pageSizeOptions: pagination?.pageSizes,
              onChange: pagination.onChange,
              onShowSizeChange: (current, size) => {
                pagination.onChange({ page: current, pageSize: size });
              },
            }}
          />
        </ComContent>
      </ComLayout>
    </Spin>
  );
};
export default Todo;
