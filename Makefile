coverage:
	@go test -cover .
coverage_html:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out