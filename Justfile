## Global Variables
# 15minute timeout for go (terratest) tests
timeout_default := "15m" 


# list out commands
default:
  @just --list


# Run all test with configurable timeout (default 15mins)
test_all timeout=timeout_default:
  go test -timeout "{{timeout}}" -v

# Test a specifc module with configurable timeout (default 15mins)
test module timeout=timeout_default:
  cd "test/{{module}}"
  go test -timeout "{{timeout}}" -v
  