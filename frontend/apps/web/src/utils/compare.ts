// 比较两个 location 对象，排除指定的属性
import { isEqual, omit } from 'lodash-es';

export const compareLocations = (loc1: any, loc2: any, excludeProps: string[] = []) => {
  // 使用 lodash 的 omit 方法过滤掉指定的属性
  const filteredLoc1 = omit(loc1, excludeProps);
  const filteredLoc2 = omit(loc2, excludeProps);

  // 使用 JSON.stringify 比较两个过滤后的对象
  return isEqual(filteredLoc1, filteredLoc2);
};
