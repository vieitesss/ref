package displayname

import "fmt"

type Props = string

func Renderer(text Props, width, height int) string {
	return fmt.Sprintf("OMG! Hello %s!", text)
}
