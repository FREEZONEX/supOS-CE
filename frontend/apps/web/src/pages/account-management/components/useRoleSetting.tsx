import { useTranslate } from '@/hooks';
import { useEffect, useRef, useState } from 'react';
import { App, Button, ConfigProvider, Divider, Flex, Form, Input, Popover, Tabs } from 'antd';
import styles from './RoleSetting.module.scss';
import { getRoleList } from '@/apis/inter-api/user-manage.ts';
import { Add, Close, UserAvatar } from '@carbon/icons-react';
import { produce } from 'immer';
import Permission from '@/pages/account-management/components/Permission.tsx';
import { addRole, deleteRole, putRole } from '@/apis/inter-api/role.ts';
import { childrenRoutes } from '@/routers';
import { validSpecialCharacter } from '@/utils/pattern';
import Loading from '@/components/loading';
import ProModal from '@/components/pro-modal';
import { useBaseStore } from '@/stores/base';
import type { ResourceProps } from '@/stores/types.ts';

export const AdminRoleId = '7ca9f922-0d35-44cf-8747-8dcfd5e66f8e';
const generalRoleId = '71dd6dc2-6b12-4273-9ec0-b44b86e5b500';
const disabledRoleList = [AdminRoleId, generalRoleId];

const AddRoleContent = ({ successBack, disabled }: { successBack: (data: any) => void; disabled?: boolean }) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    if (!open) {
      form.resetFields();
    }
  }, [open]);

  const onSave = async () => {
    const info = await form.validateFields();
    setLoading(true);
    addRole({ name: info?.roleName })
      .then((data) => {
        setOpen(false);
        message.success(formatMessage('common.optsuccess'));
        successBack?.({ roleId: data?.roleId, roleName: info?.roleName });
      })
      .finally(() => {
        setLoading(false);
      });
  };
  return (
    <Popover
      open={open}
      onOpenChange={(open) => !open && setOpen(false)}
      styles={{
        body: {
          padding: '12px 0',
        },
      }}
      content={
        <div>
          <Flex justify="space-between" align="center" gap={8} style={{ padding: '0 12px' }}>
            <UserAvatar size={24} />
            <Form form={form}>
              <Form.Item
                name="roleName"
                style={{ padding: 0, margin: 0 }}
                rules={[
                  {
                    required: true,
                    message: formatMessage('rule.required'),
                  },
                  {
                    type: 'string',
                    min: 1,
                    max: 10,
                    message: formatMessage('rule.customCharacterLimit', { length: 10 }),
                  },
                  {
                    pattern: validSpecialCharacter,
                    message: formatMessage('common.SpecialCharacterValidation', {
                      rule: '~ ` ! @ # $ % ^ & * ( ) _ + = { } [ ] \\ | ; : \' " , < > . / ?',
                    }),
                  },
                ]}
              >
                <Input
                  allowClear
                  style={{ width: 140 }}
                  size="small"
                  placeholder={formatMessage('account.addRoleName')}
                />
              </Form.Item>
            </Form>
          </Flex>

          <Divider
            style={{
              background: 'var(--supos-t-dividr-color)',
              margin: '14px auto',
            }}
          />
          <Flex justify="flex-end" align="center" style={{ padding: '0 12px' }}>
            <Button
              loading={loading}
              type="primary"
              size="small"
              style={{ width: 60 }}
              onClick={onSave}
              title={formatMessage('common.save')}
            >
              {formatMessage('common.save')}
            </Button>
          </Flex>
        </div>
      }
      arrow={false}
      placement={'bottom'}
      trigger={['click']}
    >
      <Button
        disabled={disabled}
        title={disabled ? formatMessage('account.addRoleMax') : formatMessage('account.addRole')}
        style={{ height: 26 }}
        type="primary"
        onClick={() => setOpen(true)}
      >
        {formatMessage('account.addRole')}
        <Add size={16} />
      </Button>
    </Popover>
  );
};
export interface PermissionNode {
  id: string;
  showName: string;
  code?: string;
  type: number;
  checked: boolean;
  pagePermissionChecked?: boolean;
  actionPermissionChecked?: boolean;
  actionPermissionCheckedDisabled?: boolean;
  children?: PermissionNode[];
}

