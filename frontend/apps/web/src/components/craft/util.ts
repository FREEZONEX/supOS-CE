function generateSampleData(fields: any, excludeFields: any = []) {
  const sampleData: any = {};
  if (!fields) return '';
  fields.forEach((field: any) => {
    if (!excludeFields.includes(field.name)) {
      switch (field.type) {
        case 'String':
          sampleData[field.name] = 'example ' + field.name; // 默认字符串值
          break;
        case 'Int':
          sampleData[field.name] = 1; // 默认整数值
          break;
        case 'Float':
          sampleData[field.name] = 1.23; // 默认浮动值
          break;
        case 'bigint':
          sampleData[field.name] = 1; // 默认浮动值
          break;
        case 'float8':
          sampleData[field.name] = 1; // 默认浮动值
          break;
        case 'Boolean':
          sampleData[field.name] = true; // 默认布尔值
          break;
        case 'date':
          sampleData[field.name] = '2024-11-07'; // 默认布尔值
          break;
        case 'timestamptz':
          sampleData[field.name] = '2024-11-08T01:34:58.665Z'; // 默认布尔值
          break;
        case 'timetz':
          sampleData[field.name] = '01:34:58Z'; // 默认布尔值
          break;
        default:
          sampleData[field.name] = null; // 对于未知类型，使用 null
      }
    }
  });

  // 转换为格式化的字符串
  const resultString = Object.entries(sampleData)
    .map(([key, value]) => `${key}: ${JSON.stringify(value)}`)
    .join(', ');

  return resultString;
}

export const getGeneratedPrompt = ({ codeCommand, apiUrl, fieldsName, fileds }: any) => {
  const dataObj = fileds ? generateSampleData(fileds, ['id']) : null;
  const commonPrompt = String.raw`#prompt
请生成一个html包含下面所有的客户需求，注意一定要单一html文件完成需求，所有样式文件只能以cdn的方式引入，要完整完成需求。
注意下列事项
1.尽量参考最近比较火热的css样式，完善页面的美观度
2.直接返回完整的html，不需要返回任何其他的内容。任何多余的文字，包括注释都不允许存在
3.特别注意：如果${fieldsName}是存在的，请严格使用下面提供的固定接口，以及我给的案例
`;

  const apiPrompt = `注意：接口固定地址为${apiUrl}/hasura/home/v1/graphql,`;

  const exampleListPrompt =
    apiPrompt +
    String.raw`下面是我的查询列表用例，参数案例： query MyQuery {
  ${codeCommand.databaseName} {
    ${fieldsName}
  }
}
  `;

  const exampleAddPrompt = String.raw`下面是新增到我的数据库中的案例，参数案例：mutation {
  insert_${codeCommand.databaseName}(objects: { ${dataObj} }) {
    returning {
      ${fieldsName}
    }
  }
}
  `;

  const exampleDelPrompt = String.raw`下面是删除数据库中id为0的案例，参数案例：mutation {
  delete_${codeCommand.databaseName}(where: {_id: {_eq: 0}}) {
    returning {
      ${fieldsName}
    }
  }
}
  `;

  const updatePrompt = String.raw`下面是更新数据库的案例，参数案例：mutation {
  update_${codeCommand.databaseName}(where: {_id: {_eq: "5"}}, _set: {item: "aa"}) {
    returning {
      ${fieldsName}
    }
  }
}
  `;

  const errorPrompt = String.raw`如果接口出现下面类似的报错，请弹框提示错误：{
    "errors": [
        {
            "message": "field 'id' not found in type: '_a_a_bool_exp'",
            "extensions": {
                "path": "$.selectionSet.delete__a_a.args.where.id",
                "code": "validation-failed"
            }
        }
    ]
}
  `;

  const codeCommandPrompt =
    String.raw`#写明需求
根据上面的内容，` + codeCommand.prompt;

  if (!fieldsName) return commonPrompt + codeCommandPrompt;
  return (
    commonPrompt +
    apiPrompt +
    exampleListPrompt +
    exampleAddPrompt +
    exampleDelPrompt +
    updatePrompt +
    errorPrompt +
    codeCommand.prompt
  );
};
