# Example of how to add summary jobs (optional)

# Add this to unit-tests.yml after the existing unit-tests job:

unit-tests-summary:
name: Unit Tests
needs: unit-tests
runs-on: ubuntu-latest
if: always()
steps: - name: Check unit tests result
run: |
if [ "${{ needs.unit-tests.result }}" != "success" ]; then
echo "Unit tests failed"
exit 1
fi

# Add this to integration-tests.yml after the existing integration-tests job:

integration-tests-summary:
name: Integration Tests
needs: integration-tests
runs-on: ubuntu-latest
if: always()
steps: - name: Check integration tests result
run: |
if [ "${{ needs.integration-tests.result }}" != "success" ]; then
echo "Integration tests failed"
exit 1
fi

# Then your required status checks would be:

# 1. Unit Tests (summary job)

# 2. Integration Tests (summary job)

# 3. CodeQL Security Analysis
