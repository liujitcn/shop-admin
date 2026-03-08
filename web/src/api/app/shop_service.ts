import service from "@/utils/request";
import {
  type ShopServiceService,
  type ListShopServiceResponse,
} from "@/rpc/app/shop_service";
import type { Empty } from "@/rpc/google/protobuf/empty";

const SHOP_SERVICE_URL = "/app/shop/service";

/** 商城服务 */
export class ShopServiceServiceImpl implements ShopServiceService {
  /** 查询商城服务列表 */
  ListShopService(request: Empty): Promise<ListShopServiceResponse> {
    return service<Empty, ListShopServiceResponse>({
      url: `${SHOP_SERVICE_URL}`,
      method: "get",
    });
  }
}

export const defShopServiceService = new ShopServiceServiceImpl(); 