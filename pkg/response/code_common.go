package response

type Code int

const (
	CommonCode Code = 40000 // 公共模块状态码开头

	UploadCode = 12000 // 上传相关状态码

	AuthCode               = 43100 // 用户认证相关状态码
	ProductCode            = 43500 // 单品相关状态码
	WardrobesCode          = 43600 // 衣橱相关状态码
	BrandCode              = 43700 // 品牌相关状态码
	UserCode               = 43200 // 用户相关状态码
	OutfitCode             = 43800 // 搭配相关状态码
	UserCustomCode         = 43900 // 用户自定义相关状态码
	AiOutfitCode           = 44000 // AI搭配相关状态码
	RiskCode               = 44100 // 风险检测相关状态码
	UserOutfitCalendarCode = 44200 // 用户搭配日历相关状态码
	CityCode               = 44300 // 城市相关状态码
	ShareCode              = 44400 // 分享相关状态码
	UserMuseOutfitCode     = 44500 // 用户穿搭灵感状态码
	StyleCode              = 44600 // 风格状态码
	KnowledgeCode          = 44700 // 知识库状态码
	UserStyleCode          = 44800 // 用户风格状态码
	PayCode                = 44900 // 支付相关状态码
	OrderCode              = 45000 // 订单相关状态码

)

// 通用校验错误状态码  Validate 开头 400XX
const (
	ValidateCommonError       = iota + CommonCode // 40000 表单校验失败
	ValidateVerification                          // 40001 手机验证码校验失败
	ValidateCaptcha                               // 40002 图片验证码校验失败
	UploadSignError                               // 40003 获取上传 token 失败
	UploadCallbackVerifyError                     // 40004 上传回调校验失败
	SMSError                                      // 40005 短信发送失败
	SMSLimitError                                 // 40006 短信验证码超出频率限制 (每分钟一次, 一天最多发送十次)
	ParamsCommonError                             // 40007 API请求参数错误
	ServiceCommonError                            // 40008 系统错误
	FuncCommonError                               // 40009 操作失败
	DBSelectCommonError                           // 40010 DB查询错误
	DBInsertCommonError                           // 40011 DB添加错误
	DBUpdateCommonError                           // 40012 DB更新错误
	DBDeleteCommonError                           // 40013 DB删除错误
	AuthCommonError                               // 40014 权限错误
)
