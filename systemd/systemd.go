package systemd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const tpl = `[Unit]
After=network-online.target

[Service]
LimitNOFILE=65535
Type=simple
ExecStart={{CMD}}
ExecReload=/bin/kill -s HUP $MAINPID
Restart=always

[Install]
WantedBy=multi-user.target
`

var Dir = "/lib/systemd/system/"

func AutoStart(name string) error {
	for _, arg := range os.Args {
		if arg == "--help" {
			return nil
		}
	}
	if os.Getenv("__SYSTEMD_IGNORED__") != "" {
		return nil
	}
	ctl, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(Dir, name+".service")); !os.IsNotExist(err) {
		return nil
	}
	cmd, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}
	args := append([]string{cmd}, os.Args[1:]...)
	if err = os.WriteFile(filepath.Join(Dir, name+".service"), []byte(strings.Replace(tpl, "{{CMD}}", strings.Join(args, " "), -1)), os.FileMode(644)); err != nil {
		return err
	}
	if err = exec.Command(ctl, "enable", name).Run(); err != nil {
		return err
	}
	if err = exec.Command(ctl, "start", name).Run(); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
