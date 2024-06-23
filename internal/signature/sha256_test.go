package signature

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sha256Sig_Generate(t *testing.T) {

	t.Run("1", func(t *testing.T) {
		sigSvc := NewSha256Sig("secret auth key", []byte{123})
		sig, err := sigSvc.Generate()
		require.NoError(t, err)
		require.Equal(t, "bb409314cf250f4c447cbd10e3611b189b2af6ce8aa62ca68a60917fadc8eb5e", sig)

	})
}
