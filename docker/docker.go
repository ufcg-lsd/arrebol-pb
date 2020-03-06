package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type ContainerConfig struct {
	Name string
	Image string
	Mounts []mount.Mount
}

func NewDockerClient(host string) *client.Client {
	if err := os.Setenv("DOCKER_HOST", host); err != nil {
		log.Print(err)
	}
	cli, _ := client.NewEnvClient()
	return cli
}

func ListContainer(cli *client.Client) {
	ctns, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All:true})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range ctns {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func CreateContainer(cli *client.Client, config ContainerConfig) (string, error) {
	log.Printf("Creating Container [%s]", config.Name)
	ctx := context.Background()
	hostConfig := container.HostConfig{
		Mounts: config.Mounts,
	}

	dconfig := container.Config{
		Image:   config.Image,
		Tty:     true,
	}

	b, err := cli.ContainerCreate(ctx, &dconfig, &hostConfig, nil, config.Name)

	return b.ID, err
}

func StartContainer(cli *client.Client, id string) error {
	log.Printf("Starting Container [%s]", id)
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()
	return cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func StopContainer(cli *client.Client, id string) error {
	log.Printf("Stopping Container [%s]", id)
	var timeout = 5 * time.Second
	return cli.ContainerStop(context.Background(), id, &timeout)
}

func RemoveContainer(cli *client.Client, id string) error {
	log.Printf("Removing Container [%s]", id)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
}

func Write(cli *client.Client, id string, content []string, dest string) error {
	for _, c := range content {
		c = strings.ReplaceAll(c, "'", "'\"'\"'")
		//c = AddSlashes(c)
		cmd := fmt.Sprintf(`echo -E '%s' >> %s`, c, dest)
		log.Printf("Writing [%s] on [%s] from Container [%s]", c, dest, id)
		Exec(cli, id, cmd)
	}
	return nil
}

func Copy(cli *client.Client, id, src, dest string) error {
	log.Printf("Copy [%s] to [%s] from Container [%s]", src, dest, id)
	dat, _ := ioutil.ReadFile(src)
	content := string(dat)
	content = strings.ReplaceAll(content, "'", "'\"'\"'")
	cmd := fmt.Sprintf("echo -E '%s' >| %s", content, dest)
	return Exec(cli, id, cmd)
}

func Exec(cli *client.Client, id, cmd string)  error {
	log.Printf("Executing command [%s] on container [%s]", cmd, id)
	config := types.ExecConfig{
		Cmd: []string{"/bin/bash", "-c", cmd},
	}
	rid, _ := cli.ContainerExecCreate(context.Background(), id, config)
	return cli.ContainerExecStart(context.Background(), rid.ID, types.ExecStartCheck{})
}

func Cat(cli *client.Client, id, path string) ([]byte, error) {
	log.Printf("Getting content of file [%s]", path)
	config := types.ExecConfig{
		Tty:true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:[]string{"/bin/bash", "-c", "cat " + path},
	}
	rid, _ := cli.ContainerExecCreate(context.Background(), id, config)
	hijack, _ := cli.ContainerExecAttach(context.Background(), rid.ID, types.ExecConfig{Tty: true,})
	output := read(hijack.Conn)
	return output, nil
}

func read(conn net.Conn) []byte {
	result := make([]byte, 0)
	b := make([]byte, 10)
	for ; ; {
		n, _ := conn.Read(b)
		result = append(result, b...)
		if n < len(b) {break}
		b = make([]byte, 2 * len(b))
	}
	return result
}

func Pull(cli *client.Client, image string) (io.ReadCloser, error) {
	reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	return reader, err
}

func CheckImage(cli *client.Client, image string) (exist bool, err error)  {
	exist = false
	_, _, err = cli.ImageInspectWithRaw(context.Background(), image)
	if err == nil {
		exist = true
	}
	return
}