import { type FC, useState, type Dispatch, type SetStateAction, useEffect } from 'react';
import { ChevronRight, ChevronLeft } from '@carbon/icons-react';
import { Form, App, Button, Flex } from 'antd';
import { addModel, pasteUns } from '@/apis/inter-api/uns';
import { topic2Uns } from '@/apis/inter-api/external';
import { useTranslate, useFormValue } from '@/hooks';
import dayjs from 'dayjs';
import { cloneDeep } from 'lodash-es';
import type { FieldItem, UnsTreeNode } from '@/pages/uns/types';
import { ROOT_NODE_ID } from '../../../store/treeStore';
import type { TreeStoreActions } from '../../../store/types';
import { getTargetNode } from '@/utils/uns';
import ComPopupGuide from '@/components/com-popup-guide';
import { useTreeStore } from '@/pages/uns/store/treeStore';

export interface FormStepProps {
  step: number;
  setStep: Dispatch<SetStateAction<number>>;
  handleClose: (cb?: () => void) => void;
  isCreateFolder: boolean;
  addNamespaceForAi: { [key: string]: any };
  setAddNamespaceForAi: (e: any) => void;
  successCallBack: TreeStoreActions['loadData'];
  changeCurrentPath: (node: UnsTreeNode) => void;
  setTreeMap: TreeStoreActions['setTreeMap'];
  sourceId: string;
  addModalType: string;
}

