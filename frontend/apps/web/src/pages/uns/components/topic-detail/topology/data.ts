import { TypeEnum } from './types.ts';
import { Markup } from '@antv/x6';

export const markupLine = {
  markup: Markup.getForeignObjectMarkup(),
  attrs: {
    fo: {
      width: 16,
      height: 16,
      x: -10,
      y: -8,
    },
  },
};

const commonLine = {
  sourceMarker: { name: 'circle', r: 2 },
  targetMarker: { name: 'circle', r: 2 },
  stroke: 'var(--supos-theme-color)',
  strokeDasharray: 3,
  style: {
    animation: 'ant-line 60s infinite linear',
  },
};

export const data = {
  nodes: [
    {
      id: TypeEnum.NodeRed,
      shape: TypeEnum.NodeRed,
      x: 250,
      y: 0,
      data: {
        topic: '',
        active: false,
        id: '',
        flowId: '',
        flowStatus: '',
        flowName: '',
        bindId: '',
        onBindingChange: (type: string, item: any) => {
          console.log(type, item);
        },
      },
    },
    {
      id: TypeEnum.Mqtt,
      shape: TypeEnum.Mqtt,
      x: 510,
      y: 0,
      data: {
        active: false,
      },
    },
    {
      id: TypeEnum.DataBase,
      shape: TypeEnum.DataBase,
      x: 700,
      y: 0,
      data: {
        dataType: 1,
        active: false,
        alias: '',
      },
    },
    {
      id: TypeEnum.Apps,
      shape: TypeEnum.Apps,
      x: 930,
      y: 0,
      data: {
        active: false,
        bindId: '',
        subtitle: 'Grafana',
        onBindingChange: (type: string, item: any) => {
          console.log(type, item);
        },
      },
    },
  ],
  edges: [
    {
      shape: 'edge',
      source: TypeEnum.NodeRed,
      target: TypeEnum.Mqtt,
      id: 'pushMqtt',
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: { ...commonLine },
      },
      label: {
        position: 0,
      },
    },
    {
      shape: 'edge',
      source: TypeEnum.Mqtt,
      target: TypeEnum.DataBase,
      id: 'pullMqttOrDataPersistence', //pullMqtt从mqtt拉数据或dataPersistence数据持久化
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: { ...commonLine },
      },
      label: {
        position: 0,
      },
    },
    {
      shape: 'edge',
      source: TypeEnum.DataBase,
      target: TypeEnum.Apps,
      id: TypeEnum.Apps + '1',
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: { ...commonLine },
      },
    },
  ],
};

export function findDate({ dataType, withSave2db, withDashboard, dashboardType }: any) {
  const _data = data;
  if (dataType == 3) {
    return Object.assign({}, _data, {
      nodes: _data.nodes.slice(1),
      edges: _data.edges.slice(1),
    });
  }
  if (!withSave2db) {
    if (!withDashboard) {
      return Object.assign({}, _data, {
        nodes: _data.nodes.slice(0, -2),
        edges: _data.edges.slice(0, -2),
      });
    } else {
      return Object.assign({}, _data, {
        nodes: [_data.nodes[0], _data.nodes[1], { ..._data.nodes[3], x: 700 }],
        edges: [_data.edges[0], { ..._data.edges[1], target: TypeEnum.Apps }],
      });
    }
  }
  if (withSave2db && (!dashboardType?.includes('grafana') || !withDashboard)) {
    return Object.assign({}, _data, {
      nodes: _data.nodes.slice(0, -1),
      edges: _data.edges.slice(0, -1),
    });
  }
  return _data;
}
