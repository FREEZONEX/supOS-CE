import { App, Flex, Form, Tag } from 'antd';
import useTranslate from '@/hooks/useTranslate.ts';
import BasicInfo, { type FormItemType } from './BasicInfo.tsx';
import OperationInfo from './OperationInfo.tsx';
import { AuthButton } from '@/components/auth';
import { batchDeleteResourceApi, batchEditResourceApi, postResourceApi } from '@/apis/inter-api/resource.ts';
import { useMenuStore } from '../../store/menuStore.tsx';
import { useEffect, useState } from 'react';
import { uploadAttachment } from '@/apis/inter-api';
import { passwordRegex } from '@/utils';
import {
  CUSTOM_MENU_ICON,
  CUSTOM_MENU_ICON_PRE,
  CUSTOM_MENU_ICON_PRE1,
  MENU_TARGET_PATH,
  STORAGE_PATH,
} from '@/common-types/constans.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';

const Title = () => {
  const name = Form.useWatch('name');
  const id = Form.useWatch('id');
  const formatMessage = useTranslate();
  return (
    <Flex style={{ overflow: 'hidden' }} align="center" justify="space-between" title={name}>
      <span style={{ fontSize: 30, overflow: 'hidden', whiteSpace: 'nowrap', textOverflow: 'ellipsis' }}>{name}</span>
      {!id && (
        <span>
          <Tag color="success">{formatMessage('common.new')}</Tag>
        </span>
      )}
    </Flex>
  );
};
const MenuContent = () => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const { requestMenu, setContentType, setSelectNode, selectNode, contentType, menuTree } = useMenuStore((state) => ({
    requestMenu: state.requestMenu,
    setContentType: state.setContentType,
    setSelectNode: state.setSelectNode,
    selectNode: state.selectNode,
    contentType: state.contentType,
    menuTree: state.menuTree,
  }));
  const [configs, setConfigs] = useState<FormItemType[]>([]);

  useEffect(() => {
    form.resetFields();
    // 是否是编辑
    const isEdit = ['editGroup', 'editMenu'].includes(contentType || '');
    // 设置表单项
    setConfigs(
      ['addMenu', 'editMenu'].includes(contentType || '')
        ? [
            {
              formType: 'sourceSelect',
              formProps: {
                name: 'source',
                label: formatMessage('MenuConfiguration.menuSource'),
                initialValue: {
                  routeSource: 1,
                  route: null,
                },
                validateTrigger: ['onBlur', 'onChange'],
                rules: [
                  { required: true, message: formatMessage('rule.required') },
                  {
                    validator(_, value) {
                      if (value.routeSource === 1) {
                        return Promise.resolve();
                      }
                      if (value.routeSource === 2 && value?.route) {
                        return Promise.resolve();
                      }
                      return Promise.reject(formatMessage('rule.required'));
                    },
                    validateTrigger: 'onBlur',
                  },
                ],
              },
              childProps: {
                onChange: (value: any) => {
                  form.setFieldsValue({
                    code: value?.route?.name,
                    url: value?.route?.url,
                  });
                },
              },
            },
            {
              formType: 'codeInput',
              formProps: {
                name: 'code',
                label: formatMessage('MenuConfiguration.menuCode'),
                rules: [
                  { required: true, message: formatMessage('rule.required') },
                  {
                    max: 255,
                    message: formatMessage('uns.labelMaxLength', {
                      label: formatMessage('MenuConfiguration.menuCode'),
                      length: 255,
                    }),
                  },
                  { pattern: passwordRegex, message: formatMessage('rule.password') },
                ],
              },
              childProps: {
                disabled: isEdit,
              },
            },
            {
              formType: 'input',
              formProps: {
                name: 'name',
                label: formatMessage('MenuConfiguration.menuName'),
                rules: [
                  { required: true, message: formatMessage('rule.required') },
                  {
                    max: 255,
                    message: formatMessage('uns.labelMaxLength', {
                      label: formatMessage('MenuConfiguration.menuName'),
                      length: 255,
                    }),
                  },
                ],
              },
            },
            {
              formType: 'uploadPicture',
              formProps: {
                name: 'iconFile',
                label: formatMessage('MenuConfiguration.menuIcon'),
              },
              childProps: {
                onChange: (fileList: any) => {
                  console.log(fileList);
                  if (fileList.length > 0) {
                    uploadAttachment(
                      fileList?.map((item: any) => ({ value: item?.file, name: 'files', fileName: item?.file?.name })),
                      { alias: '__templates__' }
                    ).then((data) => {
                      form.setFieldValue('icon', data?.list?.[0]?.attachmentPath);
                    });
                  } else {
                    form.setFieldValue('icon', undefined);
                  }
                },
              },
            },
            {
              formType: 'radioGroup',
              formProps: {
                initialValue: 1,
                name: 'openType',
                label: formatMessage('MenuConfiguration.openMode'),
              },
              childProps: {
                options: [
                  { label: formatMessage('MenuConfiguration.openCurrentTab'), value: 0 },
                  { label: formatMessage('MenuConfiguration.openNewTab'), value: 1 },
                ],
              },
            },
            {
              formType: 'input',
              formProps: {
                name: 'url',
                label: formatMessage('MenuConfiguration.menuUrl'),
                rules: [{ required: true, message: formatMessage('rule.required') }],
              },
            },
            {
              formType: 'radioGroup',
              formProps: {
                hidden: true,
                initialValue: 2,
                name: 'urlType',
                label: 'urlType',
              },
              childProps: {
                options: [
                  { label: '内部地址', value: 1 },
                  { label: '外部地址', value: 2 },
                ],
              },
            },
            {
              formType: 'textArea',
              formProps: {
                name: 'description',
                label: formatMessage('MenuConfiguration.menuDescription'),
              },
              childProps: {
                row: 5,
              },
            },
            {
              formType: 'checkbox',
              formProps: {
                name: 'homeEnable',
                label: formatMessage('MenuConfiguration.homepageDisplay'),
                initialValue: true,
                valuePropName: 'checked',
              },
            },
          ]
        : [
            {
              formType: 'codeInput',
              formProps: {
                name: 'code',
                label: formatMessage('MenuConfiguration.menuCode'),
                rules: [
                  { required: true, message: formatMessage('rule.required') },
                  { pattern: passwordRegex, message: formatMessage('rule.password') },
                ],
              },
              childProps: {
                disabled: isEdit,
              },
            },
            {
              formType: 'input',
              formProps: {
                name: 'name',
                label: formatMessage('MenuConfiguration.menuName'),
                rules: [{ required: true, message: formatMessage('rule.required') }],
              },
            },
            {
              formType: 'uploadPicture',
              formProps: {
                name: 'iconFile',
                label: formatMessage('MenuConfiguration.menuIcon'),
              },
              childProps: {
                onChange: (fileList: any) => {
                  if (fileList.length > 0) {
                    uploadAttachment(
                      fileList?.map((item: any) => ({ value: item?.file, name: 'files', fileName: item?.file?.name })),
                      { alias: '__templates__' }
                    ).then((data) => {
                      form.setFieldValue('icon', data?.list?.[0]?.attachmentPath);
                    });
                  } else {
                    form.setFieldValue('icon', undefined);
                  }
                },
              },
            },
            {
              formType: 'textArea',
              formProps: {
                name: 'description',
                label: formatMessage('MenuConfiguration.menuDescription'),
              },
              childProps: {
                row: 5,
              },
            },
            {
              formType: 'checkbox',
              formProps: {
                name: 'homeEnable',
                label: formatMessage('MenuConfiguration.homepageDisplay'),
                initialValue: true,
                valuePropName: 'checked',
              },
            },
          ]
    );
    if (selectNode && ['editMenu', 'editGroup'].includes(contentType || '')) {
      // 赋值
      const url = selectNode?.icon
        ? selectNode?.icon?.includes(CUSTOM_MENU_ICON_PRE) || selectNode?.icon?.includes(CUSTOM_MENU_ICON_PRE1)
          ? `${CUSTOM_MENU_ICON}?objectName=${encodeURI(selectNode.icon)}`
          : `${STORAGE_PATH}${MENU_TARGET_PATH}/${encodeURI(selectNode.icon)}`
        : '';
      form.setFieldsValue({
        ...selectNode,
        type: contentType === 'editGroup' ? 1 : 2,
        name: selectNode.showName,
        description: selectNode.showDescription,
        // 来源
        source: {
          routeSource: selectNode.routeSource,
          route:
            selectNode.routeSource === 2
              ? {
                  name: selectNode.code,
                }
              : null,
        },
        // 文件
        iconFile: selectNode.icon
          ? [
              {
                uid: '-1',
                name: 'icon.svg',
                url,
                status: 'done',
              },
            ]
          : undefined,
        operationChildren:
          selectNode?.operationChildren?.map((m) => {
            return {
              ...m,
              name: m.showName,
              description: m.showDescription,
            };
          }) || [],
      });
    } else {
      form.setFieldsValue({
        sort: selectNode?.children ? selectNode?.children?.length + 1 : (menuTree?.length ?? 0) + 1,
        type: contentType === 'addGroup' ? 1 : 2,
        parentId: selectNode?.id,
      });
    }
  }, [contentType, selectNode, menuTree]);

  const onSave = async () => {
    const info = await form.validateFields();
    setLoading(true);
    const delIds = info.delOperationChildren?.filter((f: any) => !f.id.includes('add_')).map((m: any) => m.id);
    if (info?.id && delIds?.length > 0) {
      // 删除权限
      batchDeleteResourceApi(delIds);
    }
    const fn = () => {
      postResourceApi({
        ...info,
        iconFile: undefined,
        routeSource: info?.source?.routeSource,
        source: undefined,
        children: info?.operationChildren?.map((m: any, index: number) => ({
          id: m.id.includes('add_') ? undefined : m.id,
          code: m.code,
          name: m.name,
          description: m.description,
          type: 3,
          sort: index + 1,
        })),
      })
        .then(() => {
          setContentType(null);
          setSelectNode(null);
          requestMenu();
          message.success(info.id ? formatMessage('uns.editSuccessful') : formatMessage('uns.newSuccessfullyAdded'));
        })
        .finally(() => {
          setLoading(false);
        });
    };
    if (!info?.parentId && menuTree && menuTree?.length > 0) {
      const lastInfo = menuTree[menuTree?.length - 1];
      // 更新最后一组sort
      batchEditResourceApi([
        {
          id: lastInfo.id,
          sort: menuTree?.length + 2,
        },
      ]).then(() => {
        fn();
      });
    } else {
      fn();
    }
  };
  return (
    <Form
      labelWrap
      style={{ padding: 16, overflow: 'hidden', height: '100%', display: 'flex', flexDirection: 'column' }}
      form={form}
      colon={false}
      labelCol={{ span: 4 }}
      wrapperCol={{ span: 10 }}
      labelAlign="left"
      disabled={selectNode?.editEnable === false}
    >
      <Flex gap={8} align="center" style={{ marginBottom: 16, flexShrink: 0, height: 48 }} justify="space-between">
        <Flex gap={8} align="center" style={{ overflow: 'hidden' }}>
          <Title />
          <span>
            <Tag color="success">
              level{' '}
              {(selectNode?.depth ?? 0) +
                (selectNode?.type === 4 ? 2 : selectNode && contentType?.includes('add') ? 2 : 1)}
            </Tag>
          </span>
        </Flex>
        <AuthButton
          type="primary"
          auth={contentType?.includes('edit') ? ButtonPermission['MenuConfiguration.editMenu'] : undefined}
          onClick={onSave}
          loading={loading}
          disabled={selectNode?.editEnable === false}
        >
          {formatMessage('common.save')}
        </AuthButton>
      </Flex>
      <div style={{ flex: 1, overflow: 'auto' }}>
        <BasicInfo configs={configs} />
        {selectNode?.type === 2 || selectNode?.type === 4 ? <OperationInfo /> : null}
      </div>
    </Form>
  );
};

export default MenuContent;
