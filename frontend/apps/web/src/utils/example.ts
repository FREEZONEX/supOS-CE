export const getExampleForJavaType = (type: string, name: string) => {
  switch (type) {
    case 'string':
      return name;
    case 'integer':
    case 'long':
    case 'float':
    case 'double':
      return 1;
    case 'boolean':
      return true;
    case 'datetime':
      return new Date().getTime();
    default:
      return 1;
  }
};
