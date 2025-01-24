package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_server(t *testing.T) {
	ctx := context.Background()
	e := createEcho()

	go func() {
		e.Start(":1323")
	}()

	cli := http.DefaultClient

	t.Run("OK/createUser", func(t *testing.T) {
		b, err := json.Marshal(User{
			Name:  "Smith",
			Email: "smith@test.com",
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:1323/users", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "enjoy")
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		var got User
		require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
		require.Equal(t, User{
			ID:    got.ID,
			Name:  "Smith",
			Email: "smith@test.com",
		}, got)
	})

	t.Run("OK/getUser", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:1323/users/1", nil)
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		var got User
		require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
		require.Equal(t,
			User{
				ID:    "1",
				Name:  "Smith",
				Email: "smith@test.com",
			}, got)
	})
	t.Run("OK/updateUser", func(t *testing.T) {
		b, err := json.Marshal(User{
			Name:  "Smith",
			Email: "smith@test.com",
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "http://localhost:1323/users/1", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "enjoy")
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		var got User
		require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
		require.Equal(t, User{
			ID:    "1",
			Name:  "Smith",
			Email: "smith@test.com",
		}, got)
	})
	t.Run("OK/deleteUser", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:1323/users/1", nil)
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)

		require.Empty(t, res.Body)
	})
	t.Run("NG/400", func(t *testing.T) {
		b, err := json.Marshal(User{
			Name:  "Smith",
			Email: "smith@test.com",
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:1323/users", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("NG/401", func(t *testing.T) {
		b, err := json.Marshal(User{
			Name:  "Smith",
			Email: "smith@test.com",
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:1323/users", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "invalid")
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
	t.Run("NG/404", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:1323/us", nil)
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
	t.Run("NG/415", func(t *testing.T) {
		b, err := json.Marshal(User{
			Name:  "Smith",
			Email: "smith@test.com",
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:1323/users", bytes.NewReader(b))
		req.Header.Set("x-api-key", "enjoy")
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnsupportedMediaType, res.StatusCode)

		type msg struct {
			Message string
		}
		var got msg
		require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
		require.Equal(t, msg{
			Message: "Unsupported Media Type",
		}, got)
	})

	e.Shutdown(ctx)
}

func Test_internalError(t *testing.T) {
	ctx := context.Background()
	e := createEcho()

	go func() {
		e.Start(":1323")
	}()

	cli := http.DefaultClient

	t.Run("NG/updateUser", func(t *testing.T) {
		b, err := json.Marshal(User{
			Name:  "Smith",
			Email: "smith@test.com",
		})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "http://localhost:1323/users/9999", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "enjoy")
		require.NoError(t, err)

		res, err := cli.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, res.StatusCode)

		type msg struct {
			Message string
		}
		var got msg
		require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
		require.Equal(t, msg{
			Message: "Internal Server Error",
		}, got)
	})

	e.Shutdown(ctx)
}
