import { type FC, useState, useEffect, useMemo } from 'react';
import { Form, Select, Flex, Button, message, Input } from 'antd';
import { ChevronRight, ChevronLeft } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import { getSourceList, getColumnList } from '@/apis/chat2db';
import { ds2fs, json2fs } from '@/apis/inter-api/uns';
import FieldsFormList from '@/pages/uns/components/use-create-modal/components/FieldsFormList';
import TableSelect from '@/pages/uns/components/use-create-modal/components/file/table-select';

import type { FieldItem } from '@/pages/uns/types';
import ComRadio from '@/components/com-radio';
import { useBaseStore } from '@/stores/base';
import { isEmpty } from 'lodash-es';

const { TextArea } = Input;

export interface ReverseGenerationProps {
  types?: string[];
  // 只有json情况
  onlyJson?: boolean;
  require?: boolean;
}

const ReverseGeneration: FC<ReverseGenerationProps> = ({ types, onlyJson, require }) => {
  const form = Form.useFormInstance();
  const formatMessage = useTranslate();

  const source = Form.useWatch('source', form) || form.getFieldValue('source');
  const dataSource = Form.useWatch('dataSource', form) || form.getFieldValue('dataSource');
  const next = Form.useWatch('next', form) || form.getFieldValue('next');
  const jsonList = Form.useWatch('jsonList', form) || form.getFieldValue('jsonList');
  const jsonDataPath = Form.useWatch('jsonDataPath', form) || form.getFieldValue('jsonDataPath');
  const jsonData = Form.useWatch('jsonData', form) || form.getFieldValue('jsonData');
  const systemInfo = useBaseStore((state) => state.systemInfo);

  const [sourceList, setSourceList] = useState([]);
  const [hasMoreJson, setHasMoreJson] = useState(false);

  useEffect(() => {
    if (systemInfo?.containerMap?.chat2db) {
      getSourceList().then((res: any) => {
        if (res.success) {
          setSourceList(
            res?.data?.data?.filter((source: any) =>
              ['postgresql', 'sqlserver', 'mariadb', 'mysql'].includes(source.type?.toLowerCase())
            ) || []
          );
        }
      });
    }
  }, []);

  useEffect(() => {
    setHasMoreJson(jsonList?.length > 1);
  }, [jsonList]);

  const validatorJson = (_: any, value: any) => {
    if (!value && require) return Promise.reject(new Error(formatMessage('uns.pleaseEnterJSON')));
    try {
      if (!value && !require) {
        return Promise.resolve();
      }
      const jsonVal = JSON.parse(value);
      if (['[object Object]', '[object Array]'].includes(Object.prototype.toString.call(jsonVal))) {
        return Promise.resolve();
      } else {
        return Promise.reject(new Error(formatMessage('uns.errorInTheSyntaxOfTheJSON')));
      }
      // eslint-disable-next-line
    } catch (err) {
      return Promise.reject(new Error(formatMessage('uns.errorInTheSyntaxOfTheJSON')));
    }
  };

  const sourceMap: any = {
    json: ['jsonData'],
    dataSource: ['dataSource', 'table'],
  };
  const prevStep = () => {
    setHasMoreJson(false);
  };
  const nextStep = () => {
    form.validateFields(sourceMap[source]).then(async (values) => {
      if (source === 'json') {
        const jsonObj = JSON.parse(values.jsonData);

        if (isEmpty(jsonObj)) {
          form.setFieldsValue({ jsonList: [], fields: [], jsonDataPath: undefined, next: true });
          setHasMoreJson(false);
          return;
        }
        if (hasMoreJson) {
          form.setFieldsValue({ next: true });
        } else {
          const fields: any = await json2fs(jsonObj);
          fields.forEach((e: any) => {
            e.dataPath = e.dataPath || 'default';
          });
          if (fields.length === 1) {
            form.setFieldsValue({ jsonList: [], fields: fields[0]?.fields, jsonDataPath: undefined, next: true });
            setHasMoreJson(false);
          } else {
            form.setFieldsValue({ jsonList: fields, fields: fields[0]?.fields, jsonDataPath: fields[0]?.dataPath });
            setHasMoreJson(true);
          }
        }
      }
      if (source === 'dataSource') {
        const { table } = values;
        const tableParentIds = table?.value?.split('$分隔符$');
        const hasSchema = tableParentIds.length > 3;
        const dataSourceInfo: any = sourceList.find((item: any) => item.id === dataSource) || {};
        const params = {
          dataSourceId: dataSource,
          dataSourceName: dataSourceInfo.alias,
          databaseType: dataSourceInfo.type,
          databaseName: tableParentIds[0],
          schemaName: hasSchema ? tableParentIds[1] : undefined,
          tableName: tableParentIds[hasSchema ? 3 : 2],
          refresh: true,
          pageNo: 1,
        };
        const { data }: any = await getColumnList(params);
        if (data?.length) {
          const fields = await ds2fs({ databaseType: dataSourceInfo.type, fields: data });
          if (fields) {
            fields.forEach((e: FieldItem) => {
              if (e?.type?.toLowerCase() !== 'string') delete e.maxLen;
            });
            form.setFieldsValue({ fields, next: true });
          }
        } else {
          message.warning(formatMessage('uns.TableNoFieldsTip'));
          form.setFieldsValue({ fields: undefined });
        }
      }
    });
  };

  const back = () => {
    const fields = jsonList.find((item: any) => item.dataPath === jsonDataPath)?.fields;
    form.setFieldsValue({ fields: hasMoreJson ? fields : undefined, next: false });
  };

  const options = useMemo(() => {
    if (systemInfo?.containerMap?.chat2db) {
      return [
        { label: 'JSON', value: 'json' },
        { label: formatMessage('streams.dataSource'), value: 'dataSource' },
      ];
    } else {
      return [{ label: 'JSON', value: 'json' }];
    }
  }, [systemInfo?.containerMap?.chat2db]);

  const exampleJson = `{
    "Example": {
        "PathName": {
            "TopicName": [
                {
                    "attribute1": 1380,
                    "attribute2": 1440
                }
            ]
        }
    }
}`;

  return (
    <>
      <Form.Item name="next" hidden initialValue={false}>
        <Input />
      </Form.Item>
      <Form.Item name="jsonList" hidden initialValue={[]}>
        <Input />
      </Form.Item>
      <Form.Item
        name="source"
        hidden={onlyJson}
        initialValue={onlyJson ? 'json' : undefined}
        label={formatMessage('uns.source')}
        rules={[{ required: true }]}
      >
        <Select
          options={options}
          onChange={() => {
            form.setFieldsValue({
              fields: [{}],
              dataSource: undefined,
              table: undefined,
              jsonData: undefined,
              jsonList: [],
              jsonDataPath: undefined,
              next: false,
            });
          }}
        />
      </Form.Item>
      {source === 'json' && !next && !hasMoreJson && (
        <div style={{ position: 'relative', width: '100%' }}>
          <Form.Item
            name="jsonData"
            label=""
            wrapperCol={{ span: 24 }}
            rules={[{ required: true, validator: validatorJson }]}
            validateTrigger={['onBlur', 'onChange']}
          >
            <TextArea
              placeholder={exampleJson}
              rows={8}
              allowClear
              onKeyDownCapture={(e) => {
                if (e.ctrlKey && e.code === 'Enter') {
                  if (jsonData) return;
                  e.preventDefault();
                  form.setFieldValue('jsonData', exampleJson);
                }
              }}
              onKeyDown={(e) => {
                if (e.ctrlKey && e.key === 'Enter') {
                  e.preventDefault();
                }
              }}
            />
          </Form.Item>
          {!jsonData && (
            <span
              style={{
                position: 'absolute',
                top: 6,
                right: 12,
                fontSize: '12px',
                pointerEvents: 'none',
                zIndex: 10,
                color: '#c6c6c6',
              }}
            >
              {formatMessage('uns.ctrlPQuickApplyExample')}
            </span>
          )}
        </div>
      )}
      {source === 'dataSource' && !next && (
        <>
          <Form.Item name="dataSource" label={formatMessage('streams.dataSource')} rules={[{ required: true }]}>
            <Select
              onChange={() => {
                form.setFieldsValue({
                  table: undefined,
                });
              }}
              options={sourceList}
              fieldNames={{ label: 'alias', value: 'id' }}
            />
          </Form.Item>
          <Form.Item name="table" label={formatMessage('common.table')} rules={[{ required: true }]}>
            <TableSelect sourceList={sourceList} dataSource={dataSource} formatMessage={formatMessage} />
          </Form.Item>
        </>
      )}
      {hasMoreJson && !next && (
        <Form.Item name="jsonDataPath" label={formatMessage('uns.schemaGenerated')}>
          <ComRadio
            style={{ flexWrap: 'wrap' }}
            options={jsonList.map((item: any) => ({ label: item.dataPath, value: item.dataPath }))}
            onChange={(e) => {
              const fields = jsonList.find((item: any) => item.dataPath === e.target.value)?.fields;
              form.setFieldsValue({ fields });
            }}
          />
        </Form.Item>
      )}
      {(next || hasMoreJson) && (
        <FieldsFormList types={types} disabled={!next} showMainKey={next} showWrap={false} showMoreBtn={next} />
      )}
      {source && (
        <Flex justify="flex-end" style={{ marginTop: next || hasMoreJson ? '20px' : '' }} gap={10}>
          {!next && hasMoreJson && (
            <Button color="primary" variant="filled" size="small" icon={<ChevronLeft />} onClick={prevStep}>
              {formatMessage('common.prev')}
            </Button>
          )}
          {next ? (
            <Button color="primary" variant="filled" size="small" icon={<ChevronLeft />} onClick={back}>
              {formatMessage('common.back')}
            </Button>
          ) : (
            <Button
              color="primary"
              variant="filled"
              size="small"
              icon={<ChevronRight />}
              iconPosition="end"
              onClick={nextStep}
            >
              {formatMessage('common.next')}
            </Button>
          )}
        </Flex>
      )}
    </>
  );
};

export default ReverseGeneration;