// 定义子组件暴露的 ref 类型
export interface PermissionRefProps {
  getValue: () => PermissionNode[];
  setValue: (value: PermissionNode[]) => void;
}

const getButtonGroup = (allButtonGroup: ResourceProps[], menuGroup: ResourceProps[]) => {
  const result: ResourceProps[] = [];
  allButtonGroup.forEach((item) => {
    if (item.parentId) {
      let parent = result.find((p) => p.id === item.parentId);
      if (!parent) {
        const menuParent = menuGroup.find((p) => p.id === item.parentId)!;
        // 如果父项不存在，创建并添加到结果
        parent = {
          ...menuParent,
          children: [],
        };
        result.push(parent);
      }
      parent.children?.push({ ...item, checked: false, id: 'button:' + item.code });
    }
  });
  return result;
};

const getAllMenuTree = (originMenu: ResourceProps[] = []) => {
  // 创建一个映射表，用于快速查找节点
  const map: { [key: string]: ResourceProps } = {};
  const tree: ResourceProps[] = [];
  const menu: ResourceProps[] = [];

  // 首先将所有节点存入映射表，以id为键
  originMenu.forEach((item) => {
    const parent = originMenu?.find((f) => f.id === item.parentId);
    if (!parent || parent?.type === 1) {
      map[item.id] = { ...item, children: [] };
      menu.push(item);
    }
  });
  // 遍历所有节点，根据parentId构建树
  menu.forEach((item) => {
    const node = map[item.id];

    if (!item.parentId) {
      // 如果没有parentId或parentId为0/null/undefined，则认为是根节点
      tree.push(node);
    } else {
      // 否则找到父节点，将当前节点添加到父节点的children中
      const parent = map[item.parentId];
      if (parent) {
        parent.children!.push(node);
        parent.children!.sort((a, b) => a.sort - b.sort);
      }
    }
  });
  return tree.filter((f) => f.type === 2 || (f.type === 1 && f?.children?.length)).sort((a, b) => a.sort - b.sort);
};
const useRoleSetting = ({ onSaveBack }: any) => {
  const formatMessage = useTranslate();
  const { originMenu, allButtonGroup } = useBaseStore((state) => ({
    allButtonGroup: state.allButtonGroup,
    originMenu: state.originMenu,
  }));
  const [open, setOpen] = useState(false);
  const [items, setItems] = useState<any[]>([]);
  const { message } = App.useApp();
  // 初始数据
  const initItems = useRef<any[]>([]);
  // 初始的菜单按钮配置
  const initialRolePermissionData = useRef<any[]>([]);
  // 所有button
  const allButtonData = useRef<any[]>([]);
  // 跟踪每个标签页的保存状态
  const unsavedChanges = useRef<Map<string, boolean>>(new Map());
  const [loading, setLoading] = useState(false);
  const [activeKey, setActiveKey] = useState('');
  const permissionRefs = useRef<Map<string, PermissionRefProps | null>>(new Map());

  const { modal } = App.useApp();
  const onRoleModalOpen = () => {
    setOpen(true);
  };

  const onClose = () => {
    const hasChanges = [...unsavedChanges.current.values()].some(Boolean);
    if (hasChanges) {
      modal.confirm({
        title: formatMessage('common.unsavedChanges'),
        okText: formatMessage('common.save'),
        cancelText: formatMessage('common.unSave'),
        onOk: () => {
          onSave();
          setOpen(false);
        },
        onCancel: () => {
          setOpen(false);
        },
        okButtonProps: {
          title: formatMessage('common.save'),
        },
        cancelButtonProps: {
          title: formatMessage('common.unSave'),
        },
      });
    } else {
      setOpen(false);
    }
  };
  useEffect(() => {
    if (open) {
      const menuGroup = originMenu?.filter((i) => i.type !== 3 && i.enable);
      const menuTree = getAllMenuTree(menuGroup);
      getRoleList().then((role) => {
        const buttons = getButtonGroup(allButtonGroup, menuGroup);
        initialRolePermissionData.current = mapInitialRolePermissionData(menuTree, buttons);
        allButtonData.current = extractButtonIds(buttons);
        const info =
          role?.map?.((i: any) => {
            const denyResourceButtonList = i.denyResourceList?.filter((f: any) => f.uri?.includes('button:')) || [];
            const resourceButtonList = i?.resourceList?.some((i: any) => i.uri?.includes('button:'))
              ? (allButtonData.current.filter((f: any) => !denyResourceButtonList.some((s: any) => s.uri === f)) ?? [])
              : [];
            return {
              ...i,
              resourceList: updatePermissionData(
                initialRolePermissionData.current,
                [...(i?.resourceList?.map((item: any) => item.uri) ?? []), ...resourceButtonList],
                i.roleId === AdminRoleId
              ),
            };
          }) || [];
        setItems(info);
        initItems.current = info;
        setActiveKey(role?.[0]?.roleId);
      });
    }
  }, [open, originMenu]);

  const onSave = () => {
    setLoading(true);
    const newValue = permissionRefs.current.get(activeKey)?.getValue?.();
    const { checkedFalseButtons, checkedTrueMenus } = filterMenuAndButtonItems(newValue);
    const allButton = checkedFalseButtons?.length === 0;
    putRole({
      id: activeKey,
      name: items?.find((i: any) => i.roleId === activeKey)?.roleName,
      denyResourceList: allButton ? [] : checkedFalseButtons?.map((item) => ({ uri: item })),
      allowResourceList: [...(checkedTrueMenus?.map?.((item) => ({ uri: item })) ?? []), { uri: 'button:*' }],
    })
      .then(() => {
        message.success(formatMessage('common.optsuccess'));
        setItems(
          produce(items, (draft) => {
            const info = draft.find((todo) => todo.roleId === activeKey);
            if (info) {
              info['resourceList'] = newValue;
            }
          })
        );
        initItems.current = initItems.current.map((item) => {
          if (item.roleId === activeKey) {
            return {
              ...item,
              resourceList: newValue,
            };
          } else {
            return item;
          }
        });
        unsavedChanges.current.set(activeKey, false);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const onChange = (key: string) => {
    if (unsavedChanges.current.get(activeKey)) {
      modal.confirm({
        title: formatMessage('common.unsavedChanges'),
        okText: formatMessage('common.save'),
        cancelText: formatMessage('common.unSave'),
        onOk: () => {
          onSave();
          setActiveKey(key);
        },
        onCancel: () => {
          const initPermission = initItems.current.find((item) => item.roleId === activeKey);
          permissionRefs.current.get(activeKey)?.setValue(initPermission?.resourceList);
          unsavedChanges.current.set(activeKey, false);
          setActiveKey(key);
        },
        okButtonProps: {
          title: formatMessage('common.save'),
        },
        cancelButtonProps: {
          title: formatMessage('common.unSave'),
        },
      });
    } else {
      setActiveKey(key);
    }
  };

  const RoleModal = (
    <ProModal
      afterClose={() => {
        unsavedChanges.current.clear();
        permissionRefs.current.clear();
        setItems([]);
        initItems.current = [];
      }}
      className={styles['role-setting']}
      size="sm"
      open={open}
      maskClosable={false}
      onCancel={onClose}
      title={
        <Flex justify="space-between" align="center" style={{ height: '100%' }}>
          <span>{formatMessage('account.roleSettings')}</span>
          <AddRoleContent
            successBack={(data) => {
              setItems((items) => {
                const newItems = [
                  ...items,
                  {
                    ...data,
                    resourceList: updatePermissionData(initialRolePermissionData.current, []),
                  },
                ];
                initItems.current = newItems;
                return newItems;
              });
              setActiveKey(data?.roleId);
            }}
            disabled={items?.length >= 10}
          />
        </Flex>
      }
    >
      <ConfigProvider
        theme={{
          components: {
            Tabs: {
              itemSelectedColor: 'var(--supos-theme-color)',
              zIndexPopup: 9999,
              horizontalMargin: '0 0 0 0',
            },
            Dropdown: {
              colorText: '#000',
            },
          },
        }}
      >
        <Loading spinning={loading}>
          <Tabs
            more={{
              overlayStyle: { '--supos-text-color': '#000' },
            }}
            onChange={onChange}
            activeKey={activeKey}
            items={items?.map((item: any) => {
              return {
                label: (
                  <Flex justify="space-between" align="center" gap={8}>
                    {item.roleName}
                    {!disabledRoleList.includes(item.roleId) && (
                      <Close
                        style={{ cursor: 'pointer' }}
                        onClick={(e) => {
                          e.stopPropagation();
                          modal.confirm({
                            title: formatMessage('common.deleteConfirm'),
                            onOk: async () => {
                              return await deleteRole(item?.roleId).then(() => {
                                message.success(formatMessage('common.deleteSuccessfully'));
                                onSaveBack?.();
                                setItems(
                                  produce(items, (draft) => {
                                    const index = draft.findIndex((todo) => todo.roleId === item.roleId);
                                    if (index !== -1) {
                                      draft.splice(index, 1);
                                      if (activeKey === item.roleId) {
                                        setActiveKey(draft.filter((todo) => todo.roleId !== item.roleId)?.[0]?.roleId);
                                      }
                                    }
                                  })
                                );
                              });
                            },
                            okButtonProps: {
                              title: formatMessage('common.confirm'),
                            },
                            cancelButtonProps: {
                              title: formatMessage('common.cancel'),
                            },
                          });
                        }}
                      />
                    )}
                  </Flex>
                ),
                key: item.roleId,
                children: (
                  <Permission
                    disabled={disabledRoleList.includes(item.roleId)}
                    ref={(el) => permissionRefs.current.set(item.roleId, el)}
                    initValue={item.resourceList}
                    onChange={(pre) => {
                      const currentDataString = JSON.stringify(pre);
                      const hasChanges =
                        currentDataString !==
                        JSON.stringify(initItems?.current?.find((i) => i.roleId === item.roleId)?.resourceList);
                      unsavedChanges.current.set(item.roleId, hasChanges);
                    }}
                  />
                ),
              };
            })}
          />
          <Button
            disabled={disabledRoleList.includes(activeKey)}
            onClick={onSave}
            style={{ height: 32, marginTop: 20 }}
            block
            type="primary"
            loading={loading}
            title={formatMessage('common.save')}
          >
            {formatMessage('common.save')}
          </Button>
        </Loading>
      </ConfigProvider>
    </ProModal>
  );

  return {
    RoleModal,
    onRoleModalOpen,
  };
};

// 根据前端维护的路由
const getOtherRoutes = () => {
  return childrenRoutes
    ?.filter((route) => {
      return route?.handle?.parentPath === '/_common' && !['dev', 'all'].includes(route?.handle?.type || '');
    })
    ?.map((route) => {
      const children = [
        {
          showName: route?.handle?.showName,
          code: route?.handle?.code,
          url: route?.path,
          isFrontend: true,
          isRemote: false,
          type: 2,
          urlType: 1,
        },
      ];
      return {
        ...children[0],
        children,
      };
    });
};

// 根据前端维护的button、路由以及kong维护的路由进行初始化数据整合
const mapInitialRolePermissionData = (routes: any, buttonGroup: ResourceProps[]) => {
  return [...routes, ...getOtherRoutes()]?.map((group: any) => {
    // id就是路由 keycloke配置的路由
    return {
      ...group,
      id: 'group' + group?.code,
      type: 1,
      checked: false,
      pagePermissionChecked: false,
      actionPermissionChecked: false,
      children:
        group?.children?.length > 0
          ? group?.children?.map((menu: any) => {
              const buttonList = buttonGroup?.find((f) => f.id === menu.id)?.children;
              return {
                ...menu,
                id: menu?.urlType === 1 ? menu?.url : '/' + menu?.code,
                checked: false,
                children: buttonList,
              };
            })
          : [
              {
                ...group,
                id: group?.urlType === 1 ? group?.url : '/' + group?.code,
                checked: false,
                children: buttonGroup?.find((f) => f.id === group.id)?.children,
              },
            ],
    };
  });
};

// 过滤出未选中的button和选中的menu
function filterMenuAndButtonItems(data: PermissionNode[] = []) {
  const result = {
    checkedTrueMenus: [] as string[], // type: "menu" && checked: true
    checkedFalseButtons: [] as string[], // type: "button" && checked: false
  };

  function traverse(items: PermissionNode[]) {
    items.forEach((item: PermissionNode) => {
      // 处理当前项
      if (item.type === 2 && item.checked) {
        result.checkedTrueMenus.push(item.id);
      } else if (item.type === 3 && !item.checked) {
        result.checkedFalseButtons.push(item.id);
      }

      // 递归处理子项
      if (item.children && item.children.length) {
        traverse(item.children);
      }
    });
  }

  traverse(data);
  return result;
}

// 过滤所有buttons的id
function extractButtonIds(data: any[]): string[] {
  const buttonIds: string[] = [];

  function traverse(nodes: PermissionNode[]) {
    nodes.forEach((node) => {
      if (node.type === 3) {
        buttonIds.push(node.id);
      }
      if (node.children && node.children.length) {
        traverse(node.children);
      }
    });
  }

  traverse(data);
  return buttonIds;
}

// 回显值
function updatePermissionData(data: any, idArray: string[], isAdmin: boolean = false) {
  const newData = JSON.parse(JSON.stringify(data));
  function updateChecked(items: any) {
    if (!items || !Array.isArray(items)) return;

    items.forEach((item: any) => {
      // 如果是管理员角色，直接设置所有节点为选中状态
      // 否则，检查当前项的id是否在idArray中
      if (isAdmin) {
        item.checked = true;
      } else if (idArray?.includes(item.id)) {
        item.checked = true;
      }

      if (item.children && item.children.length) {
        updateChecked(item.children);
      }
    });

    items.forEach((item: any) => {
      // 如果是group类型，检查并更新其状态
      if (item.type === 1) {
        const menuNodes = item.children?.filter((child: any) => child.type === 2) || [];
        const buttonNodes =
          item.children?.flatMap((child: any) => child.children || []).filter((child: any) => child.type === 3) || [];

        // 检查并更新pagePermissionChecked - 只受menu类型节点影响
        const allMenuChecked = menuNodes.length > 0 && menuNodes.every((menu: any) => menu.checked === true);
        // 如果任何菜单被选中，则页面权限部分选中
        item.pagePermissionChecked = allMenuChecked;

        // 检查并更新actionPermissionChecked - 只受button类型节点影响
        if (buttonNodes.length === 0) {
          item.actionPermissionChecked = false; // 设置默认值为false，确保禁用状态下始终为未选中
          item.actionPermissionCheckedDisabled = true; // 添加禁用标志
        } else {
          const allButtonsChecked = buttonNodes.every((button: any) => button.checked === true);
          item.actionPermissionChecked = allButtonsChecked;
          item.actionPermissionCheckedDisabled = false; // 有按钮时不禁用
        }
        item.checked =
          (buttonNodes.length === 0 && allMenuChecked) ||
          (menuNodes.length === 0 &&
            buttonNodes.length > 0 &&
            buttonNodes.every((button: any) => button.checked === true)) ||
          (menuNodes.length > 0 &&
            buttonNodes.length > 0 &&
            allMenuChecked &&
            buttonNodes.every((button: any) => button.checked === true));
      }
    });
  }

  updateChecked(newData);
  return newData;
}

export default useRoleSetting;
