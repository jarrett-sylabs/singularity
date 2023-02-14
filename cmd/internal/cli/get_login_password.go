package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/sylabs/singularity/docs"
	"github.com/sylabs/singularity/pkg/cmdline"
)

var (
	// GetLoginPasswordLibraryURI holds the base URI to a Sylabs library API instance
	GetLoginPasswordLibraryURI string
)

// --library
var getLoginPasswordLibraryFlag = cmdline.Flag{
	ID:           "getLoginPasswordLibraryFlag",
	Value:        &GetLoginPasswordLibraryURI,
	DefaultValue: "",
	Name:         "library",
	Usage:        "URI for library to search",
	EnvKeys:      []string{"LIBRARY"},
}

func init() {
	addCmdInit(func(cmdManager *cmdline.CommandManager) {
		cmdManager.RegisterCmd(GetLoginPasswordCmd)

		cmdManager.RegisterFlagForCmd(&getLoginPasswordLibraryFlag, GetLoginPasswordCmd)
	})
}

var GetLoginPasswordCmd = &cobra.Command{
	DisableFlagsInUseLine: true,

	Run: callShimAPIEndpoint,

	Use:     docs.GetLoginPasswordUse,
	Short:   docs.GetLoginPasswordShort,
	Long:    docs.GetLoginPasswordLong,
	Example: docs.GetLoginPasswordExample,
}

// need to get the token

func callShimAPIEndpoint(cmd *cobra.Command, args []string) {
	endPoint := "https://library.se.k3s/v1/rbac/users/current"
	defaultConfigURI := ""

	config, err := getLibraryClientConfig(defaultConfigURI)
	if err != nil {
		fmt.Errorf("config err: ", err)
	}

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		fmt.Errorf("request err: ", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", config.AuthToken))
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		fmt.Errorf("client err: %v", err)
	}
	var u User
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		fmt.Errorf("jsonerr: %v", err)
	}
	if u.OidcUserMeta.Secret != "" {
		fmt.Println(u.OidcUserMeta.Secret)
	} else {
		fmt.Errorf("failed to get secret: %v", err)
	}
}

type OidcUserMeta struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Subiss       string `json:"subiss"`
	Secret       string `json:"secret"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
}

// User
type User struct {
	Email           string       `json:"email"`
	RealName        string       `json:"realname"`
	Comment         string       `json:"comment"`
	UserId          string       `json:"user_id"`
	UserName        string       `json:"username"`
	SysAdminFlag    bool         `json:"sysadmin_flag"`
	AdminRoleInAuth string       `json:"admin_role_in_auth"`
	OidcUserMeta    OidcUserMeta `json:"oidc_user_meta"`
	CreationTime    string       `json:"creation_time"`
	UpdateTime      string       `json:"update_time"`
}
