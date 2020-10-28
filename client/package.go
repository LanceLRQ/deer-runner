package client

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/LanceLRQ/deer-executor/persistence"
	"github.com/LanceLRQ/deer-executor/persistence/problems"
	"github.com/urfave/cli/v2"
	"log"
)

var PackProblemFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     	"sign",
		Aliases:  	[]string{"s"},
		Value: 	 	false,
		Usage:    	"Enable digital sign (GPG)",
	},
	&cli.StringFlag {
		Name: 		"gpg-key",
		Aliases:  	[]string{"key"},
		Value: 		"",
		Usage: 		"GPG private key file",
	},
	&cli.StringFlag {
		Name: 		"passphrase",
		Aliases:  	[]string{"p", "password", "pwd"},
		Value: 		"",
		Usage: 		"GPG private key passphrase",
	},
}

func PackProblem(c *cli.Context) error {

	if c.String("passphrase") != "" {
		log.Println("[warn] Using a password on the command line interface can be insecure.")
	}
	passphrase := []byte(c.String("passphrase"))
	configFile := c.Args().Get(0)
	outputFile := c.Args().Get(1)

	var options persistence.CommonPersisOptions
	var err error
	var pem *persistence.DigitalSignPEM

	if c.Bool("sign") {
		pem, err = persistence.GetArmorPublicKey(c.String("gpg-key"), passphrase)
		if err != nil {
			return err
		}
	}

	options = persistence.CommonPersisOptions{
		DigitalSign: c.Bool("sign"),
		DigitalPEM:  pem,
		OutFile:     outputFile,
	}

	// problem
	session, err := executor.NewSession(configFile)
	if err != nil {
		return err
	}
	err = problems.PackProblems(session, options)
	if err != nil {
		return err
	}

	return nil
}

func ReadProblemInfo(c *cli.Context) error {

	configFile := c.Args().Get(0)
	yes, err := problems.IsProblemPackage(configFile)
	if err != nil {
		return err
	}
	// 如果是题目包文件，进行解包
	if yes {
		if c.Bool("sign") {
			g, err := problems.ReadProblemGPGInfo(configFile)
			if err != nil {
				return err
			}
			fmt.Println(g)
		} else {
			s, err := problems.ReadProblemInfo(configFile, false, "")
			if err != nil {
				return err
			}
			s.ConfigFile = ""
			fmt.Println(executor.ObjectToJSONStringFormatted(s))
		}
	} else {
		return fmt.Errorf("not deer-executor problem package file")
	}

	return nil
}
