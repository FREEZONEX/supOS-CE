export const CONTROL_WHITE_LIST = [
  'Text',
  'Number',
  'Date',
  'Textarea',
  'Radio',
  'Checkbox',
  'Select',
  'SelectMultiple',
  'Member',
  'Department',
  'SubForm',
  'Switch',
];

export const BASIC_OPERATORS = [
  { operatorName: '+', operatorValue: '+' },
  { operatorName: '-', operatorValue: '-' },
  { operatorName: '*', operatorValue: '*' },
  { operatorName: '/', operatorValue: '/' },
  { operatorName: '(', operatorValue: '(' },
  { operatorName: ')', operatorValue: ')' },
];

export const CONDITION_OPERATORS = [
  { operatorName: '>', operatorValue: '>' },
  { operatorName: '<', operatorValue: '<' },
  { operatorName: '>=', operatorValue: '>=' },
  { operatorName: '<=', operatorValue: '<=' },
  { operatorName: '==', operatorValue: '==' },
  { operatorName: '!=', operatorValue: '!=' },
];

export const FUNCTIONS_TYPE = [
  {
    cateName: 'mathFun',
    cateValue: 'math',
  },
];

export interface FunctionConfig {
  label: string;
  name: string;
  intro: string;
  usage: string;
  example: string;
}

export const FUNCTIONS_DATA: { [index: string]: FunctionConfig[] } = {
  math: [
    {
      label: 'math0label',
      name: 'ABS()',
      intro: 'math0intro',
      usage: 'math0usage',
      example: 'math0example',
    },
    {
      label: 'math1label',
      name: 'AVERAGE()',
      intro: 'math1intro',
      usage: 'math1usage',
      example: 'math1example',
    },
    {
      label: 'math2label',
      name: 'COUNT()',
      intro: 'math2intro',
      usage: 'math2usage',
      example: 'math2example',
    },
    {
      label: 'math3label',
      name: 'FIXED()',
      intro: 'math3intro',
      usage: 'math3usage',
      example: 'math3example',
    },
    {
      label: 'math4label',
      name: 'INT()',
      intro: 'math4intro',
      usage: 'math4usage',
      example: 'math4example',
    },
    {
      label: 'math5label',
      name: 'LOG()',
      intro: 'math5intro',
      usage: 'math5usage',
      example: 'math5example',
    },
    {
      label: 'math6label',
      name: 'MOD()',
      intro: 'math6intro',
      usage: 'math6usage',
      example: 'math6example',
    },
    {
      label: 'math7label',
      name: 'MAX()',
      intro: 'math7intro',
      usage: 'math7usage',
      example: 'math7example',
    },
    {
      label: 'math8label',
      name: 'MIN()',
      intro: 'math8intro',
      usage: 'math8usage',
      example: 'math8example',
    },
    {
      label: 'math9label',
      name: 'POWER()',
      intro: 'math9intro',
      usage: 'math9usage',
      example: 'math9example',
    },
    {
      label: 'math10label',
      name: 'PRODUCT()',
      intro: 'math10intro',
      usage: 'math10usage',
      example: 'math10example',
    },
    {
      label: 'math11label',
      name: 'SQRT()',
      intro: 'math11intro',
      usage: 'math11usage',
      example: 'math11example',
    },
    {
      label: 'math12label',
      name: 'SUM()',
      intro: 'math12intro',
      usage: 'math12usage',
      example: 'math12example',
    },
  ],
};

export const DATA_TYPE_MAP: { [key: string]: string } = {
  datetime: '	时间日期',
  counter: '流水号',
  string: '文本',
  long: '整数',
  decimal: '浮点',
  submodel: '子对象',
  linkmodel: '关联模型',
  linkfield: '引用字段',
  syscode: '系统编码',
  'set<file>': '文件集合',
  sysstaff: '员工',
  sysdept: '组织',
  syscompany: '公司',
  sysuser: '用户',
  systime: '时间',
  boolean: '布尔',
  double: '数值',
  option: '单选',
  'set<option>': '多选',
  'range<datetime>': '日期范围',
  'range<long>': '整数范围',
  'range<double>': '浮点范围',
};
