package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	DICDIR  = "./.jpug-doc-tool/"
)

type apiConfig struct {
	ClientID             string
	ClientSecret         string
	Name                 string
	APIAutoTranslate     string
	APIAutoTranslateType string
}

var Config apiConfig

func ignoreFileNames(fileNames []string) []string {
	var ignoreFile map[string]struct{} = map[string]struct{}{
		"jpug-doc.sgml":  {},
		"config0.sgml":   {},
		"config1.sgml":   {},
		"config2.sgml":   {},
		"config3.sgml":   {},
		"func0.sgml":     {},
		"func1.sgml":     {},
		"func2.sgml":     {},
		"func3.sgml":     {},
		"func4.sgml":     {},
		"catalogs0.sgml": {},
		"catalogs1.sgml": {},
		"catalogs2.sgml": {},
		"catalogs3.sgml": {},
		"catalogs4.sgml": {},
	}

	ret := make([]string, 0, len(fileNames))
	for _, fileName := range fileNames {
		if _, ok := ignoreFile[fileName]; ok {
			continue
		}
		ret = append(ret, fileName)
	}
	return ret
}

func targetFileName() []string {
	pattern := "./*.sgml"
	rePattern := "./*/*.sgml"

	fileNames, err := filepath.Glob(pattern)
	if err != nil {
		log.Println(err)
		return nil
	}
	reFileNames, err := filepath.Glob(rePattern)
	if err != nil {
		log.Println(err)
		return nil
	}
	fileNames = append(fileNames, reFileNames...)
	fileNames = ignoreFileNames(fileNames)
	return fileNames
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jpug-doc-tool",
	Short: "jpug-doc tool",
	Long: `
jpug-doc の翻訳を補助ツール。
前バージョンの翻訳を新しいバージョンに適用したり、
翻訳のチェックが可能です。`,
}

func ReadFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initJpug)
	_ = rootCmd.RegisterFlagCompletionFunc("", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jpug-doc-tool.yaml)")
	_ = rootCmd.RegisterFlagCompletionFunc("config", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml"}, cobra.ShellCompDirectiveFilterFileExt
	})

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".jpug-doc-tool")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Println("config file Unmarshal error")
		fmt.Println(err)
		os.Exit(1)
	}
}

func initJpug() {
	f, err := filepath.Glob("./*.sgml")
	if err != nil || len(f) == 0 {
		fmt.Fprintln(os.Stderr, "*sgmlファイルがあるディレクトリで実行してください")
		fmt.Fprintln(os.Stderr, "cd github.com/pgsql-jp/jpug-doc/doc/src/sgml")
		return
	}
	if _, err := os.Stat(DICDIR); os.IsNotExist(err) {
		os.Mkdir(DICDIR, 0o755)
	}
	refdir := DICDIR + "/ref"
	if _, err := os.Stat(refdir); os.IsNotExist(err) {
		os.Mkdir(refdir, 0o755)
	}
}
