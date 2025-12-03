export const getExampleForJavaType = (type: string, name: string) => {
  switch (type) {
    case 'STRING':
      return name;
    case 'INTEGER':
      return 42; // 使用更具代表性的数字
    case 'LONG':
      return 123456; // 使用更具代表性的数字
    case 'FLOAT':
      return 123.45; // 使用更具代表性的数字
    case 'DOUBLE':
      return '123.456789';
    case 'BOOLEAN':
      return true;
    case 'DATETIME':
      return new Date().getTime();
    default:
      return 23;
  }
};
