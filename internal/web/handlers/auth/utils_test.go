package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peltoche/zapette/internal/service/users"
	"github.com/Peltoche/zapette/internal/service/websessions"
	"github.com/Peltoche/zapette/internal/web/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Utils_Authenticator(t *testing.T) {
	t.Run("getUserAndSession success", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(&websessions.AliceWebSessionExample, nil).Once()
		usersMock.On("GetByID", mock.Anything, users.ExampleAlice.ID()).Return(&users.ExampleAlice, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, abort := auth.GetUserAndSession(w, r, AnyUser)
		assert.Equal(t, &users.ExampleAlice, user)
		assert.Equal(t, &websessions.AliceWebSessionExample, session)
		assert.False(t, abort)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("getUserAndSession with a websession error", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(nil, errors.New("some-error")).Once()

		htmlMock.On("WriteHTMLErrorPage", mock.Anything, mock.Anything, mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, abort := auth.GetUserAndSession(w, r, AnyUser)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.True(t, abort)
	})

	t.Run("getUserAndSession with a websession not found", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(nil, websessions.ErrMissingSessionToken).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, abort := auth.GetUserAndSession(w, r, AnyUser)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.True(t, abort)

		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusFound, res.StatusCode)
		assert.Equal(t, "/web/login", w.Header().Get("Location"))
	})

	t.Run("getUserAndSession with a users problem", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(&websessions.AliceWebSessionExample, nil).Once()
		usersMock.On("GetByID", mock.Anything, users.ExampleAlice.ID()).Return(nil, errors.New("some-error")).Once()

		htmlMock.On("WriteHTMLErrorPage", mock.Anything, mock.Anything, mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, abort := auth.GetUserAndSession(w, r, AnyUser)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.True(t, abort)
	})

	t.Run("getUserAndSession with a user not found", func(t *testing.T) {
		webSessionsMock := websessions.NewMockService(t)
		usersMock := users.NewMockService(t)
		htmlMock := html.NewMock(t)
		auth := NewAuthenticator(webSessionsMock, usersMock, htmlMock)

		webSessionsMock.On("GetFromReq", mock.Anything, mock.Anything).Return(&websessions.AliceWebSessionExample, nil).Once()
		usersMock.On("GetByID", mock.Anything, users.ExampleAlice.ID()).Return(nil, nil).Once()

		webSessionsMock.On("Logout", mock.Anything, mock.Anything).Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		user, session, abort := auth.GetUserAndSession(w, r, AnyUser)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.True(t, abort)
	})
}
