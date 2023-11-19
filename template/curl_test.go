package template

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCurlTemplate_compareApi(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock response"))
	}))
	defer mockServer.Close()

	tests := []struct {
		URL            *string
		expectedBody   string
		expectedStatus int
		expectedError  bool
		TestName       string
	}{
		{
			expectedBody:   "non",
			expectedStatus: 0,
			expectedError:  false,
			TestName:       "APIが指定されていなければ、エラーにならないこと",
		},
		{
			URL:            &mockServer.URL,
			expectedBody:   "Mock response",
			expectedStatus: 200,
			expectedError:  false,
			TestName:       "BodyとStatusが一致する時、エラーにならないこと",
		},
		{
			URL:            &mockServer.URL,
			expectedBody:   "Mock responses",
			expectedStatus: 200,
			expectedError:  true,
			TestName:       "Bodyが一致しない時、エラーになること",
		},
		{
			URL:            &mockServer.URL,
			expectedBody:   "Mock response",
			expectedStatus: 500,
			expectedError:  true,
			TestName:       "Statusが一致しない時、エラーになること",
		},
	}

	for _, tt := range tests {
		temp := CurlTemplate{
			URL:    mockServer.URL,
			Method: "GET",
			Expect: struct {
				Status *int    `yaml:"status"`
				Text   *string `yaml:"text"`
				Api    *string `yaml:"api"`
				File   *string `yaml:"file"`
			}(struct {
				Status *int
				Text   *string
				Api    *string
				File   *string
			}{Api: tt.URL}),
		}

		t.Run(tt.TestName, func(t *testing.T) {
			err := temp.compareApi([]byte(tt.expectedBody), tt.expectedStatus)

			if tt.expectedError && err == nil {
				t.Errorf("no error occurred")
			}
			if !tt.expectedError && err != nil {
				t.Error(err)
			}
		})
	}
}

func TestCurlTemplate_request(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock response"))
	}))
	defer mockServer.Close()

	// テストケース
	tests := []struct {
		name          string
		template      CurlTemplate
		expectedCode  int
		expectedBody  string
		expectedError bool
	}{
		{
			name: "Valid Response",
			template: CurlTemplate{
				URL:    mockServer.URL,
				Method: "GET",
			},
			expectedCode:  http.StatusOK,
			expectedBody:  "Mock response",
			expectedError: false,
		},
	}

	// テストを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, err := tt.template.request()

			if tt.expectedError && err == nil {
				t.Errorf("expected an error, but got nil")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectedCode != 0 && status != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, status)
			}

			if tt.expectedBody != "" && body != nil && string(body) != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
			}
		})
	}
}

func TestCurlTemplate_compareFile(t *testing.T) {
	// 仮のファイルとデータを作成
	path := "example.txt"
	dir := "."
	tmpFile, err := os.CreateTemp(dir, path)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // テストが終わったらファイルを削除

	fileData := []byte("Test file data")
	if _, err := tmpFile.Write(fileData); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	targetFilePath := dir + "/" + tmpFile.Name()
	tests := []struct {
		file          *string
		testName      string
		expectedBody  string
		expectedError bool
	}{
		{
			file:          &targetFilePath,
			testName:      "ファイルの中身が一致する場合はエラーにならない",
			expectedBody:  "Test file data",
			expectedError: false,
		},
		{
			file:          &targetFilePath,
			testName:      "ファイルの中身が一致しない場合はエラーになる",
			expectedBody:  "Test file dats",
			expectedError: true,
		},
		{
			testName:      "file未指定の場合はエラーにならない",
			expectedBody:  "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tmp := CurlTemplate{
				Expect: struct {
					Status *int    `yaml:"status"`
					Text   *string `yaml:"text"`
					Api    *string `yaml:"api"`
					File   *string `yaml:"file"`
				}(struct {
					Status *int
					Text   *string
					Api    *string
					File   *string
				}{File: tt.file}),
			}

			err := tmp.compareFile([]byte(tt.expectedBody))

			if tt.expectedError && err == nil {
				t.Errorf("no error occurred")
			}
			if !tt.expectedError && err != nil {
				t.Error(err)
			}
		})
	}
}
