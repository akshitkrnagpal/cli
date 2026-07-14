package debug

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/supabase/cli/pkg/pgtest"
)

func TestPostgresProxy(t *testing.T) {
	const postgresUrl = "postgresql://postgres:password@127.0.0.1:5432/postgres?sslmode=disable"

	t.Run("forwards messages between frontend and backend", func(t *testing.T) {
		// Parse connection url
		config, err := pgx.ParseConfig(postgresUrl)
		require.NoError(t, err)
		// Setup postgres mock
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Intercept(config)
		// Run test
		SetupPGX(config)
		ctx := context.Background()
		proxy, err := pgx.ConnectConfig(ctx, config)
		assert.NoError(t, err)
		assert.NoError(t, proxy.Close(ctx))
	})

	t.Run("preserves tls configuration", func(t *testing.T) {
		config, err := pgx.ParseConfig("postgresql://postgres:password@db.example.com:5432/postgres?sslmode=require")
		require.NoError(t, err)
		require.NotNil(t, config.TLSConfig)

		SetupPGX(config)

		assert.NotNil(t, config.TLSConfig)
	})

	t.Run("preserves mixed tls fallbacks", func(t *testing.T) {
		config, err := pgx.ParseConfig("postgresql://postgres:password@db.example.com:5432/postgres?sslmode=allow")
		require.NoError(t, err)
		require.Nil(t, config.TLSConfig)
		require.True(t, hasTLSFallback(config))

		SetupPGX(config)

		assert.True(t, hasTLSFallback(config))
	})
}
