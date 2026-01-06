import { useState, useEffect, useRef, forwardRef, useImperativeHandle } from 'react';
import { Form, Input, Flex, Button, Divider, App } from 'antd';
import { ChevronLeft, ChevronRight, Folder, FolderOpen, Document } from '@carbon/icons-react';
import TagSelect from '@/pages/uns/components/use-create-modal/components/TagSelect';
import FieldsFormList from '@/pages/uns/components/use-create-modal/components/FieldsFormList';
import JsonTree from './JsonTree';
import type { TreeNode, FieldItem } from './JsonTree';
import { json2fsTree, batchReverser } from '@/apis/inter-api/uns';
import { cloneDeep } from 'lodash-es';

import type { UnsTreeNode, InitTreeDataFnType } from '@/pages/uns/types';

import './index.scss';
import ComCheckbox from '@/components/com-checkbox';
import ComRadio from '@/components/com-radio';
import { generateAlias } from '@/utils/uns';
import { useBaseStore } from '@/stores/base';

const { TextArea } = Input;

export interface JsonFormRefProps {
  batchModifyDataType: (dataType: number) => void;
}

export interface JsonFormProps {
  formatMessage: any;
  types: string[];
  currentNode?: UnsTreeNode;
  close: (refreshTree?: boolean) => void;
  fullScreen: boolean;
  initTreeData: InitTreeDataFnType;
}

