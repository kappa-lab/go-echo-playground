package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TeeReader(t *testing.T) {
	dat := []byte("abcde")

	var datReader io.Reader = bytes.NewReader(dat)
	var dump bytes.Buffer
	tr := io.TeeReader(datReader, &dump)

	fmt.Println("dump:0", dump.String()) // 空

	r1, err := io.ReadAll(tr)
	require.NoError(t, err)

	fmt.Println("r1:", string(r1))       // abcde
	fmt.Println("dump:1", dump.String()) // abcde
}

func Test_HttpRequest(t *testing.T) {
	dat := []byte("{'name':'Smith'}")

	t.Run("これだと空", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(dat))
		require.NoError(t, err)

		err = request.Body.Close()
		require.NoError(t, err)

		b, err := io.ReadAll(request.Body)
		require.NoError(t, err)
		fmt.Println("b0:", string(b)) // {'name':'Smith'}

		b, err = io.ReadAll(request.Body)
		require.NoError(t, err)
		fmt.Println("b1:", string(b)) // 空
	})

	t.Run("TeeReaderに食わせると取得できる", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(dat))
		require.NoError(t, err)

		var dump bytes.Buffer
		tr := io.TeeReader(request.Body, &dump)

		// If body is of type [*bytes.Buffer], [*bytes.Reader], or
		// [*strings.Reader], the returned request's ContentLength is set to its
		// exact value (instead of -1), GetBody is populated (so 307 and 308
		// redirects can replay the body), and Body is set to [NoBody] if the
		// ContentLength is 0.
		request.Body = io.NopCloser(tr)

		b, err := io.ReadAll(request.Body)
		require.NoError(t, err)
		fmt.Println("b0:", string(b)) // {'name':'Smith'}

		b, err = io.ReadAll(request.Body)
		require.NoError(t, err)
		fmt.Println("b1:", string(b)) // 空

		fmt.Println("dump", dump.String()) // {'name':'Smith'}
	})

}
