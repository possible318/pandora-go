package main

import (
	tokens "github.com/pandora_go/exts/token"
	"github.com/pandora_go/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "PDR-Go",
	Short: "PDR golang",
	Long:  "PDR is a free",
	Run:   RunCmd,
}

func RunCmd(cmd *cobra.Command, args []string) {
	s, _ := cmd.Flags().GetString("server")
	p, _ := cmd.Flags().GetString("proxy")
	f, _ := cmd.Flags().GetString("token_file")
	a, _ := cmd.Flags().GetBool("api")
	v, _ := cmd.Flags().GetBool("verbose")
	e, _ := cmd.Flags().GetBool("sentry")

	if e {
		viper.Set("sentry", true)
	}

	if f != "" {
		tokens.InitAccessToken(f)
	}

	if p != "" {
		viper.Set("proxy", p)
	}

	if v {
		viper.Set("verbose", true)
	}

	// api 方式
	if a {
		viper.Set("api", true)
	}

	// web 模式
	if s != "" {
		web.Run(s)
	}

	// 命令行模式
	//chat.Run(p, f, a, d, v)
}

func main() {
	rootCmd.Flags().StringP("server", "s", "127.0.0.1:8008", "Start as a proxy server. Format: ip:port, default: 127.0.0.1:8008")
	rootCmd.Flags().StringP("proxy", "p", "", "Use a proxy. Format: protocol://user:pass@ip:port")
	rootCmd.Flags().StringP("token_file", "f", "access_tokens.json", "Specify an access tokens json file.")
	rootCmd.Flags().BoolP("api", "a", false, "Use gpt-3.5-turbo chat api. Note: OpenAI will bill you.")
	rootCmd.Flags().BoolP("sentry", "e", false, "Enable sentry to send error reports when errors occur.")
	rootCmd.Flags().BoolP("verbose", "v", false, "Show exception traceback.")

	//initConfig()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
