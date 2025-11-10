import fetch from 'node-fetch';
import { logger } from '../utils';

export interface CityInfo {
  province: string;
  city: string;
  county?: string;
}

export interface WeatherData {
  degree: string;
  humidity: string;
  precipitation: string;
  pressure: string;
  update_time: string;
  weather: string;
  weather_code: string;
  weather_short: string;
  wind_direction: string;
  wind_power: string;
  wind_direction_name: string;
}

export interface WeatherInfo {
  weatherData: WeatherData;
  cityInfo: CityInfo;
}

interface WeatherResponse {
  data: {
    observe: WeatherData;
  };
  message: string;
  status: number;
}

interface CitySearchResponse {
  data: {
    // "101210101": "浙江, 杭州"
    [key: string]: string;
  };
  message: string;
  status: number;
}

export class WeatherTool {
  private readonly baseUrl = 'https://wis.qq.com';

  /**
   * 搜索城市信息
   * @param city 城市名称
   * @returns 城市信息 浙江,杭州
   */
  async searchCity(city: string): Promise<CityInfo> {
    try {
      logger.log(`[WeatherTool] 搜索城市: ${city}`);

      const response = await fetch(`${this.baseUrl}/city/like?source=pc&city=${encodeURIComponent(city)}`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = (await response.json()) as CitySearchResponse;

      if (data.status !== 200) {
        throw new Error(`API error: ${data.message}`);
      }

      const addressString = Object.values(data.data)[0];

      if (!addressString) {
        throw new Error(`API error: 没有城市信息`);
      }
      const [province, _city, county] = addressString.split(',').map((m) => m.trim());
      logger.log(`[WeatherTool] 搜索城市成功`, addressString);
      return {
        province,
        city: _city,
        county,
      };
    } catch (error) {
      logger.error(`[WeatherTool] 搜索城市失败:`, error);
      throw error;
    }
  }

  /**
   * 获取天气信息
   * @param province 省份
   * @param city 城市
   * @param county 县/区
   * @returns 天气数据
   */
  async getWeather(province: string, city: string, county?: string): Promise<WeatherData> {
    try {
      const params = new URLSearchParams({
        weather_type: 'observe',
        source: 'pc',
      });
      if (province !== undefined) params.append('province', province);
      if (city !== undefined) params.append('city', city);
      if (county !== undefined) params.append('county', county);
      const response = await fetch(`${this.baseUrl}/weather/common?${params.toString()}`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = (await response.json()) as WeatherResponse;

      if (data.status !== 200) {
        throw new Error(`API error: ${data.message}`);
      }

      logger.log(`[WeatherTool] 天气数据获取成功`, data.data);
      return data.data.observe;
    } catch (error) {
      logger.error(`[WeatherTool] 获取天气信息失败:`, error);
      throw error;
    }
  }

  /**
   * 根据城市名称获取天气信息
   * @param city 城市名称
   * @returns 天气信息
   */
  async getWeatherByCity(city: string): Promise<WeatherInfo> {
    try {
      logger.log(`[WeatherTool] 开始获取城市天气: ${city}`);

      // 1. 搜索城市
      const { city: _city, county, province } = await this.searchCity(city);
      const weatherData = await this.getWeather(province, _city, county);

      // 2. 获取天气信息
      return {
        weatherData,
        cityInfo: {
          city: _city,
          county,
          province,
        },
      };
    } catch (error) {
      logger.error(`[WeatherTool] 获取城市天气失败:`, error);
      throw error;
    }
  }
}

/**
 * 格式化天气结果
 */
export function formatWeatherResult(result: WeatherInfo) {
  const { cityInfo, weatherData: weather } = result;

  return `🌤️ 天气信息：

📍 位置：${cityInfo.province}省 - ${cityInfo.city}${cityInfo.county ? `- ${cityInfo.county}` : ''}
🌡️ 温度：${weather.degree}°C
🌤️ 天气：${weather.weather} (${weather.weather_short})
💧 湿度：${weather.humidity}%
💨 风向：${weather.wind_direction_name} ${weather.wind_power}级
🌊 气压：${weather.pressure}hPa
💦 降水量：${weather.precipitation}mm
🕐 更新时间：${weather.update_time}

天气代码：${weather.weather_code}`;
}

/**
 * 格式化错误信息
 */
export function formatErrorMessage(error: any): string {
  logger.error(error);
  if (error instanceof Error) {
    if (error.message.includes('未找到城市')) {
      return `❌ 天气查询失败：未找到指定的城市，请检查城市名称是否正确。`;
    }
    if (error.message.includes('参数验证失败')) {
      return `❌ 参数错误：${error.message.replace('参数验证失败: ', '')}`;
    }
    if (error.message.includes('HTTP error') || error.message.includes('API error')) {
      return `❌ 网络请求失败：天气服务暂时不可用，请稍后重试。`;
    }
    return `❌ 天气查询失败：${error.message}`;
  }

  return `❌ 天气查询失败：未知错误`;
}