const JsonForm = forwardRef<JsonFormRefProps, JsonFormProps>(
  ({ formatMessage, types, close, fullScreen, currentNode, initTreeData }, ref) => {
    const form = Form.useFormInstance();
    const [isSave, setIsSave] = useState(false);
    const [loading, setLoading] = useState(false);
    const [bottomBtns, setBottomBtns] = useState(['next']);
    const [treeData, setTreeData] = useState<TreeNode[]>([]);
    const [selectedInfo, setSelectedInfo] = useState<TreeNode | undefined>(undefined);
    const [fileList, setFileList] = useState<TreeNode[]>([]);

    const { message } = App.useApp();

    const jsonTreeRef = useRef<any>(null);

    const {
      dashboardType,
      systemInfo: { qualityName = 'quality', timestampName = 'timeStamp', enableAutoCategorization },
    } = useBaseStore((state) => ({
      dashboardType: state.dashboardType,
      systemInfo: state.systemInfo,
    }));
    const globalDataType = Form.useWatch('dataType', form);
    const hasGrafana = dashboardType?.includes('grafana');
    const jsonData = Form.useWatch('jsonData', form);
    const globalParentDataType = Form.useWatch('parentDataType', form);
    const currentParentDataType = Form.useWatch(['currentNode', 'parentDataType'], form);

    useImperativeHandle(ref, () => ({
      batchModifyDataType,
    }));

    //校验JSON格式
    const validatorJson = (_: any, value: string) => {
      if (!value) return Promise.reject(new Error(formatMessage('uns.pleaseEnterJSON')));
      try {
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

    //递归修改树节点信息
    const modifyNodesRecursively = (nodes: TreeNode[], parentAlias?: string, parentDataPath?: string): void => {
      nodes.forEach((node) => {
        node.parentDataPath = parentDataPath;
        node.alias = generateAlias(node.name);
        node.parentAlias = parentAlias;
        // 如果存在子节点，则递归地修改每个子节点
        if (node.children && node.children.length > 0) {
          node.icon = ({ expanded }: any) => (expanded ? <FolderOpen /> : <Folder />);
          node.pathType = 0;
          modifyNodesRecursively(node.children, node.alias, node.dataPath);
        } else {
          node.icon = <Document />;
          node.pathType = 2;
          node.dataType = globalDataType || 1;
          node.parentDataType = globalParentDataType || 1;
          // node.fields = [...(node.fields || []), ...defaultFields];
        }
      });
    };

    //递归修改最小叶子节点信息
    const handleNodeInfo = (nodes: TreeNode[], field: string, value: boolean | number): void => {
      const newNodes = cloneDeep(nodes);
      const modifyNodesInfo = (nodes: TreeNode[], field: string, value: boolean | number): void => {
        nodes.forEach((node) => {
          if (node.children && node.children.length > 0) {
            modifyNodesInfo(node.children, field, value);
          } else {
            node[field] = value;
          }
        });
      };
      modifyNodesInfo(newNodes, field, value);
      setTreeData(newNodes);
    };

    //获取路径上的所有父节点
    const getAllParentPaths = (paths: string[]): string[] => {
      const allPaths: string[] = [];

      for (const path of paths) {
        const segments = path.split('.');
        let currentPath = '';

        // 构建当前路径的所有父路径
        for (const segment of segments) {
          currentPath = currentPath ? `${currentPath}.${segment}` : segment;
          if (!allPaths.includes(currentPath)) {
            allPaths.push(currentPath);
          }
        }
      }
      return allPaths;
    };

    //根据dataPath获取目标树节点信息
    const getTargetNode = (treeData: TreeNode[], targetDataPath: string): TreeNode | null => {
      for (const node of treeData) {
        if (node.dataPath === targetDataPath) {
          return { ...node };
        }
        // 如果存在子节点并且还没有找到目标节点，则递归搜索子节点
        if (node.children && node.children.length > 0) {
          const foundNode = getTargetNode(node.children, targetDataPath);
          if (foundNode) {
            return { ...foundNode }; // 找到目标节点后立即返回
          }
        }
      }
      return null; // 如果遍历完整个树都没有找到，则返回null
    };

    const prevStep = () => {
      setSelectedInfo(undefined);
      setIsSave(false);
    };

    const nextStep = () => {
      form.validateFields().then(async (values) => {
        const res: any = await json2fsTree(JSON.parse(values.jsonData || null));
        const parentAlias = currentNode?.pathType == 0 ? currentNode?.alias : currentNode?.parentAlias;
        modifyNodesRecursively(res, parentAlias);
        setTreeData(res);
        setIsSave(true);
      });
    };

    const save = async () => {
      const checked = jsonTreeRef.current?.checkedKeys;
      const newChecked = checked.slice(); //浅拷贝

      if (!checked?.length) return message.warning(formatMessage('uns.treeNoCheckedTip'));

      if (newChecked.includes(selectedInfo?.dataPath)) await form.validateFields();

      //根据dataPath数组获取树节点组成的数组
      const checkedNodes = getAllParentPaths(newChecked)
        .map((item: string) => getTargetNode(treeData, item))
        .filter((node): node is TreeNode => !!node);

      //获取最终提交的数组
      const newCheckedNodes = checkedNodes.map((node: TreeNode) => {
        const {
          pathType,
          name,
          description,
          tags,
          save2db = false,
          mainKey,
          fields,
          alias,
          parentAlias,
          addFlow = false,
          addDashBoard = false,
          dataType,
          parentDataType,
        } = node;
        if (pathType === 2 && typeof mainKey === 'number' && mainKey > -1) {
          fields?.forEach((field: FieldItem, index: number) => {
            if (dataType === 2 && index === mainKey) {
              field.unique = true;
            }
          });
        }
        return pathType === 0
          ? {
              name,
              description,
              fields,
              alias,
              parentAlias,
              pathType: 0,
            }
          : {
              dataType,
              name,
              description,
              labelNames: tags?.map((tag: any) => tag.label || tag.value) || [],
              addFlow,
              addDashBoard,
              save2db,
              fields: fields?.filter(
                (i) => !(i?.systemField || (dataType === 1 && [qualityName, timestampName].includes(i?.name)))
              ),
              alias,
              parentAlias,
              pathType: 2,
              parentDataType,
            };
      });

      //校验提交树同级节点重复名称
      // if (!noDuplicates(newCheckedNodes.map((e) => e.path))) {
      //   message.error(formatMessage('uns.topicDuplicateTip'));
      //   return;
      // }
      try {
        setLoading(true);
        const res: any = await batchReverser(newCheckedNodes);
        if (res?.code === 200) {
          message.success(formatMessage('appGui.saveSuccess'));
          setLoading(false);
          close(true);
        } else {
          setLoading(false);
        }
      } catch (err) {
        const { data, code, msg }: any = err;
        if (code === 206) {
          message.error({
            type: 'error',
            content: (
              <div>
                <div>{formatMessage('common.partialFailure')}</div>
                {Object.keys(data).map((key: string) => {
                  const dataPath = newCheckedNodes[Number(key?.split('-')?.slice(-1)?.[0])]?.name;
                  return (
                    <div style={{ textAlign: 'left' }} key={key}>
                      {dataPath}：{data[key]}
                    </div>
                  );
                })}
              </div>
            ),
            duration: 5,
          });
          initTreeData({ reset: true });
          // if (data.length !== newCheckedNodes.filter((item) => !item.path.endsWith('/')).length) close(true);
        } else {
          message.error(msg);
        }
        setLoading(false);
      }
    };

    useEffect(() => {
      const newFileList: TreeNode[] = [];
      const getAllFiles = (treeData: TreeNode[]) => {
        treeData.forEach((node: TreeNode) => {
          if (node.children && node.children.length) {
            getAllFiles(node.children);
          } else {
            newFileList.push(node);
          }
        });
      };
      getAllFiles(treeData);
      setFileList(newFileList);
    }, [treeData]);

    useEffect(() => {
      if (isSave) {
        setBottomBtns(['prev', 'save']);
      } else {
        setBottomBtns(['next']);
      }
    }, [isSave]);

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

    const getDataTypeOptions = () => {
      if (enableAutoCategorization) {
        switch (currentParentDataType) {
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

    const renderContent = () => {
      return isSave ? (
        <>
          <Flex
            gap={10}
            style={{ height: fullScreen ? 'calc(100vh - 330px)' : enableAutoCategorization ? '350px' : '400px' }}
          >
            <JsonTree
              ref={jsonTreeRef}
              treeData={treeData}
              setTreeData={setTreeData}
              selectedInfo={selectedInfo}
              setSelectedInfo={setSelectedInfo}
            />
            {selectedInfo && (
              <div style={{ flex: 1, height: '100%', overflowY: 'auto', padding: '0 10px' }}>
                <Flex align="center" gap={5}>
                  <span style={{ flexShrink: 0, height: '16px' }}>
                    {selectedInfo.pathType === 2 ? <Document /> : <Folder />}
                  </span>
                  <span style={{ wordBreak: 'break-word', minHeight: '22px' }}>
                    {getTargetNode(treeData, selectedInfo.dataPath)?.name}
                  </span>
                </Flex>
                <Divider style={{ borderColor: '#c6c6c6', margin: '10px 0' }} />
                <Form.Item
                  name={['currentNode', 'name']}
                  label={formatMessage('common.name')}
                  validateTrigger={['onBlur', 'onChange']}
                  rules={[
                    { required: true, message: formatMessage('uns.pleaseInputName') },
                    { pattern: /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/, message: formatMessage('uns.nameFormat') },
                    {
                      max: 63,
                      message: formatMessage('uns.labelMaxLength', { label: formatMessage('common.name'), length: 63 }),
                    },
                    {
                      validator: (_, value) => {
                        if (selectedInfo.pathType === 0 && ['label', 'template'].includes(value)) {
                          return Promise.reject(new Error(formatMessage('uns.prohibitKeywords')));
                        } else {
                          return Promise.resolve();
                        }
                      },
                    },
                  ]}
                >
                  <Input />
                </Form.Item>
                <Form.Item
                  name={['currentNode', 'description']}
                  label={formatMessage(`uns.${selectedInfo?.pathType === 2 ? 'fileDescription' : 'folderDescription'}`)}
                  rules={[
                    {
                      max: 255,
                      message: formatMessage('uns.labelMaxLength', {
                        label: formatMessage(
                          `uns.${selectedInfo?.pathType === 2 ? 'fileDescription' : 'folderDescription'}`
                        ),
                        length: 255,
                      }),
                    },
                  ]}
                >
                  <TextArea rows={2} />
                </Form.Item>
                {selectedInfo.pathType === 2 ? (
                  <>
                    <Form.Item name={['currentNode', 'tags']} label={formatMessage('common.label')}>
                      <TagSelect />
                    </Form.Item>
                    <Divider style={{ borderColor: '#c6c6c6' }} />

                    {enableAutoCategorization && (
                      <Form.Item
                        name={['currentNode', 'parentDataType']}
                        label={formatMessage('uns.parentDataType')}
                        initialValue={1}
                      >
                        <ComRadio
                          options={[
                            { label: formatMessage('uns.state'), value: 1 },
                            { label: formatMessage('uns.metric'), value: 3 },
                          ]}
                          onChange={(e) => {
                            switch (e.target.value) {
                              case 1:
                                form.setFieldValue(['currentNode', 'dataType'], 2);
                                break;
                              case 3:
                                form.setFieldValue(['currentNode', 'dataType'], 1);
                                break;
                              default:
                                break;
                            }
                          }}
                        />
                      </Form.Item>
                    )}

                    <Form.Item name={['currentNode', 'dataType']} label={formatMessage('uns.databaseType')}>
                      <ComRadio options={getDataTypeOptions()} />
                    </Form.Item>
                    <Divider style={{ borderColor: '#c6c6c6' }} />
                    <Form.Item
                      name={['currentNode', 'save2db']}
                      label={formatMessage('uns.persistence')}
                      valuePropName="checked"
                    >
                      <ComCheckbox />
                    </Form.Item>
                    <Form.Item
                      name={['currentNode', 'addFlow']}
                      label={formatMessage('uns.mockData')}
                      valuePropName="checked"
                    >
                      <ComCheckbox />
                    </Form.Item>
                    {hasGrafana && (
                      <Form.Item
                        name={['currentNode', 'addDashBoard']}
                        label={formatMessage('uns.autoDashboard')}
                        valuePropName="checked"
                      >
                        <ComCheckbox />
                      </Form.Item>
                    )}
                  </>
                ) : (
                  <Divider style={{ borderColor: '#c6c6c6' }} />
                )}
                <FieldsFormList
                  types={types}
                  isCreateFolder={selectedInfo.pathType === 0}
                  showTooltip={false}
                  dataTypeName={['currentNode', 'dataType']}
                  fieldsName={['currentNode', 'fields']}
                  mainKeyName={['currentNode', 'mainKey']}
                />
              </div>
            )}
          </Flex>
          <Divider style={{ borderColor: '#c6c6c6' }} />
        </>
      ) : (
        <div style={{ position: 'relative', width: '100%' }}>
          <Form.Item
            name="jsonData"
            label=""
            wrapperCol={{ span: 24 }}
            rules={[{ required: true, validator: validatorJson }]}
            validateTrigger={['onBlur', 'onChange']}
          >
            <TextArea
              allowClear
              placeholder={exampleJson}
              style={{ height: fullScreen ? 'calc(100vh - 305px)' : '300px' }}
              onKeyDownCapture={(e) => {
                if (e.ctrlKey && e.code === 'Enter') {
                  if (jsonData) return;
                  e.preventDefault();
                  form.setFieldsValue({ jsonData: exampleJson });
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
      );
    };

    const renderButtons = () => {
      return bottomBtns.map((item) => {
        switch (item) {
          case 'prev':
            return (
              <Button
                color="primary"
                variant="filled"
                size="small"
                icon={<ChevronLeft />}
                onClick={prevStep}
                key={item}
              >
                {formatMessage('common.prev')}
              </Button>
            );
          case 'next':
            return (
              <Button
                color="primary"
                variant="filled"
                size="small"
                icon={<ChevronRight />}
                iconPosition="end"
                onClick={nextStep}
                key={item}
              >
                {formatMessage('common.next')}
              </Button>
            );
          case 'save':
            return (
              <Button color="primary" variant="solid" size="small" onClick={save} loading={loading} key={item}>
                {formatMessage('common.save')}
              </Button>
            );
          default:
            return null;
        }
      });
    };

    const allPersistenceChecked = fileList.length > 0 && fileList.every((e) => e.save2db);
    const allMockDataChecked = fileList.length > 0 && fileList.every((e) => e.addFlow);
    const allAutoDashboardChecked = fileList.length > 0 && fileList.every((e) => e.addDashBoard);

    const batchPersistence = () => {
      handleNodeInfo(treeData, 'save2db', !allPersistenceChecked);
      form.setFieldValue(['currentNode', 'save2db'], !allPersistenceChecked);
    };
    const batchMockData = () => {
      handleNodeInfo(treeData, 'addFlow', !allMockDataChecked);
      form.setFieldValue(['currentNode', 'addFlow'], !allMockDataChecked);
    };
    const batchAutoDashboard = () => {
      handleNodeInfo(treeData, 'addDashBoard', !allAutoDashboardChecked);
      form.setFieldValue(['currentNode', 'addDashBoard'], !allAutoDashboardChecked);
    };
    const batchModifyDataType = (type: number, parentDataType?: number) => {
      handleNodeInfo(treeData, 'dataType', type);
      if (parentDataType) {
        handleNodeInfo(treeData, 'parentDataType', parentDataType);
        form.setFieldValue(['currentNode', 'parentDataType'], parentDataType);
      }
      form.setFieldValue(['currentNode', 'dataType'], type);
    };

    const renderBatchChecks = (hasGrafana: boolean) => {
      return isSave ? (
        <Flex gap={8}>
          <ComCheckbox
            checked={allPersistenceChecked}
            indeterminate={!allPersistenceChecked && fileList.some((e) => e.save2db)}
            onChange={batchPersistence}
          >
            {formatMessage('uns.batchPersistence')}
          </ComCheckbox>
          <ComCheckbox
            checked={allMockDataChecked}
            indeterminate={!allMockDataChecked && fileList.some((e) => e.addFlow)}
            onChange={batchMockData}
          >
            {formatMessage('uns.batchMockData')}
          </ComCheckbox>
          {hasGrafana && (
            <ComCheckbox
              checked={allAutoDashboardChecked}
              indeterminate={!allAutoDashboardChecked && fileList.some((e) => e.addDashBoard)}
              onChange={batchAutoDashboard}
            >
              {formatMessage('uns.batchAutoDashboard')}
            </ComCheckbox>
          )}
        </Flex>
      ) : (
        <div />
      );
    };

    return (
      <>
        {renderContent()}
        <Flex justify="space-between">
          {renderBatchChecks(hasGrafana)}
          <Flex gap={10}>{renderButtons()}</Flex>
        </Flex>
      </>
    );
  }
);
export default JsonForm;
