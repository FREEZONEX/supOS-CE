import { type FC, useEffect, useState, type Dispatch, type SetStateAction, useRef } from 'react';
import { Form, Select, Divider } from 'antd';
import { useTranslate } from '@/hooks';
import JsonForm from './source-form/json';
import DataSource from './source-form/data-source';
import Collector from './source-form/collector';
import { getTypes, getEmptyFolder } from '@/apis/inter-api/uns';

import type { UnsTreeNode, InitTreeDataFnType } from '@/pages/uns/types';
import ComRadio from '@/components/com-radio';
import ProModal from '@/components/pro-modal';
import { useBaseStore } from '@/stores/base';

export interface ReverseModalProps {
  currentNode?: UnsTreeNode;
  reverserOpen: boolean;
  setReverserOpen: Dispatch<SetStateAction<boolean>>;
  initTreeData: InitTreeDataFnType;
}

const ReverseModal: FC<ReverseModalProps> = ({ reverserOpen, setReverserOpen, currentNode, initTreeData }) => {
  const [form] = Form.useForm();
  const formatMessage = useTranslate();

  const [types, setTypes] = useState<string[]>([]);
  const [fullScreen, setFullScreen] = useState<boolean>(false);
  const [folderList, setFolderList] = useState([]);
  const [temporaryNode, setTemporaryNode] = useState(currentNode);

  const attributeType = Form.useWatch('attributeType', form) || form.getFieldValue('attributeType');
  const source = Form.useWatch('source', form) || form.getFieldValue('source');
  const parentDataType = Form.useWatch('parentDataType', form) || form.getFieldValue('parentDataType');

  const jsonFormRef = useRef<any>(null);

  const {
    systemInfo: { enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));

  const getSourceForm = () => {
    switch (source) {
      case 'json':
        return (
          <JsonForm
            ref={jsonFormRef}
            formatMessage={formatMessage}
            types={types}
            currentNode={currentNode}
            close={close}
            fullScreen={fullScreen}
            initTreeData={initTreeData}
          />
        );
      case 'connect':
        return <DataSource formatMessage={formatMessage} close={close} />;
      case 'grpcGateway':
        return <Collector formatMessage={formatMessage} close={close} />;
      default:
        return null;
    }
  };

  const close = (refreshTree?: boolean) => {
    if (refreshTree) initTreeData({ reset: true });
    setReverserOpen(false);
  };

  useEffect(() => {
    getTypes()
      .then((res) => {
        setTypes(res || []);
      })
      .catch((err) => {
        console.log(err);
      });
    getEmptyFolder().then((res) => {
      setFolderList(res);
    });
  }, []);

  const sourceOptionsMap: { [key: number]: { label: string; value: string }[] } = {
    1: [{ label: 'JSON', value: 'json' }],
    2: [
      { label: formatMessage('streams.dataSource'), value: 'connect' },
      // { label: formatMessage('uns.grpcGateway'), value: 'grpcGateway' },
    ],
  };

  useEffect(() => {
    if (
      attributeType === 2 &&
      temporaryNode?.pathType === 0 &&
      temporaryNode?.countChildren === 0 &&
      !(temporaryNode?.children?.length || temporaryNode?.hasChildren)
    ) {
      form.setFieldsValue({
        targetFolder: temporaryNode?.alias,
      });
      setTemporaryNode(undefined);
    }
  }, [temporaryNode, attributeType]);

  const getDataTypeOptions = () => {
    if (enableAutoCategorization) {
      switch (parentDataType) {
        case 1:
          return [{ label: formatMessage('uns.relational'), value: 2 }];
        case 3:
          return [{ label: formatMessage('uns.timeSeries'), value: 1 }];
        default:
          return [];
      }
    } else {
      return [
        { label: formatMessage('uns.timeSeries'), value: 1 },
        { label: formatMessage('uns.relational'), value: 2 },
      ];
    }
  };

  return (
    <ProModal
      title={formatMessage('uns.batchGeneration')}
      draggable={false}
      width={1000}
      open={reverserOpen}
      onCancel={() => close()}
      maskClosable={false}
      centered={false}
      keyboard={false}
      onFullScreenCallBack={(e) => {
        setFullScreen(e);
      }}
    >
      <Form
        name="reverseForm"
        form={form}
        colon={false}
        style={{ position: 'relative' }}
        initialValues={{
          dataType: enableAutoCategorization ? 2 : 1,
          source: 'json',
        }}
        labelCol={{ span: 6 }}
        wrapperCol={{ span: 18 }}
        labelAlign="left"
        labelWrap
      >
        <Form.Item name="attributeType" label={formatMessage('uns.attributeGenerationMethod')} initialValue={1}>
          <ComRadio
            options={[
              { label: formatMessage('uns.reverseGeneration'), value: 1 },
              { label: formatMessage('uns.mount'), value: 2 },
            ]}
            onChange={(e) => {
              form.setFieldsValue({
                source: e.target.value === 1 ? 'json' : 'connect',
                jsonData: undefined,

                // targetFolder: undefined,
                dataSource: undefined,
                persistence: false,
                dashboard: false,
                syncMeta: false,
                ...(enableAutoCategorization
                  ? {
                      parentDataType: e.target.value === 2 ? 3 : 1,
                      dataType: e.target.value === 2 ? 1 : 2,
                    }
                  : {}),
              });
            }}
          />
        </Form.Item>

        {enableAutoCategorization && (
          <Form.Item name="parentDataType" label={formatMessage('uns.parentDataType')} initialValue={1}>
            <ComRadio
              options={[
                { label: formatMessage('uns.state'), value: 1 },
                { label: formatMessage('uns.metric'), value: 3 },
              ]}
              onChange={(e) => {
                switch (e.target.value) {
                  case 1:
                    form.setFieldsValue({ dataType: 2 });
                    jsonFormRef?.current?.batchModifyDataType?.(2, 1);
                    break;
                  case 3:
                    form.setFieldsValue({ dataType: 1 });
                    jsonFormRef?.current?.batchModifyDataType?.(1, 3);
                    break;
                  default:
                    break;
                }
              }}
              disabled={enableAutoCategorization && attributeType === 2}
            />
          </Form.Item>
        )}

        <Form.Item name="dataType" label={formatMessage('uns.databaseType')}>
          <ComRadio
            options={getDataTypeOptions()}
            onClick={(e) => {
              jsonFormRef?.current?.batchModifyDataType?.(e.target.value);
            }}
            disabled={enableAutoCategorization && attributeType === 2}
          />
        </Form.Item>
        <Divider style={{ borderColor: '#c6c6c6' }} />
        <Form.Item name="source" label={formatMessage('uns.source')} rules={[{ required: true }]}>
          <Select placeholder={formatMessage('uns.source')} options={sourceOptionsMap[attributeType]} />
        </Form.Item>
        {attributeType === 2 && (
          <Form.Item name="targetFolder" label={formatMessage('uns.targetFolder')} rules={[{ required: true }]}>
            <Select
              placeholder={formatMessage('uns.targetFolder')}
              options={folderList}
              fieldNames={{ label: 'path', value: 'alias' }}
              showSearch
              optionFilterProp="path"
            />
          </Form.Item>
        )}
        {getSourceForm()}
      </Form>
    </ProModal>
  );
};
export default ReverseModal;
