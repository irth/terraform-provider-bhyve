package client

import (
	"errors"
	"fmt"
	"strings"
)

type ErrChecksumMismatch struct {
	Path     string
	Expected string
	Actual   string

	wrap error
}

func (e ErrChecksumMismatch) Unwrap() error {
	return e.wrap
}

func (e ErrChecksumMismatch) Error() string {
	return fmt.Sprintf("checksum mismatch: %s: expected %s, got %s", e.Path, e.Expected, e.Actual)
}

type ErrFileNotFound struct {
	Path string

	wrap error
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("resource not found: %s", e.Path)
}

func (e ErrFileNotFound) Unwrap() error {
	return e.wrap
}

func (c *Client) verifyHash(path string, checksum string) (string, error) {
	cmd := []string{}
	if checksum != "" {
		cmd = append(cmd, "-c", checksum)
	}
	cmd = append(cmd, "-q", path)

	stdout, err := c.Executor.RunCmd("sha256", cmd...)
	stdout = strings.TrimSpace(stdout)

	if err != nil {
		var cmdExecErr CommandExecutionError
		if errors.As(err, &cmdExecErr) {
			if strings.Contains(cmdExecErr.Stderr, "No such file or directory") {
				return "", ErrFileNotFound{
					Path: path,

					wrap: err,
				}
			}

			if len(stdout) == 64 {
				return stdout, ErrChecksumMismatch{
					Path:     path,
					Expected: checksum,
					Actual:   stdout,

					wrap: err,
				}
			}
		}
		return "", err
	}

	return stdout, nil
}

func (c *Client) removeFile(path string) error {
	_, err := c.Executor.RunCmd("rm", "-f", path)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %s", path, err)
	}

	return nil
}

func (c *Client) download(url string, path string, checksum string) error {
	_, err := c.verifyHash(path, checksum)

	switch {
	case err == nil:
		return nil

	case errors.As(err, &ErrFileNotFound{}):
		// do nothing

	case errors.As(err, &ErrChecksumMismatch{}):
		c.removeFile(path)

	default:
		return err
	}

	_, err = c.Executor.RunCmd("fetch", "-a", "-o", path, url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}

	_, err = c.verifyHash(path, checksum)
	if err != nil {
		c.removeFile(path)
		return err
	}

	return nil
}

func (c *Client) ISO(url string, path string, checksum string) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	isoDir := conf.ISOPath()
	_, err = c.Executor.RunCmd("mkdir", "-p", isoDir)
	if err != nil {
		return err
	}

	isoPath := conf.ISOPath(path)

	return c.download(url, isoPath, checksum)
}

func (c *Client) IMG(url string, path string, checksum string) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	isoDir := conf.IMGPath()
	_, err = c.Executor.RunCmd("mkdir", "-p", isoDir)
	if err != nil {
		return err
	}

	imgPath := conf.IMGPath(path)

	return c.download(url, imgPath, checksum)
}

func (c *Client) RemoveISO(path string) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	isoPath := conf.ISOPath(path)
	return c.removeFile(isoPath)
}

func (c *Client) RemoveIMG(path string) error {
	conf, err := c.Config()
	if err != nil {
		return err
	}

	imgPath := conf.IMGPath(path)
	return c.removeFile(imgPath)
}

func (c *Client) ChecksumISO(path string) (string, error) {
	conf, err := c.Config()
	if err != nil {
		return "", err
	}

	isoPath := conf.ISOPath(path)
	return c.verifyHash(isoPath, "")
}

func (c *Client) ChecksumIMG(path string) (string, error) {
	conf, err := c.Config()
	if err != nil {
		return "", err
	}

	imgPath := conf.IMGPath(path)
	return c.verifyHash(imgPath, "")
}
