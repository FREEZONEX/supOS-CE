export const getDefaultFields = (qualityName: string, timestampName: string) => [
  { name: qualityName, type: 'LONG', systemField: true },
  { name: timestampName, type: 'DATETIME', systemField: true },
];
