/**
 * 获取当前用户
 */
import { getUser } from '@/apis/chat2db';

export const queryChat2dbCurUser = async () => {
  // null 表示在padding，返回 void 0(undefined)表示未登录
  const curUser = (await getUser()) || void 0;
  const { data } = curUser;
  if (data) {
    // 向cookie中写入当前用户id
    const date = new Date('2030-12-30 12:30:00').toUTCString();
    document.cookie = `CHAT2DB.USER_ID=${data?.id};Expires=${date}`;
    return data;
  } else {
    return false;
  }
};
