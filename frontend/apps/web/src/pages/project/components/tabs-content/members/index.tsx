import { App, Button } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router';

import { deleteProjectMember, getProjectMembers } from '@/apis/core-api';
import type { PageProps } from '@/common-types';
import ProTable from '@/components/pro-table';
import { useActivate } from '@/contexts/tabs-lifecycle-context.ts';
import { useTranslate } from '@/hooks';
import type { ProjectMember } from '@/pages/project/types';

import TabLayout from '../../TabLayout';
import MemberModal, { type MemberModalRef } from './MemberModal';
import { createMembersColumns, type MemberItem } from './config.tsx';
import styles from './index.module.scss';

const Members: FC<PageProps> = () => {
  const { message, modal } = App.useApp();
  const { projectId: id } = useParams();
  const formatMessage = useTranslate();
  const memberModalRef = useRef<MemberModalRef>(null);
  const [editingMember, setEditingMember] = useState<MemberItem | null>(null);
  const [members, setMembers] = useState<MemberItem[]>([]);
  const [loading, setLoading] = useState(false);

  const projectId = Number(id);

  const getMembers = useCallback(async () => {
    if (!projectId || Number.isNaN(projectId)) {
      setMembers([]);
      return;
    }

    setLoading(true);
    try {
      const response = await getProjectMembers(projectId);
      const memberList: ProjectMember[] = Array.isArray(response) ? response : [];

      const normalizedMemberList: MemberItem[] = memberList.map((member) => {
        const roleIdSet = new Set<string>();
        const roleNameSet = new Set<string>();
        (Array.isArray(member.roles) ? member.roles : []).forEach((role) => {
          const roleId = role?.roleId == null ? '' : String(role.roleId);
          if (!roleId) {
            return;
          }
          roleIdSet.add(roleId);
          if (role?.roleName || role?.roleKey) {
            roleNameSet.add(role?.roleName || role?.roleKey || '');
          }
        });

        return {
          id: String(member.memberId || member.userId || ''),
          name: member.userName || member.userId || '-',
          userId: member.userId || '',
          email: member.email || '-',
          roles: Array.from(roleIdSet),
          roleNames: Array.from(roleNameSet),
          updatedAt: member.updatedAt,
        };
      });

      setMembers(normalizedMemberList);
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  useEffect(() => {
    getMembers();
  }, [getMembers]);

  useActivate(() => {
    getMembers();
  });

  // 通过 state 驱动编辑弹窗打开，避免在渲染阶段透传 ref。
  useEffect(() => {
    if (!editingMember) {
      return;
    }
    memberModalRef.current?.onOpen(editingMember);
  }, [editingMember]);

  const handleEditMember = useCallback((record: MemberItem) => {
    // 使用新对象引用，保证重复点击同一行也能触发 effect。
    setEditingMember({ ...record });
  }, []);

  const handleDelete = useCallback(
    (record: MemberItem) => {
      if (projectId) {
        deleteProjectMember(projectId, `${record.id}`).then(() => {
          message.success(formatMessage('common.deleteSuccessfully'));
          getMembers();
        });
      }
    },
    [formatMessage, getMembers, message, projectId]
  );

  const columns: any = useMemo(() => {
    return createMembersColumns({
      formatMessage,
      styles,
      onEdit: handleEditMember,
      onDelete: handleDelete,
      modal,
    });
  }, [formatMessage, handleDelete, handleEditMember, modal]);

  return (
    <>
      <TabLayout
        style={{ padding: '0 24px' }}
        toolbarActions={
          <Button
            type="primary"
            onClick={() => {
              memberModalRef.current?.onOpen();
            }}
          >
            {formatMessage('project.addMember')}
          </Button>
        }
      >
        <ProTable
          loading={loading}
          scroll={{ x: 'max-content', y: 'calc(100vh - 300px)' }}
          columns={columns}
          dataSource={members as any}
          resizeable
          pagination={false}
        />
      </TabLayout>
      <MemberModal
        ref={memberModalRef}
        refreshRequest={getMembers}
        existingMemberUserIds={members.map((m) => m.userId).filter(Boolean)}
      />
    </>
  );
};

export default Members;
