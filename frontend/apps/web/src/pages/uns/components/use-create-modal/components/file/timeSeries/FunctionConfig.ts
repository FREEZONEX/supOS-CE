interface FunctionConfig {
  label: string;
  value: string;
  minKeys: number;
  maxKeys: number;
}

export const FUNCTION_TYPES: FunctionConfig[] = [
  // 单参数函数
  // { label: 'ABS', value: 'ABS', minKeys: 1, maxKeys: 1 },
  // { label: 'FLOOR', value: 'FLOOR', minKeys: 1, maxKeys: 1 },
  // { label: 'CEIL', value: 'CEIL', minKeys: 1, maxKeys: 1 },
  // // 双参数函数
  // { label: 'MOD', value: 'MOD', minKeys: 2, maxKeys: 2 },
  // 类型转换函数
  // { label: 'CAST', value: 'CAST', minKeys: 1, maxKeys: 1 },
  // 聚合函数
  { label: 'AVG', value: 'AVG', minKeys: 1, maxKeys: 1 },
  { label: 'COUNT', value: 'COUNT', minKeys: 1, maxKeys: 1 }, // COUNT(*) 特殊处理
  { label: 'SUM', value: 'SUM', minKeys: 1, maxKeys: 1 },
  { label: 'MAX', value: 'MAX', minKeys: 1, maxKeys: 1 },
  { label: 'MIN', value: 'MIN', minKeys: 1, maxKeys: 1 },
];
