import { getInstanceInfo, getModelInfo, searchTreeData } from '@/apis/core-api/uns';
import { useTranslate } from '@/hooks';
import { Close } from '@/components/lucide-icon/carbon';
import { Button, Drawer, Form, Input, Tooltip } from 'antd';
import dayjs from 'dayjs';
import { pinyin } from 'pinyin-pro';
import { useCallback, useEffect, useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import FormContent from './form-content';
import FormStep from './form-step';
import {
  buildNamespacePreview,
  deriveParentDataTypeFromName,
  getDefaultDataTypeByParent,
  getLeafSegment,
} from './path-utils';
import './index.scss';

import type { FieldItem, InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import { useBaseStore } from '@/stores/base';
import { getExpression, parseArrayToObjects, parseTime } from '@/utils/uns';
import type { TreeStoreActions } from '../../store/types';

export interface UseOptionModalProps {
  successCallBack: InitTreeDataFnType;
  addNamespaceForAi: { [key: string]: any };
  setAddNamespaceForAi: any;
  changeCurrentPath: (node?: UnsTreeNode) => void;
  setTreeMap: TreeStoreActions['setTreeMap'];
}

const extendToArr = (extend: { [key: string]: string }) => {
  if (!extend) return undefined;
  const arr: { key: string; value: string }[] = [];
  Object.keys(extend).forEach((item) => {
    arr.push({
      key: item,
      value: extend[item],
    });
  });
  return arr;
};

const useOptionModal = ({
  successCallBack,
  addNamespaceForAi,
  setAddNamespaceForAi,
  changeCurrentPath,
  setTreeMap,
}: UseOptionModalProps) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const [uuid, setUuid] = useState('');
  const [open, setOpen] = useState(false);
  const [step, setStep] = useState(1);
  const [addModalType, setAddModalType] = useState<string>(''); //addFolder,addFile
  const [sourcePath, setSourcePath] = useState<string>(''); //父文件路径
  const [sourceId, setSourceId] = useState<string>(''); //父文件id
  const [topicType, setTopicType] = useState<number>(0); //文件夹类型

  const {
    systemInfo: { enableAutoCategorization },
  } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
  }));

  const name = Form.useWatch('name', form) || form.getFieldValue('name');
  const dataType = Form.useWatch('dataType', form) || form.getFieldValue('dataType');
  const parentDataType = Form.useWatch('parentDataType', form) || form.getFieldValue('parentDataType');

  const isCreateFolder = addModalType.includes('Folder');
  const isCreateFile = addModalType === 'addFile';

  const onClose = () => {
    setOpen(false);
    form.resetFields();
    setStep(1);
    setAddNamespaceForAi?.(null);
    setSourcePath('');
    setSourceId('');
  };

  const handleMockNextStep = async (customStep: number) => {
    setStep(() => customStep);
  };

  const changeModalType = useCallback(
    async (type?: string, targetNode?: UnsTreeNode, pasteNode?: UnsTreeNode) => {
      const {
        id,
        parentId = '',
        path = '',
        parentPath = '',
        pathType: nodeType,
        dataType,
        parentDataType,
      } = targetNode || {};
      const _folderPath = nodeType === 0 ? path : parentPath;
      const folderPath = _folderPath ? `${_folderPath}/` : '';
      const folderId = (nodeType === 0 ? id : parentId) || '';
      const _topicType = nodeType === 0 && dataType ? dataType : parentDataType || 0;
      setSourcePath(folderPath);
      setSourceId(folderId as string);
      setTopicType(_topicType);

      setOpen(true);
      if (pasteNode) {
        //数据回填
        const isPasteFolder = pasteNode.pathType === 0;
        setAddModalType(isPasteFolder ? 'pasteFolder' : 'pasteFile');
        const getInfo = isPasteFolder ? getModelInfo : getInstanceInfo;
        const detail: any = await getInfo({ id: pasteNode.id });
        if (isPasteFolder) {
          const { description, pathName: name, fields, extend, displayName } = detail || {};

          form.setFieldsValue({
            displayName,
            description,
            name,
            fields,
            extend: extendToArr(extend),
            pasteInfo: pasteNode ? { sourceId: pasteNode.id, targetId: targetNode?.id } : undefined,
            pasteNode: pasteNode,
          });
        } else {
          const {
            displayName,
            description,
            pathName: name,
            jsonFields,
            dataType,
            protocol = {},
            withFlow,
            persistence,
            expression,
            refers = [],
            dataPath,
            extend,
            extendFieldUsed,
            parentDataType,
            pathType,
          } = detail || {};
          let { fields } = detail || {};

          if (dataType === 8) {
            fields = jsonFields;
          }

          const backfillForm: { [key: string]: any } = {
            displayName,
            description,
            name,
            persistence,
            extend: extendToArr(extend),
            dataType,
            pasteInfo: pasteNode ? { sourceId: pasteNode.id, targetId: targetNode?.id } : undefined,
            parentDataType,
          };
          if (pathType === 2) {
            switch (dataType) {
              case 1:
              case 2:
                Object.assign(backfillForm, {
                  attributeType: 1,
                  addFlow: withFlow,
                  mainKey: fields.findIndex((item: FieldItem) => item.unique === true),
                  extendFieldUsed,
                });
                fields.forEach((e: FieldItem) => {
                  delete e.unique;
                });
                break;
              case 3: {
                type ReferType = {
                  id: string;
                  path: string;
                  field: string;
                  uts?: boolean;
                };
                const refersRes = await Promise.all(refers?.map((e: ReferType) => getInstanceInfo({ id: e.id })) || []);
                const _refers = refers?.map((refer: ReferType, i: number) => ({
                  ...refer,
                  refer: {
                    label: refer.path,
                    value: refer.id,
                    option: {
                      dataType: refersRes[i]?.dataType,
                    },
                  },
                  fields: refersRes[i]?.fields?.filter?.(
                    (t: FieldItem) =>
                      !t.systemField && ['INTEGER', 'LONG', 'FLOAT', 'DOUBLE', 'BOOLEAN'].includes(t.type)
                  ) || [{ name: refer.field }],
                }));
                //实时计算
                Object.assign(backfillForm, {
                  dataType: 3,
                  calculationType: 3,
                  refers: _refers,
                  expression: getExpression(refers, expression),
                  timeReference: refers?.find((item: ReferType) => item.uts)?.id,
                  extendFieldUsed,
                });
                break;
              }
              case 4: {
                //历史计算
                const {
                  window,
                  trigger = '',
                  waterMark,
                  deleteMark,
                  fillHistory,
                  ignoreUpdate,
                  ignoreExpired,
                  startTime,
                  endTime,
                } = protocol;
                Object.assign(backfillForm, {
                  dataType: 3,
                  calculationType: 4,
                  DataSource: { value: dataPath },
                  functions: parseArrayToObjects(fields.map((field: FieldItem) => field.index)),
                  whereCondition: getExpression(refers, expression, true),
                  streamOptions: { window },
                  advancedOptions:
                    !!trigger || !!waterMark || !!deleteMark || fillHistory || ignoreUpdate || ignoreExpired,
                  _advancedOptions: {
                    trigger: trigger.split(' ')[0],
                    delayTime: trigger.split(' ')[1],
                    waterMark,
                    deleteMark,
                    fillHistory,
                    ignoreUpdate,
                    ignoreExpired,
                    startTime: startTime ? dayjs(startTime, 'YYYY-MM-DD') : undefined,
                    endTime: endTime ? dayjs(endTime, 'YYYY-MM-DD') : undefined,
                  },
                });
                if (dataPath) {
                  searchTreeData({ type: 3, pageNo: 1, pageSize: 99999 }).then((res: any) => {
                    const whereFieldList =
                      res
                        ?.find((e: { topic: string }) => e.topic === dataPath)
                        ?.fields?.map(({ name, type }: FieldItem) => {
                          return { label: name, value: name, type };
                        }) || [];
                    form.setFieldsValue({ whereFieldList });
                  });
                }

                break;
              }
              case 6:
                Object.assign(backfillForm, {
                  frequency: protocol.frequency
                    ? {
                        value: parseTime(protocol.frequency)[0],
                        unit: parseTime(protocol.frequency)[1],
                      }
                    : {},
                  referIds: refers.map((item: { id: string; path: string }) => ({
                    label: item.path,
                    value: item.id,
                  })),
                });
                break;
              case 7:
                Object.assign(backfillForm, {
                  referId: {
                    label: refers?.[0].path,
                    value: refers?.[0].id,
                  },
                });
                break;
              default:
                break;
            }
          }
          console.log(backfillForm, 'backfillForm');
          form.setFieldsValue(backfillForm);
          setTimeout(() => {
            form.setFieldsValue({
              fields,
            });
          }, 500);
        }
      } else {
        setAddModalType(type || '');
        const _id = nodeType === 0 ? id : nodeType === 2 ? parentId : '';
        if (_id) {
          const detail: any = await getModelInfo({ id: _id });
          const { fields }: { fields: FieldItem[] } = detail || {};

          switch (type) {
            case 'addFolder':
              form.setFieldsValue({
                topic: folderPath,
                fields: fields,
              });
              break;
            case 'addFile':
              form.setFieldsValue({
                topic: folderPath,
                fields: [...(fields || []), {}],
                attributeType: 1,
                ...(enableAutoCategorization
                  ? {
                      parentDataType: _topicType || 1,
                      dataType: _topicType === 1 ? 8 : _topicType === 2 ? 8 : _topicType === 3 ? 1 : 8,
                    }
                  : {}),
              });
              break;
            default:
              break;
          }
        } else {
          form.setFieldsValue(
            type?.includes('File')
              ? { fields: [{}] }
              : {
                  fields: undefined,
                }
          );
        }
      }
    },
    [setSourcePath, setSourceId, setOpen, setAddModalType, form, enableAutoCategorization]
  );

  useEffect(() => {
    if (!open) return;
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setUuid(uuidv4().replace(/-/g, '').slice(0, 20));
  }, [open]);

  const nameChange = useCallback(
    (val?: string) => {
      const leafName = getLeafSegment(val);
      form.setFieldsValue({
        alias: `_${pinyin(leafName || '', { toneType: 'none', v: true })
          ?.replace(/\s+/g, '')
          ?.replace(/-/g, '_')
          .slice(0, 38)}_${uuid}`,
        topic: buildNamespacePreview(sourcePath, val),
      });
    },
    [form, sourcePath, uuid]
  );

  useEffect(() => {
    nameChange(name);
  }, [name, nameChange]);

  useEffect(() => {
    if (!open || !isCreateFile || !enableAutoCategorization) return;
    const parentDataType = deriveParentDataTypeFromName(name);
    if (!parentDataType) return;
    const nextDataType = getDefaultDataTypeByParent(parentDataType);
    if (
      form.getFieldValue('parentDataType') !== parentDataType ||
      (nextDataType && form.getFieldValue('dataType') !== nextDataType)
    ) {
      form.setFieldsValue({
        parentDataType,
        ...(nextDataType ? { dataType: nextDataType } : {}),
      });
    }
  }, [open, isCreateFile, enableAutoCategorization, name, form]);

  useEffect(() => {
    if (!open || addModalType !== 'addFile' || !name) return;
    Promise.resolve()
      .then(() => form.validateFields(['name']))
      .catch(() => undefined);
  }, [open, addModalType, name, dataType, parentDataType, form]);

  const titleMap: { [key: string]: string } = {
    addFolder: formatMessage('uns.newFolder'),
    addFile: formatMessage('uns.newFile'),
    pasteFolder: formatMessage('uns.pasteFolder'),
    pasteFile: formatMessage('uns.pasteFile'),
  };

  const Dom = (
    <Drawer
      rootClassName="optionDrawerWrap"
      title={titleMap[addModalType]}
      onClose={onClose}
      open={open}
      closable={false}
      extra={
        <Tooltip title={formatMessage('common.close')}>
          <Button color="default" variant="text" onClick={onClose} icon={<Close size={20} />} />
        </Tooltip>
      }
      maskClosable={false}
      destroyOnHidden={false}
      width={680}
      classNames={{
        body: 'newFolderOrFileModalBody',
      }}
    >
      <div className="optionContent">
        <Form
          className="useCreateModalForm"
          name="namespaceForm"
          form={form}
          colon={false}
          style={{ position: 'relative' }}
          initialValues={{
            refers: [{}],
          }}
          labelCol={{ span: 6 }}
          wrapperCol={{ span: 18 }}
          labelAlign="left"
          labelWrap
        >
          <FormContent
            step={step}
            addNamespaceForAi={addNamespaceForAi}
            setAddNamespaceForAi={setAddNamespaceForAi}
            open={open}
            addModalType={addModalType}
            topicType={topicType}
          />
          <FormStep
            step={step}
            setStep={setStep}
            handleClose={(cb) => {
              onClose();
              cb?.();
            }}
            isCreateFolder={isCreateFolder}
            addNamespaceForAi={addNamespaceForAi}
            setAddNamespaceForAi={setAddNamespaceForAi}
            successCallBack={successCallBack as any}
            changeCurrentPath={changeCurrentPath}
            setTreeMap={setTreeMap}
            sourceId={sourceId}
            addModalType={addModalType}
          />
          <Form.Item name="pasteInfo" hidden>
            <Input />
          </Form.Item>
        </Form>
      </div>
    </Drawer>
  );
  return {
    OptionModal: Dom,
    setOptionOpen: changeModalType,
    setModalClose: onClose,
    setMockNextStep: handleMockNextStep,
  };
};
export default useOptionModal;
