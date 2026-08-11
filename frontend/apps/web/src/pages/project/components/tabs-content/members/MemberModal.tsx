import { addProjectMember, getProjectRoles, updateProjectMember } from '@/apis/core-api';
import OperationForm from '@/components/operation-form';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import {
  MAX_CONCURRENT_PAGE_REQUESTS,
  PAGE_SIZE,
  useProjectCommonStore,
} from '@/pages/project/store/projectCommonStore';
import { type ProjectRole } from '@/pages/project/types';
import { fetchAllPagesWithConcurrency } from '@/utils/common';
import { App, Form } from 'antd';
import { forwardRef, useCallback, useEffect, useImperativeHandle, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router';
import styles from './index.module.scss';

export interface MemberModalRef {
  onOpen: (member?: any) => void;
  onClose: () => void;
}

export interface MemberModalProps {
  refreshRequest?: () => void;
  /** 已存在成员的 userId 列表，添加新成员时用于过滤下拉选项 */
  existingMemberUserIds?: string[];
}

type MemberFormValues = {
  id: string;
  userId: string;
  roles: string[];
};

type SelectOption = {
  label: string;
  value: string | number;
};

const MemberModal = forwardRef<MemberModalRef, MemberModalProps>(({ refreshRequest, existingMemberUserIds = [] }, ref) => {
  const { projectId: id } = useParams();
  const [visible, setVisible] = useState(false);
  const [editingMember, setEditingMember] = useState<any>(null);
  const [roleOptions, setRoleOptions] = useState<SelectOption[]>([]);
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message } = App.useApp();
  const roleOptionsRef = useRef<SelectOption[]>([]);
  const roleOptionsLoadedProjectIdRef = useRef<number | null>(null);
  const memberOptions = useProjectCommonStore((state) => state.memberOptions);
  const memberOptionsLoading = useProjectCommonStore((state) => state.memberOptionsLoading);
  const fetchMemberOptions = useProjectCommonStore((state) => state.fetchMemberOptions);
  const resetMemberOptions = useProjectCommonStore((state) => state.resetMemberOptions);

  // KeepAlive 下只有「首次进入」和「刷新 tab」才会触发 mount
  // 切换 tab 不会触发，因此 reset 只在真正需要时执行
  useEffect(() => {
    resetMemberOptions();
  }, []);

  const projectId = Number(id);

  const loadRoleOptions = useCallback(async () => {
    if (roleOptionsLoadedProjectIdRef.current === projectId) {
      return roleOptionsRef.current;
    }

    // 非法 projectId 时直接返回空，避免页面在切换路由过程中触发无效请求。
    if (!projectId || Number.isNaN(projectId)) {
      roleOptionsLoadedProjectIdRef.current = null;
      roleOptionsRef.current = [];
      return [];
    }

    const result = await fetchAllPagesWithConcurrency<ProjectRole>(
      async ({ pageNo, pageSize }) => {
        const response: any = await getProjectRoles(projectId, {
          pageNo,
          pageSize,
        });

        if (Array.isArray(response)) {
          return {
            total: response.length,
            data: response,
          };
        }

        const data = (Array.isArray(response?.data) ? response.data : []) as ProjectRole[];
        return {
          total: Number(response?.total ?? data.length),
          data,
        };
      },
      {
        pageSize: PAGE_SIZE,
        maxConcurrency: MAX_CONCURRENT_PAGE_REQUESTS,
      }
    );

    const roleMap = new Map<number, string>();
    result.data.forEach((role) => {
      if (!role?.roleId) {
        return;
      }

      roleMap.set(role.roleId, role.roleName || role.roleKey);
    });

    const options: SelectOption[] = Array.from(roleMap.entries()).map(([roleId, roleName]) => ({
      label: roleName,
      value: String(roleId),
    }));

    roleOptionsLoadedProjectIdRef.current = projectId;
    roleOptionsRef.current = options;
    return options;
  }, [projectId]);

  const renderLabel = (messageKey: string) => <span className={styles.labelStyle}>{formatMessage(messageKey)}</span>;

  const existingMemberUserIdSet = useMemo(
    () => new Set(existingMemberUserIds.map((id) => String(id))),
    [existingMemberUserIds]
  );

  const availableMemberOptions = useMemo(() => {
    // 编辑时不过滤，保持原始选项
    if (editingMember) {
      return memberOptions;
    }
    return memberOptions.filter((option) => !existingMemberUserIdSet.has(String(option.value)));
  }, [editingMember, existingMemberUserIdSet, memberOptions]);

  const onOpen = (member?: any) => {
    setEditingMember(member);
    void fetchMemberOptions();

    void loadRoleOptions()
      .then((options) => {
        setRoleOptions(options);

        if (member) {
          form.setFieldsValue({
            id: member.id,
            userId: member.userId,
            email: member.email,
            roles: member.roles || [],
          });
        }
      })
      .catch(() => {
        setRoleOptions([]);
      });

    setVisible(true);
  };

  const onClose = () => {
    form.resetFields();
    setEditingMember(null);
    setVisible(false);
  };

  const onSave = async () => {
    const isEditing = Boolean(editingMember);
    const values = (await form.validateFields()) as MemberFormValues;
    try {
      if (!projectId || Number.isNaN(projectId)) {
        return;
      }
      const userId = isEditing ? editingMember.userId : String(values.userId || '').trim();
      const payload = {
        userId,
        userName: memberOptions.find((m) => m.value == userId)?.label || '',
        roleIds: values?.roles?.map(Number),
      };

      if (!payload.userId) {
        throw new Error('member userId is required');
      }

      if (isEditing) {
        const memberId = Number(editingMember?.id);
        if (!memberId || Number.isNaN(memberId)) {
          throw new Error('member id is required');
        }
        await updateProjectMember(projectId, memberId, payload);
      } else {
        await addProjectMember(projectId, payload);
      }

      onClose();
      refreshRequest?.();
      message.success(formatMessage('common.optsuccess'));
    } catch (error: any) {
      if (error?.code == 100010) {
        message.error(formatMessage('project.memberAlreadyExists'));
      } else {
        message.error(
          (error as Error)?.message || isEditing
            ? formatMessage('project.memberEditFail')
            : formatMessage('project.memberAddFail')
        );
      }
    }
  };

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  const formOptions = [
    {
      name: 'userId',
      label: renderLabel('project.memberName'),
      type: 'Select',
      rules: [{ required: true, message: formatMessage('project.memberRequired') }],
      properties: {
        placeholder: formatMessage('project.memberNamePlaceholder'),
        allowClear: !editingMember,
        showSearch: !editingMember,
        optionFilterProp: 'label',
        options: availableMemberOptions,
        loading: memberOptionsLoading,
        disabled: !!editingMember,
        debounceTimeout: 300,
      },
    },
    {
      name: 'roles',
      label: renderLabel('project.role'),
      type: 'Select',
      properties: {
        placeholder: formatMessage('project.selectRole'),
        options: roleOptions,
        optionFilterProp: 'label',
        mode: 'multiple',
      },
    },
  ];

  return (
    <ProModal
      open={visible}
      onCancel={onClose}
      maskClosable={false}
      title={formatMessage(editingMember ? 'project.editMember' : 'project.addMember')}
      width={500}
      fullScreenable={false}
      styles={{
        body: { paddingBlockStart: 0 },
      }}
    >
      {() => {
        return (
          <OperationForm
            form={form}
            onCancel={onClose}
            onSave={onSave}
            formConfig={{
              layout: 'vertical',
              labelCol: { span: 24 },
              wrapperCol: { span: 24 },
              requiredMark: true,
            }}
            formItemOptions={formOptions}
          />
        );
      }}
    </ProModal>
  );
});

MemberModal.displayName = 'MemberModal';

export default MemberModal;
