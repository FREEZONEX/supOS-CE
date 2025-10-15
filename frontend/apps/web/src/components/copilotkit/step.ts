export const getUnsRelationalStep = (data: { [key: string]: any }, unitType: string, currentStep?: string) => {
  const steps = [];

  steps.push({
    value: null,
    title: '打开文件夹弹框',
    titleEnglish: 'Add a new folder',
    stepName: 'openFolderNewModal',
    nextStep: 'namespace',
  });

  steps.push({
    value: { name: data.Namespace },
    title: '用于文件名新增',
    titleEnglish: 'Enter new name',
    stepName: 'namespace',
    nextStep: 'modelDescription',
  });

  steps.push({
    value: { modelDescription: data.modelDescription },
    title: '用于描述模型',
    titleEnglish: 'Fill in the description of model',
    stepName: 'modelDescription',
    nextStep: 'folderFields',
  });

  steps.push({
    value: { fields: data.fields },
    title: '模型字段',
    titleEnglish: 'Choose the fields of the moder',
    stepName: 'folderFields',
    nextStep: 'saveFolder',
  });

  steps.push({
    value: { page: data.saveAvailable },
    title: '保存操作',
    titleEnglish: 'Save all',
    stepName: 'saveFolder',
    nextStep: 'openFileNewModal',
  });

  steps.push({
    value: null,
    title: '打开文件弹框',
    titleEnglish: 'Add a new file',
    stepName: 'openFileNewModal',
    nextStep: 'name',
  });

  steps.push({
    value: { name: data.FileName },
    title: '用于文件新增',
    titleEnglish: 'Enter new name',
    stepName: 'name',
    nextStep: 'instanceDescription',
  });

  steps.push({
    value: { instanceDescription: data.fileDescription },
    title: '用于描述模型',
    titleEnglish: 'Fill in the description of model',
    stepName: 'instanceDescription',
    nextStep: 'dataType',
  });

  steps.push({
    value: { dataType: data.dataType },
    title: '数据库类型',
    titleEnglish: 'Choose the type of database',
    stepName: 'dataType',
    nextStep: 'fileFields',
  });

  steps.push({
    value: { fields: data.fields },
    title: '文件字段',
    titleEnglish: 'Choose the fields of the filed',
    stepName: 'fileFields',
    nextStep: 'next',
  });

  steps.push({
    value: { button: true },
    title: '数据库类型填写完成，进行下一步',
    titleEnglish: 'Database setting completed, choose next step',
    stepName: 'next',
    nextStep: 'saveFile',
  });

  steps.push({
    value: { page: data.saveAvailable },
    title: '保存操作',
    titleEnglish: 'Save all',
    stepName: 'saveFile',
    nextStep: null,
  });

  return { steps, unitType, currentStep: currentStep ?? steps?.[0]?.stepName, rawData: data };
};

export const getTimeSeriesStep = (data: { [key: string]: any }, unitType: string, currentStep?: string) => {
  const steps = [];

  steps.push({
    value: null,
    title: '打开文件夹弹框',
    titleEnglish: 'Add a new folder',
    stepName: 'openFolderNewModal',
    nextStep: 'namespace',
  });

  steps.push({
    value: { name: data.Namespace },
    title: '用于文件名新增',
    titleEnglish: 'Enter new name',
    stepName: 'namespace',
    nextStep: 'modelDescription',
  });

  steps.push({
    value: { modelDescription: data.modelDescription },
    title: '用于描述模型',
    titleEnglish: 'Fill in the description of model',
    stepName: 'modelDescription',
    nextStep: 'folderFields',
  });

  steps.push({
    value: { fields: data.fields },
    title: '模型字段',
    titleEnglish: 'Choose the fields of the moder',
    stepName: 'folderFields',
    nextStep: 'saveFolder',
  });

  steps.push({
    value: { page: data.saveAvailable },
    title: '保存操作',
    titleEnglish: 'Save all',
    stepName: 'saveFolder',
    nextStep: 'openFileNewModal',
  });

  steps.push({
    value: null,
    title: '打开文件弹框',
    titleEnglish: 'Add a new file',
    stepName: 'openFileNewModal',
    nextStep: 'name',
  });

  steps.push({
    value: { name: data.FileName },
    title: '用于文件新增',
    titleEnglish: 'Enter new name',
    stepName: 'name',
    nextStep: 'instanceDescription',
  });

  steps.push({
    value: { instanceDescription: data.fileDescription },
    title: '用于描述模型',
    titleEnglish: 'Fill in the description of model',
    stepName: 'instanceDescription',
    nextStep: 'dataType',
  });

  steps.push({
    value: { dataType: data.dataType },
    title: '数据库类型',
    titleEnglish: 'Choose the type of database',
    stepName: 'dataType',
    nextStep: 'fileFields',
  });

  steps.push({
    value: { fields: data.fields },
    title: '文件字段',
    titleEnglish: 'Choose the fields of the filed',
    stepName: 'fileFields',
    nextStep: 'next',
  });

  steps.push({
    value: { button: true },
    title: '数据库类型填写完成，进行下一步',
    titleEnglish: 'Database setting completed, choose next step',
    stepName: 'next',
    nextStep: 'protocol',
  });

  steps.push({
    value: {
      protocol: {
        protocol: data?.protocol?.protocol,
      },
    },
    title: '选择modbus',
    titleEnglish: 'Choose the Modbus protocol',
    stepName: 'protocol',
    nextStep: 'server',
  });

  steps.push({
    value: {
      protocol: {
        serverName: data?.protocol?.serverName || 'supos_custom',
        server: {
          port: data?.protocol?.port,
          host: data?.protocol?.host,
        },
      },
    },
    title: '设置server',
    titleEnglish: 'set the modbus server',
    stepName: 'server',
    nextStep: 'unitId',
  });

  steps.push({
    value: {
      protocol: {
        unitId: data?.protocol?.unitId,
      },
    },
    title: '设置unitId',
    titleEnglish: 'set the modbus unitId',
    stepName: 'unitId',
    nextStep: 'fc',
  });

  steps.push({
    value: {
      protocol: {
        fc: data?.protocol?.fc,
      },
    },
    title: '设置fc',
    titleEnglish: 'set the modbus fc',
    stepName: 'fc',
    nextStep: 'address',
  });

  steps.push({
    value: {
      protocol: {
        address: data?.protocol?.address,
      },
    },
    title: '设置address',
    titleEnglish: 'set the modbus data start address',
    stepName: 'address',
    nextStep: 'quantity',
  });

  steps.push({
    value: {
      protocol: {
        quantity: data?.protocol?.quantity,
      },
    },
    title: '设置quantity',
    titleEnglish: 'set the modbus data quantity',
    stepName: 'quantity',
    nextStep: 'pollRate',
  });

  steps.push({
    value: {
      protocol: {
        pollRate: data?.protocol?.pollRate,
      },
    },
    title: '设置pollRate',
    titleEnglish: 'set the modbus data pollRate',
    stepName: 'pollRate',
    nextStep: 'saveFile',
  });

  steps.push({
    value: { page: data.saveAvailable },
    title: '保存操作',
    titleEnglish: 'Save all',
    stepName: 'saveFile',
    nextStep: null,
  });

  return { steps, unitType, currentStep: currentStep ?? steps?.[0]?.stepName, rawData: data };
};
