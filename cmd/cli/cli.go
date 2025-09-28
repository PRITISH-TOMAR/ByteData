package cli

import (
	"bufio"
	"fmt"
	"os"
	"github.com/PRITISH-TOMAR/byted/internal/kv"
)
// StartCLI starts the ByteData shell for a given KVEngine
func StartCLI(username, password string, engine *kv.KVEngine) {

	fmt.Printf("Welcome %s! Connected to ByteData Engine.\n", username)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Bytedata> ")
		if !reader.Scan() {
			break
		}
		cmdLine := reader.Text()
		err := ExecuteCommmand(cmdLine, engine)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
	fmt.Println("\nBye!")
}