const FormStep: FC<FormStepProps> = ({
  step,
  setStep,
  handleClose,
  isCreateFolder,
  addNamespaceForAi,
  setAddNamespaceForAi,
  successCallBack,
  changeCurrentPath,
  setTreeMap,
  sourceId,
  addModalType,
}) => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const form = Form.useFormInstance();
  const [loading, setLoading] = useState(false);

  const { operationFns, setCurrentTreeMapType, lazyTree, treeData } = useTreeStore((state) => ({
    operationFns: state.operationFns,
    setCurrentTreeMapType: state.setCurrentTreeMapType,
    lazyTree: state.lazyTree,
    treeData: state.treeData,
  }));

  //以下变量用于控制步骤按钮的显示
  const advancedOptions = useFormValue('advancedOptions', form) || false;
  const calculationType = useFormValue('calculationType', form);
  const _dataType = useFormValue('dataType', form);
  const attributeType = useFormValue('attributeType', form);
  const jsonList = useFormValue('jsonList', form);

  const isFormTopic = addModalType.includes('topic');

  const extendToObj = (extend: { key: string; value: string }[]) => {
    if (!extend) return undefined;
    const obj: { [key: string]: string } = {};
    extend.forEach((item) => {
      obj[item.key] = item.value;
    });
    return obj;
  };

  const save = () => {
    form
      .validateFields()
      .then(async () => {
        const next = form.getFieldValue('next');
        if ((attributeType === 3 || (isFormTopic && jsonList?.length > 1)) && _dataType !== 8 && !next) {
          message.error(formatMessage('uns.noFieldsTip'));
          return;
        }
        const {
          alias,
          fields,
          dataType,
          description,
          extend,
          addFlow,
          addDashBoard,
          save2db,
          calculationType,
          refers,
          expression,
          tags,
          mainKey,
          frequency,
          referIds,
          referId,
          modelId,
          name,
          createTemplate,
          displayName,
          timeReference,
          accessLevel,
          extendFieldUsed,

          functions,
          DataSource,
          streamOptions,
          whereCondition,
          havingCondition,
          advancedOptions,

          _advancedOptions,
          table,

          path,
          pasteInfo,
          parentDataType,
          pasteNode,
        } = cloneDeep(form.getFieldsValue(true));
        // 表单验证通过后的操作
        const data: { [key: string]: any } = isCreateFolder
          ? {
              name,
              displayName,
              parentId: sourceId || undefined,
              alias,
              description,
              fields,
              pathType: 0,
              extend: extendToObj(extend),
            }
          : {
              name,
              displayName,
              parentId: sourceId || undefined,
              alias,
              dataType,
              description,
              save2db,
              pathType: 2,
              extend: extendToObj(extend),
              labelNames: tags?.map(({ label, value }: { label: string; value: string }) => label || value) || [],
              fields: [1, 2, 3, 8].includes(dataType) ? fields : undefined,
              addDashBoard,
              parentDataType,
            };

        if (isCreateFolder && modelId === 'custom' && fields?.length > 0) {
          data.createTemplate = createTemplate;
        }
        data.modelId = modelId && modelId !== 'custom' ? modelId : undefined;

        if (!isCreateFolder) {
          switch (dataType) {
            case 1:
            case 2:
              if (dataType === 1) {
                data.fields = fields
                  .filter((e: FieldItem) => !e.systemField)
                  .map(
                    ({
                      name,
                      type,
                      displayName,
                      remark,
                      maxLen,
                      unit,
                      upperLimit,
                      lowerLimit,
                      decimal,
                    }: FieldItem) => ({
                      name,
                      type,
                      displayName,
                      remark,
                      maxLen,
                      unit: extendFieldUsed?.includes('unit') ? unit : undefined,
                      upperLimit: extendFieldUsed?.includes('upperLimit') ? upperLimit : undefined,
                      lowerLimit: extendFieldUsed?.includes('lowerLimit') ? lowerLimit : undefined,
                      decimal: extendFieldUsed?.includes('decimal') ? decimal : undefined,
                    })
                  );
                if (!modelId) {
                  data.extendFieldUsed = extendFieldUsed;
                }
              } else {
                if (mainKey > -1) fields[mainKey].unique = true;
                data.fields = fields.map((e: FieldItem) => ({
                  name: e.name,
                  type: e.type,
                  displayName: e.displayName,
                  remark: e.remark,
                  maxLen: e.maxLen,
                  unique: e.unique,
                }));
              }
              data.addFlow = addFlow;
              if (table?.value) {
                data.protocol = {
                  referenceDataSource: table.value
                    ?.split?.('$分隔符$')
                    ?.filter((e: string) => e !== 'tables')
                    ?.join('.'),
                };
              }
              data.accessLevel = accessLevel;
              break;
            case 3:
              if (calculationType === 3) {
                data.extendFieldUsed = extendFieldUsed;
                data.fields = fields
                  .filter((e: FieldItem) => !e.systemField)
                  .map(
                    ({
                      name,
                      type,
                      displayName,
                      remark,
                      maxLen,
                      unit,
                      upperLimit,
                      lowerLimit,
                      decimal,
                    }: FieldItem) => ({
                      name,
                      type,
                      displayName,
                      remark,
                      maxLen,
                      unit: extendFieldUsed?.includes('unit') ? unit : undefined,
                      upperLimit: extendFieldUsed?.includes('upperLimit') ? upperLimit : undefined,
                      lowerLimit: extendFieldUsed?.includes('lowerLimit') ? lowerLimit : undefined,
                      decimal: extendFieldUsed?.includes('decimal') ? decimal : undefined,
                    })
                  );

                type ReferItemType = {
                  refer: {
                    label: string;
                    value: string;
                  };
                  field: string;
                };
                //实时计算
                data.refers = refers.map((item: ReferItemType) => {
                  return { id: item?.refer?.value, field: item.field, uts: item?.refer?.value === timeReference };
                });
                data.expression = expression ? expression.replace(/\$(.*?)#/g, '$1') : '';
              }
              if (calculationType === 4) {
                //历史计算
                data.dataType = 4;

                data.referTopic = DataSource.value;
                if (functions && Array.isArray(functions) && fields && Array.isArray(fields)) {
                  data.fields = fields.map((field: FieldItem, index: number) => {
                    const func = functions[index];
                    return {
                      ...field,
                      index: `${func.functionType}(${func.key})`,
                    };
                  });
                }

                if (whereCondition) streamOptions.whereCondition = whereCondition.replace(/\$(.*?)#/g, '$1');
                if (havingCondition) streamOptions.havingCondition = havingCondition.replace(/\$(.*?)#/g, '$1');

                if (advancedOptions && _advancedOptions) {
                  //高级流选项
                  if (_advancedOptions.trigger === 'MAX_DELAY') {
                    _advancedOptions.trigger = `MAX_DELAY ${_advancedOptions.delayTime}`;
                    delete _advancedOptions.delayTime;
                  }
                  if (_advancedOptions.startTime)
                    _advancedOptions.startTime = dayjs(_advancedOptions.startTime).format('YYYY-MM-DD');
                  if (_advancedOptions.endTime)
                    _advancedOptions.endTime = dayjs(_advancedOptions.endTime).format('YYYY-MM-DD');
                  Object.keys(_advancedOptions).forEach((key: string) => {
                    if (['', undefined, null].includes(_advancedOptions[key])) delete _advancedOptions[key];
                  });
                }
                data.streamOptions = { ...streamOptions, ..._advancedOptions };
              }
              break;
            case 6:
              Object.assign(data, {
                frequency: frequency.value + frequency.unit,
                referIds: referIds.map((e: { value: string }) => e.value),
              });
              break;
            case 7:
              Object.assign(data, {
                referIds: [referId?.value],
                save2db: undefined,
                addDashBoard: undefined,
              });
              break;
            case 8:
              Object.assign(data, {
                fields: [{ name: 'json', type: 'string' }],
              });
              break;
            default:
              break;
          }
        }

        setLoading(true);
        const handleCallback = (data: { id: string; parentId: string }, queryType: string) => {
          const { id, parentId } = data;
          const hasParentNode = getTargetNode(treeData || [], parentId);

          const _parentId = hasParentNode ? parentId : sourceId ? sourceId : ROOT_NODE_ID;
          const _childId = hasParentNode || parentId === sourceId || !lazyTree ? id : parentId;

          successCallBack(
            {
              queryType,
              key: _parentId,
              newNodeKey: _childId,
              reset: !(sourceId || parentId),
              nodeDetail: pasteNode,
            },
            (_, selectInfo, opt) => {
              const currentNode = getTargetNode(_ || [], _childId);
              if (!selectInfo && !_childId) return;

              changeCurrentPath(
                selectInfo ||
                  currentNode || { key: _childId, id: _childId, pathType: queryType === 'addFolder' ? 0 : 2 }
              );
              setTreeMap(false);
              if (selectInfo) {
                // 非lasy树
                opt?.scrollTreeNode?.(id);
              }
            }
          );
        };
        const labelList =
          tags?.map(({ label, value }: { label: string; value: string | number }) => ({
            ...(label ? { id: value } : { labelName: value }),
          })) || [];

        const addRequest = isFormTopic ? topic2Uns : addModel;
        if (isFormTopic) {
          delete data.alias;
          delete data.parentId;
          data.path = path;
          data.labelList = labelList;
        }
        if (pasteInfo) {
          pasteUns({
            sourceId: pasteInfo?.sourceId || undefined,
            targetId: data?.parentId || undefined,
            newF: data,
          })
            .then(({ msg, code, data }) => {
              handleCallback(data, isCreateFolder ? 'addFolder' : 'addFile');
              handleClose(() => setLoading(false));
              if (code === 206) {
                message.warning(msg);
              } else {
                message.success(formatMessage('uns.pasteSuccess'));
              }
            })
            .catch(() => {
              setLoading(false);
            });
        } else {
          const finalData = [2, 8].includes(data?.dataType)
            ? {
                ...data,
                fields: [1, 2, 3].includes(dataType) ? fields : undefined,
                jsonFields: [8].includes(dataType) && fields?.[0]?.name ? fields : undefined,
              }
            : data;
          addRequest(finalData)
            .then((res: any) => {
              message.success(formatMessage('uns.newSuccessfullyAdded'));
              if (isFormTopic) {
                setCurrentTreeMapType('all');
                handleCallback(res, 'addFile');
                operationFns?.refreshUnusedTopicTree?.(path);
              } else {
                handleCallback(res, isCreateFolder ? 'addFolder' : 'addFile');
              }
              handleClose(() => setLoading(false));
            })
            .catch((err) => {
              setLoading(false);
              console.error(err);
              setAddNamespaceForAi?.(null);
            })
            .finally(() => {
              if (isCreateFolder && addNamespaceForAi) {
                // 如果是新增文件的
                setTimeout(() => {
                  setAddNamespaceForAi({ ...addNamespaceForAi, currentStep: 'openFileNewModal' });
                }, 500);
              }
            });
        }
      })
      .catch((info) => {
        console.error('校验失败:', info);
      })
      .finally(() => {
        if (addNamespaceForAi) {
          setAddNamespaceForAi?.(null);
        }
      });
  };

  const handleStep = async () => {
    return form.validateFields().then(() => {
      if (step === 2) {
        if (calculationType === 4 && advancedOptions) {
          setStep(() => step + 1);
        }
      } else {
        setStep(() => step + 1);
      }
    });
  };

  useEffect(() => {
    const drawerBody = document.querySelector('.newFolderOrFileModalBody');
    if (drawerBody) drawerBody.scrollTop = 0;
  }, [step]);

  return (
    <Flex align="center" justify="flex-end" gap={10}>
      {step > 1 && (
        <Button
          color="default"
          variant="filled"
          size="small"
          style={{ color: 'var(--supos-text-color)', backgroundColor: 'var(--supos-uns-button-color)' }}
          icon={<ChevronLeft />}
          onClick={() => {
            setStep(() => step - 1);
          }}
          disabled={loading}
        >
          {formatMessage('common.prev')}
        </Button>
      )}
      {(step === 1 && [1, 2, 6, 7, 8].includes(_dataType)) ||
      (step === 2 && !advancedOptions) ||
      step === 3 ||
      isCreateFolder ? (
        <ComPopupGuide
          key={isCreateFolder ? 'saveFolder' : 'saveFile'}
          currentStep={addNamespaceForAi?.currentStep}
          stepName={isCreateFolder ? 'saveFolder' : 'saveFile'}
          steps={addNamespaceForAi?.steps}
          placement="left"
          onFinish={() => {
            save?.();
          }}
        >
          <Button color="primary" variant="solid" size="small" onClick={save} loading={loading}>
            {formatMessage('common.save')}
          </Button>
        </ComPopupGuide>
      ) : (
        <ComPopupGuide
          stepName={`next`}
          steps={addNamespaceForAi?.steps}
          currentStep={addNamespaceForAi?.currentStep}
          onFinish={(_, nextStepName) => {
            handleStep()
              .then(() => {
                setAddNamespaceForAi({
                  ...addNamespaceForAi,
                  currentStep: nextStepName,
                });
              })
              .catch(() => {
                setAddNamespaceForAi(null);
              });
          }}
        >
          <Button
            color="default"
            variant="filled"
            size="small"
            icon={<ChevronRight />}
            iconPosition="end"
            onClick={handleStep}
          >
            {formatMessage('common.next')}
          </Button>
        </ComPopupGuide>
      )}
    </Flex>
  );
};
export default FormStep;
