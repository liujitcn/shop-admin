import service from "@/utils/request";
import type { Login } from "@liujitcn/shop-base";
import type { Empty } from "@/rpc/google/protobuf/empty";

type CaptchaResponse = Login.CaptchaResponse;
type RefreshTokenRequest = Login.RefreshTokenRequest;
type RefreshTokenResponse = Login.RefreshTokenResponse;
type LoginService = Login.LoginService;

const LOGIN_URL = "/login";

/** 登录公共服务 */
export class LoginServiceImpl implements LoginService {
  /** 验证码 */
  Captcha(request: Empty): Promise<CaptchaResponse> {
    return service<Empty, CaptchaResponse>({
      url: `${LOGIN_URL}/captcha`,
      method: "get",
      params: request,
    });
  }
  /** 登出 */
  Logout(request: Empty): Promise<Empty> {
    return service<Empty, Empty>({
      url: `${LOGIN_URL}/logout`,
      method: "delete",
      data: request,
    });
  }
  /** 刷新认证令牌 */
  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    return service<RefreshTokenRequest, RefreshTokenResponse>({
      url: `${LOGIN_URL}/refreshToken`,
      method: "post",
      data: request,
    });
  }
}

export const defLoginService = new LoginServiceImpl();
