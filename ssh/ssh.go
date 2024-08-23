package ssh

import (
	"fmt"
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"time"
)

// PublicKeyAuthFunc 公钥文件地址
func PublicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", zap.Error(err))
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", zap.Error(err))
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", zap.Error(err))
	}
	return ssh.PublicKeys(signer)
}

// PasswordFunc 密码
func PasswordFunc(pwd string) ssh.AuthMethod {
	return ssh.Password(pwd)
}

type ForwardInfo struct {
	Port int    //远程端口
	IP   string //远程内网IP
}

// Conn 开启 ssh连接
// forwardInfoMaps key为本地端口
func Conn(user, sshServer string, forwardInfoMaps map[int]*ForwardInfo, auth ...ssh.AuthMethod) {
	// SSH配置
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// 建立 SSH 连接
	sshClient, err := ssh.Dial("tcp", sshServer, sshConfig)
	if err != nil {
		log.Fatal("Failed to dial SSH server", zap.Error(err))
	}
	defer sshClient.Close()

	for port, info := range forwardInfoMaps {
		go func(localPort int, forwardInfo *ForwardInfo) {
			// 监听本地端口，将数据路由到 SSH 隧道
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
			if err != nil {
				log.Fatal("Failed to create listener", zap.Error(err))
			}
			defer listener.Close()
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Fatal("Failed to accept connection", zap.Error(err))
					os.Exit(1)
				}

				go func(conn net.Conn) {
					defer conn.Close()

					sshConn, err := sshClient.Dial("tcp", fmt.Sprintf("%s:%d", forwardInfo.IP, forwardInfo.Port))
					if err != nil {
						log.Fatal("Failed to establish SSH tunnel", zap.Error(err))
						return
					}
					defer sshConn.Close()

					go func() {
						_, _ = io.Copy(conn, sshConn)
					}()

					_, _ = io.Copy(sshConn, conn)
				}(conn)
			}

		}(port, info)
	}
	select {}
}
