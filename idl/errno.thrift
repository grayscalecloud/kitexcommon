
namespace go errno

enum Err {
    Success            = 200,
    NoRoute            = 405,
    NoMethod           = 404,
    BadRequest         = 500,
    ParamsErr          = 504,
    AuthorizeFail      = 403,
    TooManyRequest     = 429,
    ServiceErr         = 502,
    RecordNotFound     = 1000,
    RecordAlreadyExist = 1010,
    DirtyData          = 1020,

    RPCUserSrvErr      = 30001,
    RPCUserAdminSrvErr = 30002,
    RPCOrderSrvErr     = 30003,
    RPCProductSrvErr   = 30004,
    RPCActivitySrvErr  = 30005,
    RPCCartSrvErr      = 30006,

    UserSrvErr         = 40001,
    UserAdminSrvErr    = 40002,
    OrderSrvErr        = 40003,
    ProductSrvErr      = 40004,
    ActivitySrvErr     = 40005,
    CartSrvErr         = 40006,
} (go.type = "int64")

struct ErrNo {
  1: i64 ErrCode,
  2: string ErrMsg
}

struct Response {
  1: i64 Code,
  2: string Message,
  3: binary Data
}

const ErrNo Success = {"ErrCode": Err.Success, "ErrMsg": "success"}
const ErrNo NoRoute = {"ErrCode": Err.NoRoute, "ErrMsg": "no route"}
const ErrNo NoMethod = {"ErrCode": Err.NoMethod, "ErrMsg": "no method"}
const ErrNo BadRequest = {"ErrCode": Err.BadRequest, "ErrMsg": "bad request"}
const ErrNo ParamsErr = {"ErrCode": Err.ParamsErr, "ErrMsg": "params error"}
const ErrNo AuthorizeFail = {"ErrCode": Err.AuthorizeFail, "ErrMsg": "authorize failed"}
const ErrNo TooManyRequest = {"ErrCode": Err.TooManyRequest, "ErrMsg": "too many requests"}
const ErrNo ServiceErr = {"ErrCode": Err.ServiceErr, "ErrMsg": "service error"}
const ErrNo RPCUserSrvErr = {"ErrCode": Err.RPCUserSrvErr, "ErrMsg": "rpc user service error"}
const ErrNo UserSrvErr = {"ErrCode": Err.UserSrvErr, "ErrMsg": "user service error"}
const ErrNo RPCUserAdminSrvErr = {"ErrCode": Err.RPCUserAdminSrvErr, "ErrMsg": "rpc system service error"}
const ErrNo UserAdminSrvErr = {"ErrCode": Err.UserAdminSrvErr, "ErrMsg": "system service error"}
const ErrNo RPCProductSrvErr = {"ErrCode": Err.RPCProductSrvErr, "ErrMsg": "rpc car service error"}
const ErrNo ProductSrvErr = {"ErrCode": Err.ProductSrvErr, "ErrMsg": "car service error"}
const ErrNo RPCOrderSrvErr = {"ErrCode": Err.RPCOrderSrvErr, "ErrMsg": "rpc profile service error"}
const ErrNo OrderSrvErr = {"ErrCode": Err.OrderSrvErr, "ErrMsg": "profile service error"}
const ErrNo RPCActivitySrvErr = {"ErrCode": Err.RPCActivitySrvErr, "ErrMsg": "rpc trip service error"}
const ErrNo ActivitySrvErr = {"ErrCode": Err.ActivitySrvErr, "ErrMsg": "trip service error"}
const ErrNo RPCCartSrvErr = {"ErrCode": Err.RPCCartSrvErr, "ErrMsg": "rpc trip service error"}
const ErrNo CartSrvErr = {"ErrCode": Err.CartSrvErr, "ErrMsg": "trip service error"}
const ErrNo RecordNotFound = {"ErrCode": Err.RecordNotFound, "ErrMsg": "record not found"}
const ErrNo RecordAlreadyExist = {"ErrCode": Err.RecordAlreadyExist, "ErrMsg": "record already exist"}
const ErrNo DirtyData = {"ErrCode": Err.DirtyData, "ErrMsg": "dirty data"}