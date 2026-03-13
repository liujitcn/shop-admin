import service from "@/utils/request";
import type { Config } from "@liujitcn/shop-base";

type ConfigService = Config.ConfigService;
type ConfigRequest = Config.ConfigRequest;
type ConfigResponse = Config.ConfigResponse;

const CONFIG_URL = "/config";

/** 系统配置公共服务 */
export class ConfigServiceImpl implements ConfigService {
  /** 获取系统配置 */
  GetConfig(request: ConfigRequest): Promise<ConfigResponse> {
    return service<ConfigRequest, ConfigResponse>({
      url: `${CONFIG_URL}`,
      method: "get",
      params: request,
    });
  }
}

export const defConfigService = new ConfigServiceImpl();
