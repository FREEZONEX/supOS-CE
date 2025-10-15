export const TransformFieldNameOutData = (list: any[], fieldName: any) => {
  return list.map((item) => {
    const transformedItem: any = { ...item };
    for (const key in fieldName) {
      if (fieldName[key]) {
        transformedItem[key] = item[fieldName[key]];
        delete transformedItem[fieldName[key]];
      }
    }
    return transformedItem;
  });
};
