export const TypeEnum = {
  NodeRed: 'D-NodeRed',
  Mqtt: 'D-Mqtt',
  DataBase: 'D-DataBase',
  Apps: 'D-Apps',
} as const;

export interface NodeDataType {
  [key: string]: any;
}
