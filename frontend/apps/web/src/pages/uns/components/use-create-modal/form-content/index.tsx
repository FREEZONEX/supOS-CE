import { type FC, useState, useEffect } from 'react';
import { Form } from 'antd';
import { getTypes, getAllTemplate } from '@/apis/inter-api/uns';
import { useTranslate, useFormValue } from '@/hooks';
import FormItems from './FormItems';
import { uniqBy } from 'lodash-es';

import type { FieldItem } from '@/pages/uns/types';
import { useBaseStore } from '@/stores/base';
import { useI18nStore } from '@/stores/i18n-store';

export interface FormContentProps {
  step: number;
  addNamespaceForAi: { [key: string]: any };
  setAddNamespaceForAi: (e: any) => void;
  open: boolean;
  addModalType: string;
}

type TemplateItemType = { label: string; value: string };

type GetFormDataParamsType = {
  currentStep: number;
  dataType: number;
  modelId: string;
  fields: FieldItem[];
  attributeType: number;
  windowType: string;
  addModalType: string;
};
const FormContent: FC<FormContentProps> = ({ step, addNamespaceForAi, setAddNamespaceForAi, open, addModalType }) => {
  const form = Form.useFormInstance();
  const formatMessage = useTranslate();
  const { dashboardType } = useBaseStore((state) => ({
    dashboardType: state.dashboardType,
  }));
  const lang = useI18nStore((state) => state.lang);

  const [types, setTypes] = useState([]);
  const [templateList, setTemplateList] = useState<TemplateItemType[]>([]); //模版列表

  const topic = useFormValue('topic', form);
  const path = useFormValue('path', form);
  const calculationType = useFormValue('calculationType', form);
  const dataType = useFormValue('dataType', form);
  const modelId = useFormValue('modelId', form);
  const attributeType = useFormValue('attributeType', form);
  const fields = useFormValue('fields', form) || [];
  const windowType = useFormValue(['streamOptions', 'window', 'windowType'], form);

  useEffect(() => {
    if (!open) return;
    getTypes()
      .then((res: any) => {
        setTypes(res || []);
      })
      .catch((err) => {
        console.log(err);
      });
    getAllTemplate({ pageNo: 1, pageSize: 9999 }).then((res: any) => {
      if (res && Array.isArray(res)) {
        const _res = res.map((item) => ({
          ...item,
          label: item.name,
          value: item.id,
        }));
        setTemplateList([{ label: formatMessage('common.custom'), value: 'custom' }].concat(_res));
      }
    });
  }, [open]);

  const selectAll = (options: any[] = []) => {
    const currentReferTopics = form.getFieldValue('referIds') || [];
    const referIds = uniqBy(
      [...currentReferTopics, ...options.map((i) => ({ label: i.path, value: i.id }))].slice(0, 100),
      'value'
    );
    form.setFieldsValue({ referIds });
  };

  const getFormData = (data: GetFormDataParamsType) => {
    const { currentStep, dataType, modelId, fields, attributeType, windowType } = data;
    const isCreateFolder = addModalType.includes('Folder');
    const isFormTopic = addModalType.includes('topic');
    const formItemList = [];

    if (currentStep === 1) {
      //第一步
      if (isFormTopic) {
        formItemList.push({
          formType: 'input',
          formProps: {
            name: 'path',
            label: formatMessage('uns.path'),
          },
          childProps: { disabled: true },
        });
      } else {
        formItemList.push(
          {
            formType: 'input',
            formProps: {
              name: 'topic',
              label: formatMessage('uns.namespace'),
              tooltip: {
                title: formatMessage('uns.namespaceTooltip'),
              },
            },
            childProps: { disabled: true },
          },
          {
            formType: 'input',
            formProps: {
              name: 'alias',
              label: formatMessage('uns.alias'),
              tooltip: {
                title: formatMessage('uns.aliasTooltip'),
              },
            },
            childProps: { disabled: true },
          }
        );
      }
      formItemList.push({ formType: 'divider', formProps: { name: 'aliasDivider' } });
      if (isFormTopic) {
        formItemList.push({
          formType: 'input',
          formProps: {
            name: 'topicName',
            label: formatMessage('common.name'),
          },
          childProps: { disabled: true },
        });
      } else {
        formItemList.push({
          formType: 'input',
          formProps: {
            name: 'name',
            label: formatMessage('common.name'),
            rules: [
              { required: true, message: formatMessage('uns.pleaseInputName') },
              { pattern: /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/, message: formatMessage('uns.nameFormat') },
              {
                max: 63,
                message: formatMessage('uns.labelMaxLength', { label: formatMessage('common.name'), length: 63 }),
              },
              {
                validator: (_: any, value: string) => {
                  if (isCreateFolder && ['label', 'template'].includes(value)) {
                    return Promise.reject(new Error(formatMessage('uns.prohibitKeywords')));
                  } else {
                    return Promise.resolve();
                  }
                },
              }, // 添加自定义校验规则
            ],
          },
        });
      }

      formItemList.push(
        {
          formType: 'input',
          formProps: {
            name: 'displayName',
            label: formatMessage('uns.displayName'),
            rules: [{ max: 128 }],
          },
        },
        {
          formType: 'textArea',
          formProps: {
            name: 'description',
            label: formatMessage(isCreateFolder ? 'uns.folderDescription' : 'uns.fileDescription'),
            rules: [
              {
                max: 255,
                message: formatMessage('uns.labelMaxLength', {
                  label: formatMessage(isCreateFolder ? 'uns.folderDescription' : 'uns.fileDescription'),
                  length: 255,
                }),
              },
            ],
          },
          childProps: { rows: 2 },
        }
      );
      if ([1, 2].includes(dataType) && !isCreateFolder) {
        formItemList.push({
          formType: 'select',
          formProps: {
            label: formatMessage('uns.writDownData'),
            name: 'accessLevel',
            initialValue: 'READ_ONLY',
          },
          childProps: {
            options: [
              { label: formatMessage('uns.true'), value: 'READ_WRITE' },
              { label: formatMessage('uns.false'), value: 'READ_ONLY' },
            ],
          },
        });
      }
      formItemList.push({ formType: 'expandFormList', formProps: { name: 'expandFormList' } });

      if (isCreateFolder) {
        //创建文件夹
        formItemList.push(
          { formType: 'divider', formProps: { name: 'modelDescriptionDivider' } },
          {
            formType: 'select',
            formProps: {
              name: 'modelId',
              label: formatMessage('common.template'),
              initialValue: 'custom',
            },
            childProps: {
              showSearch: true,
              optionFilterProp: 'path',
              options: templateList,
              onChange: (modelId: string) => {
                if (modelId === 'custom' || !modelId) {
                  form.setFieldValue('fields', undefined);
                }
              },
            },
          }
        );
        if (modelId === 'custom' && fields.length) {
          formItemList.push({
            formType: 'checkbox',
            formProps: {
              name: 'createTemplate',
              label: formatMessage('uns.generationTemplate'),
              initialValue: true,
              valuePropName: 'checked',
            },
          });
        }
        formItemList.push({
          formType: 'fieldsFormList',
          formProps: {
            name: 'fields',
          },
          childProps: {
            disabled: modelId !== 'custom',
            isCreateFolder,
            showMainKey: false,
            types,
            addNamespaceForAi,
            setAddNamespaceForAi,
          },
        });
      } else {
        formItemList.push({
          formType: 'radioGroup',
          formProps: {
            name: 'dataType',
            label: formatMessage('uns.databaseType'),
            initialValue: 1,
            tooltip: {
              title: (
                <div>
                  <span>• {formatMessage('uns.timeSeries')}</span> — {formatMessage('uns.dataTypeTooltip-TimeSeries')}
                  <br />
                  <span>• {formatMessage('uns.relational')}</span> — {formatMessage('uns.dataTypeTooltip-Relational')}
                  <br />
                  {!isFormTopic && (
                    <>
                      <span>• {formatMessage('uns.calculation')}</span> —&nbsp;
                      {formatMessage('uns.dataTypeTooltip-Calculation')}
                      <br />
                      <span>• {formatMessage('uns.aggregation')}</span> —&nbsp;
                      {formatMessage('uns.dataTypeTooltip-Aggregation')}
                      <br />
                      <span>• {formatMessage('uns.reference')}</span> —&nbsp;
                      {formatMessage('uns.dataTypeTooltip-Reference')}
                    </>
                  )}
                </div>
              ),
            },
          },
          childProps: {
            style: { flexWrap: 'wrap' },
            options: [
              { label: formatMessage('uns.timeSeries'), value: 1 },
              { label: formatMessage('uns.relational'), value: 2 },
              ...(isFormTopic
                ? []
                : [
                    { label: formatMessage('uns.calculation'), value: 3 },
                    { label: formatMessage('uns.aggregation'), value: 6 },
                    { label: formatMessage('uns.reference'), value: 7 },
                  ]),
            ],
            onChange: (e: any) => {
              const resetObj = {
                refers: [{}],
                frequency: undefined,
                aggregationModel: undefined,
                referIds: undefined,
                mainKey: undefined,
                timeReference: undefined,
              };
              if (![1, 2].includes(e.target.value)) {
                Object.assign(resetObj, {
                  fields: [{}],
                  attributeType: 1,
                  modelId: undefined,
                  jsonData: undefined,
                  jsonList: [],
                  jsonDataPath: undefined,
                  source: undefined,
                  dataSource: undefined,
                  table: undefined,
                  next: false,
                });
              }
              form.setFieldsValue(resetObj);
            },
          },
        });
        if ([1, 2].includes(dataType)) {
          //选择时序或关系型
          if (isFormTopic) {
            formItemList.push({
              formType: 'topicToUnsFieldsList',
              formProps: { name: 'topicToUnsFieldsList' },
              childProps: { types },
            });
          } else {
            formItemList.push(
              { formType: 'divider', formProps: { name: 'timeSeriesDivider' } },
              {
                formType: 'radioGroup',
                formProps: {
                  name: 'attributeType',
                  label: formatMessage('uns.attributeGenerationMethod'),
                  initialValue: 1,
                  tooltip: {
                    title: (
                      <div>
                        <span>• {formatMessage('common.custom')}</span> —&nbsp;
                        {formatMessage('uns.attributeGenerationMethodTooltip-Custom')}
                        <br />
                        <span>• {formatMessage('common.template')}</span> —&nbsp;
                        {formatMessage('uns.attributeGenerationMethodTooltip-Template')}
                        <br />
                        <span>• {formatMessage('uns.reverseGeneration')}</span> —&nbsp;
                        {formatMessage('uns.attributeGenerationMethodTooltip-ReverseGeneration')}
                      </div>
                    ),
                  },
                  className: lang === 'en-US' ? 'customLabelStyle' : '',
                },
                childProps: {
                  options: [
                    { label: formatMessage('common.custom'), value: 1 },
                    { label: formatMessage('common.template'), value: 2 },
                    { label: formatMessage('uns.reverseGeneration'), value: 3 },
                  ],
                  onChange: () => {
                    form.setFieldsValue({
                      fields: [{}],

                      modelId: undefined,
                      jsonData: undefined,
                      jsonList: [],
                      jsonDataPath: undefined,
                      source: undefined,
                      dataSource: undefined,
                      table: undefined,
                      next: false,
                      mainKey: undefined,
                    });
                  },
                },
              }
            );
            if (attributeType === 1) {
              //选择自定义属性
              formItemList.push({
                formType: 'fieldsFormList',
                formProps: { name: 'fields' },
                childProps: {
                  types,
                  addNamespaceForAi,
                  setAddNamespaceForAi,
                  showMoreBtn: dataType === 1,
                },
              });
            }
            if (attributeType === 2) {
              //选择模板属性
              formItemList.push({
                formType: 'modelFieldsForm',
                formProps: { name: 'modelFieldsForm' },
                childProps: {
                  options: templateList.slice(1),
                  types,
                },
              });
            }
            if (attributeType === 3) {
              //选择逆向生成属性
              formItemList.push({
                formType: 'reverseGeneration',
                formProps: { name: 'reverseGeneration' },
                childProps: { types },
              });
            }
          }
        }
        if (dataType === 3) {
          //选择计算
          formItemList.push(
            {
              formType: 'radioGroup',
              formProps: {
                name: 'calculationType',
                label: formatMessage('uns.calculationType'),
                initialValue: 3,
              },
              childProps: {
                options: [
                  { label: formatMessage('uns.realtime'), value: 3 },
                  // { label: formatMessage('common.history'), value: 4 },
                ],
                onChange: () => {
                  form.setFieldsValue({ fields: [{}] });
                },
              },
            },
            { formType: 'divider', formProps: { name: 'calculationTypeDivider' } },
            {
              formType: 'fieldsFormList',
              formProps: { name: 'fields' },
              childProps: {
                types,
                addNamespaceForAi,
                setAddNamespaceForAi,
                showMoreBtn: true,
              },
            }
          );
        }
        if (dataType === 6) {
          //选择聚合
          formItemList.push(
            { formType: 'divider', formProps: { name: 'aggregationDivider' } },
            {
              formType: 'frequency',
              formProps: {
                name: 'frequency',
                label: formatMessage('uns.frequency'),
                style: { marginBottom: 0 },
                required: true,
                tooltip: {
                  title: formatMessage('uns.frequencyTooltip'),
                },
              },
            },
            {
              formType: 'divider',
              formProps: {
                name: 'frequencyDivider',
              },
            },
            {
              formType: 'checkbox',
              formProps: {
                name: 'save2db',
                label: formatMessage('uns.persistence'),
                initialValue: false,
                valuePropName: 'checked',
                tooltip: {
                  title: formatMessage('uns.persistenceTooltip'),
                },
              },
            },
            {
              formType: 'searchSelect',
              formProps: {
                name: 'referIds',
                label: formatMessage('uns.aggregationTarget'),
                rules: [{ required: true }],
                tooltip: {
                  title: <div>{formatMessage('uns.aggregationTargetTooltip')}</div>,
                },
              },
              childProps: {
                placeholder: formatMessage('uns.searchInstance'),
                mode: 'multiple',
                maxCount: 100,
                selectAll: selectAll,
                apiParams: { type: 2, normal: true },
                labelInValue: true,
              },
            }
          );
        }
        if (dataType === 7) {
          //选择引用
          formItemList.push(
            { formType: 'divider', formProps: { name: 'referenceDivider' } },
            {
              formType: 'searchSelect',
              formProps: {
                name: 'referId',
                label: formatMessage('uns.referenceTarget'),
                rules: [{ required: true }],
              },
              childProps: {
                placeholder: formatMessage('uns.searchInstance'),
                apiParams: { type: 2, normal: true },
                labelInValue: true,
              },
            },
            { formType: 'divider', formProps: { name: 'tagsDivider' } },
            {
              formType: 'tagSelect',
              formProps: {
                name: 'tags',
                label: formatMessage('common.label'),
                tooltip: {
                  title: formatMessage('uns.labelTooltip'),
                },
              },
              childProps: {
                tagMaxLen: 63,
              },
            }
          );
        }
      }
    }
    if (currentStep === 2) {
      //第二步
      formItemList.push(
        {
          formType: 'showTopic',
          formProps: {
            name: 'showTopic',
            label: formatMessage(isFormTopic ? 'uns.path' : 'uns.namespace'),
            initialValue: isFormTopic ? path : topic,
            tooltip: isFormTopic
              ? undefined
              : {
                  title: formatMessage('uns.namespaceTooltip'),
                },
            style: { marginBottom: 0 },
          },
        },
        {
          formType: 'divider',
          formProps: {
            name: 'topic2Divider',
          },
        }
      );
      if (dataType === 3) {
        if (calculationType === 3) {
          formItemList.push(
            {
              formType: 'expressionForm',
              formProps: { name: 'expressionForm' },
              childProps: {
                showTimeReference: true,
              },
            },
            { formType: 'divider', formProps: { name: 'expressionFormDivider' } }
          );
        }
        if (calculationType === 4) {
          formItemList.push({ formType: 'aggForm', formProps: { name: 'aggForm' } });
        }
      }
      formItemList.push(
        {
          formType: 'tagSelect',
          formProps: {
            name: 'tags',
            label: formatMessage('common.label'),
            tooltip: {
              title: formatMessage('uns.labelTooltip'),
            },
          },
          childProps: {
            tagMaxLen: 63,
          },
        },
        { formType: 'divider', formProps: { name: 'tagsDivider' } }
      );
      if (dataType !== 3) {
        formItemList.push({
          formType: 'checkbox',
          formProps: {
            name: 'addFlow',
            label: formatMessage('uns.mockData'),
            initialValue: true,
            valuePropName: 'checked',
            tooltip: {
              title: formatMessage('uns.mockDataTooltip'),
            },
          },
        });
      }
      if (dashboardType?.includes('grafana')) {
        formItemList.push({
          formType: 'checkbox',
          formProps: {
            name: 'addDashBoard',
            label: formatMessage('uns.autoDashboard'),
            initialValue: true,
            valuePropName: 'checked',
            tooltip: {
              title: formatMessage('uns.autoDashboardTooltip'),
            },
            className: lang === 'en-US' ? 'customLabelStyle' : '',
          },
        });
      }
      formItemList.push({
        formType: 'checkbox',
        formProps: {
          name: 'save2db',
          label: formatMessage('uns.persistence'),
          initialValue: false,
          valuePropName: 'checked',
          tooltip: {
            title: formatMessage('uns.persistenceTooltip'),
          },
        },
        childProps: { disabled: dataType === 3 && calculationType === 4 },
      });
      if (dataType === 3 && calculationType === 4) {
        formItemList.push({
          formType: 'checkbox',
          formProps: {
            name: 'advancedOptions',
            label: formatMessage('streams.advancedOptions'),
            initialValue: false,
            valuePropName: 'checked',
          },
          childProps: {
            onChange: () => {
              form.setFieldValue('_advancedOptions', undefined);
            },
            disabled: windowType === 'COUNT_WINDOW',
          },
        });
      }
    }
    if (currentStep === 3) {
      //第三步
      formItemList.push(
        {
          formType: 'showTopic',
          formProps: {
            name: 'showTopic',
            label: formatMessage('uns.namespace'),
            initialValue: topic,
            tooltip: {
              title: formatMessage('uns.namespaceTooltip'),
            },
            style: { marginBottom: 0 },
          },
        },
        {
          formType: 'divider',
          formProps: {
            name: 'topic3Divider',
          },
        }
      );
      if (dataType === 3 && calculationType === 4) {
        formItemList.push({
          formType: 'advancedOptions',
          formProps: { name: 'advancedOptions' },
        });
      }
    }
    formItemList.push({
      formType: 'divider',
      formProps: {
        name: 'bottomDivider',
      },
    });
    return formItemList;
  };
  return (
    <FormItems
      formData={getFormData({
        currentStep: step,
        dataType,
        modelId,
        fields,
        attributeType,
        windowType,
        addModalType,
      })}
    />
  );
};

export default FormContent;
