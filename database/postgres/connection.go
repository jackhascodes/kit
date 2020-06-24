// package postgres provides methods for connecting to a postgres database either raw or via ssh.
// Cribbed from vinzenz/dial-pq-via-ssh.go
package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func InitConnection(opts ...func(*Config)) (*sql.DB, error) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	if c.useSSH {
		return ConnectDBViaSSH(c)
	}
	return ConnectDB(c)
}

type Config struct {
	dbHost, dbUser, dbPass, dbName, sshHost, sshPort, sshUser, sshPass string
	useSSH                                                             bool
}

func WithDatabase(dbHost, dbUser, dbPass, dbName string) func(*Config) {
	return func(c *Config) {
		c.dbHost = dbHost
		c.dbUser = dbUser
		c.dbPass = dbPass
		c.dbName = dbName
	}
}

func WithSSH(sshHost, sshPort, sshUser, sshPass string) func(config *Config) {
	return func(c *Config) {
		c.useSSH = true
		c.sshHost = sshHost
		c.sshPass = sshPass
		c.sshUser = sshUser
		c.sshPort = sshPort
	}
}

func ConnectDB(c *Config) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.dbUser, c.dbPass, c.dbHost, c.dbName))
}

func ConnectDBViaSSH(c *Config) (*sql.DB, error) {
	var agentClient agent.Agent
	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close()

		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: c.sshUser,
		Auth: []ssh.AuthMethod{},
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}
	// When there's a non empty password add the password AuthMethod
	if c.sshPass != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return c.sshPass, nil
		}))
	}
	sshConfig.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	// Connect to the SSH Server
	sshp, _ := strconv.ParseInt(c.sshPort, 10, 8)
	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.sshHost, sshp), sshConfig)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Now we register the ViaSSHDialer with the ssh connection as a parameter
	sql.Register("postgres+ssh", &ViaSSHDialer{sshcon})

	// And now we can use our new driver with the regular postgres connection string tunneled through the SSH connection
	return sql.Open("postgres+ssh", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.dbUser, c.dbPass, c.dbHost, c.dbName))
}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func (self *ViaSSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.client.Dial(network, address)
}
