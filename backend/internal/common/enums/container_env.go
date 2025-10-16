package enums

// ContainerEnvEnum represents service container built-in environment variables
type ContainerEnvEnum struct {
	Name        string
	Description string
}

var (
	ContainerEnvServiceIsShow = ContainerEnvEnum{
		Name:        "service_is_show",
		Description: "服务是否显示",
	}

	ContainerEnvServiceLogo = ContainerEnvEnum{
		Name:        "service_logo",
		Description: "LOGO",
	}

	ContainerEnvServiceDescription = ContainerEnvEnum{
		Name:        "service_description",
		Description: "服务描述",
	}

	ContainerEnvServiceRedirectURL = ContainerEnvEnum{
		Name:        "service_redirect_url",
		Description: "高阶使用跳转路由",
	}

	ContainerEnvServiceAccount = ContainerEnvEnum{
		Name:        "service_account",
		Description: "帐号",
	}

	ContainerEnvServicePassword = ContainerEnvEnum{
		Name:        "service_password",
		Description: "密码",
	}
)
